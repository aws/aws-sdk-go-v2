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

func TestStackDeserializeEnd_HandleDeserializeSuccess(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, &testutils.NoopPublisher{})
	mw := GetRecordStackDeserializeEndMiddleware()

	data := metrics.Context(ctx).Data()

	data.NewAttempt()

	_, _, _ = mw.HandleDeserialize(ctx, middleware.DeserializeInput{}, testutils.NoopDeserializeHandler{})

	attempt, _ := data.LatestAttempt()
	actualTime := attempt.DeserializeEndTime
	expectedTime := sdk.NowTime()
	if actualTime != expectedTime {
		t.Errorf("Unexpected DeserializeEndTime, should be '%s' but was '%s'", expectedTime, actualTime)
	}
}

func TestStackDeserializeEnd_HandleDeserializeNoAttempt(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, &testutils.NoopPublisher{})
	mw := GetRecordStackDeserializeEndMiddleware()

	data := metrics.Context(ctx).Data()

	_, _, _ = mw.HandleDeserialize(ctx, middleware.DeserializeInput{}, testutils.NoopDeserializeHandler{})

	attempt, err := data.LatestAttempt()

	if attempt != nil {
		t.Errorf("Attempt should be nil")
	}
	if err == nil {
		t.Errorf("Err should not be nil")
	}
}
