package kitchensinktest

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

// serverSkew offsets the Date header from the client's clock, which the SDK turns
// into the clock skew that the header's ttl component is derived from.
type retryMetricsHTTP struct {
	failures   int
	serverSkew time.Duration

	headers []string
}

func (m *retryMetricsHTTP) Do(req *http.Request) (*http.Response, error) {
	m.headers = append(m.headers, req.Header.Get("Amz-Sdk-Request"))

	status, body := 200, "{}"
	if len(m.headers) <= m.failures {
		status, body = 500, mkerr("InternalServerError")
	}
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Date":         []string{sdk.NowTime().Add(m.serverSkew).UTC().Format(http.TimeFormat)},
			"Content-Type": []string{"application/x-amz-json-1.0"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}, nil
}

func TestRetryMetricsHeader(t *testing.T) {
	restoreSleep := sdk.TestingUseNopSleep()
	defer restoreSleep()
	restoreTime := sdk.TestingUseReferenceTime(time.Unix(1000, 0).UTC())
	defer restoreTime()

	const serverSkew = 30 * time.Second

	// Deadlines run on the real clock, not the frozen test clock, so the ttl is
	// derived here rather than hardcoded.
	deadline := time.Now().UTC().Add(10 * time.Minute)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	ttl := "; ttl=" + deadline.Add(serverSkew).Format("20060102T150405Z")

	mock := &retryMetricsHTTP{failures: 2, serverSkew: serverSkew}
	svc := New(Options{
		Region: "us-east-1",
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}, nil
		}),
		HTTPClient:         mock,
		EndpointResolverV2: &endpointResolver{},
	})

	if _, err := svc.GetItem(ctx, nil); err != nil {
		t.Fatal(err)
	}

	// The skew only exists once a response has been seen, so ttl appears from the
	// second attempt on.
	expect := []string{
		"attempt=1; max=3",
		"attempt=2; max=3" + ttl,
		"attempt=3; max=3" + ttl,
	}
	if !reflect.DeepEqual(expect, mock.headers) {
		t.Errorf("expect %v, got %v", expect, mock.headers)
	}
}
