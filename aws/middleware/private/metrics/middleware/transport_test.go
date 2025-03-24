package middleware

import (
	"context"
	"github.com/Enflick/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/Enflick/aws-sdk-go-v2/aws/middleware/private/metrics/testutils"
	"github.com/Enflick/aws-sdk-go-v2/internal/sdk"
	"github.com/Enflick/smithy-go/middleware"
	"testing"
	"time"
)

func TestTransportMetrics_HandleSerialize(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, &testutils.NoopPublisher{})

	data := metrics.Context(ctx).Data()

	data.NewAttempt()

	mw := GetTransportMetricsMiddleware()
	_, _, _ = mw.HandleDeserialize(ctx, middleware.DeserializeInput{}, testutils.NoopDeserializeHandler{})

	attempt, _ := data.LatestAttempt()

	actualStartTime := attempt.ServiceCallStart
	expectedStartTime := sdk.NowTime()

	if actualStartTime != expectedStartTime {
		t.Errorf("Unexpected ServiceCallStart, should be '%s' but was '%s'", expectedStartTime, expectedStartTime)
	}

	actualEndTime := attempt.ServiceCallEnd
	expectedEndTime := sdk.NowTime()

	if actualEndTime != expectedEndTime {
		t.Errorf("Unexpected ServiceCallEnd, should be '%s' but was '%s'", expectedEndTime, actualEndTime)
	}

}
