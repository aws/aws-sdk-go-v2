package attributevalue_test

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func ExampleMarshal() {
	type Record struct {
		Bytes   []byte
		MyField string
		Letters []string
		Numbers []int
	}

	r := Record{
		Bytes:   []byte{48, 49},
		MyField: "MyFieldValue",
		Letters: []string{"a", "b", "c", "d"},
		Numbers: []int{1, 2, 3},
	}
	av, err := attributevalue.Marshal(r)
	m := av.(*types.AttributeValueMemberM)
	fmt.Println("err", err)
	fmt.Println("Bytes", awsutil.Prettify(m.Value["Bytes"]))
	fmt.Println("MyField", awsutil.Prettify(m.Value["MyField"]))
	fmt.Println("Letters", awsutil.Prettify(m.Value["Letters"]))
	fmt.Println("Numbers", awsutil.Prettify(m.Value["Numbers"]))

	// Output:
	// err <nil>
	// Bytes {
	//   Value: <binary> len 2
	// }
	// MyField {
	//   Value: "MyFieldValue"
	// }
	// Letters {
	//   Value: [
	//     &{a},
	//     &{b},
	//     &{c},
	//     &{d}
	//   ]
	// }
	// Numbers {
	//   Value: [&{1},&{2},&{3}]
	// }
}

func ExampleUnmarshal() {
	type Record struct {
		Bytes   []byte
		MyField string
		Letters []string
		A2Num   map[string]int
	}

	expect := Record{
		Bytes:   []byte{48, 49},
		MyField: "MyFieldValue",
		Letters: []string{"a", "b", "c", "d"},
		A2Num:   map[string]int{"a": 1, "b": 2, "c": 3},
	}

	av := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Bytes":   &types.AttributeValueMemberB{Value: []byte{48, 49}},
			"MyField": &types.AttributeValueMemberS{Value: "MyFieldValue"},
			"Letters": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "a"},
				&types.AttributeValueMemberS{Value: "b"},
				&types.AttributeValueMemberS{Value: "c"},
				&types.AttributeValueMemberS{Value: "d"},
			}},
			"A2Num": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"a": &types.AttributeValueMemberN{Value: "1"},
				"b": &types.AttributeValueMemberN{Value: "2"},
				"c": &types.AttributeValueMemberN{Value: "3"},
			}},
		},
	}

	actual := Record{}
	err := attributevalue.Unmarshal(av, &actual)
	fmt.Println(err, reflect.DeepEqual(expect, actual))

	// Output:
	// <nil> true
}
