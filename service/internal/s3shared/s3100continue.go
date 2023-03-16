package s3shared

import (
	"context"
	"fmt"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const s3100ContinueID = "S3100Continue"
const defaultLimit int64 = 1024 * 1024 * 2

// Add100Continue add middleware, which adds {Expect: 100-continue} header for s3 client HTTP PUT request larger than 2MB
// or with unknown size streaming bodies, during operation builder step
func Add100Continue(stack *middleware.Stack, option AddContinueOption) error {
	return stack.Build.Add(&s3100Continue{
		continueHeaderThresholdBytes: option.ContinueHeaderThresholdBytes,
	}, middleware.After)
}

type s3100Continue struct {
	continueHeaderThresholdBytes int64
}

// ID returns the middleware identifier
func (m *s3100Continue) ID() string {
	return s3100ContinueID
}

func (m *s3100Continue) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	sizeLimit := defaultLimit
	switch {
	case m.continueHeaderThresholdBytes == -1:
		return next.HandleBuild(ctx, in)
	case m.continueHeaderThresholdBytes > 0:
		sizeLimit = m.continueHeaderThresholdBytes
	default:
	}

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	if req.ContentLength == -1 || (req.ContentLength == 0 && req.Body != nil) || req.ContentLength >= sizeLimit {
		req.Header.Set("Expect", "100-continue")
	}

	return next.HandleBuild(ctx, in)
}

// AddContinueOption passes user configuration of threshold size to trigger 100-continue header
type AddContinueOption struct {
	ContinueHeaderThresholdBytes int64
}
