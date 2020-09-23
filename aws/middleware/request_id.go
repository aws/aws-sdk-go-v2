package middleware

import (
	"github.com/awslabs/smithy-go/middleware"
)

// requestID is used to retrieve request id from response metadata
type requestID struct {
}

// SetRequestIDMetadata sets the provided request id over middleware metadata
func SetRequestIDMetadata(metadata *middleware.Metadata, id string) {
	metadata.Set(requestID{}, id)
}

// GetRequestIDMetadata retrieves the request id from middleware metadata
// returns string and bool indicating value of request id, whether request id was set.
func GetRequestIDMetadata(metadata middleware.Metadata) (string, bool) {
	if !metadata.Has(requestID{}) {
		return "", false
	}

	v, ok := metadata.Get(requestID{}).(string)
	if !ok {
		return "", true
	}
	return v, true
}
