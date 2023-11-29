// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics/testutils"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
	"testing"
	"time"
)

func TestEndpointResolutionStart_HandleSerialize_Success(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, &testutils.NoopPublisher{})
	mw := GetRecordEndpointResolutionStartMiddleware()
	_, _, _ = mw.HandleSerialize(ctx, middleware.SerializeInput{}, testutils.NoopSerializeHandler{})

	actualTime := metrics.Context(ctx).Data().ResolveEndpointStartTime
	expectedTime := sdk.NowTime()
	if actualTime != expectedTime {
		t.Errorf("Unexpected ResolveEndpointStartTime, should be '%s' but was '%s'", expectedTime, actualTime)
	}
}

func TestEndpointResolutionStart_HandleSerialize_Error(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	ctx := metrics.InitMetricContext(context.TODO(), &metrics.SharedConnectionCounter{}, &testutils.NoopPublisher{})

	mw := GetRecordEndpointResolutionStartMiddleware()
	_, _, _ = mw.HandleSerialize(ctx, middleware.SerializeInput{}, testutils.NoopSerializeHandler{})

	actualTime := metrics.Context(ctx).Data().ResolveEndpointStartTime
	expectedTime := sdk.NowTime()
	if actualTime != expectedTime {
		t.Errorf("Unexpected ResolveEndpointStartTime, should be '%s' but was '%s'", expectedTime, actualTime)
	}
}
