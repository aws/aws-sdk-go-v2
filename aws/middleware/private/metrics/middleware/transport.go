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

type TransportMetrics struct{}

func GetTransportMetricsMiddleware() *TransportMetrics {
	return &TransportMetrics{}
}

func (m *TransportMetrics) ID() string {
	return "TransportMetrics"
}

func (m *TransportMetrics) HandleDeserialize(
	ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler,
) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, attemptErr error,
) {

	mctx := metrics.Context(ctx)

	if attempt, e := mctx.Data().LatestAttempt(); e == nil {
		attempt.ServiceCallStart = sdk.NowTime()
		mctx.ConnectionCounter().AddActiveRequest()
	}

	out, metadata, err := next.HandleDeserialize(ctx, in)

	if attempt, e := mctx.Data().LatestAttempt(); e == nil {
		attempt.ServiceCallEnd = sdk.NowTime()
		mctx.ConnectionCounter().RemoveActiveRequest()
	}

	return out, metadata, err

}
