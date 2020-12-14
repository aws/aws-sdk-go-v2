package attributevalue

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
	"github.com/google/go-cmp/cmp"
)

type testBinarySetStruct struct {
	Binarys [][]byte `dynamodbav:",binaryset"`
}
type testNumberSetStruct struct {
	Numbers []int `dynamodbav:",numberset"`
}
type testStringSetStruct struct {
	Strings []string `dynamodbav:",stringset"`
}

type testIntAsStringStruct struct {
	Value int `dynamodbav:",string"`
}

type testOmitEmptyStruct struct {
	Value  string  `dynamodbav:",omitempty"`
	Value2 *string `dynamodbav:",omitempty"`
	Value3 int
}

type testAliasedString string
type testAliasedStringSlice []string
type testAliasedInt int
type testAliasedIntSlice []int
type testAliasedMap map[string]int
type testAliasedSlice []string
type testAliasedByteSlice []byte
type testAliasedBool bool
type testAliasedBoolSlice []bool

type testAliasedStruct struct {
	Value  testAliasedString
	Value2 testAliasedInt
	Value3 testAliasedMap
	Value4 testAliasedSlice

	Value5 testAliasedByteSlice
	Value6 []testAliasedInt
	Value7 []testAliasedString

	Value8  []testAliasedByteSlice `dynamodbav:",binaryset"`
	Value9  []testAliasedInt       `dynamodbav:",numberset"`
	Value10 []testAliasedString    `dynamodbav:",stringset"`

	Value11 testAliasedIntSlice
	Value12 testAliasedStringSlice

	Value13 testAliasedBool
	Value14 testAliasedBoolSlice

	Value15 map[testAliasedString]string
}

type testNamedPointer *int

var testDate, _ = time.Parse(time.RFC3339, "2016-05-03T17:06:26.209072Z")

var sharedTestCases = map[string]struct {
	in               types.AttributeValue
	actual, expected interface{}
	err              error
}{
	"binary slice": {
		in:       &types.AttributeValueMemberB{Value: []byte{48, 49}},
		actual:   &[]byte{},
		expected: []byte{48, 49},
	},
	"Binary slice oversized": {
		in: &types.AttributeValueMemberB{Value: []byte{48, 49}},
		actual: func() *[]byte {
			v := make([]byte, 0, 10)
			return &v
		}(),
		expected: []byte{48, 49},
	},
	"binary slice pointer": {
		in: &types.AttributeValueMemberB{Value: []byte{48, 49}},
		actual: func() **[]byte {
			v := make([]byte, 0, 10)
			v2 := &v
			return &v2
		}(),
		expected: []byte{48, 49},
	},
	"bool": {
		in:       &types.AttributeValueMemberBOOL{Value: true},
		actual:   new(bool),
		expected: true,
	},
	"list": {
		in: &types.AttributeValueMemberL{Value: []types.AttributeValue{
			&types.AttributeValueMemberN{Value: "123"},
		}},
		actual:   &[]int{},
		expected: []int{123},
	},
	"map, interface": {
		in: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"abc": &types.AttributeValueMemberN{Value: "123"},
		}},
		actual:   &map[string]int{},
		expected: map[string]int{"abc": 123},
	},
	"map, struct": {
		in: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"Abc": &types.AttributeValueMemberN{Value: "123"},
		}},
		actual:   &struct{ Abc int }{},
		expected: struct{ Abc int }{Abc: 123},
	},
	"map, struct with tags": {
		in: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"abc": &types.AttributeValueMemberN{Value: "123"},
		}},
		actual: &struct {
			Abc int `json:"abc" dynamodbav:"abc"`
		}{},
		expected: struct {
			Abc int `json:"abc" dynamodbav:"abc"`
		}{Abc: 123},
	},
	"number, int": {
		in:       &types.AttributeValueMemberN{Value: "123"},
		actual:   new(int),
		expected: 123,
	},
	"number, Float": {
		in:       &types.AttributeValueMemberN{Value: "123.1"},
		actual:   new(float64),
		expected: float64(123.1),
	},
	"null ptr": {
		in:       &types.AttributeValueMemberNULL{Value: true},
		actual:   new(*string),
		expected: nil,
	},
	"string": {
		in:       &types.AttributeValueMemberS{Value: "abc"},
		actual:   new(string),
		expected: "abc",
	},
	"empty string": {
		in:       &types.AttributeValueMemberS{Value: ""},
		actual:   new(string),
		expected: "",
	},
	"binary Set": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Binarys": &types.AttributeValueMemberBS{Value: [][]byte{{48, 49}, {50, 51}}},
			},
		},
		actual:   &testBinarySetStruct{},
		expected: testBinarySetStruct{Binarys: [][]byte{{48, 49}, {50, 51}}},
	},
	"number Set": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Numbers": &types.AttributeValueMemberNS{Value: []string{"123", "321"}},
			},
		},
		actual:   &testNumberSetStruct{},
		expected: testNumberSetStruct{Numbers: []int{123, 321}},
	},
	"string Set": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Strings": &types.AttributeValueMemberSS{Value: []string{"abc", "efg"}},
			},
		},
		actual:   &testStringSetStruct{},
		expected: testStringSetStruct{Strings: []string{"abc", "efg"}},
	},
	"int value as string": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Value": &types.AttributeValueMemberS{Value: "123"},
			},
		},
		actual:   &testIntAsStringStruct{},
		expected: testIntAsStringStruct{Value: 123},
	},
	"omitempty": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Value3": &types.AttributeValueMemberN{Value: "0"},
			},
		},
		actual:   &testOmitEmptyStruct{},
		expected: testOmitEmptyStruct{Value: "", Value2: nil, Value3: 0},
	},
	"aliased type": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Value":  &types.AttributeValueMemberS{Value: "123"},
				"Value2": &types.AttributeValueMemberN{Value: "123"},
				"Value3": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"Key": &types.AttributeValueMemberN{Value: "321"},
				}},
				"Value4": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "1"},
					&types.AttributeValueMemberS{Value: "2"},
					&types.AttributeValueMemberS{Value: "3"},
				}},
				"Value5": &types.AttributeValueMemberB{Value: []byte{0, 1, 2}},
				"Value6": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberN{Value: "1"},
					&types.AttributeValueMemberN{Value: "2"},
					&types.AttributeValueMemberN{Value: "3"},
				}},
				"Value7": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "1"},
					&types.AttributeValueMemberS{Value: "2"},
					&types.AttributeValueMemberS{Value: "3"},
				}},
				"Value8": &types.AttributeValueMemberBS{Value: [][]byte{
					{0, 1, 2}, {3, 4, 5},
				}},
				"Value9": &types.AttributeValueMemberNS{Value: []string{
					"1",
					"2",
					"3",
				}},
				"Value10": &types.AttributeValueMemberSS{Value: []string{
					"1",
					"2",
					"3",
				}},
				"Value11": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberN{Value: "1"},
					&types.AttributeValueMemberN{Value: "2"},
					&types.AttributeValueMemberN{Value: "3"},
				}},
				"Value12": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "1"},
					&types.AttributeValueMemberS{Value: "2"},
					&types.AttributeValueMemberS{Value: "3"},
				}},
				"Value13": &types.AttributeValueMemberBOOL{Value: true},
				"Value14": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberBOOL{Value: true},
					&types.AttributeValueMemberBOOL{Value: false},
					&types.AttributeValueMemberBOOL{Value: true},
				}},
				"Value15": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"TestKey": &types.AttributeValueMemberS{Value: "TestElement"},
				}},
			},
		},
		actual: &testAliasedStruct{},
		expected: testAliasedStruct{
			Value: "123", Value2: 123,
			Value3: testAliasedMap{
				"Key": 321,
			},
			Value4: testAliasedSlice{"1", "2", "3"},
			Value5: testAliasedByteSlice{0, 1, 2},
			Value6: []testAliasedInt{1, 2, 3},
			Value7: []testAliasedString{"1", "2", "3"},
			Value8: []testAliasedByteSlice{
				{0, 1, 2},
				{3, 4, 5},
			},
			Value9:  []testAliasedInt{1, 2, 3},
			Value10: []testAliasedString{"1", "2", "3"},
			Value11: testAliasedIntSlice{1, 2, 3},
			Value12: testAliasedStringSlice{"1", "2", "3"},
			Value13: true,
			Value14: testAliasedBoolSlice{true, false, true},
			Value15: map[testAliasedString]string{"TestKey": "TestElement"},
		},
	},
	"number named pointer": {
		in:       &types.AttributeValueMemberN{Value: "123"},
		actual:   new(testNamedPointer),
		expected: testNamedPointer(aws.Int(123)),
	},
	"time.Time": {
		in:       &types.AttributeValueMemberS{Value: "2016-05-03T17:06:26.209072Z"},
		actual:   new(time.Time),
		expected: testDate,
	},
	"time.Time List": {
		in: &types.AttributeValueMemberL{Value: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "2016-05-03T17:06:26.209072Z"},
			&types.AttributeValueMemberS{Value: "2016-05-04T17:06:26.209072Z"},
		}},
		actual:   new([]time.Time),
		expected: []time.Time{testDate, testDate.Add(24 * time.Hour)},
	},
	"time.Time struct": {
		in: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"abc": &types.AttributeValueMemberS{Value: "2016-05-03T17:06:26.209072Z"},
		}},
		actual: &struct {
			Abc time.Time `json:"abc" dynamodbav:"abc"`
		}{},
		expected: struct {
			Abc time.Time `json:"abc" dynamodbav:"abc"`
		}{Abc: testDate},
	},
	"time.Time ptr struct": {
		in: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"abc": &types.AttributeValueMemberS{Value: "2016-05-03T17:06:26.209072Z"},
		}},
		actual: &struct {
			Abc *time.Time `json:"abc" dynamodbav:"abc"`
		}{},
		expected: struct {
			Abc *time.Time `json:"abc" dynamodbav:"abc"`
		}{Abc: &testDate},
	},
}

var sharedListTestCases = map[string]struct {
	in               []types.AttributeValue
	actual, expected interface{}
	err              error
}{
	"union members": {
		in: []types.AttributeValue{
			&types.AttributeValueMemberB{Value: []byte{48, 49}},
			&types.AttributeValueMemberBOOL{Value: true},
			&types.AttributeValueMemberN{Value: "123"},
			&types.AttributeValueMemberS{Value: "123"},
		},
		actual: func() *[]interface{} {
			v := []interface{}{}
			return &v
		}(),
		expected: []interface{}{[]byte{48, 49}, true, 123., "123"},
	},
	"numbers": {
		in: []types.AttributeValue{
			&types.AttributeValueMemberN{Value: "1"},
			&types.AttributeValueMemberN{Value: "2"},
			&types.AttributeValueMemberN{Value: "3"},
		},
		actual:   &[]interface{}{},
		expected: []interface{}{1., 2., 3.},
	},
}

var sharedMapTestCases = map[string]struct {
	in               map[string]types.AttributeValue
	actual, expected interface{}
	err              error
}{
	"union members": {
		in: map[string]types.AttributeValue{
			"B":    &types.AttributeValueMemberB{Value: []byte{48, 49}},
			"BOOL": &types.AttributeValueMemberBOOL{Value: true},
			"N":    &types.AttributeValueMemberN{Value: "123"},
			"S":    &types.AttributeValueMemberS{Value: "123"},
		},
		actual: &map[string]interface{}{},
		expected: map[string]interface{}{
			"B": []byte{48, 49}, "BOOL": true,
			"N": 123., "S": "123",
		},
	},
}

func assertConvertTest(t *testing.T, actual, expected interface{}, err, expectedErr error) {
	t.Helper()

	if expectedErr != nil {
		if err != nil {
			if e, a := expectedErr, err; !strings.Contains(a.Error(), e.Error()) {
				t.Errorf("expect %v, got %v", e, a)
			}
		} else {
			t.Fatalf("expected error, %v", expectedErr)
		}
	} else if err != nil {
		t.Fatalf("expect no error, got %v", err)
	} else {
		if diff := cmp.Diff(ptrToValue(expected), ptrToValue(actual)); len(diff) != 0 {
			t.Errorf("expect match\n%s", diff)
		}
	}
}

func ptrToValue(in interface{}) interface{} {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		return ptrToValue(v.Interface())
	}
	return v.Interface()
}
