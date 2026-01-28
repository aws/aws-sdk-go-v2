package transfermanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var etag string = "myetag"
var vID string = "myversion"

func TestGetObject(t *testing.T) {
	cases := map[string]struct {
		data              []byte
		errReaders        []s3testing.TestErrReader
		getObjectFn       func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error)
		optFn             func(*Options)
		versionID         string
		checksumType      s3types.ChecksumType
		expectInvocations int
		expectRanges      []string
		expectVersions    []string
		expectETags       []string
		partsCount        int32
		expectParts       []int32
		expectGetErr      string
		expectReadErr     string
		dataValidationFn  func(*testing.T, []byte)
	}{
		"range download in order": {
			data:        buf20MB,
			getObjectFn: s3testing.RangeGetObjectFn,
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
				o.Concurrency = 1
			},
			expectInvocations: 3,
			expectRanges:      []string{"bytes=0-8388607", "bytes=8388608-16777215", "bytes=16777216-20971519"},
			expectETags:       []string{etag, etag, etag},
		},
		"range download zero": {
			data:        []byte{},
			getObjectFn: s3testing.NonRangeGetObjectFn,
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
			},
			expectInvocations: 1,
		},
		"range download with customized part size and versionID": {
			data:        buf20MB,
			getObjectFn: s3testing.RangeGetObjectFn,
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
				o.PartSizeBytes = 10 * 1024 * 1024
				o.Concurrency = 1
			},
			versionID:         vID,
			expectInvocations: 2,
			expectRanges:      []string{"bytes=0-10485759", "bytes=10485760-20971519"},
			expectVersions:    []string{vID, vID},
		},
		"range download with s3 error": {
			data:        buf20MB,
			getObjectFn: s3testing.ErrRangeGetObjectFn,
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
				o.Concurrency = 1
			},
			expectInvocations: 2,
			expectReadErr:     "s3 service error",
		},
		"range download with content mismatch error": {
			data:        buf20MB,
			getObjectFn: s3testing.MismatchRangeGetObjectFn,
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
				o.Concurrency = 1
			},
			expectInvocations: 2,
			expectReadErr:     "PreconditionFailed",
		},
		"range download with resp range mismatch error": {
			data:        buf20MB,
			getObjectFn: s3testing.WrongRangeGetObjectFn,
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
				o.Concurrency = 1
			},
			expectInvocations: 2,
			expectReadErr:     "range mismatch between request",
		},
		"content length download single chunk": {
			data:        buf2MB,
			getObjectFn: s3testing.NonRangeGetObjectFn,
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
			},
			expectInvocations: 1,
			expectRanges:      []string{"bytes=0-2097151"},
			expectETags:       []string{etag},
			dataValidationFn: func(t *testing.T, bytes []byte) {
				count := 0
				for _, b := range bytes {
					count += int(b)
				}
				if count != 0 {
					t.Errorf("expect 0 count, got %d", count)
				}
			},
		},
		"range download single chunk": {
			data:        buf2MB,
			getObjectFn: s3testing.RangeGetObjectFn,
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
			},
			versionID:         vID,
			expectInvocations: 1,
			expectVersions:    []string{vID},
			expectRanges:      []string{"bytes=0-2097151"},
			dataValidationFn: func(t *testing.T, bytes []byte) {
				count := 0
				for _, b := range bytes {
					count += int(b)
				}
				if count != 0 {
					t.Errorf("expect 0 count, got %d", count)
				}
			},
		},
		"range download success without retry": {
			data:        []byte("123"),
			getObjectFn: s3testing.ErrReaderFn,
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("123"), Len: 3, Err: io.EOF},
			},
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
			},
			expectInvocations: 1,
			dataValidationFn: func(t *testing.T, bytes []byte) {
				if e, a := "123", string(bytes); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
		},
		"range download fail retry": {
			data:        []byte("ab"),
			getObjectFn: s3testing.ErrReaderFn,
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("ab"), Len: 2, Err: io.ErrUnexpectedEOF},
			},
			optFn: func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
			},
			expectInvocations: 1,
			expectReadErr:     "unexpected EOF",
		},
		"parts download in order": {
			data:        buf2MB,
			getObjectFn: s3testing.PartGetObjectFn,
			optFn: func(o *Options) {
				o.Concurrency = 1
			},
			partsCount:        3,
			expectInvocations: 3,
			expectETags:       []string{etag, etag, etag},
			expectParts:       []int32{1, 2, 3},
		},
		"parts download with composite checksum type": {
			data:        buf2MB,
			getObjectFn: s3testing.PartGetObjectFn,
			optFn: func(o *Options) {
				o.Concurrency = 1
			},
			checksumType:      s3types.ChecksumTypeComposite,
			partsCount:        3,
			expectInvocations: 3,
			expectETags:       []string{etag, etag, etag},
			expectParts:       []int32{1, 2, 3},
		},
		"parts download with version ID": {
			data:              buf2MB,
			getObjectFn:       s3testing.PartGetObjectFn,
			optFn:             func(o *Options) {},
			partsCount:        3,
			versionID:         vID,
			expectInvocations: 3,
			expectVersions:    []string{vID, vID, vID},
		},
		"part download zero": {
			data:              buf2MB,
			getObjectFn:       s3testing.PartGetObjectFn,
			optFn:             func(o *Options) {},
			partsCount:        1,
			expectInvocations: 1,
			expectParts:       []int32{1},
		},
		"part download with s3 error": {
			data:        buf2MB,
			getObjectFn: s3testing.ErrPartGetObjectFn,
			optFn: func(o *Options) {
				o.Concurrency = 1
			},
			partsCount:        3,
			expectInvocations: 2,
			expectReadErr:     "s3 service error",
		},
		"part download with mismatch error": {
			data:        buf2MB,
			getObjectFn: s3testing.MismatchPartGetObjectFn,
			optFn: func(o *Options) {
				o.Concurrency = 1
			},
			partsCount:        3,
			expectInvocations: 2,
			expectReadErr:     "PreconditionFailed",
		},
		"part download single chunk": {
			data:              []byte("123"),
			getObjectFn:       s3testing.PartGetObjectFn,
			optFn:             func(o *Options) {},
			partsCount:        1,
			expectInvocations: 1,
			expectParts:       []int32{1},
			dataValidationFn: func(t *testing.T, bytes []byte) {
				if e, a := "123", string(bytes); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
		},
		"part download success without retry": {
			getObjectFn: s3testing.ErrReaderFn,
			optFn:       func(o *Options) {},
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("ab"), Len: 2, Err: io.EOF},
			},
			partsCount:        1,
			expectInvocations: 1,
			expectParts:       []int32{1},
			dataValidationFn: func(t *testing.T, bytes []byte) {
				if e, a := "ab", string(bytes); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
			},
		},
		"part download fail retry": {
			data:        []byte("ab"),
			getObjectFn: s3testing.ErrReaderFn,
			optFn:       func(o *Options) {},
			errReaders: []s3testing.TestErrReader{
				{Buf: []byte("ab"), Len: 2, Err: io.ErrUnexpectedEOF},
			},
			expectInvocations: 1,
			expectReadErr:     "unexpected EOF",
			dataValidationFn: func(t *testing.T, bytes []byte) {
				if e, a := "ab", string(bytes); e != a {
					t.Errorf("expect %q response, got %q", e, a)
				}
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
			s3Client.ChecksumType = c.checksumType
			mgr := New(s3Client, c.optFn)

			input := &GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			}
			input.VersionID = nzstring(c.versionID)

			out, err := mgr.GetObject(context.Background(), input)

			if err != nil {
				if c.expectGetErr == "" {
					t.Fatalf("expect no error when getting object, got %q", err)
				} else if e, a := c.expectGetErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %s error message to be in %s", e, a)
				}
			} else if c.expectGetErr != "" {
				t.Fatal("expect error when getting object, got nil")
			}

			if err != nil {
				return
			}

			actualBuf, err := io.ReadAll(out.Body)
			if err != nil {
				if c.expectReadErr == "" {
					t.Fatalf("expect no error when reading response, got %q", err)
				} else if e, a := c.expectReadErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %s error message to be in %s", e, a)
				}
			} else if c.expectReadErr != "" {
				t.Fatal("expect error when reading response, got nil")
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

			if c.checksumType == s3types.ChecksumTypeComposite {
				if out.ChecksumCRC32 != nil || out.ChecksumCRC32C != nil || out.ChecksumCRC64NVME != nil ||
					out.ChecksumSHA1 != nil || out.ChecksumSHA256 != nil {
					t.Errorf("expect all composite checksum value to be empty in output, got non-empty value: %s, %s, %s, %s, %s",
						aws.ToString(out.ChecksumCRC32), aws.ToString(out.ChecksumCRC32C), aws.ToString(out.ChecksumCRC64NVME),
						aws.ToString(out.ChecksumSHA1), aws.ToString(out.ChecksumSHA256))
				}
			}

			if c.dataValidationFn != nil {
				c.dataValidationFn(t, actualBuf)
			}
		})
	}
}

func TestGetAsyncWithFailure(t *testing.T) {
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
			reqCount := atomic.Int64{}

			s3Client := &s3testing.TransferManagerLoggingClient{}
			s3Client.PartsCount = 10
			s3Client.Data = buf80MB
			s3Client.GetObjectFn = func(c *s3testing.TransferManagerLoggingClient, params *s3.GetObjectInput) (out *s3.GetObjectOutput, err error) {
				count := reqCount.Load()
				switch count {
				case 1:
					// Give a chance for the multipart chunks to be queued up
					time.Sleep(1 * time.Second)
					err = fmt.Errorf("some connection error")
				default:
					var start, end int64
					if params.Range != nil {
						start, end, err = getReqRange(aws.ToString(params.Range))
						if err != nil {
							return
						}
					}
					body := bytes.NewReader(make([]byte, minPartSizeBytes))
					out = &s3.GetObjectOutput{
						Body:          io.NopCloser(body),
						ContentLength: aws.Int64(int64(body.Len())),
						ContentRange:  aws.String(fmt.Sprintf("bytes %d-%d/%d", start, end, body.Len()*10)),
					}

					if count > 0 {
						// sleep here to ensure context switching between goroutines
						time.Sleep(25 * time.Millisecond)
					}
				}

				reqCount.Add(1)
				return out, err
			}

			mgr := New(s3Client, func(o *Options) {
				o.Concurrency = 2
				o.GetObjectType = c.downloadType
			})

			// Expect this request to exit quickly after failure
			out, err := mgr.GetObject(context.Background(), &GetObjectInput{
				Bucket: aws.String("Bucket"),
				Key:    aws.String("Key"),
			})
			_, err = io.ReadAll(out.Body)

			if err == nil {
				t.Fatal("expect error, got none")
			} else if e, a := "some connection error", err.Error(); !strings.Contains(a, e) {
				t.Fatalf("expect %s error message to be in %s", e, a)
			}

			if count := reqCount.Load(); count > 3 {
				t.Errorf("expect no more than 3 requests, but received %d", count)
			}
		})
	}
}

func TestGetObjectWithContextCanceled(t *testing.T) {
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
			mgr := New(s3.New(s3.Options{
				Region: "mock-region",
			}), func(o *Options) {
				o.GetObjectType = c.downloadType
			})

			ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
			ctx.Error = fmt.Errorf("context canceled")
			close(ctx.DoneCh)

			_, err := mgr.GetObject(ctx, &GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("Key"),
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
