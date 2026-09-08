package transfermanager

import (
	"bytes"
	"context"
	"errors"

	"io"
	"math"
	"math/rand"
	"testing"
	"time"

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

// TestConcurrentReaderPartUnequalSizes exercises the parts-mode read path
// (partRead) with multipart objects whose parts have unequal sizes (#3526).
// The existing TestConcurrentReader cases don't set getType/bufferThreshold, so
// they route through rangeRead; these cases set both to drive partRead directly,
// covering large and small memory thresholds and single/multi goroutine paths.
func TestConcurrentReaderPartUnequalSizes(t *testing.T) {
	cases := map[string]struct {
		sizes           []int
		concurrency     int
		bufferThreshold int64
	}{
		"conc1 large threshold":  {sizes: []int{60, 50, 10}, concurrency: 1, bufferThreshold: 1 << 20},
		"conc5 large threshold":  {sizes: []int{60, 50, 10}, concurrency: 5, bufferThreshold: 1 << 20},
		"conc5 small threshold":  {sizes: []int{60, 50, 49, 2}, concurrency: 5, bufferThreshold: 50},
		"conc5 first part small": {sizes: []int{8, 50, 49, 2}, concurrency: 5, bufferThreshold: 50},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			partsData := make([][]byte, len(c.sizes))
			var expect []byte
			for i, s := range c.sizes {
				b := bytes.Repeat([]byte{byte('A' + i)}, s)
				partsData[i] = b
				expect = append(expect, b...)
			}

			s3Client := &s3testing.TransferManagerLoggingClient{}
			s3Client.GetObjectFn = s3testing.UnequalPartGetObjectFn
			s3Client.PartsData = partsData
			s3Client.PartsCount = int32(len(c.sizes))
			s3Client.Data = expect

			partSize := int64(c.sizes[0])
			sectionParts := int32(c.bufferThreshold / partSize)
			if sectionParts < 1 {
				sectionParts = 1
			}
			partsCount := int32(len(c.sizes))
			capacity := sectionParts
			if capacity > partsCount {
				capacity = partsCount
			}

			r := &concurrentReader{
				partSize:        partSize,
				partsCount:      partsCount,
				sectionParts:    sectionParts,
				getType:         types.GetObjectParts,
				bufferThreshold: c.bufferThreshold,
				options: Options{
					GetObjectType: types.GetObjectParts,
					Concurrency:   c.concurrency,
					S3:            s3Client,
				},
				in:         &GetObjectInput{Bucket: aws.String("bucket"), Key: aws.String("key")},
				capacity:   capacity,
				buf:        make(map[int32]*outChunk),
				ctx:        context.Background(),
				ch:         make(chan outChunk, c.concurrency),
				totalBytes: int64(len(expect)),
			}

			got, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}
			if e, a := len(expect), len(got); e != a {
				t.Fatalf("expect %d bytes, got %d", e, a)
			}
			if !bytes.Equal(expect, got) {
				t.Fatalf("expect downloaded stream to equal assembled parts")
			}
		})
	}
}

// TestConcurrentReaderPartMemoryThrottleNoDeadlock guards against the parts-mode
// deadlock where bounding memory by breaking out of the receive loop orphaned an
// in-flight download producer on a full r.ch and hung Read's r.wg.Wait() (#3526
// follow-up). It reproduces the trigger conditions: sectionParts > Concurrency
// (more parts dispatched than r.ch can buffer), unequal parts where later parts
// are much larger than part 1, and a delayed part 0 so a large later part is
// received first. Memory must be bounded by throttling dispatch, not by
// abandoning received parts, so this must complete rather than hang.
func TestConcurrentReaderPartMemoryThrottleNoDeadlock(t *testing.T) {
	sizes := []int{10, 200, 200, 200, 200, 200}
	partsData := make([][]byte, len(sizes))
	var expect []byte
	for i, s := range sizes {
		b := bytes.Repeat([]byte{byte('A' + i)}, s)
		partsData[i] = b
		expect = append(expect, b...)
	}

	s3Client := &s3testing.TransferManagerLoggingClient{}
	s3Client.PartsData = partsData
	s3Client.PartsCount = int32(len(sizes))
	s3Client.Data = expect
	// Delay part 0 so a later, larger part is received first and would trip the
	// memory budget before the consecutive part arrives.
	s3Client.GetObjectFn = func(c *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		if aws.ToInt32(in.PartNumber) == 1 {
			time.Sleep(200 * time.Millisecond)
		}
		return s3testing.UnequalPartGetObjectFn(c, in)
	}

	partSize := int64(sizes[0]) // 10
	bufferThreshold := int64(100)
	sectionParts := int32(bufferThreshold / partSize) // 10, > Concurrency below
	partsCount := int32(len(sizes))
	capacity := sectionParts
	if capacity > partsCount {
		capacity = partsCount
	}

	r := &concurrentReader{
		partSize:        partSize,
		partsCount:      partsCount,
		sectionParts:    sectionParts,
		getType:         types.GetObjectParts,
		bufferThreshold: bufferThreshold,
		options: Options{
			GetObjectType: types.GetObjectParts,
			Concurrency:   2, // < sectionParts and < in-flight parts
			S3:            s3Client,
		},
		in:         &GetObjectInput{Bucket: aws.String("bucket"), Key: aws.String("key")},
		capacity:   capacity,
		buf:        make(map[int32]*outChunk),
		ctx:        context.Background(),
		ch:         make(chan outChunk, 2),
		totalBytes: int64(len(expect)),
	}

	done := make(chan struct{})
	var got []byte
	var err error
	go func() {
		got, err = io.ReadAll(r)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("Read did not complete: likely deadlocked on r.wg.Wait() with an orphaned download producer")
	}

	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if !bytes.Equal(expect, got) {
		t.Fatalf("expect %d bytes equal to assembled parts, got %d", len(expect), len(got))
	}
}
