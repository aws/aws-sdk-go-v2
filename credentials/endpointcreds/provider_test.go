package endpointcreds_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go"
)

type mockClient func(*http.Request) (*http.Response, error)

func (m mockClient) Do(r *http.Request) (*http.Response, error) {
	return m(r)
}

func TestRetrieveRefreshableCredentials(t *testing.T) {
	orig := sdk.NowTime
	defer func() { sdk.NowTime = orig }()

	p := endpointcreds.New("http://127.0.0.1", func(o *endpointcreds.Options) {
		o.HTTPClient = mockClient(func(r *http.Request) (*http.Response, error) {
			expTime := time.Now().UTC().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z")

			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf(`{
  "AccessKeyID": "AKID",
  "SecretAccessKey": "SECRET",
  "Token": "TOKEN",
  "Expiration": "%s"
}`, expTime)))),
			}, nil
		})
	})
	creds, err := p.Retrieve(context.Background())

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "TOKEN", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if creds.Expired() {
		t.Errorf("expect not expired")
	}

	sdk.NowTime = func() time.Time {
		return time.Now().Add(2 * time.Hour)
	}
	if !creds.Expired() {
		t.Errorf("expect to be expired")
	}
}

func TestRetrieveStaticCredentials(t *testing.T) {
	orig := sdk.NowTime
	defer func() { sdk.NowTime = orig }()

	p := endpointcreds.New("http://127.0.0.1", func(o *endpointcreds.Options) {
		o.HTTPClient = mockClient(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
  "AccessKeyID": "AKID",
  "SecretAccessKey": "SECRET",
  "AccountID": "012345678901"
}`))),
			}, nil
		})
	})
	creds, err := p.Retrieve(context.Background())

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "012345678901", creds.AccountID; e != a {
		t.Errorf("expect account ID to be %v, got %v", e, a)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}

	sdk.NowTime = func() time.Time {
		return time.Date(3000, 12, 16, 1, 30, 37, 0, time.UTC)
	}

	if creds.Expired() {
		t.Errorf("expect not to be expired")
	}
}

func TestAuthTokenProvider(t *testing.T) {
	cases := map[string]struct {
		AuthToken         string
		AuthTokenProvider endpointcreds.AuthTokenProvider
		ExpectAuthToken   string
		ExpectError       bool
	}{
		"AuthToken": {
			AuthToken:       "Basic abc123",
			ExpectAuthToken: "Basic abc123",
		},
		"AuthFileToken": {
			AuthToken: "Basic abc123",
			AuthTokenProvider: endpointcreds.TokenProviderFunc(func() (string, error) {
				return "Hello %20world", nil
			}),
			ExpectAuthToken: "Hello %20world",
		},
		"RetrieveFileTokenError": {
			AuthToken: "Basic abc123",
			AuthTokenProvider: endpointcreds.TokenProviderFunc(func() (string, error) {
				return "", fmt.Errorf("test error")
			}),
			ExpectAuthToken: "Hello %20world",
			ExpectError:     true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			orig := sdk.NowTime
			defer func() { sdk.NowTime = orig }()

			var actualToken string
			p := endpointcreds.New("http://127.0.0.1", func(o *endpointcreds.Options) {
				o.HTTPClient = mockClient(func(r *http.Request) (*http.Response, error) {
					actualToken = r.Header["Authorization"][0]
					return &http.Response{
						StatusCode: 200,
						Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
  "AccessKeyID": "AKID",
  "SecretAccessKey": "SECRET"
}`))),
					}, nil
				})
				o.AuthorizationToken = c.AuthToken
				o.AuthorizationTokenProvider = c.AuthTokenProvider
			})
			creds, err := p.Retrieve(context.Background())

			if err != nil && !c.ExpectError {
				t.Errorf("expect no error, got %v", err)
			} else if err == nil && c.ExpectError {
				t.Errorf("expect error, got nil")
			}

			if c.ExpectError {
				return
			}

			if e, a := "AKID", creds.AccessKeyID; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if e, a := "SECRET", creds.SecretAccessKey; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if v := creds.SessionToken; len(v) != 0 {
				t.Errorf("expect empty, got %v", v)
			}
			if e, a := c.ExpectAuthToken, actualToken; e != a {
				t.Errorf("Expect %v, got %v", e, a)
			}

			sdk.NowTime = func() time.Time {
				return time.Date(3000, 12, 16, 1, 30, 37, 0, time.UTC)
			}

			if creds.Expired() {
				t.Errorf("expect not to be expired")
			}
		})
	}
}

func TestFailedRetrieveCredentials(t *testing.T) {
	p := endpointcreds.New("http://127.0.0.1", func(o *endpointcreds.Options) {
		o.HTTPClient = mockClient(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
  "code": "Error",
  "message": "Message"
}`))),
				Header: http.Header{
					"Content-Type": {"application/json"},
				},
			}, nil
		})
	})
	creds, err := p.Retrieve(context.Background())

	if err == nil {
		t.Fatalf("expect error, got none")
	}

	if e, a := "failed to load credentials", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v, got %v", e, a)
	}

	var apiError smithy.APIError
	if !errors.As(err, &apiError) {
		t.Fatalf("expect %T error, got %v", apiError, err)
	}
	if e, a := "Error", apiError.ErrorCode(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "Message", apiError.ErrorMessage(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	if v := creds.AccessKeyID; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if v := creds.SecretAccessKey; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if creds.Expired() {
		t.Errorf("expect empty creds not to be expired")
	}
}

type mockClientN struct {
	responses []*http.Response
	index     int
}

func (c *mockClientN) Do(r *http.Request) (*http.Response, error) {
	resp := c.responses[c.index]
	c.index++
	return resp, nil
}

func TestRetryHTTPStatusCode(t *testing.T) {
	expTime := time.Now().UTC().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z")
	credsResp := fmt.Sprintf(`{"AccessKeyID":"AKID","SecretAccessKey":"SECRET","Token":"TOKEN","Expiration":"%s"}`, expTime)

	p := endpointcreds.New("http://127.0.0.1", func(o *endpointcreds.Options) {
		o.HTTPClient = &mockClientN{
			responses: []*http.Response{
				{
					StatusCode: 429,
					Body:       io.NopCloser(strings.NewReader("You have made too many requests.")),
					Header: http.Header{
						"Content-Type": {"text/plain"},
					},
				},
				{
					StatusCode: 500,
					Body:       io.NopCloser(strings.NewReader("Internal server error.")),
					Header: http.Header{
						"Content-Type": {"text/plain"},
					},
				},
				{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(credsResp)),
					Header: http.Header{
						"Content-Type": {"application/json"},
					},
				},
			},
		}
	})

	creds, err := p.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "TOKEN", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if creds.Expired() {
		t.Errorf("expect not expired")
	}

	sdk.NowTime = func() time.Time {
		return time.Now().Add(2 * time.Hour)
	}
	if !creds.Expired() {
		t.Errorf("expect to be expired")
	}
}
