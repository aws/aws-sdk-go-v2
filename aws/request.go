package aws

import (
	"fmt"
)

const (
	// ErrCodeSerialization is the serialization error code that is received
	// during protocol unmarshaling.
	ErrCodeSerialization = "SerializationError"
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
