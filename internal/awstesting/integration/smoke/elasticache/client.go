// +build integration

//Package elasticache provides gucumber integration tests support.
package elasticache

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@elasticache", func() {
		gucumber.World["client"] = elasticache.New(integration.Config())
	})
}
