package eventbridge

import (
	"context"
	"fmt"

	smithy "github.com/aws/smithy-go"
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
	opts, err := r.resolveAuthSchemes(ctx, params)
	if err != nil {
		return nil, err
	}

	// preserve pre-SRA behavior where everything technically had anonymous
	return append(opts, &smithyauth.Option{
		SchemeID: smithyauth.SchemeIDAnonymous,
	}), nil
}

func (r *endpointAuthResolver) resolveAuthSchemes(
	ctx context.Context, params *AuthResolverParameters,
) (
	[]*smithyauth.Option, error,
) {
	endpt, err := r.EndpointResolver.ResolveEndpoint(ctx, *params.endpointParams)
	if err != nil {
		return nil, fmt.Errorf("resolve endpoint: %w", err)
	}

	if opts, ok := smithyauth.GetAuthOptions(&endpt.Properties); ok {
		return opts, nil
	}

	// endpoint rules didn't specify, fallback to sigv4
	return []*smithyauth.Option{
		{
			SchemeID: smithyauth.SchemeIDSigV4,
			SignerProperties: func() smithy.Properties {
				var props smithy.Properties
				smithyhttp.SetSigV4SigningName(&props, "events")
				smithyhttp.SetSigV4SigningRegion(&props, params.Region)
				return props
			}(),
		},
		{
			SchemeID: smithyauth.SchemeIDSigV4A,
		},
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
