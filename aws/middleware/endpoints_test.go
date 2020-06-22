package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

type mockSerializeHandler func(ctx context.Context, in middleware.SerializeInput) (out middleware.SerializeOutput, metadata middleware.Metadata, err error)

func (m mockSerializeHandler) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput,
) (out middleware.SerializeOutput, metadata middleware.Metadata, err error) {
	return m(ctx, in)
}

func TestResolveServiceEndpoint(t *testing.T) {
	initContext := func(signingName, endpointID, region string) context.Context {
		ctx := context.Background()
		ctx = SetSigningName(ctx, signingName)
		ctx = setEndpointPrefix(ctx, endpointID)
		ctx = setRegion(ctx, region)
		return ctx
	}

	cases := map[string]struct {
		Resolver      func(*testing.T) aws.EndpointResolver
		Input         middleware.SerializeInput
		Context       context.Context
		Handler       func(*testing.T) mockSerializeHandler
		ExpectedError string
	}{
		"resolves endpoint and sets signing region": {
			Resolver: func(t *testing.T) aws.EndpointResolver {
				return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					t.Helper()
					if e, a := "fooService", service; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					if e, a := "barRegion", region; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return aws.Endpoint{
						URL:           "https://foo.us-west-2.amazonaws.com",
						SigningRegion: "us-west-2",
					}, nil
				})
			},
			Context: initContext("foo", "fooService", "barRegion"),
			Input:   middleware.SerializeInput{Request: &smithyhttp.Request{Request: &http.Request{}}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return func(ctx context.Context, in middleware.SerializeInput) (out middleware.SerializeOutput, metadata middleware.Metadata, err error) {
					t.Helper()
					if e, a := "us-west-2", GetSigningRegion(ctx); e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					if e, a := "foo", GetSigningName(ctx); e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					if e, a := "https://foo.us-west-2.amazonaws.com", in.Request.(*smithyhttp.Request).URL.String(); e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return out, metadata, err
				}
			},
		},
		"error on invalid url parse": {
			Resolver: func(t *testing.T) aws.EndpointResolver {
				return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL: "://invalid",
					}, nil
				})
			},
			Context: initContext("foo", "fooService", "barRegion"),
			Input:   middleware.SerializeInput{Request: &smithyhttp.Request{Request: &http.Request{}}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return nil
			},
			ExpectedError: "failed to parse endpoint URL",
		},
		"prefers endpoint signing name if available": {
			Resolver: func(t *testing.T) aws.EndpointResolver {
				return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:         "https://foo.us-west-2.amazonaws.com",
						SigningName: "bar",
					}, nil
				})
			},
			Context: initContext("foo", "", ""),
			Input:   middleware.SerializeInput{Request: &smithyhttp.Request{Request: &http.Request{}}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return func(ctx context.Context, in middleware.SerializeInput) (out middleware.SerializeOutput, metadata middleware.Metadata, err error) {
					t.Helper()
					if e, a := "bar", GetSigningName(ctx); e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return out, metadata, err
				}
			},
		},
		"prefers endpoint signing name if previously unknown": {
			Resolver: func(t *testing.T) aws.EndpointResolver {
				return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:         "https://foo.us-west-2.amazonaws.com",
						SigningName: "bar",
					}, nil
				})
			},
			Context: initContext("", "", ""),
			Input:   middleware.SerializeInput{Request: &smithyhttp.Request{Request: &http.Request{}}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return func(ctx context.Context, in middleware.SerializeInput) (out middleware.SerializeOutput, metadata middleware.Metadata, err error) {
					t.Helper()
					if e, a := "bar", GetSigningName(ctx); e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return out, metadata, err
				}
			},
		},
		"prefers known signing name over derived value": {
			Resolver: func(t *testing.T) aws.EndpointResolver {
				return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:                "https://foo.us-west-2.amazonaws.com",
						SigningName:        "bar",
						SigningNameDerived: true,
					}, nil
				})
			},
			Context: initContext("foo", "bar", ""),
			Input:   middleware.SerializeInput{Request: &smithyhttp.Request{Request: &http.Request{}}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return func(ctx context.Context, in middleware.SerializeInput) (out middleware.SerializeOutput, metadata middleware.Metadata, err error) {
					t.Helper()
					if e, a := "foo", GetSigningName(ctx); e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return out, metadata, err
				}
			},
		},
		"errors on endpoint resolver failures": {
			Resolver: func(t *testing.T) aws.EndpointResolver {
				return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					t.Helper()
					return aws.Endpoint{}, fmt.Errorf("some resolver error")
				})
			},
			Context: context.Background(),
			Input:   middleware.SerializeInput{Request: &smithyhttp.Request{Request: &http.Request{}}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return nil
			},
			ExpectedError: "failed to resolve service endpoint",
		},
		"errors on unknown transport types": {
			Resolver: func(t *testing.T) aws.EndpointResolver {
				return nil
			},
			Context: context.Background(),
			Input:   middleware.SerializeInput{Request: struct{}{}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return nil
			},
			ExpectedError: "unknown transport type",
		},
		"errors on nil resolver": {
			Resolver: func(t *testing.T) aws.EndpointResolver {
				return nil
			},
			Context: context.Background(),
			Input:   middleware.SerializeInput{Request: &smithyhttp.Request{Request: &http.Request{}}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return nil
			},
			ExpectedError: "expected endpoint resolver to not be nil",
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			m := ResolveServiceEndpoint{Resolver: tt.Resolver(t)}
			_, _, err := m.HandleSerialize(tt.Context, tt.Input, tt.Handler(t))
			if err != nil && len(tt.ExpectedError) == 0 {
				t.Fatalf("expected no error, got %v", err)
			} else if err != nil && len(tt.ExpectedError) > 0 {
				if e, a := tt.ExpectedError, err.Error(); !strings.Contains(a, e) {
					t.Errorf("expected %v, got %v", e, a)
				}
			} else if err == nil && len(tt.ExpectedError) > 0 {
				t.Error("expected error, got none")
			}
		})
	}
}
