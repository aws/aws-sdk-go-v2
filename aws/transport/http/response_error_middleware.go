package http

import (
	"context"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/id"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// AddResponseErrorMiddleware adds response error wrapper middleware
func AddResponseErrorMiddleware(stack *middleware.Stack) error {
	// add error wrapper middleware before request id retriever middleware so that it can wrap the error response
	// returned by operation deserializers
	return stack.Deserialize.Add(&errorWrapperMiddleware{}, middleware.Before)
}

type errorWrapperMiddleware struct {
}

// ID returns the middleware identifier
func (m *errorWrapperMiddleware) ID() string {
	return id.ResponseErrorWrapper
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

	// look for request id in metadata
	reqID, _ := awsmiddleware.GetRequestIDMetadata(metadata)

	// Wrap the returned smithy error with the request id retrieved from the metadata
	err = &ResponseError{
		ResponseError: &smithyhttp.ResponseError{
			Response: resp,
			Err:      err,
		},
		RequestID: reqID,
	}

	return out, metadata, err
}
