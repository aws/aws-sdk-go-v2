package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/stretchr/testify/assert"
)

type stubProvider struct {
	creds   Credentials
	expired bool
	err     error
}

func (s *stubProvider) Retrieve() (Credentials, error) {
	s.expired = false
	s.creds.Source = "stubProvider"
	return s.creds, s.err
}
func (s *stubProvider) IsExpired() bool {
	return s.expired
}

func TestCredentialsGet(t *testing.T) {
	c := NewCredentialsLoader(&stubProvider{
		creds: Credentials{
			AccessKeyID:     "AKID",
			SecretAccessKey: "SECRET",
			SessionToken:    "",
		},
		expired: true,
	})

	creds, err := c.Get()
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, "AKID", creds.AccessKeyID, "Expect access key ID to match")
	assert.Equal(t, "SECRET", creds.SecretAccessKey, "Expect secret access key to match")
	assert.Empty(t, creds.SessionToken, "Expect session token to be empty")
}

func TestCredentialsGetWithError(t *testing.T) {
	c := NewCredentialsLoader(&stubProvider{err: awserr.New("provider error", "", nil), expired: true})

	_, err := c.Get()
	assert.Equal(t, "provider error", err.(awserr.Error).Code(), "Expected provider error")
}

func TestCredentialsExpire(t *testing.T) {
	stub := &stubProvider{}
	c := NewCredentialsLoader(stub)

	stub.expired = false
	assert.True(t, c.IsExpired(), "Expected to start out expired")
	c.Expire()
	assert.True(t, c.IsExpired(), "Expected to be expired")

	c.forceRefresh = false
	assert.False(t, c.IsExpired(), "Expected not to be expired")

	stub.expired = true
	assert.True(t, c.IsExpired(), "Expected to be expired")
}

func TestCredentialsGetWithProviderName(t *testing.T) {
	stub := &stubProvider{}

	c := NewCredentialsLoader(stub)

	creds, err := c.Get()
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, creds.Source, "stubProvider", "Expected provider name to match")
}
