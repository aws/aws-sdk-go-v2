package middleware_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"net/http"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/internal/rand"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestClientRequestID(t *testing.T) {
	oReader := rand.Reader
	defer func() {
		rand.Reader = oReader
	}()
	rand.Reader = bytes.NewReader(make([]byte, 16))

	mid := middleware.ClientRequestID{}

	in := smithymiddleware.BuildInput{Request: &smithyhttp.Request{Request: &http.Request{Header: make(http.Header)}}}
	ctx := context.Background()
	_, _, err := mid.HandleBuild(ctx, in, smithymiddleware.BuildHandlerFunc(func(ctx context.Context, input smithymiddleware.BuildInput) (
		out smithymiddleware.BuildOutput, metadata smithymiddleware.Metadata, err error,
	) {
		req := in.Request.(*smithyhttp.Request)

		value := req.Header.Get("amz-sdk-invocation-id")

		expected := "00000000-0000-4000-8000-000000000000"
		if value != expected {
			t.Errorf("expect %v, got %v", expected, value)
		}

		return out, metadata, err
	}))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	in = smithymiddleware.BuildInput{}
	_, _, err = mid.HandleBuild(ctx, in, nil)
	if err != nil {
		if e, a := "unknown transport type", err.Error(); !strings.Contains(a, e) {
			t.Errorf("expected %q, got %q", e, a)
		}
	} else {
		t.Errorf("expected error, got %q", err)
	}
}

func TestAttemptClockSkewHandler(t *testing.T) {
	cases := map[string]struct {
		Next              smithymiddleware.DeserializeHandlerFunc
		ResponseAt        func() time.Time
		ExpectAttemptSkew time.Duration
		ExpectServerTime  time.Time
		ExpectResponseAt  time.Time
	}{
		"no response": {
			Next: func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				return out, m, err
			},
			ResponseAt: func() time.Time {
				return time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
			},
			ExpectResponseAt: time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		},
		"failed response": {
			Next: func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				out.RawResponse = &smithyhttp.Response{
					Response: &http.Response{
						StatusCode: 0,
						Header:     http.Header{},
					},
				}
				return out, m, err
			},
			ResponseAt: func() time.Time {
				return time.Date(2020, 6, 7, 8, 9, 10, 0, time.UTC)
			},
			ExpectResponseAt: time.Date(2020, 6, 7, 8, 9, 10, 0, time.UTC),
		},
		"no date header response": {
			Next: func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				out.RawResponse = &smithyhttp.Response{
					Response: &http.Response{
						StatusCode: 200,
						Header:     http.Header{},
					},
				}
				return out, m, err
			},
			ResponseAt: func() time.Time {
				return time.Date(2020, 11, 12, 13, 14, 15, 0, time.UTC)
			},
			ExpectResponseAt: time.Date(2020, 11, 12, 13, 14, 15, 0, time.UTC),
		},
		"invalid date header response": {
			Next: func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				out.RawResponse = &smithyhttp.Response{
					Response: &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Date": []string{"abc123"},
						},
					},
				}
				return out, m, err
			},
			ResponseAt: func() time.Time {
				return time.Date(2020, 1, 2, 16, 17, 18, 0, time.UTC)
			},
			ExpectResponseAt: time.Date(2020, 1, 2, 16, 17, 18, 0, time.UTC),
		},
		"date response": {
			Next: func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				out.RawResponse = &smithyhttp.Response{
					Response: &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Date": []string{"Thu, 05 Mar 2020 22:25:15 GMT"},
						},
					},
				}
				return out, m, err
			},
			ResponseAt: func() time.Time {
				return time.Date(2020, 3, 5, 22, 25, 17, 0, time.UTC)
			},
			ExpectResponseAt:  time.Date(2020, 3, 5, 22, 25, 17, 0, time.UTC),
			ExpectServerTime:  time.Date(2020, 3, 5, 22, 25, 15, 0, time.UTC),
			ExpectAttemptSkew: -2 * time.Second,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if c.ResponseAt != nil {
				sdkTime := sdk.NowTime
				defer func() {
					sdk.NowTime = sdkTime
				}()
				sdk.NowTime = c.ResponseAt
			}
			mw := middleware.RecordResponseTiming{}
			_, metadata, err := mw.HandleDeserialize(context.Background(), smithymiddleware.DeserializeInput{}, c.Next)
			if err != nil {
				t.Errorf("expect no error, got %v", err)
			}

			if v, ok := middleware.GetResponseAt(metadata); ok {
				if !reflect.DeepEqual(v, c.ExpectResponseAt) {
					t.Fatalf("expected %v, got %v", c.ExpectResponseAt, v)
				}
			} else if !c.ExpectResponseAt.IsZero() {
				t.Fatal("expected response at to be set in metadata, was not")
			}

			if v, ok := middleware.GetServerTime(metadata); ok {
				if !reflect.DeepEqual(v, c.ExpectServerTime) {
					t.Fatalf("expected %v, got %v", c.ExpectServerTime, v)
				}
			} else if !c.ExpectServerTime.IsZero() {
				t.Fatal("expected server time to be set in metadata, was not")
			}

			if v, ok := middleware.GetAttemptSkew(metadata); ok {
				if !reflect.DeepEqual(v, c.ExpectAttemptSkew) {
					t.Fatalf("expected %v, got %v", c.ExpectAttemptSkew, v)
				}
			} else if c.ExpectAttemptSkew != 0 {
				t.Fatal("expected attempt skew to be set in metadata, was not")
			}
		})
	}
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Options struct {
	HTTPClient HTTPClient
	RetryMode  aws.RetryMode
	Retryer    aws.Retryer
	Offset     *atomic.Int64
}

type MockClient struct {
	options Options
}

func addRetry(stack *smithymiddleware.Stack, o Options) error {
	attempt := retry.NewAttemptMiddleware(o.Retryer, smithyhttp.RequestCloner, func(m *retry.Attempt) {
		m.LogAttempts = false
	})
	return stack.Finalize.Add(attempt, smithymiddleware.After)
}

func addOffset(stack *smithymiddleware.Stack, o Options) error {
	buildOffset := middleware.AddTimeOffsetBuildMiddleware{Offset: o.Offset}
	deserializeOffset := middleware.AddTimeOffsetDeserializeMiddleware{Offset: o.Offset}
	err := stack.Build.Add(&buildOffset, smithymiddleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&deserializeOffset, smithymiddleware.Before)
	if err != nil {
		return err
	}
	return nil
}

// Middleware to set a `Date` object that includes sdk time and offset
type MockAddDateHeader struct {
}

func (l *MockAddDateHeader) ID() string {
	return "MockAddDateHeader"
}

func (l *MockAddDateHeader) HandleFinalize(
	ctx context.Context, in smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler,
) (
	out smithymiddleware.FinalizeOutput, metadata smithymiddleware.Metadata, attemptError error,
) {
	req := in.Request.(*smithyhttp.Request)
	date := sdk.NowTime()
	skew := middleware.GetAttemptSkewContext(ctx)
	date = date.Add(skew)
	req.Header.Set("Date", date.Format(time.RFC850))
	return next.HandleFinalize(ctx, in)
}

// Middleware to deserialize the response which just says "OK" if the response is 200
type DeserializeFailIfNotHTTP200 struct {
}

func (*DeserializeFailIfNotHTTP200) ID() string {
	return "DeserializeFailIfNotHTTP200"
}

func (m *DeserializeFailIfNotHTTP200) HandleDeserialize(ctx context.Context, in smithymiddleware.DeserializeInput, next smithymiddleware.DeserializeHandler) (
	out smithymiddleware.DeserializeOutput, metadata smithymiddleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}
	response, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok {
		return out, metadata, fmt.Errorf("expected raw response to be set on testing")
	}
	if response.StatusCode != 200 {
		return out, metadata, mockRetryableError{true}
	}
	return out, metadata, err
}

func (c *MockClient) setupMiddleware(stack *smithymiddleware.Stack) error {
	err := error(nil)
	if c.options.Retryer != nil {
		err = addRetry(stack, c.options)
		if err != nil {
			return err
		}
	}
	if c.options.Offset != nil {
		err = addOffset(stack, c.options)
		if err != nil {
			return err
		}
	}
	err = stack.Finalize.Add(&MockAddDateHeader{}, smithymiddleware.After)
	if err != nil {
		return err
	}
	err = middleware.AddRecordResponseTiming(stack)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&DeserializeFailIfNotHTTP200{}, smithymiddleware.After)
	if err != nil {
		return err
	}
	return nil
}

func (c *MockClient) Do(ctx context.Context) (interface{}, error) {
	// setup middlewares
	ctx = smithymiddleware.ClearStackValues(ctx)
	stack := smithymiddleware.NewStack("stack", smithyhttp.NewStackRequest)
	err := c.setupMiddleware(stack)
	if err != nil {
		return nil, err
	}
	handler := smithymiddleware.DecorateHandler(smithyhttp.NewClientHandler(c.options.HTTPClient), stack)
	result, _, err := handler.Handle(ctx, 1)
	if err != nil {
		return nil, err
	}
	return result, err
}

type mockRetryableError struct{ b bool }

func (m mockRetryableError) RetryableError() bool { return m.b }
func (m mockRetryableError) Error() string {
	return fmt.Sprintf("mock retryable %t", m.b)
}

func failRequestIfSkewed() smithyhttp.ClientDoFunc {
	return func(req *http.Request) (*http.Response, error) {
		dateHeader := req.Header.Get("Date")
		if dateHeader == "" {
			return nil, fmt.Errorf("expected `Date` header to be set")
		}
		reqDate, err := time.Parse(time.RFC850, dateHeader)
		if err != nil {
			return nil, err
		}
		parsedReqTime := time.Now().Sub(reqDate)
		parsedReqTime = time.Duration.Abs(parsedReqTime)
		thresholdForSkewError := 4 * time.Minute
		if thresholdForSkewError-parsedReqTime <= 0 {
			return &http.Response{
				StatusCode: 403,
				Header: http.Header{
					"Date": {time.Now().Format(time.RFC850)},
				},
			}, nil
		}
		// else, return OK
		return &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
		}, nil
	}
}

func TestSdkOffsetIsSet(t *testing.T) {
	nowTime := sdk.NowTime
	defer func() {
		sdk.NowTime = nowTime
	}()
	fiveMinuteSkew := func() time.Time {
		return time.Now().Add(5 * time.Minute)
	}
	sdk.NowTime = fiveMinuteSkew
	c := MockClient{
		Options{
			HTTPClient: failRequestIfSkewed(),
		},
	}
	resp, err := c.Do(context.Background())
	if err == nil {
		t.Errorf("Expected first request to fail since clock skew logic has not run. Got %v and err %v", resp, err)
	}
}

func TestRetrySetsSkewInContext(t *testing.T) {
	defer resetDefaults(sdk.TestingUseNopSleep())
	fiveMinuteSkew := func() time.Time {
		return time.Now().Add(5 * time.Minute)
	}
	sdk.NowTime = fiveMinuteSkew
	c := MockClient{
		Options{
			HTTPClient: failRequestIfSkewed(),
			Retryer: retry.NewStandard(func(s *retry.StandardOptions) {
			}),
		},
	}
	resp, err := c.Do(context.Background())
	if err != nil {
		t.Errorf("Expected request to succeed on retry. Got %v and err %v", resp, err)
	}
}

func TestSkewIsSetOnTheWholeClient(t *testing.T) {
	defer resetDefaults(sdk.TestingUseNopSleep())
	fiveMinuteSkew := func() time.Time {
		return time.Now().Add(5 * time.Minute)
	}
	sdk.NowTime = fiveMinuteSkew
	var offset atomic.Int64
	offset.Store(0)
	c := MockClient{
		Options{
			HTTPClient: failRequestIfSkewed(),
			Retryer: retry.NewStandard(func(s *retry.StandardOptions) {
			}),
			Offset: &offset,
		},
	}
	resp, err := c.Do(context.Background())
	if err != nil {
		t.Errorf("Expected request to succeed on retry. Got %v and err %v", resp, err)
	}
	// Remove retryer so it has to succeed on first call
	c.options.Retryer = nil
	// same client, new request
	resp, err = c.Do(context.Background())
	if err != nil {
		t.Errorf("Expected second request to succeed since the skew should be set on the client. Got %v and err %v", resp, err)
	}
}

func resetDefaults(restoreSleepFunc func()) {
	sdk.NowTime = time.Now
	restoreSleepFunc()
}
