package eventbridge

import (
	"context"
	"fmt"

	smithyauth "github.com/aws/smithy-go/auth"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type endpointAuthResolver struct {
	EndpointResolver EndpointResolverV2
}

var _ AuthSchemeResolver = (*endpointAuthResolver)(nil)

func (r *endpointAuthResolver) ResolveAuthSchemes(
	ctx context.Context, params *AuthResolverParameters,
) (
	[]*smithyauth.Option, error,
) {
	endpt, err := r.EndpointResolver.ResolveEndpoint(ctx, *params.endpointParams)
	if err != nil {
		return nil, fmt.Errorf("resolve endpoint: %v", err)
	}

	if opts, ok := smithyauth.GetAuthOptions(&endpt.Properties); ok {
		return opts, nil
	}

	// endpoint rules didn't specify, fallback to sigv4
	return []*smithyauth.Option{
		smithyhttp.NewSigV4Option(func(props *smithyhttp.SigV4Properties) {
			props.SigningName = "events"
			props.SigningRegion = params.Region
		}),
	}, nil
}

func finalizeServiceEndpointAuthResolver(options *Options) {
	if _, ok := options.AuthSchemeResolver.(*defaultAuthSchemeResolver); !ok {
		return
	}

	options.AuthSchemeResolver = &endpointAuthResolver{
		EndpointResolver: options.EndpointResolverV2,
	}
}

func finalizeOperationEndpointAuthResolver(options *Options) {
	resolver, ok := options.AuthSchemeResolver.(*endpointAuthResolver)
	if !ok {
		return
	}

	if resolver.EndpointResolver == options.EndpointResolverV2 {
		return
	}

	options.AuthSchemeResolver = &endpointAuthResolver{
		EndpointResolver: options.EndpointResolverV2,
	}
}
