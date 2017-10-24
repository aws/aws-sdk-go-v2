// +build integration

//Package dynamodb provides gucumber integration tests support.
package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@dynamodb", func() {
		gucumber.World["client"] = dynamodb.New(integration.Config())
	})
}
