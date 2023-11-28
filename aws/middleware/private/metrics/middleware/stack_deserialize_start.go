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

type StackDeserializeStart struct{}

func GetRecordStackDeserializeStartMiddleware() *StackDeserializeStart {
	return &StackDeserializeStart{}
}

func (m *StackDeserializeStart) ID() string {
	return "StackDeserializeStart"
}

func (m *StackDeserializeStart) HandleDeserialize(
	ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler,
) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {

	out, metadata, err = next.HandleDeserialize(ctx, in)

	mctx := metrics.Context(ctx)

	attemptMetrics, attemptErr := mctx.Data().LatestAttempt()

	if attemptErr != nil {
		fmt.Println(err)
	} else {
		attemptMetrics.DeserializeStartTime = sdk.NowTime()
	}

	return out, metadata, err
}
