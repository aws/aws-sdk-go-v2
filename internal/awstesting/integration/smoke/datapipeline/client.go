// +build integration

//Package datapipeline provides gucumber integration tests support.
package datapipeline

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/datapipeline"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@datapipeline", func() {
		gucumber.World["client"] = datapipeline.New(integration.Config())
	})
}
