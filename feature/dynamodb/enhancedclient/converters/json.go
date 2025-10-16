package converters

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ AttributeConverter[any] = (*JsonConverter)(nil)

// JsonConverter converts between arbitrary Go values and DynamoDB AttributeValues
// that contain JSON data.
//
// Behaviour (matches the implementation in json.go):
//
// FromAttributeValue:
//   - Accepts *types.AttributeValueMemberB: treats av.Value as raw JSON bytes
//     and unmarshals into an `any` value. If av is nil, returns ErrNilValue.
//   - Accepts *types.AttributeValueMemberS: treats av.Value as a JSON string
//     and unmarshals into an `any` value. If av is nil, returns ErrNilValue.
//   - Accepts *types.AttributeValueMemberNULL: returns (nil, nil).
//   - Any other AttributeValue concrete type returns unsupportedType(...).
//
// ToAttributeValue:
//   - If the provided Go value is nil, returns *types.AttributeValueMemberNULL{Value: true}.
//   - Marshals the Go value to JSON. If `opts` contains `as=bytes`, returns
//     *types.AttributeValueMemberB with the JSON bytes; otherwise returns
//     *types.AttributeValueMemberS with the JSON string.
//   - Returns any JSON marshal error encountered.
//
// Notes:
//   - This converter only recognizes DynamoDB B and S attribute members containing
//     JSON payloads (and NULL). It does not convert arbitrary AttributeValue types.
//   - The converter returns ErrNilValue when a concrete B or S member pointer is nil,
//     and unsupportedType for mismatched AttributeValue types.
type JsonConverter struct {
}

func (j JsonConverter) FromAttributeValue(v types.AttributeValue, _ []string) (any, error) {
	switch av := v.(type) {
	case *types.AttributeValueMemberB:
		if av == nil {
			return nil, ErrNilValue
		}

		var o any
		err := json.Unmarshal(av.Value, &o)
		return o, err
	case *types.AttributeValueMemberS:
		if av == nil {
			return nil, ErrNilValue
		}

		var o any
		err := json.Unmarshal([]byte(av.Value), &o)
		return o, err
	case *types.AttributeValueMemberNULL:
		return nil, nil
	default:
		return "", unsupportedType(v, types.AttributeValueMemberS{}, types.AttributeValueMemberB{}, &types.AttributeValueMemberNULL{})
	}
}

func (j JsonConverter) ToAttributeValue(v any, opts []string) (types.AttributeValue, error) {
	as := getOpt(opts, "as")

	if v == nil {
		return &types.AttributeValueMemberNULL{Value: true}, nil
	}

	val, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	if as == "bytes" {
		return &types.AttributeValueMemberB{
			Value: val,
		}, nil
	}

	return &types.AttributeValueMemberS{
		Value: string(val),
	}, nil
}
