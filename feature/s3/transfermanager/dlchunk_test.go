package transfermanager

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	internalio "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/io"
)

const chunkTestBufSize = 16

type countingPool struct {
	t *testing.T

	mu          sync.Mutex
	outstanding map[*byte]struct{}
}

func (p *countingPool) Get() any {
	p.mu.Lock()
	defer p.mu.Unlock()

	b := make([]byte, chunkTestBufSize)
	p.outstanding[&b[0]] = struct{}{}
	return b
}

func (p *countingPool) Put(v any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	b := v.([]byte)
	if len(b) != chunkTestBufSize {
		p.t.Errorf("Put buffer with len %d, want %d", len(b), chunkTestBufSize)
		return
	}
	if _, ok := p.outstanding[&b[0]]; !ok {
		p.t.Errorf("Put buffer %p that is not outstanding", &b[0])
		return
	}
	delete(p.outstanding, &b[0])
}

type chunkSink struct {
	mu  sync.Mutex
	buf []byte
}

func (s *chunkSink) WriteAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy(s.buf[off:], p)
	return len(p), nil
}

type chunkHarness struct {
	pool *countingPool
	sink *chunkSink
	w    *internalio.AsyncWriterAt
}

func newChunkHarness(t *testing.T, size int) *chunkHarness {
	h := &chunkHarness{
		pool: &countingPool{t: t, outstanding: map[*byte]struct{}{}},
		sink: &chunkSink{buf: make([]byte, size)},
	}
	h.w = internalio.NewAsyncWriterAt(h.sink, h.pool)
	h.w.Start()
	return h
}

func (h *chunkHarness) finish(t *testing.T) {
	t.Helper()

	h.w.Wait()
	h.w.Stop()
	if err := h.w.Error(); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if n := len(h.pool.outstanding); n != 0 {
		t.Errorf("%d buffers were never returned to the pool", n)
	}
}

var chunkReaders = map[string]func(io.Reader) io.Reader{
	"plain":         func(r io.Reader) io.Reader { return r },
	"one byte":      iotest.OneByteReader,
	"half":          iotest.HalfReader,
	"data with err": iotest.DataErrReader,
}

var chunkSizes = []int{
	0,
	1,
	chunkTestBufSize - 1,
	chunkTestBufSize,
	chunkTestBufSize + 1,
	3 * chunkTestBufSize,
	3*chunkTestBufSize + 5,
}

func TestDlChunkReadFrom(t *testing.T) {
	for rname, wrap := range chunkReaders {
		for _, size := range chunkSizes {
			for _, start := range []int{0, 1000} {
				t.Run(fmt.Sprintf("%s/size=%d/start=%d", rname, size, start), func(t *testing.T) {
					data := randomBytes(int64(size), size)
					h := newChunkHarness(t, start+size+chunkTestBufSize)

					c := dlChunk{start: int64(start), sink: h.w}
					n, err := c.ReadFrom(wrap(bytes.NewReader(data)))
					h.finish(t)

					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					if n != int64(size) {
						t.Errorf("ReadFrom returned %d, want %d", n, size)
					}
					checkBytes(t, h.sink.buf[start:start+size], data)
					if !isZero(h.sink.buf[:start]) || !isZero(h.sink.buf[start+size:]) {
						t.Error("ReadFrom wrote outside [start, start+size)")
					}
				})
			}
		}
	}
}

func TestDlChunkReadFromError(t *testing.T) {
	errBody := errors.New("body failed")
	for rname, wrap := range chunkReaders {
		for _, size := range chunkSizes {
			t.Run(fmt.Sprintf("%s/failAt=%d", rname, size), func(t *testing.T) {
				const start = 100
				data := randomBytes(int64(size), size)
				h := newChunkHarness(t, start+size+chunkTestBufSize)

				c := dlChunk{start: start, sink: h.w}
				body := io.MultiReader(bytes.NewReader(data), iotest.ErrReader(errBody))
				n, err := c.ReadFrom(wrap(body))
				h.finish(t)

				if err != errBody {
					t.Fatalf("err = %v, want %v", err, errBody)
				}
				if n > int64(size) || n%chunkTestBufSize != 0 {
					t.Fatalf("ReadFrom returned %d, want a multiple of %d <= %d", n, chunkTestBufSize, size)
				}
				checkBytes(t, h.sink.buf[start:start+int(n)], data[:n])
				if !isZero(h.sink.buf[start+int(n):]) {
					t.Error("ReadFrom wrote a partial buffer")
				}
			})
		}
	}
}

// A body that ends with io.ErrUnexpectedEOF is truncated, even if it happens
// to end on a buffer boundary, and the tail buffer is never written.
func TestDlChunkReadFromUnexpectedEOF(t *testing.T) {
	const start = 100

	for _, size := range []int{chunkTestBufSize - 1, 3 * chunkTestBufSize, 3*chunkTestBufSize + 5} {
		t.Run(fmt.Sprintf("size=%d", size), func(t *testing.T) {
			data := randomBytes(int64(size), size)
			h := newChunkHarness(t, start+size+chunkTestBufSize)

			c := dlChunk{start: start, sink: h.w}
			body := io.MultiReader(bytes.NewReader(data), iotest.ErrReader(io.ErrUnexpectedEOF))
			n, err := c.ReadFrom(body)
			h.finish(t)

			if err != io.ErrUnexpectedEOF {
				t.Fatalf("err = %v, want %v", err, io.ErrUnexpectedEOF)
			}
			wantRead := int64(size / chunkTestBufSize * chunkTestBufSize)
			if n != wantRead {
				t.Fatalf("ReadFrom returned %d, want %d", n, wantRead)
			}
			checkBytes(t, h.sink.buf[start:start+int(n)], data[:n])
			if !isZero(h.sink.buf[start+int(n):]) {
				t.Error("ReadFrom wrote past the bytes it reported")
			}
		})
	}
}

// Retries reuse the chunk, so a second ReadFrom must write from start again.
func TestDlChunkReadFromRetry(t *testing.T) {
	const (
		start = 100
		size  = 3*chunkTestBufSize + 5
	)

	for _, failAt := range []int{0, 1, chunkTestBufSize, 2*chunkTestBufSize + 3} {
		t.Run(fmt.Sprintf("failAt=%d", failAt), func(t *testing.T) {
			data := randomBytes(size, size)
			h := newChunkHarness(t, start+size)

			c := dlChunk{start: start, sink: h.w}
			body := io.MultiReader(bytes.NewReader(data[:failAt]), iotest.ErrReader(errors.New("body failed")))
			if _, err := c.ReadFrom(body); err == nil {
				t.Fatal("expected first attempt to fail")
			}

			n, err := c.ReadFrom(bytes.NewReader(data))
			h.finish(t)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n != size {
				t.Errorf("ReadFrom returned %d, want %d", n, size)
			}
			checkBytes(t, h.sink.buf[start:], data)
		})
	}
}

type blockingSink struct {
	entered chan struct{}
	release chan struct{}
}

func (s *blockingSink) WriteAt(p []byte, _ int64) (int, error) {
	select {
	case s.entered <- struct{}{}:
	default:
	}
	<-s.release
	return len(p), nil
}

// ReadFrom must not report a failure until every write it queued has landed,
// otherwise a retry's writes could race with them.
func TestDlChunkReadFromErrorWaitsForWrites(t *testing.T) {
	pool := &countingPool{t: t, outstanding: map[*byte]struct{}{}}
	sink := &blockingSink{entered: make(chan struct{}, 1), release: make(chan struct{})}
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(sink.release) }) }
	w := internalio.NewAsyncWriterAt(sink, pool)
	w.Start()
	defer w.Stop()
	defer release()

	errBody := errors.New("body failed")
	c := dlChunk{sink: w}
	body := io.MultiReader(bytes.NewReader(make([]byte, 2*chunkTestBufSize)), iotest.ErrReader(errBody))

	returned := make(chan error, 1)
	go func() {
		_, err := c.ReadFrom(body)
		returned <- err
	}()

	<-sink.entered
	select {
	case err := <-returned:
		t.Fatalf("ReadFrom returned %v while its writes were still in progress", err)
	case <-time.After(20 * time.Millisecond):
	}

	release()
	select {
	case err := <-returned:
		if err != errBody {
			t.Fatalf("err = %v, want %v", err, errBody)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("ReadFrom did not return after its writes were released")
	}
}

func isZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}
