package transfermanager

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"sync"
	"testing"
	"unsafe"
)

type vectorWriteCall struct {
	off      int64
	vectors  [][]byte
	pointers []uintptr
}

type recordingVectorWriterAt struct {
	mu        sync.Mutex
	calls     []vectorWriteCall
	err       error
	alignment int64
}

func (w *recordingVectorWriterAt) writeBufferAlignment() int64 {
	if w.alignment == 0 {
		return 1
	}
	return w.alignment
}

func (w *recordingVectorWriterAt) writeVectorAt(vectors [][]byte, off int64) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	call := vectorWriteCall{
		off:      off,
		vectors:  make([][]byte, len(vectors)),
		pointers: make([]uintptr, len(vectors)),
	}
	var n int
	for i, vector := range vectors {
		call.pointers[i] = uintptr(unsafe.Pointer(unsafe.SliceData(vector)))
		call.vectors[i] = bytes.Clone(vector)
		n += len(vector)
	}
	w.calls = append(w.calls, call)
	if w.err != nil {
		return 0, w.err
	}
	return n, nil
}

func (w *recordingVectorWriterAt) snapshot() []vectorWriteCall {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]vectorWriteCall(nil), w.calls...)
}

type ownedWriteResult struct {
	n   int
	err error
}

func submitOwnedWrite(writer *groupedVectorWriterAt, buf []byte, n int, off int64) <-chan ownedWriteResult {
	done := make(chan ownedWriteResult, 1)
	writer.writeOwned(buf, int64(n), off, func(n int, err error) {
		done <- ownedWriteResult{n: n, err: err}
	})
	return done
}

func requireOwnedWrite(t *testing.T, done <-chan ownedWriteResult, wantN int) {
	t.Helper()
	result := <-done
	if result.err != nil {
		t.Fatalf("owned write: %v", result.err)
	}
	if result.n != wantN {
		t.Fatalf("owned write count = %d, want %d", result.n, wantN)
	}
}

func TestDownloadFileVectorLayout(t *testing.T) {
	if got, want := writeBehindChunkSize, 8*1024*1024; got != want {
		t.Fatalf("write-behind chunk size = %d, want %d", got, want)
	}
	if downloadFileVectorChunkSize != writeBehindChunkSize {
		t.Fatalf("vector chunk size = %d, want write-behind chunk size %d", downloadFileVectorChunkSize, writeBehindChunkSize)
	}
	if got, want := downloadFileVectorCount, 4; got != want {
		t.Fatalf("vector count = %d, want %d", got, want)
	}

	writer := newGroupedVectorWriterAt(&recordingVectorWriterAt{}, downloadFileVectorChunkSize, downloadFileVectorCount, false)
	if got, want := writer.groupSize, int64(32*1024*1024); got != want {
		t.Fatalf("pwritev group size = %d, want %d", got, want)
	}
}

func TestGroupedVectorWriterAtWritesOriginalFourBuffers(t *testing.T) {
	destination := &recordingVectorWriterAt{}
	writer := newGroupedVectorWriterAt(destination, 4, 4, false)

	writes := []struct {
		off  int64
		data string
	}{
		{off: 8, data: "ijkl"},
		{off: 0, data: "abcd"},
		{off: 12, data: "mnop"},
		{off: 4, data: "efgh"},
	}
	pointers := map[int64]uintptr{}
	var completions []<-chan ownedWriteResult
	for i, write := range writes {
		buf := []byte(write.data)
		pointers[write.off] = uintptr(unsafe.Pointer(unsafe.SliceData(buf)))
		completions = append(completions, submitOwnedWrite(writer, buf, len(buf), write.off))
		if got := len(destination.snapshot()); got != i/3 {
			t.Fatalf("calls after write %d = %d", i, got)
		}
	}

	calls := destination.snapshot()
	if len(calls) != 1 {
		t.Fatalf("calls = %d, want 1", len(calls))
	}
	if calls[0].off != 0 {
		t.Fatalf("offset = %d, want 0", calls[0].off)
	}
	if got, want := calls[0].pointers, []uintptr{pointers[0], pointers[4], pointers[8], pointers[12]}; !slices.Equal(got, want) {
		t.Fatalf("pwritev buffer pointers = %v, want original pointers %v", got, want)
	}
	if got := bytes.Join(calls[0].vectors, nil); !bytes.Equal(got, []byte("abcdefghijklmnop")) {
		t.Fatalf("written data = %q", got)
	}
	for _, done := range completions {
		requireOwnedWrite(t, done, 4)
	}
}

func TestGroupedVectorWriterAtRetainsOwnershipUntilPwritevCompletes(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	destination := &blockingVectorWriterAt{started: started, release: release}
	writer := newGroupedVectorWriterAt(destination, 1, 4, false)

	var completions []<-chan ownedWriteResult
	for i := 0; i < 3; i++ {
		completions = append(completions, submitOwnedWrite(writer, []byte{byte(i)}, 1, int64(i)))
	}
	last := make(chan (<-chan ownedWriteResult), 1)
	go func() {
		last <- submitOwnedWrite(writer, []byte{3}, 1, 3)
	}()

	<-started
	for i, done := range completions {
		select {
		case result := <-done:
			t.Fatalf("write %d completed before pwritev returned: %+v", i, result)
		default:
		}
	}
	close(release)
	completions = append(completions, <-last)
	for _, done := range completions {
		requireOwnedWrite(t, done, 1)
	}
}

type blockingVectorWriterAt struct {
	started chan<- struct{}
	release <-chan struct{}
	once    sync.Once
}

func (w *blockingVectorWriterAt) writeVectorAt(vectors [][]byte, _ int64) (int, error) {
	w.once.Do(func() { close(w.started) })
	<-w.release
	return vectorLength(vectors), nil
}

func TestGroupedVectorWriterAtFlushesContiguousPartialGroup(t *testing.T) {
	destination := &recordingVectorWriterAt{}
	writer := newGroupedVectorWriterAt(destination, 4, 4, false)
	first := submitOwnedWrite(writer, []byte("abcd"), 4, 0)
	second := submitOwnedWrite(writer, []byte("efgh"), 4, 4)

	if err := writer.flushPending(); err != nil {
		t.Fatalf("flushPending: %v", err)
	}
	requireOwnedWrite(t, first, 4)
	requireOwnedWrite(t, second, 4)

	calls := destination.snapshot()
	if len(calls) != 1 || len(calls[0].vectors) != 2 {
		t.Fatalf("calls = %+v, want one two-vector call", calls)
	}
}

func TestGroupedVectorWriterAtFlushesGapsSeparately(t *testing.T) {
	destination := &recordingVectorWriterAt{}
	writer := newGroupedVectorWriterAt(destination, 4, 4, false)
	first := submitOwnedWrite(writer, []byte("abcd"), 4, 0)
	second := submitOwnedWrite(writer, []byte("ijkl"), 4, 8)

	if err := writer.flushPending(); err != nil {
		t.Fatalf("flushPending: %v", err)
	}
	requireOwnedWrite(t, first, 4)
	requireOwnedWrite(t, second, 4)

	calls := destination.snapshot()
	if got, want := []int64{calls[0].off, calls[1].off}, []int64{0, 8}; !slices.Equal(got, want) {
		t.Fatalf("pwritev offsets = %v, want %v", got, want)
	}
}

func TestGroupedVectorWriterAtPadsOnlyDirectIOTail(t *testing.T) {
	destination := &recordingVectorWriterAt{alignment: 4}
	writer := newGroupedVectorWriterAt(destination, 4, 4, true)
	buf := []byte{'a', 'b', 0xff, 0xff}
	done := submitOwnedWrite(writer, buf, 2, 0)
	requireOwnedWrite(t, done, 2)

	calls := destination.snapshot()
	if len(calls) != 1 || len(calls[0].vectors) != 1 {
		t.Fatalf("calls = %+v, want one one-vector call", calls)
	}
	if got, want := calls[0].vectors[0], []byte{'a', 'b', 0, 0}; !bytes.Equal(got, want) {
		t.Fatalf("physical tail = %v, want %v", got, want)
	}
}

func TestGroupedVectorWriterAtPropagatesPwritevError(t *testing.T) {
	writeErr := errors.New("write failed")
	destination := &recordingVectorWriterAt{err: writeErr}
	writer := newGroupedVectorWriterAt(destination, 1, 4, false)

	var completions []<-chan ownedWriteResult
	for i := 0; i < 4; i++ {
		completions = append(completions, submitOwnedWrite(writer, []byte{byte(i)}, 1, int64(i)))
	}
	for _, done := range completions {
		result := <-done
		if !errors.Is(result.err, writeErr) {
			t.Fatalf("owned write error = %v, want %v", result.err, writeErr)
		}
	}
	if err := writer.flush(); !errors.Is(err, writeErr) {
		t.Fatalf("flush error = %v, want %v", err, writeErr)
	}
}

func TestGroupedVectorWriterAtFlushesDuplicateBeforeReplacement(t *testing.T) {
	destination := &recordingVectorWriterAt{}
	writer := newGroupedVectorWriterAt(destination, 4, 4, false)
	original := submitOwnedWrite(writer, []byte("old!"), 4, 0)
	replacement := submitOwnedWrite(writer, []byte("new!"), 4, 0)
	requireOwnedWrite(t, original, 4)
	requireOwnedWrite(t, replacement, 4)

	calls := destination.snapshot()
	if len(calls) != 2 {
		t.Fatalf("calls = %d, want 2", len(calls))
	}
	if got := string(calls[0].vectors[0]); got != "old!" {
		t.Fatalf("first write = %q, want old!", got)
	}
	if got := string(calls[1].vectors[0]); got != "new!" {
		t.Fatalf("second write = %q, want new!", got)
	}
}

type shortVectorWriterAt struct {
	offsets []int64
	first   bool
}

func (w *shortVectorWriterAt) writeVectorAt(vectors [][]byte, off int64) (int, error) {
	w.offsets = append(w.offsets, off)
	length := vectorLength(vectors)
	if !w.first {
		w.first = true
		return min(5, length), nil
	}
	return length, nil
}

func TestGroupedVectorWriterAtContinuesShortPwritev(t *testing.T) {
	destination := &shortVectorWriterAt{}
	writer := newGroupedVectorWriterAt(destination, 4, 4, false)

	var completions []<-chan ownedWriteResult
	for i, data := range []string{"abcd", "efgh", "ijkl", "mnop"} {
		completions = append(completions, submitOwnedWrite(writer, []byte(data), 4, int64(4*i)))
	}
	for _, done := range completions {
		requireOwnedWrite(t, done, 4)
	}
	if got, want := destination.offsets, []int64{0, 5}; !slices.Equal(got, want) {
		t.Fatalf("pwritev offsets = %v, want %v", got, want)
	}
}

func TestGroupedVectorWriterAtReportsProgressBeforePwritevError(t *testing.T) {
	writeErr := errors.New("disk full")
	destination := &shortThenErrorVectorWriterAt{err: writeErr}
	writer := newGroupedVectorWriterAt(destination, 4, 4, false)

	completions := make([]<-chan ownedWriteResult, 0, 4)
	for i, data := range []string{"abcd", "efgh", "ijkl", "mnop"} {
		completions = append(completions, submitOwnedWrite(writer, []byte(data), 4, int64(4*i)))
	}
	results := make([]ownedWriteResult, 4)
	for i, done := range completions {
		results[i] = <-done
	}
	if results[0].n != 4 || results[0].err != nil {
		t.Fatalf("first completion = %+v, want 4 bytes and no error", results[0])
	}
	if results[1].n != 1 || !errors.Is(results[1].err, writeErr) {
		t.Fatalf("second completion = %+v, want 1 byte and %v", results[1], writeErr)
	}
	for i := 2; i < len(results); i++ {
		if results[i].n != 0 || !errors.Is(results[i].err, writeErr) {
			t.Fatalf("completion %d = %+v, want 0 bytes and %v", i, results[i], writeErr)
		}
	}
	if got, want := destination.offsets, []int64{0, 5}; !slices.Equal(got, want) {
		t.Fatalf("pwritev offsets = %v, want %v", got, want)
	}
}

type shortThenErrorVectorWriterAt struct {
	offsets []int64
	err     error
	calls   int
}

func (w *shortThenErrorVectorWriterAt) writeVectorAt(vectors [][]byte, off int64) (int, error) {
	w.offsets = append(w.offsets, off)
	w.calls++
	if w.calls == 1 {
		return min(5, vectorLength(vectors)), nil
	}
	return -1, w.err
}

type negativeErrorVectorWriterAt struct {
	err error
}

func (w *negativeErrorVectorWriterAt) writeVectorAt([][]byte, int64) (int, error) {
	return -1, w.err
}

func TestGroupedVectorWriterAtPreservesPwritevErrorForPaddedTail(t *testing.T) {
	writeErr := errors.New("no space left on device")
	writer := newGroupedVectorWriterAt(&negativeErrorVectorWriterAt{err: writeErr}, 4, 4, true)
	buf := []byte{'a', 'b', 0xff, 0xff}
	result := <-submitOwnedWrite(writer, buf, 2, 0)
	if result.n != 0 || !errors.Is(result.err, writeErr) {
		t.Fatalf("owned write result = %+v, want 0 bytes and %v", result, writeErr)
	}
}

type noProgressVectorWriterAt struct{}

func (*noProgressVectorWriterAt) writeVectorAt([][]byte, int64) (int, error) { return 0, nil }

func TestGroupedVectorWriterAtRejectsNoProgress(t *testing.T) {
	writer := newGroupedVectorWriterAt(&noProgressVectorWriterAt{}, 1, 1, false)
	result := <-submitOwnedWrite(writer, []byte{1}, 1, 0)
	if !errors.Is(result.err, io.ErrNoProgress) {
		t.Fatalf("owned write error = %v, want %v", result.err, io.ErrNoProgress)
	}
}
