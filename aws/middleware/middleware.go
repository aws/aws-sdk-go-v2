package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/awslabs/smithy-go/middleware"
	smithyHTTP "github.com/awslabs/smithy-go/transport/http"
)

// RequestInvocationIDMiddleware is a Smithy BuildMiddleware that will generate a unique ID for logical API operation
// invocation.
type RequestInvocationIDMiddleware struct{}

// ID the identifier for the RequestInvocationIDMiddleware
func (r RequestInvocationIDMiddleware) ID() string {
	return "RequestInvocationIDMiddleware"
}

// HandleBuild attaches a unique operation invocation id for the operation to the request
func (r RequestInvocationIDMiddleware) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	const invocationIDHeader = "amz-sdk-invocation-id"

	invocationID, err := sdk.UUIDVersion4()
	if err != nil {
		return out, metadata, err
	}

	switch req := in.Request.(type) {
	case *smithyHTTP.Request:
		req.Header.Set(invocationIDHeader, invocationID)
	default:
		return out, metadata, fmt.Errorf("unknown transport type %T", req)
	}

	return next.HandleBuild(ctx, in)
}

// AttemptClockSkewMiddleware calculates the clock skew of the SDK client
// TODO: Could be a better name, since this calculates more then skew
type AttemptClockSkewMiddleware struct{}

// ID is the middleware identifier
func (a AttemptClockSkewMiddleware) ID() string {
	return "AttemptClockSkewMiddlware"
}

// HandleDeserialize calculates response metadata and clock skew
func (a AttemptClockSkewMiddleware) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	respMeta := ResponseMetadata{}

	out, metadata, err = next.HandleDeserialize(ctx, in)
	respMeta.ResponseAt = sdk.NowTime()

	switch resp := out.RawResponse.(type) {
	case *smithyHTTP.Response:
		respDateHeader := resp.Header.Get("Date")
		if len(respDateHeader) == 0 {
			break
		}
		var parseErr error
		respMeta.ServerTime, parseErr = http.ParseTime(respDateHeader)
		if parseErr != nil {
			// TODO: What should logging of errors look like?
			break
		}
	}

	if !respMeta.ServerTime.IsZero() {
		respMeta.AttemptSkew = respMeta.ServerTime.Sub(respMeta.ResponseAt)
	}

	setResponseMetadata(&metadata, respMeta)

	return out, metadata, err
}

type responseMetadataKey struct{}

// ResponseMetadata is metadata about the transport layer response
type ResponseMetadata struct {
	ResponseAt  time.Time
	ServerTime  time.Time
	AttemptSkew time.Duration
}

// GetResponseMetadata retrieves response metadata from the context, if nil returns an empty value
func GetResponseMetadata(metadata middleware.Metadata) (v ResponseMetadata) {
	v, _ = metadata.Get(responseMetadataKey{}).(ResponseMetadata)
	return v
}

// setResponseMetadata sets the ResponseMetadata on the given context
func setResponseMetadata(metadata *middleware.Metadata, responseMetadata ResponseMetadata) {
	metadata.Set(responseMetadataKey{}, responseMetadata)
}
