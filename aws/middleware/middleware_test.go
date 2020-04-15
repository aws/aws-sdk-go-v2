package middleware_test

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	smithyMiddleware "github.com/awslabs/smithy-go/middleware"
	smithyHTTP "github.com/awslabs/smithy-go/transport/http"
)

type mockBuildHandler func(context.Context, smithyMiddleware.BuildInput) (smithyMiddleware.BuildOutput, smithyMiddleware.Metadata, error)

func (f mockBuildHandler) HandleBuild(ctx context.Context, in smithyMiddleware.BuildInput) (smithyMiddleware.BuildOutput, smithyMiddleware.Metadata, error) {
	return f(ctx, in)
}

func TestRequestInvocationIDMiddleware(t *testing.T) {
	mid := middleware.RequestInvocationIDMiddleware{}

	in := smithyMiddleware.BuildInput{Request: &smithyHTTP.Request{Request: &http.Request{Header: make(http.Header)}}}
	ctx := context.Background()
	_, _, err := mid.HandleBuild(ctx, in, mockBuildHandler(func(ctx context.Context, input smithyMiddleware.BuildInput) (
		out smithyMiddleware.BuildOutput, metadata smithyMiddleware.Metadata, err error,
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

	in = smithyMiddleware.BuildInput{}
	_, _, err = mid.HandleBuild(ctx, in, nil)
	if err != nil {
		if e, a := "unknown transport type", err.Error(); !strings.Contains(a, e) {
			t.Errorf("expected %q, got %q", e, a)
		}
	} else {
		t.Errorf("expected error, got %q", err)
	}
}
