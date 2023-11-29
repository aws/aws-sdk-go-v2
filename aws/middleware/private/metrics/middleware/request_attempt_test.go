// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics/testutils"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestRegisterAttemptMetricContext_HandleFinalize(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	cases := map[string]struct {
		ConnectionCounter             *metrics.SharedConnectionCounter
		Publisher                     metrics.MetricPublisher
		ProvideResponse               func(metadata middleware.Metadata) *smithyhttp.Response
		Handler                       middleware.FinalizeHandler
		Input                         middleware.FinalizeInput
		ExpectedRequestId             string
		ExpectedExtendedRequestId     string
		ExpectedStatusCode            int
		ExpectedResponseContentLength int64
	}{
		"success": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         &testutils.NoopPublisher{},
			ProvideResponse: func(metadata middleware.Metadata) *smithyhttp.Response {
				res := smithyhttp.Response{}
				res.Response = &http.Response{
					StatusCode:    400,
					ContentLength: 1234,
					Header:        map[string][]string{},
				}
				res.Header.Set(amznRequestIdKey, "reqId")
				res.Header.Set(amznRequestId2Key, "reqId2")
				return &res
			},
			Handler:                       testutils.NoopFinalizeHandler{},
			Input:                         middleware.FinalizeInput{},
			ExpectedRequestId:             "reqId",
			ExpectedExtendedRequestId:     "reqId2",
			ExpectedStatusCode:            400,
			ExpectedResponseContentLength: 1234,
		},
		"noResInfo": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         &testutils.NoopPublisher{},
			ProvideResponse: func(metadata middleware.Metadata) *smithyhttp.Response {
				return nil
			},
			Handler:                       testutils.NoopFinalizeHandler{},
			ExpectedRequestId:             unkAmznReqId,
			ExpectedExtendedRequestId:     unkAmznReqId2,
			ExpectedStatusCode:            -1,
			ExpectedResponseContentLength: -1,
		},
		"noMetadata": {
			ConnectionCounter:             &metrics.SharedConnectionCounter{},
			Publisher:                     &testutils.NoopPublisher{},
			ProvideResponse:               getRawResponse,
			Handler:                       testutils.NoopFinalizeHandler{},
			ExpectedRequestId:             unkAmznReqId,
			ExpectedExtendedRequestId:     unkAmznReqId2,
			ExpectedStatusCode:            -1,
			ExpectedResponseContentLength: -1,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			getRawResponse = c.ProvideResponse

			ctx := metrics.InitMetricContext(context.TODO(), c.ConnectionCounter, c.Publisher)

			mw := GetRegisterAttemptMetricContextMiddleware()

			_, _, _ = mw.HandleFinalize(ctx, c.Input, c.Handler)

			latestAttempt, _ := metrics.Context(ctx).Data().LatestAttempt()
			actualRequestId := latestAttempt.RequestID
			actualExtendedRequestId := latestAttempt.ExtendedRequestID
			actualStatusCode := latestAttempt.StatusCode
			actualResponseContentLength := latestAttempt.ResponseContentLength

			if actualRequestId != c.ExpectedRequestId {
				t.Errorf("Unexpected RequestId, should be '%s' but was '%s'", c.ExpectedRequestId, actualRequestId)
			}
			if actualExtendedRequestId != c.ExpectedExtendedRequestId {
				t.Errorf("Unexpected ExtendedRequestId, should be '%s' but was '%s'", c.ExpectedExtendedRequestId, actualExtendedRequestId)
			}
			if actualStatusCode != c.ExpectedStatusCode {
				t.Errorf("Unexpected StatusCode, should be '%d' but was '%d'", c.ExpectedStatusCode, actualStatusCode)
			}
			if actualResponseContentLength != c.ExpectedResponseContentLength {
				t.Errorf("Unexpected StatusCode, should be '%d' but was '%d'", c.ExpectedResponseContentLength, actualResponseContentLength)
			}
		})
	}

}
