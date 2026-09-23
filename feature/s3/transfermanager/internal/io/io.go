package io

import (
	"io"
	"sync"
	"unsafe"
)

var Pools BufferPools

const alignedBy = 4096

// BufferPools retains a separate sync.Pool for each buffer size.
type BufferPools struct {
	mu    sync.Mutex
	pools map[int]*sync.Pool
}

// https://en.wikipedia.org/wiki/Data_structure_alignment#Computing_padding
//
// IMPORTANT: this only works when alignedBy is a power of 2
func align[T uintptr | int64](addr T) T {
	return (addr + (alignedBy - 1)) &^ (alignedBy - 1)
}

func makealigned(size int) []byte {
	p := make([]byte, size+alignedBy)
	u := unsafe.Pointer(&p[0])
	addr := uintptr(u)

	off := align(addr) - addr
	return p[off : int(off)+size : int(off)+size]
}

// Pool returns the pool for buffers of the requested size.
//
// DO NOT repeatedly call Pool() in a transfer operation. Grab the pool of the
// size you need once and retain a reference to it.
func (bps *BufferPools) Pool(size int) *sync.Pool {
	bps.mu.Lock()
	defer bps.mu.Unlock()

	if bps.pools == nil {
		bps.pools = make(map[int]*sync.Pool)
	}
	if p, ok := bps.pools[size]; ok {
		return p
	}

	p := &sync.Pool{
		New: func() any {
			return makealigned(size)
		},
	}
	bps.pools[size] = p
	return p
}

// AsyncWriterAt wraps an io.WriterAt to expose an interface which writes to it
// asynchronously.
//
// The async writer also owns all of the buffers it uses to write. The caller
// requests a buffer from it using Buffer(), and WriteAt immediately reclaims
// ownership of that buffer.
type AsyncWriterAt struct {
	w            io.WriterAt
	bufs         *sync.Pool
	startWorkers int
	maxWorkers   int

	wg       sync.WaitGroup
	queue    chan writeAtJob
	doneOnce sync.Once
	done     chan struct{}
	errOnce  sync.Once
	err      error
}

func NewAsyncWriterAt(w io.WriterAt, bufs *sync.Pool, startWorkers, maxWorkers, queueDepth int) *AsyncWriterAt {
	return &AsyncWriterAt{
		w:            w,
		bufs:         bufs,
		done:         make(chan struct{}),
		startWorkers: startWorkers,
		maxWorkers:   maxWorkers,
		queue:        make(chan writeAtJob, queueDepth),
	}
}

// Buffer returns a pooled byte buffer. The caller will fill this buffer and
// then pass it back to the dispatcher via WriteAt.
func (w *AsyncWriterAt) Buffer() []byte {
	return w.bufs.Get().([]byte)
}

// WriteAt queues the bytes for writing.
//
// Unlike the synchronous io.WriteAt:
//   - The write length is explicitly passed in n, because the job that
//     eventually performs the write needs the original slice header of p (i.e.
//     NOT a subslice) so it can return the full slice to the pool.
//   - This method retains p, callers MUST NOT retain or modify p.
func (w *AsyncWriterAt) WriteAt(p []byte, n int, off int64) {
	w.wg.Add(1)

	select {
	case <-w.done:
		w.wg.Done()
		w.bufs.Put(p)
		return
	case w.queue <- writeAtJob{p, n, off}:
	}
}

// Done returns a channel that's closed when the writer is done.
func (w *AsyncWriterAt) Done() chan struct{} {
	return w.done
}

// Error returns the write error (if any) encountered during async write.
func (w *AsyncWriterAt) Error() error {
	return w.err
}

// Start spins up write workers.
func (w *AsyncWriterAt) Start() {
	// TODO maxWorkers
	for range w.startWorkers {
		go w.doWrites()
	}
}

// Stop immediately terminates all write workers.
func (w *AsyncWriterAt) Stop() {
	w.doneOnce.Do(func() {
		close(w.done)
		for {
			select {
			case job := <-w.queue:
				w.bufs.Put(job.p)
				w.wg.Done()
			default:
				return
			}
		}
	})
}

// Wait blocks until internal WaitGroup counter of the dispatcher is zero,
// which means that there are no more write jobs queued.
//
// By contrast, Wait DOES NOT signal that no future write jobs will be
// submitted. The caller MUST first Wait() on the internal downloader's
// WaitGroup to confirm that.
func (w *AsyncWriterAt) Wait() {
	select {
	case <-w.done:
	default:
		w.wg.Wait()
	}
}

func (w *AsyncWriterAt) doWrites() {
	for {
		select {
		case <-w.done:
			return
		case job := <-w.queue:
			n, err := w.w.WriteAt(job.p[:job.n], job.off)
			if err == nil && n != job.n {
				err = io.ErrShortWrite
			}
			if err != nil {
				w.errOnce.Do(func() { w.err = err })
				w.Stop()
			}

			w.wg.Done()
			w.bufs.Put(job.p)
		}
	}
}

type writeAtJob struct {
	p   []byte
	n   int
	off int64
}

// File is a lazily-initialized download destination.
type File interface {
	io.WriterAt
	Init(int64, int64) error
	Close() error
}
