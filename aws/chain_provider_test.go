package aws

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/stretchr/testify/assert"
)

type secondStubProvider struct {
	creds   Credentials
	expired bool
	err     error
}

func (s *secondStubProvider) Retrieve() (Credentials, error) {
	s.expired = false
	s.creds.ProviderName = "secondStubProvider"
	return s.creds, s.err
}
func (s *secondStubProvider) IsExpired() bool {
	return s.expired
}

func TestChainProviderWithNames(t *testing.T) {
	p := &ChainProvider{
		Providers: []CredentialsProvider{
			&stubProvider{err: awserr.New("FirstError", "first provider error", nil)},
			&stubProvider{err: awserr.New("SecondError", "second provider error", nil)},
			&secondStubProvider{
				creds: Credentials{
					AccessKeyID:     "AKIF",
					SecretAccessKey: "NOSECRET",
					SessionToken:    "",
				},
			},
			&stubProvider{
				creds: Credentials{
					AccessKeyID:     "AKID",
					SecretAccessKey: "SECRET",
					SessionToken:    "",
				},
			},
		},
	}

	creds, err := p.Retrieve()
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, "secondStubProvider", creds.ProviderName, "Expect provider name to match")

	// Also check credentials
	assert.Equal(t, "AKIF", creds.AccessKeyID, "Expect access key ID to match")
	assert.Equal(t, "NOSECRET", creds.SecretAccessKey, "Expect secret access key to match")
	assert.Empty(t, creds.SessionToken, "Expect session token to be empty")

}

func TestChainProviderGet(t *testing.T) {
	p := &ChainProvider{
		Providers: []CredentialsProvider{
			&stubProvider{err: awserr.New("FirstError", "first provider error", nil)},
			&stubProvider{err: awserr.New("SecondError", "second provider error", nil)},
			&stubProvider{
				creds: Credentials{
					AccessKeyID:     "AKID",
					SecretAccessKey: "SECRET",
					SessionToken:    "",
				},
			},
		},
	}

	creds, err := p.Retrieve()
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, "AKID", creds.AccessKeyID, "Expect access key ID to match")
	assert.Equal(t, "SECRET", creds.SecretAccessKey, "Expect secret access key to match")
	assert.Empty(t, creds.SessionToken, "Expect session token to be empty")
}

func TestChainProviderIsExpired(t *testing.T) {
	stubProvider := &stubProvider{expired: true}
	p := &ChainProvider{
		Providers: []CredentialsProvider{
			stubProvider,
		},
	}

	assert.True(t, p.IsExpired(), "Expect expired to be true before any Retrieve")
	_, err := p.Retrieve()
	assert.Nil(t, err, "Expect no error")
	assert.False(t, p.IsExpired(), "Expect not expired after retrieve")

	stubProvider.expired = true
	assert.True(t, p.IsExpired(), "Expect return of expired provider")

	_, err = p.Retrieve()
	assert.False(t, p.IsExpired(), "Expect not expired after retrieve")
}

func TestChainProviderWithNoProvider(t *testing.T) {
	p := &ChainProvider{
		Providers: []CredentialsProvider{},
	}

	assert.True(t, p.IsExpired(), "Expect expired with no providers")
	_, err := p.Retrieve()
	if e, a := "no valid providers", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q error in %q", e, a)
	}
}

func TestChainProviderWithNoValidProvider(t *testing.T) {
	errs := []error{
		awserr.New("FirstError", "first provider error", nil),
		awserr.New("SecondError", "second provider error", nil),
	}
	p := &ChainProvider{
		Providers: []CredentialsProvider{
			&stubProvider{err: errs[0]},
			&stubProvider{err: errs[1]},
		},
	}

	assert.True(t, p.IsExpired(), "Expect expired with no providers")
	_, err := p.Retrieve()
	if e, a := "no valid providers", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q error in %q", e, a)
	}
}
