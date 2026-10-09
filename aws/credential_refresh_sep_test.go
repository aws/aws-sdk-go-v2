package aws

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/credentials/nonrecoverable"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

// These tests run the SEP test cases in testdata/credential-refresh-tests.json against CredentialsCache.

type sepCase struct {
	Documentation string   `json:"documentation"`
	Given         sepGiven `json:"given"`
	Steps         []sepStep
}

type sepGiven struct {
	CachedCredentials               string `json:"cachedCredentials"`
	AccessKeyID                     string `json:"accessKeyId"`
	ConfiguredAdvisoryWindowSeconds *int   `json:"configuredAdvisoryWindowSeconds"`
	RefreshBackoffSeconds           *int   `json:"refreshBackoffSeconds"`
}

type sepStep struct {
	Type                string       `json:"type"`
	Response            string       `json:"response"`
	LifetimeSeconds     *int         `json:"lifetimeSeconds"`
	Seconds             int          `json:"seconds"`
	RejectedAccessKeyID string       `json:"rejectedAccessKeyId"`
	Expected            *sepExpected `json:"expected"`
}

type sepExpected struct {
	Result                string `json:"result"`
	SourceContacted       bool   `json:"sourceContacted"`
	RateLimited           bool   `json:"rateLimited"`
	AdvisoryWindowSeconds *int   `json:"advisoryWindowSeconds"`
}

// sepSeedLifetime is the lifetime of the seeded credentials, and sepSeedAdvisory is the advisory window it gets.
const (
	sepSeedLifetime = 15 * time.Minute
	sepSeedAdvisory = 5 * time.Minute
)

// sepFakeProvider returns the response set for the current step and counts its calls.
type sepFakeProvider struct {
	t               *testing.T
	response        string
	lifetimeSeconds int
	calls           int
}

func (p *sepFakeProvider) Retrieve(ctx context.Context) (Credentials, error) {
	p.calls++
	switch p.response {
	case "freshCredentials":
		lifetime := time.Duration(p.lifetimeSeconds) * time.Second
		return Credentials{
			AccessKeyID: "AKID-NEW", SecretAccessKey: "S", CanExpire: true,
			Expires: sdk.NowTime().Add(lifetime),
		}, nil
	case "staleCredentials":
		return Credentials{
			AccessKeyID: "AKID-STALE", SecretAccessKey: "S", CanExpire: true,
			Expires: sdk.NowTime(),
		}, nil
	case "nonRecoverableError":
		return Credentials{}, &nonrecoverable.Failure{ActionableMessage: "non-recoverable"}
	case "error":
		return Credentials{}, errors.New("transient source error")
	default:
		p.t.Fatalf("Retrieve called with no response configured for this step")
		return Credentials{}, nil
	}
}

// ProviderSources returns IMDS so the cache uses the SEP path, not the legacy one.
func (p *sepFakeProvider) ProviderSources() []CredentialSource {
	return []CredentialSource{CredentialSourceIMDS}
}

func TestCredentialRefreshSEPVectors(t *testing.T) {
	raw, err := os.ReadFile("testdata/credential-refresh-tests.json")
	if err != nil {
		t.Fatalf("read vectors: %v", err)
	}

	var cases []sepCase
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("parse vectors: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("no test vectors loaded")
	}

	for _, tc := range cases {
		t.Run(tc.Documentation, func(t *testing.T) {
			runSEPCase(t, tc)
		})
	}
}

func runSEPCase(t *testing.T, tc sepCase) {
	origNow := sdk.NowTime
	defer func() { sdk.NowTime = origNow }()

	expires := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)

	var now time.Time
	switch tc.Given.CachedCredentials {
	case "none":
		now = expires.Add(-sepSeedLifetime)
	case "valid":
		now = expires.Add(-10 * time.Minute)
	case "advisory":
		now = expires.Add(-3 * time.Minute)
	case "mandatory":
		now = expires.Add(-30 * time.Second)
	case "expired":
		now = expires.Add(time.Second)
	default:
		t.Fatalf("unknown given.cachedCredentials %q", tc.Given.CachedCredentials)
	}
	sdk.NowTime = func() time.Time { return now }

	provider := &sepFakeProvider{t: t}
	var optFns []func(*CredentialsCacheOptions)
	if tc.Given.ConfiguredAdvisoryWindowSeconds != nil {
		seconds := *tc.Given.ConfiguredAdvisoryWindowSeconds
		optFns = append(optFns, func(o *CredentialsCacheOptions) {
			o.ExpiryWindow = time.Duration(seconds) * time.Second
		})
	}
	p := NewCredentialsCache(provider, optFns...)

	if tc.Given.RefreshBackoffSeconds != nil {
		// Pin the backoff to the value the test case expects.
		want := *tc.Given.RefreshBackoffSeconds
		r := float64(want-300) / 300
		p.rand = func() (float64, error) { return r, nil }
	}

	if tc.Given.CachedCredentials != "none" {
		akid := tc.Given.AccessKeyID
		if akid == "" {
			akid = "AKID-0"
		}
		p.creds.Store(&cachedCredentials{
			creds: Credentials{
				AccessKeyID: akid, SecretAccessKey: "S", CanExpire: true, Expires: expires,
			},
			advisoryWindow: sepSeedAdvisory,
		})
	}

	for i, step := range tc.Steps {
		switch step.Type {
		case "advanceTime":
			now = now.Add(time.Duration(step.Seconds) * time.Second)
			continue
		case "invalidate":
			p.InvalidateCredentials(Credentials{AccessKeyID: step.RejectedAccessKeyID})
			continue
		case "getCredentials":
		default:
			t.Fatalf("step %d: unknown step type %q", i, step.Type)
		}

		provider.response = step.Response
		provider.calls = 0
		if step.LifetimeSeconds != nil {
			provider.lifetimeSeconds = *step.LifetimeSeconds
		} else {
			provider.lifetimeSeconds = int(sepSeedLifetime / time.Second)
		}

		// rateLimited means a refresh was due but the backoff skipped it, so check before Retrieve.
		preState := p.creds.Load()
		gotRateLimited := preState != nil && preState.refreshNeeded(now) && preState.rateLimited(now)

		creds, err := p.Retrieve(context.Background())

		exp := step.Expected
		if exp == nil {
			t.Fatalf("step %d: getCredentials step has no expected block", i)
		}

		if got := provider.calls > 0; got != exp.SourceContacted {
			t.Errorf("step %d: sourceContacted = %v, want %v", i, got, exp.SourceContacted)
		}

		switch exp.Result {
		case "cachedCredentials":
			if err != nil {
				t.Errorf("step %d: expected cached credentials, got error: %v", i, err)
			}
		case "newCredentials":
			if err != nil {
				t.Errorf("step %d: expected new credentials, got error: %v", i, err)
			} else if creds.AccessKeyID != "AKID-NEW" {
				t.Errorf("step %d: expected new credentials, got %v", i, creds.AccessKeyID)
			}
		case "noCredentialsError":
			if err == nil {
				t.Errorf("step %d: expected an initial-fetch error, got none", i)
			}
			var nr *nonrecoverable.Failure
			if errors.As(err, &nr) {
				t.Errorf("step %d: expected a generic no-credentials error, got non-recoverable: %v", i, err)
			}
		case "nonRecoverableError":
			var nr *nonrecoverable.Failure
			if !errors.As(err, &nr) {
				t.Errorf("step %d: expected a non-recoverable error, got %v", i, err)
			}
		default:
			t.Fatalf("step %d: unknown expected.result %q", i, exp.Result)
		}

		if exp.AdvisoryWindowSeconds != nil {
			e := p.creds.Load()
			want := time.Duration(*exp.AdvisoryWindowSeconds) * time.Second
			if e == nil {
				t.Errorf("step %d: no cached credentials, want advisory window %v", i, want)
			} else if e.advisoryWindow != want {
				t.Errorf("step %d: advisory window = %v, want %v", i, e.advisoryWindow, want)
			}
		}

		if gotRateLimited != exp.RateLimited {
			t.Errorf("step %d: rateLimited = %v, want %v", i, gotRateLimited, exp.RateLimited)
		}
	}
}
