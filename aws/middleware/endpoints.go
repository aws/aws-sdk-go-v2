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

	var endpoint aws.Endpoint
	endpoint, err = r.Resolver.ResolveEndpoint(GetEndpointID(ctx), GetRegion(ctx))
	if err != nil {
		return out, metadata, fmt.Errorf("failed to resolve service endpoint ")
	}

	operationMetadata := GetOperationMetadata(ctx)
	req.URL, err = url.Parse(endpoint.URL + operationMetadata.HTTPPath)
	if err != nil {
		return out, metadata, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}

	ctx = SetSigningRegion(ctx, endpoint.SigningRegion)
	return next.HandleSerialize(ctx, in)
}

// ResolveServiceEndpointOptions is an interface for retrieving options to configure a ResolveServiceEndpoint middleware.
type ResolveServiceEndpointOptions interface {
	GetResolver() aws.EndpointResolver
}

// AddResolveServiceEndpointMiddleware creates a ResolveServiceEndpoint middleware with the provided options
// and registers it on the provided stack.
func AddResolveServiceEndpointMiddleware(stack *middleware.Stack, options ResolveServiceEndpointOptions) error {
	m := ResolveServiceEndpoint{
		Resolver: options.GetResolver(),
	}
	return stack.Serialize.Add(m, middleware.Before)
}

// RemoveResolveServiceEndpointMiddleware removes the ResolveServiceEndpoint middleware from the provided stack.
func RemoveResolveServiceEndpointMiddleware(stack *middleware.Stack) error {
	id := ResolveServiceEndpoint{}.ID()
	return stack.Serialize.Remove(id)
}
