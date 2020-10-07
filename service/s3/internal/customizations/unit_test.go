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

func Test_HeadBucket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
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
		Retryer: aws.NoOpRetryer{},
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})
	params := &s3.HeadBucketInput{Bucket: aws.String("aws-sdk-go-data")}
	_, err := client.HeadBucket(ctx, params)
	if err == nil {
		t.Error("expected error, got none")
	}

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
}
