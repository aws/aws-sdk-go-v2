package retry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/smithy-go"
	smithytraits "github.com/aws/smithy-go/traits"
)

// mockTraitError is a modeled error type. Like generated error types, each
// distinct type reports its smithy.api#retryable trait through RetryableTrait.
type mockTraitError struct {
	code   string
	status int
}

func (e *mockTraitError) Error() string                 { return e.code }
func (e *mockTraitError) ErrorCode() string             { return e.code }
func (e *mockTraitError) ErrorMessage() string          { return "" }
func (e *mockTraitError) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }
func (e *mockTraitError) HTTPStatusCode() int           { return e.status }

type mockRetryableTraitError struct{ mockTraitError }

func (*mockRetryableTraitError) RetryableTrait() (smithytraits.Retryable, bool) {
	return smithytraits.Retryable{}, true
}

type mockThrottlingTraitError struct{ mockTraitError }

func (*mockThrottlingTraitError) RetryableTrait() (smithytraits.Retryable, bool) {
	return smithytraits.Retryable{Throttling: true}, true
}

type mockNoTraitError struct{ mockTraitError }

func (*mockNoTraitError) RetryableTrait() (smithytraits.Retryable, bool) {
	return smithytraits.Retryable{}, false
}

func TestRetryableTraitChecks(t *testing.T) {
	cases := map[string]struct {
		Err            error
		ExpectRetry    bool
		ExpectThrottle bool
	}{
		"retryable trait": {
			Err:         &mockRetryableTraitError{mockTraitError{code: "SomeConflict", status: 409}},
			ExpectRetry: true,
		},
		"throttling retryable trait": {
			Err:            &mockThrottlingTraitError{mockTraitError{code: "SlowDownException", status: 400}},
			ExpectRetry:    true,
			ExpectThrottle: true,
		},
		"wrapped throttling retryable trait": {
			Err:            fmt.Errorf("wrapped, %w", &mockThrottlingTraitError{mockTraitError{code: "SlowDownException"}}),
			ExpectRetry:    true,
			ExpectThrottle: true,
		},
		"modeled without trait": {
			Err: &mockNoTraitError{mockTraitError{code: "SomeError", status: 400}},
		},
		"modeled error without RetryableTrait": {
			Err: &mockModeledAPIError{code: "ConcurrentResourcesLimitExceededException", err: &mockStatusCodeError{code: 429}},
		},
		"not an API error": {
			Err: fmt.Errorf("some error"),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.ExpectRetry, (RetryableTrait{}).IsErrorRetryable(c.Err).Bool(); e != a {
				t.Errorf("expect retryable %v, got %v", e, a)
			}
			if e, a := c.ExpectThrottle, (ThrottleTrait{}).IsErrorThrottle(c.Err).Bool(); e != a {
				t.Errorf("expect throttle %v, got %v", e, a)
			}
		})
	}
}

// TestStandard_retryableTrait verifies how the default Standard retryer
// treats modeled errors with the smithy.api#retryable trait.
func TestStandard_retryableTrait(t *testing.T) {
	cases := map[string]struct {
		Err             error
		ExpectRetryable bool
		ExpectCost      uint
		ExpectMaxDelay  time.Duration
	}{
		// A modeled 429 is excluded from the unmodeled-429 rule, so without
		// the trait it is not retried at all.
		"modeled 429 with retryable trait": {
			Err:             &mockRetryableTraitError{mockTraitError{code: "ThrottledRequestException", status: 429}},
			ExpectRetryable: true,
			ExpectCost:      defaultRetryCost,
			ExpectMaxDelay:  50 * time.Millisecond,
		},
		// The error code is not in DefaultThrottleErrorCodes, so throttle
		// handling comes only from the trait.
		"throttling retryable trait with unlisted error code": {
			Err:             &mockThrottlingTraitError{mockTraitError{code: "SlowDownException", status: 400}},
			ExpectRetryable: true,
			ExpectCost:      DefaultThrottlingRetryCost,
			ExpectMaxDelay:  time.Second,
		},
		// No trait: the code list still decides.
		"modeled throttling error code without trait": {
			Err:             &mockNoTraitError{mockTraitError{code: "ThrottlingException", status: 400}},
			ExpectRetryable: true,
			ExpectCost:      DefaultThrottlingRetryCost,
			ExpectMaxDelay:  time.Second,
		},
		"modeled 400 without trait": {
			Err: &mockNoTraitError{mockTraitError{code: "ValidationException", status: 400}},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			rl := ratelimit.NewTokenRateLimit(DefaultRetryRateTokens)
			r := NewStandard(func(o *StandardOptions) { o.RateLimiter = rl })

			if e, a := c.ExpectRetryable, r.IsErrorRetryable(c.Err); e != a {
				t.Fatalf("expect retryable %v, got %v", e, a)
			}
			if !c.ExpectRetryable {
				return
			}

			if _, err := r.GetRetryToken(context.Background(), c.Err); err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := DefaultRetryRateTokens-c.ExpectCost, rl.Remaining(); e != a {
				t.Errorf("expect %v remaining retry tokens, got %v", e, a)
			}

			backoff := r.backoff.(*ExponentialJitterBackoff)
			backoff.randFloat64 = func() (float64, error) { return 1, nil }
			delay, err := r.RetryDelay(1, c.Err)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := c.ExpectMaxDelay, delay; e != a {
				t.Errorf("expect first retry delay %v, got %v", e, a)
			}
		})
	}
}

// TestAdaptiveMode_throttlingTrait verifies the adaptive rate limiter
// treats a modeled throttling error with an unlisted code as a throttle.
func TestAdaptiveMode_throttlingTrait(t *testing.T) {
	a := NewAdaptiveMode()
	err := &mockThrottlingTraitError{mockTraitError{code: "SlowDownException", status: 400}}

	release, tokenErr := a.GetAttemptToken(context.Background())
	if tokenErr != nil {
		t.Fatalf("expect no error, got %v", tokenErr)
	}
	if releaseErr := release(err); releaseErr != nil {
		t.Fatalf("expect no error, got %v", releaseErr)
	}

	a.rateLimit.mu.Lock()
	enabled := a.rateLimit.tokenBucketEnabled
	a.rateLimit.mu.Unlock()
	if !enabled {
		t.Errorf("expect the adaptive token bucket to be enabled after a throttle")
	}
}
