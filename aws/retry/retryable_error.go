package retry

import (
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// IsErrorRetryable provides the interface of an implementation to determine if
// a error as the result of an operation is retryable.
type IsErrorRetryable interface {
	IsErrorRetryable(error) aws.Ternary
}

// IsErrorRetryables is a collection of checks to determine of the error is
// retryable.  Iterates through the checks and returns the state of retryable
// if any check returns something other than unknown.
type IsErrorRetryables []IsErrorRetryable

// IsErrorRetryable returns if the error is retryable if any of the checks in
// the list return a value other than unknown.
func (r IsErrorRetryables) IsErrorRetryable(err error) aws.Ternary {
	for _, re := range r {
		if v := re.IsErrorRetryable(err); v != aws.UnknownTernary {
			return v
		}
	}
	return aws.UnknownTernary
}

// IsErrorRetryableFunc wraps a function with the IsErrorRetryable interface.
type IsErrorRetryableFunc func(error) aws.Ternary

// IsErrorRetryable returns if the error is retryable.
func (fn IsErrorRetryableFunc) IsErrorRetryable(err error) aws.Ternary {
	return fn(err)
}

// RetryableOperationError is an IsErrorRetryable implementation which uses the
// optional interface Retryable on the error value to determine if the error is
// retryable.
type RetryableOperationError struct{}

// IsErrorRetryable returns if the error is retryable if it satisfies the
// Retryable interface, and returns if the attempt should be retried.
func (RetryableOperationError) IsErrorRetryable(err error) aws.Ternary {
	var v interface{ RetryableError() bool }

	if !errors.As(err, &v) {
		return aws.UnknownTernary
	}

	return aws.BoolTernary(v.RetryableError())
}

// NoRetryCanceledError detects if the error was an request canceled error and
// returns if so.
type NoRetryCanceledError struct{}

// IsErrorRetryable returns the error is not retryable if the request was
// canceled.
func (NoRetryCanceledError) IsErrorRetryable(err error) aws.Ternary {
	var v interface{ CanceledError() bool }

	if !errors.As(err, &v) {
		return aws.UnknownTernary
	}

	if v.CanceledError() {
		return aws.FalseTernary
	}
	return aws.UnknownTernary
}

// RetryableConnectionError determines if the underlying error is an HTTP
// connection and returns if it should be retried.
type RetryableConnectionError struct{}

// IsErrorRetryable returns if the error is caused by and HTTP connection
// error, and should be retried.
func (RetryableConnectionError) IsErrorRetryable(err error) aws.Ternary {
	if err == nil {
		return aws.UnknownTernary
	}
	var conErrVal bool

	var conErr interface{ ConnectionError() bool }
	var tempErr interface{ Temporary() bool }
	var urlErr *url.Error
	var netOpErr *net.OpError

	switch {
	case errors.As(err, &conErr):
		conErrVal = conErr.ConnectionError()

	case errors.As(err, &urlErr):
		conErrVal = strings.Contains(urlErr.Error(), "connection refused")

	case errors.As(err, &netOpErr):
		conErrVal = netOpErr.Op == "dial"
		if !conErrVal {
			conErrVal = strings.Contains(netOpErr.Error(), "connection reset")
		}

	case errors.As(err, &tempErr):
		// url.Error and net.OpError implement Temporary but have different
		// expectations for when an error is or is not temporary.
		conErrVal = tempErr.Temporary()

	default:
		if strings.Contains(err.Error(), "connection reset") {
			return aws.TrueTernary
		}
		return aws.UnknownTernary
	}

	return aws.BoolTernary(conErrVal)

}

// RetryableHTTPStatusCode provides a IsErrorRetryable based on HTTP status
// codes.
type RetryableHTTPStatusCode struct {
	Codes map[int]struct{}
}

// IsErrorRetryable return if the passed in error is retryable based on the
// HTTP status code.
func (r RetryableHTTPStatusCode) IsErrorRetryable(err error) aws.Ternary {
	var v interface{ StatusCode() int }

	if !errors.As(err, &v) {
		return aws.UnknownTernary
	}

	_, ok := r.Codes[v.StatusCode()]
	if !ok {
		return aws.UnknownTernary
	}

	return aws.TrueTernary
}

// RetryableErrorCode determines if an attempt should be retried based on the
// API error code.
type RetryableErrorCode struct {
	Codes map[string]struct{}
}

// IsErrorRetryable return if the error is retryable based on the error codes.
// Returns unknown if the error doesn't have a code or it is unknown.
func (r RetryableErrorCode) IsErrorRetryable(err error) aws.Ternary {
	var v interface{ Code() string }

	if !errors.As(err, &v) {
		return aws.UnknownTernary
	}

	_, ok := r.Codes[v.Code()]
	if !ok {
		return aws.UnknownTernary
	}

	return aws.TrueTernary
}
