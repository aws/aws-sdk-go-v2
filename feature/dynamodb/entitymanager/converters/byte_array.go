package converters

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// ByteArrayConverter converts between Go byte slices ([]byte) and DynamoDB
// binary (B) AttributeValues.
//
// Behaviour:
//
// FromAttributeValue:
//   - Accepts *types.AttributeValueMemberB and returns its Value ([]byte).
//   - Returns ErrNilValue if the AttributeValue pointer is nil.
//   - Returns unsupportedType if the AttributeValue is not of type B.
//
// ToAttributeValue:
//   - Converts a Go []byte into a *types.AttributeValueMemberB containing that byte slice.
//   - Returns ErrNilValue if the input slice is nil.
//
// This converter only supports DynamoDB binary attributes and Go []byte values.
// Any other AttributeValue type will cause an unsupportedType error.
//
// Example:
//
//	var c converters.ByteArrayConverter
//	av, _ := c.ToAttributeValue([]byte("hello"), nil)
//	// av == &types.AttributeValueMemberB{Value: []byte("hello")}
//
//	v, _ := c.FromAttributeValue(av, nil)
//	// v == []byte("hello")
type ByteArrayConverter struct {
}

// FromAttributeValue converts a DynamoDB binary (B) AttributeValue to a Go []byte.
// Returns ErrNilValue for nil pointers, or unsupportedType for unsupported AttributeValue types.
func (n ByteArrayConverter) FromAttributeValue(v types.AttributeValue, _ []string) ([]byte, error) {
	switch av := v.(type) {
	case *types.AttributeValueMemberB:
		if av == nil {
			return nil, ErrNilValue
		}
		return av.Value, nil
	default:
		return nil, unsupportedType(v, (*types.AttributeValueMemberB)(nil))
	}
}

// ToAttributeValue converts a Go []byte to a DynamoDB binary (B) AttributeValue.
// Returns ErrNilValue if the input slice is nil.
func (n ByteArrayConverter) ToAttributeValue(v []byte, _ []string) (types.AttributeValue, error) {
	if v == nil {
		return nil, ErrNilValue
	}

	return &types.AttributeValueMemberB{
		Value: v,
	}, nil
}
