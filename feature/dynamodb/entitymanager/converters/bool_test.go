package converters

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestBoolConverter_FromAttributeValue(t *testing.T) {
	cases := []struct {
		input          types.AttributeValue
		opts           []string
		expectedOutput any
		expectedError  bool
	}{
		{input: &types.AttributeValueMemberBOOL{Value: true}, opts: []string{}, expectedOutput: true, expectedError: false},
		{input: &types.AttributeValueMemberBOOL{Value: false}, opts: []string{}, expectedOutput: false, expectedError: false},
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
			tc := BoolConverter{}

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

func TestBoolConverter_ToAttributeValue(t *testing.T) {
	cases := []struct {
		input          bool
		opts           []string
		expectedOutput types.AttributeValue
		expectedError  bool
	}{
		{input: true, opts: nil, expectedOutput: &types.AttributeValueMemberBOOL{Value: true}, expectedError: false},
		{input: false, opts: nil, expectedOutput: &types.AttributeValueMemberBOOL{Value: false}, expectedError: false},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc := BoolConverter{}

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
