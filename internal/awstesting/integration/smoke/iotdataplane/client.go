// +build integration

//Package iotdataplane provides gucumber integration tests support.
package iotdataplane

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@iotdataplane", func() {
		svc := iot.New(integration.Config())
		result, err := svc.DescribeEndpointRequest(&iot.DescribeEndpointInput{}).Send()
		if err != nil {
			gucumber.World["error"] = err
			return
		}

		cfg := integration.Config()
		cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://" + *result.EndpointAddress)

		gucumber.World["client"] = iotdataplane.New(cfg)
	})
}
