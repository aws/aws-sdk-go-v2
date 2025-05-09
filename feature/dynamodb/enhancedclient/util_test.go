package enhancedclient

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestTypeToScalarAttributeType(t *testing.T) {
	cases := []struct {
		input    reflect.Type
		expected types.ScalarAttributeType
		ok       bool
	}{
		{
			input:    reflect.TypeFor[uint](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[uint8](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[uint16](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[uint32](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[uint64](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[int](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[int8](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[int16](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[int32](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[int64](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[float32](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[float64](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[complex64](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[complex128](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[string](),
			expected: types.ScalarAttributeTypeS,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[[]byte](),
			expected: types.ScalarAttributeTypeB,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[[1]byte](),
			expected: types.ScalarAttributeTypeB,
			ok:       true,
		},
		///

		{
			input:    reflect.TypeFor[*uint](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*uint8](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*uint16](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*uint32](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*uint64](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*int](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*int8](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*int16](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*int32](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*int64](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*float32](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*float64](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*complex64](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*complex128](),
			expected: types.ScalarAttributeTypeN,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*string](),
			expected: types.ScalarAttributeTypeS,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*[]byte](),
			expected: types.ScalarAttributeTypeB,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[*[1]byte](),
			expected: types.ScalarAttributeTypeB,
			ok:       true,
		},
		{
			input:    reflect.TypeFor[order](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[*order](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[[]order](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[*[]order](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[map[string]string](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[*map[string]string](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[[]map[string]string](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[*[]map[string]string](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[any](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[[]any](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[map[string]string](),
			expected: "",
			ok:       false,
		},
		{
			input:    reflect.TypeFor[chan any](),
			expected: "",
			ok:       false,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, ok := typeToScalarAttributeType(c.input)

			if diff := cmpDiff(c.expected, actual); len(diff) != 0 {
				t.Errorf("different values: %s", diff)
			}

			if diff := cmpDiff(c.ok, ok); len(diff) != 0 {
				t.Errorf("different values: %s", diff)
			}
		})
	}
}
