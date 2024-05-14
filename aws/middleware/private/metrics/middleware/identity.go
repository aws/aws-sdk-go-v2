package middleware

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
)

func timeGetIdentity(stack *middleware.Stack) error {
	if err := stack.Finalize.Insert(getIdentityStart{}, "GetIdentity", middleware.Before); err != nil {
		return err
	}
	if err := stack.Finalize.Insert(getIdentityEnd{}, "GetIdentity", middleware.After); err != nil {
		return err
	}
	return nil
}

type getIdentityStart struct{}

func (m getIdentityStart) ID() string { return "getIdentityStart" }

func (m getIdentityStart) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	out middleware.FinalizeOutput, md middleware.Metadata, err error,
) {
	mctx := metrics.Context(ctx)
	mctx.Data().GetIdentityStartTime = sdk.NowTime()
	return next.HandleFinalize(ctx, in)
}

type getIdentityEnd struct{}

func (m getIdentityEnd) ID() string { return "getIdentityEnd" }

func (m getIdentityEnd) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	out middleware.FinalizeOutput, md middleware.Metadata, err error,
) {
	mctx := metrics.Context(ctx)
	mctx.Data().GetIdentityEndTime = sdk.NowTime()
	return next.HandleFinalize(ctx, in)
}
