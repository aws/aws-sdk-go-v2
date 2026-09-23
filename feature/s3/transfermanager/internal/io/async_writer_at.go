package io

import (
	"io"
	"sync"
)

// These values are tuned around single-file download on a high-throughput disk
// since that's the most intensive thing we need scale to. Excess
// goroutines/queue depth cost mostly nothing, excess goroutines will just
// immediately park and never do real work, and queue depth is just some extra
// channel memory.
const (
	numWriteWorkers = 32
	jobQueueDepth   = 64
)

// AsyncWriterAt wraps an io.WriterAt to expose an interface which writes to it
// asynchronously.
//
// Asynchronous write dispatch allows downloading goroutines to maximize the
// amount of time they spend actually doing network i/o, which is crucial for
// achieving high sustained throughput under many concurrent connections.
//
// The async writer also owns all of the buffers it uses to write. The caller
// requests a buffer from it using Buffer(), and WriteAt immediately reclaims
// ownership of that buffer.
type AsyncWriterAt struct {
	w    io.WriterAt
	bufs *sync.Pool

	jobs    sync.WaitGroup // tracks write jobs in queue
	workers sync.WaitGroup // tracks actual write goroutines
	queue   chan writeAtJob

	// stop is the explicit "shut it down" from the outside caller
	// failed is "a write failed, shut it down" internally
	stopOnce sync.Once
	stop     chan struct{}
	failed   chan struct{}
	errOnce  sync.Once
	err      error
}

// NewAsyncWriterAt initializes an async writer for the given sink and buffer pool.
func NewAsyncWriterAt(w io.WriterAt, bufs *sync.Pool) *AsyncWriterAt {
	return &AsyncWriterAt{
		w:      w,
		bufs:   bufs,
		queue:  make(chan writeAtJob, jobQueueDepth),
		stop:   make(chan struct{}),
		failed: make(chan struct{}),
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
	w.jobs.Add(1)

	select {
	case <-w.failed:
		w.jobs.Done()
		w.bufs.Put(p)
		return
	case <-w.stop:
		w.jobs.Done()
		w.bufs.Put(p)
		return
	case w.queue <- writeAtJob{p: p, n: n, off: off}:
	}
}

// Done returns a channel that's closed when a write fails.
func (w *AsyncWriterAt) Done() <-chan struct{} {
	return w.failed
}

// Error returns the first write error, if any.
func (w *AsyncWriterAt) Error() error {
	return w.err
}

// Start spins up write workers.
func (w *AsyncWriterAt) Start() {
	w.workers.Add(numWriteWorkers)
	for range numWriteWorkers {
		go func() {
			defer w.workers.Done()
			w.doWrites()
		}()
	}
}

// Stop terminates workers and releases any jobs they did not process.
func (w *AsyncWriterAt) Stop() {
	w.stopOnce.Do(func() {
		close(w.stop)
		w.workers.Wait()
		for {
			select {
			case job := <-w.queue:
				w.bufs.Put(job.p)
				w.jobs.Done()
			default:
				return
			}
		}
	})
}

// Wait blocks until all submitted jobs have completed or been discarded.
// The caller MUST first wait for all producers to stop submitting jobs.
func (w *AsyncWriterAt) Wait() {
	w.jobs.Wait()
}

func (w *AsyncWriterAt) fail(err error) {
	w.errOnce.Do(func() {
		w.err = err
		close(w.failed)
	})
}

func (w *AsyncWriterAt) doWrites() {
	for {
		select {
		case <-w.stop:
			return
		case job := <-w.queue:
			select {
			case <-w.failed:
				w.bufs.Put(job.p)
				w.jobs.Done()
				continue
			default:
			}

			n, err := w.w.WriteAt(job.p[:job.n], job.off)
			if err == nil && n != job.n {
				err = io.ErrShortWrite
			}
			if err != nil {
				w.fail(err)
			}

			w.jobs.Done()
			w.bufs.Put(job.p)
		}
	}
}

type writeAtJob struct {
	p   []byte
	n   int
	off int64
}
