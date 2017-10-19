package aws

import "github.com/aws/aws-sdk-go-v2/aws/endpoints"

// EndpointResolver resolves an endpoint for a service and region.
// TODO need to resolve aws.EndpointResolver and endpoints.Resolver.
type EndpointResolver interface {
	// TODO move endpoint options into this package.
	EndpointFor(service, region string, opts ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error)
}

// EndpointResolverFunc is a helper utility that wraps a function so it satisfies the
// Resolver interface. This is useful when you want to add additional endpoint
// resolving logic, or stub out specific endpoints with custom values.
//
// TODO need to resolve aws.ResolverFunc and endpoints.ResolverFunc.
type EndpointResolverFunc func(service, region string, opts ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error)

// EndpointFor wraps the ResolverFunc function to satisfy the Resolver interface.
func (fn EndpointResolverFunc) EndpointFor(service, region string, opts ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	return fn(service, region, opts...)
}

// ResolveWithEndpoint allows a static Resolved Endpoint to be used as an endpoint resolver
type ResolveWithEndpoint endpoints.ResolvedEndpoint

// ResolveWithEndpointURL allows a static URL to be used as a endpoint resolver.
// TODO is this helper utility funciton needed?
func ResolveWithEndpointURL(url string) ResolveWithEndpoint {
	return ResolveWithEndpoint(endpoints.ResolvedEndpoint{URL: url})
}

// EndpointFor returns the static endpoint.
func (v ResolveWithEndpoint) EndpointFor(service, region string, opts ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	return endpoints.ResolvedEndpoint(v), nil
}

// TODO Endpoint should be in aws
//type Endpoint struct {
//	URL           string
//	SigningName   string
//	SigningRegion string
//	SigningMethod string
//}
