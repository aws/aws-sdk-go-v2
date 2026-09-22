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
			// we do aligned allocs no matter what b/c it's not that much extra
			p := make([]byte, size+alignedBy)
			u := unsafe.Pointer(&p[0])
			addr := uintptr(u)

			// https://en.wikipedia.org/wiki/Data_structure_alignment#Computing_padding
			aligned := (addr + (alignedBy - 1)) &^ (alignedBy - 1)
			off := aligned - addr
			return p[off : int(off)+size : int(off)+size]
		},
	}
	bps.pools[size] = p
	return p
}

// AsyncWriterAt wraps an io.WraterAt to expose an interface which writes to it
// asynchronously.
type AsyncWriterAt struct {
	w    io.WriterAt
	err  chan error
	done chan struct{}

	startWorkers, maxWorkers int
	queue                    chan writeAtJob
}

func NewAsyncWriterAt(w io.WriterAt, startWorkers, maxWorkers, queueDepth int) *AsyncWriterAt {
	return &AsyncWriterAt{
		w:            w,
		err:          make(chan error, 1),
		done:         make(chan struct{}, 1),
		startWorkers: startWorkers,
		maxWorkers:   maxWorkers,
		queue:        make(chan writeAtJob, queueDepth),
	}
}

// WriteAt queues the bytes for writing.
//
// WriteAt retains p, callers MUST NOT retain or modify p.
func (w *AsyncWriterAt) WriteAt(p []byte, off int64) {
	w.queue <- writeAtJob{p, off}
}

func (w *AsyncWriterAt) Error() chan error {
	return w.err
}

func (w *AsyncWriterAt) Start() {
	// TODO maxWorkers
	for range w.startWorkers {
		go w.doWrites()
	}
}

func (w *AsyncWriterAt) Stop() {
	close(w.done)
}

func (w *AsyncWriterAt) doWrites() {
	for {
		select {
		case <-w.done:
			return
		case job := <-w.queue:
			if _, err := w.w.WriteAt(job.p, job.off); err != nil {
				w.err <- err
			}
		}
	}
}

type writeAtJob struct {
	p   []byte
	off int64
}

type File interface {
	io.WriterAt
	Init(int64, int64) error
	Close() error
}
