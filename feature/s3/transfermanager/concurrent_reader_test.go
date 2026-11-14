package transfermanager

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// TestConcurrentReaderReadUsesSliceLenForBounds is a regression test for
// #3388. Read must not write past len(p) when cap(p) is larger.
func TestConcurrentReaderReadUsesSliceLenForBounds(t *testing.T) {
	for _, mode := range []types.GetObjectType{types.GetObjectParts, types.GetObjectRanges} {
		t.Run(string(mode), func(t *testing.T) {
			const partSize = 8
			parts := [][]byte{
				bytes.Repeat([]byte{'a'}, partSize),
				bytes.Repeat([]byte{'b'}, partSize),
			}
			var want []byte
			for _, b := range parts {
				want = append(want, b...)
			}

			s3Client := &s3testing.TransferManagerLoggingClient{Data: want, PartsData: parts, PartsCount: int32(len(parts))}
			if mode == types.GetObjectParts {
				s3Client.Data = parts[0]
				s3Client.GetObjectFn = s3testing.UnequalPartGetObjectFn
			} else {
				s3Client.GetObjectFn = s3testing.RangeGetObjectFn
			}

			mgr := New(s3Client, func(o *Options) {
				o.GetObjectType = mode
				o.Concurrency = 2
				o.PartSizeBytes = partSize
				o.GetObjectBufferSize = 2 * partSize
			})
			out, err := mgr.GetObject(context.Background(), &GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			})
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			var got []byte
			for {
				p := make([]byte, 4, 9)
				n, err := out.Body.Read(p)
				if n > len(p) {
					t.Fatalf("read %d bytes into a slice of len %d", n, len(p))
				}
				for _, b := range p[n:] {
					if b != 0 {
						t.Fatal("read wrote past len(p) into the slice's spare capacity")
					}
				}
				got = append(got, p[:n]...)
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("expect no read error, got %v", err)
				}
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("expect %q, got %q", want, got)
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
				partSize:   c.partSize,
				partsCount: c.partsCount,
				options:    c.options,
				in: &GetObjectInput{
					Bucket: aws.String("bucket"),
					Key:    aws.String("key"),
				},
			}
			// Set the window the way get() does, from a buffer of sectionParts parts.
			r.window = min(getWindow(c.options.Concurrency, int64(c.sectionParts)*c.partSize, c.partSize), c.partsCount)

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
			r.initReader(ctx)

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

// TestConcurrentReaderReadRepeatAfterError verifies that once a download fails,
// every Read returns the same error and no more GetObject calls are made.
func TestConcurrentReaderReadRepeatAfterError(t *testing.T) {
	ctx := context.Background()
	s3Client := &s3testing.TransferManagerLoggingClient{}
	s3Client.GetObjectFn = s3testing.ErrRangeGetObjectFn
	s3Client.Data = []byte("abcdefghijkl")

	r := &concurrentReader{
		partSize:   4,
		partsCount: 3,
		options: Options{
			GetObjectType: types.GetObjectRanges,
			Concurrency:   1,
			S3:            s3Client,
		},
		in: &GetObjectInput{
			Bucket: aws.String("bucket"),
			Key:    aws.String("key"),
		},
		totalBytes: int64(len(s3Client.Data)),
	}
	r.window = 1
	r.initReader(ctx)

	buf := make([]byte, 4)
	var readErr error
	for i := 0; i < 100 && readErr == nil; i++ {
		_, readErr = r.Read(buf)
	}
	if readErr == nil {
		t.Fatal("expected read to eventually return an error")
	}
	if readErr == io.EOF {
		t.Fatal("expected a service error, got io.EOF")
	}

	stored := r.loadErr()
	if stored == nil {
		t.Fatal("expected the reader to record the error")
	}
	if !errors.Is(readErr, stored) {
		t.Fatalf("expected read to return stored error, got %v and stored %v", readErr, stored)
	}

	invocationsAfterError := s3Client.GetObjectInvocations

	_, err := r.Read(buf)
	if !errors.Is(err, stored) {
		t.Fatalf("expected repeated read to return stored error, got %v and stored %v", err, stored)
	}
	if got := s3Client.GetObjectInvocations; got != invocationsAfterError {
		t.Fatalf("expected no new GetObject calls after error, got %d vs %d", got, invocationsAfterError)
	}
}

// TestConcurrentReaderInFlightWindow verifies the parts in flight reach, and
// never exceed, the smaller of Concurrency and GetObjectBufferSize/PartSizeBytes.
func TestConcurrentReaderInFlightWindow(t *testing.T) {
	cases := map[string]struct {
		concurrency int
		bufferParts int64
		partSize    int64
		parts       int64
		wantWindow  int32
	}{
		"limited by concurrency": {concurrency: 8, bufferParts: 32, partSize: 16, parts: 64, wantWindow: 8},
		"limited by buffer":      {concurrency: 32, bufferParts: 8, partSize: 16, parts: 64, wantWindow: 8},
		"wide window":            {concurrency: 16, bufferParts: 16, partSize: 16, parts: 64, wantWindow: 16},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			data := make([]byte, c.partSize*c.parts)
			for i := range data {
				data[i] = byte(i % 251)
			}
			s3Client := &s3testing.TransferManagerLoggingClient{Data: data}

			release := make(chan struct{})
			s3Client.GetObjectFn = blockingRangeGetObjectFn(release)

			mgr := New(s3Client, func(o *Options) {
				o.GetObjectType = types.GetObjectRanges
				o.Concurrency = c.concurrency
				o.PartSizeBytes = c.partSize
				o.GetObjectBufferSize = c.bufferParts * c.partSize
			})

			// Release the requests once the window is full, or after a
			// timeout so a regression fails instead of hanging.
			go func() {
				deadline := time.Now().Add(3 * time.Second)
				for time.Now().Before(deadline) {
					if s3Client.MaxInFlight.Load() >= c.wantWindow {
						break
					}
					time.Sleep(time.Millisecond)
				}
				close(release)
			}()

			out, err := mgr.GetObject(context.Background(), &GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			got, err := io.ReadAll(out.Body)
			if err != nil {
				t.Fatalf("expected no read error, got %v", err)
			}
			if !bytes.Equal(got, data) {
				t.Errorf("downloaded data mismatch: got %d bytes, want %d", len(got), len(data))
			}
			if a := s3Client.MaxInFlight.Load(); a != c.wantWindow {
				t.Errorf("expected max in-flight %d, got %d", c.wantWindow, a)
			}
		})
	}
}

// TestConcurrentReaderRetriesBodyError verifies a failed body read is retried
// and the download then succeeds.
func TestConcurrentReaderRetriesBodyError(t *testing.T) {
	data := []byte("hello concurrent world")
	var calls int32
	const failUntil = 2

	s3Client := &s3testing.TransferManagerLoggingClient{Data: data}
	s3Client.GetObjectFn = func(cl *s3testing.TransferManagerLoggingClient, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		n := atomic.AddInt32(&calls, 1)
		body := cl.Data
		if n <= failUntil {
			// A body shorter than ContentLength fails the read.
			body = cl.Data[:len(cl.Data)-1]
		}
		return &s3.GetObjectOutput{
			Body:          io.NopCloser(bytes.NewReader(body)),
			ContentLength: aws.Int64(int64(len(cl.Data))),
			ETag:          aws.String("myetag"),
		}, nil
	}

	mgr := New(s3Client, func(o *Options) {
		o.GetObjectType = types.GetObjectRanges
		o.Concurrency = 1
		o.PartBodyMaxRetries = 5
	})

	out, err := mgr.GetObject(context.Background(), &GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := io.ReadAll(out.Body)
	if err != nil {
		t.Fatalf("expected no read error, got %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Errorf("expected %q, got %q", data, got)
	}
	if a := atomic.LoadInt32(&calls); a != failUntil+1 {
		t.Errorf("expected %d GetObject attempts, got %d", failUntil+1, a)
	}
}

// TestConcurrentReaderClose verifies Close stops a download in progress, can
// be called twice, and makes Read fail.
func TestConcurrentReaderClose(t *testing.T) {
	partSize := int64(16)
	parts := int64(64)
	data := make([]byte, partSize*parts)
	for i := range data {
		data[i] = byte(i % 251)
	}

	s3Client := &s3testing.TransferManagerLoggingClient{Data: data}
	release := make(chan struct{})
	s3Client.GetObjectFn = blockingRangeGetObjectFn(release)

	const concurrency = 8
	mgr := New(s3Client, func(o *Options) {
		o.GetObjectType = types.GetObjectRanges
		o.Concurrency = concurrency
		o.PartSizeBytes = partSize
		o.GetObjectBufferSize = concurrency * partSize
	})

	out, err := mgr.GetObject(context.Background(), &GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	closer, ok := out.Body.(io.Closer)
	if !ok {
		t.Fatal("expected GetObject body to implement io.Closer")
	}

	readErrCh := make(chan error, 1)
	go func() {
		_, err := io.ReadAll(out.Body)
		readErrCh <- err
	}()

	// Wait until the window is full of blocked requests.
	waitFor(t, 2*time.Second, func() bool {
		return s3Client.InFlight.Load() >= concurrency
	})

	if err := closer.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if err := closer.Close(); err != nil {
		t.Fatalf("second Close returned error: %v", err)
	}

	// Unblock the requests. The workers should discard them and exit.
	close(release)

	waitFor(t, 2*time.Second, func() bool {
		return s3Client.InFlight.Load() == 0
	})

	select {
	case err := <-readErrCh:
		if !errors.Is(err, errReadClosed) {
			t.Fatalf("expect read to fail with %q after Close, got %v", errReadClosed, err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("reader did not return after Close")
	}
}

// TestConcurrentReaderContextCancel verifies canceling the context stops a
// download in progress without Close.
func TestConcurrentReaderContextCancel(t *testing.T) {
	partSize := int64(16)
	parts := int64(64)
	data := make([]byte, partSize*parts)
	for i := range data {
		data[i] = byte(i % 251)
	}

	s3Client := &s3testing.TransferManagerLoggingClient{Data: data}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s3Client.GetObjectFn = func(cl *s3testing.TransferManagerLoggingClient, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
			return nil, fmt.Errorf("request not canceled in time")
		}
	}

	const concurrency = 8
	mgr := New(s3Client, func(o *Options) {
		o.GetObjectType = types.GetObjectRanges
		o.Concurrency = concurrency
		o.PartSizeBytes = partSize
		o.GetObjectBufferSize = concurrency * partSize
	})

	out, err := mgr.GetObject(ctx, &GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	readErrCh := make(chan error, 1)
	go func() {
		_, err := io.ReadAll(out.Body)
		readErrCh <- err
	}()

	waitFor(t, 2*time.Second, func() bool {
		return s3Client.InFlight.Load() >= concurrency
	})

	cancel()

	select {
	case err := <-readErrCh:
		if err == nil {
			t.Fatal("expected a non-nil error after context cancellation")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("reader did not return after context cancellation")
	}

	waitFor(t, 2*time.Second, func() bool {
		return s3Client.InFlight.Load() == 0
	})
}

// blockingRangeGetObjectFn serves byte ranges of Data once release is closed.
func blockingRangeGetObjectFn(release <-chan struct{}) func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return func(cl *s3testing.TransferManagerLoggingClient, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		<-release
		start, fin, err := getReqRange(aws.ToString(params.Range))
		if err != nil {
			return nil, err
		}
		fin++
		if fin > int64(len(cl.Data)) {
			fin = int64(len(cl.Data))
		}
		body := cl.Data[start:fin]
		out := &s3.GetObjectOutput{
			Body:          io.NopCloser(bytes.NewReader(body)),
			ContentLength: aws.Int64(int64(len(body))),
			ETag:          aws.String("myetag"),
		}
		if len(body) != len(cl.Data) {
			out.ContentRange = aws.String(fmt.Sprintf("bytes %d-%d/%d", start, fin-1, len(cl.Data)))
		}
		return out, nil
	}
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("condition not met within %s", timeout)
}

// partsGetObject starts a GetObjectParts download of partsData, served by fn.
func partsGetObject(t *testing.T, partsData [][]byte, concurrency int, bufferSize int64,
	fn func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error),
) (io.Reader, []byte) {
	t.Helper()
	var expect []byte
	for _, b := range partsData {
		expect = append(expect, b...)
	}
	s3Client := &s3testing.TransferManagerLoggingClient{}
	s3Client.GetObjectFn = fn
	s3Client.PartsData = partsData
	s3Client.PartsCount = int32(len(partsData))
	s3Client.Data = expect

	mgr := New(s3Client, func(o *Options) {
		o.GetObjectType = types.GetObjectParts
		o.Concurrency = concurrency
		o.GetObjectBufferSize = bufferSize
	})
	out, err := mgr.GetObject(context.Background(), &GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	return out.Body, expect
}

func sizedParts(sizes []int) [][]byte {
	parts := make([][]byte, len(sizes))
	for i, s := range sizes {
		parts[i] = bytes.Repeat([]byte{byte('A' + i%26)}, s)
	}
	return parts
}

// TestConcurrentReaderPartUnequalSizes verifies parts-mode downloads of objects
// with unequal part sizes (#3526).
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
			body, expect := partsGetObject(t, sizedParts(c.sizes), c.concurrency, c.bufferThreshold, s3testing.UnequalPartGetObjectFn)
			got, err := io.ReadAll(body)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}
			if !bytes.Equal(expect, got) {
				t.Fatalf("expect %d bytes equal to assembled parts, got %d", len(expect), len(got))
			}
		})
	}
}

// TestConcurrentReaderPartMemoryThrottleNoDeadlock verifies a parts-mode
// download does not hang when a small, delayed part 1 arrives after larger
// parts (#3526 follow-up).
func TestConcurrentReaderPartMemoryThrottleNoDeadlock(t *testing.T) {
	fn := func(c *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		if aws.ToInt32(in.PartNumber) == 1 {
			time.Sleep(200 * time.Millisecond)
		}
		return s3testing.UnequalPartGetObjectFn(c, in)
	}
	body, expect := partsGetObject(t, sizedParts([]int{10, 200, 200, 200, 200, 200}), 2, 100, fn)

	done := make(chan struct{})
	var got []byte
	var err error
	go func() {
		got, err = io.ReadAll(body)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("Read did not complete: likely deadlocked with an orphaned download producer")
	}
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if !bytes.Equal(expect, got) {
		t.Fatalf("expect %d bytes equal to assembled parts, got %d", len(expect), len(got))
	}
}

// TestConcurrentReaderSmallReadsNoDeadlock verifies small, part-aligned reads
// do not deadlock a parts-mode download (#3552).
func TestConcurrentReaderSmallReadsNoDeadlock(t *testing.T) {
	const (
		partSize    = 1 << 10
		concurrency = 5
		partsCount  = 6 + concurrency + 4
		readSize    = 128
	)
	sizes := make([]int, partsCount)
	for i := range sizes {
		sizes[i] = partSize
	}
	body, expect := partsGetObject(t, sizedParts(sizes), concurrency, 6*partSize, s3testing.ReaderPartGetObjectFn)

	done := make(chan error, 1)
	var got []byte
	go func() {
		br := bufio.NewReaderSize(body, readSize)
		buf := make([]byte, readSize)
		for {
			n, err := br.Read(buf)
			got = append(got, buf[:n]...)
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				done <- err
				return
			}
		}
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("read error: %v", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("DEADLOCK: Read did not complete within 15s (GH #3552)")
	}
	if !bytes.Equal(expect, got) {
		t.Fatalf("expect %d bytes equal to assembled parts, got %d", len(expect), len(got))
	}
}
