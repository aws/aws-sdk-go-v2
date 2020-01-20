package aws_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/private/protocol/jsonrpc"
	"github.com/aws/aws-sdk-go-v2/private/protocol/rest"
)

type testData struct {
	Data string
}

func body(str string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(str)))
}

func unmarshal(req *aws.Request) {
	defer req.HTTPResponse.Body.Close()
	if req.Data != nil {
		json.NewDecoder(req.HTTPResponse.Body).Decode(req.Data)
	}
	return
}

func unmarshalError(req *aws.Request) {
	bodyBytes, err := ioutil.ReadAll(req.HTTPResponse.Body)
	if err != nil {
		req.Error = awserr.New("UnmarshaleError", req.HTTPResponse.Status, err)
		return
	}
	if len(bodyBytes) == 0 {
		req.Error = awserr.NewRequestFailure(
			awserr.New("UnmarshaleError", req.HTTPResponse.Status, fmt.Errorf("empty body")),
			req.HTTPResponse.StatusCode,
			"",
		)
		return
	}
	var jsonErr jsonErrorResponse
	if err := json.Unmarshal(bodyBytes, &jsonErr); err != nil {
		req.Error = awserr.New("UnmarshaleError", "JSON unmarshal", err)
		return
	}
	req.Error = awserr.NewRequestFailure(
		awserr.New(jsonErr.Code, jsonErr.Message, nil),
		req.HTTPResponse.StatusCode,
		"",
	)
}

type jsonErrorResponse struct {
	Code    string `json:"__type"`
	Message string `json:"message"`
}

// test that retries occur for 5xx status codes
func TestRequestRecoverRetry5xx(t *testing.T) {
	restoreSleep := mockSleep()
	defer restoreSleep()

	reqNum := 0
	reqs := []http.Response{
		{StatusCode: 500, Body: body(`{"__type":"UnknownError","message":"An error occurred."}`)},
		{StatusCode: 502, Body: body(`{"__type":"UnknownError","message":"An error occurred."}`)},
		{StatusCode: 200, Body: body(`{"data":"valid"}`)},
	}

	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 10
	})

	s := awstesting.NewClient(cfg)
	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &reqs[reqNum]
		reqNum++
	})
	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	err := r.Send()
	if err != nil {
		t.Fatalf("expect no error, but got %v", err)
	}
	if e, a := 2, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}
	if e, a := "valid", out.Data; e != a {
		t.Errorf("expect %q output got %q", e, a)
	}
}

// test that retries occur for 4xx status codes with a response type that can be retried - see `shouldRetry`
func TestRequestRecoverRetry4xxRetryable(t *testing.T) {
	restoreSleep := mockSleep()
	defer restoreSleep()

	reqNum := 0
	reqs := []http.Response{
		{StatusCode: 400, Body: body(`{"__type":"Throttling","message":"Rate exceeded."}`)},
		{StatusCode: 400, Body: body(`{"__type":"ProvisionedThroughputExceededException","message":"Rate exceeded."}`)},
		{StatusCode: 429, Body: body(`{"__type":"FooException","message":"Rate exceeded."}`)},
		{StatusCode: 200, Body: body(`{"data":"valid"}`)},
	}

	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 10
	})

	s := awstesting.NewClient(cfg)
	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &reqs[reqNum]
		reqNum++
	})
	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	err := r.Send()
	if err != nil {
		t.Fatalf("expect no error, but got %v", err)
	}
	if e, a := 3, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}
	if e, a := "valid", out.Data; e != a {
		t.Errorf("expect %q output got %q", e, a)
	}
}

// test that retries don't occur for 4xx status codes with a response type that can't be retried
func TestRequest4xxUnretryable(t *testing.T) {
	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 1
	})

	s := awstesting.NewClient(cfg)

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &http.Response{
			StatusCode: 401,
			Body:       body(`{"__type":"SignatureDoesNotMatch","message":"Signature does not match."}`),
		}
	})
	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	err := r.Send()
	if err == nil {
		t.Fatalf("expect error, but did not get one")
	}
	aerr := err.(awserr.RequestFailure)
	if e, a := 401, aerr.StatusCode(); e != a {
		t.Errorf("expect %d status code, got %d", e, a)
	}
	if e, a := "SignatureDoesNotMatch", aerr.Code(); e != a {
		t.Errorf("expect %q error code, got %q", e, a)
	}
	if e, a := "Signature does not match.", aerr.Message(); e != a {
		t.Errorf("expect %q error message, got %q", e, a)
	}
	if e, a := 0, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}
}

func TestRequestExhaustRetries(t *testing.T) {
	restoreSleep := mockSleep()
	defer restoreSleep()

	orig := sdk.SleepWithContext
	defer func() { sdk.SleepWithContext = orig }()

	var delays []time.Duration
	sdk.SleepWithContext = func(ctx context.Context, dur time.Duration) error {
		delays = append(delays, dur)
		return nil
	}

	reqNum := 0
	reqs := []http.Response{
		{StatusCode: 500, Body: body(`{"__type":"UnknownError","message":"An error occurred."}`)},
		{StatusCode: 500, Body: body(`{"__type":"UnknownError","message":"An error occurred."}`)},
		{StatusCode: 500, Body: body(`{"__type":"UnknownError","message":"An error occurred."}`)},
		{StatusCode: 500, Body: body(`{"__type":"UnknownError","message":"An error occurred."}`)},
	}

	s := awstesting.NewClient(unit.Config())

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &reqs[reqNum]
		reqNum++
	})
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
	err := r.Send()
	if err == nil {
		t.Fatalf("expect error, but did not get one")
	}
	aerr := err.(awserr.RequestFailure)
	if e, a := 500, aerr.StatusCode(); e != a {
		t.Errorf("expect %d status code, got %d", e, a)
	}
	if e, a := "UnknownError", aerr.Code(); e != a {
		t.Errorf("expect %q error code, got %q", e, a)
	}
	if e, a := "An error occurred.", aerr.Message(); e != a {
		t.Errorf("expect %q error message, got %q", e, a)
	}
	if e, a := 3, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}

	expectDelays := []struct{ min, max time.Duration }{{30, 60}, {60, 120}, {120, 240}}
	for i, v := range delays {
		min := expectDelays[i].min * time.Millisecond
		max := expectDelays[i].max * time.Millisecond
		if !(min <= v && v <= max) {
			t.Errorf("expect delay to be within range, i:%d, v:%s, min:%s, max:%s",
				i, v, min, max)
		}
	}
}

// test that the request is retried after the credentials are expired.
func TestRequest_RecoverExpiredCreds(t *testing.T) {
	reqNum := 0
	reqs := []http.Response{
		{StatusCode: 400, Body: body(`{"__type":"ExpiredTokenException","message":"expired token"}`)},
		{StatusCode: 200, Body: body(`{"data":"valid"}`)},
	}
	expectCreds := []aws.Credentials{
		{
			AccessKeyID:     "expiredKey",
			SecretAccessKey: "expiredSecret",
		},
		{
			AccessKeyID:     "AKID",
			SecretAccessKey: "SECRET",
		},
	}

	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 10
	})

	credsInvalidated := false
	credsProvider := func() aws.CredentialsProvider {
		creds := expectCreds[0]
		return awstesting.MockCredentialsProvider{
			RetrieveFn: func(ctx context.Context) (aws.Credentials, error) {
				return creds, nil
			},
			InvalidateFn: func() {
				creds = expectCreds[1]
				credsInvalidated = true
			},
		}
	}()
	cfg.Credentials = credsProvider

	s := awstesting.NewClient(cfg)

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.Build.PushFront(func(r *aws.Request) {
		creds, err := r.Config.Credentials.Retrieve(context.Background())
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
		if e, a := "expiredKey", creds.AccessKeyID; e != a {
			t.Errorf("expect %v key, got %v", e, a)
		}
		if e, a := "expiredSecret", creds.SecretAccessKey; e != a {
			t.Errorf("expect %v secret, got %v", e, a)
		}
	})

	s.Handlers.AfterRetry.PushBack(func(r *aws.Request) {
		if !credsInvalidated {
			t.Errorf("expect creds to be invalidated")
		}
	})

	s.Handlers.Sign.Clear()
	s.Handlers.Sign.PushBack(func(r *aws.Request) {
		creds, err := r.Config.Credentials.Retrieve(context.Background())
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
		if e, a := expectCreds[reqNum].AccessKeyID, creds.AccessKeyID; e != a {
			t.Errorf("expect %v key, got %v", e, a)
		}
		if e, a := expectCreds[reqNum].SecretAccessKey, creds.SecretAccessKey; e != a {
			t.Errorf("expect %v secret, got %v", e, a)
		}
	})
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &reqs[reqNum]
		reqNum++
	})

	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	err := r.Send()
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	creds, err := r.Config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if creds.Expired() {
		t.Errorf("expect valid creds after cred expired recovery")
	}

	if e, a := 1, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}
	if e, a := "valid", out.Data; e != a {
		t.Errorf("expect %q output got %q", e, a)
	}
}

func TestMakeAddtoUserAgentHandler(t *testing.T) {
	fn := aws.MakeAddToUserAgentHandler("name", "version", "extra1", "extra2")
	r := &aws.Request{HTTPRequest: &http.Request{Header: http.Header{}}}
	r.HTTPRequest.Header.Set("User-Agent", "foo/bar")
	fn(r)

	if e, a := "foo/bar name/version (extra1; extra2)", r.HTTPRequest.Header.Get("User-Agent"); !strings.HasPrefix(a, e) {
		t.Errorf("expect %q user agent, got %q", e, a)
	}
}

func TestMakeAddtoUserAgentFreeFormHandler(t *testing.T) {
	fn := aws.MakeAddToUserAgentFreeFormHandler("name/version (extra1; extra2)")
	r := &aws.Request{HTTPRequest: &http.Request{Header: http.Header{}}}
	r.HTTPRequest.Header.Set("User-Agent", "foo/bar")
	fn(r)

	if e, a := "foo/bar name/version (extra1; extra2)", r.HTTPRequest.Header.Get("User-Agent"); !strings.HasPrefix(a, e) {
		t.Errorf("expect %q user agent, got %q", e, a)
	}
}

func TestRequestUserAgent(t *testing.T) {
	cfg := unit.Config()
	cfg.Region = "us-east-1"

	s := awstesting.NewClient(cfg)

	req := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, &testData{})
	req.HTTPRequest.Header.Set("User-Agent", "foo/bar")
	if err := req.Build(); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	expectUA := fmt.Sprintf("foo/bar %s/%s (%s; %s; %s)",
		aws.SDKName, aws.SDKVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	if e, a := expectUA, req.HTTPRequest.Header.Get("User-Agent"); !strings.HasPrefix(a, e) {
		t.Errorf("expect %q user agent, got %q", e, a)
	}
}

func TestRequestThrottleRetries(t *testing.T) {
	orig := sdk.SleepWithContext
	defer func() { sdk.SleepWithContext = orig }()

	var delays []time.Duration
	sdk.SleepWithContext = func(ctx context.Context, dur time.Duration) error {
		delays = append(delays, dur)
		return nil
	}

	reqNum := 0
	reqs := []http.Response{
		{StatusCode: 500, Body: body(`{"__type":"Throttling","message":"An error occurred."}`)},
		{StatusCode: 500, Body: body(`{"__type":"Throttling","message":"An error occurred."}`)},
		{StatusCode: 500, Body: body(`{"__type":"Throttling","message":"An error occurred."}`)},
		{StatusCode: 500, Body: body(`{"__type":"Throttling","message":"An error occurred."}`)},
	}

	s := awstesting.NewClient(unit.Config())

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &reqs[reqNum]
		reqNum++
	})
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
	err := r.Send()
	if err == nil {
		t.Fatalf("expect error, but did not get one")
	}
	aerr := err.(awserr.RequestFailure)
	if e, a := 500, aerr.StatusCode(); e != a {
		t.Errorf("expect %d status code, got %d", e, a)
	}
	if e, a := "Throttling", aerr.Code(); e != a {
		t.Errorf("expect %q error code, got %q", e, a)
	}
	if e, a := "An error occurred.", aerr.Message(); e != a {
		t.Errorf("expect %q error message, got %q", e, a)
	}
	if e, a := 3, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}

	expectDelays := []struct{ min, max time.Duration }{{500, 1000}, {1000, 2000}, {2000, 4000}}
	for i, v := range delays {
		min := expectDelays[i].min * time.Millisecond
		max := expectDelays[i].max * time.Millisecond
		if !(min <= v && v <= max) {
			t.Errorf("expect delay to be within range, i:%d, v:%s, min:%s, max:%s",
				i, v, min, max)
		}
	}
}

// test that retries occur for request timeouts when response.Body can be nil
func TestRequestRecoverTimeoutWithNilBody(t *testing.T) {
	reqNum := 0
	reqs := []*http.Response{
		{StatusCode: 0, Body: nil}, // body can be nil when requests time out
		{StatusCode: 200, Body: body(`{"data":"valid"}`)},
	}
	errors := []error{
		errTimeout, nil,
	}

	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 10
	})

	s := awstesting.NewClient(cfg)

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.AfterRetry.Clear() // force retry on all errors
	s.Handlers.AfterRetry.PushBack(func(r *aws.Request) {
		if r.Error != nil {
			r.Error = nil
			r.Retryable = aws.Bool(true)
			r.RetryCount++
		}
	})
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = reqs[reqNum]
		r.Error = errors[reqNum]
		reqNum++
	})
	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	err := r.Send()
	if err != nil {
		t.Fatalf("expect no error, but got %v", err)
	}
	if e, a := 1, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}
	if e, a := "valid", out.Data; e != a {
		t.Errorf("expect %q output got %q", e, a)
	}
}

func TestRequestRecoverTimeoutWithNilResponse(t *testing.T) {
	reqNum := 0
	reqs := []*http.Response{
		nil,
		{StatusCode: 200, Body: body(`{"data":"valid"}`)},
	}
	errors := []error{
		errTimeout,
		nil,
	}

	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 10
	})

	s := awstesting.NewClient(cfg)

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.AfterRetry.Clear() // force retry on all errors
	s.Handlers.AfterRetry.PushBack(func(r *aws.Request) {
		if r.Error != nil {
			r.Error = nil
			r.Retryable = aws.Bool(true)
			r.RetryCount++
		}
	})
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = reqs[reqNum]
		r.Error = errors[reqNum]
		reqNum++
	})
	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	err := r.Send()
	if err != nil {
		t.Fatalf("expect no error, but got %v", err)
	}
	if e, a := 1, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}
	if e, a := "valid", out.Data; e != a {
		t.Errorf("expect %q output got %q", e, a)
	}
}

func TestRequest_NoBody(t *testing.T) {
	cases := []string{
		"GET", "HEAD", "DELETE",
		"PUT", "POST", "PATCH",
	}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if v := r.TransferEncoding; len(v) > 0 {
					t.Errorf("expect no body sent with Transfer-Encoding, %v", v)
				}

				outMsg := []byte(`{"Value": "abc"}`)

				if b, err := ioutil.ReadAll(r.Body); err != nil {
					t.Fatalf("expect no error reading request body, got %v", err)
				} else if n := len(b); n > 0 {
					t.Errorf("expect no request body, got %d bytes", n)
				}

				w.Header().Set("Content-Length", strconv.Itoa(len(outMsg)))
				if _, err := w.Write(outMsg); err != nil {
					t.Fatalf("expect no error writing server response, got %v", err)
				}
			}))
			defer server.Close()

			cfg := unit.Config()
			cfg.Region = "mock-region"
			cfg.Retryer = aws.NoOpRetryer{}
			cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

			s := awstesting.NewClient(cfg)

			s.Handlers.Build.PushBack(rest.Build)
			s.Handlers.Validate.Clear()
			s.Handlers.Unmarshal.PushBack(unmarshal)
			s.Handlers.UnmarshalError.PushBack(unmarshalError)

			in := struct {
				Bucket *string `location:"uri" locationName:"bucket"`
				Key    *string `location:"uri" locationName:"key"`
			}{
				Bucket: aws.String("mybucket"), Key: aws.String("myKey"),
			}

			out := struct {
				Value *string
			}{}

			r := s.NewRequest(&aws.Operation{
				Name: "OpName", HTTPMethod: c, HTTPPath: "/{bucket}/{key+}",
			}, &in, &out)

			if err := r.Send(); err != nil {
				t.Fatalf("expect no error sending request, got %v", err)
			}
		})
	}
}

func TestIsSerializationErrorRetryable(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{
			err:      awserr.New(aws.ErrCodeSerialization, "foo error", nil),
			expected: false,
		},
		{
			err:      awserr.New("ErrFoo", "foo error", nil),
			expected: false,
		},
		{
			err:      nil,
			expected: false,
		},
		{
			err:      awserr.New(aws.ErrCodeSerialization, "foo error", stubConnectionResetError),
			expected: true,
		},
	}

	for i, c := range testCases {
		r := &aws.Request{
			Error: c.err,
		}
		if r.IsErrorRetryable() != c.expected {
			t.Errorf("Case %d: expected %v, but received %v", i, c.expected, !c.expected)
		}
	}
}

func TestWithLogLevel(t *testing.T) {
	r := &aws.Request{}

	opt := aws.WithLogLevel(aws.LogDebugWithHTTPBody)
	r.ApplyOptions(opt)

	if !r.Config.LogLevel.Matches(aws.LogDebugWithHTTPBody) {
		t.Errorf("expect log level to be set, but was not, %v",
			r.Config.LogLevel)
	}
}

func TestWithGetResponseHeader(t *testing.T) {
	r := &aws.Request{}

	var val, val2 string
	r.ApplyOptions(
		aws.WithGetResponseHeader("x-a-header", &val),
		aws.WithGetResponseHeader("x-second-header", &val2),
	)

	r.HTTPResponse = &http.Response{
		Header: func() http.Header {
			h := http.Header{}
			h.Set("x-a-header", "first")
			h.Set("x-second-header", "second")
			return h
		}(),
	}
	r.Handlers.Complete.Run(r)

	if e, a := "first", val; e != a {
		t.Errorf("expect %q header value got %q", e, a)
	}
	if e, a := "second", val2; e != a {
		t.Errorf("expect %q header value got %q", e, a)
	}
}

func TestWithGetResponseHeaders(t *testing.T) {
	r := &aws.Request{}

	var headers http.Header
	opt := aws.WithGetResponseHeaders(&headers)

	r.ApplyOptions(opt)

	r.HTTPResponse = &http.Response{
		Header: func() http.Header {
			h := http.Header{}
			h.Set("x-a-header", "headerValue")
			return h
		}(),
	}
	r.Handlers.Complete.Run(r)

	if e, a := "headerValue", headers.Get("x-a-header"); e != a {
		t.Errorf("expect %q header value got %q", e, a)
	}
}

type connResetCloser struct {
}

func (rc *connResetCloser) Read(b []byte) (int, error) {
	return 0, stubConnectionResetError
}

func (rc *connResetCloser) Close() error {
	return nil
}

func TestSerializationErrConnectionReset(t *testing.T) {
	restoreSleep := mockSleep()
	defer restoreSleep()

	count := 0
	handlers := aws.Handlers{}
	handlers.Send.PushBack(func(r *aws.Request) {
		count++
		r.HTTPResponse = &http.Response{}
		r.HTTPResponse.Body = &connResetCloser{}
	})

	handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)
	handlers.AfterRetry.PushBackNamed(defaults.AfterRetryHandler)

	op := &aws.Operation{
		Name:       "op",
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	meta := aws.Metadata{
		ServiceName:   "fooService",
		SigningName:   "foo",
		SigningRegion: "foo",
		APIVersion:    "2001-01-01",
		JSONVersion:   "1.1",
		TargetPrefix:  "Foo",
	}
	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 5
	})

	req := aws.New(
		cfg,
		meta,
		handlers,
		aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
			d.NumMaxRetries = 5
		}),
		op,
		&struct{}{},
		&struct{}{},
	)

	osErr := stubConnectionResetError
	req.ApplyOptions(aws.WithResponseReadTimeout(time.Second))
	err := req.Send()
	if err == nil {
		t.Error("expected rror 'SerializationError', but received nil")
	}
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() != "SerializationError" {
		t.Errorf("expected 'SerializationError', but received %q", aerr.Code())
	} else if !ok {
		t.Errorf("expected 'awserr.Error', but received %v", reflect.TypeOf(err))
	} else if aerr.OrigErr().Error() != osErr.Error() {
		t.Errorf("expected %q, but received %q", osErr.Error(), aerr.OrigErr().Error())
	}

	if count != 6 {
		t.Errorf("expected '6', but received %d", count)
	}
}

type testRetryer struct {
	shouldRetry bool
}

func (d *testRetryer) MaxRetries() int {
	return 3
}

// RetryRules returns the delay duration before retrying this request again
func (d *testRetryer) RetryRules(r *aws.Request) time.Duration {
	return time.Duration(time.Millisecond)
}

func (d *testRetryer) ShouldRetry(r *aws.Request) bool {
	d.shouldRetry = true
	if r.Retryable != nil {
		return *r.Retryable
	}

	if r.HTTPResponse.StatusCode >= 500 {
		return true
	}
	return r.IsErrorRetryable()
}

func TestEnforceShouldRetryCheck(t *testing.T) {
	restoreSleep := mockSleep()
	defer restoreSleep()

	tp := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: 1 * time.Millisecond,
	}

	client := &http.Client{Transport: tp}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This server should wait forever. Requests will timeout and the SDK should
		// attempt to retry.
		select {}
	}))

	retryer := &testRetryer{}

	cfg := unit.Config()
	cfg.Region = "mock-region"
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)
	cfg.Retryer = retryer
	cfg.HTTPClient = client

	s := awstesting.NewClient(cfg)

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)

	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	err := r.Send()
	if err == nil {
		t.Fatalf("expect error, but got nil")
	}
	if e, a := 3, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}
	if !retryer.shouldRetry {
		t.Errorf("expect 'true' for ShouldRetry, but got %v", retryer.shouldRetry)
	}
}

type errReader struct {
	err error
}

func (reader *errReader) Read(b []byte) (int, error) {
	return 0, reader.err
}

func (reader *errReader) Close() error {
	return nil
}

func TestIsNoBodyReader(t *testing.T) {
	cases := []struct {
		reader io.ReadCloser
		expect bool
	}{
		{ioutil.NopCloser(bytes.NewReader([]byte("abc"))), false},
		{ioutil.NopCloser(bytes.NewReader(nil)), false},
		{nil, false},
		{http.NoBody, true},
	}

	for i, c := range cases {
		if e, a := c.expect, http.NoBody == c.reader; e != a {
			t.Errorf("%d, expect %t match, but was %t", i, e, a)
		}
	}
}

func TestRequest_TemporaryRetry(t *testing.T) {
	restoreSleep := mockSleep()
	defer restoreSleep()

	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1024")
		w.WriteHeader(http.StatusOK)

		w.Write(make([]byte, 100))

		f := w.(http.Flusher)
		f.Flush()

		<-done
	}))
	defer server.Close()

	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 1
	})
	cfg.HTTPClient = &http.Client{
		Timeout: 100 * time.Millisecond,
	}
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := awstesting.NewClient(cfg)

	req := svc.NewRequest(&aws.Operation{
		Name: "name", HTTPMethod: "GET", HTTPPath: "/path",
	}, &struct{}{}, &struct{}{})

	req.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		defer req.HTTPResponse.Body.Close()
		_, err := io.Copy(ioutil.Discard, req.HTTPResponse.Body)
		r.Error = awserr.New(aws.ErrCodeSerialization, "error", err)
	})

	err := req.Send()
	if err == nil {
		t.Errorf("expect error, got none")
	}
	close(done)

	aerr := err.(awserr.Error)
	if e, a := aws.ErrCodeSerialization, aerr.Code(); e != a {
		t.Errorf("expect %q error code, got %q", e, a)
	}

	if e, a := 1, req.RetryCount; e != a {
		t.Errorf("expect %d retries, got %d", e, a)
	}

	type temporary interface {
		Temporary() bool
	}

	terr, ok := aerr.OrigErr().(temporary)
	if !ok {
		t.Fatalf("expect error to implement temporary, got %T", aerr.OrigErr())
	}
	if !terr.Temporary() {
		t.Errorf("expect temporary error, was not")
	}
}

func TestSanitizeHostForHeader(t *testing.T) {
	cases := []struct {
		url                 string
		expectedRequestHost string
	}{
		{"https://estest.us-east-1.es.amazonaws.com:443", "estest.us-east-1.es.amazonaws.com"},
		{"https://estest.us-east-1.es.amazonaws.com", "estest.us-east-1.es.amazonaws.com"},
		{"https://localhost:9200", "localhost:9200"},
		{"http://localhost:80", "localhost"},
		{"http://localhost:8080", "localhost:8080"},
	}

	for _, c := range cases {
		r, _ := http.NewRequest("GET", c.url, nil)
		aws.SanitizeHostForHeader(r)

		if h := r.Host; h != c.expectedRequestHost {
			t.Errorf("expect %v host, got %q", c.expectedRequestHost, h)
		}
	}
}

func TestRequestBodySeekFails(t *testing.T) {
	s := awstesting.NewClient(unit.Config())
	s.Handlers.Validate.Clear()
	s.Handlers.Build.Clear()

	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	r.SetReaderBody(&stubSeekFail{
		Err: fmt.Errorf("failed to seek reader"),
	})
	err := r.Send()
	if err == nil {
		t.Fatal("expect error, but got none")
	}

	aerr := err.(awserr.Error)
	if e, a := aws.ErrCodeSerialization, aerr.Code(); e != a {
		t.Errorf("expect %v error code, got %v", e, a)
	}

}

func Test501NotRetrying(t *testing.T) {
	reqNum := 0
	reqs := []http.Response{
		{StatusCode: 500, Body: body(`{"__type":"UnknownError","message":"An error occurred."}`)},
		{StatusCode: 501, Body: body(`{"__type":"NotImplemented","message":"An error occurred."}`)},
		{StatusCode: 200, Body: body(`{"data":"valid"}`)},
	}

	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 10
	})
	s := awstesting.NewClient(cfg)
	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &reqs[reqNum]
		reqNum++
	})
	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	err := r.Send()
	if err == nil {
		t.Fatal("expect error, but got none")
	}

	aerr := err.(awserr.Error)
	if e, a := "NotImplemented", aerr.Code(); e != a {
		t.Errorf("expected error code %q, but received %q", e, a)
	}
	if e, a := 1, int(r.RetryCount); e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}
}

func TestRequestInvalidEndpoint(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("http://localhost:90 ")

	r := aws.New(
		cfg,
		aws.Metadata{},
		cfg.Handlers,
		aws.NewDefaultRetryer(),
		&aws.Operation{},
		nil,
		nil,
	)

	if r.Error == nil {
		t.Errorf("expect error, got none")
	}
}

func TestResetBody_WithEmptyBody(t *testing.T) {
	r := aws.Request{
		HTTPRequest: &http.Request{},
	}

	reader := strings.NewReader("")
	r.Body = reader

	r.ResetBody()

	if a, e := r.HTTPRequest.Body, http.NoBody; a != e {
		t.Errorf("expected request body to be set to reader, got %#v",
			r.HTTPRequest.Body)
	}
}

func TestRequest_FollowPUTRedirects(t *testing.T) {
	const bodySize = 1024

	redirectHit := 0
	endpointHit := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/redirect-me":
			u := *r.URL
			u.Path = "/endpoint"
			w.Header().Set("Location", u.String())
			w.WriteHeader(307)
			redirectHit++
		case "/endpoint":
			b := bytes.Buffer{}
			io.Copy(&b, r.Body)
			r.Body.Close()
			if e, a := bodySize, b.Len(); e != a {
				t.Fatalf("expect %d body size, got %d", e, a)
			}
			endpointHit++
		default:
			t.Fatalf("unexpected endpoint used, %q", r.URL.String())
		}
	}))
	defer server.Close()

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := awstesting.NewClient(cfg)

	req := svc.NewRequest(&aws.Operation{
		Name:       "Operation",
		HTTPMethod: "PUT",
		HTTPPath:   "/redirect-me",
	}, &struct{}{}, &struct{}{})
	req.SetReaderBody(bytes.NewReader(make([]byte, bodySize)))

	err := req.Send()
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := 1, redirectHit; e != a {
		t.Errorf("expect %d redirect hits, got %d", e, a)
	}
	if e, a := 1, endpointHit; e != a {
		t.Errorf("expect %d endpoint hits, got %d", e, a)
	}
}

type timeoutErr struct {
	error
}

var errTimeout = awserr.New("foo", "bar", &timeoutErr{
	errors.New("net/http: request canceled"),
})

type stubSeekFail struct {
	Err error
}

func (f *stubSeekFail) Read(b []byte) (int, error) {
	return len(b), nil
}
func (f *stubSeekFail) ReadAt(b []byte, offset int64) (int, error) {
	return len(b), nil
}
func (f *stubSeekFail) Seek(offset int64, mode int) (int64, error) {
	return 0, f.Err
}

func TestRequestEndpointConstruction(t *testing.T) {
	cases := map[string]struct {
		EndpointResolver aws.EndpointResolverFunc
		ExpectedEndpoint aws.Endpoint
	}{
		"resolved modeled endpoint": {
			EndpointResolver: func(_, _ string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "https://localhost",
					SigningName:   "foo-service",
					SigningRegion: "bar-region",
				}, nil
			},
			ExpectedEndpoint: aws.Endpoint{
				URL:           "https://localhost",
				SigningName:   "foo-service",
				SigningRegion: "bar-region",
			},
		},
		"resolved endpoint missing signing region": {
			EndpointResolver: func(_, _ string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:         "https://localhost",
					SigningName: "foo-service",
				}, nil
			},
			ExpectedEndpoint: aws.Endpoint{
				URL:           "https://localhost",
				SigningName:   "foo-service",
				SigningRegion: "meta-bar-region",
			},
		},
		"resolved endpoint missing signing name": {
			EndpointResolver: func(_, _ string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "https://localhost",
					SigningRegion: "bar-region",
				}, nil
			},
			ExpectedEndpoint: aws.Endpoint{
				URL:           "https://localhost",
				SigningName:   "meta-foo-service",
				SigningRegion: "bar-region",
			},
		},
		"resolved endpoint signing name derived": {
			EndpointResolver: func(_, _ string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:                "https://localhost",
					SigningRegion:      "bar-region",
					SigningName:        "derived-signing-name",
					SigningNameDerived: true,
				}, nil
			},
			ExpectedEndpoint: aws.Endpoint{
				URL:                "https://localhost",
				SigningName:        "meta-foo-service",
				SigningRegion:      "bar-region",
				SigningNameDerived: true,
			},
		},
		"resolved endpoint missing signing region and signing name": {
			EndpointResolver: func(_, _ string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: "https://localhost",
				}, nil
			},
			ExpectedEndpoint: aws.Endpoint{
				URL:           "https://localhost",
				SigningName:   "meta-foo-service",
				SigningRegion: "meta-bar-region",
			},
		},
	}

	meta := aws.Metadata{
		ServiceName:   "FooService",
		SigningName:   "meta-foo-service",
		SigningRegion: "meta-bar-region",
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			client := aws.NewClient(aws.Config{EndpointResolver: tt.EndpointResolver}, meta)

			request := client.NewRequest(&aws.Operation{Name: "ZapOperation", HTTPMethod: "PUT", HTTPPath: "/"}, nil, nil)

			if e, a := tt.ExpectedEndpoint, request.Endpoint; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, got %v", e, a)
			}
		})
	}
}
