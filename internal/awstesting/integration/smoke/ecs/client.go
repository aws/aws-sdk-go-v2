// +build integration

//Package ecs provides gucumber integration tests support.
package ecs

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@ecs", func() {
		// FIXME remove custom region
		cfg := integration.Config()
		cfg.Region = "us-west-2"

		gucumber.World["client"] = ecs.New(cfg)
	})
}
