package middleware

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// ResolveServiceEndpoint is a middleware that will resolve the service endpoint.
type ResolveServiceEndpoint struct {
	Resolver aws.EndpointResolver
}

// ID is a middleware identifier
func (r ResolveServiceEndpoint) ID() string {
	return "ResolveServiceEndpoint"
}

// HandleSerialize resolves the service endpoint and returns
func (r ResolveServiceEndpoint) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (out middleware.SerializeOutput, metadata middleware.Metadata, err error) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", in.Request)
	}

	if r.Resolver == nil {
		return out, metadata, fmt.Errorf("expected endpoint resolver to not be nil")
	}

	var endpoint aws.Endpoint
	endpoint, err = r.Resolver.ResolveEndpoint(GetEndpointPrefix(ctx), GetRegion(ctx))
	if err != nil {
		return out, metadata, fmt.Errorf("failed to resolve service endpoint ")
	}

	req.URL, err = url.Parse(endpoint.URL)
	if err != nil {
		return out, metadata, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}

	if len(endpoint.SigningName) > 0 && (len(GetSigningName(ctx)) == 0 || !endpoint.SigningNameDerived) {
		ctx = SetSigningName(ctx, endpoint.SigningName)
	}
	ctx = SetSigningRegion(ctx, endpoint.SigningRegion)

	return next.HandleSerialize(ctx, in)
}

// ResolveServiceEndpointOptions is an interface for retrieving options to configure a ResolveServiceEndpoint middleware.
type ResolveServiceEndpointOptions interface {
	GetEndpointResolver() aws.EndpointResolver
}

// AddResolveServiceEndpointMiddleware creates a ResolveServiceEndpoint middleware with the provided options
// and registers it on the provided stack.
func AddResolveServiceEndpointMiddleware(stack *middleware.Stack, options ResolveServiceEndpointOptions) error {
	m := ResolveServiceEndpoint{
		Resolver: options.GetEndpointResolver(),
	}
	return stack.Serialize.Add(m, middleware.Before)
}

// RemoveResolveServiceEndpointMiddleware removes the ResolveServiceEndpoint middleware from the provided stack.
func RemoveResolveServiceEndpointMiddleware(stack *middleware.Stack) error {
	id := ResolveServiceEndpoint{}.ID()
	return stack.Serialize.Remove(id)
}
