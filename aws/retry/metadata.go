package retry

import "github.com/aws/smithy-go/middleware"

type requestAttemptsKey struct {
}

type requestAttempts struct {

	// Attempts is a slice consisting metadata from all request attempts.
	Attempts []requestAttempt
}

type requestAttempt struct {

	// Response is raw response if received for the request attempt.
	Response interface{}

	// Error is the error if received for the request attempt.
	Error error

	// Retryable denotes if request will be retried.
	Retryable bool

	// Retried indicates if this request was a retried request.
	Retried bool

	// AttemptMetadata denotes existing metadata for the request attempt.
	AttemptMetadata middleware.Metadata
}

// addRequestAttemptMetadata adds request attempts metadata to middleware metadata
func addRequestAttemptMetadata(metadata *middleware.Metadata, v requestAttempts) {
	metadata.Set(requestAttemptsKey{}, v)
}
