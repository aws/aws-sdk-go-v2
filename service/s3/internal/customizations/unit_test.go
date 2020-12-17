package customizations_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/awslabs/smithy-go"
)

func Test_EmptyResponse(t *testing.T) {
	cases := map[string]struct {
		status       int
		responseBody []byte
		expectError  bool
	}{
		"success case with no response body": {
			status:       200,
			responseBody: []byte(``),
		},
		"error case with no response body": {
			status:       400,
			responseBody: []byte(``),
			expectError:  true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(c.status)
					w.Write(c.responseBody)
				}))
			defer server.Close()

			ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFn()

			cfg := aws.Config{
				Region: "us-east-1",
				EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:         server.URL,
						SigningName: "s3",
					}, nil
				}),
				Retryer: aws.NopRetryer{},
			}

			client := s3.NewFromConfig(cfg, func(options *s3.Options) {
				options.UsePathStyle = true
			})

			params := &s3.HeadBucketInput{Bucket: aws.String("aws-sdk-go-data")}
			_, err := client.HeadBucket(ctx, params)
			if c.expectError {
				var apiErr smithy.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expect error to be API error, was not, %v", err)
				}
				if len(apiErr.ErrorCode()) == 0 {
					t.Errorf("expect non-empty error code")
				}
				if len(apiErr.ErrorMessage()) == 0 {
					t.Errorf("expect non-empty error message")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err.Error())
				}
			}
		})
	}
}
