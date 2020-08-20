package aws

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
