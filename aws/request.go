package aws

import (
	"fmt"
)

// RequestCanceledError is the error that will be returned by an API request
// that was canceled. Requests given a Context may return this error when
// canceled.
type RequestCanceledError struct {
	Err error
}

// CanceledError returns true to satisfy interfaces checking for canceled errors.
func (*RequestCanceledError) CanceledError() bool { return true }

// Unwrap returns the underlying error, if there was one.
func (e *RequestCanceledError) Unwrap() error {
	return e.Err
}
func (e *RequestCanceledError) Error() string {
	return fmt.Sprintf("request canceled, %v", e.Err)
}

// RequestSendError provides a generic request transport error.
type RequestSendError struct {
	Err error
}

// ConnectionError return that the error is related to not being able to send
// the request.
func (e *RequestSendError) ConnectionError() bool {
	return true
}

// Unwrap returns the underlying error, if there was one.
func (e *RequestSendError) Unwrap() error {
	return e.Err
}

func (e *RequestSendError) Error() string {
	return fmt.Sprintf("request send failed, %v", e.Err)
}
