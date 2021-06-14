package manager_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestDownloadStreamOrder(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient(buf12MB)

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})

	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(len(buf12MB)), n; e != a {
		t.Errorf("expect %d buffer length, got %d", e, a)
	}

	if e, a := 3, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}

	expectRngs := []string{"bytes=0-5242879", "bytes=5242880-10485759", "bytes=10485760-15728639"}
	if e, a := expectRngs, *ranges; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v ranges, got %v", e, a)
	}
}

func TestDownloadStreamZero(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient([]byte{})

	d := manager.NewDownloader(c)
	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if n != 0 {
		t.Errorf("expect 0 bytes read, got %d", n)
	}
	// When sliding window is concurrent can expect up to the size of the window
	// but at least the amount that is required
	if e, a := 1, *invocations; a < e {
		t.Errorf("expect at least %v API calls, got %v", e, a)
	}

	sort.Strings(*ranges)
	expectRngs := "bytes=0-5242879"
	if e, a := expectRngs, (*ranges)[0]; a != e {
		t.Errorf("expect %v range, got %v", e, a)
	}
}

func TestDownloadStreamSetPartSize(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient([]byte{1, 2, 3})

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
		d.PartSize = 1
	})
	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(3), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 3, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}
	expectRngs := []string{"bytes=0-0", "bytes=1-1", "bytes=2-2"}
	if e, a := expectRngs, *ranges; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v ranges, got %v", e, a)
	}
	expectBytes := []byte{1, 2, 3}
	if e, a := expectBytes, w.Bytes(); !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v bytes, got %v", e, a)
	}
}

func TestDownloadStreamError(t *testing.T) {
	c, invocations, _ := newDownloadRangeClient([]byte{1, 2, 3})

	num := 0
	orig := c.GetObjectFn
	c.GetObjectFn = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		out, err := orig(ctx, params, optFns...)
		num++
		if num > 1 {
			return &s3.GetObjectOutput{}, fmt.Errorf("s3 service error")
		}
		return out, err
	}

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
		d.PartSize = 1
	})
	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err == nil {
		t.Fatalf("expect error, got none")
	}
	if e, a := "s3 service error", err.Error(); e != a {
		t.Errorf("expect %s error code, got %s", e, a)
	}
	if e, a := int64(1), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 2, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}
	expectBytes := []byte{1}
	if e, a := expectBytes, w.Bytes(); !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v bytes, got %v", e, a)
	}
}

func TestDownloadStreamNonChunk(t *testing.T) {
	c, invocations := newDownloadNonRangeClient(buf2MB)

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})
	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(len(buf2MB)), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 1, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}

	count := 0
	for _, b := range w.Bytes() {
		count += int(b)
	}
	if count != 0 {
		t.Errorf("expect 0 count, got %d", count)
	}
}

func TestDownloadStreamNoContentRangeLength(t *testing.T) {
	s, invocations, _ := newDownloadRangeClient(buf2MB)

	d := manager.NewDownloader(s, func(d *manager.Downloader) {
		d.Concurrency = 1
	})
	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(len(buf2MB)), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 1, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}

	count := 0
	for _, b := range w.Bytes() {
		count += int(b)
	}
	if count != 0 {
		t.Errorf("expect 0 count, got %d", count)
	}
}

func TestDownloadStreamContentRangeTotalAny(t *testing.T) {
	s, invocations := newDownloadContentRangeTotalAnyClient(buf2MB)

	d := manager.NewDownloader(s, func(d *manager.Downloader) {
		d.Concurrency = 1
	})
	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(len(buf2MB)), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 2, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}

	count := 0
	for _, b := range w.Bytes() {
		count += int(b)
	}
	if count != 0 {
		t.Errorf("expect 0 count, got %d", count)
	}
}

func TestDownloadStreamPartBodyRetry_SuccessNoRetry(t *testing.T) {
	c, invocations := newDownloadWithErrReaderClient([]testErrReader{
		{Buf: []byte("abc"), Len: 3, Err: io.EOF},
	})

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})

	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(3), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 1, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}
	if e, a := "abc", string(w.Bytes()); e != a {
		t.Errorf("expect %q response, got %q", e, a)
	}
}

func TestDownloadStreamPartBodyRetry_FailRetry(t *testing.T) {
	c, invocations := newDownloadWithErrReaderClient([]testErrReader{
		{Buf: []byte("ab"), Len: 3, Err: io.ErrUnexpectedEOF},
	})

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
		d.PartBodyMaxRetries = 0
	})

	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err == nil {
		t.Fatalf("expect error, got none")
	}
	if e, a := "unexpected EOF", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q error message to be in %q", e, a)
	}
	if e, a := int64(0), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 1, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}
	if e, a := "ab", string(w.Bytes()); e != a {
		t.Errorf("expect %q response, got %q", e, a)
	}
}

func TestDownloadStreamWithContextCanceled(t *testing.T) {
	d := manager.NewDownloader(s3.New(s3.Options{
		Region: "mock-region",
	}))

	params := s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("Key"),
	}

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	ctx.Error = fmt.Errorf("context canceled")
	close(ctx.DoneCh)

	_, err := d.DownloadStream(ctx, io.Discard, &params)
	if err == nil {
		t.Fatalf("expected error, did not get one")
	}
	if e, a := "canceled", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expected error message to contain %q, but did not %q", e, a)
	}
}

func TestDownloadStream_WithRange(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 10 // should be ignored
		d.PartSize = 1     // should be ignored
	})

	w := &bytes.Buffer{}
	n, err := d.DownloadStream(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
		Range:  aws.String("bytes=2-6"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(5), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 1, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}
	expectRngs := []string{"bytes=2-6"}
	if e, a := expectRngs, *ranges; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v ranges, got %v", e, a)
	}
	expectBytes := []byte{2, 3, 4, 5, 6}
	if e, a := expectBytes, w.Bytes(); !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v bytes, got %v", e, a)
	}
}

func TestDownloadStream_WithFailure(t *testing.T) {
	reqCount := int64(0)
	startingByte := 0

	client := mockDownloadCLient(func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (out *s3.GetObjectOutput, err error) {
		switch atomic.LoadInt64(&reqCount) {
		case 1:
			// Give a chance for the multipart chunks to be queued up
			time.Sleep(1 * time.Second)
			err = fmt.Errorf("some connection error")
		default:
			body := bytes.NewReader(make([]byte, manager.DefaultDownloadPartSize))
			out = &s3.GetObjectOutput{
				Body:          ioutil.NopCloser(body),
				ContentLength: int64(body.Len()),
				ContentRange:  aws.String(fmt.Sprintf("bytes %d-%d/%d", startingByte, body.Len()-1, body.Len()*10)),
			}

			startingByte += body.Len()
			if reqCount > 0 {
				// sleep here to ensure context switching between goroutines
				time.Sleep(25 * time.Millisecond)
			}
		}
		atomic.AddInt64(&reqCount, 1)
		return out, err
	})

	d := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.Concurrency = 2
	})

	params := s3.GetObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
	}

	// Expect this request to exit quickly after failure
	_, err := d.DownloadStream(context.Background(), io.Discard, &params)
	if err == nil {
		t.Fatalf("expect error, got none")
	}

	// with sliding window expect a multiple of the concurrency
	if atomic.LoadInt64(&reqCount) > 4 {
		t.Errorf("expect no more than 4 requests, but received %d", reqCount)
	}
}

func TestDownloaderStreamValidARN(t *testing.T) {
	cases := map[string]struct {
		input   s3.GetObjectInput
		wantErr bool
	}{
		"standard bucket": {
			input: s3.GetObjectInput{
				Bucket: aws.String("test-bucket"),
				Key:    aws.String("test-key"),
			},
		},
		"accesspoint": {
			input: s3.GetObjectInput{
				Bucket: aws.String("arn:aws:s3:us-west-2:123456789012:accesspoint/myap"),
				Key:    aws.String("test-key"),
			},
		},
		"outpost accesspoint": {
			input: s3.GetObjectInput{
				Bucket: aws.String("arn:aws:s3-outposts:us-west-2:012345678901:outpost/op-1234567890123456/accesspoint/myaccesspoint"),
				Key:    aws.String("test-key"),
			},
		},
		"s3-object-lambda accesspoint": {
			input: s3.GetObjectInput{
				Bucket: aws.String("arn:aws:s3-object-lambda:us-west-2:123456789012:accesspoint/myap"),
			},
			wantErr: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			client, _ := newDownloadNonRangeClient(buf2MB)

			downloader := manager.NewDownloader(client, func(downloader *manager.Downloader) {
				downloader.Concurrency = 1
			})

			_, err := downloader.DownloadStream(context.Background(), io.Discard, &tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("err: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}
