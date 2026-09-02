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
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
)

// stands in for the ttl the SDK reports, which depends on the real-clock deadline
const ttlPlaceholder = "<ttl>"

// records the Amz-Sdk-Request header of every attempt, failing with a retryable
// status until failures is exhausted. serverSkew offsets the Date header from
// the client's clock, which the SDK turns into the clock skew that the header's
// ttl component is derived from.
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

func newRetryMetricsClient(hc HTTPClient, optFns ...func(*Options)) *Client {
	return New(Options{
		Region: "us-east-1",
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}, nil
		}),
		HTTPClient:         hc,
		EndpointResolverV2: &endpointResolver{},
	}, optFns...)
}

// withDeprecatedMetricsHeader registers the deprecated retry.MetricsHeader
// middleware, which sets the same header from its own slot on the stack. Used to
// compare the two implementations.
func withDeprecatedMetricsHeader(o *Options) {
	o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
		return stack.Finalize.Insert(&retry.MetricsHeader{}, "Retry", middleware.After)
	})
}

func TestRetryMetricsHeader(t *testing.T) {
	restoreSleep := sdk.TestingUseNopSleep()
	defer restoreSleep()
	restoreTime := sdk.TestingUseReferenceTime(time.Unix(1000, 0).UTC())
	defer restoreTime()

	for name, tt := range map[string]struct {
		failures   int
		serverSkew time.Duration
		deadline   bool
		expect     []string
	}{
		"single attempt": {
			expect: []string{"attempt=1; max=3"},
		},
		"header tracks every attempt": {
			failures: 2,
			expect: []string{
				"attempt=1; max=3",
				"attempt=2; max=3",
				"attempt=3; max=3",
			},
		},
		// ttl needs both a deadline to report and a positive clock skew, and the
		// skew only exists once a response has been seen, so it appears from the
		// second attempt on.
		"ttl once skew is known": {
			failures:   1,
			serverSkew: 30 * time.Second,
			deadline:   true,
			expect: []string{
				"attempt=1; max=3",
				"attempt=2; max=3; ttl=" + ttlPlaceholder,
			},
		},
		"no ttl without a deadline": {
			failures:   1,
			serverSkew: 30 * time.Second,
			expect: []string{
				"attempt=1; max=3",
				"attempt=2; max=3",
			},
		},
		"no ttl without skew": {
			failures: 1,
			deadline: true,
			expect: []string{
				"attempt=1; max=3",
				"attempt=2; max=3",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Deadlines run on the real clock, not the frozen test clock, so the
			// ttl the SDK reports is derived here rather than hardcoded. The
			// frozen clock makes the skew exactly serverSkew.
			ctx := context.Background()
			expect := tt.expect
			if tt.deadline {
				deadline := time.Now().UTC().Add(10 * time.Minute)
				var cancel context.CancelFunc
				ctx, cancel = context.WithDeadline(ctx, deadline)
				defer cancel()

				ttl := deadline.Add(tt.serverSkew).Format("20060102T150405Z")
				expect = make([]string, len(tt.expect))
				for i, e := range tt.expect {
					expect[i] = strings.ReplaceAll(e, ttlPlaceholder, ttl)
				}
			}

			// The same scenario has to produce the same header whether the
			// deprecated middleware is on the stack or not.
			for _, deprecated := range []bool{false, true} {
				var optFns []func(*Options)
				if deprecated {
					optFns = append(optFns, withDeprecatedMetricsHeader)
				}

				mock := &retryMetricsHTTP{failures: tt.failures, serverSkew: tt.serverSkew}
				if _, err := newRetryMetricsClient(mock, optFns...).GetItem(ctx, nil); err != nil {
					t.Fatalf("deprecated=%v: %v", deprecated, err)
				}
				if !reflect.DeepEqual(expect, mock.headers) {
					t.Errorf("deprecated=%v: expect %v, got %v", deprecated, expect, mock.headers)
				}
			}
		})
	}
}
