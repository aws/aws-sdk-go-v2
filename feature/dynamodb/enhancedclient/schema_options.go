package enhancedclient

import "github.com/aws/aws-sdk-go-v2/feature/dynamodb/enhancedclient/converters"

// SchemaOptions defines configuration options for Schema behavior.
type SchemaOptions struct {
	// ErrorOnMissingField controls whether decoding should return an error
	// when a field is missing in the destination struct.
	// If true, decoding will fail when the schema field cannot be matched.
	// If false or nil, missing fields will be ignored.
	ErrorOnMissingField *bool

	// IgnoreNilValueErrors controls whether decoding should ignore errors
	// caused by nil values during schema conversion.
	// If true, fields with nil values that cause conversion errors will be skipped.
	// If false or nil, such cases will trigger an error.
	IgnoreNilValueErrors *bool

	// ConverterRegistry provides a registry of type converters used during
	// encoding and decoding operations. It will be set on both the Decoder
	// and Encoder to control how values are transformed between Go types
	// and schema representations.
	ConverterRegistry *converters.Registry
}
