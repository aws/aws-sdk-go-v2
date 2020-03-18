package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
)

const (
	// StaticCredentialsProviderName provides a name of Static provider
	StaticCredentialsProviderName = "StaticCredentialsProvider"

	// ErrCodeStaticCredentialsEmpty is emitted when static credentials are empty.
	ErrCodeStaticCredentialsEmpty = "EmptyStaticCreds"
)

// A StaticCredentialsProvider is a set of credentials which are set programmatically,
// and will never expire.
type StaticCredentialsProvider struct {
	Value Credentials
}

// NewStaticCredentialsProvider return a StaticCredentialsProvider initialized with the AWS credentials
// passed in.
func NewStaticCredentialsProvider(key, secret, session string) StaticCredentialsProvider {
	return StaticCredentialsProvider{
		Value: Credentials{
			AccessKeyID:     key,
			SecretAccessKey: secret,
			SessionToken:    session,
		},
	}
}

// Retrieve returns the credentials or error if the credentials are invalid.
func (s StaticCredentialsProvider) Retrieve(ctx context.Context) (Credentials, error) {
	v := s.Value
	if v.AccessKeyID == "" || v.SecretAccessKey == "" {
		return Credentials{Source: StaticCredentialsProviderName}, awserr.New(ErrCodeStaticCredentialsEmpty, "static credentials are empty", nil)
	}

	if len(v.Source) == 0 {
		v.Source = StaticCredentialsProviderName
	}

	return v, nil
}

// IsExpired returns if the credentials are expired.
//
// For StaticCredentialsProvider, the credentials never expired.
func (s StaticCredentialsProvider) IsExpired() bool {
	return false
}
