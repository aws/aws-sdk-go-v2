package transfermanager

import (
	"context"
	"testing"
)

type mockListener struct{}

func (mockListener) OnObjectTransferStart(context.Context, *ObjectTransferStartEvent)       {}
func (mockListener) OnObjectTransferComplete(context.Context, *ObjectTransferCompleteEvent) {}

func TestProgressListenerRegisterAndCopy(t *testing.T) {
	o := Options{}
	o.ProgressListeners.Register(mockListener{})

	expectIntEq(t, 1, len(o.ProgressListeners.ObjectTransferStart))
	expectIntEq(t, 0, len(o.ProgressListeners.ObjectBytesTransferred))
	expectIntEq(t, 1, len(o.ProgressListeners.ObjectTransferComplete))
	expectIntEq(t, 0, len(o.ProgressListeners.ObjectTransferFailed))

	cp := o.Copy()
	cp.ProgressListeners.Register(mockListener{})

	expectIntEq(t, 1, len(o.ProgressListeners.ObjectTransferStart))
	expectIntEq(t, 0, len(o.ProgressListeners.ObjectBytesTransferred))
	expectIntEq(t, 1, len(o.ProgressListeners.ObjectTransferComplete))
	expectIntEq(t, 0, len(o.ProgressListeners.ObjectTransferFailed))

	expectIntEq(t, 2, len(cp.ProgressListeners.ObjectTransferStart))
	expectIntEq(t, 0, len(cp.ProgressListeners.ObjectBytesTransferred))
	expectIntEq(t, 2, len(cp.ProgressListeners.ObjectTransferComplete))
	expectIntEq(t, 0, len(cp.ProgressListeners.ObjectTransferFailed))
}

func expectIntEq(t *testing.T, expect, actual int) {
	t.Helper()
	if expect != actual {
		t.Errorf("%v != %v", expect, actual)
	}
}
