package retry

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

var _ aws.Retryer = (*Standard)(nil)

func TestStandard_copyOptions(t *testing.T) {
	origDefaultRetryables := DefaultRetryables
	defer func() {
		DefaultRetryables = origDefaultRetryables
	}()
	DefaultRetryables = append([]IsErrorRetryable{}, DefaultRetryables...)

	origDefaultTimeouts := DefaultTimeouts
	defer func() {
		DefaultTimeouts = origDefaultTimeouts
	}()
	DefaultTimeouts = append([]IsErrorTimeout{}, DefaultTimeouts...)

	a := NewStandard(func(ao *StandardOptions) {
		ao.Retryables[0] = nil
		ao.Timeouts[0] = nil
	})

	if DefaultRetryables[0] == nil {
		t.Errorf("expect no change to global var")
	}

	if a.options.Retryables[0] != nil {
		t.Errorf("expect retryables to be changed")
	}

	if DefaultTimeouts[0] == nil {
		t.Errorf("expect no change to global var")
	}

	if a.options.Timeouts[0] != nil {
		t.Errorf("expect timeouts to be changed")
	}
}

func TestStandard_IsErrorRetryable(t *testing.T) {
	cases := map[string]struct {
		Retryable IsErrorRetryable
		Err       error
		Expect    bool
	}{
		"is retryable": {
			Expect: true,
			Err:    fmt.Errorf("expected error"),
			Retryable: IsErrorRetryableFunc(
				func(error) aws.Ternary {
					return aws.TrueTernary
				}),
		},
		"is not retryable": {
			Expect: false,
			Err:    fmt.Errorf("expected error"),
			Retryable: IsErrorRetryableFunc(
				func(error) aws.Ternary {
					return aws.FalseTernary
				}),
		},
		"unknown retryable": {
			Expect: false,
			Err:    fmt.Errorf("expected error"),
			Retryable: IsErrorRetryableFunc(
				func(error) aws.Ternary {
					return aws.UnknownTernary
				}),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := NewStandard(func(o *StandardOptions) {
				o.Retryables = []IsErrorRetryable{
					IsErrorRetryableFunc(
						func(err error) aws.Ternary {
							if e, a := c.Err, err; e != a {
								t.Fatalf("expect %v, error, got %v", e, a)
							}
							return c.Retryable.IsErrorRetryable(err)
						}),
				}
			})
			if e, a := c.Expect, r.IsErrorRetryable(c.Err); e != a {
				t.Errorf("expect %t retryable, got %t", e, a)
			}
		})
	}
}

func TestStandard_MaxAttempts(t *testing.T) {
	cases := map[string]struct {
		Max    int
		Expect int
	}{
		"defaults": {
			Expect: 3,
		},
		"custom": {
			Max:    10,
			Expect: 10,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := NewStandard(func(o *StandardOptions) {
				if c.Max != 0 {
					o.MaxAttempts = c.Max
				} else {
					c.Max = o.MaxAttempts
				}
			})
			if e, a := c.Max, r.MaxAttempts(); e != a {
				t.Errorf("expect %v max, got %v", e, a)
			}
		})
	}
}

func TestStandard_RetryDelay(t *testing.T) {
	cases := map[string]struct {
		Backoff     BackoffDelayer
		Attempt     int
		Err         error
		Assert      func(*testing.T, time.Duration, error)
		ExpectDelay time.Duration
		ExpectErr   error
	}{
		"success": {
			Attempt:     2,
			Err:         fmt.Errorf("expected error"),
			ExpectDelay: 10 * time.Millisecond,

			Backoff: BackoffDelayerFunc(
				func(attempt int, err error) (time.Duration, error) {
					return 10 * time.Millisecond, nil
				}),
		},
		"error": {
			Attempt:     2,
			Err:         fmt.Errorf("expected error"),
			ExpectDelay: 0,
			ExpectErr:   fmt.Errorf("failed get delay"),
			Backoff: BackoffDelayerFunc(
				func(attempt int, err error) (time.Duration, error) {
					return 0, fmt.Errorf("failed get delay")
				}),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := NewStandard(func(o *StandardOptions) {
				o.Backoff = BackoffDelayerFunc(
					func(attempt int, err error) (time.Duration, error) {
						if e, a := c.Err, err; e != a {
							t.Errorf("expect %v error, got %v", e, a)
						}
						if e, a := c.Attempt, attempt; e != a {
							t.Errorf("expect %v attempt, got %v", e, a)
						}
						return c.Backoff.BackoffDelay(attempt, err)
					})
			})

			delay, err := r.RetryDelay(c.Attempt, c.Err)
			if c.ExpectErr != nil {
				if e, a := c.ExpectErr.Error(), err.Error(); e != a {
					t.Errorf("expect %v error, got %v", e, a)
				}
			} else {
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
			}

			if e, a := c.ExpectDelay, delay; e != a {
				t.Errorf("expect %v delay, got %v", e, a)
			}
		})
	}
}

func TestStandard_retryEventuallySucceeds(t *testing.T) {
	ratelimiter := ratelimit.NewTokenRateLimit(500)
	retryer := NewStandard(func(o *StandardOptions) {
		o.MaxAttempts = 3
		o.RateLimiter = ratelimiter
	})
	// Use the default backoff, with a fixed jitter.
	retryer.backoff.(*ExponentialJitterBackoff).randFloat64 = func() (float64, error) {
		return 0.5, nil
	}

	attempts := []struct {
		attemptNum       int
		responseErr      error
		expectRetryable  bool
		expectRetryQuota uint
		expectDelay      time.Duration
	}{
		{
			attemptNum:       1,
			responseErr:      newStubResponseError(500),
			expectRetryable:  true,
			expectRetryQuota: 486,
			expectDelay:      25 * time.Millisecond,
		},
		{
			attemptNum:       2,
			responseErr:      newStubResponseError(500),
			expectRetryable:  true,
			expectRetryQuota: 472,
			expectDelay:      50 * time.Millisecond,
		},
		{
			attemptNum:       3,
			expectRetryQuota: 486,
			// Refund 14 cost of retry, no increment for a retry's success
		},
		{
			attemptNum:       1,
			responseErr:      newStubResponseError(500),
			expectRetryable:  true,
			expectRetryQuota: 472,
			expectDelay:      25 * time.Millisecond,
		},
		{
			attemptNum:       2,
			responseErr:      newStubTimeoutError(),
			expectRetryable:  true,
			expectRetryQuota: 458,
			expectDelay:      50 * time.Millisecond,
			// Timeouts cost the same as other errors
		},
		{
			attemptNum:       3,
			expectRetryQuota: 472,
			// Refund 14 cost of retry
		},
		{
			attemptNum:       1,
			expectRetryQuota: 473,
			// Increment 1 for a first attempt success
		},
	}

	retryToken := nopRelease
	for i, attempt := range attempts {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			ctx := setRetryMetadata(context.Background(), retryMetadata{AttemptNum: attempt.attemptNum})
			attemptToken, err := retryer.GetAttemptToken(ctx)
			if err != nil {
				t.Fatalf("failed to get attempt token, %v", err)
			}

			if err := retryToken(attempt.responseErr); err != nil {
				t.Fatalf("expect release retry token not to fail, %v", err)
			}
			if err := attemptToken(attempt.responseErr); err != nil {
				t.Fatalf("expect release attempt token not to fail, %v", err)
			}

			retryToken = nopRelease
			if attempt.responseErr != nil {
				isRetryable := retryer.IsErrorRetryable(attempt.responseErr)
				if e, a := attempt.expectRetryable, isRetryable; e != a {
					t.Errorf("expect %v retryable, got %v", e, a)
				}

				if isRetryable {
					retryToken, err = retryer.GetRetryToken(ctx, attempt.responseErr)
					if err != nil {
						t.Fatalf("expect get retry token not to fail, %v", err)
					}
				}

				retryDelay, err := retryer.RetryDelay(attempt.attemptNum-1, attempt.responseErr)
				if err != nil {
					t.Fatalf("expect no retry delay error, got %v", err)
				}
				if e, a := attempt.expectDelay, retryDelay; e != a {
					t.Fatalf("expect %v retry delay, got %v", e, a)
				}
			}

			if e, a := attempt.expectRetryQuota, ratelimiter.Remaining(); e != a {
				t.Errorf("expect %v remaining tokens, got %v", e, a)
			}
		})
	}
}

func TestStandard_defaultRetryQuota(t *testing.T) {
	r := NewStandard()
	opErr := fmt.Errorf("generic error")

	var count int
	for i := 0; i < 1000; i++ {
		if _, err := r.GetRetryToken(context.Background(), opErr); err != nil {
			var qe ratelimit.QuotaExceededError
			if !errors.As(err, &qe) {
				t.Fatalf("expect QuotaExceededError, got %v", err)
			}
			break
		}
		count++
	}

	// 500 tokens at 14 per retry.
	if e, a := 35, count; e != a {
		t.Errorf("expect %v retry tokens granted, got %v", e, a)
	}
}

func TestStandard_defaultRetryCosts(t *testing.T) {
	cases := map[string]struct {
		Err        error
		ExpectCost uint
	}{
		"generic":   {Err: fmt.Errorf("generic error"), ExpectCost: 14},
		"timeout":   {Err: newStubTimeoutError(), ExpectCost: 14},
		"5xx":       {Err: newStubResponseError(500), ExpectCost: 14},
		"429":       {Err: newStubResponseError(429), ExpectCost: 5},
		"throttle":  {Err: &mockErrorCodeError{code: "ThrottlingException"}, ExpectCost: 5},
		"throttle2": {Err: &mockErrorCodeError{code: "SlowDown"}, ExpectCost: 5},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			rl := ratelimit.NewTokenRateLimit(500)
			r := NewStandard(func(o *StandardOptions) {
				o.RateLimiter = rl
			})
			release, err := r.GetRetryToken(context.Background(), c.Err)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := 500-c.ExpectCost, rl.Remaining(); e != a {
				t.Errorf("expect %v remaining tokens, got %v", e, a)
			}
			if err := release(nil); err != nil {
				t.Fatalf("expect no release error, got %v", err)
			}
			if e, a := uint(500), rl.Remaining(); e != a {
				t.Errorf("expect %v remaining tokens after refund, got %v", e, a)
			}
		})
	}
}

func TestStandard_GetAttemptToken_retryAttemptNoIncrement(t *testing.T) {
	rl := ratelimit.NewTokenRateLimit(500)
	if _, err := rl.GetToken(context.Background(), 10); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	r := NewStandard(func(o *StandardOptions) {
		o.RateLimiter = rl
	})

	ctx := setRetryMetadata(context.Background(), retryMetadata{AttemptNum: 2})
	release, err := r.GetAttemptToken(ctx)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if err := release(nil); err != nil {
		t.Fatalf("expect no release error, got %v", err)
	}
	if e, a := uint(490), rl.Remaining(); e != a {
		t.Errorf("expect %v remaining tokens, got %v", e, a)
	}

	ctx = setRetryMetadata(context.Background(), retryMetadata{AttemptNum: 1})
	release, err = r.GetAttemptToken(ctx)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if err := release(nil); err != nil {
		t.Fatalf("expect no release error, got %v", err)
	}
	if e, a := uint(491), rl.Remaining(); e != a {
		t.Errorf("expect %v remaining tokens, got %v", e, a)
	}
}

func newStubResponseError(statusCode int) *smithyhttp.ResponseError {
	return &smithyhttp.ResponseError{
		Response: &smithyhttp.Response{
			Response: &http.Response{
				StatusCode: statusCode,
				Header:     http.Header{},
				Body:       io.NopCloser(bytes.NewReader(nil)),
			},
		},
	}
}

func newStubTimeoutError() *mockTimeoutErr {
	return &mockTimeoutErr{timeout: true}
}
