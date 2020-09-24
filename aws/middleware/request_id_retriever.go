package middleware

import (
	"context"

	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// AddRequestIDRetrieverMiddleware adds request id retriever middleware
func AddRequestIDRetrieverMiddleware(stack *middleware.Stack) {
	// add error wrapper middleware before operation deserializers so that it can wrap the error response
	// returned by operation deserializers
	stack.Deserialize.Insert(&requestIDRetrieverMiddleware{}, "OperationDeserializer", middleware.Before)
}

type requestIDRetrieverMiddleware struct {
}

// ID returns the middleware identifier
func (m *requestIDRetrieverMiddleware) ID() string {
	return "AWSRequestIDRetrieverMiddleware"
}

func (m *requestIDRetrieverMiddleware) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)

	resp, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok {
		// No raw response to wrap with.
		return out, metadata, err
	}

	// Different header which can map to request id
	requestIDHeaderList := []string{"X-Amzn-Requestid", "X-Amz-RequestId"}

	for _, h := range requestIDHeaderList {
		// check for headers known to contain Request id
		if v := resp.Header.Get(h); len(v) != 0 {
			// set reqID on metadata for successful responses.
			SetRequestIDMetadata(&metadata, v)
			break
		}
	}

	return out, metadata, err
}
