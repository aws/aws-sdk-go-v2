package enhancedclient

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Map map[string]types.AttributeValue

func (m Map) String() string {
	buff := strings.Builder{}
	buff.WriteString("Map{")
	for key, value := range m {
		buff.WriteString(fmt.Sprintf("%s: %#+v, ", key, value))
	}
	buff.WriteString("}")

	return buff.String()
}

// With function takes a key string and a value of any kind and add it to the map as the corresponding AttributeValueMemberX type
// - []byte becomes types.AttributeValueMemberB
// - bool becomes types.AttributeValueMemberBOOL
// - [][]byte becomes types.AttributeValueMemberBS
// - []any becomes types.AttributeValueMemberL
// - map[any]any becomes types.AttributeValueMemberM
// - any type of int or float becomes types.AttributeValueMemberN
// - any type of int or float array ([5]type{...}) becomes types.AttributeValueMemberNS
// - nil becomes types.AttributeValueMemberNULL{Value: true}
// - string becomes types.AttributeValueMemberS
// - [3]string becomes types.AttributeValueMemberSS
// Note: [3] and [5] are not actual values we search for, they are just examples to illustrate go arrays vs go slices
func (m Map) With(key string, value any) Map {
	v := reflect.ValueOf(value)
	t := tag{}
	if v.Kind() == reflect.Array {
		k := v.Type().Elem().Kind()
		t.AsStrSet = k == reflect.String
		t.AsNumSet = k >= reflect.Int && k <= reflect.Float64 && k != reflect.Uintptr
		// t.AsBinSet is handled in encodeSlice()
	}
	av, _ := NewEncoder[any]().encode(v, t)

	m[key] = av

	return m
}
