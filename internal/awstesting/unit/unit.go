// Package unit performs initialization and validation for unit tests
package unit

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
)

func init() {
	// TODO getting a default populated config should be in the "defaults" package
	*Config = defaults.Config()
	Config.EndpointResolver = aws.ResolveStaticEndpointURL("http://endpoint")
	Config.Region = aws.String("mock-region")
	Config.Credentials = aws.NewStaticCredentials("AKID", "SECRET", "SESSION")
}

// Config is a shared config for unit tests to use.
var Config = &aws.Config{}
