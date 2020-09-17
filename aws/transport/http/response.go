package http

import (
	"errors"
	"fmt"

	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// ResponseError provides the HTTP centric error type wrapping the underlying error
// with the HTTP response value and the deserialized RequestID.
type ResponseError struct {
	*smithyhttp.ResponseError
	RequestID string
}

func (e *ResponseError) ServiceRequestID() string { return e.RequestID }
func (e *ResponseError) Error() string {
	return fmt.Sprintf(
		"https response error StatusCode: %d, RequsetID: %s, %v",
		e.Response.StatusCode, e.RequestID, e.Err)
}

// As populates target and returns true if the type of target is a error type that
// the ResponseError embeeds, (e.g. Smithy's HTTP ResponseError)
func (e *ResponseError) As(target interface{}) bool {
	return errors.As(e.ResponseError, target)
}
