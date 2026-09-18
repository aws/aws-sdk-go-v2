package transfermanager

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

const (
	downloadFileVectorChunkSize = writeBehindChunkSize
	downloadFileVectorCount     = 4
)

type vectorWriterAt interface {
	writeVectorAt([][]byte, int64) (int, error)
}

type ownedVectorWrite struct {
	buf      []byte
	n        int64
	off      int64
	complete func(int, error)
}

type vectorWriteGroup struct {
	slots []*ownedVectorWrite
	count int
}

type groupedVectorWriterAt struct {
	w         vectorWriterAt
	chunkSize int64
	count     int
	groupSize int64
	alignment int64
	padTail   bool

	mu     sync.Mutex
	groups map[int64]*vectorWriteGroup
	closed bool
	err    error
}

func newGroupedVectorWriterAt(w vectorWriterAt, chunkSize int64, count int, padTail bool) *groupedVectorWriterAt {
	alignment := int64(1)
	if provider, ok := w.(writeBufferAlignmentProvider); ok {
		alignment = provider.writeBufferAlignment()
	}
	return &groupedVectorWriterAt{
		w:         w,
		chunkSize: chunkSize,
		count:     count,
		groupSize: chunkSize * int64(count),
		alignment: alignment,
		padTail:   padTail,
		groups:    map[int64]*vectorWriteGroup{},
	}
}

func (w *groupedVectorWriterAt) writeBufferAlignment() int64 {
	return w.alignment
}

func (w *groupedVectorWriterAt) preallocate(size int64) error {
	if preallocator, ok := w.w.(downloadFilePreallocator); ok {
		return preallocator.preallocate(size)
	}
	return nil
}

func (w *groupedVectorWriterAt) WriteAt(p []byte, off int64) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if off < 0 {
		return 0, fmt.Errorf("negative write offset %d", off)
	}
	if err := w.currentError(); err != nil {
		return 0, err
	}

	n, err := writeAllVectors(w.w, [][]byte{p}, off)
	if err != nil {
		return min(n, len(p)), w.setError(err)
	}
	return len(p), nil
}

func (w *groupedVectorWriterAt) writeOwned(buf []byte, n int64, off int64, complete func(int, error)) {
	write := &ownedVectorWrite{buf: buf, n: n, off: off, complete: complete}
	if n <= 0 || n > int64(len(buf)) {
		complete(0, fmt.Errorf("invalid owned write length %d for buffer length %d", n, len(buf)))
		return
	}
	if off < 0 {
		complete(0, fmt.Errorf("negative write offset %d", off))
		return
	}

	if n != w.chunkSize || off%w.chunkSize != 0 {
		w.writeRuns([][]*ownedVectorWrite{{write}})
		return
	}

	groupOffset := off / w.groupSize * w.groupSize
	slot := int((off - groupOffset) / w.chunkSize)

	var runs [][]*ownedVectorWrite
	w.mu.Lock()
	if w.closed {
		err := fmt.Errorf("grouped vector writer is closed")
		w.mu.Unlock()
		complete(0, err)
		return
	}
	if w.err != nil {
		err := w.err
		w.mu.Unlock()
		complete(0, err)
		return
	}

	group := w.groups[groupOffset]
	if group == nil {
		group = &vectorWriteGroup{slots: make([]*ownedVectorWrite, w.count)}
		w.groups[groupOffset] = group
	}
	if group.slots[slot] != nil {
		delete(w.groups, groupOffset)
		runs = append(runs, group.runs()...)
		runs = append(runs, []*ownedVectorWrite{write})
	} else {
		group.slots[slot] = write
		group.count++
		if group.count == w.count {
			delete(w.groups, groupOffset)
			runs = append(runs, group.runs()...)
		}
	}
	w.mu.Unlock()

	w.writeRuns(runs)
}

func (w *groupedVectorWriterAt) flushPending() error {
	return w.flushGroups(false)
}

func (w *groupedVectorWriterAt) flush() error {
	return w.flushGroups(true)
}

func (w *groupedVectorWriterAt) flushGroups(closeWriter bool) error {
	w.mu.Lock()
	if closeWriter && w.closed {
		err := w.err
		w.mu.Unlock()
		return err
	}
	if !closeWriter && w.closed {
		err := w.err
		if err == nil {
			err = fmt.Errorf("grouped vector writer is closed")
		}
		w.mu.Unlock()
		return err
	}
	if closeWriter {
		w.closed = true
	}

	offsets := make([]int64, 0, len(w.groups))
	for off := range w.groups {
		offsets = append(offsets, off)
	}
	sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })
	runs := make([][]*ownedVectorWrite, 0, len(offsets))
	for _, off := range offsets {
		runs = append(runs, w.groups[off].runs()...)
	}
	w.groups = map[int64]*vectorWriteGroup{}
	err := w.err
	w.mu.Unlock()

	if err != nil {
		completeRuns(runs, err)
		return err
	}
	return w.writeRuns(runs)
}

func (w *groupedVectorWriterAt) writeRuns(runs [][]*ownedVectorWrite) error {
	for i, run := range runs {
		if err := w.currentError(); err != nil {
			completeRuns(runs[i:], err)
			return err
		}
		if err := w.writeBatch(run); err != nil {
			err = w.setError(err)
			completeRuns(runs[i+1:], err)
			return err
		}
	}
	return nil
}

func (w *groupedVectorWriterAt) writeBatch(writes []*ownedVectorWrite) error {
	if len(writes) == 0 {
		return nil
	}

	vectors := make([][]byte, len(writes))
	expectedOffset := writes[0].off
	for i, write := range writes {
		if write.off != expectedOffset {
			err := fmt.Errorf("non-contiguous vector write at offset %d, expected %d", write.off, expectedOffset)
			completeWrites(writes, err)
			return err
		}

		physicalLen := write.n
		if remainder := physicalLen % w.alignment; remainder != 0 {
			if !w.padTail || i != len(writes)-1 {
				err := fmt.Errorf("write length %d at offset %d is not aligned to %d bytes", write.n, write.off, w.alignment)
				completeWrites(writes, err)
				return err
			}
			physicalLen += w.alignment - remainder
			if physicalLen > int64(cap(write.buf)) {
				err := fmt.Errorf("write buffer capacity %d is smaller than padded length %d", cap(write.buf), physicalLen)
				completeWrites(writes, err)
				return err
			}
			write.buf = write.buf[:physicalLen]
			clear(write.buf[write.n:])
		}
		vectors[i] = write.buf[:physicalLen]
		expectedOffset += physicalLen
	}

	written, err := writeAllVectors(w.w, vectors, writes[0].off)
	if err != nil {
		completeWritesAfterError(writes, vectors, written, err)
		return err
	}
	for _, write := range writes {
		write.complete(int(write.n), nil)
	}
	return nil
}

func completeWritesAfterError(writes []*ownedVectorWrite, vectors [][]byte, written int, err error) {
	failed := len(writes) - 1
	remaining := written
	foundPartial := false
	for i, vector := range vectors {
		if remaining < len(vector) {
			failed = i
			foundPartial = true
			break
		}
		remaining -= len(vector)
	}
	if !foundPartial {
		remaining = len(vectors[failed])
	}

	for i, write := range writes {
		switch {
		case i < failed:
			write.complete(int(write.n), nil)
		case i == failed:
			write.complete(min(remaining, int(write.n)), err)
		default:
			write.complete(0, err)
		}
	}
}

func writeAllVectors(w vectorWriterAt, vectors [][]byte, off int64) (int, error) {
	remaining := append([][]byte(nil), vectors...)
	written := 0
	total := vectorLength(vectors)
	for len(remaining) > 0 {
		n, err := w.writeVectorAt(remaining, off+int64(written))
		remainingLen := vectorLength(remaining)
		if err != nil {
			if n < 0 {
				n = 0
			}
			if n > remainingLen {
				return written, fmt.Errorf("pwritev at offset %d returned invalid count %d with error %v", off+int64(written), n, err)
			}
			written += n
			return written, fmt.Errorf("pwritev at offset %d after %d of %d bytes: %w", off, written, total, err)
		}
		if n < 0 || n > remainingLen {
			return written, fmt.Errorf("pwritev at offset %d returned invalid count %d", off+int64(written), n)
		}
		written += n
		remaining = trimWrittenVectors(remaining, n)
		if n == 0 {
			return written, fmt.Errorf("pwritev at offset %d wrote %d of %d bytes: %w", off, written, total, io.ErrNoProgress)
		}
	}
	return written, nil
}

func vectorLength(vectors [][]byte) int {
	var length int
	for _, vector := range vectors {
		length += len(vector)
	}
	return length
}

func trimWrittenVectors(vectors [][]byte, n int) [][]byte {
	for len(vectors) > 0 && n >= len(vectors[0]) {
		n -= len(vectors[0])
		vectors = vectors[1:]
	}
	if len(vectors) > 0 && n > 0 {
		vectors[0] = vectors[0][n:]
	}
	return vectors
}

func (w *groupedVectorWriterAt) currentError() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.err
}

func (w *groupedVectorWriterAt) setError(err error) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.err == nil {
		w.err = err
	}
	return w.err
}

func completeWrites(writes []*ownedVectorWrite, err error) {
	for _, write := range writes {
		write.complete(0, err)
	}
}

func completeRuns(runs [][]*ownedVectorWrite, err error) {
	for _, run := range runs {
		completeWrites(run, err)
	}
}

func (g *vectorWriteGroup) runs() [][]*ownedVectorWrite {
	var runs [][]*ownedVectorWrite
	for i := 0; i < len(g.slots); {
		for i < len(g.slots) && g.slots[i] == nil {
			i++
		}
		start := i
		for i < len(g.slots) && g.slots[i] != nil {
			i++
		}
		if start < i {
			runs = append(runs, g.slots[start:i])
		}
	}
	return runs
}
