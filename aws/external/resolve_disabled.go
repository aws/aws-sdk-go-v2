// +build disabled

package external

import "github.com/aws/aws-sdk-go-v2/aws"

// ResolveEndpointResolverFunc extracts the first instance of a EndpointResolverFunc from the config slice
// and sets the functions result on the aws.Config.EndpointResolver
func ResolveEndpointResolverFunc(cfg *aws.Config, configs Configs) error {
	endpointResolverFunc, found, err := GetEndpointResolverFunc(configs)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	cfg.EndpointResolver = endpointResolverFunc(cfg.EndpointResolver)

	return nil
}

type ec2MetadataRegionClient interface {
	Region(context.Context) (string, error)
}

// newEC2MetadataClient is the EC2 instance metadata service client, allows for swapping during testing
var newEC2MetadataClient = func(cfg aws.Config) ec2MetadataRegionClient {
	return ec2metadata.New(cfg)
}

// ResolveEC2Region attempts to resolve the region using the EC2 instance metadata service. If region is already set on
// the config no lookup occurs. If an error is returned the service is assumed unavailable.
func ResolveEC2Region(cfg *aws.Config, _ Configs) error {
	if len(cfg.Region) > 0 {
		return nil
	}

	client := newEC2MetadataClient(*cfg)

	// TODO: What does context look like with external config loading and how to handle the impact to service client config loading
	region, err := client.Region(context.Background())
	if err != nil {
		return nil
	}

	cfg.Region = region

	return nil
}
