package http

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
func GetRequestIDMetadata(metadata middleware.Metadata) string {
	if !metadata.Has(requestID{}) {
		return ""
	}

	v, ok := metadata.Get(requestID{}).(string)
	if !ok {
		return ""
	}
	return v
}
