package lexruntimeservice

import "github.com/aws/aws-sdk-go-v2/aws"

type EndpointResolver struct {
	EndpointID string
}

func (*EndpointResolver) ResolveEndpoint(service, region string) (aws.Endpoint, error) {
	return aws.Endpoint{}, nil
}

func newEndpointResolver() *EndpointResolver {
	return &EndpointResolver{
		EndpointID: "runtime.lex",
	}
}
