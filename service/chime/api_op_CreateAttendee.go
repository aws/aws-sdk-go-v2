// Code generated by smithy-go-codegen DO NOT EDIT.

package chime

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/chime/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

//	Creates a new attendee for an active Amazon Chime SDK meeting. For more
//
// information about the Amazon Chime SDK, see [Using the Amazon Chime SDK]in the Amazon Chime SDK Developer
// Guide.
//
// This API is is no longer supported and will not be updated. We recommend using
// the latest version, [CreateAttendee], in the Amazon Chime SDK.
//
// Using the latest version requires migrating to a dedicated namespace. For more
// information, refer to [Migrating from the Amazon Chime namespace]in the Amazon Chime SDK Developer Guide.
//
// Deprecated: Replaced by CreateAttendee in the Amazon Chime SDK Meetings
// Namespace
//
// [Migrating from the Amazon Chime namespace]: https://docs.aws.amazon.com/chime-sdk/latest/dg/migrate-from-chm-namespace.html
// [Using the Amazon Chime SDK]: https://docs.aws.amazon.com/chime-sdk/latest/dg/meetings-sdk.html
// [CreateAttendee]: https://docs.aws.amazon.com/chime-sdk/latest/APIReference/API_meeting-chime_CreateAttendee.html
func (c *Client) CreateAttendee(ctx context.Context, params *CreateAttendeeInput, optFns ...func(*Options)) (*CreateAttendeeOutput, error) {
	if params == nil {
		params = &CreateAttendeeInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateAttendee", params, optFns, c.addOperationCreateAttendeeMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateAttendeeOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateAttendeeInput struct {

	// The Amazon Chime SDK external user ID. An idempotency token. Links the attendee
	// to an identity managed by a builder application.
	//
	// This member is required.
	ExternalUserId *string

	// The Amazon Chime SDK meeting ID.
	//
	// This member is required.
	MeetingId *string

	// The tag key-value pairs.
	Tags []types.Tag

	noSmithyDocumentSerde
}

type CreateAttendeeOutput struct {

	// The attendee information, including attendee ID and join token.
	Attendee *types.Attendee

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateAttendeeMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpCreateAttendee{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpCreateAttendee{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CreateAttendee"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addCredentialSource(stack, options); err != nil {
		return err
	}
	if err = addOpCreateAttendeeValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateAttendee(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opCreateAttendee(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CreateAttendee",
	}
}
