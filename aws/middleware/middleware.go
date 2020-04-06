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
	return "Request Invocation ID Middleware"
}

// HandleBuild attaches a unique operation invocation id for the operation to the request
func (r RequestInvocationIDMiddleware) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (out middleware.BuildOutput, metadata middleware.Metadata, err error) {
	const invocationIDHeader = "amz-sdk-invocation-id"

	invocationID, err := sdk.UUIDVersion4()
	if err != nil {
		return out, middleware.NewMetadata(), err
	}

	switch req := in.Request.(type) {
	case *smithyHTTP.Request:
		req.Header.Set(invocationIDHeader, invocationID)
	default:
		return middleware.BuildOutput{}, middleware.NewMetadata(), fmt.Errorf("unknown transport type %T", req)
	}

	return next.HandleBuild(ctx, in)
}

// AttemptClockSkewMiddleware calculates the clock skew of the SDK client
// TODO: Could be a better name, since this calculates more then skew
type AttemptClockSkewMiddleware struct{}

func (a AttemptClockSkewMiddleware) ID() string {
	return "Attempt Clock Skew Middlware"
}

// HandleDeserialize calculates response metadata and clock skew
func (a AttemptClockSkewMiddleware) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (out middleware.DeserializeOutput, metadata middleware.Metadata, err error) {
	respMeta := ResponseMetadata{}

	deserialize, metadata, err := next.HandleDeserialize(ctx, in)
	respMeta.ResponseAt = sdk.NowTime()

	switch resp := deserialize.RawResponse.(type) {
	case *smithyHTTP.Response:
		respDateHeader := resp.Header.Get("Date")
		if len(respDateHeader) == 0 {
			break
		}
		respMeta.ServerTime, err = http.ParseTime(respDateHeader)
		if err != nil {
			// TODO: What should logging of errors look like?
			break
		}
	}

	if !respMeta.ServerTime.IsZero() {
		respMeta.AttemptSkew = respMeta.ServerTime.Sub(respMeta.ResponseAt)
	}

	SetResponseMetadata(metadata, respMeta)

	return deserialize, metadata, err
}

type responseMetadataKey struct{}

// ResponseMetadata is metadata about the transport layer response
type ResponseMetadata struct {
	ResponseAt  time.Time
	ServerTime  time.Time
	AttemptSkew time.Duration
}

// GetResponseMetadata retrieves response metadata from the context, if nil returns an empty value
func GetResponseMetadata(metadata middleware.Metadata) ResponseMetadata {
	switch v := metadata.Get(responseMetadataKey{}).(type) {
	case ResponseMetadata:
		return v
	default:
		return ResponseMetadata{}
	}
}

// SetResponseMetadata sets the ResponseMetadata on the given context
func SetResponseMetadata(metadata middleware.Metadata, responseMetadata ResponseMetadata) {
	metadata.Set(responseMetadataKey{}, responseMetadata)
}
