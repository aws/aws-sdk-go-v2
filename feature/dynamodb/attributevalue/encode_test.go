package attributevalue

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	smithydocument "github.com/aws/smithy-go/document"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalShared(t *testing.T) {
	for name, c := range sharedTestCases {
		t.Run(name, func(t *testing.T) {
			av, err := Marshal(c.expected)
			assertConvertTest(t, av, c.in, err, c.err)
		})
	}
}

func TestMarshalListShared(t *testing.T) {
	for name, c := range sharedListTestCases {
		t.Run(name, func(t *testing.T) {
			av, err := MarshalList(c.expected)
			assertConvertTest(t, av, c.in, err, c.err)
		})
	}
}

func TestMarshalMapShared(t *testing.T) {
	for name, c := range sharedMapTestCases {
		t.Run(name, func(t *testing.T) {
			av, err := MarshalMap(c.expected)
			assertConvertTest(t, av, c.in, err, c.err)
		})
	}
}

type marshalMarshaler struct {
	Value  string
	Value2 int
	Value3 bool
	Value4 time.Time
}

func (m *marshalMarshaler) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"abc": &types.AttributeValueMemberS{Value: m.Value},
			"def": &types.AttributeValueMemberN{Value: strconv.Itoa(m.Value2)},
			"ghi": &types.AttributeValueMemberBOOL{Value: m.Value3},
			"jkl": &types.AttributeValueMemberS{Value: m.Value4.Format(time.RFC3339Nano)},
		},
	}, nil
}

func TestMarshalMashaler(t *testing.T) {
	m := &marshalMarshaler{
		Value:  "value",
		Value2: 123,
		Value3: true,
		Value4: testDate,
	}

	expect := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"abc": &types.AttributeValueMemberS{Value: "value"},
			"def": &types.AttributeValueMemberN{Value: "123"},
			"ghi": &types.AttributeValueMemberBOOL{Value: true},
			"jkl": &types.AttributeValueMemberS{Value: "2016-05-03T17:06:26.209072Z"},
		},
	}

	actual, err := Marshal(m)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if e, a := expect, actual; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}

type testOmitEmptyElemListStruct struct {
	Values []string `dynamodbav:",omitemptyelem"`
}

type testOmitEmptyElemMapStruct struct {
	Values map[string]interface{} `dynamodbav:",omitemptyelem"`
}

func TestMarshalListOmitEmptyElem(t *testing.T) {
	expect := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Values": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "abc"},
				&types.AttributeValueMemberS{Value: "123"},
			}},
		},
	}

	m := testOmitEmptyElemListStruct{Values: []string{"abc", "", "123"}}

	actual, err := Marshal(m)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	if diff := cmp.Diff(expect, actual, cmpopts.IgnoreTypes(smithydocument.NoSerde{})); len(diff) != 0 {
		t.Errorf("expect match\n%s", diff)
	}
}

func TestMarshalMapOmitEmptyElem(t *testing.T) {
	expect := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Values": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"abc": &types.AttributeValueMemberN{Value: "123"},
				"hij": &types.AttributeValueMemberS{Value: ""},
				"klm": &types.AttributeValueMemberS{Value: "abc"},
				"qrs": &types.AttributeValueMemberS{Value: "abc"},
			}},
		},
	}

	m := testOmitEmptyElemMapStruct{Values: map[string]interface{}{
		"abc": 123.,
		"efg": nil,
		"hij": "",
		"klm": "abc",
		"nop": func() interface{} {
			var v *string
			return v
		}(),
		"qrs": func() interface{} {
			v := "abc"
			return &v
		}(),
	}}

	actual, err := Marshal(m)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	if diff := cmp.Diff(expect, actual, cmpopts.IgnoreTypes(smithydocument.NoSerde{})); len(diff) != 0 {
		t.Errorf("expect match\n%s", diff)
	}
}

type testNullEmptyElemListStruct struct {
	Values []string `dynamodbav:",nullemptyelem"`
}

type testNullEmptyElemMapStruct struct {
	Values map[string]interface{} `dynamodbav:",nullemptyelem"`
}

func TestMarshalListNullEmptyElem(t *testing.T) {
	expect := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Values": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "abc"},
				&types.AttributeValueMemberNULL{Value: true},
				&types.AttributeValueMemberS{Value: "123"},
			}},
		},
	}

	m := testNullEmptyElemListStruct{Values: []string{"abc", "", "123"}}

	actual, err := Marshal(m)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	if diff := cmp.Diff(expect, actual, cmpopts.IgnoreTypes(smithydocument.NoSerde{})); len(diff) != 0 {
		t.Errorf("expect match\n%s", diff)
	}
}

func TestMarshalMapNullEmptyElem(t *testing.T) {
	expect := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Values": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"abc": &types.AttributeValueMemberN{Value: "123"},
				"efg": &types.AttributeValueMemberNULL{Value: true},
				"hij": &types.AttributeValueMemberS{Value: ""},
				"klm": &types.AttributeValueMemberS{Value: "abc"},
				"nop": &types.AttributeValueMemberNULL{Value: true},
				"qrs": &types.AttributeValueMemberS{Value: "abc"},
			}},
		},
	}

	m := testNullEmptyElemMapStruct{Values: map[string]interface{}{
		"abc": 123.,
		"efg": nil,
		"hij": "",
		"klm": "abc",
		"nop": func() interface{} {
			var v *string
			return v
		}(),
		"qrs": func() interface{} {
			v := "abc"
			return &v
		}(),
	}}

	actual, err := Marshal(m)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	if diff := cmp.Diff(expect, actual, cmpopts.IgnoreTypes(smithydocument.NoSerde{})); len(diff) != 0 {
		t.Errorf("expect match\n%s", diff)
	}
}

type testOmitEmptyScalar struct {
	IntZero       int  `dynamodbav:",omitempty"`
	IntPtrNil     *int `dynamodbav:",omitempty"`
	IntPtrSetZero *int `dynamodbav:",omitempty"`
}

func TestMarshalOmitEmpty(t *testing.T) {
	expect := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"IntPtrSetZero": &types.AttributeValueMemberN{Value: "0"},
		},
	}

	m := testOmitEmptyScalar{IntPtrSetZero: aws.Int(0)}

	actual, err := Marshal(m)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	if e, a := expect, actual; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEncodeEmbeddedPointerStruct(t *testing.T) {
	type B struct {
		Bint int
	}
	type C struct {
		Cint int
	}
	type A struct {
		Aint int
		*B
		*C
	}
	a := A{Aint: 321, B: &B{123}}
	if e, a := 321, a.Aint; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 123, a.Bint; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if a.C != nil {
		t.Errorf("expect nil, got %v", a.C)
	}

	actual, err := Marshal(a)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	expect := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Aint": &types.AttributeValueMemberN{Value: "321"},
			"Bint": &types.AttributeValueMemberN{Value: "123"},
		},
	}
	if e, a := expect, actual; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEncodeUnixTime(t *testing.T) {
	type A struct {
		Normal time.Time
		Tagged time.Time `dynamodbav:",unixtime"`
		Typed  UnixTime
	}

	a := A{
		Normal: time.Unix(123, 0).UTC(),
		Tagged: time.Unix(456, 0),
		Typed:  UnixTime(time.Unix(789, 0)),
	}

	actual, err := Marshal(a)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	expect := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Normal": &types.AttributeValueMemberS{Value: "1970-01-01T00:02:03Z"},
			"Tagged": &types.AttributeValueMemberN{Value: "456"},
			"Typed":  &types.AttributeValueMemberN{Value: "789"},
		},
	}
	if e, a := expect, actual; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestUnixTimeString(t *testing.T) {
	gotime := time.Date(2016, time.May, 03, 17, 06, 26, 0, time.UTC)
	ddbtime := UnixTime(gotime)
	if fmt.Sprint(gotime) != fmt.Sprint(ddbtime) {
		t.Error("UnixTime.String not equal to time.Time.String")
	}
}

type AliasedTime time.Time

func TestEncodeAliasedUnixTime(t *testing.T) {
	type A struct {
		Normal AliasedTime
		Tagged AliasedTime `dynamodbav:",unixtime"`
	}

	a := A{
		Normal: AliasedTime(time.Unix(123, 0).UTC()),
		Tagged: AliasedTime(time.Unix(456, 0)),
	}

	actual, err := Marshal(a)
	if err != nil {
		t.Errorf("expect no err, got %v", err)
	}
	expect := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Normal": &types.AttributeValueMemberS{Value: "1970-01-01T00:02:03Z"},
			"Tagged": &types.AttributeValueMemberN{Value: "456"},
		},
	}
	if e, a := expect, actual; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestMarshalTime_S(t *testing.T) {
	type A struct {
		TimeField   time.Time
		TimeFieldsL []time.Time
	}
	cases := map[string]struct {
		input      time.Time
		expect     string
		encodeTime func(time.Time) (types.AttributeValue, error)
	}{
		"String RFC3339Nano (Default)": {
			input:  time.Unix(123, 10000000).UTC(),
			expect: "1970-01-01T00:02:03.01Z",
		},
		"String UnixDate": {
			input:  time.Unix(123, 0).UTC(),
			expect: "Thu Jan  1 00:02:03 UTC 1970",
			encodeTime: func(t time.Time) (types.AttributeValue, error) {
				return &types.AttributeValueMemberS{
					Value: t.Format(time.UnixDate),
				}, nil
			},
		},
		"String RFC3339 millis keeping zeroes": {
			input:  time.Unix(123, 10000000).UTC(),
			expect: "1970-01-01T00:02:03.010Z",
			encodeTime: func(t time.Time) (types.AttributeValue, error) {
				return &types.AttributeValueMemberS{
					Value: t.Format("2006-01-02T15:04:05.000Z07:00"), // Would be RFC3339 millis with zeroes
				}, nil
			},
		},
		"String RFC822": {
			input:  time.Unix(120, 0).UTC(),
			expect: "01 Jan 70 00:02 UTC",
			encodeTime: func(t time.Time) (types.AttributeValue, error) {
				return &types.AttributeValueMemberS{
					Value: t.Format(time.RFC822),
				}, nil
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			inputValue := A{
				TimeField:   c.input,
				TimeFieldsL: []time.Time{c.input},
			}
			actual, err := MarshalWithOptions(inputValue, func(eo *EncoderOptions) {
				if c.encodeTime != nil {
					eo.EncodeTime = c.encodeTime
				}
			})
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			expectedValue := &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"TimeField": &types.AttributeValueMemberS{Value: c.expect},
					"TimeFieldsL": &types.AttributeValueMemberL{Value: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: c.expect},
					}},
				},
			}
			if diff := cmp.Diff(expectedValue, actual, getIgnoreAVUnexportedOptions()...); diff != "" {
				t.Errorf("expect attribute value match\n%s", diff)
			}
		})
	}
}

func TestMarshalTime_N(t *testing.T) {
	type A struct {
		TimeField   time.Time
		TimeFieldsL []time.Time
	}
	cases := map[string]struct {
		input      time.Time
		expect     string
		encodeTime func(time.Time) (types.AttributeValue, error)
	}{
		"Number Unix seconds": {
			input:  time.Unix(123, 10000000).UTC(),
			expect: "123",
			encodeTime: func(t time.Time) (types.AttributeValue, error) {
				return &types.AttributeValueMemberN{
					Value: strconv.Itoa(int(t.Unix())),
				}, nil
			},
		},
		"Number Unix milli": {
			input:  time.Unix(123, 10000000).UTC(),
			expect: "123010",
			encodeTime: func(t time.Time) (types.AttributeValue, error) {
				return &types.AttributeValueMemberN{
					Value: strconv.Itoa(int(t.UnixNano() / int64(time.Millisecond))),
				}, nil
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			inputValue := A{
				TimeField:   c.input,
				TimeFieldsL: []time.Time{c.input},
			}
			actual, err := MarshalWithOptions(inputValue, func(eo *EncoderOptions) {
				if c.encodeTime != nil {
					eo.EncodeTime = c.encodeTime
				}
			})
			if err != nil {
				t.Errorf("expect no err, got %v", err)
			}
			expectedValue := &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"TimeField": &types.AttributeValueMemberN{Value: c.expect},
					"TimeFieldsL": &types.AttributeValueMemberL{Value: []types.AttributeValue{
						&types.AttributeValueMemberN{Value: c.expect},
					}},
				},
			}
			if diff := cmp.Diff(expectedValue, actual, getIgnoreAVUnexportedOptions()...); diff != "" {
				t.Errorf("expect attribute value match\n%s", diff)
			}
		})
	}
}

func TestEncoderFieldByIndex(t *testing.T) {
	type (
		Middle struct{ Inner int }
		Outer  struct{ *Middle }
	)

	// nil embedded struct
	outer := Outer{}
	outerFields := unionStructFields(reflect.TypeOf(outer), structFieldOptions{})
	innerField, _ := outerFields.FieldByName("Inner")

	_, found := encoderFieldByIndex(reflect.ValueOf(&outer).Elem(), innerField.Index)
	if found != false {
		t.Error("expected found to be false when embedded struct is nil")
	}

	// non-nil embedded struct
	outer = Outer{Middle: &Middle{Inner: 3}}
	outerFields = unionStructFields(reflect.TypeOf(outer), structFieldOptions{})
	innerField, _ = outerFields.FieldByName("Inner")

	f, found := encoderFieldByIndex(reflect.ValueOf(&outer).Elem(), innerField.Index)
	if !found {
		t.Error("expected found to be true")
	}
	if f.Kind() != reflect.Int || f.Int() != int64(outer.Inner) {
		t.Error("expected f to be of kind Int with value equal to outer.Inner")
	}
}

func TestMarshalMap_keyTypes(t *testing.T) {
	type StrAlias string
	type IntAlias int
	type BoolAlias bool

	cases := map[string]struct {
		input    interface{}
		expectAV map[string]types.AttributeValue
	}{
		"string key": {
			input: map[string]interface{}{
				"a": 123,
				"b": "efg",
			},
			expectAV: map[string]types.AttributeValue{
				"a": &types.AttributeValueMemberN{Value: "123"},
				"b": &types.AttributeValueMemberS{Value: "efg"},
			},
		},
		"string alias key": {
			input: map[StrAlias]interface{}{
				"a": 123,
				"b": "efg",
			},
			expectAV: map[string]types.AttributeValue{
				"a": &types.AttributeValueMemberN{Value: "123"},
				"b": &types.AttributeValueMemberS{Value: "efg"},
			},
		},
		"Number key": {
			input: map[Number]interface{}{
				Number("1"): 123,
				Number("2"): "efg",
			},
			expectAV: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberN{Value: "123"},
				"2": &types.AttributeValueMemberS{Value: "efg"},
			},
		},
		"int key": {
			input: map[int]interface{}{
				1: 123,
				2: "efg",
			},
			expectAV: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberN{Value: "123"},
				"2": &types.AttributeValueMemberS{Value: "efg"},
			},
		},
		"int alias key": {
			input: map[IntAlias]interface{}{
				1: 123,
				2: "efg",
			},
			expectAV: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberN{Value: "123"},
				"2": &types.AttributeValueMemberS{Value: "efg"},
			},
		},
		"bool key": {
			input: map[bool]interface{}{
				true:  123,
				false: "efg",
			},
			expectAV: map[string]types.AttributeValue{
				"true":  &types.AttributeValueMemberN{Value: "123"},
				"false": &types.AttributeValueMemberS{Value: "efg"},
			},
		},
		"bool alias key": {
			input: map[BoolAlias]interface{}{
				true:  123,
				false: "efg",
			},
			expectAV: map[string]types.AttributeValue{
				"true":  &types.AttributeValueMemberN{Value: "123"},
				"false": &types.AttributeValueMemberS{Value: "efg"},
			},
		},
		"textMarshaler key": {
			input: map[testTextMarshaler]interface{}{
				{Foo: "1"}: 123,
				{Foo: "2"}: "efg",
			},
			expectAV: map[string]types.AttributeValue{
				"Foo:1": &types.AttributeValueMemberN{Value: "123"},
				"Foo:2": &types.AttributeValueMemberS{Value: "efg"},
			},
		},
		"textMarshaler ptr key": {
			input: map[*testTextMarshaler]interface{}{
				{Foo: "1"}: 123,
				{Foo: "2"}: "efg",
			},
			expectAV: map[string]types.AttributeValue{
				"Foo:1": &types.AttributeValueMemberN{Value: "123"},
				"Foo:2": &types.AttributeValueMemberS{Value: "efg"},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			av, err := MarshalMap(c.input)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if diff := cmp.Diff(c.expectAV, av, getIgnoreAVUnexportedOptions()...); diff != "" {
				t.Errorf("expect attribute value match\n%s", diff)
			}
		})
	}
}
