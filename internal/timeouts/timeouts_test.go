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

func setInternalGate(t *testing.T, on bool, rollout map[string]bool) {
	t.Helper()

	prevEnabled := enableReadTimeout2026
	prevRollout := readTimeout2026Rollout
	enableReadTimeout2026 = on
	readTimeout2026Rollout = rollout
	t.Cleanup(func() {
		enableReadTimeout2026 = prevEnabled
		readTimeout2026Rollout = prevRollout
	})
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

func TestGetServiceReadTimeout_InternalGate(t *testing.T) {
	for name, tt := range map[string]struct {
		serviceID string
		rollout   map[string]bool
		envGate   bool
		expect    time.Duration
		expectOK  bool
	}{
		"service in rollout applies default, envGate irrelevant": {
			serviceID: "DynamoDB",
			rollout:   map[string]bool{"DynamoDB": true},
			envGate:   false,
			expect:    5 * time.Minute, expectOK: true,
		},
		"service in rollout applies default even with envGate on": {
			serviceID: "DynamoDB",
			rollout:   map[string]bool{"DynamoDB": true},
			envGate:   true,
			expect:    5 * time.Minute, expectOK: true,
		},
		"service not in rollout gets nothing": {
			serviceID: "DynamoDB",
			rollout:   map[string]bool{"SQS": true},
			envGate:   false,
		},
		"nil rollout map gets nothing": {
			serviceID: "DynamoDB",
			rollout:   nil,
			envGate:   false,
		},
		"exempt service in rollout still gets nothing": {
			serviceID: "S3",
			rollout:   map[string]bool{"S3": true},
			envGate:   false,
		},
		"long-hold service in rollout gets the relaxed value": {
			serviceID: "SQS",
			rollout:   map[string]bool{"SQS": true},
			envGate:   false,
			expect:    15 * time.Minute, expectOK: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			setGate(t, tt.envGate)
			setInternalGate(t, true, tt.rollout)

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

func TestGetServiceReadTimeout_InternalGateOff(t *testing.T) {
	setGate(t, true)
	setInternalGate(t, false, map[string]bool{"DynamoDB": true})

	got, ok := GetServiceReadTimeout("DynamoDB")
	if !ok {
		t.Fatalf("expect ok true, got false")
	}
	if got != 5*time.Minute {
		t.Errorf("expect %v, got %v", 5*time.Minute, got)
	}
}
