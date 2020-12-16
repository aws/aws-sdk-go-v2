package customizations_test

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestPutObject_PresignURL(t *testing.T) {
	cases := map[string]struct {
		input                  s3.PutObjectInput
		options                s3.PresignOptions
		expectPresignedURLHost string
		expectRequestURIQuery  []string
		expectSignedHeader     http.Header
		expectMethod           string
		expectError            string
	}{
		"standard case": {
			input: s3.PutObjectInput{
				Bucket: aws.String("mock-bucket"),
				Key:    aws.String("mockkey"),
				Body:   strings.NewReader("hello-world"),
			},
			expectPresignedURLHost: "https://mock-bucket.s3.us-west-2.amazonaws.com/mockkey?",
			expectRequestURIQuery: []string{
				"X-Amz-Expires=900",
				"X-Amz-Credential",
				"X-Amz-Date",
				"x-id=PutObject",
				"X-Amz-Signature",
			},
			expectMethod: "PUT",
			expectSignedHeader: http.Header{
				"content-length": []string{"11"},
				"content-type":   []string{"application/octet-stream"},
				"host":           []string{"mock-bucket.s3.us-west-2.amazonaws.com"},
			},
		},
		"seekable payload": {
			input: s3.PutObjectInput{
				Bucket: aws.String("mock-bucket"),
				Key:    aws.String("mockkey"),
				Body:   bytes.NewReader([]byte("hello-world")),
			},
			expectPresignedURLHost: "https://mock-bucket.s3.us-west-2.amazonaws.com/mockkey?",
			expectRequestURIQuery: []string{
				"X-Amz-Expires=900",
				"X-Amz-Credential",
				"X-Amz-Date",
				"x-id=PutObject",
				"X-Amz-Signature",
			},
			expectMethod: "PUT",
			expectSignedHeader: http.Header{
				"content-length": []string{"11"},
				"content-type":   []string{"application/octet-stream"},
				"host":           []string{"mock-bucket.s3.us-west-2.amazonaws.com"},
			},
		},
		"unseekable payload": {
			// unseekable payload succeeds as we disable content sha256 computation for streaming input
			input: s3.PutObjectInput{
				Bucket: aws.String("mock-bucket"),
				Key:    aws.String("mockkey"),
				Body:   bytes.NewBuffer([]byte(`hello-world`)),
			},
			expectPresignedURLHost: "https://mock-bucket.s3.us-west-2.amazonaws.com/mockkey?",
			expectRequestURIQuery: []string{
				"X-Amz-Expires=900",
				"X-Amz-Credential",
				"X-Amz-Date",
				"x-id=PutObject",
				"X-Amz-Signature",
			},
			expectMethod: "PUT",
			expectSignedHeader: http.Header{
				"content-length": []string{"11"},
				"content-type":   []string{"application/octet-stream"},
				"host":           []string{"mock-bucket.s3.us-west-2.amazonaws.com"},
			},
		},
		"empty body": {
			input: s3.PutObjectInput{
				Bucket: aws.String("mock-bucket"),
				Key:    aws.String("mockkey"),
				Body:   bytes.NewReader([]byte(``)),
			},
			expectPresignedURLHost: "https://mock-bucket.s3.us-west-2.amazonaws.com/mockkey?",
			expectRequestURIQuery: []string{
				"X-Amz-Expires=900",
				"X-Amz-Credential",
				"X-Amz-Date",
				"x-id=PutObject",
				"X-Amz-Signature",
			},
			expectMethod: "PUT",
			expectSignedHeader: http.Header{
				"host": []string{"mock-bucket.s3.us-west-2.amazonaws.com"},
			},
		},
		"nil body": {
			input: s3.PutObjectInput{
				Bucket: aws.String("mock-bucket"),
				Key:    aws.String("mockkey"),
			},
			expectPresignedURLHost: "https://mock-bucket.s3.us-west-2.amazonaws.com/mockkey?",
			expectRequestURIQuery: []string{
				"X-Amz-Expires=900",
				"X-Amz-Credential",
				"X-Amz-Date",
				"x-id=PutObject",
				"X-Amz-Signature",
			},
			expectMethod: "PUT",
			expectSignedHeader: http.Header{
				"host": []string{"mock-bucket.s3.us-west-2.amazonaws.com"},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			cfg := aws.Config{
				Region:      "us-west-2",
				Credentials: unit.StubCredentialsProvider{},
				Retryer:     aws.NopRetryer{},
			}
			presignClient := s3.NewPresignClient(s3.NewFromConfig(cfg), func(options *s3.PresignOptions) {
				options = &c.options
			})

			req, err := presignClient.PresignPutObject(ctx, &c.input)
			if err != nil {
				if len(c.expectError) == 0 {
					t.Fatalf("expected no error, got %v", err)
				}
				// if expect error, match error and skip rest
				if e, a := c.expectError, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expected error to be %s, got %s", e, a)
				}
			} else {
				if len(c.expectError) != 0 {
					t.Fatalf("expected error to be %v, got none", c.expectError)
				}
			}

			if e, a := c.expectPresignedURLHost, req.URL; !strings.Contains(a, e) {
				t.Fatalf("expected presigned url to contain host %s, got %s", e, a)
			}

			if len(c.expectRequestURIQuery) != 0 {
				for _, label := range c.expectRequestURIQuery {
					if e, a := label, req.URL; !strings.Contains(a, e) {
						t.Fatalf("expected presigned url to contain %v label in url: %v", label, req.URL)
					}
				}
			}

			if e, a := c.expectSignedHeader, req.SignedHeader; len(cmp.Diff(e, a)) != 0 {
				t.Fatalf("expected signed header to be %s, got %s, \n diff : %s", e, a, cmp.Diff(e, a))
			}

			if e, a := c.expectMethod, req.Method; !strings.EqualFold(e, a) {
				t.Fatalf("expected presigning Method to be %s, got %s", e, a)
			}

		})
	}
}
