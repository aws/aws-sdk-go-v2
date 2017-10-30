package awstesting

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

// NewClient creates and initializes a generic service client for testing.
func NewClient(cfg aws.Config) *aws.Client {
	if cfg.Retryer == nil {
		cfg.Retryer = aws.DefaultRetryer{NumMaxRetries: 3}
	}
	return aws.NewClient(cfg, aws.Metadata{ServiceName: "mockService"})
}
