package customizations

import (
	"context"
	"fmt"
	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"
)

// removeBucketFromPathMiddleware needs to be executed after serialize step is performed
type removeBucketFromPathMiddleware struct {
}

func (m *removeBucketFromPathMiddleware) ID() string {
	return "S3:RemoveBucketFromPathMiddleware"
}

func (m *removeBucketFromPathMiddleware) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	// check if a bucket removal from HTTP path is required
	bucket, ok := getRemoveBucketFromPath(ctx)
	if !ok {
		return next.HandleSerialize(ctx, in)
	}

	req, ok := in.Request.(*http.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	removeBucketFromPath(req.URL, bucket)
	return next.HandleSerialize(ctx, in)
}

type removeBucketKey struct {
	bucket string
}

func setBucketToRemoveOnContext(ctx context.Context, bucket string) context.Context {
	return context.WithValue(ctx, removeBucketKey{}, bucket)
}

func getRemoveBucketFromPath(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(removeBucketKey{}).(string)
	return v, ok
}
