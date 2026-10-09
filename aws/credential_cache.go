package aws

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/credentials/nonrecoverable"
	sdkrand "github.com/aws/aws-sdk-go-v2/internal/rand"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/internal/sync/singleflight"
	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/middleware"
)

// CredentialsCacheOptions are the options
type CredentialsCacheOptions struct {

	// ExpiryWindow will allow the credentials to trigger refreshing prior to
	// the credentials actually expiring. This is beneficial so race conditions
	// with expiring credentials do not cause request to fail unexpectedly
	// due to ExpiredTokenException exceptions.
	//
	// An ExpiryWindow of 10s would cause calls to IsExpired() to return true
	// 10 seconds before the credentials are actually expired. This can cause an
	// increased number of requests to refresh the credentials to occur.
	//
	// If ExpiryWindow is 0 or less it will be ignored.
	ExpiryWindow time.Duration

	// ExpiryWindowJitterFrac provides a mechanism for randomizing the
	// expiration of credentials within the configured ExpiryWindow by a random
	// percentage. Valid values are between 0.0 and 1.0.
	//
	// As an example if ExpiryWindow is 60 seconds and ExpiryWindowJitterFrac
	// is 0.5 then credentials will be set to expire between 30 to 60 seconds
	// prior to their actual expiration time.
	//
	// If ExpiryWindow is 0 or less then ExpiryWindowJitterFrac is ignored.
	// If ExpiryWindowJitterFrac is 0 then no randomization will be applied to the window.
	// If ExpiryWindowJitterFrac < 0 the value will be treated as 0.
	// If ExpiryWindowJitterFrac > 1 the value will be treated as 1.
	ExpiryWindowJitterFrac float64
}

// CredentialsCache provides caching and concurrency safe credentials retrieval
// via the provider's retrieve method.
//
// CredentialsCache will look for optional interfaces on the Provider to adjust
// how the credential cache handles credentials caching.
//
//   - HandleFailRefreshCredentialsCacheStrategy - Allows provider to handle
//     credential refresh failures. This could return an updated Credentials
//     value, or attempt another means of retrieving credentials.
//
//   - AdjustExpiresByCredentialsCacheStrategy - Allows provider to adjust how
//     credentials Expires is modified. This could modify how the Credentials
//     Expires is adjusted based on the CredentialsCache ExpiryWindow option.
//     Such as providing a floor not to reduce the Expires below.
type CredentialsCache struct {
	provider CredentialsProvider

	options CredentialsCacheOptions
	creds   atomic.Pointer[cachedCredentials]
	sf      singleflight.Group

	cachedErr  atomic.Pointer[cachedError]
	refreshing atomic.Bool
	legacy     bool
	rand       func() (float64, error)
}

type cachedCredentials struct {
	creds                Credentials
	advisoryWindow       time.Duration
	nextRefreshAllowedAt time.Time
}

type cachedError struct {
	err       error
	expiresAt time.Time
}

// NewCredentialsCache returns a CredentialsCache that wraps provider. Provider
// is expected to not be nil. A variadic list of one or more functions can be
// provided to modify the CredentialsCache configuration. This allows for
// configuration of credential expiry window and jitter.
func NewCredentialsCache(provider CredentialsProvider, optFns ...func(options *CredentialsCacheOptions)) *CredentialsCache {
	options := CredentialsCacheOptions{}

	for _, fn := range optFns {
		fn(&options)
	}

	if options.ExpiryWindow < 0 {
		options.ExpiryWindow = 0
	}

	if options.ExpiryWindowJitterFrac < 0 {
		options.ExpiryWindowJitterFrac = 0
	} else if options.ExpiryWindowJitterFrac > 1 {
		options.ExpiryWindowJitterFrac = 1
	}

	return &CredentialsCache{
		provider: provider,
		options:  options,
		legacy:   !inScope(provider),
		rand:     sdkrand.CryptoRandFloat64,
	}
}

// picks an legacy, initial, advisory, or mandatory refresh based on where the cached credentials are in their lifetime.
// Or dont try to refresh for nonRecoverableErr
// Retrieve returns the credentials. If the credentials have already been
// retrieved, and not expired the cached credentials will be returned. If the
// credentials have not been retrieved yet, or expired the provider's Retrieve
// method will be called.
//
// Returns and error if the provider's retrieve method returns an error.
func (p *CredentialsCache) Retrieve(ctx context.Context) (Credentials, error) {
	if p.legacy {
		return p.legacyRetrieve(ctx)
	}

	now := sdk.NowTime().Round(0)

	e := p.creds.Load()
	if e == nil || !e.creds.HasKeys() {
		return p.waitForRetrieve(ctx, retrieveInitial)
	}
	if !e.refreshNeeded(now) {
		return e.creds, nil
	}
	if e.rateLimited(now) {
		return e.creds, nil
	}

	if err, ok := p.nonRecoverableErrCached(now); ok {
		return Credentials{}, err
	}
	if !e.mandatoryRefreshNeeded(now) {
		return p.advisoryRefresh(ctx, e)
	}
	return p.waitForRetrieve(ctx, refreshMandatory)
}

type refreshType int

const (
	retrieveInitial refreshType = iota
	refreshAdvisory
	refreshMandatory
)

// advisoryRefresh starts a singleRetrieve if none is in flight; callers that find one in flight don't wait and get the cached credentials.
func (p *CredentialsCache) advisoryRefresh(ctx context.Context, currCreds *cachedCredentials) (Credentials, error) {
	if !p.refreshing.CompareAndSwap(false, true) {
		return currCreds.creds, nil
	}

	ch := p.sf.DoChan("", func() (interface{}, error) {
		return p.singleRetrieve(&suppressedContext{ctx}, refreshAdvisory)
	})
	select {
	case res := <-ch:
		if res.Err != nil {
			return Credentials{}, res.Err
		}
		return res.Val.(Credentials), nil
	case <-ctx.Done():
		return currCreds.creds, nil
	}
}

// waitForRetrieve starts a retrieve if none is in flight and waits for its result; used for the initial and mandatory refresh.
func (p *CredentialsCache) waitForRetrieve(ctx context.Context, typ refreshType) (Credentials, error) {
	p.refreshing.Store(true)
	ch := p.sf.DoChan("", func() (interface{}, error) {
		return p.singleRetrieve(&suppressedContext{ctx}, typ)
	})
	select {
	case res := <-ch:
		if res.Err != nil {
			return Credentials{}, res.Err
		}
		return res.Val.(Credentials), nil
	case <-ctx.Done():
		return Credentials{}, &RequestCanceledError{Err: ctx.Err()}
	}
}

// singleRetrieve asks the provider for new credentials and caches the result. If the provider fails, it keeps
// the old credentials and waits before trying again, or raises the error right away if retrying won't help.
func (p *CredentialsCache) singleRetrieve(ctx context.Context, typ refreshType) (Credentials, error) {
	defer p.refreshing.Store(false)

	now := sdk.NowTime().Round(0)
	currCreds := p.creds.Load()
	if currCreds == nil || !currCreds.creds.HasKeys() {
		typ = retrieveInitial
	}

	switch typ {
	case retrieveInitial:
		if currCreds != nil && currCreds.creds.HasKeys() {
			return currCreds.creds, nil
		}
	case refreshAdvisory:
		if currCreds != nil && !currCreds.refreshNeeded(now) {
			return currCreds.creds, nil
		}
	case refreshMandatory:
		if currCreds != nil && !currCreds.mandatoryRefreshNeeded(now) {
			return currCreds.creds, nil
		}
	}
	if err, ok := p.nonRecoverableErrCached(now); ok {
		return Credentials{}, err
	}
	if currCreds != nil && currCreds.rateLimited(now) {
		return currCreds.creds, nil
	}

	newCreds, err := p.provider.Retrieve(ctx)

	if err == nil && newCreds.CanExpire && !newCreds.Expires.After(now) {
		err = errors.New("credential source returned already-expired credentials")
	}

	if err != nil {
		var nr *nonrecoverable.Failure
		if errors.As(err, &nr) {
			p.cachedErr.Store(&cachedError{err: err, expiresAt: now.Add(p.jitter(1, 5))})
			return Credentials{}, err
		}

		if typ == retrieveInitial {
			return Credentials{}, fmt.Errorf("failed to refresh cached credentials, %w", err)
		}

		backoff := p.jitter(300, 600)
		p.creds.Store(currCreds.withBackoff(now.Add(backoff)))
		logFailedRefresh(ctx, err, backoff)
		return currCreds.creds, nil
	}

	p.creds.Store(&cachedCredentials{
		creds:          newCreds,
		advisoryWindow: p.advisoryWindowFor(now, newCreds),
	})
	p.cachedErr.Store(nil)
	return newCreds, nil
}

// logFailedRefresh logs a warning that a refresh failed and the cached credentials will be used until the backoff ends.
func logFailedRefresh(ctx context.Context, err error, backoff time.Duration) {
	logger := middleware.GetLogger(ctx)
	logger.Logf(logging.Warn,
		"Credential refresh failed: %v. The SDK will continue using cached credentials. "+
			"A refresh of these credentials will be attempted again after %d seconds.",
		err, int(backoff/time.Second))
}

// legacyRetrieve keeps the old caching behavior for out of scope providers, like process and custom providers.
func (p *CredentialsCache) legacyRetrieve(ctx context.Context) (Credentials, error) {
	if e := p.creds.Load(); e != nil && e.creds.HasKeys() && !e.creds.Expired() {
		return e.creds, nil
	}

	resCh := p.sf.DoChan("", func() (interface{}, error) {
		return p.legacySingleRetrieve(&suppressedContext{ctx})
	})
	select {
	case res := <-resCh:
		return res.Val.(Credentials), res.Err
	case <-ctx.Done():
		return Credentials{}, &RequestCanceledError{Err: ctx.Err()}
	}
}

func (p *CredentialsCache) legacySingleRetrieve(ctx context.Context) (interface{}, error) {
	var currCreds Credentials
	if cur := p.creds.Load(); cur != nil {
		currCreds = cur.creds
	}
	if currCreds.HasKeys() && !currCreds.Expired() {
		return currCreds, nil
	}

	newCreds, err := p.provider.Retrieve(ctx)
	if err != nil {
		handleFailToRefresh := defaultHandleFailToRefresh
		if cs, ok := p.provider.(HandleFailRefreshCredentialsCacheStrategy); ok {
			handleFailToRefresh = cs.HandleFailToRefresh
		}
		newCreds, err = handleFailToRefresh(ctx, currCreds, err)
		if err != nil {
			return Credentials{}, fmt.Errorf("failed to refresh cached credentials, %w", err)
		}
	}

	if newCreds.CanExpire && p.options.ExpiryWindow > 0 {
		adjustExpiresBy := defaultAdjustExpiresBy
		if cs, ok := p.provider.(AdjustExpiresByCredentialsCacheStrategy); ok {
			adjustExpiresBy = cs.AdjustExpiresBy
		}

		randFloat64, err := sdkrand.CryptoRandFloat64()
		if err != nil {
			return Credentials{}, fmt.Errorf("failed to get random provider, %w", err)
		}

		var jitter time.Duration
		if p.options.ExpiryWindowJitterFrac > 0 {
			jitter = time.Duration(randFloat64 *
				p.options.ExpiryWindowJitterFrac * float64(p.options.ExpiryWindow))
		}

		newCreds, err = adjustExpiresBy(newCreds, -(p.options.ExpiryWindow - jitter))
		if err != nil {
			return Credentials{}, fmt.Errorf("failed to adjust credentials expires, %w", err)
		}
	}

	p.creds.Store(&cachedCredentials{creds: newCreds})
	return newCreds, nil
}

// ProviderSources returns a list of where the underlying credential provider
// has been sourced, if available. Returns empty if the provider doesn't implement
// the interface
func (p *CredentialsCache) ProviderSources() []CredentialSource {
	asSource, ok := p.provider.(CredentialProviderSource)
	if !ok {
		return []CredentialSource{}
	}
	return asSource.ProviderSources()
}

// Invalidate will invalidate the cached credentials. The next call to Retrieve
// will cause the provider's Retrieve method to be called.
func (p *CredentialsCache) Invalidate() {
	p.creds.Store(nil)
	p.cachedErr.Store(nil)
}

// InvalidateCredentials marks the cached credentials for refresh if they're the ones a service rejected, keeping them and any backoff in place.
func (p *CredentialsCache) InvalidateCredentials(rejected Credentials) {
	for {
		v := p.creds.Load()
		if v == nil || !v.creds.HasKeys() {
			return
		}
		if v.creds.AccessKeyID != rejected.AccessKeyID {
			return
		}
		if p.creds.CompareAndSwap(v, v.withExpiration(sdk.NowTime())) {
			return
		}
	}
}

// IsCredentialsProvider returns whether credential provider wrapped by CredentialsCache
// matches the target provider type.
func (p *CredentialsCache) IsCredentialsProvider(target CredentialsProvider) bool {
	return IsCredentialsProvider(p.provider, target)
}

// HandleFailRefreshCredentialsCacheStrategy is an interface for
// CredentialsCache to allow CredentialsProvider  how failed to refresh
// credentials is handled.
type HandleFailRefreshCredentialsCacheStrategy interface {
	// Given the previously cached Credentials, if any, and refresh error, may
	// returns new or modified set of Credentials, or error.
	//
	// Credential caches may use default implementation if nil.
	HandleFailToRefresh(context.Context, Credentials, error) (Credentials, error)
}

// defaultHandleFailToRefresh returns the passed in error.
func defaultHandleFailToRefresh(ctx context.Context, _ Credentials, err error) (Credentials, error) {
	return Credentials{}, err
}

// AdjustExpiresByCredentialsCacheStrategy is an interface for CredentialCache
// to allow CredentialsProvider to intercept adjustments to Credentials expiry
// based on expectations and use cases of CredentialsProvider.
//
// Credential caches may use default implementation if nil.
type AdjustExpiresByCredentialsCacheStrategy interface {
	// Given a Credentials as input, applying any mutations and
	// returning the potentially updated Credentials, or error.
	AdjustExpiresBy(Credentials, time.Duration) (Credentials, error)
}

// defaultAdjustExpiresBy adds the duration to the passed in credentials Expires,
// and returns the updated credentials value. If Credentials value's CanExpire
// is false, the passed in credentials are returned unchanged.
func defaultAdjustExpiresBy(creds Credentials, dur time.Duration) (Credentials, error) {
	if !creds.CanExpire {
		return creds, nil
	}

	creds.Expires = creds.Expires.Add(dur)
	return creds, nil
}

// inScopeSources is a map of in-scope providers.
var inScopeSources = map[CredentialSource]struct{}{
	CredentialSourceIMDS:                 {},
	CredentialSourceHTTP:                 {},
	CredentialSourceSTSAssumeRole:        {},
	CredentialSourceSTSAssumeRoleWebID:   {},
	CredentialSourceEnvVarsSTSWebIDToken: {},
	CredentialSourceProfileSTSWebIDToken: {},
	CredentialSourceSSO:                  {},
	CredentialSourceProfileSSO:           {},
	CredentialSourceSSOLegacy:            {},
	CredentialSourceProfileSSOLegacy:     {},
	CredentialSourceLogin:                {},
	CredentialSourceProfileLogin:         {},
}

// inScope reports whether the provider is in inScopeSources.
func inScope(provider CredentialsProvider) bool {
	asSource, ok := provider.(CredentialProviderSource)
	if !ok {
		return false
	}
	sources := asSource.ProviderSources()
	if len(sources) == 0 {
		return false
	}
	_, ok = inScopeSources[sources[len(sources)-1]]
	return ok
}

// mandatoryWindow returns 1 minute, or the advisory window if that's shorter.
func (e *cachedCredentials) mandatoryWindow() time.Duration {
	return min(time.Minute, e.advisoryWindow)
}

func (e *cachedCredentials) refreshNeeded(now time.Time) bool {
	if !e.creds.CanExpire {
		return false
	}
	return !now.Before(e.creds.Expires.Add(-e.advisoryWindow))
}

func (e *cachedCredentials) mandatoryRefreshNeeded(now time.Time) bool {
	if !e.creds.CanExpire {
		return false
	}
	return !now.Before(e.creds.Expires.Add(-e.mandatoryWindow()))
}

func (e *cachedCredentials) rateLimited(now time.Time) bool {
	return !e.nextRefreshAllowedAt.IsZero() && now.Before(e.nextRefreshAllowedAt)
}

// withBackoff returns a copy of the cached credentials with the backoff set to until.
func (e *cachedCredentials) withBackoff(until time.Time) *cachedCredentials {
	next := *e
	next.nextRefreshAllowedAt = until
	return &next
}

// withExpiration returns a copy of the cached credentials with the expiration set to t.
func (e *cachedCredentials) withExpiration(t time.Time) *cachedCredentials {
	next := *e
	next.creds.Expires = t
	return &next
}

// nonRecoverableErrCached returns the cached non-recoverable error if it hasn't expired yet.
func (p *CredentialsCache) nonRecoverableErrCached(now time.Time) (error, bool) {
	ce := p.cachedErr.Load()
	if ce == nil || !now.Before(ce.expiresAt) {
		return nil, false
	}
	return ce.err, true
}

// advisoryWindowFor returns the customer's ExpiryWindow if set, otherwise a window based on how long the credentials last.
func (p *CredentialsCache) advisoryWindowFor(now time.Time, creds Credentials) time.Duration {
	if p.options.ExpiryWindow > 0 {
		return p.jitteredConfiguredWindow()
	}
	if !creds.CanExpire {
		return 0
	}
	switch lifetime := creds.Expires.Sub(now); {
	case lifetime <= 20*time.Minute:
		return 5 * time.Minute
	case lifetime < 90*time.Minute:
		return 15 * time.Minute
	default:
		return 60 * time.Minute
	}
}

// jitteredConfiguredWindow returns ExpiryWindow shortened by a random amount, up to ExpiryWindowJitterFrac of it.
func (p *CredentialsCache) jitteredConfiguredWindow() time.Duration {
	window := p.options.ExpiryWindow
	if p.options.ExpiryWindowJitterFrac <= 0 {
		return window
	}
	frac := p.options.ExpiryWindowJitterFrac
	if frac > 1 {
		frac = 1
	}
	r, err := p.rand()
	if err != nil {
		return window
	}
	return window - time.Duration(r*frac*float64(window))
}

// jitter returns a uniformly random duration in [min, max] seconds, falling back to min on an read failure.
func (p *CredentialsCache) jitter(min, max int) time.Duration {
	r, err := p.rand()
	if err != nil {
		return time.Duration(min) * time.Second
	}
	return time.Duration(min)*time.Second + time.Duration(r*float64(max-min))*time.Second
}
