package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// ResponseError provides the HTTP centric error type wrapping the underlying
// error with the HTTP response value and the deserialized RequestID.
type ResponseError struct {
	*smithyhttp.ResponseError
	RequestID string
}

// ServiceRequestID returns the HTTP response wrapping the underlying error.
func (e *ResponseError) ServiceRequestID() string { return e.RequestID }

func (e *ResponseError) Error() string {
	return fmt.Sprintf(
		"http response error StatusCode: %d, RequestID: %s, %v",
		e.Response.StatusCode, e.RequestID, e.Err)
}

// As populates target and returns true if the type of target is a error type
// that the ResponseError embeds, (e.g. Smithy's HTTP ResponseError)
func (e *ResponseError) As(target interface{}) bool {
	return errors.As(e.ResponseError, target)
}

// Client provides the interface for the minimum HTTP client behavior to
// invoke a HTTP request.
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientHandler provides a smithy middleware handler wrapper for an standard
// HTTP client.
type ClientHandler struct {
	Client Client
}

// Handle invokes the HTTP client with the provided HTTP request. Returns the
// response, or error if the request failed.
func (h ClientHandler) Handle(ctx context.Context, input interface{}) (
	out interface{}, metadata middleware.Metadata, err error,
) {
	sReq, ok := input.(*smithyhttp.Request)
	if !ok {
		return nil, metadata, fmt.Errorf(
			"invalid input type for HTTP client handlers, expect %T, got %T",
			input, (*smithyhttp.Request)(nil))
	}

	hReq := sReq.Build(ctx)

	resp, err := h.Client.Do(hReq)
	return resp, metadata, err
}
