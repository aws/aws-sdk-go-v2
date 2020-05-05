package lexruntimeservice

import (
	"fmt"

	"github.com/awslabs/smithy-go"
)

// InvalidParameterExceptionInterface provides the interface for InvalidParameterException
// to be extended in the future without a breaking change.
type InvalidParameterExceptionInterface interface {
	smithy.APIError
	isInvalidParameterException()
}

// InvalidParameterException is an example of an modeled exception with additional
// members beyond the message and fault.
type InvalidParameterException struct {
	Message string
}

func (e *InvalidParameterException) isInvalidParameterException() {}

// ErrorCode returns the code for the API exception.
func (e *InvalidParameterException) ErrorCode() string { return "InvalidParameterException" }

// ErrorMessage returns the message from the exception.
func (e *InvalidParameterException) ErrorMessage() string { return e.Message }

// ErrorFault returns the source of the error.
func (e *InvalidParameterException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }
func (e *InvalidParameterException) Error() string {
	return fmt.Sprintf("api error %s: %s", e.ErrorCode(), e.ErrorMessage())
}

var _ smithy.APIError = (*InvalidParameterException)(nil)

var _ InvalidParameterExceptionInterface = (*InvalidParameterException)(nil)
