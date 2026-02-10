package converters

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestByteArrayConverter_FromAttributeValue(t *testing.T) {
	cases := []struct {
		input          types.AttributeValue
		opts           []string
		expectedOutput any
		expectedError  bool
	}{
		{input: &types.AttributeValueMemberB{Value: nil}, opts: []string{}, expectedOutput: ([]byte)(nil), expectedError: false},
		{input: &types.AttributeValueMemberB{Value: []byte("test")}, opts: []string{}, expectedOutput: []byte("test"), expectedError: false},
		// errors
		{input: nil, opts: nil, expectedOutput: nil, expectedError: true},
		{input: (types.AttributeValue)(nil), opts: nil, expectedOutput: nil, expectedError: true},
		{input: (*types.AttributeValueMemberN)(nil), opts: nil, expectedOutput: nil, expectedError: true},
		{input: (*types.AttributeValueMemberS)(nil), opts: nil, expectedOutput: nil, expectedError: true},
		{input: (*types.AttributeValueMemberBOOL)(nil), opts: nil, expectedOutput: nil, expectedError: true},
		{input: &types.AttributeValueMemberS{Value: "true"}, opts: nil, expectedOutput: nil, expectedError: true},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc := ByteArrayConverter{}

			actualOutput, actualError := tc.FromAttributeValue(c.input, c.opts)

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

func TestByteArrayConverter_ToAttributeValue(t *testing.T) {
	cases := []struct {
		input          []uint8
		opts           []string
		expectedOutput types.AttributeValue
		expectedError  bool
	}{
		{input: nil, opts: nil, expectedOutput: nil, expectedError: true},
		{input: ([]byte)(nil), opts: nil, expectedOutput: nil, expectedError: true},
		{input: []byte("tests"), opts: nil, expectedOutput: &types.AttributeValueMemberB{Value: []byte("tests")}, expectedError: false},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc := ByteArrayConverter{}

			actualOutput, actualError := tc.ToAttributeValue(c.input, c.opts)

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
