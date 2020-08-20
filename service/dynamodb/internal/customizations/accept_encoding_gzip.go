package customizations

import "github.com/awslabs/smithy-go/middleware"

// AddAcceptEncodingGzip explicitly adds handling for accept-encoding GZIP
// middleware to the operation stack. This allows checksums to be correctly
// computed without disabling GZIP support.
func AddAcceptEncodingGzip(stack *middleware.Stack) {
	// TODO add middleware to enable GZIP response, but use a custom decoder
	// that will happen AFTER checksum validation.

	// This must happen BEFORE deserialization
}
