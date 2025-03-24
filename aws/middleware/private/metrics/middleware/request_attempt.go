// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"

	awsmiddleware "github.com/Enflick/aws-sdk-go-v2/aws/middleware"
	"github.com/Enflick/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/Enflick/smithy-go/middleware"
	"github.com/Enflick/smithy-go/transport/http"
)

const (
	amznRequestIdKey  = "X-Amz-Request-Id"
	amznRequestId2Key = "X-Amz-Id-2"
	unkAmznReqId      = "unk"
	unkAmznReqId2     = "unk"
)

type RegisterAttemptMetricContext struct{}

func GetRegisterAttemptMetricContextMiddleware() *RegisterAttemptMetricContext {
	return &RegisterAttemptMetricContext{}
}

func (m *RegisterAttemptMetricContext) ID() string {
	return "RegisterAttemptMetricContext"
}

var getRawResponse = func(metadata middleware.Metadata) *http.Response {
	switch res := awsmiddleware.GetRawResponse(metadata).(type) {
	case *http.Response:
		return res
	default:
		return nil
	}
}

func (m *RegisterAttemptMetricContext) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {

	mctx := metrics.Context(ctx)
	mctx.Data().NewAttempt()

	out, metadata, err = next.HandleFinalize(ctx, in)

	res := getRawResponse(metadata)

	attemptMetrics, _ := mctx.Data().LatestAttempt()

	if res != nil {
		attemptMetrics.RequestID = res.Header.Get(amznRequestIdKey)
		attemptMetrics.ExtendedRequestID = res.Header.Get(amznRequestId2Key)
		attemptMetrics.StatusCode = res.StatusCode
		attemptMetrics.ResponseContentLength = res.ContentLength
	} else {
		attemptMetrics.RequestID = unkAmznReqId
		attemptMetrics.ExtendedRequestID = unkAmznReqId2
		attemptMetrics.StatusCode = -1
		attemptMetrics.ResponseContentLength = -1
	}

	return out, metadata, err
}
