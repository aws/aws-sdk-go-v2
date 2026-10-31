package retry

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// IsErrorThrottle provides the interface of an implementation to determine if
// a error response from an operation is a throttling error.
type IsErrorThrottle interface {
	IsErrorThrottle(error) aws.Ternary
}

// IsErrorThrottles is a collection of checks to determine of the error a
// throttle error. Iterates through the checks and returns the state of
// throttle if any check returns something other than unknown.
type IsErrorThrottles []IsErrorThrottle

// IsErrorThrottle returns if the error is a throttle error if any of the
// checks in the list return a value other than unknown.
func (r IsErrorThrottles) IsErrorThrottle(err error) aws.Ternary {
	for _, re := range r {
		if v := re.IsErrorThrottle(err); v != aws.UnknownTernary {
			return v
		}
	}
	return aws.UnknownTernary
}

// IsErrorThrottleFunc wraps a function with the IsErrorThrottle interface.
type IsErrorThrottleFunc func(error) aws.Ternary

// IsErrorThrottle returns if the error is a throttle error.
func (fn IsErrorThrottleFunc) IsErrorThrottle(err error) aws.Ternary {
	return fn(err)
}

// ThrottleErrorCode determines if an attempt should be retried based on the
// API error code.
type ThrottleErrorCode struct {
	Codes map[string]struct{}
}

// IsErrorThrottle return if the error is a throttle error based on the error
// codes. Returns unknown if the error doesn't have a code or it is unknown.
func (r ThrottleErrorCode) IsErrorThrottle(err error) aws.Ternary {
	var v interface{ ErrorCode() string }

	if !errors.As(err, &v) {
		return aws.UnknownTernary
	}

	_, ok := r.Codes[v.ErrorCode()]
	if !ok {
		return aws.UnknownTernary
	}

	return aws.TrueTernary
}

// DefaultThrottleHTTPStatusCodes is the default set of HTTP status codes the
// SDK considers throttle errors.
var DefaultThrottleHTTPStatusCodes = map[int]struct{}{
	429: {},
}

// ThrottleHTTPStatusCode determines if an attempt failed with a throttle
// error based on the HTTP status code of the response.
//
// An HTTP 429 response is only considered a throttle error when the error did
// not deserialize to a modeled error type. A modeled error is a throttle error
// only if its error code is a throttle error code.
type ThrottleHTTPStatusCode struct {
	Codes map[int]struct{}
}

// IsErrorThrottle returns if the error is a throttle error based on the HTTP
// status code. Returns unknown if the error doesn't have a status code or the
// code is not in the set.
func (r ThrottleHTTPStatusCode) IsErrorThrottle(err error) aws.Ternary {
	var v interface{ HTTPStatusCode() int }

	if !errors.As(err, &v) {
		return aws.UnknownTernary
	}

	code := v.HTTPStatusCode()
	_, ok := r.Codes[code]
	if !ok {
		return aws.UnknownTernary
	}

	if code == http.StatusTooManyRequests && isModeledAPIError(err) {
		return aws.UnknownTernary
	}

	return aws.TrueTernary
}
