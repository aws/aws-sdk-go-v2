package attributevalue

import (
	"reflect"
	"strconv"
	"testing"
	"time"

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
	if diff := cmp.Diff(expect, actual); len(diff) != 0 {
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
	if diff := cmp.Diff(expect, actual); len(diff) != 0 {
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
	if diff := cmp.Diff(expect, actual); len(diff) != 0 {
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
	if diff := cmp.Diff(expect, actual); len(diff) != 0 {
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
