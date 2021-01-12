package retry

import (
	"fmt"
	"github.com/aws/smithy-go/middleware"
)

// attemptResultsKey is a metadata accessor key to retrieve metadata
// for all request attempts.
type attemptResultsKey struct {
}

// GetAttemptResults retrieves attempts results from middleware metadata.
func GetAttemptResults(metadata middleware.Metadata) (AttemptResults, error) {
	m, ok := metadata.Get(attemptResultsKey{}).(AttemptResults)
	if !ok {
		return AttemptResults{},
			fmt.Errorf("failed to fetch attempt results from metadata")
	}
	return m, nil
}

// AttemptResults represents struct containing metadata returned by all request attempts.
type AttemptResults struct {

	// Results is a slice consisting attempt result from all request attempts.
	// Results are stored in order request attempt is made.
	Results []AttemptResult
}

// AttemptResult represents attempt result returned by a single request attempt.
type AttemptResult struct {

	// Response is raw response if received for the request attempt.
	Response interface{}

	// Err is the error if received for the request attempt.
	Err error

	// Retryable denotes if request may be retried. This states if an
	// error is considered retryable.
	Retryable bool

	// Retried indicates if this request was retried.
	Retried bool

	// ResponseMetadata is any existing metadata passed via the response middlewares.
	ResponseMetadata middleware.Metadata
}

// addAttemptResults adds attempt results to middleware metadata
func addAttemptResults(metadata *middleware.Metadata, v AttemptResults) {
	metadata.Set(attemptResultsKey{}, v)
}
