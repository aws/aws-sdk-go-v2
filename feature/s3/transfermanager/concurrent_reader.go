package transfermanager

import (
	"io"
	"sync"
)

// ConcurrentReader receives object parts from working goroutines, composes those chunks in order and read
// to user's buffer. ConcurrentReader limits the max number of chunks it could receive and read at the same
// time so getter won't send following parts' request to s3 until user reads all current chunks, which avoids
// too much memory consumption when caching large object parts
type ConcurrentReader struct {
	ch  chan outChunk
	buf map[int32]*outChunk

	partsCount int32
	capacity   int32
	count      int32
	read       int32

	written  int64
	partSize int64

	m sync.Mutex

	err error
}

// NewConcurrentReader returns a ConcurrentReader used in GetObject input
func NewConcurrentReader() *ConcurrentReader {
	return &ConcurrentReader{
		buf:      make(map[int32]*outChunk),
		partSize: 1, // just a placeholder value
	}
}

// Read implements io.Reader to compose object parts in order and read to p.
// It will receive up to r.capacity chunks, read them to p if any chunk index
// fits into p scope, otherwise it will buffer those chunks and read them in
// following calls
func (r *ConcurrentReader) Read(p []byte) (int, error) {
	if cap(p) == 0 {
		return 0, nil
	}

	var written int

	for r.count < r.getCapacity() {
		if e := r.getErr(); e != nil && e != io.EOF {
			r.written += int64(written)
			r.clean()
			return written, r.getErr()
		}
		if written >= cap(p) {
			r.written += int64(written)
			return written, r.getErr()
		}

		oc, ok := <-r.ch
		if !ok {
			r.written += int64(written)
			return written, r.getErr()
		}

		r.count++
		index := r.getPartSize()*int64(oc.index) - r.written

		if index < int64(cap(p)) {
			n, err := oc.body.Read(p[index:])
			oc.cur += int64(n)
			written += n
			if err != nil && err != io.EOF {
				r.setErr(err)
				r.clean()
				r.written += int64(written)
				return written, r.getErr()
			}
		}
		if oc.cur < oc.length {
			r.buf[oc.index] = &oc
		} else {
			r.incrRead(1)
			if r.getRead() >= r.partsCount {
				r.setErr(io.EOF)
			}
		}
	}

	partSize := r.getPartSize()
	minIndex := int32(r.written / partSize)
	maxIndex := min(int32((r.written+int64(cap(p))-1)/partSize), r.getCapacity()-1)
	for i := minIndex; i <= maxIndex; i++ {
		if e := r.getErr(); e != nil && e != io.EOF {
			r.written += int64(written)
			r.clean()
			return written, r.getErr()
		}

		c, ok := r.buf[i]
		if ok {
			index := int64(i)*partSize + c.cur - r.written
			n, err := c.body.Read(p[index:])
			c.cur += int64(n)
			written += n
			if err != nil && err != io.EOF {
				r.setErr(err)
				r.clean()
				r.written += int64(written)
				return written, r.getErr()
			}
			if c.cur >= c.length {
				r.incrRead(1)
				delete(r.buf, i)
				if r.getRead() >= r.partsCount {
					r.setErr(io.EOF)
				}
			}
		}
	}

	r.written += int64(written)
	return written, r.getErr()
}

func (r *ConcurrentReader) setPartSize(n int64) {
	r.m.Lock()
	defer r.m.Unlock()

	r.partSize = n
}

func (r *ConcurrentReader) getPartSize() int64 {
	r.m.Lock()
	defer r.m.Unlock()

	return r.partSize
}

func (r *ConcurrentReader) setCapacity(n int32) {
	r.m.Lock()
	defer r.m.Unlock()

	r.capacity = n
}

func (r *ConcurrentReader) getCapacity() int32 {
	r.m.Lock()
	defer r.m.Unlock()

	return r.capacity
}

func (r *ConcurrentReader) setPartsCount(n int32) {
	r.m.Lock()
	defer r.m.Unlock()

	r.partsCount = n
}

func (r *ConcurrentReader) incrRead(n int32) {
	r.m.Lock()
	defer r.m.Unlock()

	r.read += n
}

func (r *ConcurrentReader) getRead() int32 {
	r.m.Lock()
	defer r.m.Unlock()

	return r.read
}

func (r *ConcurrentReader) setErr(err error) {
	r.m.Lock()
	defer r.m.Unlock()

	r.err = err
}

func (r *ConcurrentReader) getErr() error {
	r.m.Lock()
	defer r.m.Unlock()

	return r.err
}

func (r *ConcurrentReader) clean() {
	for {
		_, ok := <-r.ch
		if !ok {
			break
		}
	}
}
