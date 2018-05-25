// +build integration

//Package autoscalingplans provides gucumber integration tests support.
package autoscalingplans

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/autoscalingplans"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@autoscalingplans", func() {
		gucumber.World["client"] = autoscalingplans.New(integration.Config())
	})
}
