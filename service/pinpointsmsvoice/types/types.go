// Code generated by smithy-go-codegen DO NOT EDIT.

package types

import (
	smithydocument "github.com/aws/smithy-go/document"
)

// An object that defines a message that contains text formatted using Amazon
// Pinpoint Voice Instructions markup.
type CallInstructionsMessageType struct {

	// The language to use when delivering the message. For a complete list of
	// supported languages, see the Amazon Polly Developer Guide.
	Text *string

	noSmithyDocumentSerde
}

// An object that contains information about an event destination that sends data
// to Amazon CloudWatch Logs.
type CloudWatchLogsDestination struct {

	// The Amazon Resource Name (ARN) of an Amazon Identity and Access Management
	// (IAM) role that is able to write event data to an Amazon CloudWatch destination.
	IamRoleArn *string

	// The name of the Amazon CloudWatch Log Group that you want to record events in.
	LogGroupArn *string

	noSmithyDocumentSerde
}

// An object that defines an event destination.
type EventDestination struct {

	// An object that contains information about an event destination that sends data
	// to Amazon CloudWatch Logs.
	CloudWatchLogsDestination *CloudWatchLogsDestination

	// Indicates whether or not the event destination is enabled. If the event
	// destination is enabled, then Amazon Pinpoint sends response data to the
	// specified event destination.
	Enabled *bool

	// An object that contains information about an event destination that sends data
	// to Amazon Kinesis Data Firehose.
	KinesisFirehoseDestination *KinesisFirehoseDestination

	// An array of EventDestination objects. Each EventDestination object includes
	// ARNs and other information that define an event destination.
	MatchingEventTypes []EventType

	// A name that identifies the event destination configuration.
	Name *string

	// An object that contains information about an event destination that sends data
	// to Amazon SNS.
	SnsDestination *SnsDestination

	noSmithyDocumentSerde
}

// An object that defines a single event destination.
type EventDestinationDefinition struct {

	// An object that contains information about an event destination that sends data
	// to Amazon CloudWatch Logs.
	CloudWatchLogsDestination *CloudWatchLogsDestination

	// Indicates whether or not the event destination is enabled. If the event
	// destination is enabled, then Amazon Pinpoint sends response data to the
	// specified event destination.
	Enabled *bool

	// An object that contains information about an event destination that sends data
	// to Amazon Kinesis Data Firehose.
	KinesisFirehoseDestination *KinesisFirehoseDestination

	// An array of EventDestination objects. Each EventDestination object includes
	// ARNs and other information that define an event destination.
	MatchingEventTypes []EventType

	// An object that contains information about an event destination that sends data
	// to Amazon SNS.
	SnsDestination *SnsDestination

	noSmithyDocumentSerde
}

// An object that contains information about an event destination that sends data
// to Amazon Kinesis Data Firehose.
type KinesisFirehoseDestination struct {

	// The Amazon Resource Name (ARN) of an IAM role that can write data to an Amazon
	// Kinesis Data Firehose stream.
	DeliveryStreamArn *string

	// The Amazon Resource Name (ARN) of the Amazon Kinesis Data Firehose destination
	// that you want to use in the event destination.
	IamRoleArn *string

	noSmithyDocumentSerde
}

// An object that defines a message that contains unformatted text.
type PlainTextMessageType struct {

	// The language to use when delivering the message. For a complete list of
	// supported languages, see the Amazon Polly Developer Guide.
	LanguageCode *string

	// The plain (not SSML-formatted) text to deliver to the recipient.
	Text *string

	// The name of the voice that you want to use to deliver the message. For a
	// complete list of supported voices, see the Amazon Polly Developer Guide.
	VoiceId *string

	noSmithyDocumentSerde
}

// An object that contains information about an event destination that sends data
// to Amazon SNS.
type SnsDestination struct {

	// The Amazon Resource Name (ARN) of the Amazon SNS topic that you want to publish
	// events to.
	TopicArn *string

	noSmithyDocumentSerde
}

// An object that defines a message that contains SSML-formatted text.
type SSMLMessageType struct {

	// The language to use when delivering the message. For a complete list of
	// supported languages, see the Amazon Polly Developer Guide.
	LanguageCode *string

	// The SSML-formatted text to deliver to the recipient.
	Text *string

	// The name of the voice that you want to use to deliver the message. For a
	// complete list of supported voices, see the Amazon Polly Developer Guide.
	VoiceId *string

	noSmithyDocumentSerde
}

// An object that contains a voice message and information about the recipient
// that you want to send it to.
type VoiceMessageContent struct {

	// An object that defines a message that contains text formatted using Amazon
	// Pinpoint Voice Instructions markup.
	CallInstructionsMessage *CallInstructionsMessageType

	// An object that defines a message that contains unformatted text.
	PlainTextMessage *PlainTextMessageType

	// An object that defines a message that contains SSML-formatted text.
	SSMLMessage *SSMLMessageType

	noSmithyDocumentSerde
}

type noSmithyDocumentSerde = smithydocument.NoSerde
