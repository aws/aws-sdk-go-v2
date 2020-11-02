package customizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/internal/s3shared"
)

// processOutpostIDMiddleware is special customization middleware to be applied for operations
// CreateBucket, ListRegionalBuckets which must resolve endpoint to s3-outposts.{region}.amazonaws.com
// with region as client region and signed by s3-control if an outpost id is provided.
type processOutpostIDMiddleware struct {
	// GetOutpostID points to a function that processes an input and returns an outpostID as string ptr,
	// and bool indicating if outpostID is supported or set.
	GetOutpostID func(interface{}) (*string, bool)

	// UseDualStack indicates of dual stack endpoints should be used
	UseDualstack bool
}

// ID returns the middleware ID.
func (*processOutpostIDMiddleware) ID() string { return "S3Control:ProcessOutpostIDMiddleware" }

// HandleSerialize adds a serialize step, this has to be before operation serializer and arn endpoint logic.
// Ideally this step will be ahead of ARN customization for CreateBucket, ListRegionalBucket operation.
func (m *processOutpostIDMiddleware) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	// if host name is immutable, skip this customization
	if smithyhttp.GetHostnameImmutable(ctx) {
		return next.HandleSerialize(ctx, in)
	}

	// attempt to fetch an outpost id
	outpostID, ok := m.GetOutpostID(in.Parameters)
	if !ok {
		return next.HandleSerialize(ctx, in)
	}

	// check if outpostID was not set or is empty
	if outpostID == nil || len(strings.TrimSpace(*outpostID)) == 0 {
		return next.HandleSerialize(ctx, in)
	}

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	serviceEndpointLabel := "s3-outposts."
	requestRegion := awsmiddleware.GetRegion(ctx)

	// validate if fips
	if s3shared.IsFIPS(requestRegion) {
		return out, metadata, fmt.Errorf("unsupported fips region provided for outposts request")
	}
	// validate if dualstack
	if m.UseDualstack {
		return out, metadata, fmt.Errorf("dualstack is not supported for outposts request")
	}

	// set request url
	req.URL.Host = serviceEndpointLabel + requestRegion + ".amazonaws.com"

	// Disable endpoint host prefix for s3-control
	ctx = smithyhttp.DisableEndpointHostPrefix(ctx, true)

	// redirect signer
	ctx = awsmiddleware.SetSigningName(ctx, "s3-outposts")
	ctx = awsmiddleware.SetSigningRegion(ctx, requestRegion)

	return next.HandleSerialize(ctx, in)
}
