package id

const (
	// RegisterServiceMetadata is the slot ID for registering service metadata
	RegisterServiceMetadata = "RegisterServiceMetadata"
	// ResolveEndpoint is the slot ID for endpoint resolver middleware.
	ResolveEndpoint = "ResolveEndpoint"
	// ClientRequestID is the slot ID for middleware that generate/set a client-side request-id
	ClientRequestID = "ClientRequestID"
	// ComputePayloadHash is the slot ID for middleware that compute hashes of the transport body payload
	ComputePayloadHash = "ComputePayloadHash"
	// Retry is the slot ID for middleware that handles retrying on transport requests.
	Retry = "Retry"
	// Signing is the slot ID for middleware that handles signing a transport request.
	Signing = "Signing"
	// ResponseErrorWrapper is the slot ID for middleware that wraps or decorates middleware errors.
	ResponseErrorWrapper = "ResponseErrorWrapper"
	// RequestIDRetriever is the slot ID for middleware that retrieves the response request-id.
	RequestIDRetriever = "RequestIDRetriever"
	// UserAgent is the slot ID for middleware the configures the client user-agent for the given transport.
	UserAgent = "UserAgent"
)
