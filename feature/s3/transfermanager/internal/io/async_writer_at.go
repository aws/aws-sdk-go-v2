package io

import (
	"fmt"
	"io"
	"slices"
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
	bufs BufferPool

	jobs    sync.WaitGroup // tracks write jobs in queue
	workers sync.WaitGroup // tracks actual write goroutines
	queue   chan writeAtJob

	// stop is the explicit "shut it down" from the outside caller
	// failed is "a write failed, shut it down" internally
	startOnce sync.Once
	stopOnce  sync.Once
	stop      chan struct{}
	failed    chan struct{}
	errOnce   sync.Once
	err       error

	// landed maps the offset of each completed write to its length. A chunk
	// retry rewrites the offsets of the attempt it replaces, so the latest
	// write at an offset wins.
	landedMu sync.Mutex
	landed   map[int64]int
}

// NewAsyncWriterAt initializes an async writer for the given sink and buffer pool.
func NewAsyncWriterAt(w io.WriterAt, bufs BufferPool) *AsyncWriterAt {
	return &AsyncWriterAt{
		w:      w,
		bufs:   bufs,
		queue:  make(chan writeAtJob, jobQueueDepth),
		stop:   make(chan struct{}),
		failed: make(chan struct{}),
		landed: make(map[int64]int),
	}
}

// Buffer returns a pooled byte buffer. The caller will fill this buffer and
// then pass it back to the dispatcher via WriteAt.
func (w *AsyncWriterAt) Buffer() []byte {
	return w.bufs.Get().([]byte)
}

// Release returns a buffer obtained from Buffer that will not be passed to
// WriteAt.
func (w *AsyncWriterAt) Release(p []byte) {
	w.bufs.Put(p)
}

// WriteAt queues the bytes for writing.
//
// Unlike the synchronous io.WriteAt:
//   - The write length is explicitly passed in n, because the job that
//     eventually performs the write needs the original slice header of p (i.e.
//     NOT a subslice) so it can return the full slice to the pool.
//   - This method retains p, callers MUST NOT retain or modify p.
//   - If done is non-nil, done.Done() is called once the job has been written
//     or discarded. Callers use this to wait for a subset of their writes to
//     land without waiting on the whole writer.
//
// WriteAt MUST NOT be called concurrently with or after Stop. Calls made
// before Start are queued, but at most jobQueueDepth of them can be pending
// before WriteAt blocks until Start is called.
func (w *AsyncWriterAt) WriteAt(p []byte, n int, off int64, done *sync.WaitGroup) {
	job := writeAtJob{p: p, n: n, off: off, done: done}
	w.jobs.Add(1)
	if done != nil {
		done.Add(1)
	}

	select {
	case <-w.failed:
		w.finish(job)
	case <-w.stop:
		w.finish(job)
	case w.queue <- job:
	}
}

// Done returns a channel that's closed when a write fails.
func (w *AsyncWriterAt) Done() <-chan struct{} {
	return w.failed
}

// Error returns the first write error, if any.
//
// Error is only safe to call after Done is closed or Stop has returned.
func (w *AsyncWriterAt) Error() error {
	return w.err
}

// Start spins up write workers. Subsequent calls are no-ops.
func (w *AsyncWriterAt) Start() {
	w.startOnce.Do(func() {
		w.workers.Add(numWriteWorkers)
		for range numWriteWorkers {
			go func() {
				defer w.workers.Done()
				w.doWrites()
			}()
		}
	})
}

// Stop terminates workers and releases any jobs they did not process. Writes
// already in progress complete before Stop returns, and no writes to the
// underlying io.WriterAt occur after it returns.
//
// All WriteAt calls MUST have returned before Stop is called. Subsequent calls
// are no-ops.
func (w *AsyncWriterAt) Stop() {
	w.stopOnce.Do(func() {
		close(w.stop)
		w.workers.Wait()
		for {
			select {
			case job := <-w.queue:
				w.finish(job)
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

// CheckContiguous verifies that the writes which landed, ordered by offset,
// exactly tile [0, size): the first begins at 0, each begins where the
// previous one ended, and the last ends at size.
//
// The caller MUST first Wait for all submitted jobs. This is meaningful only
// when no write failed or was discarded.
func (w *AsyncWriterAt) CheckContiguous(size int64) error {
	w.landedMu.Lock()
	offs := make([]int64, 0, len(w.landed))
	for off := range w.landed {
		offs = append(offs, off)
	}
	slices.Sort(offs)

	var next int64
	for _, off := range offs {
		if off != next {
			w.landedMu.Unlock()
			return fmt.Errorf("write at offset %d, expected %d", off, next)
		}
		next += int64(w.landed[off])
	}
	w.landedMu.Unlock()

	if next != size {
		return fmt.Errorf("writes end at offset %d, expected %d", next, size)
	}
	return nil
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
				w.finish(job)
				continue
			default:
			}

			n, err := w.w.WriteAt(job.p[:job.n], job.off)
			if err == nil && n != job.n {
				err = io.ErrShortWrite
			}
			if err != nil {
				w.fail(err)
			} else if job.n > 0 {
				w.landedMu.Lock()
				w.landed[job.off] = job.n
				w.landedMu.Unlock()
			}

			w.finish(job)
		}
	}
}

// Every job, whether written or discarded, MUST be released through finish.
func (w *AsyncWriterAt) finish(job writeAtJob) {
	w.bufs.Put(job.p)
	if job.done != nil {
		job.done.Done()
	}
	w.jobs.Done()
}

type writeAtJob struct {
	p    []byte
	n    int
	off  int64
	done *sync.WaitGroup
}
