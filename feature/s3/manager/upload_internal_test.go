package manager

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/manager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type testReader struct {
	br *bytes.Reader
	m  sync.Mutex
}

func (r *testReader) Read(p []byte) (n int, err error) {
	r.m.Lock()
	defer r.m.Unlock()
	return r.br.Read(p)
}

func TestUploadByteSlicePool(t *testing.T) {
	cases := map[string]struct {
		PartSize      int64
		FileSize      int64
		Concurrency   int
		ExAllocations uint64
	}{
		"single part, single concurrency": {
			PartSize:      sdkio.MebiByte * 5,
			FileSize:      sdkio.MebiByte * 5,
			ExAllocations: 2,
			Concurrency:   1,
		},
		"multi-part, single concurrency": {
			PartSize:      sdkio.MebiByte * 5,
			FileSize:      sdkio.MebiByte * 10,
			ExAllocations: 2,
			Concurrency:   1,
		},
		"multi-part, multiple concurrency": {
			PartSize:      sdkio.MebiByte * 5,
			FileSize:      sdkio.MebiByte * 20,
			ExAllocations: 3,
			Concurrency:   2,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			var p *recordedPartPool

			unswap := swapByteSlicePool(func(sliceSize int64) byteSlicePool {
				p = newRecordedPartPool(sliceSize)
				return p
			})
			defer unswap()

			client, _, _ := s3testing.NewUploadLoggingClient(nil)

			uploader := NewUploader(client, func(u *Uploader) {
				u.PartSize = tt.PartSize
				u.Concurrency = tt.Concurrency
			})

			expected := s3testing.GetTestBytes(int(tt.FileSize))
			_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   &testReader{br: bytes.NewReader(expected)},
			})
			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}

			if v := atomic.LoadInt64(&p.recordedOutstanding); v != 0 {
				t.Fatalf("expected zero outsnatding pool parts, got %d", v)
			}

			gets, allocs := atomic.LoadUint64(&p.recordedGets), atomic.LoadUint64(&p.recordedAllocs)

			t.Logf("total gets %v, total allocations %v", gets, allocs)
			if e, a := tt.ExAllocations, allocs; a > e {
				t.Errorf("expected %v allocations, got %v", e, a)
			}
		})
	}
}

func TestUploadByteSlicePool_Failures(t *testing.T) {
	const (
		putObject               = "PutObject"
		createMultipartUpload   = "CreateMultipartUpload"
		uploadPart              = "UploadPart"
		completeMultipartUpload = "CompleteMultipartUpload"
	)

	cases := map[string]struct {
		PartSize   int64
		FileSize   int64
		Operations []string
	}{
		"single part": {
			PartSize: sdkio.MebiByte * 5,
			FileSize: sdkio.MebiByte * 4,
			Operations: []string{
				putObject,
			},
		},
		"multi-part": {
			PartSize: sdkio.MebiByte * 5,
			FileSize: sdkio.MebiByte * 10,
			Operations: []string{
				createMultipartUpload,
				uploadPart,
				completeMultipartUpload,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			for _, operation := range tt.Operations {
				t.Run(operation, func(t *testing.T) {
					var p *recordedPartPool

					unswap := swapByteSlicePool(func(sliceSize int64) byteSlicePool {
						p = newRecordedPartPool(sliceSize)
						return p
					})
					defer unswap()

					client, _, _ := s3testing.NewUploadLoggingClient(nil)

					switch operation {
					case putObject:
						client.PutObjectFn = func(*s3testing.UploadLoggingClient, *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
							return nil, fmt.Errorf("put object failure")
						}
					case createMultipartUpload:
						client.CreateMultipartUploadFn = func(*s3testing.UploadLoggingClient, *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
							return nil, fmt.Errorf("create multipart upload failure")
						}
					case uploadPart:
						client.UploadPartFn = func(*s3testing.UploadLoggingClient, *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
							return nil, fmt.Errorf("upload part failure")
						}
					case completeMultipartUpload:
						client.CompleteMultipartUploadFn = func(*s3testing.UploadLoggingClient, *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
							return nil, fmt.Errorf("complete multipart upload failure")
						}
					}

					uploader := NewUploader(client, func(u *Uploader) {
						u.Concurrency = 1
						u.PartSize = tt.PartSize
					})

					expected := s3testing.GetTestBytes(int(tt.FileSize))
					_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
						Bucket: aws.String("bucket"),
						Key:    aws.String("key"),
						Body:   &testReader{br: bytes.NewReader(expected)},
					})
					if err == nil {
						t.Fatalf("expected error but got none")
					}

					if v := atomic.LoadInt64(&p.recordedOutstanding); v != 0 {
						t.Fatalf("expected zero outsnatding pool parts, got %d", v)
					}
				})
			}
		})
	}
}

func TestUploadByteSlicePoolConcurrentMultiPartSize(t *testing.T) {
	var (
		pools []*recordedPartPool
		mtx   sync.Mutex
	)

	unswap := swapByteSlicePool(func(sliceSize int64) byteSlicePool {
		mtx.Lock()
		defer mtx.Unlock()
		b := newRecordedPartPool(sliceSize)
		pools = append(pools, b)
		return b
	})
	defer unswap()

	client, _, _ := s3testing.NewUploadLoggingClient(nil)

	uploader := NewUploader(client, func(u *Uploader) {
		u.PartSize = 5 * sdkio.MebiByte
		u.Concurrency = 2
	})

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			expected := s3testing.GetTestBytes(int(15 * sdkio.MebiByte))
			_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   &testReader{br: bytes.NewReader(expected)},
			})
			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
		}()
		go func() {
			defer wg.Done()
			expected := s3testing.GetTestBytes(int(15 * sdkio.MebiByte))
			_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   &testReader{br: bytes.NewReader(expected)},
			}, func(u *Uploader) {
				u.PartSize = 6 * sdkio.MebiByte
			})
			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
		}()
	}

	wg.Wait()

	if e, a := 3, len(pools); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	for _, p := range pools {
		if v := atomic.LoadInt64(&p.recordedOutstanding); v != 0 {
			t.Fatalf("expected zero outsnatding pool parts, got %d", v)
		}

		t.Logf("total gets %v, total allocations %v",
			atomic.LoadUint64(&p.recordedGets),
			atomic.LoadUint64(&p.recordedAllocs))
	}
}

func BenchmarkPools(b *testing.B) {
	cases := []struct {
		PartSize      int64
		FileSize      int64
		Concurrency   int
		ExAllocations uint64
	}{
		0: {
			PartSize:    sdkio.MebiByte * 5,
			FileSize:    sdkio.MebiByte * 5,
			Concurrency: 1,
		},
		1: {
			PartSize:    sdkio.MebiByte * 5,
			FileSize:    sdkio.MebiByte * 10,
			Concurrency: 1,
		},
		2: {
			PartSize:    sdkio.MebiByte * 5,
			FileSize:    sdkio.MebiByte * 20,
			Concurrency: 2,
		},
		3: {
			PartSize:    sdkio.MebiByte * 5,
			FileSize:    sdkio.MebiByte * 250,
			Concurrency: 10,
		},
	}

	client, _, _ := s3testing.NewUploadLoggingClient(nil)

	pools := map[string]func(sliceSize int64) byteSlicePool{
		"sync.Pool": func(sliceSize int64) byteSlicePool {
			return newSyncSlicePool(sliceSize)
		},
		"custom": func(sliceSize int64) byteSlicePool {
			return newMaxSlicePool(sliceSize)
		},
	}

	for name, poolFunc := range pools {
		b.Run(name, func(b *testing.B) {
			unswap := swapByteSlicePool(poolFunc)
			defer unswap()
			for i, c := range cases {
				b.Run(strconv.Itoa(i), func(b *testing.B) {
					uploader := NewUploader(client, func(u *Uploader) {
						u.PartSize = c.PartSize
						u.Concurrency = c.Concurrency
					})

					expected := s3testing.GetTestBytes(int(c.FileSize))
					b.ResetTimer()
					_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
						Bucket: aws.String("bucket"),
						Key:    aws.String("key"),
						Body:   &testReader{br: bytes.NewReader(expected)},
					})
					if err != nil {
						b.Fatalf("expected no error, but got %v", err)
					}
				})
			}
		})
	}
}
