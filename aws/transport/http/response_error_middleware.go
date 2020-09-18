package http

import (
	"context"

	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// AddResponseErrorWrapper adds response error wrapper middleware
func AddResponseErrorWrapper(stack *middleware.Stack) {
	// add error wrapper middleware before operation deserializers so that it can wrap the error response
	// returned by operation deserializers
	stack.Deserialize.Insert(&errorWrapperMiddleware{}, "OperationDeserializer", middleware.Before)
}

type errorWrapperMiddleware struct {
}

// ID returns the middleware identifier
func (m *errorWrapperMiddleware) ID() string {
	return "AWSResponseErrorWrapperMiddleware"
}

func (m *errorWrapperMiddleware) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err == nil {
		// Nothing to do when there is no error.
		return out, metadata, err
	}

	resp, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok {
		// No raw response to wrap with.
		return out, metadata, err
	}

	// Wrap the returned smithy error with the request id retrieved from the metadata
	if err != nil {
		err = &ResponseError{
			ResponseError: &smithyhttp.ResponseError{
				Response: resp,
				Err:      err,
			},
			RequestID: GetRequestIDMetadata(metadata),
		}
	}

	return out, metadata, err
}
