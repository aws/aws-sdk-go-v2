package manager_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type invalidRangeClient struct {
}

func TestDownload_RangeMismatch(t *testing.T) {
	c, _, _ := newDownloadBadRangeClient(buf12MB)

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})

	w := manager.NewWriteAtBuffer(make([]byte, len(buf12MB)))
	_, err := d.Download(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	if err == nil {
		t.Fatalf("expect err, got none")
	}
	if !strings.Contains(err.Error(), "invalid content range") {
		t.Errorf("error mismatch:\n%v !=\n%v", err, "invalid content range")
	}
}

func TestDownload_RangeMismatchDisabled(t *testing.T) {
	c, _, _ := newDownloadBadRangeClient(buf12MB)

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
		d.DisableValidateParts = true
	})

	w := manager.NewWriteAtBuffer(make([]byte, len(buf12MB)))
	_, err := d.Download(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Fatalf("expect no err, got %v", err)
	}
}
