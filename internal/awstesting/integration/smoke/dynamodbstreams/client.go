// +build integration

//Package dynamodbstreams provides gucumber integration tests support.
package dynamodbstreams

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@dynamodbstreams", func() {
		gucumber.World["client"] = dynamodbstreams.New(integration.Config())
	})
}
