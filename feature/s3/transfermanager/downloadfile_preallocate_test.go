package transfermanager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type recordingPreallocator struct {
	mu               sync.Mutex
	preallocatedSize int64
	preallocateCalls int
	preallocateErr   error
	writes           int
	data             []byte
}

func (w *recordingPreallocator) preallocate(size int64) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.preallocateCalls++
	w.preallocatedSize = size
	return w.preallocateErr
}

func (w *recordingPreallocator) WriteAt(p []byte, off int64) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.preallocateCalls == 0 {
		return 0, errors.New("write occurred before preallocation")
	}
	end := int(off) + len(p)
	if end > len(w.data) {
		w.data = append(w.data, make([]byte, end-len(w.data))...)
	}
	w.writes++
	return copy(w.data[off:], p), nil
}

func (w *recordingPreallocator) snapshot() (size int64, calls, writes int, data []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.preallocatedSize, w.preallocateCalls, w.writes, bytes.Clone(w.data)
}

func TestDownloaderPreallocatesBeforeFirstWrite(t *testing.T) {
	data := bytes.Repeat([]byte("p"), 2*directIOAlignment+17)
	client := &s3testing.TransferManagerLoggingClient{
		Data:        data,
		GetObjectFn: s3testing.RangeGetObjectFn,
	}
	writer := &recordingPreallocator{}

	out, err := New(client, func(o *Options) {
		o.GetObjectType = types.GetObjectRanges
		o.PartSizeBytes = directIOAlignment
		o.Concurrency = 2
	}).DownloadObject(context.Background(), &DownloadObjectInput{
		Bucket:   aws.String("bucket"),
		Key:      aws.String("key"),
		WriterAt: writer,
	})
	if err != nil {
		t.Fatalf("DownloadObject: %v", err)
	}

	size, calls, writes, got := writer.snapshot()
	if size != int64(len(data)) {
		t.Fatalf("preallocated size = %d, want %d", size, len(data))
	}
	if calls != 1 {
		t.Fatalf("preallocate calls = %d, want 1", calls)
	}
	if writes == 0 {
		t.Fatal("no writes occurred")
	}
	if !bytes.Equal(got, data) {
		t.Fatal("written data does not match source")
	}
	if aws.ToInt64(out.ContentLength) != int64(len(data)) {
		t.Fatalf("ContentLength = %d, want %d", aws.ToInt64(out.ContentLength), len(data))
	}
}

func TestDownloaderPreallocatesRequestedRangeSize(t *testing.T) {
	data := bytes.Repeat([]byte("r"), 3*directIOAlignment)
	client := &s3testing.TransferManagerLoggingClient{
		Data:        data,
		GetObjectFn: s3testing.RangeGetObjectFn,
	}
	writer := &recordingPreallocator{}

	_, err := New(client, func(o *Options) {
		o.GetObjectType = types.GetObjectRanges
		o.PartSizeBytes = directIOAlignment
		o.Concurrency = 1
	}).DownloadObject(context.Background(), &DownloadObjectInput{
		Bucket:   aws.String("bucket"),
		Key:      aws.String("key"),
		Range:    aws.String(fmt.Sprintf("bytes=%d-%d", directIOAlignment, 3*directIOAlignment-1)),
		WriterAt: writer,
	})
	if err != nil {
		t.Fatalf("DownloadObject: %v", err)
	}

	size, calls, _, _ := writer.snapshot()
	if size != 2*directIOAlignment {
		t.Fatalf("preallocated size = %d, want %d", size, 2*directIOAlignment)
	}
	if calls != 1 {
		t.Fatalf("preallocate calls = %d, want 1", calls)
	}
}

func TestDownloaderSkipsPreallocationForZeroLength(t *testing.T) {
	client := &s3testing.TransferManagerLoggingClient{
		GetObjectFn: func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body:          io.NopCloser(strings.NewReader("")),
				ContentLength: aws.Int64(0),
			}, nil
		},
	}
	writer := &recordingPreallocator{}

	_, err := New(client, func(o *Options) {
		o.GetObjectType = types.GetObjectRanges
		o.PartSizeBytes = directIOAlignment
	}).DownloadObject(context.Background(), &DownloadObjectInput{
		Bucket:   aws.String("bucket"),
		Key:      aws.String("key"),
		WriterAt: writer,
	})
	if err != nil {
		t.Fatalf("DownloadObject: %v", err)
	}

	_, calls, writes, _ := writer.snapshot()
	if calls != 0 {
		t.Fatalf("preallocate calls = %d, want 0", calls)
	}
	if writes != 0 {
		t.Fatalf("writes = %d, want 0", writes)
	}
}

func TestDownloaderPreallocationErrorStopsWrites(t *testing.T) {
	preallocateErr := errors.New("fallocate failed")
	client := &s3testing.TransferManagerLoggingClient{
		GetObjectFn: func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body:          io.NopCloser(strings.NewReader("data")),
				ContentLength: aws.Int64(4),
			}, nil
		},
	}
	writer := &recordingPreallocator{preallocateErr: preallocateErr}

	_, err := New(client, func(o *Options) {
		o.GetObjectType = types.GetObjectRanges
		o.PartSizeBytes = directIOAlignment
	}).DownloadObject(context.Background(), &DownloadObjectInput{
		Bucket:   aws.String("bucket"),
		Key:      aws.String("key"),
		WriterAt: writer,
	})
	if err == nil || !strings.Contains(err.Error(), preallocateErr.Error()) {
		t.Fatalf("DownloadObject error = %v, want %v", err, preallocateErr)
	}

	_, calls, writes, _ := writer.snapshot()
	if calls != 1 {
		t.Fatalf("preallocate calls = %d, want 1", calls)
	}
	if writes != 0 {
		t.Fatalf("writes after preallocation error = %d, want 0", writes)
	}
}
