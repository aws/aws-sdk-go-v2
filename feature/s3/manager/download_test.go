package manager_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	managertesting "github.com/aws/aws-sdk-go-v2/feature/s3/manager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type downloadCaptureClient struct {
	GetObjectFn          func(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	GetObjectInvocations int

	RetrievedRanges []string

	lock sync.Mutex
}

func (c *downloadCaptureClient) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.GetObjectInvocations++

	if params.Range != nil {
		c.RetrievedRanges = append(c.RetrievedRanges, aws.ToString(params.Range))
	}

	return c.GetObjectFn(ctx, params, optFns...)
}

var rangeValueRegex = regexp.MustCompile(`bytes=(\d+)-(\d+)`)

func parseRange(rangeValue string) (start, fin int64) {
	rng := rangeValueRegex.FindStringSubmatch(rangeValue)
	start, _ = strconv.ParseInt(rng[1], 10, 64)
	fin, _ = strconv.ParseInt(rng[2], 10, 64)
	return start, fin
}

func newDownloadRangeClient(data []byte) (*downloadCaptureClient, *int, *[]string) {
	capture := &downloadCaptureClient{}

	capture.GetObjectFn = func(_ context.Context, params *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		start, fin := parseRange(aws.ToString(params.Range))
		fin++

		if fin >= int64(len(data)) {
			fin = int64(len(data))
		}

		bodyBytes := data[start:fin]

		return &s3.GetObjectOutput{
			Body:          ioutil.NopCloser(bytes.NewReader(bodyBytes)),
			ContentRange:  aws.String(fmt.Sprintf("bytes %d-%d/%d", start, fin-1, len(data))),
			ContentLength: int64(len(bodyBytes)),
		}, nil
	}

	return capture, &capture.GetObjectInvocations, &capture.RetrievedRanges
}

func newDownloadNonRangeClient(data []byte) (*downloadCaptureClient, *int) {
	capture := &downloadCaptureClient{}

	capture.GetObjectFn = func(_ context.Context, params *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		return &s3.GetObjectOutput{
			Body:          ioutil.NopCloser(bytes.NewReader(data[:])),
			ContentLength: int64(len(data)),
		}, nil
	}

	return capture, &capture.GetObjectInvocations
}

type mockHTTPStatusError struct {
	StatusCode int
}

func (m *mockHTTPStatusError) Error() string {
	return fmt.Sprintf("http status code: %v", m.StatusCode)
}

func (m *mockHTTPStatusError) HTTPStatusCode() int {
	return m.StatusCode
}

func newDownloadContentRangeTotalAnyClient(data []byte) (*downloadCaptureClient, *int) {
	capture := &downloadCaptureClient{}
	completed := false

	capture.GetObjectFn = func(_ context.Context, params *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		if completed {
			return nil, &mockHTTPStatusError{StatusCode: 416}
		}

		start, fin := parseRange(aws.ToString(params.Range))
		fin++

		if fin >= int64(len(data)) {
			fin = int64(len(data))
			completed = true
		}

		bodyBytes := data[start:fin]

		return &s3.GetObjectOutput{
			Body:         ioutil.NopCloser(bytes.NewReader(bodyBytes)),
			ContentRange: aws.String(fmt.Sprintf("bytes %d-%d/*", start, fin-1)),
		}, nil
	}

	return capture, &capture.GetObjectInvocations
}

func newDownloadWithErrReaderClient(cases []testErrReader) (*downloadCaptureClient, *int) {
	var index int

	c := &downloadCaptureClient{}
	c.GetObjectFn = func(_ context.Context, params *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		c := cases[index]
		out := &s3.GetObjectOutput{
			Body:          ioutil.NopCloser(&c),
			ContentRange:  aws.String(fmt.Sprintf("bytes %d-%d/%d", 0, c.Len-1, c.Len)),
			ContentLength: c.Len,
		}
		index++
		return out, nil
	}

	return c, &c.GetObjectInvocations
}

func TestDownloadOrder(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient(buf12MB)

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})

	w := manager.NewWriteAtBuffer(make([]byte, len(buf12MB)))
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
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

func TestDownloadZero(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient([]byte{})

	d := manager.NewDownloader(c)
	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if n != 0 {
		t.Errorf("expect 0 bytes read, got %d", n)
	}
	if e, a := 1, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}

	expectRngs := []string{"bytes=0-5242879"}
	if e, a := expectRngs, *ranges; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v ranges, got %v", e, a)
	}
}

func TestDownloadSetPartSize(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient([]byte{1, 2, 3})

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
		d.PartSize = 1
	})
	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
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

func TestDownloadError(t *testing.T) {
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
	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
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

func TestDownloadNonChunk(t *testing.T) {
	c, invocations := newDownloadNonRangeClient(buf2MB)

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})
	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
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

func TestDownloadNoContentRangeLength(t *testing.T) {
	s, invocations, _ := newDownloadRangeClient(buf2MB)

	d := manager.NewDownloader(s, func(d *manager.Downloader) {
		d.Concurrency = 1
	})
	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
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

func TestDownloadContentRangeTotalAny(t *testing.T) {
	s, invocations := newDownloadContentRangeTotalAnyClient(buf2MB)

	d := manager.NewDownloader(s, func(d *manager.Downloader) {
		d.Concurrency = 1
	})
	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
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

func TestDownloadPartBodyRetry_SuccessRetry(t *testing.T) {
	c, invocations := newDownloadWithErrReaderClient([]testErrReader{
		{Buf: []byte("ab"), Len: 3, Err: io.ErrUnexpectedEOF},
		{Buf: []byte("123"), Len: 3, Err: io.EOF},
	})

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})

	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(3), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 2, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}
	if e, a := "123", string(w.Bytes()); e != a {
		t.Errorf("expect %q response, got %q", e, a)
	}
}

func TestDownloadPartBodyRetry_SuccessNoRetry(t *testing.T) {
	c, invocations := newDownloadWithErrReaderClient([]testErrReader{
		{Buf: []byte("abc"), Len: 3, Err: io.EOF},
	})

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
	})

	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
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

func TestDownloadPartBodyRetry_FailRetry(t *testing.T) {
	c, invocations := newDownloadWithErrReaderClient([]testErrReader{
		{Buf: []byte("ab"), Len: 3, Err: io.ErrUnexpectedEOF},
	})

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 1
		d.PartBodyMaxRetries = 0
	})

	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})

	if err == nil {
		t.Fatalf("expect error, got none")
	}
	if e, a := "unexpected EOF", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q error message to be in %q", e, a)
	}
	if e, a := int64(2), n; e != a {
		t.Errorf("expect %d bytes read, got %d", e, a)
	}
	if e, a := 1, *invocations; e != a {
		t.Errorf("expect %v API calls, got %v", e, a)
	}
	if e, a := "ab", string(w.Bytes()); e != a {
		t.Errorf("expect %q response, got %q", e, a)
	}
}

func TestDownloadWithContextCanceled(t *testing.T) {
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

	w := &manager.WriteAtBuffer{}

	_, err := d.Download(ctx, w, &params)
	if err == nil {
		t.Fatalf("expected error, did not get one")
	}
	if e, a := "canceled", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expected error message to contain %q, but did not %q", e, a)
	}
}

func TestDownload_WithRange(t *testing.T) {
	c, invocations, ranges := newDownloadRangeClient([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	d := manager.NewDownloader(c, func(d *manager.Downloader) {
		d.Concurrency = 10 // should be ignored
		d.PartSize = 1     // should be ignored
	})

	w := &manager.WriteAtBuffer{}
	n, err := d.Download(context.Background(), w, &s3.GetObjectInput{
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

type mockDownloadCLient func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)

func (m mockDownloadCLient) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m(ctx, params, optFns...)
}

func TestDownload_WithFailure(t *testing.T) {
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

	w := &manager.WriteAtBuffer{}
	params := s3.GetObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
	}

	// Expect this request to exit quickly after failure
	_, err := d.Download(context.Background(), w, &params)
	if err == nil {
		t.Fatalf("expect error, got none")
	}

	if atomic.LoadInt64(&reqCount) > 3 {
		t.Errorf("expect no more than 3 requests, but received %d", reqCount)
	}
}

func TestDownloadBufferStrategy(t *testing.T) {
	cases := map[string]struct {
		partSize     int64
		strategy     *recordedWriterReadFromProvider
		expectedSize int64
	}{
		"no strategy": {
			partSize:     manager.DefaultDownloadPartSize,
			expectedSize: 10 * sdkio.MebiByte,
		},
		"partSize modulo bufferSize == 0": {
			partSize: 5 * sdkio.MebiByte,
			strategy: &recordedWriterReadFromProvider{
				WriterReadFromProvider: manager.NewPooledBufferedWriterReadFromProvider(int(sdkio.MebiByte)), // 1 MiB
			},
			expectedSize: 10 * sdkio.MebiByte, // 10 MiB
		},
		"partSize modulo bufferSize > 0": {
			partSize: 5 * 1024 * 1204, // 5 MiB
			strategy: &recordedWriterReadFromProvider{
				WriterReadFromProvider: manager.NewPooledBufferedWriterReadFromProvider(2 * int(sdkio.MebiByte)), // 2 MiB
			},
			expectedSize: 10 * sdkio.MebiByte, // 10 MiB
		},
	}

	for name, tCase := range cases {
		t.Run(name, func(t *testing.T) {
			expected := managertesting.GetTestBytes(int(tCase.expectedSize))

			client, _, _ := newDownloadRangeClient(expected)

			d := manager.NewDownloader(client, func(d *manager.Downloader) {
				d.PartSize = tCase.partSize
				if tCase.strategy != nil {
					d.BufferProvider = tCase.strategy
				}
			})

			buffer := manager.NewWriteAtBuffer(make([]byte, len(expected)))

			n, err := d.Download(context.Background(), buffer, &s3.GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			})
			if err != nil {
				t.Errorf("failed to download: %v", err)
			}

			if e, a := len(expected), int(n); e != a {
				t.Errorf("expected %v, got %v downloaded bytes", e, a)
			}

			if e, a := expected, buffer.Bytes(); !bytes.Equal(e, a) {
				t.Errorf("downloaded bytes did not match expected")
			}

			if tCase.strategy != nil {
				if e, a := tCase.strategy.callbacksVended, tCase.strategy.callbacksExecuted; e != a {
					t.Errorf("expected %v, got %v", e, a)
				}
			}
		})
	}
}

type testErrReader struct {
	Buf []byte
	Err error
	Len int64

	off int
}

func (r *testErrReader) Read(p []byte) (int, error) {
	to := len(r.Buf) - r.off

	n := copy(p, r.Buf[r.off:to])
	r.off += n

	if n < len(p) {
		return n, r.Err

	}

	return n, nil
}

func TestDownloadBufferStrategy_Errors(t *testing.T) {
	expected := managertesting.GetTestBytes(int(10 * sdkio.MebiByte))

	client, _, _ := newDownloadRangeClient(expected)
	strat := &recordedWriterReadFromProvider{
		WriterReadFromProvider: manager.NewPooledBufferedWriterReadFromProvider(int(2 * sdkio.MebiByte)),
	}

	seenOps := make(map[string]struct{})
	orig := client.GetObjectFn
	client.GetObjectFn = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		out, err := orig(ctx, params, optFns...)

		fingerPrint := fmt.Sprintf("%s/%s/%s", *params.Bucket, *params.Key, *params.Range)
		if _, ok := seenOps[fingerPrint]; ok {
			return out, err
		}
		seenOps[fingerPrint] = struct{}{}

		_, _ = io.Copy(ioutil.Discard, out.Body)

		out.Body = ioutil.NopCloser(&badReader{err: io.ErrUnexpectedEOF})

		return out, err
	}

	d := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.PartSize = 5 * sdkio.MebiByte
		d.BufferProvider = strat
		d.Concurrency = 1
	})

	buffer := manager.NewWriteAtBuffer(make([]byte, len(expected)))

	n, err := d.Download(context.Background(), buffer, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Errorf("failed to download: %v", err)
	}

	if e, a := len(expected), int(n); e != a {
		t.Errorf("expected %v, got %v downloaded bytes", e, a)
	}

	if e, a := expected, buffer.Bytes(); !bytes.Equal(e, a) {
		t.Errorf("downloaded bytes did not match expected")
	}

	if e, a := strat.callbacksVended, strat.callbacksExecuted; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}
}

func TestDownloaderValidARN(t *testing.T) {
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

			_, err := downloader.Download(context.Background(), &awstesting.DiscardAt{}, &tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("err: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

type recordedWriterReadFromProvider struct {
	callbacksVended   uint32
	callbacksExecuted uint32
	manager.WriterReadFromProvider
}

func (r *recordedWriterReadFromProvider) GetReadFrom(writer io.Writer) (manager.WriterReadFrom, func()) {
	w, cleanup := r.WriterReadFromProvider.GetReadFrom(writer)

	atomic.AddUint32(&r.callbacksVended, 1)
	return w, func() {
		atomic.AddUint32(&r.callbacksExecuted, 1)
		cleanup()
	}
}

type badReader struct {
	err error
}

func (b *badReader) Read(p []byte) (int, error) {
	tb := managertesting.GetTestBytes(len(p))
	copy(p, tb)

	return len(p), b.err
}
