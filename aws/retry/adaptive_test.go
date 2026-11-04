package retry

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
)

func TestAdaptiveMode_defaultOptions(t *testing.T) {
	a := NewAdaptiveMode()

	s, ok := a.retryer.(*Standard)
	if !ok || s == nil {
		t.Fatalf("expect nested retryer %T, got none", s)
	}

	if e, a := false, a.options.FailOnNoAttemptTokens; e != a {
		t.Errorf("expect %v default fast fail, got %v", e, a)
	}

	if e, a := DefaultMaxAttempts, s.options.MaxAttempts; e != a {
		t.Errorf("expect %v default max attempts, got %v", e, a)
	}
}

func TestAdaptiveMode_customOptions(t *testing.T) {
	a := NewAdaptiveMode(func(ao *AdaptiveModeOptions) {
		ao.FailOnNoAttemptTokens = true
		ao.StandardOptions = append(ao.StandardOptions, func(so *StandardOptions) {
			so.MaxAttempts = 10
		})
	})

	s, ok := a.retryer.(*Standard)
	if !ok || s == nil {
		t.Fatalf("expect nested retryer %T, got none", s)
	}

	if e, a := true, a.options.FailOnNoAttemptTokens; e != a {
		t.Errorf("expect %v custom fast fail, got %v", e, a)
	}

	if e, a := 10, s.options.MaxAttempts; e != a {
		t.Errorf("expect %v custom max attempts, got %v", e, a)
	}
}

func TestAdaptiveMode_copyOptions(t *testing.T) {
	origDefaultThrottles := DefaultThrottles
	defer func() {
		DefaultThrottles = origDefaultThrottles
	}()
	DefaultThrottles = append([]IsErrorThrottle{}, DefaultThrottles...)

	a := NewAdaptiveMode(func(ao *AdaptiveModeOptions) {
		ao.Throttles[0] = nil
	})

	if DefaultThrottles[0] == nil {
		t.Errorf("expect no change to global var")
	}

	if a.options.Throttles[0] != nil {
		t.Errorf("expect throttles to be changed")
	}
}

func TestAdaptiveMode_attemptTokenRefillsRetryQuota(t *testing.T) {
	cases := map[string]struct {
		opErr       error
		expectRetry bool
	}{
		"success": {
			opErr:       nil,
			expectRetry: true,
		},
		"failure": {
			opErr:       newStubResponseError(500),
			expectRetry: false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			retryer := NewAdaptiveMode(func(ao *AdaptiveModeOptions) {
				ao.StandardOptions = append(ao.StandardOptions, func(o *StandardOptions) {
					o.RateLimiter = ratelimit.NewTokenRateLimit(1)
					o.RetryCost = 1
					o.NoRetryIncrement = 1
				})
			})

			// Trigger a failed request to reduce the retry tokens to zero.
			opErr := newStubResponseError(500)
			if _, err := retryer.GetRetryToken(context.Background(), opErr); err != nil {
				t.Fatalf("expect get retry token not to fail, %v", err)
			}

			// Execute processing based on whether the request succeeded or failed
			release, _ := retryer.GetAttemptToken(context.Background())
			if err := release(c.opErr); err != nil {
				t.Fatalf("expect release attempt token not to fail, %v", err)
			}

			// Verify whether the retry token has been restored.
			_, err := retryer.GetRetryToken(context.Background(), opErr)
			if e, a := c.expectRetry, err == nil; e != a {
				t.Errorf("expect retry allowed to be %v, got %v, %v", e, a, err)
			}
		})
	}
}
