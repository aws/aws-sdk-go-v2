// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
)

type EndpointResolutionStart struct{}

func GetRecordEndpointResolutionStartMiddleware() *EndpointResolutionStart {
	return &EndpointResolutionStart{}
}

func (m *EndpointResolutionStart) ID() string {
	return "EndpointResolutionStart"
}

// Deprecated: Endpoint resolution now occurs in Finalize. The ResolveEndpoint
// middleware remains in serialize but is largely a no-op.
func (m *EndpointResolutionStart) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {

	mctx := metrics.Context(ctx)
	mctx.Data().ResolveEndpointStartTime = sdk.NowTime()

	out, metadata, err = next.HandleSerialize(ctx, in)

	return out, metadata, err
}

func (m *EndpointResolutionStart) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	middleware.FinalizeOutput, middleware.Metadata, error,
) {
	mctx := metrics.Context(ctx)
	mctx.Data().ResolveEndpointStartTime = sdk.NowTime()
	return next.HandleFinalize(ctx, in)
}
