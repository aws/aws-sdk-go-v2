package transfermanager

import (
	"context"
	"sync"
	"testing"
)

type mockPartialListener struct{}

func (mockPartialListener) OnObjectTransferStart(context.Context, *ObjectTransferStartEvent)       {}
func (mockPartialListener) OnObjectTransferComplete(context.Context, *ObjectTransferCompleteEvent) {}

func TestProgressListenerRegisterAndCopy(t *testing.T) {
	o := Options{}
	o.ObjectProgressListeners.Register(mockPartialListener{})

	expectIntEq(t, 1, len(o.ObjectProgressListeners.ObjectTransferStart))
	expectIntEq(t, 0, len(o.ObjectProgressListeners.ObjectBytesTransferred))
	expectIntEq(t, 1, len(o.ObjectProgressListeners.ObjectTransferComplete))
	expectIntEq(t, 0, len(o.ObjectProgressListeners.ObjectTransferFailed))

	cp := o.Copy()
	cp.ObjectProgressListeners.Register(mockPartialListener{})

	expectIntEq(t, 1, len(o.ObjectProgressListeners.ObjectTransferStart))
	expectIntEq(t, 0, len(o.ObjectProgressListeners.ObjectBytesTransferred))
	expectIntEq(t, 1, len(o.ObjectProgressListeners.ObjectTransferComplete))
	expectIntEq(t, 0, len(o.ObjectProgressListeners.ObjectTransferFailed))

	expectIntEq(t, 2, len(cp.ObjectProgressListeners.ObjectTransferStart))
	expectIntEq(t, 0, len(cp.ObjectProgressListeners.ObjectBytesTransferred))
	expectIntEq(t, 2, len(cp.ObjectProgressListeners.ObjectTransferComplete))
	expectIntEq(t, 0, len(cp.ObjectProgressListeners.ObjectTransferFailed))
}

func expectIntEq(t *testing.T, expect, actual int) {
	t.Helper()
	if expect != actual {
		t.Errorf("%v != %v", expect, actual)
	}
}

type mockListener struct {
	mu sync.Mutex

	start    []*ObjectTransferStartEvent
	transfer []*ObjectBytesTransferredEvent
	complete []*ObjectTransferCompleteEvent
	failed   []*ObjectTransferFailedEvent
}

func (m *mockListener) OnObjectTransferStart(ctx context.Context, event *ObjectTransferStartEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.start = append(m.start, event)
}

func (m *mockListener) OnObjectBytesTransferred(ctx context.Context, event *ObjectBytesTransferredEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.transfer = append(m.transfer, event)
}

func (m *mockListener) OnObjectTransferComplete(ctx context.Context, event *ObjectTransferCompleteEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.complete = append(m.complete, event)
}

func (m *mockListener) OnObjectTransferFailed(ctx context.Context, event *ObjectTransferFailedEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.failed = append(m.failed, event)
}

func (m *mockListener) expectComplete(t *testing.T, in, out any) {
	t.Helper()

	if len(m.start) != 1 {
		t.Fatalf("transfer start was called %d times instead of once", len(m.start))
	}
	if len(m.complete) != 1 {
		t.Fatalf("transfer complete was called %d times instead of once", len(m.complete))
	}
	if len(m.failed) != 0 {
		t.Fatalf("transfer failed was called on expected completion: %v", m.failed[0].Error)
	}

	start := m.start[0]
	complete := m.complete[0]

	// input/output are all literal equality checks because what we emit in
	// progress AND return should be the same pointer
	if in != start.Input {
		t.Errorf("transfer start: input %v != %v", in, start.Input)
	}
	if in != complete.Input {
		t.Errorf("transfer complete: input %v != %v", in, complete.Input)
	}
	if out != complete.Output {
		t.Errorf("transfer complete: output %v != %v", out, complete.Output)
	}
}

func (m *mockListener) expectFailed(t *testing.T, in any, err error) {
	t.Helper()

	if len(m.start) != 1 {
		t.Fatalf("transfer start was called %d times instead of once", len(m.start))
	}
	if len(m.complete) != 0 {
		t.Fatalf("transfer complete was called on expected failure: %v", m.complete[0])
	}
	if len(m.failed) != 1 {
		t.Fatalf("transfer failed was %d times instead of once", len(m.failed))
	}

	start := m.start[0]
	failed := m.failed[0]

	if in != start.Input {
		t.Errorf("transfer start: input %v != %v", in, start.Input)
	}
	if in != failed.Input {
		t.Errorf("transfer failed: input %v != %v", in, failed.Input)
	}
	if err != failed.Error {
		t.Errorf("transfer complete: output %v != %v", err, failed.Error)
	}
}

func (m *mockListener) expectStartTotalBytes(t *testing.T, expect int64) {
	t.Helper()

	if len(m.start) != 1 {
		t.Fatalf("transfer start was called %d times instead of once", len(m.start))
	}

	start := m.start[0]
	if expect != start.TotalBytes {
		t.Errorf("transfer start: total bytes %v != %v", expect, start.TotalBytes)
	}
}

func (m *mockListener) expectCompleteTotalBytes(t *testing.T, expect int64) {
	t.Helper()

	if len(m.complete) != 1 {
		t.Fatalf("transfer complete was called %d times instead of once", len(m.complete))
	}

	complete := m.complete[0]
	if expect != complete.TotalBytes {
		t.Errorf("transfer complete: total bytes %v != %v", expect, complete.TotalBytes)
	}
}

func (m *mockListener) expectByteTransfers(t *testing.T, expect ...int64) {
	t.Helper()

	if len(m.start) != 1 {
		t.Fatalf("transfer start was called %d times instead of once", len(m.start))
	}
	if len(m.transfer) != len(expect) {
		t.Fatalf("bytes transferred was called %d times instead of expected %d times", len(m.transfer), len(expect))
	}

	for i, ex := range expect {
		if ex != m.transfer[i].BytesTransferred {
			t.Errorf("transfer call %d: byte count %d != %d", i, ex, m.transfer[i].BytesTransferred)
		}
	}
}

type mockDirectoryListener struct {
	mu sync.Mutex

	start    []*ObjectsTransferStartEvent
	transfer []*ObjectsTransferredEvent
	complete []*ObjectsTransferCompleteEvent
	failed   []*ObjectsTransferFailedEvent
}

func (m *mockDirectoryListener) OnObjectsTransferStart(ctx context.Context, event *ObjectsTransferStartEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.start = append(m.start, event)
}

func (m *mockDirectoryListener) OnObjectsTransferred(ctx context.Context, event *ObjectsTransferredEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.transfer = append(m.transfer, event)
}

func (m *mockDirectoryListener) OnObjectsTransferComplete(ctx context.Context, event *ObjectsTransferCompleteEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.complete = append(m.complete, event)
}

func (m *mockDirectoryListener) OnObjectsTransferFailed(ctx context.Context, event *ObjectsTransferFailedEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.failed = append(m.failed, event)
}

func (m *mockDirectoryListener) expectStart(t *testing.T, in any) {
	t.Helper()

	if len(m.start) != 1 {
		t.Fatalf("transfer start was called %d times instead of once", len(m.start))
	}

	start := m.start[0]
	if in != start.Input {
		t.Errorf("transfer start: input %v != %v", in, start.Input)
	}
}

func (m *mockDirectoryListener) expectComplete(t *testing.T, in, out any, expectTotalObjects int64) {
	t.Helper()

	if len(m.start) != 1 {
		t.Fatalf("transfer start was called %d times instead of once", len(m.start))
	}
	if len(m.complete) != 1 {
		t.Fatalf("transfer complete was called %d times instead of once", len(m.complete))
	}
	if len(m.failed) != 0 {
		t.Fatalf("transfer failed was called on expected completion: %v", m.failed[0].Error)
	}

	start := m.start[0]
	complete := m.complete[0]

	// input/output are all literal equality checks because what we emit in
	// progress AND return should be the same pointer
	if in != start.Input {
		t.Errorf("transfer start: input %v != %v", in, start.Input)
	}
	if in != complete.Input {
		t.Errorf("transfer complete: input %v != %v", in, complete.Input)
	}
	if out != complete.Output {
		t.Errorf("transfer complete: output %v != %v", out, complete.Output)
	}
	if expectTotalObjects != complete.ObjectsTransferred {
		t.Errorf("transfer complete: total objects %v != %v", expectTotalObjects, complete.ObjectsTransferred)
	}
}

func (m *mockDirectoryListener) expectFailed(t *testing.T, in any, err error) {
	t.Helper()

	// considering uncertain order of traversing directory, it is possible
	// that an error occurred before any put object call succeeded, so there's
	// no start input check for failure case
	if len(m.complete) != 0 {
		t.Fatalf("transfer complete was called on expected failure: %v", m.complete[0])
	}
	if len(m.failed) != 1 {
		t.Fatalf("transfer failed was %d times instead of once", len(m.failed))
	}

	failed := m.failed[0]

	if in != failed.Input {
		t.Errorf("transfer failed: input %v != %v", in, failed.Input)
	}
	if err != failed.Error {
		t.Errorf("transfer complete: output %v != %v", err, failed.Error)
	}
}

func (m *mockDirectoryListener) expectObjectsTransferred(t *testing.T, expectObjectsCount ...int64) {
	t.Helper()

	if len(m.start) != 1 {
		t.Fatalf("transfer start was called %d times instead of once", len(m.start))
	}
	if len(m.transfer) != len(expectObjectsCount) {
		t.Fatalf("objects transferred was called %d times instead of expected %d times", len(m.transfer), len(expectObjectsCount))
	}

	for i, ex := range expectObjectsCount {
		if ex != m.transfer[i].ObjectsTransferred {
			t.Errorf("transfer call %d: object count %d != %d", i, ex, m.transfer[i].ObjectsTransferred)
		}
	}
}
