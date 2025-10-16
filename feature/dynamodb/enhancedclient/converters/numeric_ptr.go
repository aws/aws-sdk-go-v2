package converters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ AttributeConverter[*uint] = (*NumericPtrConverter[uint])(nil)
var _ AttributeConverter[*uint8] = (*NumericPtrConverter[uint8])(nil)
var _ AttributeConverter[*uint16] = (*NumericPtrConverter[uint16])(nil)
var _ AttributeConverter[*uint32] = (*NumericPtrConverter[uint32])(nil)
var _ AttributeConverter[*uint64] = (*NumericPtrConverter[uint64])(nil)
var _ AttributeConverter[*int] = (*NumericPtrConverter[int])(nil)
var _ AttributeConverter[*int8] = (*NumericPtrConverter[int8])(nil)
var _ AttributeConverter[*int16] = (*NumericPtrConverter[int16])(nil)
var _ AttributeConverter[*int32] = (*NumericPtrConverter[int32])(nil)
var _ AttributeConverter[*int64] = (*NumericPtrConverter[int64])(nil)
var _ AttributeConverter[*float32] = (*NumericPtrConverter[float32])(nil)
var _ AttributeConverter[*float64] = (*NumericPtrConverter[float64])(nil)

// NumericPtrConverter converts between Go pointer-to-number values (*T)
// and DynamoDB number (N) AttributeValues.
//
// It is a generic converter parameterized by T, which must satisfy the internal
// `number` constraint (e.g. uint, int64, float64, etc.).
//
// Behaviour:
//
// FromAttributeValue:
//   - Accepts *types.AttributeValueMemberN and parses av.Value into a new value
//     of type T, returning a pointer to it (*T).
//   - Returns ErrNilValue if the AttributeValue pointer is nil.
//   - Returns unsupportedType if the AttributeValue is not a number type.
//
// ToAttributeValue:
//   - Converts a Go *T into a *types.AttributeValueMemberN containing the string
//     representation of the pointed-to numeric value.
//   - Returns ErrNilValue if the input pointer is nil.
//
// This converter only supports DynamoDB number attributes and Go numeric pointer
// types. Any other AttributeValue type or nil handling violation will result in
// an error.
//
// Example:
//
//	var c converters.NumericPtrConverter[int]
//	v := 42
//	av, _ := c.ToAttributeValue(&v, nil)
//	// av == &types.AttributeValueMemberN{Value: "42"}
//
//	out, _ := c.FromAttributeValue(av, nil)
//	// *out == 42
type NumericPtrConverter[T number] struct {
}

func (n NumericPtrConverter[T]) FromAttributeValue(v types.AttributeValue, _ []string) (*T, error) {
	out := new(T)

	switch av := v.(type) {
	case *types.AttributeValueMemberN:
		if strings.Contains(av.Value, ".") {
			f, err := strconv.ParseFloat(av.Value, 64)
			if err != nil {
				return nil, err
			}

			*out = T(f)
		} else {
			i, err := strconv.ParseInt(av.Value, 10, 64)
			if err != nil {
				return nil, err
			}

			*out = T(i)
		}
	default:
		return nil, unsupportedType(v, (*types.AttributeValueMemberN)(nil))
	}

	return out, nil
}

func (n NumericPtrConverter[T]) ToAttributeValue(v *T, _ []string) (types.AttributeValue, error) {
	if v == nil {
		return nil, ErrNilValue
	}

	return &types.AttributeValueMemberN{
		Value: fmt.Sprintf("%v", *v),
	}, nil
}
