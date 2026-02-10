package converters

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestNumericPtrConverter_FromAttributeValue(t *testing.T) {
	cases := []struct {
		input          types.AttributeValue
		options        []string
		expectedOutput any
		expectedError  bool
	}{
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: aws.Uint(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "-123"},
			expectedOutput: aws.Uint(18446744073709551493),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: aws.Uint(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: aws.Uint8(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "1234"},
			expectedOutput: aws.Uint8(210),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "-1234"},
			expectedOutput: aws.Uint8(46),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: aws.Uint8(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: aws.Uint16(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123456"},
			expectedOutput: aws.Uint16(57920),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: aws.Uint16(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: aws.Uint32(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "12345678901"},
			expectedOutput: aws.Uint32(3755744309),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: aws.Uint32(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123"},
			expectedOutput: aws.Uint64(123),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123456789012345678901234567890123"},
			expectedOutput: aws.Uint64(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: aws.Uint64(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123.10"},
			expectedOutput: aws.Float32(123.10),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "-123.10"},
			expectedOutput: aws.Float32(-123.10),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123456789012345678901234567890123"},
			expectedOutput: aws.Float32(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123,10"},
			expectedOutput: aws.Float32(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: aws.Float32(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123.10"},
			expectedOutput: aws.Float64(123.10),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "-123.10"},
			expectedOutput: aws.Float64(-123.10),
			expectedError:  false,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123456789012345678901234567890123"},
			expectedOutput: aws.Float64(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "123,10"},
			expectedOutput: aws.Float64(0),
			expectedError:  true,
		},
		{
			input:          &types.AttributeValueMemberN{Value: "test"},
			expectedOutput: aws.Float64(0),
			expectedError:  true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var actualError error
			var actualOutput any

			switch v := c.expectedOutput.(type) {
			case *uint:
				cvt := NumericPtrConverter[uint]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *uint8:
				cvt := NumericPtrConverter[uint8]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *uint16:
				cvt := NumericPtrConverter[uint16]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *uint32:
				cvt := NumericPtrConverter[uint32]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *uint64:
				cvt := NumericPtrConverter[uint64]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *int:
				cvt := NumericPtrConverter[int]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *int8:
				cvt := NumericPtrConverter[int8]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *int16:
				cvt := NumericPtrConverter[int16]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *int32:
				cvt := NumericPtrConverter[int32]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *int64:
				cvt := NumericPtrConverter[int64]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *float32:
				cvt := NumericPtrConverter[float32]{}
				actualOutput, actualError = cvt.FromAttributeValue(c.input, c.options)
			case *float64:
				cvt := NumericPtrConverter[float64]{}
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

func TestNumericPtrConverter_ToAttributeValue(t *testing.T) {
	cases := []struct {
		input          any
		options        []string
		expectedOutput types.AttributeValue
		expectedError  bool
	}{
		{
			input:          aws.Uint(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Uint8(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Uint16(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Uint32(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Uint64(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Int(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Int8(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Int16(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Int32(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Int64(123),
			expectedOutput: &types.AttributeValueMemberN{Value: "123"},
			expectedError:  false,
		},
		{
			input:          aws.Float32(123.456),
			expectedOutput: &types.AttributeValueMemberN{Value: "123.456"},
			expectedError:  false,
		},
		{
			input:          aws.Float64(123.456),
			expectedOutput: &types.AttributeValueMemberN{Value: "123.456"},
			expectedError:  false,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var actualError error
			var actualOutput any

			switch v := c.input.(type) {
			case *uint:
				cvt := NumericPtrConverter[uint]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *uint8:
				cvt := NumericPtrConverter[uint8]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *uint16:
				cvt := NumericPtrConverter[uint16]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *uint32:
				cvt := NumericPtrConverter[uint32]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *uint64:
				cvt := NumericPtrConverter[uint64]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *int:
				cvt := NumericPtrConverter[int]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *int8:
				cvt := NumericPtrConverter[int8]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *int16:
				cvt := NumericPtrConverter[int16]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *int32:
				cvt := NumericPtrConverter[int32]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *int64:
				cvt := NumericPtrConverter[int64]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *float32:
				cvt := NumericPtrConverter[float32]{}
				actualOutput, actualError = cvt.ToAttributeValue(v, c.options)
			case *float64:
				cvt := NumericPtrConverter[float64]{}
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
