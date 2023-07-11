package s3

import (
	"context"
	"fmt"
	"net/url"
	"path"

	smithy "github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// serializeImmutableHostnameBucketMiddleware handles injecting the bucket name into
// "immutable" hostnames resolved via v1 EndpointResolvers. This CANNOT be done in
// serialization, since v2 endpoint resolution requires removing the {Bucket} path
// segment from all S3 requests.
type serializeImmutableHostnameBucketMiddleware struct{}

func (*serializeImmutableHostnameBucketMiddleware) ID() string {
	return "serializeImmutableHostnameBucket"
}

func (m *serializeImmutableHostnameBucketMiddleware) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	request, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, &smithy.SerializationError{Err: fmt.Errorf("unknown transport type %T", in.Request)}
	}
	if !smithyhttp.GetHostnameImmutable(ctx) {
		return next.HandleSerialize(ctx, in)
	}

	if bucket, ok := bucketFromInput(in.Parameters); ok {
		request.URL.Path = path.Join(request.URL.Path, bucket)
		request.URL.RawPath = path.Join(request.URL.RawPath, url.PathEscape(bucket))
	}

	return next.HandleSerialize(ctx, in)
}

func addSerializeImmutableHostnameBucketMiddleware(stack *middleware.Stack) error {
	return stack.Serialize.Insert(
		&serializeImmutableHostnameBucketMiddleware{},
		"ResolveEndpointV2",
		middleware.After,
	)
}
