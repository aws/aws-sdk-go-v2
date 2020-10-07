package testing

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// EndpointResolverFunc is a mock s3 endpoint resolver that wraps the given function
type EndpointResolverFunc func(region string, options s3.ResolverOptions) (aws.Endpoint, error)

func (m EndpointResolverFunc) ResolveEndpoint(region string, options s3.ResolverOptions) (aws.Endpoint, error) {
	return m(region, options)
}
