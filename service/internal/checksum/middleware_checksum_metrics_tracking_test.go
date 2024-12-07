//go:build go1.16
// +build go1.16

package checksum

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"net/http"
	"testing"
)

func TestRequestChecksumMetricsTracking(t *testing.T) {
	cases := map[string]struct {
		requestChecksumCalculation aws.RequestChecksumCalculation
		reqHeaders                 http.Header
		expectedUserAgentHeader    string
	}{
		//"default": {
		//	requestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
		//	reqHeaders:                 map[string][]string{},
		//	expectedUserAgentHeader:    "m/Z",
		//},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ua := awsmiddleware.NewRequestUserAgent()
			req := smithyhttp.NewStackRequest().(*smithyhttp.Request)
			req.Header = c.reqHeaders
			mw := RequestChecksumMetricsTracking{
				RequestChecksumCalculation: c.requestChecksumCalculation,
				UserAgent:                  ua,
			}
			mw.HandleBuild(context.Background(),
				middleware.BuildInput{Request: req},
				middleware.BuildHandlerFunc(func(ctx context.Context, in middleware.BuildInput) (out middleware.BuildOutput, metadata middleware.Metadata, err error) {
					return
				}))

			ua.HandleBuild(context.Background(), middleware.BuildInput{Request: req},
				middleware.BuildHandlerFunc(func(ctx context.Context, in middleware.BuildInput) (out middleware.BuildOutput, metadata middleware.Metadata, err error) {
					return
				}))

			if e, a := c.expectedUserAgentHeader, req.Header["User-Agent"][0]; e != a {
				t.Errorf("expected user agent header to be %s, got %s", e, a)
			}
		})
	}
}
