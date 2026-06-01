package transfermanager

import (
	"bytes"
	"context"
	"errors"

	"io"
	"math"
	"math/rand"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestConcurrentReaderReadUsesSliceLenForBounds(t *testing.T) {
	tests := map[string]func(*concurrentReader){
		"buffered chunk": func(r *concurrentReader) {
			r.capacity = 2
			r.receiveCount = 2
			r.buf[1] = &outChunk{
				body:   bytes.NewReader([]byte("chunk")),
				index:  1,
				length: int64(len("chunk")),
			}
		},
		"received chunk": func(r *concurrentReader) {
			r.capacity = 2
			r.receiveCount = 1
			r.ch <- outChunk{
				body:   bytes.NewReader([]byte("chunk")),
				index:  1,
				length: int64(len("chunk")),
			}
		},
	}

	for name, setup := range tests {
		t.Run(name, func(t *testing.T) {
			r := &concurrentReader{
				partSize:   8,
				partsCount: 2,
				buf:        make(map[int32]*outChunk),
				ch:         make(chan outChunk, 1),
			}
			setup(r)

			p := make([]byte, 4, 9)
			n, err := r.read(p)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if n != 0 {
				t.Fatalf("expect no bytes read into short slice, got %d", n)
			}

			chunk, ok := r.buf[1]
			if !ok {
				t.Fatal("expect chunk to remain buffered")
			}
			if chunk.cur != 0 {
				t.Fatalf("expect chunk cursor to remain unchanged, got %d", chunk.cur)
			}
		})
	}
}

func TestConcurrentReader(t *testing.T) {
	cases := map[string]struct {
		partSize     int64
		partsCount   int32
		sectionParts int32
		getObjectFn  func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error)
		options      Options
	}{
		"part get single goroutine": {
			partSize:     10,
			partsCount:   1000,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectParts,
				Concurrency:   1,
			},
			getObjectFn: s3testing.ReaderPartGetObjectFn,
		},
		"part get single goroutine with only one section": {
			partSize:     1000,
			partsCount:   5,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectParts,
				Concurrency:   3,
			},
			getObjectFn: s3testing.ReaderPartGetObjectFn,
		},
		"part get single goroutine with only one part": {
			partSize:     1000,
			partsCount:   1,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectParts,
				Concurrency:   3,
			},
			getObjectFn: s3testing.ReaderPartGetObjectFn,
		},
		"part get multiple goroutines": {
			partSize:     10,
			partsCount:   1000,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectParts,
				Concurrency:   5,
			},
			getObjectFn: s3testing.ReaderPartGetObjectFn,
		},
		"part get multiple goroutines with only one section": {
			partSize:     10,
			partsCount:   6,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectParts,
				Concurrency:   5,
			},
			getObjectFn: s3testing.ReaderPartGetObjectFn,
		},
		"part get multiple goroutines with only one part": {
			partSize:     10,
			partsCount:   1,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectParts,
				Concurrency:   5,
			},
			getObjectFn: s3testing.ReaderPartGetObjectFn,
		},
		"part get multiple goroutines with large part size": {
			partSize:     10000,
			partsCount:   10000,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectParts,
				Concurrency:   5,
			},
			getObjectFn: s3testing.ReaderPartGetObjectFn,
		},
		"range get single goroutine": {
			partSize:     10,
			partsCount:   1000,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectRanges,
				Concurrency:   1,
			},
			getObjectFn: s3testing.RangeGetObjectFn,
		},
		"range get single goroutine with only one section": {
			partSize:     1000,
			partsCount:   5,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectRanges,
				Concurrency:   3,
			},
			getObjectFn: s3testing.RangeGetObjectFn,
		},
		"range get single goroutine with only one part": {
			partSize:     1000,
			partsCount:   1,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectRanges,
				Concurrency:   3,
			},
			getObjectFn: s3testing.RangeGetObjectFn,
		},
		"range get multiple goroutines": {
			partSize:     10,
			partsCount:   1000,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectRanges,
				Concurrency:   5,
			},
			getObjectFn: s3testing.RangeGetObjectFn,
		},
		"range get multiple goroutines with only one section": {
			partSize:     10,
			partsCount:   6,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectRanges,
				Concurrency:   5,
			},
			getObjectFn: s3testing.RangeGetObjectFn,
		},
		"range get multiple goroutines with only one part": {
			partSize:     10,
			partsCount:   1,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectRanges,
				Concurrency:   5,
			},
			getObjectFn: s3testing.RangeGetObjectFn,
		},
		"range get multiple goroutines with large part size": {
			partSize:     10000,
			partsCount:   10000,
			sectionParts: 6,
			options: Options{
				GetObjectType: types.GetObjectRanges,
				Concurrency:   5,
			},
			getObjectFn: s3testing.RangeGetObjectFn,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			s3Client := &s3testing.TransferManagerLoggingClient{}
			s3Client.GetObjectFn = c.getObjectFn
			r := &concurrentReader{
				partSize:     c.partSize,
				partsCount:   c.partsCount,
				sectionParts: c.sectionParts,
				options:      c.options,
				in: &GetObjectInput{
					Bucket: aws.String("bucket"),
					Key:    aws.String("key"),
				},
				capacity: int32(math.Min(float64(c.sectionParts), float64(c.partsCount))),
				buf:      make(map[int32]*outChunk),
				ctx:      ctx,
				ch:       make(chan outChunk, c.options.Concurrency),
			}

			expectBuf := make([]byte, 0)
			expectPartsData := make([][]byte, c.partsCount)
			for i := int32(0); i < c.partsCount; i++ {
				b := make([]byte, c.partSize)
				if i == c.partsCount-1 {
					b = make([]byte, rand.Intn(int(c.partSize))+1)
				}
				rand.Read(b)
				expectBuf = append(expectBuf, b...)
				expectPartsData[i] = b
			}
			s3Client.Data = expectBuf
			s3Client.PartsData = expectPartsData
			r.options.S3 = s3Client
			r.totalBytes = int64(len(expectBuf))

			actualBuf, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("expect no error when reading, got %v", err)
			}

			if e, a := len(expectBuf), len(actualBuf); e != a {
				t.Errorf("expect data sent to have length %d, but got %d", e, a)
			}
			if e, a := expectBuf, actualBuf; !bytes.Equal(e, a) {
				t.Errorf("expect data sent to be %v, got %v", e, a)
			}
		})
	}
}

func TestConcurrentReaderReadRepeatAfterError(t *testing.T) {
	ctx := context.Background()
	s3Client := &s3testing.TransferManagerLoggingClient{}
	s3Client.GetObjectFn = s3testing.ErrRangeGetObjectFn
	s3Client.Data = []byte("abcdefghijkl")

	r := &concurrentReader{
		partSize:     4,
		partsCount:   3,
		sectionParts: 3,
		options: Options{
			GetObjectType: types.GetObjectRanges,
			Concurrency:   1,
			S3:            s3Client,
		},
		in: &GetObjectInput{
			Bucket: aws.String("bucket"),
			Key:    aws.String("key"),
		},
		capacity:   3,
		buf:        make(map[int32]*outChunk),
		ctx:        ctx,
		ch:         make(chan outChunk, 1),
		totalBytes: int64(len(s3Client.Data)),
	}

	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err == nil {
		t.Fatal("expected first read to return an error")
	}
	if !errors.Is(err, r.getErr()) {
		t.Fatalf("expected first read to return stored error, got %v and stored %v", err, r.getErr())
	}

	firstReadInvocations := s3Client.GetObjectInvocations
	if firstReadInvocations != 2 {
		t.Fatalf("expected first read to schedule 2 GetObject calls, got %d", firstReadInvocations)
	}

	_, err = r.Read(buf)
	if !errors.Is(err, r.getErr()) {
		t.Fatalf("expected repeated read to return stored error, got %v and stored %v", err, r.getErr())
	}

	if got := s3Client.GetObjectInvocations; got != firstReadInvocations {
		t.Fatalf("expected repeated read not to schedule more downloads, got %d GetObject calls after %d on first read", got, firstReadInvocations)
	}
}
