// +build integration

//Package mediastoredata provides gucumber integration tests support.
package mediastoredata

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	_ "github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/smoke"
	"github.com/aws/aws-sdk-go-v2/service/mediastore"
	"github.com/aws/aws-sdk-go-v2/service/mediastoredata"
	"github.com/gucumber/gucumber"
)

func init() {
	const containerName = "awsgosdkteamintegcontainer"
	gucumber.Before("@mediastoredata", func() {
		mediastoreSvc := mediastore.New(integration.Config())
		resp, err := mediastoreSvc.DescribeContainerRequest(
			&mediastore.DescribeContainerInput{
				ContainerName: aws.String(containerName),
			}).Send()
		if err != nil {
			gucumber.World["error"] = fmt.Errorf(
				"failed to get mediastore container endpoint for test, %v",
				err)
			return
		}

		cfg := integration.Config()
		cfg.EndpointResolver = aws.ResolveWithEndpointURL(*resp.Container.Endpoint)
		gucumber.World["client"] = mediastoredata.New(cfg)
	})
}
