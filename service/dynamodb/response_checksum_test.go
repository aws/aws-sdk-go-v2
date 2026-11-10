package dynamodb

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestResponseChecksumValidation(t *testing.T) {
	const body = `{"TableNames":["t1"]}`
	good := fmt.Sprint(crc32.ChecksumIEEE([]byte(body)))

	cases := map[string]struct {
		header      []string
		disable     bool
		expectErr   string
		expectCalls int32
	}{
		"matching":          {header: []string{good}, expectCalls: 1},
		"absent":            {expectCalls: 1},
		"mismatch":          {header: []string{"12345"}, expectErr: "did not match expected checksum, 12345", expectCalls: 3},
		"mismatch disabled": {header: []string{"12345"}, disable: true, expectCalls: 1},
		"first mismatch":    {header: []string{"12345", good}, expectErr: "did not match expected checksum", expectCalls: 3},
		"unparsable":        {header: []string{"abc"}, expectErr: "unable to parse checksum header", expectCalls: 1},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var calls atomic.Int32
			httpClient := smithyhttp.ClientDoFunc(func(*http.Request) (*http.Response, error) {
				calls.Add(1)
				header := http.Header{"Content-Type": {"application/x-amz-json-1.0"}}
				for _, v := range c.header {
					header.Add("X-Amz-Crc32", v)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     header,
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			})

			client := New(Options{
				Region:       "us-west-2",
				BaseEndpoint: aws.String("https://dynamodb.test"),
				HTTPClient:   httpClient,
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}, nil
				}),
				Retryer: retry.NewStandard(func(o *retry.StandardOptions) {
					o.MaxAttempts = 3
					o.Backoff = retry.BackoffDelayerFunc(func(int, error) (time.Duration, error) { return 0, nil })
				}),
				DisableValidateResponseChecksum: c.disable,
			})

			out, err := client.ListTables(context.Background(), &ListTablesInput{})
			if c.expectErr == "" {
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
				if e, a := 1, len(out.TableNames); e != a {
					t.Errorf("expect %d table names, got %d", e, a)
				}
			} else {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if !strings.Contains(err.Error(), c.expectErr) {
					t.Errorf("expect error to contain %q, got %v", c.expectErr, err)
				}
			}

			if e, a := c.expectCalls, calls.Load(); e != a {
				t.Errorf("expect %d calls, got %d", e, a)
			}
		})
	}
}
