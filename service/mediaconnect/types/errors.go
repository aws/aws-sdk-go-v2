// Code generated by smithy-go-codegen DO NOT EDIT.

package types

import (
	"fmt"
	smithy "github.com/aws/smithy-go"
)

// Exception raised by Elemental MediaConnect when adding the flow output. See the
// error message for the operation for more information on the cause of this
// exception.
type AddFlowOutputs420Exception struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *AddFlowOutputs420Exception) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *AddFlowOutputs420Exception) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *AddFlowOutputs420Exception) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "AddFlowOutputs420Exception"
	}
	return *e.ErrorCodeOverride
}
func (e *AddFlowOutputs420Exception) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// This exception is thrown if the request contains a semantic error. The precise
// meaning depends on the API, and is documented in the error message.
type BadRequestException struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *BadRequestException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *BadRequestException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *BadRequestException) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "BadRequestException"
	}
	return *e.ErrorCodeOverride
}
func (e *BadRequestException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The requested operation would cause a conflict with the current state of a
// service resource associated with the request. Resolve the conflict before
// retrying this request.
type ConflictException struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *ConflictException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *ConflictException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *ConflictException) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "ConflictException"
	}
	return *e.ErrorCodeOverride
}
func (e *ConflictException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// Exception raised by Elemental MediaConnect when creating the bridge. See the
// error message for the operation for more information on the cause of this
// exception.
type CreateBridge420Exception struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *CreateBridge420Exception) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *CreateBridge420Exception) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *CreateBridge420Exception) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "CreateBridge420Exception"
	}
	return *e.ErrorCodeOverride
}
func (e *CreateBridge420Exception) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// Exception raised by Elemental MediaConnect when creating the flow. See the
// error message for the operation for more information on the cause of this
// exception.
type CreateFlow420Exception struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *CreateFlow420Exception) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *CreateFlow420Exception) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *CreateFlow420Exception) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "CreateFlow420Exception"
	}
	return *e.ErrorCodeOverride
}
func (e *CreateFlow420Exception) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// Exception raised by Elemental MediaConnect when creating the gateway. See the
// error message for the operation for more information on the cause of this
// exception.
type CreateGateway420Exception struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *CreateGateway420Exception) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *CreateGateway420Exception) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *CreateGateway420Exception) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "CreateGateway420Exception"
	}
	return *e.ErrorCodeOverride
}
func (e *CreateGateway420Exception) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// You do not have sufficient access to perform this action.
type ForbiddenException struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *ForbiddenException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *ForbiddenException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *ForbiddenException) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "ForbiddenException"
	}
	return *e.ErrorCodeOverride
}
func (e *ForbiddenException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// Exception raised by Elemental MediaConnect when granting the entitlement. See
// the error message for the operation for more information on the cause of this
// exception.
type GrantFlowEntitlements420Exception struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *GrantFlowEntitlements420Exception) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *GrantFlowEntitlements420Exception) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *GrantFlowEntitlements420Exception) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "GrantFlowEntitlements420Exception"
	}
	return *e.ErrorCodeOverride
}
func (e *GrantFlowEntitlements420Exception) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The server encountered an internal error and is unable to complete the request.
type InternalServerErrorException struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *InternalServerErrorException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *InternalServerErrorException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *InternalServerErrorException) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "InternalServerErrorException"
	}
	return *e.ErrorCodeOverride
}
func (e *InternalServerErrorException) ErrorFault() smithy.ErrorFault { return smithy.FaultServer }

// One or more of the resources in the request does not exist in the system.
type NotFoundException struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *NotFoundException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *NotFoundException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *NotFoundException) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "NotFoundException"
	}
	return *e.ErrorCodeOverride
}
func (e *NotFoundException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The service is currently unavailable or busy.
type ServiceUnavailableException struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *ServiceUnavailableException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *ServiceUnavailableException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *ServiceUnavailableException) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "ServiceUnavailableException"
	}
	return *e.ErrorCodeOverride
}
func (e *ServiceUnavailableException) ErrorFault() smithy.ErrorFault { return smithy.FaultServer }

// The request was denied due to request throttling.
type TooManyRequestsException struct {
	Message *string

	ErrorCodeOverride *string

	noSmithyDocumentSerde
}

func (e *TooManyRequestsException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *TooManyRequestsException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *TooManyRequestsException) ErrorCode() string {
	if e == nil || e.ErrorCodeOverride == nil {
		return "TooManyRequestsException"
	}
	return *e.ErrorCodeOverride
}
func (e *TooManyRequestsException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }
