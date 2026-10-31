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
			Attempt:    0,
			Expect:     0,
		},
		"min delay ceiling": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  func() (float64, error) { return maxB, nil },
			Attempt:    0,
			Expect:     timeconv.FloatSecondsDur(maxB * float64(2)),
		},
		"attempt delay": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  func() (float64, error) { return 0.5, nil },
			Attempt:    1,
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

// TestExponentialJitterBackoff_legacyDelaysUnchanged asserts that an explicitly
// configured NewExponentialJitterBackoff waits as long through the retry
// middleware as it did before the middleware passed the 0-based retry index.
// Before, the middleware called BackoffDelay(attemptNum, err). Now it calls
// BackoffDelay(attemptNum-1, err).
func TestExponentialJitterBackoff_legacyDelaysUnchanged(t *testing.T) {
	// Legacy formula on the 1-based attempt number: b * 2^attempt, or
	// maxBackoff once attempt passes log2(maxBackoff / 1s).
	legacy := func(b float64, attemptNum int, maxBackoff time.Duration) time.Duration {
		if attemptNum > int(math.Log2(float64(maxBackoff)/float64(time.Second))) {
			return maxBackoff
		}
		return timeconv.FloatSecondsDur(b * float64(int64(1)<<uint64(attemptNum)))
	}

	for _, maxBackoff := range []time.Duration{time.Second, 20 * time.Second} {
		for _, b := range []float64{0, 0.5, 1} {
			for attemptNum := 1; attemptNum <= 8; attemptNum++ {
				j := NewExponentialJitterBackoff(maxBackoff)
				j.randFloat64 = func() (float64, error) { return b, nil }

				d, err := j.BackoffDelay(attemptNum-1, nil)
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
				if e, a := legacy(b, attemptNum, maxBackoff), d; e != a {
					t.Errorf("max %v, b %v, attempt %d: expect %v delay, got %v",
						maxBackoff, b, attemptNum, e, a)
				}
			}
		}
	}

	// Values for the SDK's 20s default max backoff with b = 1.
	expect := []time.Duration{
		2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second,
		20 * time.Second, 20 * time.Second,
	}
	for i, e := range expect {
		j := NewExponentialJitterBackoff(20 * time.Second)
		j.randFloat64 = func() (float64, error) { return 1, nil }
		d, err := j.BackoffDelay(i, nil)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
		if d != e {
			t.Errorf("retry index %d: expect %v delay, got %v", i, e, d)
		}
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
			Attempt:    0,
			Expect:     25 * time.Millisecond,
		},
		"third retry": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  0.5,
			Attempt:    2,
			Expect:     100 * time.Millisecond,
		},
		"max backoff applied before jitter": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  0.5,
			Attempt:    20,
			Expect:     10 * time.Second,
		},
		"throttle uses 1s base": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  0.5,
			Attempt:    0,
			Err:        throttleErr,
			Expect:     500 * time.Millisecond,
		},
		"throttle second retry": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  1,
			Attempt:    1,
			Err:        throttleErr,
			Expect:     2 * time.Second,
		},
		"429 uses 1s base": {
			MaxBackoff: 20 * time.Second,
			RandFloat:  1,
			Attempt:    0,
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
