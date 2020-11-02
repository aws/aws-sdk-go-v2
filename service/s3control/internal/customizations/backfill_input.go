package customizations

import (
	"context"
	"fmt"

	"github.com/awslabs/smithy-go/middleware"

	"github.com/aws/aws-sdk-go-v2/service/internal/s3shared"
)

// BackfillInputMiddleware validates and backfill's values from ARN into request serializable input.
// This middleware must be executed after `InitARNLookupMiddleware` and before `inputValidationMiddleware`.
type BackfillInputMiddleware struct {

	// BackfillAccountID points to a function that validates the input for accountID. If absent, it populates the
	// accountID and returns a copy. If present, but different than passed in accountID value throws an error
	BackfillAccountID func(interface{}, string) (interface{}, error)
}

// ID representing the middleware
func (m *BackfillInputMiddleware) ID() string {
	return "S3Control:BackfillInputMiddleware"
}

// HandleInitialize handles the middleware behavior in an Initialize step.
func (m *BackfillInputMiddleware) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	// fetch arn from context
	av, ok := s3shared.GetARNResourceFromContext(ctx)
	if !ok {
		return next.HandleInitialize(ctx, in)
	}

	input, err := m.BackfillAccountID(in.Parameters, av.AccountID)
	if err != nil {
		return out, metadata, fmt.Errorf("invalid ARN, %w", err)
	}

	// assign the modified copy of input to in.Parameters
	in.Parameters = input

	return next.HandleInitialize(ctx, in)
}
