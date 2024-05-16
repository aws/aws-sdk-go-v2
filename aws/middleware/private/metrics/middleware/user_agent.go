package middleware

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type captureUserAgent struct{}

func (*captureUserAgent) ID() string { return "captureUserAgent" }

func (*captureUserAgent) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, md middleware.Metadata, err error,
) {
	r, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, md, fmt.Errorf("unexpected transport type %T", in.Request)
	}

	mctx := metrics.Context(ctx)
	mctx.Data().UserAgent = r.Header.Get("User-Agent")
	return next.HandleBuild(ctx, in)
}
