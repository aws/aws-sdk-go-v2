package converters

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestNumericConverter_FromAttributeValue(t *testing.T) {
	cases := []struct {
		input          types.AttributeValue
		options        []string
		expectedOutput any
		expectedError  bool
	}{
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: uint(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "-123"},
			expectedOutput: uint(18446744073709551493),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: uint(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: uint8(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "1234"},
			expectedOutput: uint8(210),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "-1234"},
			expectedOutput: uint8(46),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: uint8(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: uint16(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123456"},
			expectedOutput: uint16(57920),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: uint16(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: uint32(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "12345678901"},
			expectedOutput: uint32(3755744309),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: uint32(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: uint64(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123456789012345678901234567890123"},
			expectedOutput: uint64(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: uint64(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123.10"},
			expectedOutput: float32(123.10),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "-123.10"},
			expectedOutput: float32(-123.10),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123456789012345678901234567890123"},
			expectedOutput: float32(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123,10"},
			expectedOutput: float32(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: float32(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123.10"},
			expectedOutput: float64(123.10),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "-123.10"},
			expectedOutput: float64(-123.10),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123456789012345678901234567890123"},
			expectedOutput: float64(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123,10"},
			expectedOutput: float64(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: float64(0),
			expectedError:  true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var actualError error
			var actualOutput any

			switch v := c.expectedOutput.(type) {
			case uint:
				cvt := NumericConverter[uint]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case uint8:
				cvt := NumericConverter[uint8]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case uint16:
				cvt := NumericConverter[uint16]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case uint32:
				cvt := NumericConverter[uint32]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case uint64:
				cvt := NumericConverter[uint64]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case int:
				cvt := NumericConverter[int]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case int8:
				cvt := NumericConverter[int8]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case int16:
				cvt := NumericConverter[int16]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case int32:
				cvt := NumericConverter[int32]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case int64:
				cvt := NumericConverter[int64]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case float32:
				cvt := NumericConverter[float32]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case float64:
				cvt := NumericConverter[float64]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			default:
				t.Errorf("unsupported type: %T", v)
			}

			if actualError == nil && c.expectedError {
				t.Fatalf("expected error, got none")
			}

			if actualError != nil && !c.expectedError {
				t.Fatalf("unexpected error, got: %v", actualError)
			}

			if actualError != nil && c.expectedError {
				return
			}

			if !reflect.DeepEqual(c.expectedOutput, actualOutput) {
				t.Fatalf("%#+v != %#+v", c.expectedOutput, actualOutput)
			}
		})
	}
}

func TestNumericConverter_ToAttributeValue(t *testing.T) {
	cases := []struct {
		input          any
		options        []string
		expectedOutput types.AttributeValue
		expectedError  bool
	}{
		{
			input:          uint(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          uint8(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          uint16(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          uint32(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          uint64(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          int(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          int8(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          int16(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          int32(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          int64(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          float32(123.456),
			expectedOutput: &types.AttributeValueMemberN{Value: "123.456"},
			expectedError:  false,
		},
		{
			input:          float64(123.456),
			expectedOutput: &types.AttributeValueMemberN{Value: "123.456"},
			expectedError:  false,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var actualError error
			var actualOutput any

			switch v := c.input.(type) {
			case uint:
				cvt := NumericConverter[uint]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case uint8:
				cvt := NumericConverter[uint8]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case uint16:
				cvt := NumericConverter[uint16]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case uint32:
				cvt := NumericConverter[uint32]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case uint64:
				cvt := NumericConverter[uint64]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case int:
				cvt := NumericConverter[int]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case int8:
				cvt := NumericConverter[int8]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case int16:
				cvt := NumericConverter[int16]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case int32:
				cvt := NumericConverter[int32]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case int64:
				cvt := NumericConverter[int64]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case float32:
				cvt := NumericConverter[float32]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case float64:
				cvt := NumericConverter[float64]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			default:
				t.Errorf("unsupported type: %T", v)
			}

			if actualError == nil && c.expectedError {
				t.Fatalf("expected error, got none")
			}

			if actualError != nil && !c.expectedError {
				t.Fatalf("unexpected error, got: %v", actualError)
			}

			if actualError != nil && c.expectedError {
				return
			}

			if !reflect.DeepEqual(c.expectedOutput, actualOutput) {
				t.Fatalf("%#+v != %#+v", c.expectedOutput, actualOutput)
			}
		})
	}
}
