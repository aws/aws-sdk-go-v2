//go:build go1.21
// +build go1.21

package checksum

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"net/http"
	"strings"
	"testing"
)

func TestRequestChecksumMetricsTracking(t *testing.T) {
	cases := map[string]struct {
		requestChecksumCalculation aws.RequestChecksumCalculation
		reqHeaders                 http.Header
		expectedUserAgentHeader    string
	}{
		"default": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
			reqHeaders:                 map[string][]string{},
			expectedUserAgentHeader:    " m/Z",
		},
		"calculate checksum when required": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
			reqHeaders:                 map[string][]string{},
			expectedUserAgentHeader:    " m/a",
		},
		"default with crc32 checksum": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
			reqHeaders: map[string][]string{
				"X-Amz-Checksum-Crc32": {"aa"},
			},
			expectedUserAgentHeader: " m/U,Z",
		},
		"calculate checksum when required with sha256 checksum": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
			reqHeaders: map[string][]string{
				"X-Amz-Checksum-Sha256": {"aa"},
			},
			expectedUserAgentHeader: " m/Y,a",
		},
		"default with crc32c and crc64": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
			reqHeaders: map[string][]string{
				"X-Amz-Checksum-Crc32c":    {"aa"},
				"X-Amz-Checksum-Crc64nvme": {"aa"},
			},
			expectedUserAgentHeader: " m/V,W,Z",
		},
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

			if e, a := c.expectedUserAgentHeader, req.Header["User-Agent"][0]; !strings.Contains(a, e) {
				t.Errorf("expected user agent header to include %s, got %s", e, a)
			}
		})
	}
}

func TestResponseChecksumMetricsTracking(t *testing.T) {
	cases := map[string]struct {
		responseChecksumValidation aws.ResponseChecksumValidation
		expectedUserAgentHeader    string
	}{
		"default": {
			responseChecksumValidation: aws.ResponseChecksumValidationWhenSupported,
			expectedUserAgentHeader:    " m/b",
		},
		"validate checksum when required": {
			responseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
			expectedUserAgentHeader:    " m/c",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ua := awsmiddleware.NewRequestUserAgent()
			req := smithyhttp.NewStackRequest().(*smithyhttp.Request)
			mw := ResponseChecksumMetricsTracking{
				ResponseChecksumValidation: c.responseChecksumValidation,
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

			if e, a := c.expectedUserAgentHeader, req.Header["User-Agent"][0]; !strings.Contains(a, e) {
				t.Errorf("expected user agent header to contain %s, got %s", e, a)
			}
		})
	}
}
