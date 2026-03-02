package converters

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ AttributeConverter[string] = (*StringConverter)(nil)

// StringConverter implements AttributeConverter for the Go string type.
//
// DynamoDB stores string values in AttributeValueMemberS. This converter:
//   - Validates that the provided AttributeValue is of type *types.AttributeValueMemberS
//   - Returns an error if a nil AttributeValueMemberS is encountered
//   - Performs no special handling for empty strings ("" is a valid value)
//
// Tag options: Currently ignored. The parameter is accepted to keep a uniform
// signature with other converters that may use tag-derived options.
//
// Usage:
//
//	conv := StringConverter{}
//	av, _ := conv.ToAttributeValue("hello", nil)
//	s, _ := conv.FromAttributeValue(av, nil)
//
// The converter is stateless and safe for concurrent use.
type StringConverter struct {
}

// FromAttributeValue converts a DynamoDB string (S) AttributeValue to a Go string.
// Returns ErrNilValue for nil pointers, or unsupportedType for unsupported AttributeValue types.
func (n StringConverter) FromAttributeValue(v types.AttributeValue, _ []string) (string, error) {
	switch av := v.(type) {
	case *types.AttributeValueMemberS:
		if av == nil {
			return "", ErrNilValue
		}

		return av.Value, nil
	default:
		return "", unsupportedType(v, (*types.AttributeValueMemberS)(nil))
	}
}

// ToAttributeValue converts a Go string to a DynamoDB string (S) AttributeValue.
func (n StringConverter) ToAttributeValue(v string, _ []string) (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: v,
	}, nil
}
