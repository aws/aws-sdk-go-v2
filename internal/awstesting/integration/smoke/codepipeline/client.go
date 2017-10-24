// +build integration

//Package codepipeline provides gucumber integration tests support.
package codepipeline

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@codepipeline", func() {
		gucumber.World["client"] = codepipeline.New(integration.Config())
	})
}
