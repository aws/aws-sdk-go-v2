package transfermanager

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

func TestUploadMultipartAbortError(t *testing.T) {
	cases := []struct {
		name       string
		cause      error
		abortFails bool
	}{
		{"abort-succeeds", errors.New("original upload failure"), false},
		{"abort-fails", errors.New("original upload failure"), true},
		{"canceled-abort-fails", context.Canceled, true},
		{"deadline-abort-fails", context.DeadlineExceeded, true},
	}
	for _, operation := range []string{"UploadPart", "CompleteMultipartUpload"} {
		for _, tc := range cases {
			t.Run(operation+"/"+tc.name, func(t *testing.T) {
				cause, abortFails := tc.cause, tc.abortFails
				primary := &smithy.OperationError{ServiceID: "S3", OperationName: operation, Err: cause}
				abortError := errors.New("abort request failed")
				c, _, _ := s3testing.NewUploadLoggingClient(nil)
				if operation == "UploadPart" {
					c.UploadPartFn = func(context.Context, *s3testing.TransferManagerLoggingClient, *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
						return nil, primary
					}
				} else {
					c.CompleteMultipartUploadFn = func(context.Context, *s3testing.TransferManagerLoggingClient, *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
						return nil, primary
					}
				}
				aborted := false
				c.AbortMultipartUploadFn = func(_ *s3testing.TransferManagerLoggingClient, input *s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error) {
					aborted = true
					if aws.ToString(input.UploadId) != "UPLOAD-ID" {
						t.Errorf("unexpected aborted upload ID: %v", input.UploadId)
					}
					if abortFails {
						return nil, abortError
					}
					return &s3.AbortMultipartUploadOutput{}, nil
				}
				mgr := New(c, func(o *Options) { o.Concurrency = 1 })
				out, err := mgr.UploadObject(context.Background(), &UploadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"), Body: bytes.NewReader(buf20MB),
				})
				if out != nil || err == nil || !aborted {
					t.Fatalf("expected an aborted failed upload, got output=%v error=%v aborted=%v", out, err, aborted)
				}
				var multipartError MultipartUploadError
				if !errors.As(err, &multipartError) || multipartError.UploadID() != "UPLOAD-ID" {
					t.Fatalf("multipart upload ID is unavailable: %v", err)
				}
				if !errors.Is(err, cause) {
					t.Errorf("original upload cause is not wrapped: %v", err)
				}
				var operationError *smithy.OperationError
				if !errors.As(err, &operationError) || operationError != primary {
					t.Errorf("original operation error is not preserved: %v", err)
				}
				if got := strings.Contains(err.Error(), abortError.Error()); got != abortFails {
					t.Errorf("abort failure reported=%v, want %v: %v", got, abortFails, err)
				}
			})
		}
	}
}
