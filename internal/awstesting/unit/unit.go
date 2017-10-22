// Package unit performs initialization and validation for unit tests
package unit

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
)

func init() {
	// TODO getting a default populated config should be in the "defaults" package
	*Config = defaults.Config()
	Config.Region = aws.String("mock-region")
	Config.CredentialsLoader = aws.NewCredentialsLoader(aws.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION",
			Source: "unit test credentials",
		},
	})
}

// Config is a shared config for unit tests to use.
var Config = &aws.Config{}
