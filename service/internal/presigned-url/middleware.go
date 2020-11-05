package presignedurl

import (
	"context"
	"fmt"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"

	"github.com/awslabs/smithy-go/middleware"
)

// ParameterAccessor provides an collection of accessor to for retrieving and
// setting the values needed to PresignedURL generation
type ParameterAccessor struct {
	GetPresignedURL func(interface{}) (string, bool, error)
	GetSourceRegion func(interface{}) (string, bool, error)

	CopyInput            func(interface{}) (interface{}, error)
	SetDestinationRegion func(interface{}, string) error
	SetPresignedURL      func(interface{}, string) error

	// GetPresigner fetches a request presigner
	// The function should be of format:
	// func $name (ctx context.Context, srcRegion string, params interface{}) (v4.HTTPRequest, error)
	GetPresigner func(context.Context, string, interface{}) (v4.PresignedHTTPRequest, error)
}

// Options provides the set of options needed by the presigned URL middleware.
type Options struct {
	Accessor ParameterAccessor
}

// AddMiddleware adds the Presign URL middleware to the middleware stack.
func AddMiddleware(stack *middleware.Stack, opts Options) error {
	return stack.Initialize.Add(&presignMiddleware{options: opts}, middleware.Before)
}

// RemoveMiddleware removes the Presign URL middleware from the stack.
func RemoveMiddleware(stack *middleware.Stack) error {
	return stack.Initialize.Remove((*presignMiddleware)(nil).ID())
}

type presignMiddleware struct {
	options Options
}

func (m *presignMiddleware) ID() string { return "PresignURLCustomization" }

func (m *presignMiddleware) HandleInitialize(
	ctx context.Context, input middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	// If PresignedURL is already set ignore middleware.
	if _, ok, err := m.options.Accessor.GetPresignedURL(input.Parameters); err != nil {
		return out, metadata, fmt.Errorf("presign middleware failed, %w", err)
	} else if ok {
		return next.HandleInitialize(ctx, input)
	}

	// If have source region is not set ignore middleware.
	srcRegion, ok, err := m.options.Accessor.GetSourceRegion(input.Parameters)
	if err != nil {
		return out, metadata, fmt.Errorf("presign middleware failed, %w", err)
	} else if !ok || len(srcRegion) == 0 {
		return next.HandleInitialize(ctx, input)
	}

	// Create a copy of the original input so the destination region value can
	// be added. This ensures that value does not leak into the original
	// request parameters.
	paramCpy, err := m.options.Accessor.CopyInput(input.Parameters)
	if err != nil {
		return out, metadata, fmt.Errorf("unable to create presigned URL, %w", err)
	}

	// Destination region is the API client's configured region.
	dstRegion := awsmiddleware.GetRegion(ctx)
	if err = m.options.Accessor.SetDestinationRegion(paramCpy, dstRegion); err != nil {
		return out, metadata, fmt.Errorf("presign middleware failed, %w", err)
	}

	presignedReq, err := m.options.Accessor.GetPresigner(ctx, srcRegion, paramCpy)
	if err != nil {
		return out, metadata, fmt.Errorf("unable to create presigned URL, %w", err)
	}

	// Update the original input with the presigned URL value.
	if err = m.options.Accessor.SetPresignedURL(input.Parameters, presignedReq.URL); err != nil {
		return out, metadata, fmt.Errorf("presign middleware failed, %w", err)
	}

	return next.HandleInitialize(ctx, input)
}
