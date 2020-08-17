package middleware

import (
	"context"
	"net/http"
	"runtime"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type mockBuilderHandler func(context.Context, middleware.BuildInput) (middleware.BuildOutput, middleware.Metadata, error)

func (m mockBuilderHandler) HandleBuild(ctx context.Context, in middleware.BuildInput) (out middleware.BuildOutput, metadata middleware.Metadata, err error) {
	return m(ctx, in)
}

func TestRequestUserAgent_HandleBuild(t *testing.T) {
	userAgent := aws.SDKName + "/" + aws.SDKVersion + " GOOS/" + runtime.GOOS + " GOARCH/" + runtime.GOARCH + " GO/" + runtime.Version()
	cases := map[string]struct {
		In     middleware.BuildInput
		Next   func(*testing.T, middleware.BuildInput) middleware.BuildHandler
		Expect middleware.BuildInput
		Err    bool
	}{
		"adds product information": {
			In: middleware.BuildInput{Request: &smithyhttp.Request{
				Request: &http.Request{Header: map[string][]string{}},
			}},
			Expect: middleware.BuildInput{Request: &smithyhttp.Request{
				Request: &http.Request{Header: map[string][]string{
					"User-Agent": {userAgent},
				}},
			}},
			Next: func(t *testing.T, expect middleware.BuildInput) middleware.BuildHandler {
				return mockBuilderHandler(func(ctx context.Context, input middleware.BuildInput) (o middleware.BuildOutput, m middleware.Metadata, err error) {
					if diff := cmp.Diff(input, expect, cmpopts.IgnoreUnexported(http.Request{}, smithyhttp.Request{})); len(diff) > 0 {
						t.Error(diff)
					}
					return o, m, err
				})
			},
		},
		"appends to existing": {
			In: middleware.BuildInput{Request: &smithyhttp.Request{
				Request: &http.Request{Header: map[string][]string{
					"User-Agent": {"previously set"},
				}},
			}},
			Expect: middleware.BuildInput{Request: &smithyhttp.Request{
				Request: &http.Request{Header: map[string][]string{
					"User-Agent": {"previously set " + userAgent},
				}},
			}},
			Next: func(t *testing.T, expect middleware.BuildInput) middleware.BuildHandler {
				return mockBuilderHandler(func(ctx context.Context, input middleware.BuildInput) (o middleware.BuildOutput, m middleware.Metadata, err error) {
					if diff := cmp.Diff(input, expect, cmpopts.IgnoreUnexported(http.Request{}, smithyhttp.Request{})); len(diff) > 0 {
						t.Error(diff)
					}
					return o, m, err
				})
			},
		},
		"errors for unknown type": {
			In: middleware.BuildInput{Request: struct{}{}},
			Next: func(t *testing.T, input middleware.BuildInput) middleware.BuildHandler {
				return nil
			},
			Err: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			b := NewRequestUserAgent()
			_, _, err := b.HandleBuild(context.Background(), tt.In, tt.Next(t, tt.Expect))
			if (err != nil) != tt.Err {
				t.Errorf("error %v, want error %v", err, tt.Err)
				return
			}
		})
	}
}
