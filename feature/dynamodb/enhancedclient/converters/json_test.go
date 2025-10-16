package converters

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestJsonConverter_FromAttributeValue(t *testing.T) {
	cases := []struct {
		input          types.AttributeValue
		expectedOutput any
		expectedError  bool
	}{
		{
			input:          &types.AttributeValueMemberNULL{Value: true},
			expectedOutput: nil,
		},
		{
			input:          &types.AttributeValueMemberNULL{Value: false},
			expectedOutput: nil,
		},
		{
			input:          &types.AttributeValueMemberS{Value: "[]"},
			expectedOutput: []any{},
		},
		{
			input:          &types.AttributeValueMemberS{Value: "{}"},
			expectedOutput: map[string]any{},
		},
		{
			input: &types.AttributeValueMemberS{Value: `{"test":"test"}`},
			expectedOutput: map[string]any{
				"test": "test",
			},
		},
		{
			input: &types.AttributeValueMemberS{Value: `[{"test":"test"}]`},
			expectedOutput: []any{
				map[string]any{
					"test": "test",
				},
			},
		},
		{
			input:          &types.AttributeValueMemberS{Value: `"test"`},
			expectedOutput: "test",
		},
		{
			input:         &types.AttributeValueMemberS{Value: `[`},
			expectedError: true,
		},
		{
			input:         &types.AttributeValueMemberS{Value: `[{"test":"test}"`},
			expectedError: true,
		},
		{
			input:         &types.AttributeValueMemberS{Value: ""},
			expectedError: true,
		},
		{
			input:         &types.AttributeValueMemberB{Value: []byte{}},
			expectedError: true,
		},
		{
			input:         &types.AttributeValueMemberM{},
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			j := JsonConverter{}
			actualOutput, actualError := j.FromAttributeValue(c.input, nil)

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

func TestJsonConverter_ToAttributeValue(t *testing.T) {
	cases := []struct {
		input          any
		expectedOutput types.AttributeValue
		expectedError  bool
	}{
		{
			expectedOutput: &types.AttributeValueMemberNULL{Value: true},
		},
		{
			input:          "test",
			expectedOutput: &types.AttributeValueMemberS{Value: `"test"`},
		},
		{
			input:          []any{},
			expectedOutput: &types.AttributeValueMemberS{Value: `[]`},
		},
		{
			input:          map[string]any{},
			expectedOutput: &types.AttributeValueMemberS{Value: `{}`},
		},
		{
			input:          map[string]any{"test": "test"},
			expectedOutput: &types.AttributeValueMemberS{Value: `{"test":"test"}`},
		},
		{
			input:          []string{"test"},
			expectedOutput: &types.AttributeValueMemberS{Value: `["test"]`},
		},
		{
			input: struct {
				Test   string `json:"test"`
				hidden string `json:"hidden"`
			}{
				Test:   "test",
				hidden: "you-can't-see-me",
			},
			expectedOutput: &types.AttributeValueMemberS{Value: `{"test":"test"}`},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			j := JsonConverter{}
			actualOutput, actualError := j.ToAttributeValue(c.input, nil)

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
