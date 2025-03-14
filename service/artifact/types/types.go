// Code generated by smithy-go-codegen DO NOT EDIT.

package types

import (
	smithydocument "github.com/aws/smithy-go/document"
	"time"
)

// Account settings for the customer.
type AccountSettings struct {

	// Notification subscription status of the customer.
	NotificationSubscriptionStatus NotificationSubscriptionStatus

	noSmithyDocumentSerde
}

// Summary for customer-agreement resource.
type CustomerAgreementSummary struct {

	// Terms required to accept the agreement resource.
	AcceptanceTerms []string

	// ARN of the agreement resource the customer-agreement resource represents.
	AgreementArn *string

	// ARN of the customer-agreement resource.
	Arn *string

	// AWS account Id that owns the resource.
	AwsAccountId *string

	// Description of the resource.
	Description *string

	// Timestamp indicating when the agreement was terminated.
	EffectiveEnd *time.Time

	// Timestamp indicating when the agreement became effective.
	EffectiveStart *time.Time

	// Identifier of the customer-agreement resource.
	Id *string

	// Name of the customer-agreement resource.
	Name *string

	// ARN of the organization that owns the resource.
	OrganizationArn *string

	// State of the resource.
	State CustomerAgreementState

	// Terms required to terminate the customer-agreement resource.
	TerminateTerms []string

	// Type of the customer-agreement resource.
	Type AgreementType

	noSmithyDocumentSerde
}

// Full detail for report resource metadata.
type ReportDetail struct {

	// Acceptance type for report.
	AcceptanceType AcceptanceType

	// ARN for the report resource.
	Arn *string

	// Category for the report resource.
	Category *string

	// Associated company name for the report resource.
	CompanyName *string

	// Timestamp indicating when the report resource was created.
	CreatedAt *time.Time

	// Timestamp indicating when the report resource was deleted.
	DeletedAt *time.Time

	// Description for the report resource.
	Description *string

	// Unique resource ID for the report resource.
	Id *string

	// Timestamp indicating when the report resource was last modified.
	LastModifiedAt *time.Time

	// Name for the report resource.
	Name *string

	// Timestamp indicating the report resource effective end.
	PeriodEnd *time.Time

	// Timestamp indicating the report resource effective start.
	PeriodStart *time.Time

	// Associated product name for the report resource.
	ProductName *string

	// Sequence number to enforce optimistic locking.
	SequenceNumber *int64

	// Series for the report resource.
	Series *string

	// Current state of the report resource
	State PublishedState

	// The message associated with the current upload state.
	StatusMessage *string

	// Unique resource ARN for term resource.
	TermArn *string

	// The current state of the document upload.
	UploadState UploadState

	// Version for the report resource.
	Version *int64

	noSmithyDocumentSerde
}

// Summary for report resource.
type ReportSummary struct {

	// Acceptance type for report.
	AcceptanceType AcceptanceType

	// ARN for the report resource.
	Arn *string

	// Category for the report resource.
	Category *string

	// Associated company name for the report resource.
	CompanyName *string

	// Description for the report resource.
	Description *string

	// Unique resource ID for the report resource.
	Id *string

	// Name for the report resource.
	Name *string

	// Timestamp indicating the report resource effective end.
	PeriodEnd *time.Time

	// Timestamp indicating the report resource effective start.
	PeriodStart *time.Time

	// Associated product name for the report resource.
	ProductName *string

	// Series for the report resource.
	Series *string

	// Current state of the report resource.
	State PublishedState

	// The message associated with the current upload state.
	StatusMessage *string

	// The current state of the document upload.
	UploadState UploadState

	// Version for the report resource.
	Version *int64

	noSmithyDocumentSerde
}

// Validation exception message and name.
type ValidationExceptionField struct {

	// Message describing why the field failed validation.
	//
	// This member is required.
	Message *string

	// Name of validation exception.
	//
	// This member is required.
	Name *string

	noSmithyDocumentSerde
}

type noSmithyDocumentSerde = smithydocument.NoSerde
