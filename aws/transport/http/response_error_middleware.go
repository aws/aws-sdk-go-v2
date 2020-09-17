package http

import (
	"context"
	"errors"

	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// AddResponseErrorWrapper adds response error wrapper middleware
func AddResponseErrorWrapper(stack *middleware.Stack) {
	// add error wrapper middleware before operation deserializers so that it can wrap the error response
	// returned by operation deserializers
	stack.Deserialize.Insert(&errorWrapperMiddleware{}, "OperationDeserializer", middleware.Before)
}

// RequestID is the identifier for Request Id set in response metadata
const RequestID = "RequestID"

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

	_, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok {
		// No raw response to wrap with.
		return out, metadata, err
	}

	var reqID string
	// TODO: modify protocol deserializers to set RequestID on metadata when
	//  available in raw response.
	if metadata.Has(RequestID) {
		if v, ok := metadata.Get(RequestID).(string); ok {
			reqID = v
		}
	}

	// Wrap the returned smithy error with the request id retrieved from the metadata
	if err != nil {
		var respErr *smithyhttp.ResponseError
		if errors.As(err, &respErr) {
			err = &ResponseError{
				ResponseError: respErr,
				RequestID:     reqID,
			}
		}
	}

	return out, metadata, err
}
