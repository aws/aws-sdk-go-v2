package ec2rolecreds

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/awslabs/smithy-go"
)

const credsRespTmpl = `{
  "Code": "Success",
  "Type": "AWS-HMAC",
  "AccessKeyId" : "accessKey",
  "SecretAccessKey" : "secret",
  "Token" : "token",
  "Expiration" : "%s",
  "LastUpdated" : "2009-11-23T0:00:00Z"
}`

const credsFailRespTmpl = `{
  "Code": "ErrorCode",
  "Message": "ErrorMsg",
  "LastUpdated": "2009-11-23T0:00:00Z"
}`

type mockClient struct {
	t          *testing.T
	roleName   string
	failAssume bool
	expireOn   string
}

func (c mockClient) GetMetadata(
	ctx context.Context, params *imds.GetMetadataInput, optFns ...func(*imds.Options),
) (
	*imds.GetMetadataOutput, error,
) {
	switch params.Path {
	case iamSecurityCredsPath:
		return &imds.GetMetadataOutput{
			Content: ioutil.NopCloser(strings.NewReader(c.roleName)),
		}, nil

	case iamSecurityCredsPath + c.roleName:
		var w strings.Builder
		if c.failAssume {
			fmt.Fprintf(&w, credsFailRespTmpl)
		} else {
			fmt.Fprintf(&w, credsRespTmpl, c.expireOn)
		}
		return &imds.GetMetadataOutput{
			Content: ioutil.NopCloser(strings.NewReader(w.String())),
		}, nil
	default:
		return nil, fmt.Errorf("unexpected path, %v", params.Path)
	}
}

func TestProvider(t *testing.T) {
	orig := sdk.NowTime
	defer func() { sdk.NowTime = orig }()

	p := New(func(options *Options) {
		options.Client = mockClient{
			roleName:   "RoleName",
			failAssume: false,
			expireOn:   "2014-12-16T01:51:37Z",
		}
	})

	creds, err := p.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "accessKey", creds.AccessKeyID; e != a {
		t.Errorf("Expect access key ID to match")
	}
	if e, a := "secret", creds.SecretAccessKey; e != a {
		t.Errorf("Expect secret access key to match")
	}
	if e, a := "token", creds.SessionToken; e != a {
		t.Errorf("Expect session token to match")
	}

	sdk.NowTime = func() time.Time {
		return time.Date(2014, 12, 16, 0, 55, 37, 0, time.UTC)
	}

	if creds.Expired() {
		t.Errorf("Expect not expired")
	}
}

func TestProvider_FailAssume(t *testing.T) {
	p := New(func(options *Options) {
		options.Client = mockClient{
			roleName:   "RoleName",
			failAssume: true,
			expireOn:   "2014-12-16T01:51:37Z",
		}
	})

	creds, err := p.Retrieve(context.Background())
	if err == nil {
		t.Fatalf("expect error, got none")
	}

	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expect %T error, got %v", apiErr, err)
	}
	if e, a := "ErrorCode", apiErr.ErrorCode(); e != a {
		t.Errorf("expect %v code, got %v", e, a)
	}
	if e, a := "ErrorMsg", apiErr.ErrorMessage(); e != a {
		t.Errorf("expect %v message, got %v", e, a)
	}

	nestedErr := errors.Unwrap(apiErr)
	if nestedErr != nil {
		t.Fatalf("expect no nested error, got %v", err)
	}

	if e, a := "", creds.AccessKeyID; e != a {
		t.Errorf("Expect access key ID to match")
	}
	if e, a := "", creds.SecretAccessKey; e != a {
		t.Errorf("Expect secret access key to match")
	}
	if e, a := "", creds.SessionToken; e != a {
		t.Errorf("Expect session token to match")
	}
}

func TestProvider_IsExpired(t *testing.T) {
	orig := sdk.NowTime
	defer func() { sdk.NowTime = orig }()

	p := New(func(options *Options) {
		options.Client = mockClient{
			roleName:   "RoleName",
			failAssume: false,
			expireOn:   "2014-12-16T01:51:37Z",
		}
	})

	sdk.NowTime = func() time.Time {
		return time.Date(2014, 12, 16, 0, 55, 37, 0, time.UTC)
	}

	creds, err := p.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if creds.Expired() {
		t.Errorf("expect not to be expired")
	}

	sdk.NowTime = func() time.Time {
		return time.Date(2014, 12, 16, 1, 55, 37, 0, time.UTC)
	}

	if !creds.Expired() {
		t.Errorf("expect to be expired")
	}
}
