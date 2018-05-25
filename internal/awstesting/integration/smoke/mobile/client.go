// +build integration

//Package mobile provides gucumber integration tests support.
package mobile

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/mobile"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@mobile", func() {
		gucumber.World["client"] = mobile.New(integration.Config())
	})
}
