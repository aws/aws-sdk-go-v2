package retry

import (
	"math"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/timeconv"
)

func TestExponentialJitterBackoff_AttemptDelay(t *testing.T) {
	maxB := 1 - 1/float64(1<<53)

	cases := map[string]struct {
		MaxBackoff time.Duration
		RandFloat  func() (float64, error)
		Attempt    int
		Expect     time.Duration
	}{
		"min delay floor": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  func() (float64, error) { return 0, nil },
			Attempt:    1,
			Expect:     0,
		},
		"min delay ceiling": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  func() (float64, error) { return maxB, nil },
			Attempt:    1,
			Expect:     timeconv.FloatSecondsDur(maxB * float64(2)),
		},
		"attempt delay": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  func() (float64, error) { return 0.5, nil },
			Attempt:    2,
			Expect:     timeconv.FloatSecondsDur(0.5 * float64(1<<2)),
		},
		"max delay": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  func() (float64, error) { return maxB, nil },
			Attempt:    2147483647,
			Expect:     timeconv.FloatSecondsDur(maxB * math.Exp2(math.Log2(20))),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			j := NewExponentialJitterBackoff(c.MaxBackoff)
			j.randFloat64 = c.RandFloat

			d, err := j.BackoffDelay(c.Attempt, nil)
			if err != nil {
				t.Fatalf("expect not error, %v", err)
			}

			if e, a := c.Expect, d; e != a {
				t.Errorf("expect %v delay, got %v", e, a)
			}
		})
	}
}

func TestExponentialJitterBackoff_retries2026AttemptDelay(t *testing.T) {
	throttleErr := &mockErrorCodeError{code: "Throttling"}

	cases := map[string]struct {
		MaxBackoff time.Duration
		RandFloat  float64
		Attempt    int
		Err        error
		Expect     time.Duration
	}{
		"first retry": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  0.5,
			Attempt:    1,
			Expect:     25 * time.Millisecond,
		},
		"third retry": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  0.5,
			Attempt:    3,
			Expect:     100 * time.Millisecond,
		},
		"max backoff applied before jitter": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  0.5,
			Attempt:    21,
			Expect:     10 * time.Second,
		},
		"throttle uses 1s base": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  0.5,
			Attempt:    1,
			Err:        throttleErr,
			Expect:     500 * time.Millisecond,
		},
		"throttle second retry": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  1,
			Attempt:    2,
			Err:        throttleErr,
			Expect:     2 * time.Second,
		},
		"429 uses 1s base": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  1,
			Attempt:    1,
			Err:        &mockStatusCodeError{code: 429},
			Expect:     1 * time.Second,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			j := newExponentialJitterBackoffWithOptions(c.MaxBackoff,
				withBaseDelay(50*time.Millisecond),
				withThrottleCheck(IsErrorThrottles(DefaultThrottles)),
			)
			j.randFloat64 = func() (float64, error) { return c.RandFloat, nil }

			d, err := j.BackoffDelay(c.Attempt, c.Err)
			if err != nil {
				t.Fatalf("expect not error, %v", err)
			}

			if e, a := c.Expect, d; e != a {
				t.Errorf("expect %v delay, got %v", e, a)
			}
		})
	}
}
