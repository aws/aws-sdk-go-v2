// +build integration

//Package elasticbeanstalk provides gucumber integration tests support.
package elasticbeanstalk

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@elasticbeanstalk", func() {
		gucumber.World["client"] = elasticbeanstalk.New(integration.Config())
	})
}
