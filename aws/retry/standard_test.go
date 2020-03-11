package retry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
)

var _ aws.Retryer = (*retry.Standard)(nil)

func TestStandard_IsErrorRetryable(t *testing.T) {
	cases := map[string]struct {
		Retryable retry.IsErrorRetryable
		Err       error
		Expect    bool
	}{
		"is retryable": {
			Expect: true,
			Err:    fmt.Errorf("expected error"),
			Retryable: retry.IsErrorRetryableFunc(
				func(error) aws.Ternary {
					return aws.TrueTernary
				}),
		},
		"is not retryable": {
			Expect: false,
			Err:    fmt.Errorf("expected error"),
			Retryable: retry.IsErrorRetryableFunc(
				func(error) aws.Ternary {
					return aws.FalseTernary
				}),
		},
		"unknown retryable": {
			Expect: false,
			Err:    fmt.Errorf("expected error"),
			Retryable: retry.IsErrorRetryableFunc(
				func(error) aws.Ternary {
					return aws.UnknownTernary
				}),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := retry.NewStandard(func(o *retry.StandardOptions) {
				o.Retryables = []retry.IsErrorRetryable{
					retry.IsErrorRetryableFunc(
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
			r := retry.NewStandard(func(o *retry.StandardOptions) {
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
		Backoff     retry.BackoffDelayer
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

			Backoff: retry.BackoffDelayerFunc(
				func(attempt int, err error) (time.Duration, error) {
					return 10 * time.Millisecond, nil
				}),
		},
		"error": {
			Attempt:     2,
			Err:         fmt.Errorf("expected error"),
			ExpectDelay: 0,
			ExpectErr:   fmt.Errorf("failed get delay"),
			Backoff: retry.BackoffDelayerFunc(
				func(attempt int, err error) (time.Duration, error) {
					return 0, fmt.Errorf("failed get delay")
				}),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := retry.NewStandard(func(o *retry.StandardOptions) {
				o.Backoff = retry.BackoffDelayerFunc(
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
