package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/smithy-go"
)

// TestWrapAsClockSkew verifies clock skew error classification per the Clock
// Skew Correction / Retry Behavior SEPs: a known clock skew code is retryable
// only when a candidate skew was observed from the response Date header and
// |attemptSkew - candidateSkew| exceeds the 4 minute detection threshold,
// i.e. the signing time (now + attemptSkew) diverges from the server time
// (now + candidateSkew) by more than the threshold.
func TestWrapAsClockSkew(t *testing.T) {
	apiErr := func(code string) error {
		return &smithy.GenericAPIError{Code: code, Message: "x"}
	}

	above := 5 * time.Minute
	below := 2 * time.Minute

	cases := map[string]struct {
		err              error
		candidateSkew    time.Duration
		hasCandidateSkew bool
		attemptSkew      time.Duration
		wantRetryable    bool
	}{
		"InvalidSignatureException above threshold": {apiErr("InvalidSignatureException"), above, true, 0, true},
		"SignatureDoesNotMatch above threshold":     {apiErr("SignatureDoesNotMatch"), above, true, 0, true},
		"AuthFailure above threshold":               {apiErr("AuthFailure"), above, true, 0, true},
		"RequestTimeTooSkewed above threshold":      {apiErr("RequestTimeTooSkewed"), above, true, 0, true},
		"AccessDeniedException above threshold":     {apiErr("AccessDeniedException"), above, true, 0, true},

		"negative skew above threshold": {apiErr("InvalidSignatureException"), -above, true, 0, true},
		"skew below threshold":          {apiErr("InvalidSignatureException"), below, true, 0, false},
		"skew exactly at threshold":     {apiErr("InvalidSignatureException"), skewThreshold, true, 0, false},
		"no candidate skew":             {apiErr("InvalidSignatureException"), 0, false, 0, false},
		"unknown code above threshold":  {apiErr("ValidationException"), above, true, 0, false},
		"non-api error above threshold": {errors.New("boom"), above, true, 0, false},

		// removed from the set by the SEP; must no longer be treated as skew
		"RequestExpired not skew":     {apiErr("RequestExpired"), above, true, 0, false},
		"RequestInTheFuture not skew": {apiErr("RequestInTheFuture"), above, true, 0, false},

		// stale offset healing: |attemptSkew - candidateSkew| > 4min
		"stale offset heals when clocks realign":    {apiErr("RequestTimeTooSkewed"), 0, true, above, true},
		"stale offset with low candidate not enough": {apiErr("InvalidSignatureException"), below, true, above, false},
		"stale offset no date header":               {apiErr("RequestTimeTooSkewed"), 0, false, above, false},
		"stale offset small attempt skew":           {apiErr("RequestTimeTooSkewed"), 0, true, below, false},
		"stale offset unknown code":                 {apiErr("ValidationException"), 0, true, above, false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := wrapAsClockSkew(c.err, c.candidateSkew, c.hasCandidateSkew, c.attemptSkew)

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
