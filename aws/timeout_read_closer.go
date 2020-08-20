package aws

import (
	"fmt"
	"io"
	"time"
)

type readResult struct {
	n   int
	err error
}

// ResponseTimeoutError is an error when the reads from the response are
// delayed longer than the timeout the read was configured for.
type ResponseTimeoutError struct {
	TimeoutDur time.Duration
}

// Timeout returns that the error is was caused by a timeout, and can be
// retried.
func (*ResponseTimeoutError) Timeout() bool { return true }

func (e *ResponseTimeoutError) Error() string {
	return fmt.Sprintf("read on body reach timeout limit, %v", e.TimeoutDur)
}

// timeoutReadCloser will handle body reads that take too long.
// We will return a ErrReadTimeout error if a timeout occurs.
type timeoutReadCloser struct {
	reader   io.ReadCloser
	duration time.Duration
}

// Read will spin off a goroutine to call the reader's Read method. We will
// select on the timer's channel or the read's channel. Whoever completes first
// will be returned.
func (r *timeoutReadCloser) Read(b []byte) (int, error) {
	timer := time.NewTimer(r.duration)
	c := make(chan readResult, 1)

	go func() {
		n, err := r.reader.Read(b)
		timer.Stop()
		c <- readResult{n: n, err: err}
	}()

	select {
	case data := <-c:
		return data.n, data.err
	case <-timer.C:
		return 0, &ResponseTimeoutError{TimeoutDur: r.duration}
	}
}

func (r *timeoutReadCloser) Close() error {
	return r.reader.Close()
}
