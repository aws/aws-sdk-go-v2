package s3

import (
	"context"
	"errors"
	"fmt"
	"io"

	presignedurl "github.com/aws/aws-sdk-go-v2/service/internal/presigned-url"
	"github.com/aws/smithy-go/middleware"
)

var errNoContentLength = errors.New(
	"The operation input had an undefined content length. PutObject MUST have a " +
		"derivable content length from either (1) an explicit value for the " +
		"ContentLength input member (2) the Body input member implementing io.Seeker " +
		"such that the SDK can derive a value.",
)

// PutObject MUST have a derivable content length for the body in some form,
// since the service does not implement chunked transfer-encoding (and
// aws-chunked encoding requires the length anyway).
//
// We gate this constraint at the client level through additional validation
// rather than letting the request through, which would fail with a 501.
type validatePutObjectContentLength struct{}

func (*validatePutObjectContentLength) ID() string {
	return "validatePutObjectContentLength"
}

func (*validatePutObjectContentLength) HandleInitialize(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if presignedurl.GetIsPresigning(ctx) { // won't have a body
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*PutObjectInput)
	if !ok {
		return out, metadata, fmt.Errorf("unknown input parameters type %T", in.Parameters)
	}

	_, ok = input.Body.(io.Seeker)
	if !ok && input.ContentLength == nil {
		return out, metadata, errNoContentLength
	}
	return next.HandleInitialize(ctx, in)
}

func addValidatePutObjectContentLength(stack *middleware.Stack) error {
	return stack.Initialize.Add(&validatePutObjectContentLength{}, middleware.After)
}
