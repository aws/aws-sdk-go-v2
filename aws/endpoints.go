package aws

import (
	"fmt"
)

// Endpoint represents the endpoint a service client should make requests to.
type Endpoint struct {
	// The URL of the endpoint.
	URL string

	// The endpoint partition
	PartitionID string

	// The service name that should be used for signing the requests to the
	// endpoint.
	SigningName string

	// The region that should be used for signing the request to the endpoint.
	SigningRegion string

	// The signing method that should be used for signing the requests to the
	// endpoint.
	SigningMethod string
}

// EndpointNotFound is a sentinel error to indicate that the EndpointResolver implementation was unable
// to resolve an endpoint for the given service and region. Resolvers should use this to indicate that
// a client should fallback and attempt to use it's default resolver to resolve the endpoint.
type EndpointNotFound struct {
	Err error
}

// Error is the error message.
func (e *EndpointNotFound) Error() string {
	return fmt.Sprintf("endpoint not found, %v", e.Err)
}

// Unwrap returns the underlying error.
func (e *EndpointNotFound) Unwrap() error {
	return e.Err
}

// EndpointResolver is an endpoint resolver that can be used to provide or override an endpoint for the given
// service and region. Clients will attempt to use the EndpointResolver first to resolve an endpoint if available.
// If the EndpointResolver returns an EndpointNotFound error, clients will fallback to attempting to resolve the endpoint
// using their default endpoint resolver.
type EndpointResolver interface {
	ResolveEndpoint(service, region string) (Endpoint, error)
}
