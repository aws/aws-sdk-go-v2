// +build integration

//Package simpledb provides gucumber integration tests support.
package simpledb

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/simpledb"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@simpledb", func() {
		gucumber.World["client"] = simpledb.New(integration.Config())
	})
}
