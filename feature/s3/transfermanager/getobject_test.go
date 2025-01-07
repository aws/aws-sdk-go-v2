package transfermanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestDownloadObject(t *testing.T) {
	cases := map[string]struct {
		data              []byte
		options           Options
		loggingClientFn   func(*s3testing.TransferManagerLoggingClient)
		downloadInputFn   func(*GetObjectInput)
		expectInvocations int
		expectRanges      []string
		expectErr         string
		dataValidationFn  func(t *testing.T, w *types.WriteAtBuffer)
	}{
		"range download in order": {
			options: Options{
				Concurrency:           1,
				MultipartDownloadType: types.MultipartDownloadTypeRange,
			},
			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.RangeGetObjectFn
				c.Data = buf20MB
			},
			expectInvocations: 3,
			expectRanges:      []string{"bytes=0-8388607", "bytes=8388608-16777215", "bytes=16777216-20971519"},
		},
		"range download zero": {
			options: Options{
				MultipartDownloadType: types.MultipartDownloadTypeRange,
			},
			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.RangeGetObjectFn
				c.Data = []byte{}
			},
			expectInvocations: 1,
			expectRanges:      []string{"bytes=0-8388607"},
		},
		"range download with customized part size": {
			options: Options{
				Concurrency:           1,
				MultipartDownloadType: types.MultipartDownloadTypeRange,
				PartSizeBytes:         10 * 1024 * 1024,
			},
			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.RangeGetObjectFn
				c.Data = buf20MB
			},
			expectInvocations: 2,
			expectRanges:      []string{"bytes=0-10485759", "bytes=10485760-20971519"},
		},
		"range download with s3 error": {
			options: Options{
				Concurrency:           1,
				MultipartDownloadType: types.MultipartDownloadTypeRange,
			},
			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.ErrGetObjectFn
				c.Data = buf20MB
			},
			expectInvocations: 2,
			expectErr:         "s3 service error",
		},
		"content length download single chunk": {
			options: Options{
				MultipartDownloadType: types.MultipartDownloadTypeRange,
			},
			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.NonRangeGetObjectFn
				c.Data = buf2MB
			},
			expectInvocations: 1,
			expectRanges:      []string{"bytes=0-8388607"},
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				count := 0
				for _, b := range w.Bytes() {
					count += int(b)
				}
				if count != 0 {
					t.Errorf("expect 0 count, got %d", count)
				}
			},
		},
		"range download single chunk": {
			options: Options{
				Concurrency:           1,
				MultipartDownloadType: types.MultipartDownloadTypeRange,
			},
			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.RangeGetObjectFn
				c.Data = buf2MB
			},
			expectInvocations: 1,
			expectRanges:      []string{"bytes=0-8388607"},
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				count := 0
				for _, b := range w.Bytes() {
					count += int(b)
				}
				if count != 0 {
					t.Errorf("expect 0 count, got %d", count)
				}
			},
		},
		"range download with success retry": {
			options: Options{
				Concurrency:           1,
				MultipartDownloadType: types.MultipartDownloadTypeRange,
			},

			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.ErrReaderFn
				c.ErrReaders = []s3testing.TestErrReader{
					{Buf: []byte("ab"), Len: 3, Err: io.ErrUnexpectedEOF},
					{Buf: []byte("123"), Len: 3, Err: io.EOF},
				}
			},
			expectInvocations: 2,
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "123", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
		},
		"range download success without retry": {
			options: Options{
				Concurrency:           1,
				MultipartDownloadType: types.MultipartDownloadTypeRange,
			},

			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.ErrReaderFn
				c.ErrReaders = []s3testing.TestErrReader{
					{Buf: []byte("123"), Len: 3, Err: io.EOF},
				}
			},
			expectInvocations: 1,
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "123", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
		},
		"range download fail retry": {
			options: Options{
				Concurrency:           1,
				PartBodyMaxRetries:    1,
				MultipartDownloadType: types.MultipartDownloadTypeRange,
			},
			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.ErrReaderFn
				c.ErrReaders = []s3testing.TestErrReader{
					{Buf: []byte("ab"), Len: 3, Err: io.ErrUnexpectedEOF},
				}
			},
			expectInvocations: 1,
			expectErr:         "unexpected EOF",
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "ab", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
		},
		"range download a range of object": {
			options: Options{
				Concurrency:           1,
				MultipartDownloadType: types.MultipartDownloadTypeRange,
			},
			loggingClientFn: func(c *s3testing.TransferManagerLoggingClient) {
				c.GetObjectFn = s3testing.RangeGetObjectFn
				c.Data = buf20MB
			},
			downloadInputFn: func(input *GetObjectInput) {
				input.Range = "bytes=0-10485759"
			},
			expectInvocations: 2,
			expectRanges:      []string{"bytes=0-8388607", "bytes=8388608-10485759"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s3Client, invocations, ranges := s3testing.NewDownloadClient()
			if c.loggingClientFn != nil {
				c.loggingClientFn(s3Client)
			}
			mgr := New(s3Client, c.options)
			w := types.NewWriteAtBuffer(make([]byte, 0))

			ctx := context.Background()
			input := &GetObjectInput{
				Bucket: "bucket",
				Key:    "key",
			}
			if c.downloadInputFn != nil {
				c.downloadInputFn(input)
			}
			_, err := mgr.DownloadObject(ctx, w, input)
			if err != nil {
				if c.expectErr == "" {
					t.Fatalf("expect no error, got %v", err)
				} else if e, a := c.expectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %s error message to be in %s", e, a)
				}
			} else {
				if c.expectErr != "" {
					t.Fatal("expect error, got nil")
				}
			}

			if e, a := c.expectInvocations, *invocations; e != a {
				t.Errorf("expect %v API calls, got %v", e, a)
			}

			if len(c.expectRanges) > 0 {
				if e, a := c.expectRanges, *ranges; !reflect.DeepEqual(e, a) {
					t.Errorf("expect %v ranges, got %v", e, a)
				}
			}

			if c.dataValidationFn != nil {
				c.dataValidationFn(t, w)
			}
		})
	}
}

func TestDownload_WithFailure(t *testing.T) {
	startingByte := 0
	reqCount := int64(0)

	s3Client, _, _ := s3testing.NewDownloadClient()
	s3Client.GetObjectFn = func(c *s3testing.TransferManagerLoggingClient, params *s3.GetObjectInput) (out *s3.GetObjectOutput, err error) {
		switch atomic.LoadInt64(&reqCount) {
		case 1:
			// Give a chance for the multipart chunks to be queued up
			time.Sleep(1 * time.Second)
			err = fmt.Errorf("some connection error")
		default:
			body := bytes.NewReader(make([]byte, minPartSizeBytes))
			out = &s3.GetObjectOutput{
				Body:          ioutil.NopCloser(body),
				ContentLength: aws.Int64(int64(body.Len())),
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
	}

	d := New(s3Client, Options{
		Concurrency:           2,
		MultipartDownloadType: types.MultipartDownloadTypeRange,
	})

	w := types.NewWriteAtBuffer(make([]byte, 0))

	// Expect this request to exit quickly after failure
	_, err := d.DownloadObject(context.Background(), w, &GetObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
	})
	if err == nil {
		t.Fatal("expect error, got none")
	} else if e, a := "some connection error", err.Error(); !strings.Contains(a, e) {
		t.Fatalf("expect %s error message to be in %s", e, a)
	}

	if atomic.LoadInt64(&reqCount) > 3 {
		t.Errorf("expect no more than 3 requests, but received %d", reqCount)
	}
}
