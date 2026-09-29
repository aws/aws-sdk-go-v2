package io

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	testTimeout = 10 * time.Second
	poison      = 0xde
)

var errInjected = errors.New("injected write failure")

// testPool reuses buffers aggressively and poisons them on Put, so a buffer
// released while its write is still in progress shows up as corrupted output.
type testPool struct {
	t    *testing.T
	size int

	mu          sync.Mutex
	free        [][]byte
	outstanding map[*byte]struct{}
}

func newTestPool(t *testing.T, size int) *testPool {
	return &testPool{t: t, size: size, outstanding: map[*byte]struct{}{}}
}

func (p *testPool) Get() any {
	p.mu.Lock()
	defer p.mu.Unlock()

	var b []byte
	if n := len(p.free); n > 0 {
		b = p.free[n-1]
		p.free = p.free[:n-1]
	} else {
		b = make([]byte, p.size)
	}
	p.outstanding[&b[0]] = struct{}{}
	return b
}

func (p *testPool) Put(v any) {
	b, ok := v.([]byte)
	if !ok {
		p.t.Errorf("Put(%T), want []byte", v)
		return
	}
	if len(b) != p.size || cap(b) != p.size {
		p.t.Errorf("Put buffer with len=%d cap=%d, want %d", len(b), cap(b), p.size)
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.outstanding[&b[0]]; !ok {
		p.t.Errorf("Put buffer %p that is not outstanding", &b[0])
		return
	}
	delete(p.outstanding, &b[0])
	for i := range b {
		b[i] = poison
	}
	p.free = append(p.free, b)
}

func (p *testPool) numOutstanding() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.outstanding)
}

// memWriterAt fails the test on overlapping or out-of-range writes, and on any
// write after close.
type memWriterAt struct {
	t      *testing.T
	jitter bool

	failOff atomic.Int64 // offset whose write fails, -1 for none

	mu      sync.Mutex
	buf     []byte
	written []bool
	calls   int
	closed  bool
}

func newMemWriterAt(t *testing.T, size int) *memWriterAt {
	m := &memWriterAt{t: t, buf: make([]byte, size), written: make([]bool, size)}
	m.failOff.Store(-1)
	return m
}

func (m *memWriterAt) WriteAt(p []byte, off int64) (int, error) {
	if m.jitter {
		jitter()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls++
	if m.closed {
		m.t.Errorf("WriteAt(off=%d, len=%d) after close", off, len(p))
	}
	if off == m.failOff.Load() {
		return 0, errInjected
	}
	if off < 0 || off+int64(len(p)) > int64(len(m.buf)) {
		m.t.Errorf("WriteAt(off=%d, len=%d) out of range [0, %d)", off, len(p), len(m.buf))
		return 0, errors.New("out of range")
	}
	for i := range p {
		if m.written[off+int64(i)] {
			m.t.Errorf("WriteAt(off=%d, len=%d) overlaps a previous write at %d", off, len(p), off+int64(i))
			break
		}
	}
	for i := range p {
		m.written[off+int64(i)] = true
	}
	copy(m.buf[off:], p)
	return len(p), nil
}

func (m *memWriterAt) close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
}

func (m *memWriterAt) numCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls
}

func (m *memWriterAt) checkComplete(t *testing.T, seed int64) {
	t.Helper()

	m.mu.Lock()
	defer m.mu.Unlock()

	for i, ok := range m.written {
		if !ok {
			t.Fatalf("offset %d was never written", i)
		}
	}
	checkContent(t, m.buf, seed)
}

func (m *memWriterAt) checkWritten(t *testing.T, seed int64) {
	t.Helper()

	m.mu.Lock()
	defer m.mu.Unlock()

	for i, ok := range m.written {
		if ok && m.buf[i] != patternByte(seed, int64(i)) {
			t.Fatalf("content mismatch at offset %d: got %#x, want %#x", i, m.buf[i], patternByte(seed, int64(i)))
		}
	}
}

type writerAtFunc func([]byte, int64) (int, error)

func (f writerAtFunc) WriteAt(p []byte, off int64) (int, error) {
	return f(p, off)
}

type gatedCall struct {
	p    []byte
	off  int64
	resp chan gatedResp
}

type gatedResp struct {
	n   int
	err error
}

func (c *gatedCall) ok() {
	c.resp <- gatedResp{n: len(c.p)}
}

func (c *gatedCall) fail(err error) {
	c.resp <- gatedResp{err: err}
}

// gatedWriterAt blocks every write until the test responds to it.
type gatedWriterAt struct {
	calls chan *gatedCall
}

func newGatedWriterAt() *gatedWriterAt {
	return &gatedWriterAt{calls: make(chan *gatedCall, numWriteWorkers+jobQueueDepth)}
}

func (g *gatedWriterAt) WriteAt(p []byte, off int64) (int, error) {
	c := &gatedCall{p: p, off: off, resp: make(chan gatedResp)}
	g.calls <- c
	r := <-c.resp
	return r.n, r.err
}

func (g *gatedWriterAt) next(t *testing.T) *gatedCall {
	t.Helper()

	select {
	case c := <-g.calls:
		return c
	case <-time.After(testTimeout):
		t.Fatalf("no write reached the sink within %v\n%s", testTimeout, goroutines())
		return nil
	}
}

// autoRespond succeeds every subsequent write. The returned func stops
// responding and reports how many writes were answered.
func (g *gatedWriterAt) autoRespond() func() int {
	var n int
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case c := <-g.calls:
				n++
				c.ok()
			case <-stop:
				return
			}
		}
	}()
	return func() int {
		close(stop)
		<-done
		return n
	}
}

// fillWorkers blocks every worker in the sink and fills the queue behind them.
func fillWorkers(t *testing.T, w *AsyncWriterAt, g *gatedWriterAt, bufSize int) []*gatedCall {
	t.Helper()

	withTimeout(t, "filling queue", func() {
		for i := range numWriteWorkers + jobQueueDepth {
			w.WriteAt(w.Buffer(), bufSize, int64(i*bufSize))
		}
	})

	inflight := make([]*gatedCall, numWriteWorkers)
	for i := range inflight {
		inflight[i] = g.next(t)
	}
	if n := len(w.queue); n != jobQueueDepth {
		t.Fatalf("queue depth = %d, want %d", n, jobQueueDepth)
	}
	return inflight
}

func patternByte(seed, off int64) byte {
	x := uint64(seed)*0x9e3779b97f4a7c15 ^ uint64(off>>3)
	x ^= x >> 30
	x *= 0xbf58476d1ce4e5b9
	x ^= x >> 27
	x *= 0x94d049bb133111eb
	x ^= x >> 31
	return byte(x >> (8 * (off & 7)))
}

func checkContent(t *testing.T, got []byte, seed int64) {
	t.Helper()

	for i := range got {
		if got[i] == patternByte(seed, int64(i)) {
			continue
		}

		lo, hi := max(0, i-8), min(len(got), i+8)
		want := make([]byte, hi-lo)
		for j := range want {
			want[j] = patternByte(seed, int64(lo+j))
		}
		t.Fatalf("content mismatch at offset %d (showing [%d, %d)):\ngot  % x\nwant % x", i, lo, hi, got[lo:hi], want)
	}
}

type chunk struct {
	off int64
	n   int
}

func planChunks(r *rand.Rand, size, bufSize int, shortN bool) []chunk {
	var cs []chunk
	for off := 0; off < size; {
		n := min(bufSize, size-off)
		if shortN {
			n = 1 + r.Intn(n)
		}
		cs = append(cs, chunk{off: int64(off), n: n})
		off += n
	}
	r.Shuffle(len(cs), func(i, j int) { cs[i], cs[j] = cs[j], cs[i] })
	return cs
}

func produce(w *AsyncWriterAt, seed int64, chunks []chunk, producers int, jit bool) {
	ch := make(chan chunk)
	var wg sync.WaitGroup
	for range producers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for c := range ch {
				buf := w.Buffer()
				for i := range c.n {
					buf[i] = patternByte(seed, c.off+int64(i))
				}
				if jit {
					jitter()
				}
				w.WriteAt(buf, c.n, c.off)
			}
		}()
	}
	for _, c := range chunks {
		ch <- c
	}
	close(ch)
	wg.Wait()
}

func jitter() {
	switch n := rand.Intn(64); {
	case n == 0:
		time.Sleep(50 * time.Microsecond)
	case n < 16:
		for range n {
			runtime.Gosched()
		}
	}
}

func goroutines() string {
	buf := make([]byte, 1<<20)
	return string(buf[:runtime.Stack(buf, true)])
}

func withTimeout(t *testing.T, what string, fn func()) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatalf("%s did not return within %v\n%s", what, testTimeout, goroutines())
	}
}

func waitClosed(t *testing.T, what string, ch <-chan struct{}) {
	t.Helper()
	withTimeout(t, what, func() { <-ch })
}

func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

// checkGoroutines must be called before anything in the test starts goroutines.
func checkGoroutines(t *testing.T) {
	base := runtime.NumGoroutine()
	t.Cleanup(func() {
		deadline := time.Now().Add(testTimeout)
		for runtime.NumGoroutine() > base {
			if time.Now().After(deadline) {
				t.Errorf("goroutines leaked: %d, want %d\n%s", runtime.NumGoroutine(), base, goroutines())
				return
			}
			time.Sleep(time.Millisecond)
		}
	})
}

func shutdown(t *testing.T, w *AsyncWriterAt, pool *testPool) {
	t.Helper()

	withTimeout(t, "Wait", w.Wait)
	withTimeout(t, "Stop", w.Stop)
	if n := pool.numOutstanding(); n != 0 {
		t.Errorf("%d buffers were never returned to the pool", n)
	}
}

func iterations(n int) int {
	if testing.Short() {
		return max(1, n/10)
	}
	return n
}

func TestAsyncWriterAtIntegrity(t *testing.T) {
	for _, bufSize := range []int{1, 7, 4096} {
		seen := map[int]bool{}
		for _, size := range []int{0, 1, bufSize - 1, bufSize, bufSize + 1, 37*bufSize + 13, 1024 * bufSize} {
			if seen[size] {
				continue
			}
			seen[size] = true

			for _, shortN := range []bool{false, true} {
				for _, producers := range []int{1, numWriteWorkers, 4 * numWriteWorkers} {
					for _, jit := range []bool{false, true} {
						name := fmt.Sprintf("buf=%d/size=%d/shortN=%v/producers=%d/jitter=%v", bufSize, size, shortN, producers, jit)
						t.Run(name, func(t *testing.T) {
							n := iterations(20)
							if size > 64*1024 {
								n = iterations(3)
							}
							for seed := range int64(n) {
								testIntegrity(t, seed, size, bufSize, shortN, producers, jit)
							}
						})
					}
				}
			}
		}
	}
}

func testIntegrity(t *testing.T, seed int64, size, bufSize int, shortN bool, producers int, jit bool) {
	checkGoroutines(t)

	pool := newTestPool(t, bufSize)
	sink := newMemWriterAt(t, size)
	sink.jitter = jit
	w := NewAsyncWriterAt(sink, pool)
	w.Start()

	chunks := planChunks(rand.New(rand.NewSource(seed)), size, bufSize, shortN)
	withTimeout(t, "producers", func() { produce(w, seed, chunks, producers, jit) })
	shutdown(t, w, pool)
	sink.close()

	if err := w.Error(); err != nil {
		t.Fatalf("seed %d: unexpected error: %v", seed, err)
	}
	if isClosed(w.Done()) {
		t.Fatalf("seed %d: Done closed on success", seed)
	}
	if got, want := sink.numCalls(), len(chunks); got != want {
		t.Fatalf("seed %d: %d writes reached the sink, want %d", seed, got, want)
	}
	sink.checkComplete(t, seed)
}

func TestAsyncWriterAtBuffer(t *testing.T) {
	pool := newTestPool(t, 32)
	w := NewAsyncWriterAt(newMemWriterAt(t, 0), pool)

	b := w.Buffer()
	if len(b) != 32 || cap(b) != 32 {
		t.Fatalf("Buffer() len=%d cap=%d, want 32", len(b), cap(b))
	}
	if n := pool.numOutstanding(); n != 1 {
		t.Fatalf("outstanding buffers = %d, want 1", n)
	}

	w.Start()
	w.WriteAt(b, 0, 0)
	shutdown(t, w, pool)
}

func TestAsyncWriterAtRelease(t *testing.T) {
	pool := newTestPool(t, 8)
	w := NewAsyncWriterAt(newMemWriterAt(t, 0), pool)

	w.Release(w.Buffer())
	if n := pool.numOutstanding(); n != 0 {
		t.Errorf("%d buffers were never returned to the pool", n)
	}
}

func TestAsyncWriterAtWritesPrefix(t *testing.T) {
	checkGoroutines(t)

	pool := newTestPool(t, 16)
	var got []byte
	var gotOff int64
	w := NewAsyncWriterAt(writerAtFunc(func(p []byte, off int64) (int, error) {
		got, gotOff = p, off
		return len(p), nil
	}), pool)
	w.Start()

	buf := w.Buffer()
	w.WriteAt(buf, 5, 100)
	withTimeout(t, "Wait", w.Wait)

	if len(got) != 5 {
		t.Errorf("sink got len %d, want 5", len(got))
	}
	if len(got) > 0 && &got[0] != &buf[0] {
		t.Error("sink got a different backing array than the submitted buffer")
	}
	if gotOff != 100 {
		t.Errorf("sink got offset %d, want 100", gotOff)
	}
	shutdown(t, w, pool)
}

func TestAsyncWriterAtZeroLengthWrite(t *testing.T) {
	checkGoroutines(t)

	pool := newTestPool(t, 16)
	sink := newMemWriterAt(t, 16)
	w := NewAsyncWriterAt(sink, pool)
	w.Start()

	w.WriteAt(w.Buffer(), 0, 0)
	w.WriteAt(w.Buffer(), 0, 16)
	shutdown(t, w, pool)

	if err := w.Error(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, ok := range sink.written {
		if ok {
			t.Fatalf("offset %d was written by a zero-length write", i)
		}
	}
}

func TestAsyncWriterAtWriteBeforeStart(t *testing.T) {
	checkGoroutines(t)

	const bufSize = 8
	pool := newTestPool(t, bufSize)
	sink := newMemWriterAt(t, jobQueueDepth*bufSize)
	w := NewAsyncWriterAt(sink, pool)

	chunks := planChunks(rand.New(rand.NewSource(1)), jobQueueDepth*bufSize, bufSize, false)
	withTimeout(t, "WriteAt before Start", func() { produce(w, 1, chunks, 1, false) })
	if n := sink.numCalls(); n != 0 {
		t.Fatalf("%d writes reached the sink before Start", n)
	}

	w.Start()
	shutdown(t, w, pool)
	sink.checkComplete(t, 1)
}

func TestAsyncWriterAtStartIdempotent(t *testing.T) {
	checkGoroutines(t)

	pool := newTestPool(t, 8)
	w := NewAsyncWriterAt(newMemWriterAt(t, 0), pool)

	base := runtime.NumGoroutine()
	w.Start()
	w.Start()
	w.Start()
	if got := runtime.NumGoroutine() - base; got != numWriteWorkers {
		t.Errorf("started %d workers, want %d", got, numWriteWorkers)
	}
	shutdown(t, w, pool)
}

func TestAsyncWriterAtStopIdempotent(t *testing.T) {
	checkGoroutines(t)

	pool := newTestPool(t, 8)
	w := NewAsyncWriterAt(newMemWriterAt(t, 0), pool)
	w.Start()

	withTimeout(t, "concurrent Stop", func() {
		var wg sync.WaitGroup
		for range 8 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w.Stop()
			}()
		}
		wg.Wait()
	})
	withTimeout(t, "second Stop", w.Stop)
	withTimeout(t, "Wait", w.Wait)
}

func TestAsyncWriterAtStopWithoutStart(t *testing.T) {
	checkGoroutines(t)

	pool := newTestPool(t, 8)
	sink := newMemWriterAt(t, 8*jobQueueDepth)
	w := NewAsyncWriterAt(sink, pool)

	for i := range jobQueueDepth {
		w.WriteAt(w.Buffer(), 8, int64(i*8))
	}
	withTimeout(t, "Stop", w.Stop)
	withTimeout(t, "Wait", w.Wait)

	if n := pool.numOutstanding(); n != 0 {
		t.Errorf("%d buffers were never returned to the pool", n)
	}
	if n := sink.numCalls(); n != 0 {
		t.Errorf("%d writes reached the sink", n)
	}
}

func TestAsyncWriterAtStartAfterStop(t *testing.T) {
	checkGoroutines(t)

	pool := newTestPool(t, 8)
	w := NewAsyncWriterAt(newMemWriterAt(t, 0), pool)
	w.Stop()
	w.Start()
	withTimeout(t, "Wait", w.Wait)
}

func TestAsyncWriterAtWaitWithoutJobs(t *testing.T) {
	checkGoroutines(t)

	pool := newTestPool(t, 8)
	w := NewAsyncWriterAt(newMemWriterAt(t, 0), pool)
	w.Start()
	shutdown(t, w, pool)
	if isClosed(w.Done()) {
		t.Error("Done closed without a failure")
	}
}

func TestAsyncWriterAtWriteError(t *testing.T) {
	cases := map[string]struct {
		n       func(len int) int
		err     error
		wantErr error
	}{
		"error":                 {n: func(int) int { return 0 }, err: errInjected, wantErr: errInjected},
		"error with full n":     {n: func(l int) int { return l }, err: errInjected, wantErr: errInjected},
		"short write":           {n: func(l int) int { return l - 1 }, wantErr: io.ErrShortWrite},
		"zero write":            {n: func(int) int { return 0 }, wantErr: io.ErrShortWrite},
		"write longer than p":   {n: func(l int) int { return l + 1 }, wantErr: io.ErrShortWrite},
		"short write and error": {n: func(l int) int { return l / 2 }, err: errInjected, wantErr: errInjected},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			checkGoroutines(t)

			pool := newTestPool(t, 8)
			w := NewAsyncWriterAt(writerAtFunc(func(p []byte, _ int64) (int, error) {
				return c.n(len(p)), c.err
			}), pool)
			w.Start()

			w.WriteAt(w.Buffer(), 8, 0)
			waitClosed(t, "Done", w.Done())
			if err := w.Error(); !errors.Is(err, c.wantErr) {
				t.Fatalf("Error() = %v, want %v", err, c.wantErr)
			}

			shutdown(t, w, pool)
			if err := w.Error(); !errors.Is(err, c.wantErr) {
				t.Fatalf("Error() after Stop = %v, want %v", err, c.wantErr)
			}
		})
	}
}

func TestAsyncWriterAtErrorDiscardsQueued(t *testing.T) {
	checkGoroutines(t)

	g := newGatedWriterAt()
	pool := newTestPool(t, 8)
	w := NewAsyncWriterAt(g, pool)
	w.Start()

	inflight := fillWorkers(t, w, g, 8)
	inflight[0].fail(errInjected)
	waitClosed(t, "Done", w.Done())

	for _, c := range inflight[1:] {
		c.ok()
	}
	extra := g.autoRespond()
	shutdown(t, w, pool)

	if n := extra(); n != 0 {
		t.Errorf("%d queued writes reached the sink after failure", n)
	}
	if err := w.Error(); !errors.Is(err, errInjected) {
		t.Errorf("Error() = %v, want %v", err, errInjected)
	}
}

func TestAsyncWriterAtConcurrentErrors(t *testing.T) {
	for seed := range int64(iterations(500)) {
		t.Run(fmt.Sprint(seed), func(t *testing.T) {
			checkGoroutines(t)

			g := newGatedWriterAt()
			pool := newTestPool(t, 8)
			w := NewAsyncWriterAt(g, pool)
			w.Start()

			inflight := fillWorkers(t, w, g, 8)
			errs := make(map[error]bool, len(inflight))
			release := make(chan struct{})
			var wg sync.WaitGroup
			for i, c := range inflight {
				err := fmt.Errorf("write %d: %w", i, errInjected)
				errs[err] = true
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-release
					c.fail(err)
				}()
			}
			close(release)
			wg.Wait()

			waitClosed(t, "Done", w.Done())
			first := w.Error()
			if !errs[first] {
				t.Fatalf("Error() = %v, want one of the injected errors", first)
			}

			extra := g.autoRespond()
			shutdown(t, w, pool)
			if n := extra(); n != 0 {
				t.Errorf("%d queued writes reached the sink after failure", n)
			}
			if err := w.Error(); err != first {
				t.Errorf("Error() changed from %v to %v", first, err)
			}
		})
	}
}

func TestAsyncWriterAtErrorUnblocksProducers(t *testing.T) {
	checkGoroutines(t)

	g := newGatedWriterAt()
	pool := newTestPool(t, 8)
	w := NewAsyncWriterAt(g, pool)
	w.Start()

	inflight := fillWorkers(t, w, g, 8)

	// the queue is full, so these block until something gives
	const producers = 8
	var wg sync.WaitGroup
	for i := range producers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.WriteAt(w.Buffer(), 8, int64((numWriteWorkers+jobQueueDepth+i)*8))
		}()
	}

	inflight[0].fail(errInjected)
	withTimeout(t, "blocked producers", wg.Wait)

	for _, c := range inflight[1:] {
		c.ok()
	}
	extra := g.autoRespond()
	shutdown(t, w, pool)

	if n := extra(); n != 0 {
		t.Errorf("%d writes reached the sink after failure", n)
	}
}

func TestAsyncWriterAtWriteAfterError(t *testing.T) {
	checkGoroutines(t)

	var calls atomic.Int64
	pool := newTestPool(t, 8)
	w := NewAsyncWriterAt(writerAtFunc(func([]byte, int64) (int, error) {
		calls.Add(1)
		return 0, errInjected
	}), pool)
	w.Start()

	w.WriteAt(w.Buffer(), 8, 0)
	waitClosed(t, "Done", w.Done())

	withTimeout(t, "WriteAt after failure", func() {
		for i := range 1000 {
			w.WriteAt(w.Buffer(), 8, int64(i*8))
		}
	})
	shutdown(t, w, pool)

	if n := calls.Load(); n != 1 {
		t.Errorf("%d writes reached the sink, want 1", n)
	}
}

// Stop has to wait for the failing write, and Error must reflect it once Stop
// returns.
func TestAsyncWriterAtStopDuringFailingWrite(t *testing.T) {
	checkGoroutines(t)

	g := newGatedWriterAt()
	pool := newTestPool(t, 8)
	w := NewAsyncWriterAt(g, pool)
	w.Start()

	w.WriteAt(w.Buffer(), 8, 0)
	c := g.next(t)

	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		w.Stop()
	}()
	waitClosed(t, "stop signal", w.stop)
	if isClosed(stopped) {
		t.Fatal("Stop returned while a write was in progress")
	}

	c.fail(errInjected)
	waitClosed(t, "Stop", stopped)

	if err := w.Error(); !errors.Is(err, errInjected) {
		t.Errorf("Error() = %v, want %v", err, errInjected)
	}
	if !isClosed(w.Done()) {
		t.Error("Done not closed after failure")
	}
	withTimeout(t, "Wait", w.Wait)
	if n := pool.numOutstanding(); n != 0 {
		t.Errorf("%d buffers were never returned to the pool", n)
	}
}

// Stop without Wait (the downloader's error path): in-flight writes finish,
// queued ones may or may not be written, and nothing is written after Stop
// returns.
func TestAsyncWriterAtStopDiscardsQueued(t *testing.T) {
	for seed := range int64(iterations(500)) {
		t.Run(fmt.Sprint(seed), func(t *testing.T) {
			checkGoroutines(t)

			g := newGatedWriterAt()
			pool := newTestPool(t, 8)
			w := NewAsyncWriterAt(g, pool)
			w.Start()

			inflight := fillWorkers(t, w, g, 8)

			stopped := make(chan struct{})
			go func() {
				defer close(stopped)
				w.Stop()
			}()
			waitClosed(t, "stop signal", w.stop)

			extra := g.autoRespond()
			for _, c := range inflight {
				c.ok()
			}
			waitClosed(t, "Stop", stopped)
			n := extra()

			withTimeout(t, "Wait", w.Wait)
			if n := pool.numOutstanding(); n != 0 {
				t.Errorf("%d buffers were never returned to the pool", n)
			}
			if n := len(g.calls); n != 0 {
				t.Errorf("%d writes reached the sink after Stop returned", n)
			}
			if n > jobQueueDepth {
				t.Errorf("%d queued writes processed, but only %d were queued", n, jobQueueDepth)
			}
			if err := w.Error(); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// Randomized mix of producers, failures, and both shutdown orders. Outcomes
// aren't predicted exactly, only the invariants that must always hold.
func TestAsyncWriterAtStress(t *testing.T) {
	for seed := range int64(iterations(1000)) {
		t.Run(fmt.Sprint(seed), func(t *testing.T) {
			checkGoroutines(t)

			r := rand.New(rand.NewSource(seed))
			size := r.Intn(64 * 1024)
			bufSize := 1 + r.Intn(512)
			producers := 1 + r.Intn(4*numWriteWorkers)
			chunks := planChunks(r, size, bufSize, r.Intn(2) == 0)
			stopEarly := r.Intn(2) == 0

			pool := newTestPool(t, bufSize)
			sink := newMemWriterAt(t, size)
			sink.jitter = true
			injectFailure := len(chunks) > 0 && r.Intn(2) == 0
			if injectFailure {
				sink.failOff.Store(chunks[r.Intn(len(chunks))].off)
			}

			w := NewAsyncWriterAt(sink, pool)
			w.Start()
			withTimeout(t, "producers", func() { produce(w, seed, chunks, producers, true) })

			if stopEarly {
				withTimeout(t, "Stop", w.Stop)
				sink.close()
				withTimeout(t, "Wait", w.Wait)
			} else {
				withTimeout(t, "Wait", w.Wait)
				withTimeout(t, "Stop", w.Stop)
				sink.close()
			}

			if n := pool.numOutstanding(); n != 0 {
				t.Fatalf("%d buffers were never returned to the pool", n)
			}
			if n := sink.numCalls(); n > len(chunks) {
				t.Fatalf("%d writes reached the sink, but only %d were submitted", n, len(chunks))
			}
			sink.checkWritten(t, seed)

			err := w.Error()
			switch {
			case injectFailure:
				if !errors.Is(err, errInjected) && !(stopEarly && err == nil) {
					t.Fatalf("Error() = %v, want %v", err, errInjected)
				}
				if err != nil && !isClosed(w.Done()) {
					t.Fatal("Done not closed after failure")
				}
			case err != nil:
				t.Fatalf("unexpected error: %v", err)
			case !stopEarly:
				sink.checkComplete(t, seed)
			}
		})
	}
}
