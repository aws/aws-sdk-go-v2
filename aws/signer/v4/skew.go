package v4

import (
	"context"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
)

// DetectSkewMiddleware applies clock skew to instances of Signer.
type DetectSkewMiddleware struct {
	Signer *Signer
}

// ID identifies DetectSkewMiddleware.
func (*DetectSkewMiddleware) ID() string {
	return "aws.signer.v4#DetectSkew"
}

// HandleFinalize applies clock skew.
func (m *DetectSkewMiddleware) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleFinalize(ctx, in)
	if skew, ok := awsmiddleware.GetAttemptSkew(metadata); ok {
		m.Signer.clockSkew.Store(int64(skew))
	}

	return out, metadata, err
}
