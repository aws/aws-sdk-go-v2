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
	"runtime"
	"strings"
	"testing"
)

var expectedAgent = aws.SDKName + "/" + aws.SDKVersion +
	" ua/2.1" +
	" os/" + getNormalizedOSName() +
	" lang/go#" + strings.Map(rules, languageVersion) + // normalize as the user-agent builder will
	" md/GOOS#" + runtime.GOOS +
	" md/GOARCH#" + runtime.GOARCH

var languageVersion = strings.TrimPrefix(runtime.Version(), "go")

func getNormalizedOSName() (os string) {
	switch runtime.GOOS {
	case "android":
		os = "android"
	case "linux":
		os = "linux"
	case "windows":
		os = "windows"
	case "darwin":
		os = "macos"
	case "ios":
		os = "ios"
	default:
		os = "other"
	}
	return os
}

var validChars = map[rune]bool{
	'!': true, '#': true, '$': true, '%': true, '&': true, '\'': true, '*': true, '+': true,
	'-': true, '.': true, '^': true, '_': true, '`': true, '|': true, '~': true,
}

func rules(r rune) rune {
	switch {
	case r >= '0' && r <= '9':
		return r
	case r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z':
		return r
	case validChars[r]:
		return r
	default:
		return '-'
	}
}

func TestRequestChecksumMetricsTracking(t *testing.T) {
	cases := map[string]struct {
		requestChecksumCalculation aws.RequestChecksumCalculation
		reqHeaders                 http.Header
		expectedUserAgentHeader    string
	}{
		"default": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
			reqHeaders:                 map[string][]string{},
			expectedUserAgentHeader:    expectedAgent + " m/Z",
		},
		"calculate checksum when required": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
			reqHeaders:                 map[string][]string{},
			expectedUserAgentHeader:    expectedAgent + " m/a",
		},
		"default with crc32 checksum": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
			reqHeaders: map[string][]string{
				"X-Amz-Checksum-Crc32": {"aa"},
			},
			expectedUserAgentHeader: expectedAgent + " m/U,Z",
		},
		"calculate checksum when required with sha256 checksum": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
			reqHeaders: map[string][]string{
				"X-Amz-Checksum-Sha256": {"aa"},
			},
			expectedUserAgentHeader: expectedAgent + " m/Y,a",
		},
		"default with crc32c and crc64": {
			requestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
			reqHeaders: map[string][]string{
				"X-Amz-Checksum-Crc32c":    {"aa"},
				"X-Amz-Checksum-Crc64nvme": {"aa"},
			},
			expectedUserAgentHeader: expectedAgent + " m/V,W,Z",
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

			if e, a := c.expectedUserAgentHeader, req.Header["User-Agent"][0]; e != a {
				t.Errorf("expected user agent header to be %s, got %s", e, a)
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
			expectedUserAgentHeader:    expectedAgent + " m/b",
		},
		"validate checksum when required": {
			responseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
			expectedUserAgentHeader:    expectedAgent + " m/c",
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

			if e, a := c.expectedUserAgentHeader, req.Header["User-Agent"][0]; e != a {
				t.Errorf("expected user agent header to be %s, got %s", e, a)
			}
		})
	}
}
