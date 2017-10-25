// +build integration

//Package configservice provides gucumber integration tests support.
package configservice

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@configservice", func() {
		gucumber.World["client"] = configservice.New(integration.Config())
	})
}
