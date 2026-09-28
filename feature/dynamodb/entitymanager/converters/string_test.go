package converters

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestStringConverter_FromAttributeValue(t *testing.T) {
	cases := []struct {
		input          types.AttributeValue
		options        []string
		expectedOutput any
		expectedError  bool
	}{
		{
			input:          &types.AttributeValueMemberS{},
			expectedOutput: "",
		},
		{
			input: &types.AttributeValueMemberS{
				Value: "test",
			},
			expectedOutput: "test",
		},
		{
			input:         &types.AttributeValueMemberB{},
			expectedError: true,
		},
		{
			input:         (*types.AttributeValueMemberB)(nil),
			expectedError: true,
		},
		{
			input:         (*types.AttributeValueMemberS)(nil),
			expectedError: true,
		},
		{
			input:         nil,
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cvt := &StringConverter{}
			actualOutput, actualError := cvt.FromAttributeValue(c.input, c.options)

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

func TestStringConverter_ToAttributeValue(t *testing.T) {
	cases := []struct {
		input          string
		options        []string
		expectedOutput types.AttributeValue
		expectedError  bool
	}{
		{
			input:          "",
			expectedOutput: &types.AttributeValueMemberS{},
		},
		{
			input: "test",
			expectedOutput: &types.AttributeValueMemberS{
				Value: "test",
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cvt := &StringConverter{}
			actualOutput, actualError := cvt.ToAttributeValue(c.input, c.options)

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
