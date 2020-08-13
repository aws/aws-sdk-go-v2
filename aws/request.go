package aws

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	// ErrCodeSerialization is the serialization error code that is received
	// during protocol unmarshaling.
	ErrCodeSerialization = "SerializationError"

	// ErrCodeRead is an error that is returned during HTTP reads.
	ErrCodeRead = "ReadError"
)

// MaxAttemptsError provides the error when the maximum number of attempts have
// been exceeded.
type MaxAttemptsError struct {
	Attempt int
	Err     error
}

func (e *MaxAttemptsError) Error() string {
	return fmt.Sprintf("exceeded maximum number of attempts, %d, %v", e.Attempt, e.Err)
}

// Unwrap returns the nested error causing the max attempts error. Provides the
// implementation for errors.Is and errors.As to unwrap nested errors.
func (e *MaxAttemptsError) Unwrap() error {
	return e.Err
}

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

// SanitizeHostForHeader removes default port from host and updates request.Host
func SanitizeHostForHeader(r *http.Request) {
	host := getHost(r)
	port := portOnly(host)
	if port != "" && isDefaultPort(r.URL.Scheme, port) {
		r.Host = stripPort(host)
	}
}

// Returns host from request
func getHost(r *http.Request) string {
	if r.Host != "" {
		return r.Host
	}

	return r.URL.Host
}

// Hostname returns u.Host, without any port number.
//
// If Host is an IPv6 literal with a port number, Hostname returns the
// IPv6 literal without the square brackets. IPv6 literals may include
// a zone identifier.
//
// Copied from the Go 1.8 standard library (net/url)
func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

// Port returns the port part of u.Host, without the leading colon.
// If u.Host doesn't contain a port, Port returns an empty string.
//
// Copied from the Go 1.8 standard library (net/url)
func portOnly(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return ""
	}
	if i := strings.Index(hostport, "]:"); i != -1 {
		return hostport[i+len("]:"):]
	}
	if strings.Contains(hostport, "]") {
		return ""
	}
	return hostport[colon+len(":"):]
}

// Returns true if the specified URI is using the standard port
// (i.e. port 80 for HTTP URIs or 443 for HTTPS URIs)
func isDefaultPort(scheme, port string) bool {
	if port == "" {
		return true
	}

	lowerCaseScheme := strings.ToLower(scheme)
	if (lowerCaseScheme == "http" && port == "80") || (lowerCaseScheme == "https" && port == "443") {
		return true
	}

	return false
}
