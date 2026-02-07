package credentials

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	// StaticCredentialsName provides a name of Static provider
	StaticCredentialsName = "StaticCredentials"
)

// StaticCredentialsEmptyError is emitted when static credentials are empty.
type StaticCredentialsEmptyError struct{}

func (*StaticCredentialsEmptyError) Error() string {
	return "static credentials are empty"
}

// AccessKeyIDInvalidWhitespaceError is emitted when AccessKeyID contains invalid whitespace.
type AccessKeyIDInvalidWhitespaceError struct{}

func (*AccessKeyIDInvalidWhitespaceError) Error() string {
	return "AccessKeyID contains invalid whitespace"
}

// SecretAccessKeyInvalidWhitespaceError is emitted when SecretAccessKey contains invalid whitespace.
type SecretAccessKeyInvalidWhitespaceError struct{}

func (*SecretAccessKeyInvalidWhitespaceError) Error() string {
	return "SecretAccessKey contains invalid whitespace"
}

// SessionTokenInvalidWhitespaceError is emitted when SessionToken contains invalid whitespace.
type SessionTokenInvalidWhitespaceError struct{}

func (*SessionTokenInvalidWhitespaceError) Error() string {
	return "SessionToken contains invalid whitespace"
}

// A StaticCredentialsProvider is a set of credentials which are set, and will
// never expire.
type StaticCredentialsProvider struct {
	Value aws.Credentials
	// These values are for reporting purposes and are not meant to be set up directly
	Source []aws.CredentialSource
}

// ProviderSources returns the credential chain that was used to construct this provider
func (s StaticCredentialsProvider) ProviderSources() []aws.CredentialSource {
	if s.Source == nil {
		return []aws.CredentialSource{aws.CredentialSourceCode} // If no source has been set, assume this is used directly which means hardcoded creds
	}
	return s.Source
}

// NewStaticCredentialsProvider return a StaticCredentialsProvider initialized with the AWS
// credentials passed in.
func NewStaticCredentialsProvider(key, secret, session string) StaticCredentialsProvider {
	return StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     key,
			SecretAccessKey: secret,
			SessionToken:    session,
		},
	}
}

// Retrieve returns the credentials or error if the credentials are invalid.
func (s StaticCredentialsProvider) Retrieve(_ context.Context) (aws.Credentials, error) {
	v := s.Value

	if strings.ContainsAny(v.AccessKeyID, " \t\r\n") {
		return aws.Credentials{
			Source: StaticCredentialsName,
		}, &AccessKeyIDInvalidWhitespaceError{}
	}
	if strings.ContainsAny(v.SecretAccessKey, " \t\r\n") {
		return aws.Credentials{
			Source: StaticCredentialsName,
		}, &SecretAccessKeyInvalidWhitespaceError{}
	}
	if strings.ContainsAny(v.SessionToken, " \t\r\n") {
		return aws.Credentials{
			Source: StaticCredentialsName,
		}, &SessionTokenInvalidWhitespaceError{}
	}

	if v.AccessKeyID == "" || v.SecretAccessKey == "" {
		return aws.Credentials{
			Source: StaticCredentialsName,
		}, &StaticCredentialsEmptyError{}
	}

	if len(v.Source) == 0 {
		v.Source = StaticCredentialsName
	}

	return v, nil
}
