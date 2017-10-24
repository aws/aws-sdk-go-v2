// +build integration

//Package sns provides gucumber integration tests support.
package sns

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@sns", func() {
		gucumber.World["client"] = sns.New(integration.Config())
	})
}
