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

/*
func TestRetryThrottleStatusCodes(t *testing.T) {
	cases := []struct {
		expectThrottle bool
		expectRetry    bool
		r              Request
	}{
		{
			false,
			false,
			Request{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
		},
		{
			true,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 429},
			},
		},
		{
			true,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 502},
			},
		},
		{
			true,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 503},
			},
		},
		{
			true,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 504},
			},
		},
		{
			false,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
		},
	}

	d := NewDefaultRetryer(func(d *DefaultRetryer) {
		d.NumMaxRetries = 100
	})
	for i, c := range cases {
		throttle := c.r.IsErrorThrottle()
		retry := d.ShouldRetry(&c.r)

		if e, a := c.expectThrottle, throttle; e != a {
			t.Errorf("%d: expected %v, but received %v", i, e, a)
		}

		if e, a := c.expectRetry, retry; e != a {
			t.Errorf("%d: expected %v, but received %v", i, e, a)
		}
	}
}

func TestGetRetryAfterDelay(t *testing.T) {
	cases := []struct {
		r Request
		e bool
	}{
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			false,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			false,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 429},
			},
			true,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 503},
			},
			true,
		},
	}

	for i, c := range cases {
		a := canUseRetryAfterHeader(&c.r)
		if c.e != a {
			t.Errorf("%d: expected %v, but received %v", i, c.e, a)
		}
	}
}

func TestGetRetryDelay(t *testing.T) {
	cases := []struct {
		r     Request
		e     time.Duration
		equal bool
		ok    bool
	}{
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 429, Header: http.Header{"Retry-After": []string{"3600"}}},
			},
			3600 * time.Second,
			true,
			true,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": []string{"120"}}},
			},
			120 * time.Second,
			true,
			true,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": []string{"120"}}},
			},
			1 * time.Second,
			false,
			true,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": []string{""}}},
			},
			0 * time.Second,
			true,
			false,
		},
	}

	for i, c := range cases {
		a, ok := getRetryAfterDelay(&c.r)
		if c.ok != ok {
			t.Errorf("%d: expected %v, but received %v", i, c.ok, ok)
		}

		if (c.e != a) == c.equal {
			t.Errorf("%d: expected %v, but received %v", i, c.e, a)
		}
	}
}

func TestRetryDelay(t *testing.T) {
	d := NewDefaultRetryer(func(d *DefaultRetryer) {
		d.NumMaxRetries = 100
	})
	r := Request{}
	for i := 0; i < 100; i++ {
		rTemp := r
		rTemp.HTTPResponse = &http.Response{StatusCode: 500, Header: http.Header{"Retry-After": []string{"299"}}}
		rTemp.RetryCount = i
		a := d.RetryRules(&rTemp)
		if a > 5*time.Minute {
			t.Errorf("retry delay should never be greater than five minutes, received %s for retrycount %d", a, i)
		}
	}

	for i := 0; i < 100; i++ {
		rTemp := r
		rTemp.RetryCount = i
		rTemp.HTTPResponse = &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": []string{""}}}
		a := d.RetryRules(&rTemp)
		if a > 5*time.Minute {
			t.Errorf("retry delay should not be greater than five minutes, received %s for retrycount %d", a, i)
		}
	}

	rTemp := r
	rTemp.RetryCount = 1
	rTemp.HTTPResponse = &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": []string{"300"}}}
	a := d.RetryRules(&rTemp)
	if a < 5*time.Minute {
		t.Errorf("retry delay should not be less than retry-after duration, received %s for retrycount %d", a, 1)
	}
}
*/
