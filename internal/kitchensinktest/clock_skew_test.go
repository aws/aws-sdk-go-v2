package kitchensinktest

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

// simulates a skew-aware server, set ServerTime and Do will reject if the
// request time is outside the threshold
type skewHTTP struct {
	ServerTime time.Time
	attempts   int
}

func (m *skewHTTP) Do(req *http.Request) (*http.Response, error) {
	m.attempts++

	reqTime, err := time.Parse("20060102T150405Z", req.Header.Get("X-Amz-Date"))
	if err != nil {
		return nil, err
	}

	if m.ServerTime.Sub(reqTime).Abs() > 4*time.Minute {
		return m.skewed(), nil
	}
	return m.ok(), nil
}

func (m *skewHTTP) skewed() *http.Response {
	return &http.Response{
		StatusCode: 400,
		Header: http.Header{
			"Date":         []string{m.ServerTime.UTC().Format(http.TimeFormat)},
			"Content-Type": []string{"application/x-amz-json-1.0"},
		},
		Body: io.NopCloser(strings.NewReader(mkerr("RequestTimeTooSkewed"))),
	}
}

func (m *skewHTTP) ok() *http.Response {
	return &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Date":         []string{m.ServerTime.UTC().Format(http.TimeFormat)},
			"Content-Type": []string{"application/x-amz-json-1.0"},
		},
		Body: io.NopCloser(strings.NewReader("{}")),
	}
}

func (m *skewHTTP) reset() { m.attempts = 0 }

func mkerr(typ string) string {
	return `{"__type":"` + typ + `","message":"error message"}`
}

func newSkewClient(hc HTTPClient, optfn ...func(*Options)) *Client {
	opts := Options{
		Region: "us-east-1",
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}, nil
		}),
		HTTPClient:         hc,
		EndpointResolverV2: &endpointResolver{},
	}
	return New(opts, optfn...)
}

func TestClockSkew_OK(t *testing.T) {
	tm := time.Unix(1000, 0).UTC()
	restore := sdk.TestingUseReferenceTime(tm)
	defer restore()

	mock := &skewHTTP{ServerTime: tm}
	svc := newSkewClient(mock)

	_, err := svc.GetItem(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if mock.attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", mock.attempts)
	}
}

func TestClockSkew_DefiniteCode(t *testing.T) {
	clientTime := time.Unix(1000, 0).UTC()
	restore := sdk.TestingUseReferenceTime(clientTime)
	defer restore()

	mock := &skewHTTP{ServerTime: clientTime.Add(-5 * time.Minute)}
	svc := newSkewClient(mock)

	// Op 1: skew triggers retry, corrects on attempt 2
	if _, err := svc.GetItem(context.Background(), nil); err != nil {
		t.Fatalf("op 1: %v", err)
	}
	if mock.attempts != 2 {
		t.Errorf("op 1: expected 2 attempts, got %d", mock.attempts)
	}

	// Op 2: stored offset lets first attempt succeed
	mock.reset()
	if _, err := svc.GetItem(context.Background(), nil); err != nil {
		t.Fatalf("op 2: %v", err)
	}
	if mock.attempts != 1 {
		t.Errorf("op 2: expected 1 attempt, got %d", mock.attempts)
	}

	// Op 3: server clock fixed, stale offset causes a retry then heals
	mock.ServerTime = clientTime
	mock.reset()
	if _, err := svc.GetItem(context.Background(), nil); err != nil {
		t.Fatalf("op 3: %v", err)
	}
	if mock.attempts != 2 {
		t.Errorf("op 3: expected 2 attempts, got %d", mock.attempts)
	}

	// Op 4: healed offset sticks
	mock.reset()
	if _, err := svc.GetItem(context.Background(), nil); err != nil {
		t.Fatalf("op 4: %v", err)
	}
	if mock.attempts != 1 {
		t.Errorf("op 4: expected 1 attempt, got %d", mock.attempts)
	}
}

func TestClockSkew_HealWithMaxAttempts1(t *testing.T) {
	clientTime := time.Unix(1000, 0).UTC()
	restore := sdk.TestingUseReferenceTime(clientTime)
	defer restore()

	mock := &skewHTTP{ServerTime: clientTime.Add(-5 * time.Minute)}
	svc := newSkewClient(mock, func(o *Options) {
		o.RetryMaxAttempts = 1
	})

	// Op 1: skew detected, no retry (max attempts 1), but offset is saved
	_, err := svc.GetItem(context.Background(), nil)
	if err == nil {
		t.Fatal("op 1: expected error")
	}
	if mock.attempts != 1 {
		t.Errorf("op 1: expected 1 attempt, got %d", mock.attempts)
	}

	// Op 2: saved offset applied, succeeds
	mock.reset()
	if _, err := svc.GetItem(context.Background(), nil); err != nil {
		t.Fatalf("op 2: %v", err)
	}
	if mock.attempts != 1 {
		t.Errorf("op 2: expected 1 attempt, got %d", mock.attempts)
	}

	// Op 3: server clock fixed, stale offset causes skew, fails but heals
	mock.ServerTime = clientTime
	mock.reset()
	_, err = svc.GetItem(context.Background(), nil)
	if err == nil {
		t.Fatal("op 3: expected error")
	}
	if mock.attempts != 1 {
		t.Errorf("op 3: expected 1 attempt, got %d", mock.attempts)
	}

	// Op 4: healed offset sticks
	mock.reset()
	if _, err := svc.GetItem(context.Background(), nil); err != nil {
		t.Fatalf("op 4: %v", err)
	}
	if mock.attempts != 1 {
		t.Errorf("op 4: expected 1 attempt, got %d", mock.attempts)
	}
}
