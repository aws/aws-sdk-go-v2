package middleware

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/smithy-go/middleware"
)

// WrapSendErrorMiddleware provides a deserialize stack middleware that will
// wrap the error returned by the HTTP client with a ResponseError type. This
// wrapping distinguishes the requested failed due to a connection error with
// the HTTP client unable to successfully make a request.
//
// Also wraps the underlying error with a RequestCanceledError if the
// operation's `context` was canceled.
type WrapSendErrorMiddleware struct{}

// ID is the ID if the middleware.
func (*WrapSendErrorMiddleware) ID() string { return "WrapSendErrorMiddleware" }

// HandleDeserialize implements the DeserializeMiddleware interface.
func (m *WrapSendErrorMiddleware) HandleDeserialize(
	ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler,
) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		err = &aws.RequestSendError{Err: err}

		// Override the error with a context canceled error, if that was canceled.
		select {
		case <-ctx.Done():
			err = &aws.RequestCanceledError{Err: ctx.Err()}
		default:
		}
	}

	return out, metadata, err
}
