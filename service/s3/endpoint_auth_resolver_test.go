package s3

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	smithy "github.com/aws/smithy-go"
	smithyauth "github.com/aws/smithy-go/auth"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type mockEndpointResolver struct {
	endpoint smithyendpoints.Endpoint
	err      error
}

func (m *mockEndpointResolver) ResolveEndpoint(_ context.Context, _ EndpointParameters) (smithyendpoints.Endpoint, error) {
	return m.endpoint, m.err
}

func makeEndpoint(opts []*smithyauth.Option) smithyendpoints.Endpoint {
	uri, _ := url.Parse("https://mybucket.s3.us-west-2.amazonaws.com")
	var props smithy.Properties
	smithyauth.SetAuthOptions(&props, opts)
	return smithyendpoints.Endpoint{
		URI:        *uri,
		Headers:    http.Header{},
		Properties: props,
	}
}

func TestEndpointAuthResolver_BareSchemeID(t *testing.T) {
	endpt := makeEndpoint([]*smithyauth.Option{
		{
			SchemeID: "sigv4",
			SignerProperties: func() smithy.Properties {
				var sp smithy.Properties
				smithyhttp.SetSigV4SigningName(&sp, "s3")
				smithyhttp.SetSigV4ASigningName(&sp, "s3")
				smithyhttp.SetSigV4SigningRegion(&sp, "us-west-2")
				smithyhttp.SetDisableDoubleEncoding(&sp, true)
				return sp
			}(),
		},
	})

	resolver := &endpointAuthResolver{
		EndpointResolver: &mockEndpointResolver{endpoint: endpt},
	}

	region := "us-west-2"
	params := &AuthResolverParameters{
		Operation:      "GetObject",
		Region:         region,
		endpointParams: &EndpointParameters{Region: &region},
	}

	opts, err := resolver.ResolveAuthSchemes(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sigv4Opt := findOptByScheme(opts, "sigv4")
	if sigv4Opt == nil {
		t.Fatal("expected sigv4 option in resolved schemes")
	}

	if v, ok := smithyhttp.GetDisableDoubleEncoding(&sigv4Opt.SignerProperties); !ok || !v {
		t.Error("expected DisableDoubleEncoding from endpoint option")
	}
	if v, ok := smithyhttp.GetSigV4SigningRegion(&sigv4Opt.SignerProperties); !ok || v != "us-west-2" {
		t.Errorf("expected signing region %q, got %q", "us-west-2", v)
	}
	if v, ok := smithyhttp.GetSigV4SigningName(&sigv4Opt.SignerProperties); !ok || v != "s3" {
		t.Errorf("expected signing name %q, got %q", "s3", v)
	}
}

func TestEndpointAuthResolver_RebaseProps(t *testing.T) {
	endpt := makeEndpoint([]*smithyauth.Option{
		{
			SchemeID: "sigv4",
			SignerProperties: func() smithy.Properties {
				var sp smithy.Properties
				// "forget" the signing region
				smithyhttp.SetSigV4SigningName(&sp, "s3")
				smithyhttp.SetDisableDoubleEncoding(&sp, true)
				return sp
			}(),
		},
	})

	resolver := &endpointAuthResolver{
		EndpointResolver: &mockEndpointResolver{endpoint: endpt},
	}

	region := "us-west-2"
	params := &AuthResolverParameters{
		Operation:      "GetObject",
		Region:         region,
		endpointParams: &EndpointParameters{Region: &region},
	}

	opts, err := resolver.ResolveAuthSchemes(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sigv4Opt := findOptByScheme(opts, "sigv4")
	if sigv4Opt == nil {
		t.Fatal("expected sigv4 option in resolved schemes")
	}

	// should come from the base props
	if v, ok := smithyhttp.GetSigV4SigningRegion(&sigv4Opt.SignerProperties); !ok || v != "us-west-2" {
		t.Errorf("expected signing region %q from base rebase, got %q (ok=%v)", "us-west-2", v, ok)
	}
	// rebase should preserve endpoint props
	if v, ok := smithyhttp.GetDisableDoubleEncoding(&sigv4Opt.SignerProperties); !ok || !v {
		t.Error("expected DisableDoubleEncoding preserved from endpoint option")
	}
}

func findOptByScheme(opts []*smithyauth.Option, schemeID string) *smithyauth.Option {
	for _, opt := range opts {
		if opt.SchemeID == schemeID {
			return opt
		}
	}
	return nil
}
