// +build integration

//Package s3crypto provides gucumber integration tests support.
package s3crypto

import (
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3crypto"

	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@s3crypto", func() {
		cfg := integration.Config()
		cfg.Region = "us-west-2"

		encryptionClient := s3crypto.NewEncryptionClient(cfg, nil,
			func(c *s3crypto.EncryptionClient) {},
		)
		gucumber.World["encryptionClient"] = encryptionClient

		decryptionClient := s3crypto.NewDecryptionClient(cfg)
		gucumber.World["decryptionClient"] = decryptionClient

		gucumber.World["client"] = s3.New(cfg)
	})
}
