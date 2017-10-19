package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
)

// StaticProviderName provides a name of Static provider
const StaticProviderName = "StaticProvider"

var (
	// ErrStaticCredentialsEmpty is emitted when static credentials are empty.
	ErrStaticCredentialsEmpty = awserr.New("EmptyStaticCreds", "static credentials are empty", nil)
)

// A StaticProvider is a set of credentials which are set programmatically,
// and will never expire.
type StaticProvider struct {
	Value Credentials
}

// NewStaticProvider return a StaticProvider initialized with the AWS credentials
// passed in.
func NewStaticProvider(key, secret, session string) StaticProvider {
	return StaticProvider{
		Value: Credentials{
			AccessKeyID:     key,
			SecretAccessKey: secret,
			SessionToken:    session,
		},
	}
}

// Retrieve returns the credentials or error if the credentials are invalid.
func (s StaticProvider) Retrieve() (Credentials, error) {
	v := s.Value
	if v.AccessKeyID == "" || v.SecretAccessKey == "" {
		return Credentials{ProviderName: StaticProviderName}, ErrStaticCredentialsEmpty
	}

	if len(v.ProviderName) == 0 {
		v.ProviderName = StaticProviderName
	}

	return v, nil
}

// IsExpired returns if the credentials are expired.
//
// For StaticProvider, the credentials never expired.
func (s StaticProvider) IsExpired() bool {
	return false
}
