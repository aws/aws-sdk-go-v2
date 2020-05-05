package lexruntimeservice

import "github.com/aws/aws-sdk-go-v2/aws"

// EndpointResolver is a stub implementation of a per API client endpoint lookup.
type EndpointResolver struct {
	EndpointID string
}

// ResolveEndpoint returns an endpoint for the API.
func (*EndpointResolver) ResolveEndpoint(service, region string) (aws.Endpoint, error) {
	return aws.Endpoint{}, nil
}

func newEndpointResolver() *EndpointResolver {
	return &EndpointResolver{
		EndpointID: "runtime.lex",
	}
}
