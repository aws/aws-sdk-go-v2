package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
	initContext := func(endpointID, region string, metadata OperationMetadata) context.Context {
		ctx := context.Background()
		ctx = setEndpointID(ctx, endpointID)
		ctx = setRegion(ctx, region)
		ctx = setOperationMetadata(ctx, metadata)
		return ctx
	}

	cases := []struct {
		Resolver      func(*testing.T) aws.EndpointResolver
		Input         middleware.SerializeInput
		Context       context.Context
		Handler       func(*testing.T) mockSerializeHandler
		ExpectedError string
	}{
		{
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
			Context: initContext("fooService", "barRegion", OperationMetadata{HTTPPath: "/foo"}),
			Input:   middleware.SerializeInput{Request: &smithyhttp.Request{Request: &http.Request{}}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return func(ctx context.Context, in middleware.SerializeInput) (out middleware.SerializeOutput, metadata middleware.Metadata, err error) {
					t.Helper()
					if e, a := "us-west-2", GetSigningRegion(ctx); e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					req := in.Request.(*smithyhttp.Request)
					if e, a := "https://foo.us-west-2.amazonaws.com/foo", req.URL.String(); e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return out, metadata, err
				}
			},
		},
		{
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
		{
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
		{
			Resolver: func(t *testing.T) aws.EndpointResolver {
				return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					t.Helper()
					return aws.Endpoint{URL: "://malformed/foo"}, nil
				})
			},
			Context: context.Background(),
			Input:   middleware.SerializeInput{Request: &smithyhttp.Request{Request: &http.Request{}}},
			Handler: func(t *testing.T) mockSerializeHandler {
				return nil
			},
			ExpectedError: "failed to parse endpoint URL",
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
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
