// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics/testutils"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
)

func TestGetSetupMetricCollectionMiddleware(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	cases := map[string]struct {
		ConnectionCounter *metrics.SharedConnectionCounter
		Publisher         metrics.MetricPublisher
		Input             middleware.InitializeInput
		Handler           middleware.InitializeHandler
		ExpectedStartTime time.Time
		ExpectedEndTime   time.Time
		ExpectedSuccess   uint8
	}{
		"success": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         &testutils.NoopPublisher{},
			Input:             middleware.InitializeInput{},
			Handler:           testutils.NoopInitializeHandler{},
			ExpectedStartTime: time.Unix(1234, 0),
			ExpectedEndTime:   time.Unix(1234, 0),
			ExpectedSuccess:   1,
		},
		"!success": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         &testutils.NoopPublisher{},
			Input:             middleware.InitializeInput{},
			Handler:           testutils.ErrorInitializeHandler{},
			ExpectedStartTime: time.Unix(1234, 0),
			ExpectedEndTime:   time.Unix(1234, 0),
			ExpectedSuccess:   0,
		},
		"publisherFailure": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         &testutils.ErrorPublisher{},
			Input:             middleware.InitializeInput{},
			Handler:           testutils.NoopInitializeHandler{},
			ExpectedStartTime: time.Unix(1234, 0),
			ExpectedEndTime:   time.Unix(1234, 0),
			ExpectedSuccess:   1,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			ctx := metrics.InitMetricContext(context.TODO(), c.ConnectionCounter, c.Publisher)

			mw := GetSetupMetricCollectionMiddleware(c.ConnectionCounter, c.Publisher)

			_, _, _ = mw.HandleInitialize(ctx, c.Input, c.Handler)

			mctx := metrics.Context(ctx)

			actualRequestStartTime := mctx.Data().RequestStartTime
			actualRequestEndTime := mctx.Data().RequestEndTime

			if actualRequestStartTime != c.ExpectedStartTime {
				t.Errorf("Unexpected RequestStartTime, should be '%s' but was '%s'", c.ExpectedStartTime, actualRequestStartTime)
			}
			if actualRequestEndTime != c.ExpectedEndTime {
				t.Errorf("Unexpected RequestEndTime, should be '%s' but was '%s'", c.ExpectedEndTime, actualRequestEndTime)
			}
			if mctx.Data().Success != c.ExpectedSuccess {
				t.Errorf("Unexpected Success status, should be '%d' but was '%d'", c.ExpectedSuccess, mctx.Data().Success)
			}

		})
	}

}
