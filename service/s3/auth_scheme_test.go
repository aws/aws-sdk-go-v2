package s3

import (
	"context"
	"testing"

	smithy "github.com/aws/smithy-go"
	"github.com/aws/smithy-go/auth"
	smithyauth "github.com/aws/smithy-go/auth"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// FUTURE: move to smithy-go, see https://github.com/aws/smithy-go/issues/528

type mockAuthScheme struct {
	id         string
	configured bool
}

var _ smithyhttp.AuthScheme = (*mockAuthScheme)(nil)

func (m *mockAuthScheme) SchemeID() string { return m.id }

func (m *mockAuthScheme) IdentityResolver(_ auth.IdentityResolverOptions) auth.IdentityResolver {
	if m.configured {
		return &mockIdentityResolver{}
	}
	return nil
}

func (*mockAuthScheme) Signer() smithyhttp.Signer { return nil }

type mockIdentityResolver struct{}

var _ smithyauth.IdentityResolver = (*mockIdentityResolver)(nil)

func (*mockIdentityResolver) GetIdentity(ctx context.Context, props smithy.Properties) (smithyauth.Identity, error) {
	return nil, nil
}

func contains(have []string, want string) bool {
	for _, v := range have {
		if v == want {
			return true
		}
	}
	return false
}

func TestSelectScheme(t *testing.T) {
	for name, c := range map[string]struct {
		Supported  []string
		Configured []string
		Expect     string
	}{
		"support(sigv4, bearer) + cfg(sigv4, bearer) = sigv4": {
			Supported: []string{
				smithyauth.SchemeIDSigV4,
				smithyauth.SchemeIDHTTPBearer,
			},
			Configured: []string{
				smithyauth.SchemeIDSigV4,
				smithyauth.SchemeIDHTTPBearer,
			},
			Expect: smithyauth.SchemeIDSigV4,
		},
		"support(sigv4, bearer) + cfg(bearer) = bearer": {
			Supported: []string{
				smithyauth.SchemeIDSigV4,
				smithyauth.SchemeIDHTTPBearer,
			},
			Configured: []string{
				smithyauth.SchemeIDHTTPBearer,
			},
			Expect: smithyauth.SchemeIDHTTPBearer,
		},
		"support(sigv4, bearer) + cfg(sigv4) = sigv4": {
			Supported: []string{
				smithyauth.SchemeIDSigV4,
				smithyauth.SchemeIDHTTPBearer,
			},
			Configured: []string{
				smithyauth.SchemeIDSigV4,
			},
			Expect: smithyauth.SchemeIDSigV4,
		},
		"support(sigv4, bearer) + cfg(n/a) = error": {
			Supported: []string{
				smithyauth.SchemeIDSigV4,
				smithyauth.SchemeIDHTTPBearer,
			},
			Configured: []string{},
			Expect:     "",
		},
		"support(sigv4) + cfg(bearer) = error": {
			Supported: []string{
				smithyauth.SchemeIDSigV4,
			},
			Configured: []string{
				smithyauth.SchemeIDHTTPBearer,
			},
			Expect: "",
		},
		"support(anon) + cfg(sigv4) = anon": {
			Supported: []string{
				smithyauth.SchemeIDAnonymous,
			},
			Configured: []string{
				smithyauth.SchemeIDSigV4,
			},
			Expect: smithyauth.SchemeIDAnonymous,
		},
		"support(anon) + cfg(n/a) = anon": {
			Supported: []string{
				smithyauth.SchemeIDAnonymous,
			},
			Configured: []string{},
			Expect:     smithyauth.SchemeIDAnonymous,
		},
	} {
		t.Run(name, func(t *testing.T) {
			authopts := []*smithyauth.Option{}
			for _, id := range c.Supported {
				authopts = append(authopts, &smithyauth.Option{
					SchemeID: id,
				})
			}

			m := &resolveAuthSchemeMiddleware{
				options: Options{
					AuthSchemes: []smithyhttp.AuthScheme{
						&mockAuthScheme{
							id:         smithyauth.SchemeIDSigV4,
							configured: contains(c.Configured, smithyauth.SchemeIDSigV4),
						},
						&mockAuthScheme{
							id:         smithyauth.SchemeIDHTTPBearer,
							configured: contains(c.Configured, smithyauth.SchemeIDHTTPBearer),
						},
						&mockAuthScheme{
							id:         smithyauth.SchemeIDAnonymous,
							configured: true,
						},
					},
				},
			}

			scheme, ok := m.selectScheme(authopts)
			if c.Expect != "" {
				if !ok {
					t.Errorf("expected scheme '%s', got none", c.Expect)
				}
				if actual := scheme.Scheme.SchemeID(); c.Expect != actual {
					t.Errorf("expected scheme '%s', got '%s'", c.Expect, actual)
				}
			} else {
				if ok {
					t.Errorf("expected no scheme, got '%s'", scheme.Scheme.SchemeID())
				}
			}
		})
	}
}
