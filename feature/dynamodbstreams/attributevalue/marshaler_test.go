package attributevalue

import (
	"math"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
	"github.com/google/go-cmp/cmp"
)

type simpleMarshalStruct struct {
	Byte      []byte
	String    string
	PtrString *string
	Int       int
	Uint      uint
	Float32   float32
	Float64   float64
	Bool      bool
	Null      *interface{}
}

type complexMarshalStruct struct {
	Simple []simpleMarshalStruct
}

type myByteStruct struct {
	Byte []byte
}

type myByteSetStruct struct {
	ByteSet [][]byte
}

type marshallerTestInput struct {
	input    interface{}
	expected interface{}
	err      error
}

var trueValue = true
var falseValue = false

var marshalerScalarInputs = map[string]marshallerTestInput{
	"nil": {
		input:    nil,
		expected: &types.AttributeValueMemberNULL{Value: true},
	},
	"string": {
		input:    "some string",
		expected: &types.AttributeValueMemberS{Value: "some string"},
	},
	"bool": {
		input:    true,
		expected: &types.AttributeValueMemberBOOL{Value: true},
	},
	"bool false": {
		input:    false,
		expected: &types.AttributeValueMemberBOOL{Value: false},
	},
	"float": {
		input:    3.14,
		expected: &types.AttributeValueMemberN{Value: "3.14"},
	},
	"max float32": {
		input:    math.MaxFloat32,
		expected: &types.AttributeValueMemberN{Value: "340282346638528860000000000000000000000"},
	},
	"max float64": {
		input:    math.MaxFloat64,
		expected: &types.AttributeValueMemberN{Value: "179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
	},
	"integer": {
		input:    12,
		expected: &types.AttributeValueMemberN{Value: "12"},
	},
	"number integer": {
		input:    Number("12"),
		expected: &types.AttributeValueMemberN{Value: "12"},
	},
	"zero values": {
		input: simpleMarshalStruct{},
		expected: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Byte":      &types.AttributeValueMemberNULL{Value: true},
				"Bool":      &types.AttributeValueMemberBOOL{Value: false},
				"Float32":   &types.AttributeValueMemberN{Value: "0"},
				"Float64":   &types.AttributeValueMemberN{Value: "0"},
				"Int":       &types.AttributeValueMemberN{Value: "0"},
				"Null":      &types.AttributeValueMemberNULL{Value: true},
				"String":    &types.AttributeValueMemberS{Value: ""},
				"PtrString": &types.AttributeValueMemberNULL{Value: true},
				"Uint":      &types.AttributeValueMemberN{Value: "0"},
			},
		},
	},
}

var marshallerMapTestInputs = map[string]marshallerTestInput{
	// Scalar tests
	"nil": {
		input:    nil,
		expected: map[string]types.AttributeValue{},
	},
	"string": {
		input:    map[string]interface{}{"string": "some string"},
		expected: map[string]types.AttributeValue{"string": &types.AttributeValueMemberS{Value: "some string"}},
	},
	"bool": {
		input:    map[string]interface{}{"bool": true},
		expected: map[string]types.AttributeValue{"bool": &types.AttributeValueMemberBOOL{Value: true}},
	},
	"bool false": {
		input:    map[string]interface{}{"bool": false},
		expected: map[string]types.AttributeValue{"bool": &types.AttributeValueMemberBOOL{Value: false}},
	},
	"null": {
		input:    map[string]interface{}{"null": nil},
		expected: map[string]types.AttributeValue{"null": &types.AttributeValueMemberNULL{Value: true}},
	},
	"float": {
		input:    map[string]interface{}{"float": 3.14},
		expected: map[string]types.AttributeValue{"float": &types.AttributeValueMemberN{Value: "3.14"}},
	},
	"float32": {
		input:    map[string]interface{}{"float": math.MaxFloat32},
		expected: map[string]types.AttributeValue{"float": &types.AttributeValueMemberN{Value: "340282346638528860000000000000000000000"}},
	},
	"float64": {
		input:    map[string]interface{}{"float": math.MaxFloat64},
		expected: map[string]types.AttributeValue{"float": &types.AttributeValueMemberN{Value: "179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"}},
	},
	"decimal number": {
		input:    map[string]interface{}{"num": 12.},
		expected: map[string]types.AttributeValue{"num": &types.AttributeValueMemberN{Value: "12"}},
	},
	"byte": {
		input:    map[string]interface{}{"byte": []byte{48, 49}},
		expected: map[string]types.AttributeValue{"byte": &types.AttributeValueMemberB{Value: []byte{48, 49}}},
	},
	"nested blob": {
		input:    struct{ Byte []byte }{Byte: []byte{48, 49}},
		expected: map[string]types.AttributeValue{"Byte": &types.AttributeValueMemberB{Value: []byte{48, 49}}},
	},
	"map nested blob": {
		input:    map[string]interface{}{"byte_set": [][]byte{{48, 49}, {50, 51}}},
		expected: map[string]types.AttributeValue{"byte_set": &types.AttributeValueMemberBS{Value: [][]byte{{48, 49}, {50, 51}}}},
	},
	"bytes set": {
		input:    struct{ ByteSet [][]byte }{ByteSet: [][]byte{{48, 49}, {50, 51}}},
		expected: map[string]types.AttributeValue{"ByteSet": &types.AttributeValueMemberBS{Value: [][]byte{{48, 49}, {50, 51}}}},
	},
	"list": {
		input: map[string]interface{}{"list": []interface{}{"a string", 12., 3.14, true, nil, false}},
		expected: map[string]types.AttributeValue{
			"list": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "a string"},
					&types.AttributeValueMemberN{Value: "12"},
					&types.AttributeValueMemberN{Value: "3.14"},
					&types.AttributeValueMemberBOOL{Value: true},
					&types.AttributeValueMemberNULL{Value: true},
					&types.AttributeValueMemberBOOL{Value: false},
				},
			},
		},
	},
	"map": {
		input: map[string]interface{}{"map": map[string]interface{}{"nestednum": 12.}},
		expected: map[string]types.AttributeValue{
			"map": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"nestednum": &types.AttributeValueMemberN{Value: "12"},
				},
			},
		},
	},
	"struct": {
		input: simpleMarshalStruct{},
		expected: map[string]types.AttributeValue{
			"Byte":      &types.AttributeValueMemberNULL{Value: true},
			"Bool":      &types.AttributeValueMemberBOOL{Value: false},
			"Float32":   &types.AttributeValueMemberN{Value: "0"},
			"Float64":   &types.AttributeValueMemberN{Value: "0"},
			"Int":       &types.AttributeValueMemberN{Value: "0"},
			"Null":      &types.AttributeValueMemberNULL{Value: true},
			"String":    &types.AttributeValueMemberS{Value: ""},
			"PtrString": &types.AttributeValueMemberNULL{Value: true},
			"Uint":      &types.AttributeValueMemberN{Value: "0"},
		},
	},
	"nested struct": {
		input: complexMarshalStruct{},
		expected: map[string]types.AttributeValue{
			"Simple": &types.AttributeValueMemberNULL{Value: true},
		},
	},
	"nested nil slice": {
		input: struct {
			Simple []string `dynamodbav:"simple"`
		}{},
		expected: map[string]types.AttributeValue{
			"simple": &types.AttributeValueMemberNULL{Value: true},
		},
	},
	"nested nil slice omit empty": {
		input: struct {
			Simple []string `dynamodbav:"simple,omitempty"`
		}{},
		expected: map[string]types.AttributeValue{},
	},
	"nested ignored field": {
		input: struct {
			Simple []string `dynamodbav:"-"`
		}{},
		expected: map[string]types.AttributeValue{},
	},
	"complex struct members with zero": {
		input: complexMarshalStruct{Simple: []simpleMarshalStruct{{Int: -2}, {Uint: 5}}},
		expected: map[string]types.AttributeValue{
			"Simple": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"Byte":      &types.AttributeValueMemberNULL{Value: true},
							"Bool":      &types.AttributeValueMemberBOOL{Value: false},
							"Float32":   &types.AttributeValueMemberN{Value: "0"},
							"Float64":   &types.AttributeValueMemberN{Value: "0"},
							"Int":       &types.AttributeValueMemberN{Value: "-2"},
							"Null":      &types.AttributeValueMemberNULL{Value: true},
							"String":    &types.AttributeValueMemberS{Value: ""},
							"PtrString": &types.AttributeValueMemberNULL{Value: true},
							"Uint":      &types.AttributeValueMemberN{Value: "0"},
						},
					},
					&types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"Byte":      &types.AttributeValueMemberNULL{Value: true},
							"Bool":      &types.AttributeValueMemberBOOL{Value: false},
							"Float32":   &types.AttributeValueMemberN{Value: "0"},
							"Float64":   &types.AttributeValueMemberN{Value: "0"},
							"Int":       &types.AttributeValueMemberN{Value: "0"},
							"Null":      &types.AttributeValueMemberNULL{Value: true},
							"String":    &types.AttributeValueMemberS{Value: ""},
							"PtrString": &types.AttributeValueMemberNULL{Value: true},
							"Uint":      &types.AttributeValueMemberN{Value: "5"},
						},
					},
				},
			},
		},
	},
}

var marshallerListTestInputs = map[string]marshallerTestInput{
	"nil": {
		input:    nil,
		expected: []types.AttributeValue{},
	},
	"empty interface": {
		input:    []interface{}{},
		expected: []types.AttributeValue{},
	},
	"empty struct": {
		input:    []simpleMarshalStruct{},
		expected: []types.AttributeValue{},
	},
	"various types": {
		input: []interface{}{"a string", 12., 3.14, true, nil, false},
		expected: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "a string"},
			&types.AttributeValueMemberN{Value: "12"},
			&types.AttributeValueMemberN{Value: "3.14"},
			&types.AttributeValueMemberBOOL{Value: true},
			&types.AttributeValueMemberNULL{Value: true},
			&types.AttributeValueMemberBOOL{Value: false},
		},
	},
	"nested zero values": {
		input: []simpleMarshalStruct{{}},
		expected: []types.AttributeValue{
			&types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"Byte":      &types.AttributeValueMemberNULL{Value: true},
					"Bool":      &types.AttributeValueMemberBOOL{Value: false},
					"Float32":   &types.AttributeValueMemberN{Value: "0"},
					"Float64":   &types.AttributeValueMemberN{Value: "0"},
					"Int":       &types.AttributeValueMemberN{Value: "0"},
					"Null":      &types.AttributeValueMemberNULL{Value: true},
					"String":    &types.AttributeValueMemberS{Value: ""},
					"PtrString": &types.AttributeValueMemberNULL{Value: true},
					"Uint":      &types.AttributeValueMemberN{Value: "0"},
				},
			},
		},
	},
}

func Test_New_Marshal(t *testing.T) {
	for name, test := range marshalerScalarInputs {
		t.Run(name, func(t *testing.T) {
			actual, err := Marshal(test.input)
			if test.err != nil {
				if err == nil {
					t.Errorf("Marshal with input %#v returned %#v, expected error `%s`",
						test.input, actual, test.err)
				} else if err.Error() != test.err.Error() {
					t.Errorf("Marshal with input %#v returned error `%s`, expected error `%s`",
						test.input, err, test.err)
				}
			} else {
				if err != nil {
					t.Errorf("Marshal with input %#v returned error `%s`", test.input, err)
				}
				compareObjects(t, test.expected, actual)
			}
		})
	}
}

func testMarshal(t *testing.T, test marshallerTestInput) {
}

func Test_New_Unmarshal(t *testing.T) {
	// Using the same inputs from Marshal, test the reverse mapping.
	for name, test := range marshalerScalarInputs {
		t.Run(name, func(t *testing.T) {
			if test.input == nil {
				t.Skip()
			}
			actual := reflect.New(reflect.TypeOf(test.input)).Interface()
			if err := Unmarshal(test.expected.(types.AttributeValue), actual); err != nil {
				t.Errorf("Unmarshal with input %#v returned error `%s`", test.expected, err)
			}
			compareObjects(t, test.input, reflect.ValueOf(actual).Elem().Interface())
		})
	}
}

func Test_New_UnmarshalError(t *testing.T) {
	// Test that we get an error using Unmarshal to convert to a nil value.
	expected := &InvalidUnmarshalError{Type: reflect.TypeOf(nil)}
	if err := Unmarshal(nil, nil); err == nil {
		t.Errorf("Unmarshal with input %T returned no error, expected error `%v`", nil, expected)
	} else if err.Error() != expected.Error() {
		t.Errorf("Unmarshal with input %T returned error `%v`, expected error `%v`", nil, err, expected)
	}

	// Test that we get an error using Unmarshal to convert to a non-pointer value.
	var actual map[string]interface{}
	expected = &InvalidUnmarshalError{Type: reflect.TypeOf(actual)}
	if err := Unmarshal(nil, actual); err == nil {
		t.Errorf("Unmarshal with input %T returned no error, expected error `%v`", actual, expected)
	} else if err.Error() != expected.Error() {
		t.Errorf("Unmarshal with input %T returned error `%v`, expected error `%v`", actual, err, expected)
	}

	// Test that we get an error using Unmarshal to convert to nil struct.
	var actual2 *struct{ A int }
	expected = &InvalidUnmarshalError{Type: reflect.TypeOf(actual2)}
	if err := Unmarshal(nil, actual2); err == nil {
		t.Errorf("Unmarshal with input %T returned no error, expected error `%v`", actual2, expected)
	} else if err.Error() != expected.Error() {
		t.Errorf("Unmarshal with input %T returned error `%v`, expected error `%v`", actual2, err, expected)
	}
}

func Test_New_MarshalMap(t *testing.T) {
	for name, test := range marshallerMapTestInputs {
		t.Run(name, func(t *testing.T) {
			actual, err := MarshalMap(test.input)
			if test.err != nil {
				if err == nil {
					t.Errorf("MarshalMap with input %#v returned %#v, expected error `%s`",
						test.input, actual, test.err)
				} else if err.Error() != test.err.Error() {
					t.Errorf("MarshalMap with input %#v returned error `%s`, expected error `%s`",
						test.input, err, test.err)
				}
			} else {
				if err != nil {
					t.Errorf("MarshalMap with input %#v returned error `%s`", test.input, err)
				}
				compareObjects(t, test.expected, actual)
			}
		})
	}
}

func Test_New_UnmarshalMap(t *testing.T) {
	// Using the same inputs from MarshalMap, test the reverse mapping.
	for name, test := range marshallerMapTestInputs {
		t.Run(name, func(t *testing.T) {
			if test.input == nil {
				t.Skip()
			}
			actual := reflect.New(reflect.TypeOf(test.input)).Interface()
			if err := UnmarshalMap(test.expected.(map[string]types.AttributeValue), actual); err != nil {
				t.Errorf("Unmarshal with input %#v returned error `%s`", test.expected, err)
			}
			compareObjects(t, test.input, reflect.ValueOf(actual).Elem().Interface())
		})
	}
}

func Test_New_UnmarshalMapError(t *testing.T) {
	// Test that we get an error using UnmarshalMap to convert to a nil value.
	expected := &InvalidUnmarshalError{Type: reflect.TypeOf(nil)}
	if err := UnmarshalMap(nil, nil); err == nil {
		t.Errorf("UnmarshalMap with input %T returned no error, expected error `%v`", nil, expected)
	} else if err.Error() != expected.Error() {
		t.Errorf("UnmarshalMap with input %T returned error `%v`, expected error `%v`", nil, err, expected)
	}

	// Test that we get an error using UnmarshalMap to convert to a non-pointer value.
	var actual map[string]interface{}
	expected = &InvalidUnmarshalError{Type: reflect.TypeOf(actual)}
	if err := UnmarshalMap(nil, actual); err == nil {
		t.Errorf("UnmarshalMap with input %T returned no error, expected error `%v`", actual, expected)
	} else if err.Error() != expected.Error() {
		t.Errorf("UnmarshalMap with input %T returned error `%v`, expected error `%v`", actual, err, expected)
	}

	// Test that we get an error using UnmarshalMap to convert to nil struct.
	var actual2 *struct{ A int }
	expected = &InvalidUnmarshalError{Type: reflect.TypeOf(actual2)}
	if err := UnmarshalMap(nil, actual2); err == nil {
		t.Errorf("UnmarshalMap with input %T returned no error, expected error `%v`", actual2, expected)
	} else if err.Error() != expected.Error() {
		t.Errorf("UnmarshalMap with input %T returned error `%v`, expected error `%v`", actual2, err, expected)
	}
}

func Test_New_MarshalList(t *testing.T) {
	for name, c := range marshallerListTestInputs {
		t.Run(name, func(t *testing.T) {
			actual, err := MarshalList(c.input)
			if c.err != nil {
				if err == nil {
					t.Fatalf("marshalList with input %#v returned %#v, expected error `%s`",
						c.input, actual, c.err)
				} else if err.Error() != c.err.Error() {
					t.Fatalf("marshalList with input %#v returned error `%s`, expected error `%s`",
						c.input, err, c.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("MarshalList with input %#v returned error `%s`", c.input, err)
			}

			compareObjects(t, c.expected, actual)

		})
	}
}

func Test_New_UnmarshalList(t *testing.T) {
	// Using the same inputs from MarshalList, test the reverse mapping.
	for name, c := range marshallerListTestInputs {
		t.Run(name, func(t *testing.T) {
			if c.input == nil {
				t.Skip()
			}

			iv := reflect.ValueOf(c.input)

			actual := reflect.New(iv.Type())
			if iv.Kind() == reflect.Slice {
				actual.Elem().Set(reflect.MakeSlice(iv.Type(), iv.Len(), iv.Cap()))
			}

			if err := UnmarshalList(c.expected.([]types.AttributeValue), actual.Interface()); err != nil {
				t.Errorf("unmarshal with input %#v returned error `%s`", c.expected, err)
			}
			compareObjects(t, c.input, actual.Elem().Interface())
		})
	}
}

func Test_New_UnmarshalListError(t *testing.T) {
	// Test that we get an error using UnmarshalList to convert to a nil value.
	expected := &InvalidUnmarshalError{Type: reflect.TypeOf(nil)}
	if err := UnmarshalList(nil, nil); err == nil {
		t.Errorf("UnmarshalList with input %T returned no error, expected error `%v`", nil, expected)
	} else if err.Error() != expected.Error() {
		t.Errorf("UnmarshalList with input %T returned error `%v`, expected error `%v`", nil, err, expected)
	}

	// Test that we get an error using UnmarshalList to convert to a non-pointer value.
	var actual map[string]interface{}
	expected = &InvalidUnmarshalError{Type: reflect.TypeOf(actual)}
	if err := UnmarshalList(nil, actual); err == nil {
		t.Errorf("UnmarshalList with input %T returned no error, expected error `%v`", actual, expected)
	} else if err.Error() != expected.Error() {
		t.Errorf("UnmarshalList with input %T returned error `%v`, expected error `%v`", actual, err, expected)
	}

	// Test that we get an error using UnmarshalList to convert to nil struct.
	var actual2 *struct{ A int }
	expected = &InvalidUnmarshalError{Type: reflect.TypeOf(actual2)}
	if err := UnmarshalList(nil, actual2); err == nil {
		t.Errorf("UnmarshalList with input %T returned no error, expected error `%v`", actual2, expected)
	} else if err.Error() != expected.Error() {
		t.Errorf("UnmarshalList with input %T returned error `%v`, expected error `%v`", actual2, err, expected)
	}
}

func compareObjects(t *testing.T, expected interface{}, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		ev := reflect.ValueOf(expected)
		av := reflect.ValueOf(actual)
		if diff := cmp.Diff(expected, actual); len(diff) != 0 {
			t.Errorf("expect kind(%s, %T) match actual kind(%s, %T)\n%s",
				ev.Kind(), ev.Interface(), av.Kind(), av.Interface(), diff)
		}
	}
}

func BenchmarkMarshalOneMember(b *testing.B) {
	fieldCache = fieldCacher{}

	simple := simpleMarshalStruct{
		String:  "abc",
		Int:     123,
		Uint:    123,
		Float32: 123.321,
		Float64: 123.321,
		Bool:    true,
		Null:    nil,
	}
	type MyCompositeStruct struct {
		A simpleMarshalStruct `dynamodbav:"a"`
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Marshal(MyCompositeStruct{
				A: simple,
			}); err != nil {
				b.Error("unexpected error:", err)
			}
		}
	})
}

func BenchmarkMarshalTwoMembers(b *testing.B) {
	fieldCache = fieldCacher{}

	simple := simpleMarshalStruct{
		String:  "abc",
		Int:     123,
		Uint:    123,
		Float32: 123.321,
		Float64: 123.321,
		Bool:    true,
		Null:    nil,
	}

	type MyCompositeStruct struct {
		A simpleMarshalStruct `dynamodbav:"a"`
		B simpleMarshalStruct `dynamodbav:"b"`
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Marshal(MyCompositeStruct{
				A: simple,
				B: simple,
			}); err != nil {
				b.Error("unexpected error:", err)
			}
		}
	})
}

func BenchmarkUnmarshalOneMember(b *testing.B) {
	fieldCache = fieldCacher{}

	myStructAVMap, _ := Marshal(simpleMarshalStruct{
		String:  "abc",
		Int:     123,
		Uint:    123,
		Float32: 123.321,
		Float64: 123.321,
		Bool:    true,
		Null:    nil,
	})

	type MyCompositeStructOne struct {
		A simpleMarshalStruct `dynamodbav:"a"`
	}
	var out MyCompositeStructOne
	avMap := map[string]types.AttributeValue{
		"a": myStructAVMap,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := Unmarshal(&types.AttributeValueMemberM{Value: avMap}, &out); err != nil {
				b.Error("unexpected error:", err)
			}
		}
	})
}

func BenchmarkUnmarshalTwoMembers(b *testing.B) {
	fieldCache = fieldCacher{}

	myStructAVMap, _ := Marshal(simpleMarshalStruct{
		String:  "abc",
		Int:     123,
		Uint:    123,
		Float32: 123.321,
		Float64: 123.321,
		Bool:    true,
		Null:    nil,
	})

	type MyCompositeStructTwo struct {
		A simpleMarshalStruct `dynamodbav:"a"`
		B simpleMarshalStruct `dynamodbav:"b"`
	}
	var out MyCompositeStructTwo
	avMap := map[string]types.AttributeValue{
		"a": myStructAVMap,
		"b": myStructAVMap,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := Unmarshal(&types.AttributeValueMemberM{Value: avMap}, &out); err != nil {
				b.Error("unexpected error:", err)
			}
		}
	})
}

func Test_Encode_YAML_TagKey(t *testing.T) {
	input := struct {
		String      string         `yaml:"string"`
		EmptyString string         `yaml:"empty"`
		OmitString  string         `yaml:"omitted,omitempty"`
		Ignored     string         `yaml:"-"`
		Byte        []byte         `yaml:"byte"`
		Float32     float32        `yaml:"float32"`
		Float64     float64        `yaml:"float64"`
		Int         int            `yaml:"int"`
		Uint        uint           `yaml:"uint"`
		Slice       []string       `yaml:"slice"`
		Map         map[string]int `yaml:"map"`
		NoTag       string
	}{
		String:  "String",
		Ignored: "Ignored",
		Slice:   []string{"one", "two"},
		Map: map[string]int{
			"one": 1,
			"two": 2,
		},
		NoTag: "NoTag",
	}

	expected := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"string":  &types.AttributeValueMemberS{Value: "String"},
			"empty":   &types.AttributeValueMemberS{Value: ""},
			"byte":    &types.AttributeValueMemberNULL{Value: true},
			"float32": &types.AttributeValueMemberN{Value: "0"},
			"float64": &types.AttributeValueMemberN{Value: "0"},
			"int":     &types.AttributeValueMemberN{Value: "0"},
			"uint":    &types.AttributeValueMemberN{Value: "0"},
			"slice": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "one"},
					&types.AttributeValueMemberS{Value: "two"},
				},
			},
			"map": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"one": &types.AttributeValueMemberN{Value: "1"},
					"two": &types.AttributeValueMemberN{Value: "2"},
				},
			},
			"NoTag": &types.AttributeValueMemberS{Value: "NoTag"},
		},
	}

	enc := NewEncoder(func(o *EncoderOptions) {
		o.TagKey = "yaml"
	})

	actual, err := enc.Encode(input)
	if err != nil {
		t.Errorf("Encode with input %#v returned error `%s`, expected nil", input, err)
	}

	compareObjects(t, expected, actual)
}
