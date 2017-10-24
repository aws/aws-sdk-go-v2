package aws

import (
	"sync"
	"time"
)

// AnonymousCredentials is an empty Credential object that can be used as
// dummy placeholder credentials for requests that do not need signed.
//
// This CredentialsLoader can be used to configure a service to not sign requests
// when making service API calls. For example, when accessing public
// s3 objects.
//
//     cfg, err := external.LoadDefaultAWSconfig()
//     cfg.CredentialsLoader = aws.AnonymousCredentials
//
//     svc := s3.New(cfg)
//     // Get Public S3 Object
var AnonymousCredentials = NewCredentialsLoader(StaticCredentialsProvider{
	Value: Credentials{Source: "AnonymousCredentials"},
})

// A Credentials is the AWS credentials value for individual credential fields.
type Credentials struct {
	// AWS Access key ID
	AccessKeyID string

	// AWS Secret Access Key
	SecretAccessKey string

	// AWS Session Token
	SessionToken string

	// Source of the credentials
	Source string
}

// HasKeys returns if the credentials keys are set.
func (v Credentials) HasKeys() bool {
	return len(v.AccessKeyID) > 0 && len(v.SecretAccessKey) > 0
}

// A CredentialsProvider is the interface for any component which will provide credentials
// Credentials. A CredentialsProvider is required to manage its own Expired state, and what to
// be expired means.
//
// The CredentialsProvider should not need to implement its own mutexes, because
// that will be managed by CredentialsLoader.
type CredentialsProvider interface {
	// Retrieve returns nil if it successfully retrieved the value.
	// Error is returned if the value were not obtainable, or empty.
	Retrieve() (Credentials, error)
}

// Expirer is an interface for expiring values.
type Expirer interface {
	// Retrunes when the value expires.
	Expires() time.Time
}

// A Expiry provides shared expiration logic to be used by credentials
// providers to implement expiry functionality.
//
// The best method to use this struct is as an anonymous field within the
// CredentialsProvider's struct.
type Expiry struct {
	// The date/time when to expire on
	expiration time.Time
}

// SetExpiration sets the expiration IsExpired will check when called.
//
// If window is greater than 0 the expiration time will be reduced by the
// window value.
//
// Using a window is helpful to trigger credentials to expire sooner than
// the expiration time given to ensure no requests are made with expired
// tokens.
func (e *Expiry) SetExpiration(expiration time.Time, window time.Duration) {
	e.expiration = expiration
	if window > 0 {
		e.expiration = e.expiration.Add(-window)
	}
}

// Expires returns when the credentials will expire.
func (e *Expiry) Expires() time.Time {
	return e.expiration
}

// A CredentialsLoader provides synchronous safe retrieval of AWS credentials Credentials.
// CredentialsLoader will cache the credentials value until they expire. Once the value
// expires the next Get will attempt to retrieve valid credentials.
//
// CredentialsLoader is safe to use across multiple goroutines and will manage the
// synchronous state so the Providers do not need to implement their own
// synchronization.
//
// The first CredentialsLoader.Get() will always call CredentialsProvider.Retrieve() to get the
// first instance of the credentials Credentials. All calls to Get() after that
// will return the cached credentials Credentials until IsExpired() returns true.
type CredentialsLoader struct {
	creds        Credentials
	forceRefresh bool
	m            sync.Mutex

	Provider    CredentialsProvider
	CurrentTime func() time.Time
}

// NewCredentialsLoader returns a pointer to a new CredentialsLoader with the CredentialsProvider set.
func NewCredentialsLoader(Provider CredentialsProvider) *CredentialsLoader {
	return &CredentialsLoader{
		Provider:     Provider,
		forceRefresh: true,
	}
}

// Get returns the credentials value, or error if the credentials Credentials failed
// to be retrieved.
//
// Will return the cached credentials Credentials if it has not expired. If the
// credentials Credentials has expired the CredentialsProvider's Retrieve() will be called
// to refresh the credentials.
//
// If CredentialsLoader.Expire() was called the credentials Credentials will be force
// expired, and the next call to Get() will cause them to be refreshed.
func (c *CredentialsLoader) Get() (Credentials, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.isExpired() {
		creds, err := c.Provider.Retrieve()
		if err != nil {
			return Credentials{}, err
		}
		c.creds = creds
		c.forceRefresh = false
	}

	return c.creds, nil
}

// Expire expires the credentials and forces them to be retrieved on the
// next call to Get().
//
// This will override the CredentialsProvider's expired state, and force CredentialsLoader
// to call the CredentialsProvider's Retrieve().
func (c *CredentialsLoader) Expire() {
	c.m.Lock()
	defer c.m.Unlock()

	c.forceRefresh = true
}

// IsExpired returns if the credentials are no longer valid, and need
// to be retrieved.
//
// If the CredentialsLoader were forced to be expired with Expire() this will
// reflect that override.
func (c *CredentialsLoader) IsExpired() bool {
	c.m.Lock()
	defer c.m.Unlock()

	return c.isExpired()
}

func (c *CredentialsLoader) isExpired() bool {
	if c.forceRefresh {
		return true
	}

	e, ok := c.Provider.(Expirer)
	if !ok {
		return false
	}

	nowTime := c.CurrentTime
	if nowTime == nil {
		nowTime = time.Now
	}

	return e.Expires().Before(nowTime())
}
