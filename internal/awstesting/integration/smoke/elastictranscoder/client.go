// +build integration

//Package elastictranscoder provides gucumber integration tests support.
package elastictranscoder

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/elastictranscoder"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@elastictranscoder", func() {
		gucumber.World["client"] = elastictranscoder.New(integration.Config())
	})
}
