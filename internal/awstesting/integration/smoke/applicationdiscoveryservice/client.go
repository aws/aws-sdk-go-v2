// +build integration

//Package applicationdiscoveryservice provides gucumber integration tests support.
package applicationdiscoveryservice

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@applicationdiscoveryservice", func() {
		cfg := integration.Config()
		cfg.Region = "us-west-2"

		gucumber.World["client"] = applicationdiscoveryservice.New(cfg)
	})
}
