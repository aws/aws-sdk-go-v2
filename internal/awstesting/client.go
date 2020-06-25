package awstesting

import (
	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/retry"
)

// NewClient creates and initializes a generic service client for testing.
func NewClient(cfg aws.Config) *aws.Client {
	if cfg.Retryer == nil {
		cfg.Retryer = retry.NewStandard()
	}
	return aws.NewClient(cfg, aws.Metadata{ServiceName: "mockService"})
}
