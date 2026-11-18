package transfermanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Exercise bounded ranges through both public APIs with S3-style EOF truncation.
func TestGetAndDownloadObjectBoundedRanges(t *testing.T) {
	data := bytes.Repeat([]byte("0123456789"), 100)
	for _, start := range []int64{0, 1, 500, 999} {
		seen := make(map[int64]bool)
		for _, end := range []int64{start, min(start+100, 999), 999, 1999, math.MaxInt64} {
			if seen[end] {
				continue
			}
			seen[end] = true
			for _, partSize := range []int64{128, 512, 8 << 20} {
				for _, download := range []bool{false, true} {
					t.Run(fmt.Sprintf("start%d-end%d-part%d-download%t", start, end, partSize, download), func(t *testing.T) {
						client, _, _, _, _, _ := s3testing.NewDownloadClient()
						client.Data = data
						client.GetObjectFn = func(_ *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
							var first, last int64
							if _, err := fmt.Sscanf(aws.ToString(in.Range), "bytes=%d-%d", &first, &last); err != nil {
								return nil, err
							}
							if first >= int64(len(data)) {
								return nil, fmt.Errorf("416 range starts past EOF: %d", first)
							}
							last = min(last, int64(len(data))-1)
							body := data[first : last+1]
							return &s3.GetObjectOutput{
								Body: io.NopCloser(bytes.NewReader(body)), ContentLength: aws.Int64(int64(len(body))),
								ContentRange: aws.String(fmt.Sprintf("bytes %d-%d/%d", first, last, len(data))),
								ETag:         aws.String("myetag"),
							}, nil
						}
						manager := New(client, func(o *Options) {
							o.GetObjectType = types.GetObjectRanges
							o.Concurrency = 3
							o.PartSizeBytes = partSize
						})
						rng := aws.String(fmt.Sprintf("bytes=%d-%d", start, end))
						last := min(end, int64(len(data))-1)
						want := data[start : last+1]
						if download {
							buffer := types.NewWriteAtBuffer(nil)
							out, err := manager.DownloadObject(context.Background(), &DownloadObjectInput{Bucket: aws.String("test"), Key: aws.String("test"), Range: rng, WriterAt: buffer})
							if err != nil {
								t.Fatal(err)
							}
							if aws.ToInt64(out.ContentLength) != int64(len(want)) || !bytes.Equal(buffer.Bytes(), want) {
								t.Fatal("wrong downloaded length or bytes")
							}
							if got := aws.ToString(out.ContentRange); got != fmt.Sprintf("bytes=%d-%d", start, last) {
								t.Fatalf("unexpected download ContentRange %q", got)
							}
						} else {
							out, err := manager.GetObject(context.Background(), &GetObjectInput{Bucket: aws.String("test"), Key: aws.String("test"), Range: rng})
							if err != nil {
								t.Fatal(err)
							}
							body, err := io.ReadAll(out.Body)
							if err != nil {
								t.Fatal(err)
							}
							if aws.ToInt64(out.ContentLength) != int64(len(want)) || !bytes.Equal(body, want) {
								t.Fatalf("length=%d bytes=%d want=%d", aws.ToInt64(out.ContentLength), len(body), len(want))
							}
							if got := aws.ToString(out.ContentRange); got != fmt.Sprintf("bytes %d-%d/%d", start, last, len(data)) {
								t.Fatalf("unexpected get ContentRange %q", got)
							}
						}
					})
				}
			}
		}
	}
}

func TestDownloadRangeEndRejectsNonEOFMismatch(t *testing.T) {
	client, _, _, _, _, _ := s3testing.NewDownloadClient()
	client.Data = bytes.Repeat([]byte("x"), 1000)
	client.GetObjectFn = func(c *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		out, err := s3testing.RangeGetObjectFn(c, in)
		if err != nil {
			return nil, err
		}
		if aws.ToString(in.Range) == "bytes=513-999" {
			_ = out.Body.Close()
			out.Body = io.NopCloser(bytes.NewReader(c.Data[513:900]))
			out.ContentLength = aws.Int64(387)
			out.ContentRange = aws.String("bytes 513-899/1000")
		}
		return out, nil
	}
	manager := New(client, func(o *Options) {
		o.GetObjectType = types.GetObjectRanges
		o.PartSizeBytes = 512
		o.Concurrency = 1
	})
	_, err := manager.DownloadObject(context.Background(), &DownloadObjectInput{
		Bucket: aws.String("test"), Key: aws.String("test"), Range: aws.String("bytes=1-1999"), WriterAt: types.NewWriteAtBuffer(nil),
	})
	if err == nil || !strings.Contains(err.Error(), "range mismatch between request 513-999 and response 513-899") {
		t.Fatalf("expected inconsistent later range to fail, got %v", err)
	}
}

func TestDownloadRangeEndAcceptsLaterEOF(t *testing.T) {
	client, _, _, _, _, _ := s3testing.NewDownloadClient()
	client.Data = bytes.Repeat([]byte("x"), 1000)
	client.GetObjectFn = func(c *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		out, err := s3testing.RangeGetObjectFn(c, in)
		if err != nil {
			return nil, err
		}
		if aws.ToString(in.Range) == "bytes=513-999" {
			_ = out.Body.Close()
			out.Body = io.NopCloser(bytes.NewReader(c.Data[513:900]))
			out.ContentLength = aws.Int64(387)
			out.ContentRange = aws.String("bytes 513-899/900")
		}
		return out, nil
	}
	manager := New(client, func(o *Options) {
		o.GetObjectType = types.GetObjectRanges
		o.PartSizeBytes = 512
		o.Concurrency = 1
	})
	buffer := types.NewWriteAtBuffer(nil)
	out, err := manager.DownloadObject(context.Background(), &DownloadObjectInput{
		Bucket: aws.String("test"), Key: aws.String("test"), Range: aws.String("bytes=1-1999"), WriterAt: buffer,
	})
	if err != nil {
		t.Fatalf("expected EOF-truncated later range to succeed, got %v", err)
	}
	want := client.Data[1:900]
	if aws.ToInt64(out.ContentLength) != int64(len(want)) || !bytes.Equal(buffer.Bytes(), want) {
		t.Fatalf("wrong EOF-truncated output: length=%d bytes=%d want=%d", aws.ToInt64(out.ContentLength), len(buffer.Bytes()), len(want))
	}
}

func TestDownloadTotalBytesWithoutContentRange(t *testing.T) {
	for _, tc := range []struct {
		name   string
		known  int64
		length int64
		want   int64
	}{
		{name: "unknown total", known: -1, length: 100, want: 100},
		{name: "shorter response", known: 200, length: 100, want: 100},
		{name: "shorter requested range", known: 50, length: 100, want: 50},
		{name: "equal lengths", known: 100, length: 100, want: 100},
		{name: "empty response", known: 100, length: 0, want: 0},
		{name: "unknown empty response", known: -1, length: 0, want: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := &downloader{totalBytes: tc.known}
			d.setTotalBytes(&s3.GetObjectOutput{ContentLength: aws.Int64(tc.length)})
			if d.err != nil {
				t.Fatal(d.err)
			}
			if d.totalBytes != tc.want {
				t.Fatalf("totalBytes = %d, want %d", d.totalBytes, tc.want)
			}
		})
	}
}
