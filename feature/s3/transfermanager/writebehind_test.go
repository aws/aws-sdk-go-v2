package transfermanager

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"unsafe"
)

type alignmentCheckingWriterAt struct {
	alignment int64
}

func (w *alignmentCheckingWriterAt) writeBufferAlignment() int64 {
	return w.alignment
}

func (w *alignmentCheckingWriterAt) WriteAt(p []byte, off int64) (int, error) {
	if uintptr(unsafe.Pointer(unsafe.SliceData(p)))%uintptr(w.alignment) != 0 {
		return 0, fmt.Errorf("buffer is not aligned to %d", w.alignment)
	}
	if off%w.alignment != 0 {
		return 0, fmt.Errorf("offset is not aligned to %d", w.alignment)
	}
	return len(p), nil
}

type scalingBlockingWriterAt struct {
	release <-chan struct{}
}

func (w *scalingBlockingWriterAt) WriteAt(p []byte, _ int64) (int, error) {
	<-w.release
	return len(p), nil
}

func TestWriteBehindUsesFixedChunkSize(t *testing.T) {
	writer := newWriteBehindWriterAt(&alignmentCheckingWriterAt{alignment: 1}, 0)
	if got, want := writer.chunkSize, int64(8*1024*1024); got != want {
		t.Fatalf("chunk size = %d, want %d", got, want)
	}
	if err := writer.drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}
}

func TestWriteBehindUsesDestinationBufferAlignment(t *testing.T) {
	destination := &alignmentCheckingWriterAt{alignment: directIOAlignment}
	writer := newWriteBehindWriterAt(destination, directIOAlignment)

	if _, err := writer.WriteAt(make([]byte, directIOAlignment), 0); err != nil {
		t.Fatalf("WriteAt: %v", err)
	}
	if err := writer.drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}
}

func TestWriteBehindDefaultWorkerPolicy(t *testing.T) {
	if got, want := writeBehindInitialWorkers, 16; got != want {
		t.Fatalf("initial workers = %d, want %d", got, want)
	}
	if got, want := writeBehindWorkerBatch, 16; got != want {
		t.Fatalf("worker batch = %d, want %d", got, want)
	}
	if got, want := writeBehindMaxWorkers, 64; got != want {
		t.Fatalf("max workers = %d, want %d", got, want)
	}
}

func TestWriteBehindStartsInitialWorkers(t *testing.T) {
	writer := newWriteBehindWriterAt(&alignmentCheckingWriterAt{alignment: 1}, 1)
	if got := writer.workerCount.Load(); got != writeBehindInitialWorkers {
		t.Fatalf("worker count = %d, want %d", got, writeBehindInitialWorkers)
	}

	if err := writer.drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}
	if got := writer.workerCount.Load(); got != 0 {
		t.Fatalf("worker count after drain = %d, want 0", got)
	}
}

func TestWriteBehindScalesWorkersWhileQueueRemainsFull(t *testing.T) {
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseWrites := func() { releaseOnce.Do(func() { close(release) }) }

	writer := newWriteBehindWriterAtWithConfig(&scalingBlockingWriterAt{release: release}, 1, writeBehindWorkerConfig{
		queueDepth:         4,
		initialWorkers:     2,
		workerBatch:        2,
		maxWorkers:         6,
		scaleCheckInterval: time.Millisecond,
		saturationDuration: 3 * time.Millisecond,
	})
	t.Cleanup(func() {
		releaseWrites()
		_ = writer.drain()
	})

	enqueued := make(chan error, 1)
	go func() {
		for i := 0; i < 32; i++ {
			buf := writer.getBuffer()
			if _, err := writer.enqueue(buf, 1, int64(i)); err != nil {
				enqueued <- err
				return
			}
		}
		enqueued <- nil
	}()

	waitForWorkerCount(t, writer, 6)
	time.Sleep(10 * time.Millisecond)
	if got := writer.workerCount.Load(); got != 6 {
		t.Fatalf("worker count = %d, want capped count 6", got)
	}

	releaseWrites()
	select {
	case err := <-enqueued:
		if err != nil {
			t.Fatalf("enqueue: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("enqueue did not finish after writes were released")
	}

	if err := writer.drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}
	if got := writer.workerCount.Load(); got != 0 {
		t.Fatalf("worker count after drain = %d, want 0", got)
	}
}

func TestWriteBehindOwnedWriterBatchesOriginalBuffers(t *testing.T) {
	destination := &recordingVectorWriterAt{}
	grouped := newGroupedVectorWriterAt(destination, 4, 4, false)
	writer := newWriteBehindWriterAtWithConfig(grouped, 4, testWriteBehindConfig(4))

	pointers := map[int64]uintptr{}
	var lastSeq uint64
	for _, off := range []int64{8, 0, 12, 4} {
		buf := writer.getBuffer()
		copy(buf, []byte{byte(off), byte(off + 1), byte(off + 2), byte(off + 3)})
		pointers[off] = uintptr(unsafe.Pointer(unsafe.SliceData(buf)))
		seq, err := writer.enqueue(buf, 4, off)
		if err != nil {
			t.Fatalf("enqueue at %d: %v", off, err)
		}
		lastSeq = seq
	}
	if err := writer.waitThrough(lastSeq); err != nil {
		t.Fatalf("waitThrough: %v", err)
	}
	if err := writer.drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}

	calls := destination.snapshot()
	if len(calls) != 1 {
		t.Fatalf("calls = %d, want 1", len(calls))
	}
	got := calls[0].pointers
	want := []uintptr{pointers[0], pointers[4], pointers[8], pointers[12]}
	if len(got) != len(want) {
		t.Fatalf("pwritev vector count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("pwritev pointer %d = %#x, want original buffer %#x", i, got[i], want[i])
		}
	}
}

func TestWriteBehindOwnedWriterCompletesOnlyAfterPwritev(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	grouped := newGroupedVectorWriterAt(&blockingVectorWriterAt{started: started, release: release}, 1, 4, false)
	writer := newWriteBehindWriterAtWithConfig(grouped, 1, testWriteBehindConfig(4))

	var lastSeq uint64
	for i := 0; i < 4; i++ {
		buf := writer.getBuffer()
		buf[0] = byte(i)
		seq, err := writer.enqueue(buf, 1, int64(i))
		if err != nil {
			t.Fatalf("enqueue %d: %v", i, err)
		}
		lastSeq = seq
	}

	<-started
	writer.completeMu.Lock()
	completed := writer.completedThrough
	writer.completeMu.Unlock()
	if completed != 0 {
		close(release)
		t.Fatalf("completed through = %d before pwritev returned, want 0", completed)
	}

	close(release)
	if err := writer.waitThrough(lastSeq); err != nil {
		t.Fatalf("waitThrough: %v", err)
	}
	if err := writer.drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}
}

func TestWriteBehindWaitThroughFlushesPartialOwnedGroup(t *testing.T) {
	destination := &recordingVectorWriterAt{}
	grouped := newGroupedVectorWriterAt(destination, 4, 4, false)
	writer := newWriteBehindWriterAtWithConfig(grouped, 4, testWriteBehindConfig(1))

	buf := writer.getBuffer()
	copy(buf, "data")
	seq, err := writer.enqueue(buf, 4, 0)
	if err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if err := writer.waitThrough(seq); err != nil {
		t.Fatalf("waitThrough: %v", err)
	}
	if got := len(destination.snapshot()); got != 1 {
		t.Fatalf("pwritev calls = %d after barrier, want 1", got)
	}
	if err := writer.drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}
}

func TestWriteBehindSynchronousWriteFlushesOwnedGroup(t *testing.T) {
	destination := &recordingVectorWriterAt{}
	grouped := newGroupedVectorWriterAt(destination, 4, 4, false)
	writer := newWriteBehindWriterAtWithConfig(grouped, 4, testWriteBehindConfig(1))

	if n, err := writer.WriteAt([]byte("data"), 0); err != nil {
		t.Fatalf("WriteAt: %v", err)
	} else if n != 4 {
		t.Fatalf("WriteAt count = %d, want 4", n)
	}
	if got := len(destination.snapshot()); got != 1 {
		t.Fatalf("pwritev calls = %d, want 1", got)
	}
	if err := writer.drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}
}

func testWriteBehindConfig(workers int) writeBehindWorkerConfig {
	return writeBehindWorkerConfig{
		queueDepth:         8,
		initialWorkers:     workers,
		workerBatch:        1,
		maxWorkers:         workers,
		scaleCheckInterval: time.Hour,
		saturationDuration: time.Hour,
	}
}

func waitForWorkerCount(t *testing.T, writer *writeBehindWriterAt, want int64) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if writer.workerCount.Load() == want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("worker count = %d, want %d", writer.workerCount.Load(), want)
}
