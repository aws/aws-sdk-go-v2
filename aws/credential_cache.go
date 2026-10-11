package aws

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	internalcontext "github.com/aws/aws-sdk-go-v2/internal/context"
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
	provider  CredentialsProvider
	refresher refresher
}

// refresher is the refresh strategy NewCredentialsCache picks for the provider.
type refresher interface {
	retrieve(ctx context.Context) (Credentials, error)
	invalidate()
	invalidateCredentials(rejected Credentials)
}

// staticStableRefresher refreshes in-scope providers with advisory and mandatory windows, a refresh backoff, and a cached non-recoverable error.
type staticStableRefresher struct {
	provider CredentialsProvider
	options  CredentialsCacheOptions
	creds    atomic.Pointer[cacheEntry]
	sf       singleflight.Group
}

type cacheEntry struct {
	creds                Credentials
	advisoryWindow       time.Duration
	nextRefreshAllowedAt time.Time

	err          error
	errExpiresAt time.Time
}

// outOfScopeRefresher keeps the existing caching behavior, including the provider's cache strategy hooks, for providers outside the refresh lifecycle, like process and custom providers.
type outOfScopeRefresher struct {
	provider CredentialsProvider
	options  CredentialsCacheOptions
	creds    atomic.Value
	sf       singleflight.Group
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

	var r refresher = &outOfScopeRefresher{provider: provider, options: options}
	if inScope(provider) {
		r = &staticStableRefresher{provider: provider, options: options}
	}
	return &CredentialsCache{
		provider:  provider,
		refresher: r,
	}
}

// Retrieve returns the credentials. If the credentials have already been
// retrieved, and not expired the cached credentials will be returned. If the
// credentials have not been retrieved yet, or expired the provider's Retrieve
// method will be called.
//
// Returns and error if the provider's retrieve method returns an error.
func (p *CredentialsCache) Retrieve(ctx context.Context) (Credentials, error) {
	return p.refresher.retrieve(ctx)
}

// retrieve returns the cached credentials until they expire, then refreshes them once for all waiting callers.
func (p *outOfScopeRefresher) retrieve(ctx context.Context) (Credentials, error) {
	if creds, ok := p.getCreds(); ok && !creds.Expired() {
		return creds, nil
	}

	resCh := p.sf.DoChan("", func() (interface{}, error) {
		return p.singleRetrieve(&suppressedContext{ctx})
	})
	select {
	case res := <-resCh:
		return res.Val.(Credentials), res.Err
	case <-ctx.Done():
		return Credentials{}, &RequestCanceledError{Err: ctx.Err()}
	}
}

// singleRetrieve asks the provider for new credentials, applies the provider's cache strategy hooks, and caches the result.
func (p *outOfScopeRefresher) singleRetrieve(ctx context.Context) (interface{}, error) {
	currCreds, ok := p.getCreds()
	if ok && !currCreds.Expired() {
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

	p.creds.Store(&newCreds)
	return newCreds, nil
}

// getCreds returns the currently stored credentials and true. Returning false
// if no credentials were stored.
func (p *outOfScopeRefresher) getCreds() (Credentials, bool) {
	v := p.creds.Load()
	if v == nil {
		return Credentials{}, false
	}

	c := v.(*Credentials)
	if c == nil || !c.HasKeys() {
		return Credentials{}, false
	}

	return *c, true
}

// invalidate drops the cached credentials so the next retrieve calls the provider.
func (p *outOfScopeRefresher) invalidate() {
	p.creds.Store((*Credentials)(nil))
}

// invalidateCredentials is a no-op for providers outside the refresh lifecycle.
func (p *outOfScopeRefresher) invalidateCredentials(Credentials) {}

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
	p.refresher.invalidate()
}

// InvalidateCredentials marks the cached credentials for refresh if they're the ones a service rejected, keeping them and any backoff in place.
func (p *CredentialsCache) InvalidateCredentials(rejected Credentials) {
	p.refresher.invalidateCredentials(rejected)
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

// Refresh timing values for the credentials cache.
const (
	refreshBackoffMin = 300 * time.Second
	refreshBackoffMax = 600 * time.Second

	nonRecoverableErrTTLMin = 1 * time.Second
	nonRecoverableErrTTLMax = 5 * time.Second

	mandatoryRefreshWindow = time.Minute

	shortLifetimeMax  = 20 * time.Minute
	mediumLifetimeMax = 90 * time.Minute

	shortAdvisoryWindow  = 5 * time.Minute
	mediumAdvisoryWindow = 15 * time.Minute
	longAdvisoryWindow   = 60 * time.Minute
)

// retrieve picks an initial, advisory, or mandatory refresh based on where the cached credentials are in their lifetime.
func (p *staticStableRefresher) retrieve(ctx context.Context) (Credentials, error) {
	now := sdk.NowTime()

	currCreds := p.creds.Load()
	if currCreds == nil || !currCreds.creds.HasKeys() {
		return p.waitForRetrieve(ctx)
	}
	if !currCreds.refreshNeeded(now) {
		return currCreds.creds, nil
	}
	if currCreds.rateLimited(now) {
		return currCreds.creds, nil
	}

	if !currCreds.mandatoryRefreshNeeded(now) {
		return p.advisoryRefresh(ctx, currCreds)
	}
	return p.waitForRetrieve(ctx)
}

// advisoryRefresh starts a singleRetrieve if none is in flight and waits for it; callers that find one in flight get the cached credentials right away.
func (p *staticStableRefresher) advisoryRefresh(ctx context.Context, currCreds *cacheEntry) (Credentials, error) {
	ch, started := p.sf.TryDoChan("", func() (interface{}, error) {
		return p.singleRetrieve(&suppressedContext{ctx})
	})
	if !started {
		return currCreds.creds, nil
	}
	return waitForResult(ctx, ch, currCreds)
}

// waitForRetrieve starts a singleRetrieve, or joins the one in flight, and waits for it; used for the initial and mandatory refresh.
func (p *staticStableRefresher) waitForRetrieve(ctx context.Context) (Credentials, error) {
	ch := p.sf.DoChan("", func() (interface{}, error) {
		return p.singleRetrieve(&suppressedContext{ctx})
	})
	return waitForResult(ctx, ch, nil)
}

// waitForResult returns the refresh result; if ctx ends first it returns the fallback credentials, or a RequestCanceledError when there are none.
func waitForResult(ctx context.Context, ch <-chan singleflight.Result, fallback *cacheEntry) (Credentials, error) {
	select {
	case res := <-ch:
		if res.Err != nil {
			return Credentials{}, res.Err
		}
		return res.Val.(Credentials), nil
	case <-ctx.Done():
		if fallback != nil {
			return fallback.creds, nil
		}
		return Credentials{}, &RequestCanceledError{Err: ctx.Err()}
	}
}

// singleRetrieve asks the provider for new credentials and caches the result. If the provider fails, it keeps
// the old credentials and waits before trying again, or raises the error right away if retrying won't help.
func (p *staticStableRefresher) singleRetrieve(ctx context.Context) (Credentials, error) {
	now := sdk.NowTime()
	currCreds := p.creds.Load()
	hasCreds := currCreds != nil && currCreds.creds.HasKeys()
	if hasCreds && !currCreds.refreshNeeded(now) {
		return currCreds.creds, nil
	}
	if err := currCreds.nonRecoverableErr(now); err != nil {
		return Credentials{}, err
	}
	if hasCreds && currCreds.rateLimited(now) {
		return currCreds.creds, nil
	}

	if hasCreds {
		ctx = internalcontext.SetStaleAccessKeyID(ctx, currCreds.creds.AccessKeyID)
	}
	newCreds, err := p.provider.Retrieve(ctx)
	now = sdk.NowTime()

	if err == nil && newCreds.CanExpire && !newCreds.Expires.After(now) {
		err = errors.New("credential source returned already-expired credentials")
	}

	if err != nil {
		var nr *nonrecoverable.Failure
		if errors.As(err, &nr) {
			p.creds.Store(currCreds.withErr(err, now.Add(jitter(nonRecoverableErrTTLMin, nonRecoverableErrTTLMax))))
			return Credentials{}, err
		}

		if !hasCreds {
			return Credentials{}, fmt.Errorf("failed to refresh cached credentials, %w", err)
		}

		backoff := jitter(refreshBackoffMin, refreshBackoffMax)
		p.creds.Store(currCreds.withBackoff(now.Add(backoff)))
		logFailedRefresh(ctx, err, backoff)
		return currCreds.creds, nil
	}

	p.creds.Store(&cacheEntry{
		creds:          newCreds,
		advisoryWindow: p.advisoryWindowFor(now, newCreds),
	})
	return newCreds, nil
}

// invalidate drops the cached credentials and any cached error.
func (p *staticStableRefresher) invalidate() {
	p.creds.Store(nil)
}

// invalidateCredentials expires the cached credentials if they're the ones a service rejected.
func (p *staticStableRefresher) invalidateCredentials(rejected Credentials) {
	v := p.creds.Load()
	if v == nil || !v.creds.HasKeys() || v.creds.AccessKeyID != rejected.AccessKeyID {
		return
	}
	p.creds.CompareAndSwap(v, v.withExpiration(sdk.NowTime()))
}

// logFailedRefresh logs a warning that a refresh failed and the cached credentials will be used until the backoff ends.
func logFailedRefresh(ctx context.Context, err error, backoff time.Duration) {
	logger := middleware.GetLogger(ctx)
	logger.Logf(logging.Warn,
		"Credential refresh failed: %v. The SDK will continue using cached credentials. "+
			"A refresh of these credentials will be attempted again after %d seconds.",
		err, int(backoff/time.Second))
}

// mandatoryWindow returns 1 minute, or the advisory window if that's shorter.
func (e *cacheEntry) mandatoryWindow() time.Duration {
	return min(mandatoryRefreshWindow, e.advisoryWindow)
}

func (e *cacheEntry) refreshNeeded(now time.Time) bool {
	if !e.creds.CanExpire {
		return false
	}
	return !now.Before(e.creds.Expires.Add(-e.advisoryWindow))
}

func (e *cacheEntry) mandatoryRefreshNeeded(now time.Time) bool {
	if !e.creds.CanExpire {
		return false
	}
	return !now.Before(e.creds.Expires.Add(-e.mandatoryWindow()))
}

func (e *cacheEntry) rateLimited(now time.Time) bool {
	return !e.nextRefreshAllowedAt.IsZero() && now.Before(e.nextRefreshAllowedAt)
}

// withBackoff returns a copy of the cached credentials with the backoff set to until.
func (e *cacheEntry) withBackoff(until time.Time) *cacheEntry {
	next := *e
	next.nextRefreshAllowedAt = until
	return &next
}

// withExpiration returns a copy of the cached credentials with the expiration set to t.
func (e *cacheEntry) withExpiration(t time.Time) *cacheEntry {
	next := *e
	next.creds.Expires = t
	return &next
}

// withErr returns a copy of the cache entry with err cached until the given time.
func (e *cacheEntry) withErr(err error, until time.Time) *cacheEntry {
	var next cacheEntry
	if e != nil {
		next = *e
	}
	next.err = err
	next.errExpiresAt = until
	return &next
}

// nonRecoverableErr returns the cached non-recoverable error, or nil if none is cached or it has expired.
func (e *cacheEntry) nonRecoverableErr(now time.Time) error {
	if e == nil || e.err == nil || now.After(e.errExpiresAt) {
		return nil
	}
	return e.err
}

// advisoryWindowFor returns the customer's ExpiryWindow if set, otherwise a window based on how long the credentials last.
func (p *staticStableRefresher) advisoryWindowFor(now time.Time, creds Credentials) time.Duration {
	if p.options.ExpiryWindow > 0 {
		return p.jitteredConfiguredWindow()
	}
	if !creds.CanExpire {
		return 0
	}
	switch lifetime := creds.Expires.Sub(now); {
	case lifetime <= shortLifetimeMax:
		return shortAdvisoryWindow
	case lifetime < mediumLifetimeMax:
		return mediumAdvisoryWindow
	default:
		return longAdvisoryWindow
	}
}

// jitteredConfiguredWindow returns ExpiryWindow shortened by a random amount, up to ExpiryWindowJitterFrac of it.
func (p *staticStableRefresher) jitteredConfiguredWindow() time.Duration {
	window := p.options.ExpiryWindow
	if p.options.ExpiryWindowJitterFrac <= 0 {
		return window
	}
	frac := p.options.ExpiryWindowJitterFrac
	if frac > 1 {
		frac = 1
	}
	r, err := sdkrand.CryptoRandFloat64()
	if err != nil {
		return window
	}
	return window - time.Duration(r*frac*float64(window))
}

// jitter returns a uniformly random duration in [lo, hi), falling back to lo on a read failure.
func jitter(lo, hi time.Duration) time.Duration {
	r, err := sdkrand.CryptoRandFloat64()
	if err != nil {
		return lo
	}
	return lo + time.Duration(r*float64(hi-lo))
}
