package timeouts

import (
	"testing"
	"time"
)

// setGate forces the feature gate and restores it afterward.
func setGate(t *testing.T, on bool) {
	t.Helper()

	prev := enableFromEnv
	enableFromEnv = func() bool { return on }
	t.Cleanup(func() { enableFromEnv = prev })
}

func TestGetServiceReadTimeout(t *testing.T) {
	for name, tt := range map[string]struct {
		serviceID string
		gate      bool
		expect    time.Duration
		expectOK  bool
	}{
		"gate off applies nothing": {
			serviceID: "DynamoDB",
		},
		"unlisted service gets the package default": {
			serviceID: "DynamoDB", gate: true,
			expect: 5 * time.Minute, expectOK: true,
		},
		"long-hold service gets the relaxed value": {
			serviceID: "SQS", gate: true,
			expect: 15 * time.Minute, expectOK: true,
		},
		"fully exempt service gets nothing": {
			serviceID: "S3", gate: true,
		},
		"fully exempt streaming service gets nothing": {
			serviceID: "Transcribe Streaming", gate: true,
		},
		"unknown service gets the package default": {
			serviceID: "Not A Real Service", gate: true,
			expect: 5 * time.Minute, expectOK: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			setGate(t, tt.gate)

			got, ok := GetServiceReadTimeout(tt.serviceID)
			if ok != tt.expectOK {
				t.Fatalf("expect ok %v, got %v", tt.expectOK, ok)
			}
			if ok && got != tt.expect {
				t.Errorf("expect %v, got %v", tt.expect, got)
			}
		})
	}
}
