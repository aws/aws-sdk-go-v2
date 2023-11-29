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

type EndpointResolutionEnd struct{}

func GetRecordEndpointResolutionEndMiddleware() *EndpointResolutionEnd {
	return &EndpointResolutionEnd{}
}

func (m *EndpointResolutionEnd) ID() string {
	return "EndpointResolutionEnd"
}

func (m *EndpointResolutionEnd) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {

	mctx := metrics.Context(ctx)
	mctx.Data().ResolveEndpointEndTime = sdk.NowTime()

	out, metadata, err = next.HandleSerialize(ctx, in)

	return out, metadata, err
}
