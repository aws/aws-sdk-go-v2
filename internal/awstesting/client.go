package awstesting

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
)

// NewClient creates and initializes a generic service client for testing.
func NewClient(cfgs ...*aws.Config) *aws.Client {
	cfg := defaults.Config()
	cfg.EndpointResolver = aws.ResolveStaticEndpointURL("http://endpoint")
	cfg.MergeIn(cfgs...)

	endpoint, _ := cfg.EndpointResolver.EndpointFor("mock-client", aws.StringValue(cfg.Region))
	info := aws.ClientInfo{
		Endpoint: endpoint.URL,
	}

	return aws.NewClient(cfg, info, cfg.Handlers)
}
