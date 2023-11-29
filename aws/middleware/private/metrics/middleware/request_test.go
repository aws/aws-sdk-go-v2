// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics/testutils"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestRegisterMetricContext_HandleFinalize(t *testing.T) {

	cases := map[string]struct {
		ConnectionCounter *metrics.SharedConnectionCounter
		Publisher         metrics.MetricPublisher
		ProvideInput      func() middleware.FinalizeInput
		Handler           middleware.FinalizeHandler
		ExpectedCrId      string
	}{
		"success": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         &testutils.NoopPublisher{},
			ProvideInput: func() middleware.FinalizeInput {
				req := smithyhttp.NewStackRequest().(*smithyhttp.Request)
				req.Header.Set(clientRequestIdKey, "crid")
				return middleware.FinalizeInput{Request: req}
			},
			Handler:      testutils.NoopFinalizeHandler{},
			ExpectedCrId: "crid",
		},
		"noCrIdHeader": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         &testutils.NoopPublisher{},
			ProvideInput: func() middleware.FinalizeInput {
				req := smithyhttp.NewStackRequest().(*smithyhttp.Request)
				return middleware.FinalizeInput{Request: req}
			},
			Handler:      testutils.NoopFinalizeHandler{},
			ExpectedCrId: unkClientId,
		},
		"nilRequest": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         &testutils.NoopPublisher{},
			ProvideInput: func() middleware.FinalizeInput {
				return middleware.FinalizeInput{Request: nil}
			},
			Handler:      testutils.NoopFinalizeHandler{},
			ExpectedCrId: unkClientId,
		},
		"wrongRequestType": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         &testutils.NoopPublisher{},
			ProvideInput: func() middleware.FinalizeInput {
				return middleware.FinalizeInput{Request: "nil"}
			},
			Handler:      testutils.NoopFinalizeHandler{},
			ExpectedCrId: unkClientId,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			ctx := metrics.InitMetricContext(context.TODO(), c.ConnectionCounter, c.Publisher)

			mw := GetRegisterRequestMetricContextMiddleware()

			_, _, _ = mw.HandleFinalize(ctx, c.ProvideInput(), c.Handler)

			mctx := metrics.Context(ctx)
			actualCrId := mctx.Data().ClientRequestID

			if actualCrId != c.ExpectedCrId {
				t.Errorf("Unexpected ClientRequestId, should be '%s' but was '%s'", c.ExpectedCrId, actualCrId)
			}

		})
	}

}
