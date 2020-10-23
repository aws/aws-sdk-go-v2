package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/id"
	"github.com/aws/aws-sdk-go-v2/internal/rand"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/awslabs/smithy-go/middleware"
	smithymwid "github.com/awslabs/smithy-go/middleware/id"
	smithyrand "github.com/awslabs/smithy-go/rand"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// ClientRequestID is a Smithy BuildMiddleware that will generate a unique ID for logical API operation
// invocation.
type ClientRequestID struct{}

// ID the identifier for the ClientRequestID
func (r *ClientRequestID) ID() string {
	return id.ClientRequestID
}

// HandleBuild attaches a unique operation invocation id for the operation to the request
func (r ClientRequestID) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", req)
	}

	invocationID, err := smithyrand.NewUUID(rand.Reader).GetUUID()
	if err != nil {
		return out, metadata, err
	}

	const invocationIDHeader = "Amz-Sdk-Invocation-Id"
	req.Header[invocationIDHeader] = append(req.Header[invocationIDHeader][:0], invocationID)

	return next.HandleBuild(ctx, in)
}

// AttemptClockSkewMiddleware calculates the clock skew of the SDK client
// TODO: Could be a better name, since this calculates more then skew
type AttemptClockSkewMiddleware struct{}

// ID is the middleware identifier
func (a *AttemptClockSkewMiddleware) ID() string {
	return "AttemptClockSkewMiddleware"
}

// HandleDeserialize calculates response metadata and clock skew
func (a AttemptClockSkewMiddleware) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	respMeta := ResponseMetadata{}

	out, metadata, err = next.HandleDeserialize(ctx, in)
	respMeta.ResponseAt = sdk.NowTime()

	switch resp := out.RawResponse.(type) {
	case *smithyhttp.Response:
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

// AddClientRequestID adds ClientRequestID to the middleware stack
func AddClientRequestIDMiddleware(stack *middleware.Stack) error {
	m := &ClientRequestID{}
	return stack.Build.Add(m, middleware.After)
}

// AddAttemptClockSkewMiddleware adds AttemptClockSkewMiddleware to the middleware stack
func AddAttemptClockSkewMiddleware(stack *middleware.Stack) error {
	attemptClockSkeyMiddleware := &AttemptClockSkewMiddleware{}
	return stack.Deserialize.Insert(attemptClockSkeyMiddleware, smithymwid.OperationDeserializer, middleware.After)
}
