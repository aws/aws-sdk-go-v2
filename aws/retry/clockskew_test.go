package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/smithy-go"
)

// TestWrapAsClockSkew verifies clock skew error classification per the Clock
// Skew Correction / Retry Behavior SEPs: a known clock skew code is retryable
// only when a skew was observed from the response Date header and its absolute
// value exceeds the 4 minute detection threshold.
func TestWrapAsClockSkew(t *testing.T) {
	apiErr := func(code string) error {
		return &smithy.GenericAPIError{Code: code, Message: "x"}
	}

	above := 5 * time.Minute
	below := 2 * time.Minute

	cases := map[string]struct {
		err             error
		observedSkew    time.Duration
		hasObservedSkew bool
		wantRetryable   bool
	}{
		"InvalidSignatureException above threshold": {apiErr("InvalidSignatureException"), above, true, true},
		"SignatureDoesNotMatch above threshold":     {apiErr("SignatureDoesNotMatch"), above, true, true},
		"AuthFailure above threshold":               {apiErr("AuthFailure"), above, true, true},
		"RequestTimeTooSkewed above threshold":      {apiErr("RequestTimeTooSkewed"), above, true, true},
		"AccessDeniedException above threshold":     {apiErr("AccessDeniedException"), above, true, true},

		"negative skew above threshold": {apiErr("InvalidSignatureException"), -above, true, true},
		"skew below threshold":          {apiErr("InvalidSignatureException"), below, true, false},
		"skew exactly at threshold":     {apiErr("InvalidSignatureException"), skewThreshold, true, false},
		"no observed skew":              {apiErr("InvalidSignatureException"), 0, false, false},
		"unknown code above threshold":  {apiErr("ValidationException"), above, true, false},
		"non-api error above threshold": {errors.New("boom"), above, true, false},

		// removed from the set by the SEP; must no longer be treated as skew
		"RequestExpired not skew":     {apiErr("RequestExpired"), above, true, false},
		"RequestInTheFuture not skew": {apiErr("RequestInTheFuture"), above, true, false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := wrapAsClockSkew(c.err, c.observedSkew, c.hasObservedSkew)

			var rcse *retryableClockSkewError
			isRetryable := errors.As(got, &rcse)
			if isRetryable != c.wantRetryable {
				t.Fatalf("retryable = %v, want %v (got %T: %v)", isRetryable, c.wantRetryable, got, got)
			}

			if !c.wantRetryable && got != c.err {
				t.Fatalf("expected original error returned unchanged, got %v", got)
			}

			if c.wantRetryable && !errors.Is(got, c.err) {
				t.Fatalf("wrapped error must unwrap to original")
			}
		})
	}
}

// ensure GenericAPIError satisfies the ErrorCode interface used by the
// classifier (guards against an smithy API change).
var _ interface{ ErrorCode() string } = (*smithy.GenericAPIError)(nil)
