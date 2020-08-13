package middleware_test

import (
	"context"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	smithymiddleware "github.com/awslabs/smithy-go/middleware"
	smithyHTTP "github.com/awslabs/smithy-go/transport/http"
)

type mockBuildHandler func(context.Context, smithymiddleware.BuildInput) (smithymiddleware.BuildOutput, smithymiddleware.Metadata, error)

func (f mockBuildHandler) HandleBuild(ctx context.Context, in smithymiddleware.BuildInput) (smithymiddleware.BuildOutput, smithymiddleware.Metadata, error) {
	return f(ctx, in)
}

func TestRequestInvocationIDMiddleware(t *testing.T) {
	mid := middleware.RequestInvocationIDMiddleware{}

	in := smithymiddleware.BuildInput{Request: &smithyHTTP.Request{Request: &http.Request{Header: make(http.Header)}}}
	ctx := context.Background()
	_, _, err := mid.HandleBuild(ctx, in, mockBuildHandler(func(ctx context.Context, input smithymiddleware.BuildInput) (
		out smithymiddleware.BuildOutput, metadata smithymiddleware.Metadata, err error,
	) {
		req := in.Request.(*smithyHTTP.Request)

		value := req.Header.Get("amz-sdk-invocation-id")

		match, err := regexp.MatchString(`[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$`, value)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !match {
			t.Errorf("invocation id was not a UUIDv4")
		}

		return out, metadata, err
	}))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	in = smithymiddleware.BuildInput{}
	_, _, err = mid.HandleBuild(ctx, in, nil)
	if err != nil {
		if e, a := "unknown transport type", err.Error(); !strings.Contains(a, e) {
			t.Errorf("expected %q, got %q", e, a)
		}
	} else {
		t.Errorf("expected error, got %q", err)
	}
}

type mockDeserializeHandler func(ctx context.Context, in smithymiddleware.DeserializeInput) (smithymiddleware.DeserializeOutput, smithymiddleware.Metadata, error)

func (m mockDeserializeHandler) HandleDeserialize(ctx context.Context, in smithymiddleware.DeserializeInput) (smithymiddleware.DeserializeOutput, smithymiddleware.Metadata, error) {
	return m(ctx, in)
}

func TestAttemptClockSkewHandler(t *testing.T) {
	cases := map[string]struct {
		Next       mockDeserializeHandler
		Expect     middleware.ResponseMetadata
		ResponseAt func() time.Time
	}{
		"no response": {
			Next: mockDeserializeHandler(func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				return out, m, err
			}),
		},
		"failed response": {
			Next: mockDeserializeHandler(func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				out.RawResponse = &smithyHTTP.Response{
					Response: &http.Response{
						StatusCode: 0,
						Header:     http.Header{},
					},
				}
				return out, m, err
			}),
		},
		"no date header response": {
			Next: mockDeserializeHandler(func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				out.RawResponse = &smithyHTTP.Response{
					Response: &http.Response{
						StatusCode: 200,
						Header:     http.Header{},
					},
				}
				return out, m, err
			}),
		},
		"invalid date header response": {
			Next: mockDeserializeHandler(func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				out.RawResponse = &smithyHTTP.Response{
					Response: &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Date": []string{"abc123"},
						},
					},
				}
				return out, m, err
			}),
		},
		"date response": {
			Next: mockDeserializeHandler(func(ctx context.Context, in smithymiddleware.DeserializeInput,
			) (out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error) {
				out.RawResponse = &smithyHTTP.Response{
					Response: &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Date": []string{"Thu, 05 Mar 2020 22:25:15 GMT"},
						},
					},
				}
				return out, m, err
			}),
			ResponseAt: func() time.Time {
				return time.Date(2020, 3, 5, 22, 25, 17, 0, time.UTC)
			},
			Expect: middleware.ResponseMetadata{
				AttemptSkew: -2 * time.Second,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			mw := middleware.AttemptClockSkewMiddleware{}
			_, metadata, err := mw.HandleDeserialize(context.Background(), smithymiddleware.DeserializeInput{}, c.Next)
			if err != nil {
				t.Errorf("expect no error, got %v", err)
			}

			if e, a := &c, metadata; !reflect.DeepEqual(e, a) {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
