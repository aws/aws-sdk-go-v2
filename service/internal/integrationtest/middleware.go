package integrationtest

import "github.com/Enflick/smithy-go/middleware"

// RemoveOperationInputValidationMiddleware removes the validation middleware
// from the stack.
func RemoveOperationInputValidationMiddleware(stack *middleware.Stack) error {
	stack.Initialize.Remove("OperationInputValidation")
	return nil
}
