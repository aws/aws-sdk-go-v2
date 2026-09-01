package converters

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Compile-time assertion that StringPtrConverter satisfies AttributeConverter[*string].
var _ AttributeConverter[*string] = (*StringPtrConverter)(nil)

// StringPtrConverter implements AttributeConverter for the Go *string type.
//
// DynamoDB stores string values in AttributeValueMemberS. This converter:
//   - Validates that the provided AttributeValue is of type *types.AttributeValueMemberS
//   - Returns ErrNilValue if the concrete *types.AttributeValueMemberS is nil
//   - Distinguishes between a nil *string (error when marshaling) and a pointer to an empty string (valid value)
//
// Tag options: Currently ignored. The parameter is accepted to keep a uniform
// signature with other converters that may use tag-derived options.
//
// Usage:
//
//	conv := StringPtrConverter{}
//	val := "hello"
//	av, _ := conv.ToAttributeValue(&val, nil)
//	s, _ := conv.FromAttributeValue(av, nil)
//
// The converter is stateless and safe for concurrent use.
type StringPtrConverter struct{}

// FromAttributeValue converts a DynamoDB string (S) AttributeValue to a Go *string.
// Returns ErrNilValue for nil pointers, or unsupportedType for unsupported AttributeValue types.
func (n StringPtrConverter) FromAttributeValue(v types.AttributeValue, _ []string) (*string, error) {
	switch av := v.(type) {
	case *types.AttributeValueMemberS:
		if av == nil {
			return nil, ErrNilValue
		}

		return &av.Value, nil
	default:
		return nil, unsupportedType(v, (*types.AttributeValueMemberS)(nil))
	}
}

// ToAttributeValue converts a Go *string to a DynamoDB string (S) AttributeValue.
// Returns ErrNilValue if the input pointer is nil.
func (n StringPtrConverter) ToAttributeValue(v *string, _ []string) (types.AttributeValue, error) {
	if v == nil {
		return nil, ErrNilValue
	}

	return &types.AttributeValueMemberS{
		Value: *v,
	}, nil
}
