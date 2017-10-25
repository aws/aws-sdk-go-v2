// +build integration

//Package acm provides gucumber integration tests support.
package acm

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@acm", func() {
		gucumber.World["client"] = acm.New(integration.Config())
	})
}
