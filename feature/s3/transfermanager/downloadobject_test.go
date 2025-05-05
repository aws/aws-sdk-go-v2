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
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const megabyte = 1024 * 1024

func TestDownloadObject(t *testing.T) {
	cases := map[string]struct {
		data                 []byte
		errReaders           []s3testing.TestErrReader
		getObjectFn          func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error)
		options              Options
		downloadRange        string
		expectInvocations    int
		expectRanges         []string
		partNumber           int32
		versionID            string
		partsCount           int32
		expectParts          []int32
		expectVersions       []string
		expectETags          []string
		expectErr            string
		dataValidationFn     func(*testing.T, *types.WriteAtBuffer)
		listenerValidationFn func(*testing.T, *mockListener, any, any, error)
	}{
		"range download in order": {
			data:        buf20MB,
			getObjectFn: s3testing.RangeGetObjectFn,
			options: Options{
				Concurrency:   1,
				GetObjectType: types.GetObjectRanges,
			},
			expectInvocations: 3,
			expectRanges:      []string{"bytes=0-8388607", "bytes=8388608-16777215", "bytes=16777216-20971519"},
			expectETags:       []string{"", etag, etag},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectStartTotalBytes(t, 20*megabyte)
				l.expectByteTransfers(t,
					8*megabyte, 16*megabyte, 20*megabyte)
			},
		},
		"range download zero": {
			data:        []byte{},
			getObjectFn: s3testing.RangeGetObjectFn,
			options: Options{
				GetObjectType: types.GetObjectRanges,
			},
			expectInvocations: 1,
			expectRanges:      []string{"bytes=0-8388607"},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectStartTotalBytes(t, 0)
				l.expectByteTransfers(t, 0)
			},
		},
		"range download with customized part size with version ID": {
			data:        buf20MB,
			getObjectFn: s3testing.RangeGetObjectFn,
			options: Options{
				Concurrency:   1,
				GetObjectType: types.GetObjectRanges,
				PartSizeBytes: 10 * 1024 * 1024,
			},
			versionID:         vID,
			expectInvocations: 2,
			expectRanges:      []string{"bytes=0-10485759", "bytes=10485760-20971519"},
			expectVersions:    []string{vID, vID},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectStartTotalBytes(t, 20*megabyte)
				l.expectByteTransfers(t,
					10*megabyte, 20*megabyte)
			},
		},
		"range download with s3 error": {
			data:        buf20MB,
			getObjectFn: s3testing.ErrRangeGetObjectFn,
			options: Options{
				Concurrency:   1,
				GetObjectType: types.GetObjectRanges,
			},
			expectInvocations: 2,
			expectErr:         "s3 service error",
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectStartTotalBytes(t, 20*megabyte)
				l.expectFailed(t, in, err)
			},
		},
		"range download with mismatch error": {
			data:        buf20MB,
			getObjectFn: s3testing.MismatchRangeGetObjectFn,
			options: Options{
				Concurrency:   1,
				GetObjectType: types.GetObjectRanges,
			},
			expectInvocations: 2,
			expectErr:         "PreconditionFailed",
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectStartTotalBytes(t, 20*megabyte)
				l.expectFailed(t, in, err)
			},
		},
		"content length download single chunk": {
			data:        buf2MB,
			getObjectFn: s3testing.NonRangeGetObjectFn,
			options: Options{
				GetObjectType: types.GetObjectRanges,
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
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 2*megabyte)
			},
		},
		"range download single chunk with version ID": {
			data:        buf2MB,
			getObjectFn: s3testing.RangeGetObjectFn,
			options: Options{
				Concurrency:   1,
				GetObjectType: types.GetObjectRanges,
			},
			versionID:         vID,
			expectInvocations: 1,
			expectVersions:    []string{vID},
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
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 2*megabyte)
			},
		},
		"range download with success retry": {
			getObjectFn: s3testing.ErrReaderFn,
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("ab"), Len: 3, Err: io.ErrUnexpectedEOF},
				{Buf: []byte("123"), Len: 3, Err: io.EOF},
			},
			options: Options{
				Concurrency:   1,
				GetObjectType: types.GetObjectRanges,
			},
			expectInvocations: 2,
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "123", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 3)
			},
		},
		"range download success without retry": {
			getObjectFn: s3testing.ErrReaderFn,
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("123"), Len: 3, Err: io.EOF},
			},
			options: Options{
				Concurrency:   1,
				GetObjectType: types.GetObjectRanges,
			},
			expectInvocations: 1,
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "123", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 3)
			},
		},
		"range download fail retry": {
			getObjectFn: s3testing.ErrReaderFn,
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("ab"), Len: 3, Err: io.ErrUnexpectedEOF},
			},
			options: Options{
				Concurrency:        1,
				PartBodyMaxRetries: 1,
				GetObjectType:      types.GetObjectRanges,
			},
			expectInvocations: 1,
			expectErr:         "unexpected EOF",
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "ab", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectStartTotalBytes(t, 3)
				// no transferred because the first chunk blows up
				l.expectFailed(t, in, err)
			},
		},
		"range download a range of object": {
			data:        buf20MB,
			getObjectFn: s3testing.RangeGetObjectFn,
			options: Options{
				Concurrency:   1,
				GetObjectType: types.GetObjectRanges,
			},
			downloadRange:     "bytes=1-10485759",
			expectInvocations: 2,
			expectRanges:      []string{"bytes=1-8388608", "bytes=8388609-10485759"},
			expectETags:       []string{"", etag},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 8*megabyte, 10*megabyte-1)
			},
		},
		"range download a range of object with version ID": {
			data:        buf20MB,
			getObjectFn: s3testing.RangeGetObjectFn,
			options: Options{
				GetObjectType: types.GetObjectRanges,
			},
			downloadRange:     "bytes=0-10485759",
			versionID:         vID,
			expectInvocations: 2,
			expectVersions:    []string{vID, vID},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 8*megabyte, 10*megabyte)
			},
		},
		"parts download in order": {
			data:        buf2MB,
			getObjectFn: s3testing.PartGetObjectFn,
			options: Options{
				Concurrency: 1,
			},
			partsCount:        3,
			versionID:         vID,
			expectInvocations: 3,
			expectVersions:    []string{vID, vID, vID},
			expectParts:       []int32{1, 2, 3},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 2*megabyte, 4*megabyte, 6*megabyte)
			},
		},
		"part download zero": {
			data:              buf2MB,
			getObjectFn:       s3testing.PartGetObjectFn,
			options:           Options{},
			partsCount:        1,
			expectInvocations: 1,
			expectParts:       []int32{1},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 2*megabyte)
			},
		},
		"part download with s3 error": {
			data:        buf2MB,
			getObjectFn: s3testing.ErrPartGetObjectFn,
			options: Options{
				Concurrency: 1,
			},
			partsCount:        3,
			expectInvocations: 2,
			expectErr:         "s3 service error",
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectFailed(t, in, err)
			},
		},
		"part download with mismatch error": {
			data:        buf2MB,
			getObjectFn: s3testing.MismatchPartGetObjectFn,
			options: Options{
				Concurrency: 1,
			},
			partsCount:        3,
			expectInvocations: 2,
			expectErr:         "PreconditionFailed",
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectFailed(t, in, err)
			},
		},
		"part download single chunk": {
			data:              []byte("123"),
			getObjectFn:       s3testing.PartGetObjectFn,
			options:           Options{},
			partsCount:        1,
			expectInvocations: 1,
			expectParts:       []int32{1},
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "123", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 3)
			},
		},
		"part download with success retry": {
			getObjectFn: s3testing.ErrReaderFn,
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("ab"), Len: 3, Err: io.ErrUnexpectedEOF},
				{Buf: []byte("123"), Len: 3, Err: io.EOF},
			},
			options: Options{
				Concurrency: 1,
			},
			partsCount:        1,
			expectInvocations: 2,
			expectParts:       []int32{1, 1},
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "123", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 3)
			},
		},
		"part download success without retry": {
			getObjectFn: s3testing.ErrReaderFn,
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("ab"), Len: 3, Err: io.EOF},
			},
			options: Options{
				Concurrency: 1,
			},
			partsCount:        1,
			expectInvocations: 1,
			expectParts:       []int32{1},
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "ab", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 2)
			},
		},
		"part download fail retry": {
			getObjectFn: s3testing.ErrReaderFn,
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("ab"), Len: 3, Err: io.ErrUnexpectedEOF},
			},
			options: Options{
				Concurrency:        1,
				PartBodyMaxRetries: 1,
			},
			expectInvocations: 1,
			expectErr:         "unexpected EOF",
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "ab", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectStartTotalBytes(t, 3)
				// no transferred because the first chunk blows up
				l.expectFailed(t, in, err)
			},
		},
		"parts download with range input": {
			data:              []byte("123"),
			getObjectFn:       s3testing.PartGetObjectFn,
			options:           Options{},
			downloadRange:     "bytes=0-100",
			partsCount:        3,
			expectInvocations: 1,
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "123", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 3)
			},
		},
		"parts download with part number input": {
			data:              []byte("ab"),
			getObjectFn:       s3testing.PartGetObjectFn,
			options:           Options{},
			partsCount:        3,
			partNumber:        5,
			expectInvocations: 1,
			dataValidationFn: func(t *testing.T, w *types.WriteAtBuffer) {
				if e, a := "ab", string(w.Bytes()); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
			listenerValidationFn: func(t *testing.T, l *mockListener, in, out any, err error) {
				l.expectComplete(t, in, out)
				l.expectByteTransfers(t, 2)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s3Client, invocations, parts, ranges, versions, etags := s3testing.NewDownloadClient()
			s3Client.Data = c.data
			s3Client.GetObjectFn = c.getObjectFn
			s3Client.ErrReaders = c.errReaders
			s3Client.PartsCount = c.partsCount
			mgr := New(s3Client, c.options)
			w := types.NewWriteAtBuffer(make([]byte, 0))

			input := &DownloadObjectInput{
				Bucket:     "bucket",
				Key:        "key",
				WriterAt:   w,
				Range:      c.downloadRange,
				PartNumber: c.partNumber,
				VersionID:  c.versionID,
			}

			listener := &mockListener{}

			output, err := mgr.DownloadObject(context.Background(), input, func(o *Options) {
				o.ProgressListeners.Register(listener)
			})
			if err != nil {
				if c.expectErr == "" {
					t.Fatalf("expect no error, got %q", err)
				} else if e, a := c.expectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %s error message to be in %s", e, a)
				}
			} else {
				if c.expectErr != "" {
					t.Fatal("expect error, got nil")
				}
			}

			// there can be expects in both success and failure
			if c.listenerValidationFn != nil {
				c.listenerValidationFn(t, listener, input, output, err)
			}

			if err != nil {
				return
			}

			if e, a := c.expectInvocations, *invocations; e != a {
				t.Errorf("expect %v API calls, got %v", e, a)
			}

			if len(c.expectParts) > 0 {
				if e, a := c.expectParts, *parts; !reflect.DeepEqual(e, a) {
					t.Errorf("expect %v parts, got %v", e, a)
				}
			}
			if len(c.expectRanges) > 0 {
				if e, a := c.expectRanges, *ranges; !reflect.DeepEqual(e, a) {
					t.Errorf("expect %v ranges, got %v", e, a)
				}
			}
			if len(c.expectVersions) > 0 {
				if e, a := c.expectVersions, *versions; !reflect.DeepEqual(e, a) {
					t.Errorf("expect %v versions, got %v", e, a)
				}
			}
			if len(c.expectETags) > 0 {
				if e, a := c.expectETags, *etags; !reflect.DeepEqual(e, a) {
					t.Errorf("expect %v etags, got %v", e, a)
				}
			}

			if c.dataValidationFn != nil {
				c.dataValidationFn(t, w)
			}
		})
	}
}

func TestDownloadAsyncWithFailure(t *testing.T) {
	cases := map[string]struct {
		downloadType types.GetObjectType
	}{
		"part download by default": {},
		"range download": {
			downloadType: types.GetObjectRanges,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			startingByte := 0
			reqCount := int64(0)

			s3Client := &s3testing.TransferManagerLoggingClient{}
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
						PartsCount:    aws.Int32(10),
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
				Concurrency:   2,
				GetObjectType: c.downloadType,
			})

			w := types.NewWriteAtBuffer(make([]byte, 0))

			// Expect this request to exit quickly after failure
			_, err := d.DownloadObject(context.Background(), &DownloadObjectInput{
				Bucket:   "Bucket",
				Key:      "Key",
				WriterAt: w,
			})
			if err == nil {
				t.Fatal("expect error, got none")
			} else if e, a := "some connection error", err.Error(); !strings.Contains(a, e) {
				t.Fatalf("expect %s error message to be in %s", e, a)
			}

			if atomic.LoadInt64(&reqCount) > 3 {
				t.Errorf("expect no more than 3 requests, but received %d", reqCount)
			}
		})
	}
}

func TestDownloadObjectWithContextCanceled(t *testing.T) {
	cases := map[string]struct {
		downloadType types.GetObjectType
	}{
		"part download by default": {},
		"range download": {
			downloadType: types.GetObjectRanges,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			d := New(s3.New(s3.Options{
				Region: "mock-region",
			}), Options{
				GetObjectType: c.downloadType,
			})

			ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
			ctx.Error = fmt.Errorf("context canceled")
			close(ctx.DoneCh)

			w := types.NewWriteAtBuffer(make([]byte, 0))

			_, err := d.DownloadObject(ctx, &DownloadObjectInput{
				Bucket:   "bucket",
				Key:      "Key",
				WriterAt: w,
			})
			if err == nil {
				t.Fatalf("expected error, did not get one")
			}
			if e, a := "canceled", err.Error(); !strings.Contains(a, e) {
				t.Errorf("expected error message to contain %q, but did not %q", e, a)
			}
		})
	}
}
