package middleware

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var expectedAgent = aws.SDKName + "/" + aws.SDKVersion + " GOOS/" + runtime.GOOS + " GOARCH/" + runtime.GOARCH + " GO/" + runtime.Version()

type mockBuilderHandler func(context.Context, middleware.BuildInput) (middleware.BuildOutput, middleware.Metadata, error)

func (m mockBuilderHandler) HandleBuild(ctx context.Context, in middleware.BuildInput) (out middleware.BuildOutput, metadata middleware.Metadata, err error) {
	return m(ctx, in)
}

func TestRequestUserAgent_HandleBuild(t *testing.T) {
	cases := map[string]struct {
		Env    map[string]string
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
					"User-Agent": {expectedAgent},
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
					"User-Agent": {"previously set " + expectedAgent},
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
		"adds exec-env if present": {
			Env: map[string]string{
				execEnvVar: "TestCase",
			},
			In: middleware.BuildInput{Request: &smithyhttp.Request{
				Request: &http.Request{Header: map[string][]string{}},
			}},
			Expect: middleware.BuildInput{Request: &smithyhttp.Request{
				Request: &http.Request{Header: map[string][]string{
					"User-Agent": {expectedAgent + " exec-env/TestCase"},
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
			if len(tt.Env) > 0 {
				environ := os.Environ()
				os.Clearenv()
				defer func() {
					os.Clearenv()
					for _, v := range environ {
						split := strings.SplitN(v, "=", 2)
						key, value := split[0], split[1]
						os.Setenv(key, value)
					}
				}()
				for k, v := range tt.Env {
					os.Setenv(k, v)
				}
			}

			b := NewRequestUserAgent()
			_, _, err := b.HandleBuild(context.Background(), tt.In, tt.Next(t, tt.Expect))
			if (err != nil) != tt.Err {
				t.Errorf("error %v, want error %v", err, tt.Err)
				return
			}
		})
	}
}

func TestAddUserAgentKey(t *testing.T) {
	b := NewRequestUserAgent()
	stack := middleware.NewStack("testStack", smithyhttp.NewStackRequest)
	err := stack.Build.Add(b, middleware.After)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	err = AddUserAgentKey("foo")(stack)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	bi := middleware.BuildInput{Request: &smithyhttp.Request{Request: &http.Request{Header: map[string][]string{}}}}
	_, _, err = b.HandleBuild(context.Background(), bi, mockBuilderHandler(func(ctx context.Context, input middleware.BuildInput) (o middleware.BuildOutput, m middleware.Metadata, err error) {
		return o, m, err
	}))
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	ua, ok := bi.Request.(*smithyhttp.Request).Header["User-Agent"]
	if !ok {
		t.Fatalf("expect User-Agent to be present")
	}
	if ua[0] != expectedAgent+" foo" {
		t.Error("User-Agent did not match expected")
	}
}

func TestAddUserAgentKeyValue(t *testing.T) {
	b := NewRequestUserAgent()
	stack := middleware.NewStack("testStack", smithyhttp.NewStackRequest)
	err := stack.Build.Add(b, middleware.After)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	err = AddUserAgentKeyValue("foo", "bar")(stack)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	bi := middleware.BuildInput{Request: &smithyhttp.Request{Request: &http.Request{Header: map[string][]string{}}}}
	_, _, err = b.HandleBuild(context.Background(), bi, mockBuilderHandler(func(ctx context.Context, input middleware.BuildInput) (o middleware.BuildOutput, m middleware.Metadata, err error) {
		return o, m, err
	}))
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	ua, ok := bi.Request.(*smithyhttp.Request).Header["User-Agent"]
	if !ok {
		t.Fatalf("expect User-Agent to be present")
	}
	if ua[0] != expectedAgent+" foo/bar" {
		t.Error("User-Agent did not match expected")
	}
}
