package transfermanager

import (
	"bytes"
	"context"
	"errors"
	"math"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestNewFromConfig(t *testing.T) {
	baseHTTPClient := awshttp.NewBuildableClient().WithTransportOptions(func(transport *http.Transport) {
		transport.MaxIdleConnsPerHost = 7
	})
	cfg := &aws.Config{
		Region:     "us-west-2",
		HTTPClient: baseHTTPClient,
	}

	client := NewFromConfig(cfg)
	if got, want := client.options.Concurrency, defaultConfigTransferConcurrency; got != want {
		t.Fatalf("expected concurrency %d, got %d", want, got)
	}

	s3Client, ok := client.options.S3.(*s3.Client)
	if !ok {
		t.Fatalf("expected an S3 client, got %T", client.options.S3)
	}
	configuredHTTPClient, ok := s3Client.Options().HTTPClient.(*awshttp.BuildableClient)
	if !ok {
		t.Fatalf("expected a buildable HTTP client, got %T", s3Client.Options().HTTPClient)
	}
	if got, want := configuredHTTPClient.GetTransport().MaxIdleConnsPerHost, math.MaxInt; got != want {
		t.Errorf("expected max idle connections per host %d, got %d", want, got)
	}
	if got := configuredHTTPClient.GetTransport().MaxIdleConns; got != 0 {
		t.Errorf("expected unlimited max idle connections, got %d", got)
	}

	if cfg.HTTPClient != baseHTTPClient {
		t.Error("expected input config HTTP client to remain unchanged")
	}
	if got, want := baseHTTPClient.GetTransport().MaxIdleConnsPerHost, 7; got != want {
		t.Errorf("expected input transport max idle connections per host %d, got %d", want, got)
	}
}

func TestNewFromConfigOptionsOverrideDefaults(t *testing.T) {
	client := NewFromConfig(&aws.Config{}, func(o *Options) {
		o.Concurrency = 9
	})

	if got, want := client.options.Concurrency, 9; got != want {
		t.Fatalf("expected concurrency %d, got %d", want, got)
	}
}

func TestNewClientOptions(t *testing.T) {
	s3Client, _, _ := s3testing.NewUploadLoggingClient(nil)
	var clientOptionCalls atomic.Int32
	client := New(s3Client, func(o *Options) {
		o.ClientOptions = append(o.ClientOptions, func(*s3.Options) {
			clientOptionCalls.Add(1)
		})
	})

	var operationOptionCalls atomic.Int32
	_, err := client.UploadObject(context.Background(), &UploadObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
		Body:   bytes.NewReader(nil),
	}, func(o *Options) {
		o.ClientOptions = append(o.ClientOptions, func(*s3.Options) {
			operationOptionCalls.Add(1)
		})
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got, want := clientOptionCalls.Load(), int32(1); got != want {
		t.Errorf("expected client option to run %d time, got %d", want, got)
	}
	if got, want := operationOptionCalls.Load(), int32(1); got != want {
		t.Errorf("expected operation option to run %d time, got %d", want, got)
	}
	if got, want := len(client.options.ClientOptions), 1; got != want {
		t.Errorf("expected original client to retain %d client option, got %d", want, got)
	}
}

type testHTTPClient struct{}

func (*testHTTPClient) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("unexpected request")
}

func TestNewFromConfigPreservesCustomHTTPClient(t *testing.T) {
	httpClient := &testHTTPClient{}
	client := NewFromConfig(&aws.Config{HTTPClient: httpClient})
	s3Client, ok := client.options.S3.(*s3.Client)
	if !ok {
		t.Fatalf("expected an S3 client, got %T", client.options.S3)
	}

	if got := s3Client.Options().HTTPClient; got != httpClient {
		t.Fatalf("expected configured HTTP client to be preserved, got %T", got)
	}
}

func TestClientOptionsAppliedToMultipartAbort(t *testing.T) {
	s3Client, invocations, _ := s3testing.NewUploadLoggingClient(nil)
	s3Client.UploadPartFn = func(context.Context, *s3testing.TransferManagerLoggingClient, *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
		return nil, errors.New("upload part failed")
	}

	var clientOptionCalls atomic.Int32
	client := New(s3Client, func(o *Options) {
		o.Concurrency = 1
		o.PartSizeBytes = 5
		o.MultipartUploadThreshold = 5
		o.ClientOptions = append(o.ClientOptions, func(*s3.Options) {
			clientOptionCalls.Add(1)
		})
	})

	_, err := client.UploadObject(context.Background(), &UploadObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
		Body:   bytes.NewReader(make([]byte, 6)),
	})
	if err == nil {
		t.Fatal("expected multipart upload to fail")
	}

	if got, want := clientOptionCalls.Load(), int32(len(*invocations)); got != want {
		t.Errorf("expected client option on all %d S3 calls, got %d", want, got)
	}
	if got := (*invocations)[len(*invocations)-1]; got != "AbortMultipartUpload" {
		t.Errorf("expected final operation to abort multipart upload, got %s", got)
	}
}
