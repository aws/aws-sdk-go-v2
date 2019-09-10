package aws

import "time"

// NoOpRetryer provides a retryer that performs no retries.
// It should be used when we do not want retries to be performed.
type NoOpRetryer struct{}

// MaxRetries returns the number of maximum returns the service will use to make
// an individual API; For NoOpRetryer the MaxRetries will always be zero.
func (d NoOpRetryer) MaxRetries() int {
	return 0
}

// NoOpRetryer should never retry; so ShouldRetry will always return false.
func (d NoOpRetryer) ShouldRetry(_ *Request) bool {
	return false
}

// RetryRules returns the delay duration before retrying this request again;
// since NoOpRetryer does not retry, RetryRules always returns 0.
func (d NoOpRetryer) RetryRules(_ *Request) time.Duration {
	return 0
}
