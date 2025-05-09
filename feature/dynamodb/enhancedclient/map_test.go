package enhancedclient

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestMap(t *testing.T) {
	cases := map[string]struct {
		input    Map
		expected map[string]types.AttributeValue
	}{
		"string to map[string]AttributeValueMemberS": {
			input: Map{}.With("k", "k"),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberS{Value: "k"},
			},
		},
		"string slice to map[string]AttributeValueMemberL": {
			input: Map{}.With("k", []string{"k"}),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberS{
							Value: "k",
						},
					},
				},
			},
		},
		"string array to map[string]AttributeValueMemberSS": {
			input: Map{}.With("k", [1]string{"k"}),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberSS{Value: []string{"k"}},
			},
		},
		"int to map[string]AttributeValueMemberN": {
			input: Map{}.With("k", 1),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberN{Value: "1"},
			},
		},
		"float to map[string]AttributeValueMemberN": {
			input: Map{}.With("k", 1.23),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberN{Value: "1.23"},
			},
		},
		"int array to map[string]AttributeValueMemberNS": {
			input: Map{}.With("k", [1]int{1}),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberNS{Value: []string{"1"}},
			},
		},
		"byte slice to map[string]AttributeValueMemberB": {
			input: Map{}.With("k", []byte("k")),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberB{
					Value: []byte("k"),
				},
			},
		},
		"byte array to map[string]AttributeValueMemberBS": {
			input: Map{}.With("k", [][]byte{[]byte("k")}),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberBS{
					Value: [][]byte{[]byte("k")},
				},
			},
		},
		"nil to map[string]AttributeValueMemberNULL": {
			input: Map{}.With("k", nil),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		"map slice to map[string]AttributeValueMemberM": {
			input: Map{}.With("k", map[string]string{"k": "v"}),
			expected: map[string]types.AttributeValue{
				"k": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"k": &types.AttributeValueMemberS{
							Value: "v",
						},
					},
				},
			},
		},
		"bool to map[string]AttributeValueMemberBOOL": {
			input: Map{}.With("true", true).With("false", false),
			expected: map[string]types.AttributeValue{
				"true":  &types.AttributeValueMemberBOOL{Value: true},
				"false": &types.AttributeValueMemberBOOL{Value: false},
			},
		},
	}

	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			if diff := cmpDiff(c.input, Map(c.expected)); len(diff) > 0 {
				t.Fatalf("unexpected diff: %s", diff)
			}
		})
	}
}
