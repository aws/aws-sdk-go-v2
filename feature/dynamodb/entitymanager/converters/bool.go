package converters

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ AttributeConverter[bool] = (*BoolConverter)(nil)

// BoolConverter converts between Go bool values and DynamoDB BOOL AttributeValues.
//
// BoolConverter only supports the DynamoDB BOOL attribute type:
//
//   - FromAttributeValue accepts *types.AttributeValueMemberBOOL and returns the contained bool.
//     If the provided AttributeValue is nil (the pointer to the member type is nil),
//     it returns ErrNilValue.
//     If the provided AttributeValue is any other AttributeValue type, it returns unsupportedType.
//
//   - ToAttributeValue converts a Go bool into *types.AttributeValueMemberBOOL.
//
// This converter does NOT interpret strings or numbers as booleans — any non-BOOL
// AttributeValue will cause an unsupportedType error.
type BoolConverter struct{}

// FromAttributeValue converts a DynamoDB BOOL AttributeValue to a Go bool.
// Returns ErrNilValue for nil pointers, or unsupportedType for unsupported AttributeValue types.
func (n BoolConverter) FromAttributeValue(v types.AttributeValue, _ []string) (bool, error) {
	switch av := v.(type) {
	case *types.AttributeValueMemberBOOL:
		if av == nil {
			return false, ErrNilValue
		}

		return av.Value, nil
	default:
		return false, unsupportedType(v, types.AttributeValueMemberBOOL{})
	}
}

// ToAttributeValue converts a Go bool to a DynamoDB BOOL AttributeValue.
func (n BoolConverter) ToAttributeValue(v bool, _ []string) (types.AttributeValue, error) {
	return &types.AttributeValueMemberBOOL{
		Value: v,
	}, nil
}
