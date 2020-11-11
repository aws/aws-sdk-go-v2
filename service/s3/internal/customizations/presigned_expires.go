package customizations

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// AddExpiresOnPresignedURL  represents a build middleware used to assign
// expiration on a presigned URL
type AddExpiresOnPresignedURL struct {
	// Expires is time.Duration within which presigned url should be expired
	Expires time.Duration
}

// ID representing the middleware
func (*AddExpiresOnPresignedURL) ID() string {
	return "S3:AddExpiresOnPresignedURL"
}

// HandleBuild handles the build step middleware behavior
func (m *AddExpiresOnPresignedURL) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	// if expiration is unset skip this middleware
	if m.Expires == 0 {
		return next.HandleBuild(ctx, in)
	}

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", req)
	}

	// set S3 X-AMZ-Expires header
	query := req.URL.Query()
	query.Set("X-Amz-Expires", strconv.FormatInt(int64(m.Expires/time.Second), 10))
	req.URL.RawQuery = query.Encode()

	return next.HandleBuild(ctx, in)
}
