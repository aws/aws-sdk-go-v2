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
			input: reflect.TypeFor[complex64](),
			ok:    false,
		},
		{
			input: reflect.TypeFor[complex128](),
			ok:    false,
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
			input: reflect.TypeFor[*complex64](),
			ok:    false,
		},
		{
			input: reflect.TypeFor[*complex128](),
			ok:    false,
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
func TestPointer(t *testing.T) {
	type foo struct{ X int }
	cases := []struct {
		name  string
		input any
		want  any
	}{
		{"int", 42, 42},
		{"string", "hello", "hello"},
		{"struct", foo{X: 7}, foo{X: 7}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			switch v := c.input.(type) {
			case int:
				p := pointer(v)
				if p == nil || *p != c.want.(int) {
					t.Errorf("pointer(int): got %v, want pointer to %v", p, c.want)
				}
			case string:
				p := pointer(v)
				if p == nil || *p != c.want.(string) {
					t.Errorf("pointer(string): got %v, want pointer to %v", p, c.want)
				}
			case foo:
				p := pointer(v)
				if p == nil || *p != c.want.(foo) {
					t.Errorf("pointer(struct): got %v, want pointer to %+v", p, c.want)
				}
			default:
				t.Fatalf("unsupported type: %T", v)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	type foo struct{ X int }
	i := 99
	s := "world"
	f := foo{X: 123}

	cases := []struct {
		name  string
		input any
		want  any
	}{
		{"*int", &i, i},
		{"nil *int", (*int)(nil), 0},
		{"*string", &s, s},
		{"nil *string", (*string)(nil), ""},
		{"*struct", &f, f},
		{"nil *struct", (*foo)(nil), foo{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			switch v := c.input.(type) {
			case *int:
				got := unwrap(v)
				if got != c.want.(int) {
					t.Errorf("unwrap(*int): got %v, want %v", got, c.want)
				}
			case *string:
				got := unwrap(v)
				if got != c.want.(string) {
					t.Errorf("unwrap(*string): got %v, want %v", got, c.want)
				}
			case *foo:
				got := unwrap(v)
				if got != c.want.(foo) {
					t.Errorf("unwrap(*struct): got %+v, want %+v", got, c.want)
				}
			default:
				t.Fatalf("unsupported type: %T", v)
			}
		})
	}
}
