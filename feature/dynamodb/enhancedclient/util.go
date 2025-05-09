package enhancedclient

import (
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func pointer[T any](v T) *T {
	return &v
}

func unwrap[T any](v *T) T {
	if v != nil {
		return *v
	}

	return *new(T)
}

func typeToScalarAttributeType(t reflect.Type) (types.ScalarAttributeType, bool) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return types.ScalarAttributeTypeS, true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return types.ScalarAttributeTypeN, true
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return types.ScalarAttributeTypeB, true
		}
		fallthrough
	default:
		return "", false // unknown or unsupported kind
	}
}
