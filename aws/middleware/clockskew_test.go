package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// These tests replay the cross-SDK Clock Skew Correction SEP test vectors
// (testdata/clock-skew-test-cases.json) against the SDK's skew computation. For
// each operation it seeds AttemptSkew from ClientSkew, signs each attempt at
// clientTimeAtSend+AttemptSkew (asserting expectedSigningTime), computes the
// candidate skew via computeClockSkew, records it unconditionally when present,
// and asserts the resulting ClientSkew (expectedClientSkew).
type skewVectorFile struct {
	Tests []skewVectorTest `json:"tests"`
}

type skewVectorTest struct {
	Description string         `json:"description"`
	Operations  []skewVectorOp `json:"operations"`
}

type skewVectorOp struct {
	InitialClientSkew  int64               `json:"initialClientSkew"`
	MaxAttempts        int                 `json:"maxAttempts"`
	Attempts           []skewVectorAttempt `json:"attempts"`
	ExpectedClientSkew int64               `json:"expectedClientSkew"`
	ExpectedOutcome    string              `json:"expectedOutcome"`
}

type skewVectorAttempt struct {
	ClientTimeAtSend    string         `json:"clientTimeAtSend"`
	ClientTimeAtReceive string         `json:"clientTimeAtReceive"`
	ExpectedSigningTime string         `json:"expectedSigningTime"`
	Response            skewVectorResp `json:"response"`
}

type skewVectorResp struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	ErrorCode  string            `json:"errorCode"`
}

func TestClockSkewSEPVectors(t *testing.T) {
	raw, err := os.ReadFile("testdata/clock-skew-test-cases.json")
	if err != nil {
		t.Fatalf("read vectors: %v", err)
	}

	var file skewVectorFile
	if err := json.Unmarshal(raw, &file); err != nil {
		t.Fatalf("parse vectors: %v", err)
	}

	if len(file.Tests) == 0 {
		t.Fatal("no test vectors loaded")
	}

	for _, tc := range file.Tests {
		t.Run(tc.Description, func(t *testing.T) {
			var prevExpected int64
			for oi, op := range tc.Operations {
				if oi > 0 && op.InitialClientSkew != prevExpected {
					t.Fatalf("op %d: initialClientSkew %d does not match prior expectedClientSkew %d",
						oi, op.InitialClientSkew, prevExpected)
				}

				clientSkew := time.Duration(op.InitialClientSkew) * time.Second
				attemptSkew := clientSkew

				for ai, at := range op.Attempts {
					send := mustParseISO(t, at.ClientTimeAtSend)
					recv := mustParseISO(t, at.ClientTimeAtReceive)

					signing := send.Add(attemptSkew)
					if want := mustParseISO(t, at.ExpectedSigningTime); !signing.Equal(want) {
						t.Fatalf("op %d attempt %d: signing time = %s, want %s", oi, ai, signing.Format(time.RFC3339), want.Format(time.RFC3339))
					}

					var serverTime time.Time
					if d := at.Response.Headers["Date"]; d != "" {
						st, perr := smithyhttp.ParseTime(d)
						if perr != nil {
							t.Fatalf("op %d attempt %d: parse Date %q: %v", oi, ai, d, perr)
						}
						serverTime = st
					}
					hasAge := at.Response.Headers["Age"] != ""

					if skew, ok := computeClockSkew(serverTime, send, recv, hasAge); ok {
						attemptSkew = skew
						clientSkew = skew
					}
				}

				if want := time.Duration(op.ExpectedClientSkew) * time.Second; clientSkew != want {
					t.Fatalf("op %d: ClientSkew = %s, want %s", oi, clientSkew, want)
				}
				prevExpected = op.ExpectedClientSkew
			}
		})
	}
}

func mustParseISO(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("parse time %q: %v", s, err)
	}

	return v
}

// TestRecordResponseTimingClockSkewDisabled verifies that the disable flag
// suppresses skew recording while a default (enabled) middleware records it.
func TestRecordResponseTimingClockSkewDisabled(t *testing.T) {
	responseAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	restore := sdk.TestingUseReferenceTime(responseAt)
	defer restore()

	dateStr := responseAt.Add(5 * time.Minute).UTC().Format(http.TimeFormat)
	next := smithymiddleware.DeserializeHandlerFunc(func(ctx context.Context, in smithymiddleware.DeserializeInput) (
		out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error,
	) {
		out.RawResponse = &smithyhttp.Response{Response: &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Date": []string{dateStr}},
		}}
		return out, m, nil
	})

	_, mdDisabled, err := RecordResponseTiming{DisableClockSkewCorrection: true}.
		HandleDeserialize(context.Background(), smithymiddleware.DeserializeInput{}, next)
	if err != nil {
		t.Fatalf("disabled: unexpected error: %v", err)
	}
	if skew, ok := GetAttemptSkew(mdDisabled); ok {
		t.Fatalf("disabled: expected no recorded skew, got %v", skew)
	}

	_, mdEnabled, err := RecordResponseTiming{}.
		HandleDeserialize(context.Background(), smithymiddleware.DeserializeInput{}, next)
	if err != nil {
		t.Fatalf("enabled: unexpected error: %v", err)
	}
	if skew, ok := GetAttemptSkew(mdEnabled); !ok || skew != 5*time.Minute {
		t.Fatalf("enabled: expected 5m recorded skew, got %v ok=%v", skew, ok)
	}
}
