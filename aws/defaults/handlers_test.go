package defaults_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
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
	var expected *aws.MissingRegionError
	if !errors.As(err, &expected) {
		t.Fatalf("expected %T, got %T", expected, err)
	}
}

func TestShouldRetry_WithContext(t *testing.T) {
	c := awstesting.NewClient(unit.Config())
	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{}, 0)}

	req := c.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
	req.SetContext(ctx)

	req.Error = &aws.HTTPResponseError{
		Response: &http.Response{
			StatusCode: 500,
			Header:     http.Header{},
		},
	}

	defaults.RetryableCheckHandler.Fn(req)

	if req.RetryDelay == 0 {
		t.Fatalf("expect retry delay got none")
	}
}

func TestSendWithContextCanceled(t *testing.T) {
	c := awstesting.NewClient(unit.Config())

	req := c.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{}, 0)}
	req.SetContext(ctx)

	close(ctx.DoneCh)
	ctx.Error = fmt.Errorf("context canceled")

	defaults.SendHandler.Fn(req)

	if req.Error == nil {
		t.Fatalf("expect error but didn't receive one")
	}

	var aerr *aws.RequestCanceledError
	if !errors.As(req.Error, &aerr) {
		t.Fatalf("expect %T error, got %v", aerr, req.Error)
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

	if _, err := req.Send(context.Background()); err != nil {
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

	if _, err := req.Send(context.Background()); err != nil {
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

	if _, err := req.Send(context.Background()); err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestSendHandler_HEADNoBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cfg := defaults.Config()
	cfg.Region = "mock-region"
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := s3.New(cfg)
	svc.ForcePathStyle = true

	req := svc.HeadObjectRequest(&s3.HeadObjectInput{
		Bucket: aws.String("bucketname"),
		Key:    aws.String("keyname"),
	})

	if e, a := http.NoBody, req.HTTPRequest.Body; e != a {
		t.Fatalf("expect %T request body, got %T", e, a)
	}

	_, err := req.Send(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := http.StatusOK, req.HTTPResponse.StatusCode; e != a {
		t.Errorf("expect %d status code, got %d", e, a)
	}
}

func TestRequestInvocationIDHeaderHandler(t *testing.T) {
	cfg := unit.Config()
	cfg.Handlers.Send.Clear()

	var invokeID string
	cfg.Handlers.Build.PushBack(func(r *aws.Request) {
		invokeID = r.InvocationID
		if len(invokeID) == 0 {
			t.Fatalf("expect non-empty invocation id")
		}
	})
	cfg.Handlers.Send.PushBack(func(r *aws.Request) {
		if e, a := invokeID, r.InvocationID; e != a {
			t.Errorf("expect %v invoke ID, got %v", e, a)
		}
		r.Error = &aws.RequestSendError{Err: io.ErrUnexpectedEOF}
	})
	retryer := retry.NewStandard(func(o *retry.StandardOptions) {
		o.MaxAttempts = 3
	})
	r := aws.New(cfg, aws.Metadata{}, cfg.Handlers, retryer, &aws.Operation{},
		&struct{}{}, struct{}{})

	if len(r.InvocationID) == 0 {
		t.Fatalf("expect invocation id, got none")
	}

	err := r.Send()
	if err == nil {
		t.Fatalf("expect error got on")
	}
	var maxErr *aws.MaxAttemptsError
	if !errors.As(err, &maxErr) {
		t.Fatalf("expect max errors, got %v", err)
	} else {
		if e, a := 3, maxErr.Attempt; e != a {
			t.Errorf("expect %v attempts, got %v", e, a)
		}
	}
	if len(invokeID) == 0 {
		t.Fatalf("expect non-empty invocation id")
	}

	if e, a := r.InvocationID, r.HTTPRequest.Header.Get("amz-sdk-invocation-id"); e != a {
		t.Errorf("expect %v invocation id, got %v", e, a)
	}
}

func TestRetryMetricHeaderHandler(t *testing.T) {
	nowTime := sdk.NowTime
	defer func() {
		sdk.NowTime = nowTime
	}()
	sdk.NowTime = func() time.Time {
		return time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC)
	}

	cases := map[string]struct {
		Attempt           int
		MaxAttempts       int
		Client            aws.HTTPClient
		ContextDeadline   time.Time
		AttemptClockSkews []time.Duration
		Expect            string
	}{
		"first attempt": {
			Attempt: 1, MaxAttempts: 3,
			Expect: "attempt=1; max=3",
		},
		"last attempt": {
			Attempt: 3, MaxAttempts: 3,
			Expect: "attempt=3; max=3",
		},
		"no max attempt": {
			Attempt: 10,
			Expect:  "attempt=10",
		},
		"with ttl client timeout": {
			Attempt: 2, MaxAttempts: 3,
			AttemptClockSkews: []time.Duration{
				10 * time.Second,
			},
			Client: func() aws.HTTPClient {
				c := &aws.BuildableHTTPClient{}
				return c.WithTimeout(10 * time.Second)
			}(),
			Expect: "attempt=2; max=3; ttl=20200202T000020Z",
		},
		"with ttl context deadline": {
			Attempt: 1, MaxAttempts: 3,
			AttemptClockSkews: []time.Duration{
				10 * time.Second,
			},
			ContextDeadline: sdk.NowTime().Add(10 * time.Second),
			Expect:          "attempt=1; max=3; ttl=20200202T000020Z",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			cfg := unit.Config()
			if c.Client != nil {
				cfg.HTTPClient = c.Client
			}
			r := aws.New(cfg, aws.Metadata{}, cfg.Handlers, aws.NoOpRetryer{},
				&aws.Operation{}, &struct{}{}, struct{}{})
			if !c.ContextDeadline.IsZero() {
				ctx, cancel := context.WithDeadline(r.Context(), c.ContextDeadline)
				defer cancel()
				r.SetContext(ctx)
			}

			r.AttemptNum = c.Attempt
			r.AttemptClockSkews = c.AttemptClockSkews
			r.Retryer = retry.AddWithMaxAttempts(r.Retryer, c.MaxAttempts)

			defaults.RetryMetricHeaderHandler.Fn(r)
			if r.Error != nil {
				t.Fatalf("expect no error, got %v", r.Error)
			}

			if e, a := c.Expect, r.HTTPRequest.Header.Get("amz-sdk-request"); e != a {
				t.Errorf("expect %q metric, got %q", e, a)
			}
		})
	}
}

func TestAttemptClockSkewHandler(t *testing.T) {
	cases := map[string]struct {
		Req    *aws.Request
		Expect []time.Duration
	}{
		"no response": {Req: &aws.Request{}},
		"failed response": {
			Req: &aws.Request{
				HTTPResponse: &http.Response{StatusCode: 0, Header: http.Header{}},
			},
		},
		"no date header response": {
			Req: &aws.Request{
				HTTPResponse: &http.Response{StatusCode: 200, Header: http.Header{}},
			},
		},
		"invalid date header response": {
			Req: &aws.Request{
				HTTPResponse: &http.Response{
					StatusCode: 200,
					Header:     http.Header{"Date": []string{"abc123"}},
				},
			},
		},
		"response at unset": {
			Req: &aws.Request{
				HTTPResponse: &http.Response{
					StatusCode: 200,
					Header: http.Header{
						"Date": []string{"Thu, 05 Mar 2020 22:25:15 GMT"},
					},
				},
			},
		},
		"first date response": {
			Req: &aws.Request{
				HTTPResponse: &http.Response{
					StatusCode: 200,
					Header: http.Header{
						"Date": []string{"Thu, 05 Mar 2020 22:25:15 GMT"},
					},
				},
				ResponseAt: time.Date(2020, 3, 5, 22, 25, 17, 0, time.UTC),
			},
			Expect: []time.Duration{-2 * time.Second},
		},
		"subsequent date response": {
			Req: &aws.Request{
				HTTPResponse: &http.Response{
					StatusCode: 200,
					Header: http.Header{
						"Date": []string{"Thu, 05 Mar 2020 22:25:15 GMT"},
					},
				},
				ResponseAt: time.Date(2020, 3, 5, 22, 25, 14, 0, time.UTC),
				AttemptClockSkews: []time.Duration{
					-2 * time.Second,
				},
			},
			Expect: []time.Duration{
				-2 * time.Second,
				1 * time.Second,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := new(aws.Request)
			*r = *c.Req

			defaults.AttemptClockSkewHandler.Fn(r)
			if err := r.Error; err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := len(c.Expect), len(r.AttemptClockSkews); e != a {
				t.Errorf("expect %v skews, got %v", e, a)
			}

			for i, s := range r.AttemptClockSkews {
				if e, a := c.Expect[i], s; e != a {
					t.Errorf("%d, expect %v skew, got %v", i, e, a)
				}
			}
		})
	}
}
