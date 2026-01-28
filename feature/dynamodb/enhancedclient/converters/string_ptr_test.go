package converters

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestStringPtrConverter_FromAttributeValue(t *testing.T) {
	cases := []struct {
		input          types.AttributeValue
		options        []string
		expectedOutput any
		expectedError  bool
	}{
		{
			input:          &types.AttributeValueMemberS{},
			expectedOutput: aws.String(""),
		},
		{
			input: &types.AttributeValueMemberS{
				Value: "test",
			},
			expectedOutput: aws.String("test"),
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
			cvt := &StringPtrConverter{}
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

func TestStringPtrConverter_ToAttributeValue(t *testing.T) {
	cases := []struct {
		input          *string
		options        []string
		expectedOutput types.AttributeValue
		expectedError  bool
	}{
		{
			input:          aws.String(""),
			expectedOutput: &types.AttributeValueMemberS{},
		},
		{
			input: aws.String("test"),
			expectedOutput: &types.AttributeValueMemberS{
				Value: "test",
			},
		},
		{
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cvt := &StringPtrConverter{}
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
