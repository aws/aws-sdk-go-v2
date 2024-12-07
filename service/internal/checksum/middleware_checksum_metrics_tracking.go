package checksum

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

var supportedChecksumFeatures = map[Algorithm]awsmiddleware.UserAgentFeature{
	AlgorithmCRC32:     awsmiddleware.UserAgentFeatureRequestChecksumCRC32,
	AlgorithmCRC32C:    awsmiddleware.UserAgentFeatureRequestChecksumCRC32C,
	AlgorithmSHA1:      awsmiddleware.UserAgentFeatureRequestChecksumSHA1,
	AlgorithmSHA256:    awsmiddleware.UserAgentFeatureRequestChecksumSHA256,
	AlgorithmCRC64NVME: awsmiddleware.UserAgentFeatureRequestChecksumCRC64,
}

type RequestChecksumMetricsTracking struct {
	RequestChecksumCalculation aws.RequestChecksumCalculation
	UserAgent                  *awsmiddleware.RequestUserAgent
}

func (m *RequestChecksumMetricsTracking) ID() string {
	return "AWSChecksum:RequestMetricsTracking"
}

func (m *RequestChecksumMetricsTracking) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	switch m.RequestChecksumCalculation {
	case aws.RequestChecksumCalculationWhenSupported:
		m.UserAgent.AddUserAgentFeature(awsmiddleware.UserAgentFeatureRequestChecksumWhenSupported)
	case aws.RequestChecksumCalculationWhenRequired:
		m.UserAgent.AddUserAgentFeature(awsmiddleware.UserAgentFeatureRequestChecksumWhenRequired)
	}

	for algo, feat := range supportedChecksumFeatures {
		checksumHeader := AlgorithmHTTPHeader(algo)
		if checksum := req.Header.Get(checksumHeader); checksum != "" {
			m.UserAgent.AddUserAgentFeature(feat)
		}
	}

	return next.HandleBuild(ctx, in)
}

type ResponseChecksumMetricsTracking struct {
	ResponseChecksumValidation aws.ResponseChecksumValidation
	UserAgent                  *awsmiddleware.RequestUserAgent
}

func (m *ResponseChecksumMetricsTracking) ID() string {
	return "AWSChecksum:ResponseMetricsTracking"
}

func (m *ResponseChecksumMetricsTracking) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	switch m.ResponseChecksumValidation {
	case aws.ResponseChecksumValidationWhenSupported:
		m.UserAgent.AddUserAgentFeature(awsmiddleware.UserAgentFeatureResponseChecksumWhenSupported)
	case aws.ResponseChecksumValidationWhenRequired:
		m.UserAgent.AddUserAgentFeature(awsmiddleware.UserAgentFeatureResponseChecksumWhenRequired)
	}

	return next.HandleBuild(ctx, in)
}
