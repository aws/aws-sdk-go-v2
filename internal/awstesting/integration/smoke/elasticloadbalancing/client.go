// +build integration

//Package elasticloadbalancing provides gucumber integration tests support.
package elasticloadbalancing

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/elb"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@elasticloadbalancing", func() {
		gucumber.World["client"] = elb.New(integration.Config())
	})
}
