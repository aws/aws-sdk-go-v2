package expression

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestCompositeKey(t *testing.T) {
	cases := []struct {
		input       KeyConditionBuilder
		expected    Expression
		expectError bool
	}{
		{
			input:       CompositeKey(),
			expectError: true,
		},
		{
			input: CompositeKey().
				AddPartitionKey(Key("pk1").Equal(Value("1"))),
			expected: Expression{
				expressionMap: map[expressionType]string{
					keyCondition: "((#0 = :0))",
				},
				namesMap: map[string]string{
					"#0": "pk1",
				},
				valuesMap: map[string]types.AttributeValue{
					":0": &types.AttributeValueMemberS{Value: "1"},
				},
			},
		},
		{
			input: CompositeKey().
				AddPartitionKey(Key("pk1").Equal(Value("1"))).
				AddSortKey(Key("sk1").Equal(Value("1"))),
			expected: Expression{
				expressionMap: map[expressionType]string{
					keyCondition: "((#0 = :0) AND (#1 = :1))",
				},
				namesMap: map[string]string{
					"#0": "pk1",
					"#1": "sk1",
				},
				valuesMap: map[string]types.AttributeValue{
					":0": &types.AttributeValueMemberS{Value: "1"},
					":1": &types.AttributeValueMemberS{Value: "1"},
				},
			},
		},
		{
			input: CompositeKey().
				AddPartitionKey(Key("pk1").Equal(Value("1"))).
				AddPartitionKey(Key("pk2").Equal(Value("1"))).
				AddPartitionKey(Key("pk3").Equal(Value("1"))).
				AddPartitionKey(Key("pk4").Equal(Value("1"))),
			expected: Expression{
				expressionMap: map[expressionType]string{
					keyCondition: "((#0 = :0) AND (#1 = :1) AND (#2 = :2) AND (#3 = :3))",
				},
				namesMap: map[string]string{
					"#0": "pk1",
					"#1": "pk2",
					"#2": "pk3",
					"#3": "pk4",
				},
				valuesMap: map[string]types.AttributeValue{
					":0": &types.AttributeValueMemberS{Value: "1"},
					":1": &types.AttributeValueMemberS{Value: "1"},
					":2": &types.AttributeValueMemberS{Value: "1"},
					":3": &types.AttributeValueMemberS{Value: "1"},
				},
			},
		},
		{
			input: CompositeKey().
				AddPartitionKey(Key("pk1").Equal(Value("1"))).
				AddPartitionKey(Key("pk2").Equal(Value("1"))).
				AddPartitionKey(Key("pk3").Equal(Value("1"))).
				AddPartitionKey(Key("pk4").Equal(Value("1"))).
				AddSortKey(Key("sk1").Equal(Value("1"))).
				AddSortKey(Key("sk2").Equal(Value("1"))).
				AddSortKey(Key("sk3").Equal(Value("1"))).
				AddSortKey(Key("sk4").Equal(Value("1"))),
			expected: Expression{
				expressionMap: map[expressionType]string{
					keyCondition: "((#0 = :0) AND (#1 = :1) AND (#2 = :2) AND (#3 = :3) AND (#4 = :4) AND (#5 = :5) AND (#6 = :6) AND (#7 = :7))",
				},
				namesMap: map[string]string{
					"#0": "pk1",
					"#1": "pk2",
					"#2": "pk3",
					"#3": "pk4",
					"#4": "sk1",
					"#5": "sk2",
					"#6": "sk3",
					"#7": "sk4",
				},
				valuesMap: map[string]types.AttributeValue{
					":0": &types.AttributeValueMemberS{Value: "1"},
					":1": &types.AttributeValueMemberS{Value: "1"},
					":2": &types.AttributeValueMemberS{Value: "1"},
					":3": &types.AttributeValueMemberS{Value: "1"},
					":4": &types.AttributeValueMemberS{Value: "1"},
					":5": &types.AttributeValueMemberS{Value: "1"},
					":6": &types.AttributeValueMemberS{Value: "1"},
					":7": &types.AttributeValueMemberS{Value: "1"},
				},
			},
		},
		{
			input: CompositeKey().
				AddSortKey(Key("sk1").Equal(Value("1"))),
			expectError: true,
		},
		{
			input: CompositeKey().
				AddSortKey(Key("sk1").Equal(Value("1"))).
				AddSortKey(Key("sk2").Equal(Value("1"))).
				AddSortKey(Key("sk3").Equal(Value("1"))).
				AddSortKey(Key("sk4").Equal(Value("1"))),
			expectError: true,
		},
		{
			input: CompositeKey().
				AddPartitionKey(Key("pk1").Equal(Value("1"))).
				And(Key("sk1").Equal(Value("1"))),
			expectError: true,
		},
		{
			input: CompositeKey().
				AddPartitionKey(Key("pk1").GreaterThanEqual(Value("1"))).
				AddPartitionKey(Key("pk1").LessThanEqual(Value("2"))),
			expected: Expression{
				expressionMap: map[expressionType]string{
					keyCondition: "((#0 >= :0) AND (#0 <= :1))",
				},
				namesMap: map[string]string{
					"#0": "pk1",
				},
				valuesMap: map[string]types.AttributeValue{
					":0": &types.AttributeValueMemberS{Value: "1"},
					":1": &types.AttributeValueMemberS{Value: "2"},
				},
			},
		},
	}

	for idx, c := range cases {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			if !c.input.IsSet() {
				t.Errorf("IsSet() is false")
				t.Fail()
			}

			actual, err := NewBuilder().WithKeyCondition(c.input).Build()

			if c.expectError && err == nil {
				t.Error("expected error")
				t.Fail()
			}
			if !c.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
				t.Fail()
			}

			if c.expectError && err != nil {
				t.Logf("found expected error: %v", err)
				return
			}

			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("unexpted diff: %v vs %v", c.expected, actual)
			}
		})
	}
}
