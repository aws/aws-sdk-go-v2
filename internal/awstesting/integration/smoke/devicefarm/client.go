// +build integration

//Package devicefarm provides gucumber integration tests support.
package devicefarm

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/devicefarm"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@devicefarm", func() {
		// FIXME remove custom region
		cfg := integration.Config()
		cfg.Region = "us-west-2"

		gucumber.World["client"] = devicefarm.New(cfg)
	})
}
