package protocol_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/private/protocol"
)

func TestCanSetIdempotencyToken(t *testing.T) {
	cases := []struct {
		CanSet bool
		Case   interface{}
	}{
		{
			true,
			struct {
				Field *string `idempotencyToken:"true"`
			}{},
		},
		{
			true,
			struct {
				Field string `idempotencyToken:"true"`
			}{},
		},
		{
			false,
			struct {
				Field *string `idempotencyToken:"true"`
			}{Field: new(string)},
		},
		{
			false,
			struct {
				Field string `idempotencyToken:"true"`
			}{Field: "value"},
		},
		{
			false,
			struct {
				Field *int `idempotencyToken:"true"`
			}{},
		},
		{
			false,
			struct {
				Field *string
			}{},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			v := reflect.Indirect(reflect.ValueOf(c.Case))
			ty := v.Type()
			canSet := protocol.CanSetIdempotencyToken(v.Field(0), ty.Field(0))
			if e, a := c.CanSet, canSet; e != a {
				t.Errorf("expect %v can set, got %v", e, a)
			}
		})
	}
}

func TestSetIdempotencyToken(t *testing.T) {
	cases := []struct {
		Case interface{}
	}{
		{
			&struct {
				Field *string `idempotencyToken:"true"`
			}{},
		},
		{
			&struct {
				Field string `idempotencyToken:"true"`
			}{},
		},
		{
			&struct {
				Field *string `idempotencyToken:"true"`
			}{Field: new(string)},
		},
		{
			&struct {
				Field string `idempotencyToken:"true"`
			}{Field: ""},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			v := reflect.Indirect(reflect.ValueOf(c.Case))

			protocol.SetIdempotencyToken(v.Field(0))
			switch tv := v.Field(0).Interface().(type) {
			case *string:
				if tv == nil || len(*tv) == 0 {
					t.Errorf("expect to be set")
				}
			case string:
				if len(tv) == 0 {
					t.Errorf("expect to be set")
				}
			default:
				t.Errorf("value is not a string")
			}
		})
	}
}

func TestUUIDVersion4(t *testing.T) {
	uuid := protocol.UUIDVersion4(make([]byte, 16))
	if e, a := `00000000-0000-4000-8000-000000000000`, uuid; e != a {
		t.Errorf("expect %v uuid, got %v", e, a)
	}

	b := make([]byte, 16)
	for i := 0; i < len(b); i++ {
		b[i] = 1
	}
	uuid = protocol.UUIDVersion4(b)
	if e, a := `01010101-0101-4101-8101-010101010101`, uuid; e != a {
		t.Errorf("expect %v uuid, got %v", e, a)
	}
}
