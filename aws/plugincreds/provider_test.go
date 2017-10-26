// +build go1.8,awsinclude

package plugincreds

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

func TestProvider_Passthrough(t *testing.T) {
	p := Provider{}
	p.RetrieveFn = buildRetrieveFn(
		func() (k, s, t string, ext time.Time, err error) {
			return "key", "secret", "session",
				sdk.NowTime().Add(2 * time.Hour),
				nil
		},
	)

	creds, err := p.Retrieve()
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "key", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "secret", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "session", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := true, creds.CanExpire; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if creds.Expired() {
		t.Errorf("expect not to be expired")
	}
}

func TestProvider_Error(t *testing.T) {
	expectErr := fmt.Errorf("expect error")

	p := Provider{}
	p.RetrieveFn = buildRetrieveFn(
		func() (k, s, t string, ext time.Time, err error) {
			return "", "", "",
				time.Time{},
				expectErr
		},
	)

	creds, err := p.Retrieve()
	if err == nil {
		t.Fatalf("expect error, got none")
	}

	aerr := err.(awserr.Error)
	if e, a := ErrCodePluginProviderRetrieve, aerr.Code(); e != a {
		t.Errorf("expect %s error code, got %s", e, a)
	}

	if e, a := expectErr, aerr.OrigErr(); e != a {
		t.Errorf("expect %v cause error, got %v", e, a)
	}

	expect := aws.Credentials{}
	if expect != creds {
		t.Errorf("expect %+v credentials, got %+v", expect, creds)
	}
}
