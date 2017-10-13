package awstesting

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	client "github.com/aws/aws-sdk-go-v2/aws"
	metadata "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
)

// NewClient creates and initializes a generic service client for testing.
func NewClient(cfgs ...*aws.Config) *client.Client {
	info := metadata.ClientInfo{
		Endpoint:    "http://endpoint",
		SigningName: "",
	}
	cfg := defaults.Config()
	cfg.MergeIn(cfgs...)

	if v := aws.StringValue(cfg.Endpoint); len(v) > 0 {
		info.Endpoint = v
	}

	return aws.NewClient(cfg, info, cfg.Handlers)
}
