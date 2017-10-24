// +build integration

//Package directoryservice provides gucumber integration tests support.
package directoryservice

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@directoryservice", func() {
		gucumber.World["client"] = directoryservice.New(integration.Config())
	})
}
