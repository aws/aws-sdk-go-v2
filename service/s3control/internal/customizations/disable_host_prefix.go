package customizations

import (
	"context"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// DisableHostPrefixMiddleware disables host prefix serialization when v2
// endpoint resolution was used.
type DisableHostPrefixMiddleware struct{}

// ID identifies the middleware.
func (*DisableHostPrefixMiddleware) ID() string {
	return "S3ControlDisableHostPrefix"
}

// HandleSerialize controls whether to serialize modeled host prefixes based on
// the type of endpoint resolution used.
func (m *DisableHostPrefixMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, md middleware.Metadata, err error,
) {
	if awsmiddleware.GetRequiresLegacyEndpoints(ctx) {
		return next.HandleSerialize(ctx, in)
	}

	ctx = smithyhttp.DisableEndpointHostPrefix(ctx, true)
	return next.HandleSerialize(ctx, in)
}

// AddDisableHostPrefixMiddleware adds the middleware to the stack.
func AddDisableHostPrefixMiddleware(s *middleware.Stack) error {
	return s.Serialize.Insert(&DisableHostPrefixMiddleware{}, "ResolveEndpointV2", middleware.After)
}
