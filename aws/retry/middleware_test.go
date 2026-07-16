package retry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	awsmiddle "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
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
			_, _, err := retryMiddleware.HandleFinalize(ctx, tt.input, middleware.FinalizeHandlerFunc(
				func(ctx context.Context, in middleware.FinalizeInput) (
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

type testRequest struct {
	DisableRewind bool
}

func (r testRequest) RewindStream() error {
	if r.DisableRewind {
		return fmt.Errorf("rewind disabled")
	}
	return nil
}

func TestAttemptMiddleware(t *testing.T) {
	restoreSleep := sdk.TestingUseNopSleep()
	defer restoreSleep()

	sdkTime := sdk.NowTime
	defer func() {
		sdk.NowTime = sdkTime
	}()

	cases := map[string]struct {
		Request       testRequest
		Next          func(retries *[]retryMetadata) middleware.FinalizeHandler
		Expect        []retryMetadata
		Err           error
		ExpectResults AttemptResults
	}{
		"no error, no response in a single attempt": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				return middleware.FinalizeHandlerFunc(
					func(ctx context.Context, in middleware.FinalizeInput) (
						out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
					) {
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
			ExpectResults: AttemptResults{Results: []AttemptResult{
				{},
			}},
		},
		"no error in a single attempt": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				return middleware.FinalizeHandlerFunc(
					func(ctx context.Context, in middleware.FinalizeInput) (
						out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
					) {
						m, ok := getRetryMetadata(ctx)
						if ok {
							*retries = append(*retries, m)
						}
						setMockRawResponse(&metadata, "mockResponse")
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
			ExpectResults: AttemptResults{Results: []AttemptResult{
				{
					ResponseMetadata: func() middleware.Metadata {
						m := middleware.Metadata{}
						setMockRawResponse(&m, "mockResponse")
						return m
					}(),
				},
			}},
		},
		"retries errors": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				num := 0
				reqsErrs := []error{
					mockRetryableError{b: true},
					mockRetryableError{b: true},
					nil,
				}
				return middleware.FinalizeHandlerFunc(
					func(ctx context.Context, in middleware.FinalizeInput) (
						out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
					) {
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
					AttemptNum: 2,
					// note that here and everywhere else, time goes up two
					// additional minutes because of the metrics calling
					// sdk.NowTime twice
					AttemptTime: time.Date(2020, 8, 19, 10, 23, 30, 0, time.UTC),
					MaxAttempts: 3,
				},
				{
					AttemptNum:  3,
					AttemptTime: time.Date(2020, 8, 19, 10, 26, 30, 0, time.UTC),
					MaxAttempts: 3,
				},
			},
			ExpectResults: AttemptResults{Results: []AttemptResult{
				{
					Err:       mockRetryableError{b: true},
					Retryable: true,
					Retried:   true,
				},
				{
					Err:       mockRetryableError{b: true},
					Retryable: true,
					Retried:   true,
				},
				{},
			}},
		},
		"stops after max attempts": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				num := 0
				reqsErrs := []error{
					mockRetryableError{b: true},
					mockRetryableError{b: true},
					mockRetryableError{b: true},
				}
				return middleware.FinalizeHandlerFunc(
					func(ctx context.Context, in middleware.FinalizeInput) (
						out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
					) {
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
			ExpectResults: AttemptResults{Results: []AttemptResult{
				{
					Err:       mockRetryableError{b: true},
					Retryable: true,
					Retried:   true,
				},
				{
					Err:       mockRetryableError{b: true},
					Retryable: true,
					Retried:   true,
				},
				{
					Err:       &MaxAttemptsError{Attempt: 3, Err: mockRetryableError{b: true}},
					Retryable: true,
				},
			}},
		},
		"stops on rewind error": {
			Request: testRequest{DisableRewind: true},
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				return middleware.FinalizeHandlerFunc(
					func(ctx context.Context, in middleware.FinalizeInput) (
						out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
					) {
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
			ExpectResults: AttemptResults{Results: []AttemptResult{
				{
					Err:       mockRetryableError{b: true},
					Retryable: true,
					Retried:   true,
				},
				{
					Err: fmt.Errorf(
						"failed to rewind transport stream for retry, %w",
						fmt.Errorf("rewind disabled"),
					),
				},
			}},
		},
		"stops on non-retryable errors": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				return middleware.FinalizeHandlerFunc(
					func(ctx context.Context, in middleware.FinalizeInput) (
						out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
					) {
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
			ExpectResults: AttemptResults{Results: []AttemptResult{
				{
					Err: fmt.Errorf("some error"),
				},
			}},
		},
		"nested metadata and valid response": {
			Next: func(retries *[]retryMetadata) middleware.FinalizeHandler {
				num := 0
				reqsErrs := []error{
					mockRetryableError{b: true},
					nil,
				}
				return middleware.FinalizeHandlerFunc(
					func(ctx context.Context, in middleware.FinalizeInput) (
						out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
					) {
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

						if err != nil {
							metadata.Set("testKey", "testValue")
						} else {
							setMockRawResponse(&metadata, "mockResponse")
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
					AttemptTime: time.Date(2020, 8, 19, 10, 23, 30, 0, time.UTC),
					MaxAttempts: 3,
				},
			},
			ExpectResults: AttemptResults{Results: []AttemptResult{
				{
					Err:       mockRetryableError{b: true},
					Retryable: true,
					Retried:   true,
					ResponseMetadata: func() middleware.Metadata {
						m := middleware.Metadata{}
						m.Set("testKey", "testValue")
						return m
					}(),
				},
				{
					ResponseMetadata: func() middleware.Metadata {
						m := middleware.Metadata{}
						setMockRawResponse(&m, "mockResponse")
						return m
					}(),
				},
			}},
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
			_, metadata, err := am.HandleFinalize(context.Background(),
				middleware.FinalizeInput{
					Request: tt.Request,
				},
				tt.Next(&recorded),
			)
			if err != nil && tt.Err == nil {
				t.Errorf("expect no error, got %v", err)
			} else if err == nil && tt.Err != nil {
				t.Errorf("expect error, got none")
			} else if err != nil && tt.Err != nil {
				if !strings.Contains(err.Error(), tt.Err.Error()) {
					t.Errorf("expect %v, got %v", tt.Err, err)
				}
			}
			if diff := cmpDiff(recorded, tt.Expect); len(diff) > 0 {
				t.Error(diff)
			}

			attemptResults, ok := GetAttemptResults(metadata)
			if !ok {
				t.Fatalf("expected metadata to contain attempt results, got none")
			}
			if e, a := tt.ExpectResults, attemptResults; !reflect.DeepEqual(e, a) {
				t.Fatalf("expected %v, got %v", e, a)
			}

			for i, attempt := range attemptResults.Results {
				_, ok := GetAttemptResults(attempt.ResponseMetadata)
				if ok {
					t.Errorf("expect no attempt to include AttemptResults metadata, %v does, %#v", i, attempt)
				}
			}
		})
	}
}

func TestAttemptReleaseRetryLock(t *testing.T) {
	standard := NewStandard(func(s *StandardOptions) {
		s.MaxAttempts = 3
		s.RateLimiter = ratelimit.NewTokenRateLimit(10)
		s.RetryCost = 10
	})
	am := NewAttemptMiddleware(standard, func(i interface{}) interface{} {
		return i
	})
	f := func(retries *[]retryMetadata) middleware.FinalizeHandler {
		num := 0
		return middleware.FinalizeHandlerFunc(
			func(ctx context.Context, in middleware.FinalizeInput) (
				out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
			) {
				m, ok := getRetryMetadata(ctx)
				if ok {
					*retries = append(*retries, m)
				}
				if num > 0 {
					return out, metadata, err
				}
				num++
				return out, metadata, mockRetryableError{b: true}
			})
	}
	var recorded []retryMetadata
	_, _, err := am.HandleFinalize(context.Background(), middleware.FinalizeInput{}, f(&recorded))
	if err != nil {
		t.Fatal(err)
	}
	_, err = standard.GetRetryToken(context.Background(), errors.New("retryme"))
	if err != nil {
		t.Fatal(err)
	}
}

type errorCodeImplementer struct {
	errorCode string
}

func (e errorCodeImplementer) Error() string {
	return e.errorCode
}

func (e errorCodeImplementer) ErrorCode() string {
	return e.errorCode
}

func TestClockSkew(t *testing.T) {
	restoreSleep := sdk.TestingUseNopSleep()
	defer restoreSleep()

	cases := map[string]struct {
		clientSkew   time.Duration
		observedSkew time.Duration
		err          error
		disabled     bool
		shouldRetry  bool
	}{
		"no skew and any error no retry": {
			clientSkew:   time.Duration(0),
			observedSkew: time.Duration(0),
			err:          fmt.Errorf("any error"),
			shouldRetry:  false,
		},
		"no skew wrong error code no retry": {
			clientSkew:   time.Duration(0),
			observedSkew: time.Duration(0),
			err:          errorCodeImplementer{"any"},
			shouldRetry:  false,
		},
		"skewed retryable error code does retry": {
			clientSkew:   0,
			observedSkew: 5 * time.Minute,
			err:          errorCodeImplementer{"SignatureDoesNotMatch"},
			shouldRetry:  true,
		},
		"low skew retryable error code no retry": {
			clientSkew:   0,
			observedSkew: 3 * time.Minute,
			err:          errorCodeImplementer{"SignatureDoesNotMatch"},
			shouldRetry:  false,
		},
		"stale offset retries when candidate diverges from server": {
			clientSkew:   5 * time.Minute,
			observedSkew: 0,
			err:          errorCodeImplementer{"SignatureDoesNotMatch"},
			shouldRetry:  true,
		},
		"matching offset does not retry": {
			clientSkew:   5 * time.Minute,
			observedSkew: 5 * time.Minute,
			err:          errorCodeImplementer{"SignatureDoesNotMatch"},
			shouldRetry:  false,
		},
		"disabled does not retry skew error": {
			clientSkew:   0,
			observedSkew: 5 * time.Minute,
			err:          errorCodeImplementer{"SignatureDoesNotMatch"},
			disabled:     true,
			shouldRetry:  false,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			skew := &atomic.Int64{}
			skew.Store(tt.clientSkew.Nanoseconds())
			am := NewAttemptMiddleware(NewStandard(func(s *StandardOptions) {
			}), func(i any) any { return i }, func(m *Attempt) {
				m.ClientSkew = skew
				m.DisableClockSkewCorrection = tt.disabled
			})
			ctx := context.Background()
			_, metadata, err := am.HandleFinalize(ctx, middleware.FinalizeInput{}, middleware.FinalizeHandlerFunc(
				func(ctx context.Context, in middleware.FinalizeInput) (
					out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
				) {
					return out, observedSkewMetadata(t, tt.observedSkew), tt.err
				}))
			if err == nil {
				t.Errorf("Exected return from next middleware, got none")
			}
			attemptResults, ok := GetAttemptResults(metadata)
			if !ok {
				t.Errorf("Got no expected results from metadata. Metadata was %v", metadata)
			}
			if len(attemptResults.Results) == 0 {
				t.Errorf("Expected attempt results, got no results. Attempt was %v", attemptResults)
			}
			wasRetried := attemptResults.Results[0].Retried
			if tt.shouldRetry && !wasRetried {
				t.Errorf("Expected retries, found none %v", attemptResults.Results)
			}
			if !tt.shouldRetry && wasRetried {
				t.Errorf("Expected retries, found none. Results %v", attemptResults.Results)
			}
		})
	}
}

// observedSkewMetadata returns response metadata populated as the
// RecordResponseTiming middleware would for a response whose Date header
// implies the given skew, so the retry layer observes that skew.
func observedSkewMetadata(t *testing.T, skew time.Duration) middleware.Metadata {
	t.Helper()

	responseAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	restore := sdk.TestingUseReferenceTime(responseAt)
	defer restore()

	dateStr := responseAt.Add(skew).UTC().Format(http.TimeFormat)
	_, md, _ := awsmiddle.RecordResponseTiming{}.HandleDeserialize(
		context.Background(),
		middleware.DeserializeInput{},
		middleware.DeserializeHandlerFunc(func(ctx context.Context, in middleware.DeserializeInput) (
			out middleware.DeserializeOutput, m middleware.Metadata, err error,
		) {
			out.RawResponse = &smithyhttp.Response{Response: &http.Response{
				StatusCode: 200,
				Header:     http.Header{"Date": []string{dateStr}},
			}}
			return out, m, nil
		}),
	)

	return md
}

// TestLegacyRetryDelayAttemptValue verifies that when the new retries flag is
// NOT set, the attempt value passed to RetryDelay is unchanged from legacy
// behavior (1-based: first retry passes 1, second passes 2, etc.).
func TestLegacyRetryDelayAttemptValue(t *testing.T) {
	restoreSleep := sdk.TestingUseNopSleep()
	defer restoreSleep()

	var recordedAttempts []int
	retryer := NewStandard(func(o *StandardOptions) {
		o.MaxAttempts = 4
		o.Backoff = BackoffDelayerFunc(func(attempt int, err error) (time.Duration, error) {
			recordedAttempts = append(recordedAttempts, attempt)
			return 0, nil
		})
	})

	am := NewAttemptMiddleware(retryer, func(i interface{}) interface{} {
		return i
	})

	num := 0
	handler := middleware.FinalizeHandlerFunc(
		func(ctx context.Context, in middleware.FinalizeInput) (
			out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
		) {
			num++
			if num < 4 {
				return out, metadata, mockRetryableError{b: true}
			}
			return out, metadata, nil
		})

	_, _, err := am.HandleFinalize(context.Background(), middleware.FinalizeInput{}, handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Legacy behavior: attemptNum is passed directly (1-based).
	// After attempt 1 fails, RetryDelay is called with attemptNum=1.
	// After attempt 2 fails, RetryDelay is called with attemptNum=2.
	// After attempt 3 fails, RetryDelay is called with attemptNum=3.
	expected := []int{1, 2, 3}
	if len(recordedAttempts) != len(expected) {
		t.Fatalf("expected %d RetryDelay calls, got %d", len(expected), len(recordedAttempts))
	}
	for i, exp := range expected {
		if recordedAttempts[i] != exp {
			t.Errorf("RetryDelay call %d: expected attempt=%d, got attempt=%d", i, exp, recordedAttempts[i])
		}
	}
}

// mockRawResponseKey is used to test the behavior when response metadata is
// nested within the attempt request.
type mockRawResponseKey struct{}

func setMockRawResponse(m *middleware.Metadata, v interface{}) {
	m.Set(mockRawResponseKey{}, v)
}

func cmpDiff(e, a interface{}) string {
	if !reflect.DeepEqual(e, a) {
		return fmt.Sprintf("%v != %v", e, a)
	}
	return ""
}
