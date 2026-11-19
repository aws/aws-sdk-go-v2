package converters

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ AttributeConverter[*bool] = (*BoolPtrConverter)(nil)

// BoolPtrConverter converts between Go *bool values and DynamoDB BOOL AttributeValues.
//
// BoolPtrConverter behaves similarly to BoolConverter but operates on pointer
// values instead of plain bools. This allows encoding and decoding of nullable
// boolean fields.
//
// Supported conversions:
//
//   - FromAttributeValue accepts *types.AttributeValueMemberBOOL and returns
//     a Go *bool pointing to the contained value.
//     If the provided AttributeValue is nil, it returns ErrNilValue.
//     If the provided AttributeValue is any other type, it returns unsupportedType.
//
//   - ToAttributeValue converts a Go *bool into *types.AttributeValueMemberBOOL.
//     If the input pointer is nil, it returns ErrNilValue.
//
// This converter only handles the DynamoDB BOOL type. All other AttributeValue
// variants will result in an unsupportedType error.
type BoolPtrConverter struct{}

// FromAttributeValue converts a DynamoDB BOOL AttributeValue to a Go *bool.
// Returns ErrNilValue for nil pointers, or unsupportedType for unsupported AttributeValue types.
func (n BoolPtrConverter) FromAttributeValue(v types.AttributeValue, _ []string) (*bool, error) {
	switch av := v.(type) {
	case *types.AttributeValueMemberBOOL:
		if av == nil {
			return nil, ErrNilValue
		}

		return &av.Value, nil
	default:
		return nil, unsupportedType(v, (*types.AttributeValueMemberBOOL)(nil))
	}
}

// ToAttributeValue converts a Go *bool to a DynamoDB BOOL AttributeValue.
// Returns ErrNilValue if the input pointer is nil.
func (n BoolPtrConverter) ToAttributeValue(v *bool, _ []string) (types.AttributeValue, error) {
	if v == nil {
		return nil, ErrNilValue
	}

	return &types.AttributeValueMemberBOOL{
		Value: *v,
	}, nil
}
