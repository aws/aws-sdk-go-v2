package transfermanager

import "context"

// ProgressListeners holds various "transfer progress" hooks that a caller can
// supply to receive progress updates for potentially long-running transfer
// manager operations.
//
// Progress listeners are invoked synchronously within the outer transfer
// operation. Callers SHOULD NOT perform long-lived operations in these hooks,
// such as submitting the progress snapshot to some other network agent.
type ProgressListeners struct {
	ObjectTransferStart    []ObjectTransferStartListener
	ObjectBytesTransferred []ObjectBytesTransferredListener
	ObjectTransferComplete []ObjectTransferCompleteListener
	ObjectTransferFailed   []ObjectTransferFailedListener
}

// Register registers the input with all progress listener hooks that it implements.
//
// If the input does not implement a specific listener, it is a no-op for one
// instance. Callers should generally use compile-time type assertions to
// verify that their implementations satisfy the desired listener interfaces.
func (p *ProgressListeners) Register(v any) {
	if l, ok := v.(ObjectTransferStartListener); ok {
		p.ObjectTransferStart = append(p.ObjectTransferStart, l)
	}
	if l, ok := v.(ObjectBytesTransferredListener); ok {
		p.ObjectBytesTransferred = append(p.ObjectBytesTransferred, l)
	}
	if l, ok := v.(ObjectTransferCompleteListener); ok {
		p.ObjectTransferComplete = append(p.ObjectTransferComplete, l)
	}
	if l, ok := v.(ObjectTransferFailedListener); ok {
		p.ObjectTransferFailed = append(p.ObjectTransferFailed, l)
	}
}

// Copy creates a clone where all hook lists are deep-copied.
func (p *ProgressListeners) Copy() ProgressListeners {
	objectTransferStart := make([]ObjectTransferStartListener, len(p.ObjectTransferStart))
	objectBytesTransferred := make([]ObjectBytesTransferredListener, len(p.ObjectBytesTransferred))
	objectTransferComplete := make([]ObjectTransferCompleteListener, len(p.ObjectTransferComplete))
	objectTransferFailed := make([]ObjectTransferFailedListener, len(p.ObjectTransferFailed))
	copy(objectTransferStart, p.ObjectTransferStart)
	copy(objectBytesTransferred, p.ObjectBytesTransferred)
	copy(objectTransferComplete, p.ObjectTransferComplete)
	copy(objectTransferFailed, p.ObjectTransferFailed)
	return ProgressListeners{
		ObjectTransferStart:    objectTransferStart,
		ObjectBytesTransferred: objectBytesTransferred,
		ObjectTransferComplete: objectTransferComplete,
		ObjectTransferFailed:   objectTransferFailed,
	}
}

// ObjectTransferStartListener is invoked when a single-object transfer begins.
type ObjectTransferStartListener interface {
	OnObjectTransferStart(context.Context, *ObjectTransferStartEvent)
}

// ObjectTransferStartEvent is the event payload for object transfer start.
type ObjectTransferStartEvent struct {
	Request    any
	TotalBytes int64
}

// ObjectBytesTransferredListener is invoked on progress in a single-object
// transfer.
//
// This hook is ALWAYS invoked at least once for an operation, even if the
// operation only does one intermediate transfer (e.g. an object read that does
// not actually need multiple range/part requests).
type ObjectBytesTransferredListener interface {
	OnObjectBytesTransferred(context.Context, *ObjectBytesTransferredEvent)
}

// ObjectBytesTransferredEvent is the event payload for object bytes
// transferred.
type ObjectBytesTransferredEvent struct {
	Input            any
	BytesTransferred int64
	TotalBytes       int64
}

// ObjectTransferCompleteListener is invoked when a single-object transfer
// completes without error.
type ObjectTransferCompleteListener interface {
	OnObjectTransferComplete(context.Context, *ObjectTransferCompleteEvent)
}

// ObjectTransferCompleteEvent is the event payload for object transfer
// complete.
type ObjectTransferCompleteEvent struct {
	Input      any
	Output     any
	TotalBytes int64
}

// ObjectTransferFailedListener is invoked when a single-object transfer fails.
//
// This hook is only invoked for overall operation failure.
type ObjectTransferFailedListener interface {
	OnObjectTransferFailed(context.Context, *ObjectTransferFailedEvent)
}

// ObjectTransferFailedEvent is the event payload for object transfer failure.
type ObjectTransferFailedEvent struct {
	Input            any
	Output           any
	Error            error
	BytesTransferred int64
	TotalBytes       int64
}
