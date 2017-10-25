// +build integration

//Package cognitosync provides gucumber integration tests support.
package cognitosync

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/cognitosync"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@cognitosync", func() {
		gucumber.World["client"] = cognitosync.New(integration.Config())
	})
}
