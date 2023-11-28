// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
)

type MetricCollection struct {
	cc        *metrics.SharedConnectionCounter
	publisher metrics.MetricPublisher
}

func GetSetupMetricCollectionMiddleware(
	counter *metrics.SharedConnectionCounter, publisher metrics.MetricPublisher,
) *MetricCollection {
	return &MetricCollection{
		cc:        counter,
		publisher: publisher,
	}
}

func (m *MetricCollection) ID() string {
	return "MetricCollection"
}

func (m *MetricCollection) HandleInitialize(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {

	ctx = metrics.InitMetricContext(ctx, m.cc, m.publisher)

	mctx := metrics.Context(ctx)
	metricData := mctx.Data()

	metricData.RequestStartTime = sdk.NowTime()

	out, metadata, err = next.HandleInitialize(ctx, in)

	metricData.RequestEndTime = sdk.NowTime()

	if err == nil {
		metricData.Success = 1
	} else {
		metricData.Success = 0
	}

	metricData.ComputeRequestMetrics()

	publishErr := m.publisher.PostRequestMetrics(metrics.Context(ctx).Data())
	if publishErr != nil {
		fmt.Println("Failed to post request metrics")
	}

	return out, metadata, err
}
