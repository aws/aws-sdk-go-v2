package attributevalue

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
			err: &UnmarshalTypeError{
				Value: `map key "BOOL"`,
				Type:  reflect.TypeOf(int(0)),
			},
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

func (u *unmarshalUnmarshaler) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
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
func TestDecodeAliasType(t *testing.T) {
	type Str string
	type Int int
	type Uint uint
	type TT struct {
		A Str
		B Int
		C Uint
		S Str
	}

	expect := TT{
		A: "12345",
		B: 12345,
		C: 12345,
		S: "string",
	}
	m := map[string]types.AttributeValue{
		"A": &types.AttributeValueMemberN{
			Value: "12345",
		},
		"B": &types.AttributeValueMemberN{
			Value: "12345",
		},
		"C": &types.AttributeValueMemberN{
			Value: "12345",
		},
		"S": &types.AttributeValueMemberS{
			Value: "string",
		},
	}

	var actual TT
	err := UnmarshalMap(m, &actual)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("expect:\n%v\nactual:\n%v", expect, actual)
	}
}

type testUnmarshalMapKeyComplex struct {
	Foo string
}

func (t *testUnmarshalMapKeyComplex) UnmarshalText(b []byte) error {
	t.Foo = string(b)
	return nil
}
func (t *testUnmarshalMapKeyComplex) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	avM, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("unexpected AttributeValue type %T, %v", av, av)
	}
	avFoo, ok := avM.Value["foo"]
	if !ok {
		return nil
	}

	avS, ok := avFoo.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("unexpected Foo AttributeValue type, %T, %v", avM, avM)
	}

	t.Foo = avS.Value

	return nil
}

func TestUnmarshalTime_S_SS(t *testing.T) {
	type A struct {
		TimeField   time.Time
		TimeFields  []time.Time
		TimeFieldsL []time.Time
	}
	cases := map[string]struct {
		input       string
		expect      time.Time
		decodeTimeS func(string) (time.Time, error)
	}{
		"String RFC3339Nano (Default)": {
			input:  "1970-01-01T00:02:03.01Z",
			expect: time.Unix(123, 10000000).UTC(),
		},
		"String UnixDate": {
			input:  "Thu Jan  1 00:02:03 UTC 1970",
			expect: time.Unix(123, 0).UTC(),
			decodeTimeS: func(v string) (time.Time, error) {
				t, err := time.Parse(time.UnixDate, v)
				if err != nil {
					return time.Time{}, &UnmarshalError{Err: err, Value: v, Type: timeType}
				}
				return t, nil
			},
		},
		"String RFC3339 millis keeping zeroes": {
			input:  "1970-01-01T00:02:03.010Z",
			expect: time.Unix(123, 10000000).UTC(),
			decodeTimeS: func(v string) (time.Time, error) {
				t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", v)
				if err != nil {
					return time.Time{}, &UnmarshalError{Err: err, Value: v, Type: timeType}
				}
				return t, nil
			},
		},
		"String RFC822": {
			input:  "01 Jan 70 00:02 UTC",
			expect: time.Unix(120, 0).UTC(),
			decodeTimeS: func(v string) (time.Time, error) {
				t, err := time.Parse(time.RFC822, v)
				if err != nil {
					return time.Time{}, &UnmarshalError{Err: err, Value: v, Type: timeType}
				}
				return t, nil
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			inputMap := &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"TimeField":  &types.AttributeValueMemberS{Value: c.input},
					"TimeFields": &types.AttributeValueMemberSS{Value: []string{c.input}},
					"TimeFieldsL": &types.AttributeValueMemberL{Value: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: c.input},
					}},
				},
			}
			expectedValue := A{
				TimeField:   c.expect,
				TimeFields:  []time.Time{c.expect},
				TimeFieldsL: []time.Time{c.expect},
			}

			var actualValue A
			if err := UnmarshalWithOptions(inputMap, &actualValue, func(options *DecoderOptions) {
				if c.decodeTimeS != nil {
					options.DecodeTime.S = c.decodeTimeS
				}
			}); err != nil {
				t.Errorf("expect no error, got %v", err)
			}
			if diff := cmp.Diff(expectedValue, actualValue, getIgnoreAVUnexportedOptions()...); diff != "" {
				t.Errorf("expect attribute value match\n%s", diff)
			}
		})
	}
}

func TestUnmarshalTime_N_NS(t *testing.T) {
	type A struct {
		TimeField   time.Time
		TimeFields  []time.Time
		TimeFieldsL []time.Time
	}
	cases := map[string]struct {
		input       string
		expect      time.Time
		decodeTimeN func(string) (time.Time, error)
	}{
		"Number Unix seconds (Default)": {
			input:  "123",
			expect: time.Unix(123, 0),
		},
		"Number Unix milli": {
			input:  "123010",
			expect: time.Unix(123, 10000000),
			decodeTimeN: func(v string) (time.Time, error) {
				n, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return time.Time{}, &UnmarshalError{
						Err: err, Value: v, Type: timeType,
					}
				}
				return time.Unix(0, n*int64(time.Millisecond)), nil
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			inputMap := &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"TimeField":  &types.AttributeValueMemberN{Value: c.input},
					"TimeFields": &types.AttributeValueMemberNS{Value: []string{c.input}},
					"TimeFieldsL": &types.AttributeValueMemberL{Value: []types.AttributeValue{
						&types.AttributeValueMemberN{Value: c.input},
					}},
				},
			}
			expectedValue := A{
				TimeField:   c.expect,
				TimeFields:  []time.Time{c.expect},
				TimeFieldsL: []time.Time{c.expect},
			}

			var actualValue A
			if err := UnmarshalWithOptions(inputMap, &actualValue, func(options *DecoderOptions) {
				if c.decodeTimeN != nil {
					options.DecodeTime.N = c.decodeTimeN
				}
			}); err != nil {
				t.Errorf("expect no error, got %v", err)
			}
			if diff := cmp.Diff(expectedValue, actualValue, getIgnoreAVUnexportedOptions()...); diff != "" {
				t.Errorf("expect attribute value match\n%s", diff)
			}
		})
	}
}

func TestCustomDecodeSAndDefaultDecodeN(t *testing.T) {
	type A struct {
		TimeFieldS time.Time
		TimeFieldN time.Time
	}
	inputMap := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"TimeFieldS": &types.AttributeValueMemberS{Value: "01 Jan 70 00:02 UTC"},
			"TimeFieldN": &types.AttributeValueMemberN{Value: "123"},
		},
	}
	expectedValue := A{
		TimeFieldS: time.Unix(120, 0).UTC(),
		TimeFieldN: time.Unix(123, 0).UTC(),
	}

	var actualValue A
	if err := UnmarshalWithOptions(inputMap, &actualValue, func(options *DecoderOptions) {
		// overriding only the S time decoder will keep the default N time decoder
		options.DecodeTime.S = func(v string) (time.Time, error) {
			t, err := time.Parse(time.RFC822, v)
			if err != nil {
				return time.Time{}, &UnmarshalError{Err: err, Value: v, Type: timeType}
			}
			return t, nil
		}
	}); err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if diff := cmp.Diff(expectedValue, actualValue, getIgnoreAVUnexportedOptions()...); diff != "" {
		t.Errorf("expect attribute value match\n%s", diff)
	}
}

func TestCustomDecodeNAndDefaultDecodeS(t *testing.T) {
	type A struct {
		TimeFieldS time.Time
		TimeFieldN time.Time
	}
	inputMap := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"TimeFieldS": &types.AttributeValueMemberS{Value: "1970-01-01T00:02:03.01Z"},
			"TimeFieldN": &types.AttributeValueMemberN{Value: "123010"},
		},
	}
	expectedValue := A{
		TimeFieldS: time.Unix(123, 10000000).UTC(),
		TimeFieldN: time.Unix(123, 10000000).UTC(),
	}

	var actualValue A
	if err := UnmarshalWithOptions(inputMap, &actualValue, func(options *DecoderOptions) {
		// overriding only the N time decoder will keep the default S time decoder
		options.DecodeTime.N = func(v string) (time.Time, error) {
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return time.Time{}, &UnmarshalError{
					Err: err, Value: v, Type: timeType,
				}
			}
			return time.Unix(0, n*int64(time.Millisecond)), nil
		}
	}); err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if diff := cmp.Diff(expectedValue, actualValue, getIgnoreAVUnexportedOptions()...); diff != "" {
		t.Errorf("expect attribute value match\n%s", diff)
	}
}

func TestUnmarshalMap_keyTypes(t *testing.T) {
	type StrAlias string
	type IntAlias int
	type BoolAlias bool

	cases := map[string]struct {
		input      map[string]types.AttributeValue
		expectVal  interface{}
		expectType func() interface{}
	}{
		"string key": {
			input: map[string]types.AttributeValue{
				"a": &types.AttributeValueMemberN{Value: "123"},
				"b": &types.AttributeValueMemberS{Value: "efg"},
			},
			expectType: func() interface{} { return map[string]interface{}{} },
			expectVal: map[string]interface{}{
				"a": 123.,
				"b": "efg",
			},
		},
		"string alias key": {
			input: map[string]types.AttributeValue{
				"a": &types.AttributeValueMemberN{Value: "123"},
				"b": &types.AttributeValueMemberS{Value: "efg"},
			},
			expectType: func() interface{} { return map[StrAlias]interface{}{} },
			expectVal: map[StrAlias]interface{}{
				"a": 123.,
				"b": "efg",
			},
		},
		"Number key": {
			input: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberN{Value: "123"},
				"2": &types.AttributeValueMemberS{Value: "efg"},
			},
			expectType: func() interface{} { return map[Number]interface{}{} },
			expectVal: map[Number]interface{}{
				Number("1"): 123.,
				Number("2"): "efg",
			},
		},
		"int key": {
			input: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberN{Value: "123"},
				"2": &types.AttributeValueMemberS{Value: "efg"},
			},
			expectType: func() interface{} { return map[int]interface{}{} },
			expectVal: map[int]interface{}{
				1: 123.,
				2: "efg",
			},
		},
		"int alias key": {
			input: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberN{Value: "123"},
				"2": &types.AttributeValueMemberS{Value: "efg"},
			},
			expectType: func() interface{} { return map[IntAlias]interface{}{} },
			expectVal: map[IntAlias]interface{}{
				1: 123.,
				2: "efg",
			},
		},
		"bool key": {
			input: map[string]types.AttributeValue{
				"true":  &types.AttributeValueMemberN{Value: "123"},
				"false": &types.AttributeValueMemberS{Value: "efg"},
			},
			expectType: func() interface{} { return map[bool]interface{}{} },
			expectVal: map[bool]interface{}{
				true:  123.,
				false: "efg",
			},
		},
		"bool alias key": {
			input: map[string]types.AttributeValue{
				"true":  &types.AttributeValueMemberN{Value: "123"},
				"false": &types.AttributeValueMemberS{Value: "efg"},
			},
			expectType: func() interface{} { return map[BoolAlias]interface{}{} },
			expectVal: map[BoolAlias]interface{}{
				true:  123.,
				false: "efg",
			},
		},
		"textMarshaler key": {
			input: map[string]types.AttributeValue{
				"Foo:1": &types.AttributeValueMemberN{Value: "123"},
				"Foo:2": &types.AttributeValueMemberS{Value: "efg"},
			},
			expectType: func() interface{} { return map[testTextMarshaler]interface{}{} },
			expectVal: map[testTextMarshaler]interface{}{
				{Foo: "1"}: 123.,
				{Foo: "2"}: "efg",
			},
		},
		"textMarshaler DDBAvMarshaler key": {
			input: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberN{Value: "123"},
				"2": &types.AttributeValueMemberS{Value: "efg"},
			},
			expectType: func() interface{} { return map[testUnmarshalMapKeyComplex]interface{}{} },
			expectVal: map[testUnmarshalMapKeyComplex]interface{}{
				{Foo: "1"}: 123.,
				{Foo: "2"}: "efg",
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actualVal := c.expectType()
			err := UnmarshalMap(c.input, &actualVal)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			t.Logf("expectType, %T", actualVal)

			if diff := cmp.Diff(c.expectVal, actualVal); diff != "" {
				t.Errorf("expect value match\n%s", diff)
			}
		})
	}
}

func TestUnmarshalMap_keyPtrTypes(t *testing.T) {
	input := map[string]types.AttributeValue{
		"Foo:1": &types.AttributeValueMemberN{Value: "123"},
		"Foo:2": &types.AttributeValueMemberS{Value: "efg"},
	}

	expectVal := map[*testTextMarshaler]interface{}{
		{Foo: "1"}: 123.,
		{Foo: "2"}: "efg",
	}

	actualVal := map[*testTextMarshaler]interface{}{}
	err := UnmarshalMap(input, &actualVal)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	t.Logf("expectType, %T", actualVal)

	if e, a := len(expectVal), len(actualVal); e != a {
		t.Errorf("expect %v values, got %v", e, a)
	}

	for k, v := range expectVal {
		var found bool
		for ak, av := range actualVal {
			if *k == *ak {
				found = true
				if diff := cmp.Diff(v, av); diff != "" {
					t.Errorf("expect value match\n%s", diff)
				}
			}
		}
		if !found {
			t.Errorf("expect %v key not found", *k)
		}
	}

}
