package manager_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/manager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestObjectWriterSingePartUpload(t *testing.T) {
	s, ops, args := s3testing.NewUploadLoggingClient(nil)

	mgr := manager.NewUploader(s, func(u *manager.Uploader) {
		u.PartSize = 1024 * 1024 * 7
		u.Concurrency = 1
	})

	wr := mgr.NewWriter(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if _, err := wr.Write(buf2MB); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := wr.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vals := []string{"PutObject"}
	if !reflect.DeepEqual(vals, *ops) {
		t.Errorf("expect %v, got %v", vals, *ops)
	}

	// Part lengths
	if e, a := int64(1024*1024*2), getReaderLength((*args)[0].(*s3.PutObjectInput).Body); e != a {
		t.Errorf("expect %d, got %d", e, a)
	}
}

func TestObjectWriterMultipartUpload(t *testing.T) {
	s, ops, args := s3testing.NewUploadLoggingClient(nil)

	mgr := manager.NewUploader(s, func(u *manager.Uploader) {
		u.PartSize = 1024 * 1024 * 7
		u.Concurrency = 1
	})

	wr := mgr.NewWriter(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if _, err := wr.Write(buf12MB); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := wr.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vals := []string{"CreateMultipartUpload", "UploadPart", "UploadPart", "CompleteMultipartUpload"}
	if !reflect.DeepEqual(vals, *ops) {
		t.Errorf("expect %v, got %v", vals, *ops)
	}

	// Part lengths
	if e, a := int64(1024*1024*7), getReaderLength((*args)[1].(*s3.UploadPartInput).Body); e != a {
		t.Errorf("expect %d, got %d", e, a)
	}

	if e, a := int64(1024*1024*5), getReaderLength((*args)[2].(*s3.UploadPartInput).Body); e != a {
		t.Errorf("expect %d, got %d", e, a)
	}
}
