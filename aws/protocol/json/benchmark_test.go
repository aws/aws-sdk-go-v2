package json_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/json"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
	v1Encoder "github.com/aws/aws-sdk-go-v2/private/protocol/json"
	reflectEncoder "github.com/aws/aws-sdk-go-v2/private/protocol/json/jsonutil"
)

type testEnum string

type testOperationInput struct {
	_            struct{}          `type:"structure"`
	StringValue  *string           `locationName:"stringValue" type:"string"`
	IntegerValue *int64            `locationName:"integerValue" type:"integer"`
	EnumValue    testEnum          `locationName:"enumValue" type:"string" enum:"true"`
	FloatValue   *float64          `locationName:"floatValue" type:"double"`
	ListValue    []nestedShape     `locationName:"listValue" type:"list"`
	ShapeValue   *nestedShape      `locationName:"shapeValue" type:"structure"`
	MapValue     map[string]string `locationName:"mapValue" type:"map"`
	ByteSlice    []byte            `locationName:"byteSlice" type:"blob"`
}

func (v *testOperationInput) MarshalFields(encoder protocol.FieldEncoder) error {
	meta := protocol.Metadata{}

	if v.StringValue != nil {
		encoder.SetValue(protocol.BodyTarget, "stringValue", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*v.StringValue)}, meta)
	}

	if v.IntegerValue != nil {
		encoder.SetValue(protocol.BodyTarget, "integerValue", protocol.Int64Value(*v.IntegerValue), meta)
	}

	if len(v.EnumValue) > 0 {
		encoder.SetValue(protocol.BodyTarget, "enumValue", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v.EnumValue)}, meta)
	}

	if v.FloatValue != nil {
		encoder.SetValue(protocol.BodyTarget, "floatValue", protocol.Float64Value(*v.FloatValue), meta)
	}

	if v.ListValue != nil {
		listEncoder := encoder.List(protocol.BodyTarget, "listValue", meta)
		listEncoder.Start()
		for i := range v.ListValue {
			listEncoder.ListAddFields(&v.ListValue[i])
		}
		listEncoder.End()
	}

	if v.ShapeValue != nil {
		encoder.SetFields(protocol.BodyTarget, "shapeValue", v.ShapeValue, meta)
	}

	if v.MapValue != nil {
		mapEncoder := encoder.Map(protocol.BodyTarget, "mapValue", meta)
		mapEncoder.Start()
		for k := range v.MapValue {
			mapEncoder.MapSetValue(k, protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v.MapValue[k])})
		}
		mapEncoder.End()
	}

	if v.ByteSlice != nil {
		encoder.SetValue(protocol.BodyTarget, "byteSlice", protocol.QuotedValue{ValueMarshaler: protocol.BytesValue(v.ByteSlice)}, meta)
	}

	return nil
}

type nestedShape struct {
	_           struct{} `type:"structure"`
	StringValue *string  `locationName:"stringValue" type:"string"`
}

func (v *nestedShape) MarshalFields(encoder protocol.FieldEncoder) error {
	meta := protocol.Metadata{}

	if v.StringValue != nil {
		encoder.SetValue(protocol.BodyTarget, "stringValue", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*v.StringValue)}, meta)
	}

	return nil
}

func MarshalTestOperationInputAWSJSON(v *testOperationInput, e *json.Encoder) []byte {
	marshalTestOperationInputAWSJSON(v, &e.Value)

	return e.Bytes()
}

func marshalTestOperationInputAWSJSON(v *testOperationInput, j *json.Value) {
	object := j.Object()
	defer object.Close()

	if v.StringValue != nil {
		object.Key("stringValue").String(*v.StringValue)
	}

	if v.IntegerValue != nil {
		object.Key("integerValue").Integer(*v.IntegerValue)
	}

	if len(v.EnumValue) > 0 {
		object.Key("enumValue").String(string(v.EnumValue))
	}

	if v.FloatValue != nil {
		object.Key("floatValue").Float(*v.FloatValue)
	}

	if v.ListValue != nil {
		value := object.Key("listValue")
		marshalListValueShapeAWSREST(v.ListValue, &value)
	}

	if v.ShapeValue != nil {
		value := object.Key("shapeValue")
		marshalNestedShapeAWSREST(v.ShapeValue, &value)
	}

	if v.MapValue != nil {
		value := object.Key("mapValue")
		marshalMapShapeAWSREST(v.MapValue, &value)
	}

	if v.ByteSlice != nil {
		object.Key("byteSlice").ByteSlice(v.ByteSlice)
	}
}

func marshalListValueShapeAWSREST(v []nestedShape, j *json.Value) {
	array := j.Array()
	defer array.Close()

	for i := range v {
		arrayValue := array.Value()
		marshalNestedShapeAWSREST(&v[i], &arrayValue)
	}
}

func marshalNestedShapeAWSREST(v *nestedShape, j *json.Value) {
	object := j.Object()
	defer object.Close()

	if v.StringValue != nil {
		object.Key("stringValue").String(*v.StringValue)
	}
}

func marshalMapShapeAWSREST(v map[string]string, j *json.Value) {
	object := j.Object()
	defer object.Close()

	for k := range v {
		object.Key(k).String(v[k])
	}
}

var testOperationCases = [...]*testOperationInput{
	0: {},
	1: {
		StringValue:  aws.String("someStringValue1"),
		IntegerValue: aws.Int64(42),
		EnumValue:    "SOME_ENUM",
		FloatValue:   aws.Float64(3.14),
		ListValue: []nestedShape{
			{StringValue: aws.String("someStringValue2")},
			{},
		},
		ShapeValue: &nestedShape{StringValue: aws.String("someStringValue3")},
		MapValue: map[string]string{
			"someMapKey": "someMapValue",
		},
		ByteSlice: make([]byte, 1024),
	},
}

func BenchmarkEncoderV2(b *testing.B) {
	for i, operationCase := range testOperationCases {
		b.Run(fmt.Sprintf("Case%d", i), func(b *testing.B) {
			encoder := json.NewEncoder()
			_ = MarshalTestOperationInputAWSJSON(operationCase, encoder)
		})
	}
}

func BenchmarkEncoderV1(b *testing.B) {
	for i, operationCase := range testOperationCases {
		b.Run(fmt.Sprintf("Case%d", i), func(b *testing.B) {
			encoder := v1Encoder.NewEncoder()
			err := operationCase.MarshalFields(encoder)
			b.StopTimer()
			if err != nil {
				b.Fatal(err)
			}
			_, err = encoder.Encode()
			if err != nil {
				b.Fatal(err)
			}
		})
	}
}

func BenchmarkEncoderReflection(b *testing.B) {
	for i, operationCase := range testOperationCases {
		b.Run(fmt.Sprintf("Case%d", i), func(b *testing.B) {
			_, err := reflectEncoder.BuildJSON(operationCase)
			if err != nil {
				b.Fatal(err)
			}
		})
	}
}
