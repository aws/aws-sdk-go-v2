package retry

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type retryOutcome int

const (
	retryOutcomeSuccess retryOutcome = iota
	retryOutcomeRetryRequest
	retryOutcomeMaxAttemptsExceeded
	retryOutcomeRetryQuotaExceeded
	retryOutcomeFailRequest
)

type retryTestResponse struct {
	statusCode int
	errorCode  string
	headers    map[string]string
}

type retryTestExpected struct {
	outcome    retryOutcome
	retryQuota *uint
	delay      time.Duration
}

type retryTestStep struct {
	response retryTestResponse
	expected retryTestExpected
}

type retryTestCase struct {
	name            string
	maxAttempts     int
	initialTokens   *uint
	exponentialBase float64
	maxBackoffTime  time.Duration
	service         string
	longPolling     bool
	steps           []retryTestStep
}

func uintPtr(v uint) *uint { return &v }

func newRetry2026Retryer(tc retryTestCase) (aws.RetryerV2, *ratelimit.TokenRateLimit) {
	tokens := uint(500)
	if tc.initialTokens != nil {
		tokens = *tc.initialTokens
	}

	rl := ratelimit.NewTokenRateLimit(tokens)

	maxAttempts := 3
	if tc.maxAttempts != 0 {
		maxAttempts = tc.maxAttempts
	}

	maxBackoff := 20 * time.Second
	if tc.maxBackoffTime != 0 {
		maxBackoff = tc.maxBackoffTime
	}

	baseDelay := 50 * time.Millisecond
	if tc.service == "dynamodb" {
		baseDelay = 25 * time.Millisecond
		if tc.maxAttempts == 0 {
			maxAttempts = 4
		}
	}

	backoff := newExponentialJitterBackoffWithOptions(maxBackoff,
		withBaseDelay(baseDelay),
		withThrottleCheck(IsErrorThrottles(DefaultThrottles)),
	)
	if tc.exponentialBase != 0 {
		backoff.randFloat64 = func() (float64, error) {
			return tc.exponentialBase, nil
		}
	}

	r := NewStandard(func(o *StandardOptions) {
		o.MaxAttempts = maxAttempts
		o.MaxBackoff = maxBackoff
		o.RateLimiter = rl
		o.Backoff = backoff
	})

	if tc.longPolling {
		return AddWithLongPolling(r).(aws.RetryerV2), rl
	}
	return r, rl
}

func newRetry2026Error(resp retryTestResponse) error {
	httpResp := &http.Response{
		StatusCode: resp.statusCode,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewReader(nil)),
	}
	for k, v := range resp.headers {
		httpResp.Header.Set(k, v)
	}

	if resp.errorCode != "" {
		return &mockAPIError{
			ResponseError: &smithyhttp.ResponseError{
				Response: &smithyhttp.Response{Response: httpResp},
			},
			code: resp.errorCode,
		}
	}

	return &smithyhttp.ResponseError{
		Response: &smithyhttp.Response{Response: httpResp},
	}
}

type mockAPIError struct {
	*smithyhttp.ResponseError
	code string
}

func (e *mockAPIError) ErrorCode() string    { return e.code }
func (e *mockAPIError) ErrorMessage() string { return e.code }
func (e *mockAPIError) ErrorFault() int      { return 0 }
func (e *mockAPIError) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.ResponseError.Error())
}
func (e *mockAPIError) Unwrap() error { return e.ResponseError }

func TestRetry2026StandardMode(t *testing.T) {
	t.Setenv("AWS_NEW_RETRIES_2026", "true")
	cases := []retryTestCase{
		{
			name:            "retry eventually succeeds",
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(472), delay: 100 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 200},
					expected: retryTestExpected{outcome: retryOutcomeSuccess, retryQuota: uintPtr(486)},
				},
			},
		},
		{
			name:            "fail due to max attempts reached",
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 502},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 502},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(472), delay: 100 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 502},
					expected: retryTestExpected{outcome: retryOutcomeMaxAttemptsExceeded, retryQuota: uintPtr(472)},
				},
			},
		},
		{
			name:            "retry quota reached after a single retry",
			initialTokens:   uintPtr(14),
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(0), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryQuotaExceeded, retryQuota: uintPtr(0)},
				},
			},
		},
		{
			name:            "no retries at all if retry quota is 0",
			initialTokens:   uintPtr(0),
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryQuotaExceeded, retryQuota: uintPtr(0)},
				},
			},
		},
		{
			name:            "verifying exponential backoff timing",
			maxAttempts:     5,
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(472), delay: 100 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(458), delay: 200 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(444), delay: 400 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeMaxAttemptsExceeded, retryQuota: uintPtr(444)},
				},
			},
		},
		{
			name:            "verify max backoff time",
			maxAttempts:     5,
			exponentialBase: 1,
			maxBackoffTime:  200 * time.Millisecond,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(472), delay: 100 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(458), delay: 200 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(444), delay: 200 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeMaxAttemptsExceeded, retryQuota: uintPtr(444)},
				},
			},
		},
		{
			name:            "retry stops after retry quota exhaustion",
			maxAttempts:     5,
			initialTokens:   uintPtr(20),
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(6), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 502},
					expected: retryTestExpected{outcome: retryOutcomeRetryQuotaExceeded, retryQuota: uintPtr(6)},
				},
			},
		},
		{
			name:            "retry quota recovery after successful responses",
			maxAttempts:     5,
			initialTokens:   uintPtr(30),
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(16), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 502},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(2), delay: 100 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 200},
					expected: retryTestExpected{outcome: retryOutcomeSuccess, retryQuota: uintPtr(16)},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(2), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 200},
					expected: retryTestExpected{outcome: retryOutcomeSuccess, retryQuota: uintPtr(16)},
				},
			},
		},
		{
			name:            "throttling error token bucket drain and backoff duration",
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 400, errorCode: "Throttling"},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(495), delay: 1 * time.Second},
				},
				{
					response: retryTestResponse{statusCode: 200},
					expected: retryTestExpected{outcome: retryOutcomeSuccess, retryQuota: uintPtr(500)},
				},
			},
		},
		{
			name:            "dynamodb base backoff and increased retries",
			service:         "dynamodb",
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 25 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(472), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(458), delay: 100 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeMaxAttemptsExceeded, retryQuota: uintPtr(458)},
				},
			},
		},
		{
			name:            "long-polling backoff when token bucket empty",
			service:         "sqs",
			longPolling:     true,
			initialTokens:   uintPtr(0),
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryQuotaExceeded, retryQuota: uintPtr(0), delay: 50 * time.Millisecond},
				},
			},
		},
		{
			name:            "long-polling backoff after throttling error when token bucket empty",
			service:         "sqs",
			longPolling:     true,
			initialTokens:   uintPtr(0),
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 400, errorCode: "Throttling"},
					expected: retryTestExpected{outcome: retryOutcomeRetryQuotaExceeded, retryQuota: uintPtr(0), delay: 50 * time.Millisecond},
				},
			},
		},
		{
			name:            "long-polling max attempts exceeded must not delay",
			service:         "sqs",
			longPolling:     true,
			maxAttempts:     2,
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeMaxAttemptsExceeded, retryQuota: uintPtr(486)},
				},
			},
		},
		{
			name:            "long-polling success must not delay",
			service:         "sqs",
			longPolling:     true,
			maxAttempts:     2,
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 500},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 200},
					expected: retryTestExpected{outcome: retryOutcomeSuccess, retryQuota: uintPtr(500)},
				},
			},
		},
		{
			name:            "long-polling non-retriable errors must not delay",
			service:         "sqs",
			longPolling:     true,
			maxAttempts:     2,
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{statusCode: 404},
					expected: retryTestExpected{outcome: retryOutcomeFailRequest},
				},
			},
		},
		{
			name: "honor x-amz-retry-after header",
			steps: []retryTestStep{
				{
					response: retryTestResponse{
						statusCode: 500,
						headers:    map[string]string{"x-amz-retry-after": "1500"},
					},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 1500 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 200},
					expected: retryTestExpected{outcome: retryOutcomeSuccess, retryQuota: uintPtr(500)},
				},
			},
		},
		{
			name:            "x-amz-retry-after minimum is exponential backoff duration",
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{
						statusCode: 500,
						headers:    map[string]string{"x-amz-retry-after": "0"},
					},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 200},
					expected: retryTestExpected{outcome: retryOutcomeSuccess, retryQuota: uintPtr(500)},
				},
			},
		},
		{
			name:            "x-amz-retry-after maximum is 5+exponential backoff duration",
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{
						statusCode: 500,
						headers:    map[string]string{"x-amz-retry-after": "10000"},
					},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 5050 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 200},
					expected: retryTestExpected{outcome: retryOutcomeSuccess, retryQuota: uintPtr(500)},
				},
			},
		},
		{
			name:            "invalid x-amz-retry-after falls back to exponential backoff",
			exponentialBase: 1,
			steps: []retryTestStep{
				{
					response: retryTestResponse{
						statusCode: 500,
						headers:    map[string]string{"x-amz-retry-after": "invalid"},
					},
					expected: retryTestExpected{outcome: retryOutcomeRetryRequest, retryQuota: uintPtr(486), delay: 50 * time.Millisecond},
				},
				{
					response: retryTestResponse{statusCode: 200},
					expected: retryTestExpected{outcome: retryOutcomeSuccess, retryQuota: uintPtr(500)},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			retryer, rl := newRetry2026Retryer(tc)

			attemptNum := 0
			retryToken := nopRelease

			for i, step := range tc.steps {
				attemptNum++

				attemptToken, err := retryer.GetAttemptToken(context.Background())
				if err != nil {
					t.Fatalf("step %d: failed to get attempt token: %v", i, err)
				}

				var opErr error
				if step.response.statusCode >= 300 {
					opErr = newRetry2026Error(step.response)
				}

				if err := retryToken(opErr); err != nil {
					t.Fatalf("step %d: failed to release retry token: %v", i, err)
				}
				retryToken = nopRelease

				// NoRetryIncrement only applies to the initial attempt.
				if attemptNum == 1 {
					if err := attemptToken(opErr); err != nil {
						t.Fatalf("step %d: failed to release attempt token: %v", i, err)
					}
				}

				switch step.expected.outcome {
				case retryOutcomeSuccess:
					if opErr != nil {
						t.Fatalf("step %d: expected success but got error", i)
					}

				case retryOutcomeRetryRequest:
					if !retryer.IsErrorRetryable(opErr) {
						t.Fatalf("step %d: expected retryable error", i)
					}
					if attemptNum >= retryer.MaxAttempts() {
						t.Fatalf("step %d: expected retry but max attempts reached", i)
					}

					retryToken, err = retryer.GetRetryToken(context.Background(), opErr)
					if err != nil {
						t.Fatalf("step %d: expected retry token, got error: %v", i, err)
					}

					delay, delayErr := retryer.RetryDelay(attemptNum-1, opErr)
					if delayErr != nil {
						t.Fatalf("step %d: unexpected delay error: %v", i, delayErr)
					}
					delay = adjustForRetryAfterHeader(delay, opErr, nil, false)
					if e, a := step.expected.delay, delay; e != a {
						t.Errorf("step %d: expect delay %v, got %v", i, e, a)
					}

				case retryOutcomeMaxAttemptsExceeded:
					if !retryer.IsErrorRetryable(opErr) {
						t.Fatalf("step %d: expected retryable error", i)
					}
					if attemptNum < retryer.MaxAttempts() {
						t.Errorf("step %d: expected max attempts reached at %d, but only at %d", i, retryer.MaxAttempts(), attemptNum)
					}

				case retryOutcomeRetryQuotaExceeded:
					if !retryer.IsErrorRetryable(opErr) {
						t.Fatalf("step %d: expected retryable error", i)
					}
					_, getTokenErr := retryer.GetRetryToken(context.Background(), opErr)
					if getTokenErr == nil {
						t.Fatalf("step %d: expected retry quota exceeded error", i)
					}

					if step.expected.delay != 0 && tc.longPolling {
						// Mirrors the middleware quota-exceeded path:
						// - attemptNum-1: backoff exponent is 0-based (first failure = 2^0).
						// - nil error: forces non-throttle base delay (50ms), even if the
						//   actual error was a throttle. The real error is still passed to
						//   adjustForRetryAfterHeader to honor the response header.
						delay, delayErr := retryer.RetryDelay(attemptNum-1, nil)
						if delayErr != nil {
							t.Fatalf("step %d: unexpected delay error: %v", i, delayErr)
						}
						delay = adjustForRetryAfterHeader(delay, opErr, nil, false)
						if e, a := step.expected.delay, delay; e != a {
							t.Errorf("step %d: expect long-polling delay %v, got %v", i, e, a)
						}
					}

				case retryOutcomeFailRequest:
					// Non-retryable error — should not retry or delay.
					if retryer.IsErrorRetryable(opErr) {
						t.Fatalf("step %d: expected non-retryable error", i)
					}
				}

				if step.expected.retryQuota != nil {
					if e, a := *step.expected.retryQuota, rl.Remaining(); e != a {
						t.Errorf("step %d: expect retry quota %d, got %d", i, e, a)
					}
				}

				if step.expected.outcome == retryOutcomeSuccess {
					attemptNum = 0
				}
			}
		})
	}
}
