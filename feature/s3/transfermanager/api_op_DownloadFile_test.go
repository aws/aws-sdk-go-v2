package transfermanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestDownloadFileRangesAndTruncatesPadding(t *testing.T) {
	data := bytes.Repeat([]byte("abcd"), 2*directIOAlignment+31)
	client := &s3testing.TransferManagerLoggingClient{
		Data:        data,
		GetObjectFn: s3testing.RangeGetObjectFn,
	}
	path := filepath.Join(t.TempDir(), "download.bin")

	out, err := New(client, func(o *Options) {
		o.PartSizeBytes = directIOAlignment
		o.Concurrency = 2
	}).DownloadFile(context.Background(), &DownloadFileInput{
		Bucket:   aws.String("bucket"),
		Key:      aws.String("key"),
		FilePath: path,
	})
	if err != nil {
		t.Fatalf("DownloadFile: %v", err)
	}
	if got, want := aws.ToInt64(out.ContentLength), int64(len(data)); got != want {
		t.Fatalf("ContentLength = %d, want %d", got, want)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("downloaded data does not match source")
	}
	if len(client.RetrievedParts) != 0 || len(client.RetrievedRanges) < 2 {
		t.Fatalf("parts = %v, ranges = %v; want multiple range requests", client.RetrievedParts, client.RetrievedRanges)
	}
}

func TestDownloadFileRetriesShortBody(t *testing.T) {
	data := bytes.Repeat([]byte("r"), 2*directIOAlignment+123)
	client := &s3testing.TransferManagerLoggingClient{Data: data}
	client.GetObjectFn = func(c *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		if c.GetObjectInvocations != 1 {
			return s3testing.RangeGetObjectFn(c, in)
		}

		start, end, err := getReqRange(aws.ToString(in.Range))
		if err != nil {
			return nil, err
		}
		bodyBytes := append([]byte(nil), c.Data[start:start+127]...)
		return &s3.GetObjectOutput{
			Body:          io.NopCloser(bytes.NewReader(bodyBytes)),
			ContentLength: aws.Int64(end - start + 1),
			ContentRange:  aws.String(fmt.Sprintf("bytes %d-%d/%d", start, end, len(c.Data))),
		}, nil
	}
	path := filepath.Join(t.TempDir(), "retry.bin")

	_, err := New(client, func(o *Options) {
		o.PartSizeBytes = directIOAlignment
		o.Concurrency = 1
	}).DownloadFile(context.Background(), &DownloadFileInput{
		Bucket:   aws.String("bucket"),
		Key:      aws.String("key"),
		FilePath: path,
	})
	if err != nil {
		t.Fatalf("DownloadFile: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("downloaded data does not match source after retry")
	}
}

func TestDownloadFileValidatesDirectIOConstraints(t *testing.T) {
	tests := map[string]struct {
		update  func(*Options)
		wantErr string
	}{
		"ranges only": {
			update:  func(o *Options) { o.GetObjectType = types.GetObjectParts },
			wantErr: "GetObjectType must be GetObjectRanges",
		},
		"part alignment": {
			update: func(o *Options) {
				o.PartSizeBytes = directIOAlignment + 1
			},
			wantErr: "PartSizeBytes must be a positive multiple",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := New(&s3testing.TransferManagerLoggingClient{}).DownloadFile(
				context.Background(),
				&DownloadFileInput{FilePath: filepath.Join(t.TempDir(), "download.bin")},
				tt.update,
			)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}
