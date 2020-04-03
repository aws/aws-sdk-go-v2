package middleware

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"
)

// RequestInvocationIDMiddleware is a Smithy BuildMiddleware that will generate a unique ID for logical API operation
// invocation.
type RequestInvocationIDMiddleware struct{}

// ID the identifier for the RequestInvocationIDMiddleware
func (r RequestInvocationIDMiddleware) ID() string {
	return "Request Invocation ID Middleware"
}

// HandleBuild attaches a unique operation invocation id for the operation to the request
func (r RequestInvocationIDMiddleware) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (out middleware.BuildOutput, err error) {
	const invocationIDHeader = "amz-sdk-invocation-id"

	invocationID, err := sdk.UUIDVersion4()
	if err != nil {
		return out, err
	}

	switch req := in.Request.(type) {
	case *http.Request:
		req.Header.Set(invocationIDHeader, invocationID)
	default:
		return middleware.BuildOutput{}, fmt.Errorf("unknown transport type %T", req)
	}

	return next.HandleBuild(ctx, in)
}
