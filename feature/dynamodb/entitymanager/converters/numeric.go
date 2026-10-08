package converters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

var _ AttributeConverter[uint] = (*NumericConverter[uint])(nil)
var _ AttributeConverter[uint8] = (*NumericConverter[uint8])(nil)
var _ AttributeConverter[uint16] = (*NumericConverter[uint16])(nil)
var _ AttributeConverter[uint32] = (*NumericConverter[uint32])(nil)
var _ AttributeConverter[uint64] = (*NumericConverter[uint64])(nil)
var _ AttributeConverter[int] = (*NumericConverter[int])(nil)
var _ AttributeConverter[int8] = (*NumericConverter[int8])(nil)
var _ AttributeConverter[int16] = (*NumericConverter[int16])(nil)
var _ AttributeConverter[int32] = (*NumericConverter[int32])(nil)
var _ AttributeConverter[int64] = (*NumericConverter[int64])(nil)
var _ AttributeConverter[float32] = (*NumericConverter[float32])(nil)
var _ AttributeConverter[float64] = (*NumericConverter[float64])(nil)

// NumericConverter converts between Go numeric types and DynamoDB number (N)
// AttributeValues.
//
// It is a generic converter parameterized by T, which must satisfy the internal
// `number` constraint (e.g. uint, int64, float64, etc.).
//
// Behaviour:
//
// FromAttributeValue:
//   - Accepts *types.AttributeValueMemberN and parses av.Value into the target
//     numeric type T using strconv.
//   - Returns ErrNilValue if the AttributeValue pointer is nil.
//   - Returns unsupportedType if the AttributeValue is not a number type.
//
// ToAttributeValue:
//   - Converts a Go numeric value of type T into a *types.AttributeValueMemberN
//     with its string representation (via strconv.FormatFloat / strconv.FormatInt).
//
// This converter only supports DynamoDB number attributes and Go numeric types;
// any other AttributeValue type or Go kind will result in an error.
//
// Example:
//
//	var c converters.NumericConverter[int]
//	av, _ := c.ToAttributeValue(42, nil)
//	// av == &types.AttributeValueMemberN{Value: "42"}
//
//	v, _ := c.FromAttributeValue(av, nil)
//	// v == int(42)
type NumericConverter[T number] struct {
}

func (n NumericConverter[T]) FromAttributeValue(v types.AttributeValue, i []string) (T, error) {
	out := *new(T)

	switch av := v.(type) {
	case *types.AttributeValueMemberN:
		if strings.Contains(av.Value, ".") {
			f, err := strconv.ParseFloat(av.Value, 64)
			if err != nil {
				return T(0), err
			}

			out = T(f)
		} else {
			i, err := strconv.ParseInt(av.Value, 10, 64)
			if err != nil {
				return T(0), err
			}

			out = T(i)
		}
	default:
		return T(0), unsupportedType(v, (*types.AttributeValueMemberN)(nil))
	}

	return out, nil
}

func (n NumericConverter[T]) ToAttributeValue(t T, i []string) (types.AttributeValue, error) {
	return &types.AttributeValueMemberN{
		Value: fmt.Sprintf("%v", t),
	}, nil
}
