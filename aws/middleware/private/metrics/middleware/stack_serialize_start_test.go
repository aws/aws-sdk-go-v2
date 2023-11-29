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

func TestStartSerializeStart_HandleSerialize(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, &testutils.NoopPublisher{})
	mw := GetRecordStackSerializeStartMiddleware()
	_, _, _ = mw.HandleSerialize(ctx, middleware.SerializeInput{}, testutils.NoopSerializeHandler{})

	actualTime := metrics.Context(ctx).Data().SerializeStartTime
	expectedTime := sdk.NowTime()
	if actualTime != expectedTime {
		t.Errorf("Unexpected SerializeStartTime, should be '%s' but was '%s'", expectedTime, actualTime)
	}
}
