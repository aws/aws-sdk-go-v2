package attributevalue

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
	"github.com/google/go-cmp/cmp"
)

func TestUnmarshalShared(t *testing.T) {
	for name, c := range sharedTestCases {
		t.Run(name, func(t *testing.T) {
			err := Unmarshal(c.in, c.actual)
			assertConvertTest(t, c.actual, c.expected, err, c.err)
		})
	}
}

func TestUnmarshal(t *testing.T) {
	cases := []struct {
		in               types.AttributeValue
		actual, expected interface{}
		err              error
	}{
		//------------
		// Sets
		//------------
		{
			in: &types.AttributeValueMemberBS{Value: [][]byte{
				{48, 49}, {50, 51},
			}},
			actual:   &[][]byte{},
			expected: [][]byte{{48, 49}, {50, 51}},
		},
		{
			in: &types.AttributeValueMemberNS{Value: []string{
				"123", "321",
			}},
			actual:   &[]int{},
			expected: []int{123, 321},
		},
		{
			in: &types.AttributeValueMemberNS{Value: []string{
				"123", "321",
			}},
			actual:   &[]interface{}{},
			expected: []interface{}{123., 321.},
		},
		{
			in: &types.AttributeValueMemberSS{Value: []string{
				"abc", "123",
			}},
			actual:   &[]string{},
			expected: &[]string{"abc", "123"},
		},
		{
			in: &types.AttributeValueMemberSS{Value: []string{
				"abc", "123",
			}},
			actual:   &[]*string{},
			expected: &[]*string{aws.String("abc"), aws.String("123")},
		},
		//------------
		// Interfaces
		//------------
		{
			in: &types.AttributeValueMemberB{Value: []byte{48, 49}},
			actual: func() interface{} {
				var v interface{}
				return &v
			}(),
			expected: []byte{48, 49},
		},
		{
			in: &types.AttributeValueMemberBS{Value: [][]byte{
				{48, 49}, {50, 51},
			}},
			actual: func() interface{} {
				var v interface{}
				return &v
			}(),
			expected: [][]byte{{48, 49}, {50, 51}},
		},
		{
			in: &types.AttributeValueMemberBOOL{Value: true},
			actual: func() interface{} {
				var v interface{}
				return &v
			}(),
			expected: bool(true),
		},
		{
			in: &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "abc"},
				&types.AttributeValueMemberS{Value: "123"},
			}},
			actual: func() interface{} {
				var v interface{}
				return &v
			}(),
			expected: []interface{}{"abc", "123"},
		},
		{
			in: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"123": &types.AttributeValueMemberS{Value: "abc"},
				"abc": &types.AttributeValueMemberS{Value: "123"},
			}},
			actual: func() interface{} {
				var v interface{}
				return &v
			}(),
			expected: map[string]interface{}{"123": "abc", "abc": "123"},
		},
		{
			in: &types.AttributeValueMemberN{Value: "123"},
			actual: func() interface{} {
				var v interface{}
				return &v
			}(),
			expected: float64(123),
		},
		{
			in: &types.AttributeValueMemberNS{Value: []string{
				"123", "321",
			}},
			actual: func() interface{} {
				var v interface{}
				return &v
			}(),
			expected: []float64{123., 321.},
		},
		{
			in: &types.AttributeValueMemberS{Value: "123"},
			actual: func() interface{} {
				var v interface{}
				return &v
			}(),
			expected: "123",
		},
		{
			in: &types.AttributeValueMemberNULL{Value: true},
			actual: func() interface{} {
				var v string
				return &v
			}(),
			expected: "",
		},
		{
			in: &types.AttributeValueMemberNULL{Value: true},
			actual: func() interface{} {
				v := new(string)
				return &v
			}(),
			expected: nil,
		},
		{
			in: &types.AttributeValueMemberS{Value: ""},
			actual: func() interface{} {
				v := new(string)
				return &v
			}(),
			expected: aws.String(""),
		},
		{
			in: &types.AttributeValueMemberSS{Value: []string{
				"123", "321",
			}},
			actual: func() interface{} {
				var v interface{}
				return &v
			}(),
			expected: []string{"123", "321"},
		},
		{
			in: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"abc": &types.AttributeValueMemberS{Value: "123"},
				"Cba": &types.AttributeValueMemberS{Value: "321"},
			}},
			actual:   &struct{ Abc, Cba string }{},
			expected: struct{ Abc, Cba string }{Abc: "123", Cba: "321"},
		},
		{
			in:     &types.AttributeValueMemberN{Value: "512"},
			actual: new(uint8),
			err: &UnmarshalTypeError{
				Value: fmt.Sprintf("number overflow, 512"),
				Type:  reflect.TypeOf(uint8(0)),
			},
		},
		// -------
		// Empty Values
		// -------
		{
			in:       &types.AttributeValueMemberB{Value: []byte{}},
			actual:   &[]byte{},
			expected: []byte{},
		},
		{
			in:       &types.AttributeValueMemberBS{Value: [][]byte{}},
			actual:   &[][]byte{},
			expected: [][]byte{},
		},
		{
			in:       &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
			actual:   &[]interface{}{},
			expected: []interface{}{},
		},
		{
			in:       &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
			actual:   &map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			in:     &types.AttributeValueMemberN{Value: ""},
			actual: new(int),
			err:    fmt.Errorf("invalid syntax"),
		},
		{
			in:       &types.AttributeValueMemberNS{Value: []string{}},
			actual:   &[]string{},
			expected: []string{},
		},
		{
			in:       &types.AttributeValueMemberS{Value: ""},
			actual:   new(string),
			expected: "",
		},
		{
			in:       &types.AttributeValueMemberSS{Value: []string{}},
			actual:   &[]string{},
			expected: []string{},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d/%d", i, len(cases)), func(t *testing.T) {
			err := Unmarshal(c.in, c.actual)
			assertConvertTest(t, c.actual, c.expected, err, c.err)
		})
	}
}

func TestInterfaceInput(t *testing.T) {
	var v interface{}
	expected := []interface{}{"abc", "123"}
	err := Unmarshal(&types.AttributeValueMemberL{Value: []types.AttributeValue{
		&types.AttributeValueMemberS{Value: "abc"},
		&types.AttributeValueMemberS{Value: "123"},
	}}, &v)
	assertConvertTest(t, v, expected, err, nil)
}

func TestUnmarshalError(t *testing.T) {
	cases := map[string]struct {
		in               types.AttributeValue
		actual, expected interface{}
		err              error
	}{
		"invalid unmarshal": {
			in:       nil,
			actual:   int(0),
			expected: nil,
			err:      &InvalidUnmarshalError{Type: reflect.TypeOf(int(0))},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := Unmarshal(c.in, c.actual)
			assertConvertTest(t, c.actual, c.expected, err, c.err)
		})
	}
}

func TestUnmarshalListShared(t *testing.T) {
	for name, c := range sharedListTestCases {
		t.Run(name, func(t *testing.T) {
			err := UnmarshalList(c.in, c.actual)
			assertConvertTest(t, c.actual, c.expected, err, c.err)
		})
	}
}

func TestUnmarshalListError(t *testing.T) {
	cases := map[string]struct {
		in               []types.AttributeValue
		actual, expected interface{}
		err              error
	}{
		"invalid unmarshal": {
			in:       []types.AttributeValue{},
			actual:   []interface{}{},
			expected: nil,
			err:      &InvalidUnmarshalError{Type: reflect.TypeOf([]interface{}{})},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := UnmarshalList(c.in, c.actual)
			assertConvertTest(t, c.actual, c.expected, err, c.err)
		})
	}
}

func TestUnmarshalMapShared(t *testing.T) {
	for name, c := range sharedMapTestCases {
		t.Run(name, func(t *testing.T) {
			err := UnmarshalMap(c.in, c.actual)
			assertConvertTest(t, c.actual, c.expected, err, c.err)
		})
	}
}

func TestUnmarshalMapError(t *testing.T) {
	cases := []struct {
		in               map[string]types.AttributeValue
		actual, expected interface{}
		err              error
	}{
		{
			in:       map[string]types.AttributeValue{},
			actual:   map[string]interface{}{},
			expected: nil,
			err:      &InvalidUnmarshalError{Type: reflect.TypeOf(map[string]interface{}{})},
		},
		{
			in: map[string]types.AttributeValue{
				"BOOL": &types.AttributeValueMemberBOOL{Value: true},
			},
			actual:   &map[int]interface{}{},
			expected: nil,
			err:      &UnmarshalTypeError{Value: "map string key", Type: reflect.TypeOf(int(0))},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := UnmarshalMap(c.in, c.actual)
			assertConvertTest(t, c.actual, c.expected, err, c.err)
		})
	}
}

func TestUnmarshalListOfMaps(t *testing.T) {
	type testItem struct {
		Value  string
		Value2 int
	}

	cases := map[string]struct {
		in               []map[string]types.AttributeValue
		actual, expected interface{}
		err              error
	}{
		"simple map conversion": {
			in: []map[string]types.AttributeValue{
				{
					"Value": &types.AttributeValueMemberBOOL{Value: true},
				},
			},
			actual: &[]map[string]interface{}{},
			expected: []map[string]interface{}{
				{
					"Value": true,
				},
			},
		},
		"attribute to struct": {
			in: []map[string]types.AttributeValue{
				{
					"Value":  &types.AttributeValueMemberS{Value: "abc"},
					"Value2": &types.AttributeValueMemberN{Value: "123"},
				},
			},
			actual: &[]testItem{},
			expected: []testItem{
				{
					Value:  "abc",
					Value2: 123,
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := UnmarshalListOfMaps(c.in, c.actual)
			assertConvertTest(t, c.actual, c.expected, err, c.err)
		})
	}
}

type unmarshalUnmarshaler struct {
	Value  string
	Value2 int
	Value3 bool
	Value4 time.Time
}

func (u *unmarshalUnmarshaler) UnmarshalDynamoDBStreamsAttributeValue(av types.AttributeValue) error {
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok || m == nil {
		return fmt.Errorf("expected AttributeValue to be map")
	}

	if v, ok := m.Value["abc"]; !ok {
		return fmt.Errorf("expected `abc` map key")
	} else if vv, kk := v.(*types.AttributeValueMemberS); !kk || vv == nil {
		return fmt.Errorf("expected `abc` map value string")
	} else {
		u.Value = vv.Value
	}

	if v, ok := m.Value["def"]; !ok {
		return fmt.Errorf("expected `def` map key")
	} else if vv, kk := v.(*types.AttributeValueMemberN); !kk || vv == nil {
		return fmt.Errorf("expected `def` map value number")
	} else {
		n, err := strconv.ParseInt(vv.Value, 10, 64)
		if err != nil {
			return err
		}
		u.Value2 = int(n)
	}

	if v, ok := m.Value["ghi"]; !ok {
		return fmt.Errorf("expected `ghi` map key")
	} else if vv, kk := v.(*types.AttributeValueMemberBOOL); !kk || vv == nil {
		return fmt.Errorf("expected `ghi` map value number")
	} else {
		u.Value3 = vv.Value
	}

	if v, ok := m.Value["jkl"]; !ok {
		return fmt.Errorf("expected `jkl` map key")
	} else if vv, kk := v.(*types.AttributeValueMemberS); !kk || vv == nil {
		return fmt.Errorf("expected `jkl` map value string")
	} else {
		t, err := time.Parse(time.RFC3339, vv.Value)
		if err != nil {
			return err
		}
		u.Value4 = t
	}

	return nil
}

func TestUnmarshalUnmashaler(t *testing.T) {
	u := &unmarshalUnmarshaler{}
	av := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"abc": &types.AttributeValueMemberS{Value: "value"},
			"def": &types.AttributeValueMemberN{Value: "123"},
			"ghi": &types.AttributeValueMemberBOOL{Value: true},
			"jkl": &types.AttributeValueMemberS{Value: "2016-05-03T17:06:26.209072Z"},
		},
	}

	err := Unmarshal(av, u)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if e, a := "value", u.Value; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 123, u.Value2; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := true, u.Value3; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := testDate, u.Value4; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestDecodeUseNumber(t *testing.T) {
	u := map[string]interface{}{}
	av := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"abc": &types.AttributeValueMemberS{Value: "value"},
			"def": &types.AttributeValueMemberN{Value: "123"},
			"ghi": &types.AttributeValueMemberBOOL{Value: true},
		},
	}

	decoder := NewDecoder(func(o *DecoderOptions) {
		o.UseNumber = true
	})
	err := decoder.Decode(av, &u)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if e, a := "value", u["abc"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	n := u["def"].(Number)
	if e, a := "123", n.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := true, u["ghi"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestDecodeUseNumberNumberSet(t *testing.T) {
	u := map[string]interface{}{}
	av := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"ns": &types.AttributeValueMemberNS{
				Value: []string{
					"123", "321",
				},
			},
		},
	}

	decoder := NewDecoder(func(o *DecoderOptions) {
		o.UseNumber = true
	})
	err := decoder.Decode(av, &u)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	ns := u["ns"].([]Number)

	if e, a := "123", ns[0].String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "321", ns[1].String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestDecodeEmbeddedPointerStruct(t *testing.T) {
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
	av := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Aint": &types.AttributeValueMemberN{Value: "321"},
			"Bint": &types.AttributeValueMemberN{Value: "123"},
		},
	}
	decoder := NewDecoder()
	a := A{}
	err := decoder.Decode(av, &a)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := 321, a.Aint; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	// Embedded pointer struct can be created automatically.
	if e, a := 123, a.Bint; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	// But not for absent fields.
	if a.C != nil {
		t.Errorf("expect nil, got %v", a.C)
	}
}

func TestDecodeBooleanOverlay(t *testing.T) {
	type BooleanOverlay bool

	av := &types.AttributeValueMemberBOOL{Value: true}

	decoder := NewDecoder()

	var v BooleanOverlay

	err := decoder.Decode(av, &v)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := BooleanOverlay(true), v; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestDecodeUnixTime(t *testing.T) {
	type A struct {
		Normal time.Time
		Tagged time.Time `dynamodbav:",unixtime"`
		Typed  UnixTime
	}

	expect := A{
		Normal: time.Unix(123, 0).UTC(),
		Tagged: time.Unix(456, 0),
		Typed:  UnixTime(time.Unix(789, 0)),
	}

	input := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Normal": &types.AttributeValueMemberS{Value: "1970-01-01T00:02:03Z"},
			"Tagged": &types.AttributeValueMemberN{Value: "456"},
			"Typed":  &types.AttributeValueMemberN{Value: "789"},
		},
	}
	actual := A{}

	err := Unmarshal(input, &actual)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := expect, actual; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestDecodeAliasedUnixTime(t *testing.T) {
	type A struct {
		Normal AliasedTime
		Tagged AliasedTime `dynamodbav:",unixtime"`
	}

	expect := A{
		Normal: AliasedTime(time.Unix(123, 0).UTC()),
		Tagged: AliasedTime(time.Unix(456, 0)),
	}

	input := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"Normal": &types.AttributeValueMemberS{Value: "1970-01-01T00:02:03Z"},
			"Tagged": &types.AttributeValueMemberN{Value: "456"},
		},
	}
	actual := A{}

	err := Unmarshal(input, &actual)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if expect != actual {
		t.Errorf("expect %v, got %v", expect, actual)
	}
}

// see github issue #1594
func TestDecodeArrayType(t *testing.T) {
	cases := []struct {
		to, from interface{}
	}{
		{
			&[2]int{1, 2},
			&[2]int{},
		},
		{
			&[2]int64{1, 2},
			&[2]int64{},
		},
		{
			&[2]byte{1, 2},
			&[2]byte{},
		},
		{
			&[2]bool{true, false},
			&[2]bool{},
		},
		{
			&[2]string{"1", "2"},
			&[2]string{},
		},
		{
			&[2][]string{{"1", "2"}},
			&[2][]string{},
		},
	}

	for _, c := range cases {
		marshaled, err := Marshal(c.to)
		if err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		if err = Unmarshal(marshaled, c.from); err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		if diff := cmp.Diff(c.to, c.from); len(diff) != 0 {
			t.Errorf("expected match\n:%s", diff)
		}
	}
}

func TestDecoderFieldByIndex(t *testing.T) {
	type (
		Middle struct{ Inner int }
		Outer  struct{ *Middle }
	)
	var outer Outer

	outerType := reflect.TypeOf(outer)
	outerValue := reflect.ValueOf(&outer)
	outerFields := unionStructFields(outerType, structFieldOptions{})
	innerField, _ := outerFields.FieldByName("Inner")

	f := decoderFieldByIndex(outerValue.Elem(), innerField.Index)
	if outer.Middle == nil {
		t.Errorf("expected outer.Middle to be non-nil")
	}
	if f.Kind() != reflect.Int || f.Int() != int64(outer.Inner) {
		t.Error("expected f to be an int with value equal to outer.Inner")
	}
}
