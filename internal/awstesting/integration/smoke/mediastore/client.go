// +build integration

//Package mediastore provides gucumber integration tests support.
package mediastore

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/mediastore"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@mediastore", func() {
		gucumber.World["client"] = mediastore.New(integration.Config())
	})
}
