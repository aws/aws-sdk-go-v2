package middleware

import (
	"context"

	internalcontext "github.com/aws/aws-sdk-go-v2/internal/context"
	"github.com/aws/smithy-go/middleware"
)

// invocationIDKey is the accessor key used to store the operation invocation
// id on the context, alongside the Amz-Sdk-Invocation-Id header it is sent as.
type invocationIDKey struct{}

// setInvocationID records the operation invocation id on the context. The id
// ties a logged request to the response returned for it, and is shared by all
// retry attempts of a single operation invocation.
//
// Scoped to stack values. Use github.com/aws/smithy-go/middleware#ClearStackValues
// to clear all stack values.
func setInvocationID(ctx context.Context, value string) context.Context {
	return middleware.WithStackValue(ctx, invocationIDKey{}, value)
}

// GetInvocationID retrieves the operation invocation id from the context. The
// id matches the value of the Amz-Sdk-Invocation-Id header sent with the
// request, and is shared by all retry attempts of a single operation
// invocation.
//
// Returns an empty string if the id is not available, i.e. before the request
// has been built, or for stacks that do not record it. The recorded value is
// kept until the stack values are cleared, so it is still reported for entries
// logged after the invocation returned.
//
// Scoped to stack values. Use github.com/aws/smithy-go/middleware#ClearStackValues
// to clear all stack values.
func GetInvocationID(ctx context.Context) (v string) {
	v, _ = middleware.GetStackValue(ctx, invocationIDKey{}).(string)
	return v
}

// GetAttemptNumber retrieves the retry attempt number of the operation
// invocation held on the context, counted from 1. Recorded by the SDK's retryer
// for every attempt, alongside the attempt number reported in the
// Amz-Sdk-Request header and in the attempt results of the response metadata.
//
// Returns 0 if the attempt number is not available, i.e. for stacks without a
// retryer, and for middleware running outside of an attempt.
//
// Scoped to stack values. Use github.com/aws/smithy-go/middleware#ClearStackValues
// to clear all stack values.
func GetAttemptNumber(ctx context.Context) (v int) {
	return internalcontext.GetAttemptNumberContext(ctx)
}
