package kitchensinktest

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

func TestCredentialRefreshRetry(t *testing.T) {
	sdk.NowTime = func() time.Time { return time.Unix(0, 0) }
	defer func() {
		sdk.NowTime = time.Now
	}()

	credvalues := []aws.Credentials{
		{AccessKeyID: "foo", SecretAccessKey: "bar"},
		{AccessKeyID: "baz", SecretAccessKey: "qux"},
	}
	mhttp := &mockHTTP{
		resps: []*http.Response{
			{StatusCode: 500, Body: http.NoBody},
			{StatusCode: 200, Body: http.NoBody},
		},
	}

	svc := New(Options{
		Region:     "us-east-1",
		HTTPClient: mhttp,
		Credentials: &mockCredentials{
			credentials: credvalues,
		},
	})

	_, err := svc.GetItem(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(mhttp.reqs) != 2 {
		t.Fatalf("there should have been 2 requests, but there were %d", len(mhttp.reqs))
	}

	// verify that creds actually made it into the signature
	req0 := mhttp.reqs[0]
	auth0 := req0.Header.Get("Authorization")
	if !strings.Contains(auth0, "Credential=foo/19700101/us-east-1/awsjson1kitchensink/aws4_request") {
		t.Errorf("1st request should have been AKID=foo, signature was %q", auth0)
	}

	req1 := mhttp.reqs[1]
	auth1 := req1.Header.Get("Authorization")
	if !strings.Contains(auth1, "Credential=baz/19700101/us-east-1/awsjson1kitchensink/aws4_request") {
		t.Errorf("2nd request should have been AKID=baz, signature was %q", auth1)
	}
}
