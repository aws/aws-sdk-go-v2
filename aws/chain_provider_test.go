package aws

import (
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
)

func TestChainProvider_WithNames(t *testing.T) {
	p := NewChainProvider(
		[]CredentialsProvider{
			&stubProvider{err: awserr.New("FirstError", "first provider error", nil)},
			&stubProvider{err: awserr.New("SecondError", "second provider error", nil)},
			&stubProvider{
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
	)

	creds, err := p.Retrieve()
	if err != nil {
		t.Fatalf("expect no error")
	}
	if e, a := "stubProvider", creds.Source; e != a {
		t.Errorf("expect provider name to match")
	}

	// Also check credentials
	if e, a := "AKIF", creds.AccessKeyID; e != a {
		t.Errorf("expect access key ID to match")
	}
	if e, a := "NOSECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect secret access key to match")
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect session token to be empty")
	}
}

func TestChainProvider_Retrieve(t *testing.T) {
	p := NewChainProvider(
		[]CredentialsProvider{
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
	)

	creds, err := p.Retrieve()
	if err != nil {
		t.Fatalf("expect no error")
	}
	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect access key ID to match")
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect secret access key to match")
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect session token to be empty")
	}
}

func TestChainProvider_IsExpired(t *testing.T) {
	p := NewChainProvider(
		[]CredentialsProvider{
			&stubProvider{expires: time.Now().Add(-5 * time.Minute)},
		},
	)

	creds, err := p.Retrieve()
	if err != nil {
		t.Fatalf("expect no error")
	}
	if !creds.Expired() {
		t.Errorf("expect expired, %v", creds)
	}
}

func TestChainProvider_WithNoProvider(t *testing.T) {
	p := NewChainProvider([]CredentialsProvider{})

	_, err := p.Retrieve()
	if e, a := "no valid providers", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q error in %q", e, a)
	}
}

func TestChainProvider_WithNoValidProvider(t *testing.T) {
	errs := []error{
		awserr.New("FirstError", "first provider error", nil),
		awserr.New("SecondError", "second provider error", nil),
	}
	p := NewChainProvider(
		[]CredentialsProvider{
			&stubProvider{err: errs[0]},
			&stubProvider{err: errs[1]},
		},
	)

	_, err := p.Retrieve()
	if e, a := "no valid providers", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q error in %q", e, a)
	}
}
