package defaults_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestValidateEndpointHandler(t *testing.T) {
	os.Clearenv()

	cfg := unit.Config()
	cfg.Region = "us-west-2"

	svc := awstesting.NewClient(cfg)
	svc.Handlers.Clear()
	svc.Handlers.Validate.PushBackNamed(defaults.ValidateEndpointHandler)

	req := svc.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
	err := req.Build()

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestValidateEndpointHandlerErrorRegion(t *testing.T) {
	os.Clearenv()

	cfg := unit.Config()
	cfg.Region = ""

	svc := awstesting.NewClient(cfg)
	svc.Handlers.Clear()
	svc.Handlers.Validate.PushBackNamed(defaults.ValidateEndpointHandler)

	req := svc.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
	err := req.Build()

	if err == nil {
		t.Errorf("expect error, got none")
	}
	if e, a := aws.ErrMissingRegion, err; e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
}

type mockCredsProvider struct {
	retrieveCalled   bool
	invalidateCalled bool
}

func (m *mockCredsProvider) Retrieve() (aws.Credentials, error) {
	m.retrieveCalled = true
	return aws.Credentials{Source: "mockCredsProvider"}, nil
}

func (m *mockCredsProvider) Invalidate() {
	m.invalidateCalled = true
}

func TestAfterRetry_RefreshCreds(t *testing.T) {
	orig := sdk.SleepWithContext
	defer func() { sdk.SleepWithContext = orig }()
	sdk.SleepWithContext = func(context.Context, time.Duration) error { return nil }

	credProvider := &mockCredsProvider{}

	cfg := unit.Config()
	cfg.Credentials = credProvider

	svc := awstesting.NewClient(cfg)
	req := svc.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
	req.Retryable = aws.Bool(true)
	req.Error = awserr.New("ExpiredTokenException", "", nil)
	req.HTTPResponse = &http.Response{
		StatusCode: 403,
	}

	defaults.AfterRetryHandler.Fn(req)

	if !credProvider.invalidateCalled {
		t.Errorf("expect credentials to be invalidated")
	}
}

func TestAfterRetry_NoPanicRefreshStaticCreds(t *testing.T) {
	orig := sdk.SleepWithContext
	defer func() { sdk.SleepWithContext = orig }()
	sdk.SleepWithContext = func(context.Context, time.Duration) error { return nil }

	credProvider := aws.NewStaticCredentialsProvider("AKID", "SECRET", "")

	cfg := unit.Config()
	cfg.Credentials = credProvider

	svc := awstesting.NewClient(cfg)
	req := svc.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
	req.Retryable = aws.Bool(true)
	req.Error = awserr.New("ExpiredTokenException", "", nil)
	req.HTTPResponse = &http.Response{
		StatusCode: 403,
	}

	defaults.AfterRetryHandler.Fn(req)
}

func TestAfterRetryWithContextCanceled(t *testing.T) {
	c := awstesting.NewClient(unit.Config())

	req := c.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{}, 0)}
	req.SetContext(ctx)

	req.Error = fmt.Errorf("some error")
	req.Retryable = aws.Bool(true)
	req.HTTPResponse = &http.Response{
		StatusCode: 500,
	}

	close(ctx.DoneCh)
	ctx.Error = fmt.Errorf("context canceled")

	defaults.AfterRetryHandler.Fn(req)

	if req.Error == nil {
		t.Fatalf("expect error but didn't receive one")
	}

	aerr := req.Error.(awserr.Error)

	if e, a := aws.ErrCodeRequestCanceled, aerr.Code(); e != a {
		t.Errorf("expect %q, error code got %q", e, a)
	}
}

func TestAfterRetryWithContext(t *testing.T) {
	c := awstesting.NewClient(unit.Config())

	req := c.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{}, 0)}
	req.SetContext(ctx)

	req.Error = fmt.Errorf("some error")
	req.Retryable = aws.Bool(true)
	req.HTTPResponse = &http.Response{
		StatusCode: 500,
	}

	defaults.AfterRetryHandler.Fn(req)

	if req.Error != nil {
		t.Fatalf("expect no error, got %v", req.Error)
	}
	if e, a := 1, req.RetryCount; e != a {
		t.Errorf("expect retry count to be %d, got %d", e, a)
	}
}

func TestSendWithContextCanceled(t *testing.T) {
	c := awstesting.NewClient(unit.Config())

	req := c.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{}, 0)}
	req.SetContext(ctx)

	req.Error = fmt.Errorf("some error")
	req.Retryable = aws.Bool(true)
	req.HTTPResponse = &http.Response{
		StatusCode: 500,
	}

	close(ctx.DoneCh)
	ctx.Error = fmt.Errorf("context canceled")

	defaults.SendHandler.Fn(req)

	if req.Error == nil {
		t.Fatalf("expect error but didn't receive one")
	}

	aerr := req.Error.(awserr.Error)

	if e, a := aws.ErrCodeRequestCanceled, aerr.Code(); e != a {
		t.Errorf("expect %q, error code got %q", e, a)
	}
}

type testSendHandlerTransport struct{}

func (t *testSendHandlerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("mock error")
}

func TestSendHandlerError(t *testing.T) {
	cfg := unit.Config()
	cfg.HTTPClient = &http.Client{
		Transport: &testSendHandlerTransport{},
	}
	svc := awstesting.NewClient(cfg)
	svc.Handlers.Clear()
	svc.Handlers.Send.PushBackNamed(defaults.SendHandler)
	r := svc.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)

	r.Send()

	if r.Error == nil {
		t.Errorf("expect error, got none")
	}
	if r.HTTPResponse == nil {
		t.Errorf("expect response, got none")
	}
}

func TestSendWithoutFollowRedirects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/original":
			w.Header().Set("Location", "/redirected")
			w.WriteHeader(301)
		case "/redirected":
			t.Fatalf("expect not to redirect, but was")
		}
	}))

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := awstesting.NewClient(cfg)
	svc.Handlers.Clear()
	svc.Handlers.Send.PushBackNamed(defaults.SendHandler)

	r := svc.NewRequest(&aws.Operation{
		Name:     "Operation",
		HTTPPath: "/original",
	}, nil, nil)
	r.DisableFollowRedirects = true

	err := r.Send()
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := 301, r.HTTPResponse.StatusCode; e != a {
		t.Errorf("expect %d status code, got %d", e, a)
	}
}

func TestValidateReqSigHandler(t *testing.T) {
	cases := []struct {
		Req    *aws.Request
		Resign bool
	}{
		{
			Req: &aws.Request{
				Config: aws.Config{Credentials: aws.AnonymousCredentials},
				Time:   time.Now().Add(-15 * time.Minute),
			},
			Resign: false,
		},
		{
			Req: &aws.Request{
				Time: time.Now().Add(-15 * time.Minute),
			},
			Resign: true,
		},
		{
			Req: &aws.Request{
				Time: time.Now().Add(-1 * time.Minute),
			},
			Resign: false,
		},
	}

	for i, c := range cases {
		resigned := false
		c.Req.Handlers.Sign.PushBack(func(r *aws.Request) {
			resigned = true
		})

		defaults.ValidateReqSigHandler.Fn(c.Req)

		if c.Req.Error != nil {
			t.Errorf("expect no error, got %v", c.Req.Error)
		}
		if e, a := c.Resign, resigned; e != a {
			t.Errorf("%d, expect %v to be %v", i, e, a)
		}
	}
}

func setupContentLengthTestServer(t *testing.T, hasContentLength bool, contentLength int64) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Header["Content-Length"]
		if e, a := hasContentLength, ok; e != a {
			t.Errorf("expect %v to be %v", e, a)
		}
		if hasContentLength {
			if e, a := contentLength, r.ContentLength; e != a {
				t.Errorf("expect %v to be %v", e, a)
			}
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
		r.Body.Close()

		authHeader := r.Header.Get("Authorization")
		if hasContentLength {
			if e, a := "content-length", authHeader; !strings.Contains(a, e) {
				t.Errorf("expect %v to be in %v", e, a)
			}
		} else {
			if e, a := "content-length", authHeader; strings.Contains(a, e) {
				t.Errorf("expect %v to not be in %v", e, a)
			}
		}

		if e, a := contentLength, int64(len(b)); e != a {
			t.Errorf("expect %v to be %v", e, a)
		}
	}))

	return server
}

func TestBuildContentLength_ZeroBody(t *testing.T) {
	server := setupContentLengthTestServer(t, false, 0)

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := s3.New(cfg)
	svc.ForcePathStyle = true
	req := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("bucketname"),
		Key:    aws.String("keyname"),
	})

	if _, err := req.Send(); err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestBuildContentLength_NegativeBody(t *testing.T) {
	server := setupContentLengthTestServer(t, false, 0)

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := s3.New(cfg)
	svc.ForcePathStyle = true
	req := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("bucketname"),
		Key:    aws.String("keyname"),
	})

	req.HTTPRequest.Header.Set("Content-Length", "-1")

	if _, err := req.Send(); err != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}
}

func TestBuildContentLength_WithBody(t *testing.T) {
	server := setupContentLengthTestServer(t, true, 1024)

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := s3.New(cfg)
	svc.ForcePathStyle = true
	req := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("bucketname"),
		Key:    aws.String("keyname"),
		Body:   bytes.NewReader(make([]byte, 1024)),
	})

	if _, err := req.Send(); err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
