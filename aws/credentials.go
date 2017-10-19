// Package aws need to fillout this docs
//
// The CredentialsLoader is the primary method of getting access to and managing
// credentials Values. Using dependency injection retrieval of the credential
// values is handled by a object which satisfies the CredentialsProvider interface.
//
// By default the CredentialsLoader.Get() will cache the successful result of a
// CredentialsProvider's Retrieve() until CredentialsProvider.IsExpired() returns true. At which
// point CredentialsLoader will call CredentialsProvider's Retrieve() to get new credential Credentials.
//
// The CredentialsProvider is responsible for determining when credentials Credentials have expired.
// It is also important to note that CredentialsLoader will always call Retrieve the
// first time CredentialsLoader.Get() is called.
//
// Example of using the environment variable credentials.
//
//     creds := aws.NewEnvCredentials()
//
//     // Retrieve the credentials value
//     credValue, err := creds.Get()
//     if err != nil {
//         // handle error
//     }
//
// Example of forcing credentials to expire and be refreshed on the next Get().
// This may be helpful to proactively expire credentials and refresh them sooner
// than they would naturally expire on their own.
//
//     creds := aws.NewCredentials(&ec2rolecreds.EC2RoleProvider{})
//     creds.Expire()
//     credsValue, err := creds.Get()
//     // New credentials will be retrieved instead of from cache.
//
//
// Custom CredentialsProvider
//
// Each CredentialsProvider built into this package also provides a helper method to generate
// a CredentialsLoader pointer setup with the CredentialsProvider. To use a custom CredentialsProvider just
// create a type which satisfies the CredentialsProvider interface and pass it to the
// NewCredentials method.
//
//     type MyProvider struct{}
//     func (m *MyProvider) Retrieve() (Credentials, error) {...}
//     func (m *MyProvider) IsExpired() bool {...}
//
//     creds := aws.NewCredentials(&MyProvider{})
//     credValue, err := creds.Get()
//
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
// s3 buckets.
//
//     svc := s3.New(session.Must(session.NewSession(&aws.Config{
//       CredentialsLoader: aws.AnonymousCredentials,
//     })))
//     // Access public S3 buckets.
//
// @readonly
var AnonymousCredentials = NewCredentialsLoader(StaticProvider{
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

	// IsExpired returns if the credentials are no longer valid, and need
	// to be retrieved.
	IsExpired() bool
}

// An ErrorProvider is a stub credentials CredentialsProvider that always returns an error
// this is used by the SDK when construction a known CredentialsProvider is not possible
// due to an error.
type ErrorProvider struct {
	// The error to be returned from Retrieve
	Err error

	// The CredentialsProvider name to set on the Retrieved returned Credentials
	Source string
}

// Retrieve will always return the error that the ErrorProvider was created with.
func (p ErrorProvider) Retrieve() (Credentials, error) {
	return Credentials{Source: p.Source}, p.Err
}

// IsExpired will always return not expired.
func (p ErrorProvider) IsExpired() bool {
	return false
}

// A Expiry provides shared expiration logic to be used by credentials
// providers to implement expiry functionality.
//
// The best method to use this struct is as an anonymous field within the
// CredentialsProvider's struct.
//
// Example:
//     type EC2RoleProvider struct {
//         Expiry
//         ...
//     }
type Expiry struct {
	// The date/time when to expire on
	expiration time.Time

	// If set will be used by IsExpired to determine the current time.
	// Defaults to time.Now if CurrentTime is not set.  Available for testing
	// to be able to mock out the current time.
	CurrentTime func() time.Time
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

// IsExpired returns if the credentials are expired.
func (e *Expiry) IsExpired() bool {
	if e.CurrentTime == nil {
		e.CurrentTime = time.Now
	}
	return e.expiration.Before(e.CurrentTime())
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

	Provider CredentialsProvider
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

// isExpired helper method wrapping the definition of expired credentials.
func (c *CredentialsLoader) isExpired() bool {
	return c.forceRefresh || c.Provider.IsExpired()
}
