package manager_test

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestObjectReaderSinglePart(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient(buf2MB)

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})

	rd := d.NewReader(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	n, err := io.Copy(io.Discard, rd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if e, a := int64(len(buf2MB)), n; e != a {
		t.Errorf("expected %d buffer length, got %d", e, a)
	}

	if e, a := 2, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}

	expectRngs := []string{"bytes=0-0", "bytes=0-2097152"}
	if e, a := expectRngs, *ranges; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v ranges, got %v", e, a)
	}
}

func TestObjectReaderMultiPart(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient(buf12MB)

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})

	rd := d.NewReader(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	n, err := io.Copy(io.Discard, rd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if e, a := int64(len(buf12MB)), n; e != a {
		t.Errorf("expected %d buffer length, got %d", e, a)
	}

	if e, a := 4, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}

	expectRngs := []string{"bytes=0-0", "bytes=0-5242880", "bytes=5242881-10485761", "bytes=10485762-12582912"}
	if e, a := expectRngs, *ranges; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v ranges, got %v", e, a)
	}
}
