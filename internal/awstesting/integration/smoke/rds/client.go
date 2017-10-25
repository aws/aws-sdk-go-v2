// +build integration

//Package rds provides gucumber integration tests support.
package rds

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@rds", func() {
		gucumber.World["client"] = rds.New(integration.Config())
	})
}
