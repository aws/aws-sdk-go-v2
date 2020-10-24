package presignedurl

import (
	"context"
	"fmt"
	"net/http"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/id"
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
}

// URLPresigner provides the interface to presign the input parameters in to a
// presigned URL.
type URLPresigner interface {
	PresignURL(ctx context.Context, srcRegion string, params interface{}) (
		presignedURL string, signedHeader http.Header, err error,
	)
}

// Options provides the set of options needed by the presigned URL middleware.
type Options struct {
	Accessor  ParameterAccessor
	Presigner URLPresigner
}

// AddMiddleware adds the Presign URL middleware to the middleware stack.
func AddMiddleware(stack *middleware.Stack, opts Options) error {
	return stack.Initialize.Add(&presign{options: opts}, middleware.Before)
}

// RemoveMiddleware removes the Presign URL middleware from the stack.
func RemoveMiddleware(stack *middleware.Stack) error {
	return stack.Initialize.Remove((*presign)(nil).ID())
}

type presign struct {
	options Options
}

func (m *presign) ID() string { return id.Presigning }

func (m *presign) HandleInitialize(
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

	presignedURL, _, err := m.options.Presigner.PresignURL(ctx, srcRegion, paramCpy)
	if err != nil {
		return out, metadata, fmt.Errorf("unable to create presigned URL, %w", err)
	}

	// Update the original input with the presigned URL value.
	if err = m.options.Accessor.SetPresignedURL(input.Parameters, presignedURL); err != nil {
		return out, metadata, fmt.Errorf("presign middleware failed, %w", err)
	}

	return next.HandleInitialize(ctx, input)
}
