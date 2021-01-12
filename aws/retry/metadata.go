package retry

import (
	"fmt"
	"github.com/aws/smithy-go/middleware"
)

// requestAttemptsKey is a metadata accessor key to
// retrieve metadata for all request attempts.
type requestAttemptsKey struct {
}

// GetRequestAttemptsMetadata retrieves request attempts metadata from middleware metadata.
func GetRequestAttemptsMetadata(metadata middleware.Metadata) (RequestAttemptsMetadata, error) {
	m, ok := metadata.Get(requestAttemptsKey{}).(RequestAttemptsMetadata)
	if !ok {
		return RequestAttemptsMetadata{},
			fmt.Errorf("failed to fetch request attempts metadata")
	}
	return m, nil
}

// RequestAttemptsMetadata represents struct containing
// metadata returned by all request attempts.
type RequestAttemptsMetadata struct {

	// Attempts is a slice consisting metadata from all request attempts.
	// Attempts are stored in last in first order i.e. the last attempt
	// would be at the top.
	Attempts []RequestAttemptMetadata
}

// RequestAttemptMetadata represents metadata returned by a request attempt.
type RequestAttemptMetadata struct {

	// Response is raw response if received for the request attempt.
	Response interface{}

	// Err is the error if received for the request attempt.
	Err error

	// Retryable denotes if request may be retried. This states if an
	// error was retryable.
	Retryable bool

	// Retried indicates if this request was retried.
	Retried bool

	// AttemptMetadata denotes existing metadata for the request attempt.
	AttemptMetadata middleware.Metadata
}

// addRequestAttemptMetadata adds request attempts metadata to middleware metadata
func addRequestAttemptMetadata(metadata *middleware.Metadata, v RequestAttemptsMetadata) {
	metadata.Set(requestAttemptsKey{}, v)
}
