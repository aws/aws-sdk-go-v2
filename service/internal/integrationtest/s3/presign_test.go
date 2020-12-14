// +build integration

package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-cmp/cmp"
)

func TestInteg_PresignURL(t *testing.T) {
	cases := map[string]struct {
		key                  string
		body                 io.Reader
		expires              time.Duration
		sha256Header         string
		expectedSignedHeader http.Header
	}{
		"standard": {
			body: bytes.NewReader([]byte("Hello-world")),
			expectedSignedHeader: http.Header{
				"content-type":   {"application/octet-stream"},
				"content-length": {"11"},
			},
		},
		"special characters": {
			key: "some_value_(1).foo",
		},
		"nil-body": {
			expectedSignedHeader: http.Header{},
		},
		"empty-body": {
			body:                 bytes.NewReader([]byte("")),
			expectedSignedHeader: http.Header{},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			key := c.key
			if len(key) == 0 {
				key = integrationtest.UniqueID()
			}

			ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFn()

			cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
			if err != nil {
				t.Fatalf("failed to load config, %v", err)
			}

			client := s3.NewFromConfig(cfg)

			// construct a put object
			putObjectInput := &s3.PutObjectInput{
				Bucket: &setupMetadata.Buckets.Source.Name,
				Key:    &key,
				Body:   c.body,
			}

			presignerClient := s3.NewPresignClient(client, func(options *s3.PresignOptions) {
				options.Expires = 600 * time.Second
			})

			presignRequest, err := presignerClient.PresignPutObject(ctx, putObjectInput)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			for k, v := range c.expectedSignedHeader {
				value := presignRequest.SignedHeader[k]
				if len(value) == 0 {
					t.Fatalf("expected %v header to be present in presigned url, got %v", k, presignRequest.SignedHeader)
				}

				if diff := cmp.Diff(v, value); len(diff) != 0 {
					t.Fatalf("expected %v header value to be %v got %v", k, v, value)
				}
			}

			resp, err := sendHTTPRequest(presignRequest, putObjectInput.Body)
			if err != nil {
				t.Errorf("expect no error while sending HTTP request using presigned url, got %v", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("failed to put S3 object, %d:%s", resp.StatusCode, resp.Status)
			}

			// construct a get object
			getObjectInput := &s3.GetObjectInput{
				Bucket: &setupMetadata.Buckets.Source.Name,
				Key:    &key,
			}

			presignRequest, err = presignerClient.PresignGetObject(ctx, getObjectInput)
			if err != nil {
				t.Errorf("expect no error, got %v", err)
			}

			resp, err = sendHTTPRequest(presignRequest, nil)
			if err != nil {
				t.Errorf("expect no error while sending HTTP request using presigned url, got %v", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("failed to get S3 object, %d:%s", resp.StatusCode, resp.Status)
			}
		})
	}
}

func sendHTTPRequest(presignRequest *v4.PresignedHTTPRequest, body io.Reader) (*http.Response, error) {
	// create a http request
	req, err := http.NewRequest(presignRequest.Method, presignRequest.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build presigned request, %v", err)
	}

	// assign the signed headers onto the http request
	for k, vs := range presignRequest.SignedHeader {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	// Need to ensure that the content length member is set of the HTTP Request
	// or the request will NOT be transmitted correctly with a content length
	// value across the wire.
	if contLen := req.Header.Get("Content-Length"); len(contLen) > 0 {
		req.ContentLength, _ = strconv.ParseInt(contLen, 10, 64)
	}

	// assign the request body if not nil
	if body != nil {
		req.Body = ioutil.NopCloser(body)
		if req.ContentLength == 0 {
			req.Body = nil
		}
	}

	// Upload the object to S3.
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}
