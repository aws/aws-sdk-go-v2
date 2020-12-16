package retry

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
	"github.com/google/go-cmp/cmp"
)

func TestMetricsHeaderMiddleware(t *testing.T) {
	cases := []struct {
		input          middleware.FinalizeInput
		ctx            context.Context
		expectedHeader string
		expectedErr    string
	}{
		{
			input: middleware.FinalizeInput{Request: &smithyhttp.Request{Request: &http.Request{Header: make(http.Header)}}},
			ctx: func() context.Context {
				return setRetryMetadata(context.Background(), retryMetadata{
					AttemptNum:       0,
					AttemptTime:      time.Date(2020, 01, 02, 03, 04, 05, 0, time.UTC),
					MaxAttempts:      5,
					AttemptClockSkew: 0,
				})
			}(),
			expectedHeader: "attempt=0; max=5",
		},
		{
			input: middleware.FinalizeInput{Request: &smithyhttp.Request{Request: &http.Request{Header: make(http.Header)}}},
			ctx: func() context.Context {
				attemptTime := time.Date(2020, 01, 02, 03, 04, 05, 0, time.UTC)
				ctx, cancel := context.WithDeadline(context.Background(), attemptTime.Add(time.Minute))
				defer cancel()
				return setRetryMetadata(ctx, retryMetadata{
					AttemptNum:       1,
					AttemptTime:      attemptTime,
					MaxAttempts:      5,
					AttemptClockSkew: time.Second * 1,
				})
			}(),
			expectedHeader: "attempt=1; max=5; ttl=20200102T030506Z",
		},
		{
			ctx: func() context.Context {
				return setRetryMetadata(context.Background(), retryMetadata{})
			}(),
			expectedErr: "unknown transport type",
		},
	}

	retryMiddleware := MetricsHeader{}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := tt.ctx
			_, _, err := retryMiddleware.HandleFinalize(ctx, tt.input, middleware.FinalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (
				out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
			) {
				req := in.Request.(*smithyhttp.Request)

				if e, a := tt.expectedHeader, req.Header.Get("amz-sdk-request"); e != a {
					t.Errorf("expected %v, got %v", e, a)
				}

				return out, metadata, err
			}))
			if err != nil && len(tt.expectedErr) == 0 {
				t.Fatalf("expected no error, got %q", err)
			} else if err != nil && len(tt.expectedErr) != 0 {
				if e, a := tt.expectedErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expected %q, got %q", e, a)
				}
			} else if err == nil && len(tt.expectedErr) != 0 {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

type retryProvider struct {
	Retryer Retryer
}

func (t retryProvider) GetRetryer() Retryer {
	return t.Retryer
}

type mockHandler func(context.Context, interface{}) (interface{}, middleware.Metadata, error)

func (m mockHandler) Handle(ctx context.Context, input interface{}) (output interface{}, metadata middleware.Metadata, err error) {
	return m(ctx, input)
}

func (m mockHandler) ID() string {
	return fmt.Sprintf("%T", m)
}

type testRequest struct {
	DisableRewind bool
}

func (r testRequest) RewindStream() error {
	if r.DisableRewind {
		return fmt.Errorf("rewind disabled")
	}
	return nil
}

type mockRetryableError struct{ b bool }

func (m mockRetryableError) RetryableError() bool { return m.b }
func (m mockRetryableError) Error() string {
	return fmt.Sprintf("mock retryable %t", m.b)
}

func TestAttemptMiddleware(t *testing.T) {
	restoreSleep := sdk.TestingUseNopSleep()
	defer restoreSleep()

	sdkTime := sdk.NowTime
	defer func() {
		sdk.NowTime = sdkTime
	}()

	cases := map[string]struct {
		Request testRequest
		Next    func(retries *[]retryMetadata) middleware.FinalizeHandler
		Expect  []retryMetadata
		Err     error
	}{
		"no error single attempt": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				return middleware.FinalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
					m, ok := getRetryMetadata(ctx)
					if ok {
						*retries = append(*retries, m)
					}
					return out, metadata, err
				})
			},
			Expect: []retryMetadata{
				{
					AttemptNum:  1,
					AttemptTime: time.Date(2020, 8, 19, 10, 20, 30, 0, time.UTC),
					MaxAttempts: 3,
				},
			},
		},
		"retries errors": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				num := 0
				reqsErrs := []error{
					mockRetryableError{b: true},
					mockRetryableError{b: true},
					nil,
				}
				return middleware.FinalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
					m, ok := getRetryMetadata(ctx)
					if ok {
						*retries = append(*retries, m)
					}
					if num >= len(reqsErrs) {
						err = fmt.Errorf("more requests then expected")
					} else {
						err = reqsErrs[num]
						num++
					}
					return out, metadata, err
				})
			},
			Expect: []retryMetadata{
				{
					AttemptNum:  1,
					AttemptTime: time.Date(2020, 8, 19, 10, 20, 30, 0, time.UTC),
					MaxAttempts: 3,
				},
				{
					AttemptNum:  2,
					AttemptTime: time.Date(2020, 8, 19, 10, 21, 30, 0, time.UTC),
					MaxAttempts: 3,
				},
				{
					AttemptNum:  3,
					AttemptTime: time.Date(2020, 8, 19, 10, 22, 30, 0, time.UTC),
					MaxAttempts: 3,
				},
			},
		},
		"stops after max attempts": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				num := 0
				reqsErrs := []error{
					mockRetryableError{b: true},
					mockRetryableError{b: true},
					mockRetryableError{b: true},
				}
				return middleware.FinalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
					if num >= len(reqsErrs) {
						err = fmt.Errorf("more requests then expected")
					} else {
						err = reqsErrs[num]
						num++
					}
					return out, metadata, err
				})
			},
			Err: fmt.Errorf("exceeded maximum number of attempts"),
		},
		"stops on rewind error": {
			Request: testRequest{DisableRewind: true},
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				return middleware.FinalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
					m, ok := getRetryMetadata(ctx)
					if ok {
						*retries = append(*retries, m)
					}
					return out, metadata, mockRetryableError{b: true}
				})
			},
			Expect: []retryMetadata{
				{
					AttemptNum:  1,
					AttemptTime: time.Date(2020, 8, 19, 10, 20, 30, 0, time.UTC),
					MaxAttempts: 3,
				},
			},
			Err: fmt.Errorf("failed to rewind transport stream for retry"),
		},
		"stops on non-retryable errors": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				return middleware.FinalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
					m, ok := getRetryMetadata(ctx)
					if ok {
						*retries = append(*retries, m)
					}
					return out, metadata, fmt.Errorf("some error")
				})
			},
			Expect: []retryMetadata{
				{
					AttemptNum:  1,
					AttemptTime: time.Date(2020, 8, 19, 10, 20, 30, 0, time.UTC),
					MaxAttempts: 3,
				},
			},
			Err: fmt.Errorf("some error"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			sdk.NowTime = func() func() time.Time {
				base := time.Date(2020, 8, 19, 10, 20, 30, 0, time.UTC)
				num := 0
				return func() time.Time {
					t := base.Add(time.Minute * time.Duration(num))
					num++
					return t
				}
			}()

			am := NewAttemptMiddleware(NewStandard(func(s *StandardOptions) {
				s.MaxAttempts = 3
			}), func(i interface{}) interface{} {
				return i
			})

			var recorded []retryMetadata
			_, _, err := am.HandleFinalize(context.Background(), middleware.FinalizeInput{Request: tt.Request}, tt.Next(&recorded))
			if err != nil && tt.Err == nil {
				t.Errorf("expect no error, got %v", err)
			} else if err == nil && tt.Err != nil {
				t.Errorf("expect error, got none")
			} else if err != nil && tt.Err != nil {
				if !strings.Contains(err.Error(), tt.Err.Error()) {
					t.Errorf("expect %v, got %v", tt.Err, err)
				}
			}
			if diff := cmp.Diff(recorded, tt.Expect); len(diff) > 0 {
				t.Error(diff)
			}
		})
	}
}
