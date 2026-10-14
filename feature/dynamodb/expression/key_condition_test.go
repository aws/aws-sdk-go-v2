package expression

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// keyCondErrorMode will help with error cases and checking error types
type keyCondErrorMode string

const (
	noKeyConditionError keyCondErrorMode = ""
	// unsetKeyCondition error will occur when buildTree() is called on an empty
	// KeyConditionBuilder
	unsetKeyCondition = "unset parameter: KeyConditionBuilder"
	// invalidKeyConditionOperand error will occur when an invalid OperandBuilder is used as
	// an argument
	invalidKeyConditionOperand = "BuildOperand error"
	// invalidKeyConditionFormat error will occur when the first key condition is not an equal
	// clause or if more then one And condition is provided
	invalidKeyConditionFormat = "buildKeyCondition error: invalid key condition constructed"
)

// IsSet
func TestKeyIsSet(t *testing.T) {
	cases := []struct {
		name     string
		input    KeyConditionBuilder
		expected bool
	}{
		{
			name:     "set",
			input:    KeyEqual(Key("foo"), Value("bar")),
			expected: true,
		},
		{
			name:     "unset",
			input:    KeyConditionBuilder{},
			expected: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if actual := c.input.IsSet(); actual != c.expected {
				t.Errorf("expected %t, got %t", c.expected, actual)
			}
		})
	}
}

func TestKeyCompare(t *testing.T) {
	cases := []struct {
		name         string
		input        KeyConditionBuilder
		expectedNode exprNode
		err          keyCondErrorMode
	}{
		{
			name:  "key equal",
			input: Key("foo").Equal(Value(5)),
			expectedNode: exprNode{
				children: []exprNode{
					{
						names:   []string{"foo"},
						fmtExpr: "$n",
					},
					{
						values: []types.AttributeValue{
							&types.AttributeValueMemberN{Value: "5"},
						},
						fmtExpr: "$v",
					},
				},
				fmtExpr: "$c = $c",
			},
		},
		{
			name:  "key less than",
			input: Key("foo").LessThan(Value(5)),
			expectedNode: exprNode{
				children: []exprNode{
					{
						names:   []string{"foo"},
						fmtExpr: "$n",
					},
					{
						values: []types.AttributeValue{
							&types.AttributeValueMemberN{Value: "5"},
						},
						fmtExpr: "$v",
					},
				},
				fmtExpr: "$c < $c",
			},
		},
		{
			name:  "key less than equal",
			input: Key("foo").LessThanEqual(Value(5)),
			expectedNode: exprNode{
				children: []exprNode{
					{
						names:   []string{"foo"},
						fmtExpr: "$n",
					},
					{
						values: []types.AttributeValue{
							&types.AttributeValueMemberN{Value: "5"},
						},
						fmtExpr: "$v",
					},
				},
				fmtExpr: "$c <= $c",
			},
		},
		{
			name:  "key greater than",
			input: Key("foo").GreaterThan(Value(5)),
			expectedNode: exprNode{
				children: []exprNode{
					{
						names:   []string{"foo"},
						fmtExpr: "$n",
					},
					{
						values: []types.AttributeValue{
							&types.AttributeValueMemberN{Value: "5"},
						},
						fmtExpr: "$v",
					},
				},
				fmtExpr: "$c > $c",
			},
		},
		{
			name:  "key greater than equal",
			input: Key("foo").GreaterThanEqual(Value(5)),
			expectedNode: exprNode{
				children: []exprNode{
					{
						names:   []string{"foo"},
						fmtExpr: "$n",
					},
					{
						values: []types.AttributeValue{
							&types.AttributeValueMemberN{Value: "5"},
						},
						fmtExpr: "$v",
					},
				},
				fmtExpr: "$c >= $c",
			},
		},
		{
			name:  "unset KeyConditionBuilder",
			input: KeyConditionBuilder{},
			err:   unsetKeyCondition,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := c.input.buildTree()
			if c.err != noKeyConditionError {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := string(c.err), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expectedNode, actual; !reflect.DeepEqual(a, e) {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}

func TestKeyBetween(t *testing.T) {
	cases := []struct {
		name         string
		input        KeyConditionBuilder
		expectedNode exprNode
		err          keyCondErrorMode
	}{
		{
			name:  "key between",
			input: Key("foo").Between(Value(5), Value(10)),
			expectedNode: exprNode{
				children: []exprNode{
					{
						names:   []string{"foo"},
						fmtExpr: "$n",
					},
					{
						values: []types.AttributeValue{
							&types.AttributeValueMemberN{Value: "5"},
						},
						fmtExpr: "$v",
					},
					{
						values: []types.AttributeValue{
							&types.AttributeValueMemberN{Value: "10"},
						},
						fmtExpr: "$v",
					},
				},
				fmtExpr: "$c BETWEEN $c AND $c",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := c.input.buildTree()
			if c.err != noKeyConditionError {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := string(c.err), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expectedNode, actual; !reflect.DeepEqual(a, e) {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}

func TestKeyBeginsWith(t *testing.T) {
	cases := []struct {
		name         string
		input        KeyConditionBuilder
		expectedNode exprNode
		err          keyCondErrorMode
	}{
		{
			name:  "key begins with",
			input: Key("foo").BeginsWith("bar"),
			expectedNode: exprNode{
				children: []exprNode{
					{
						names:   []string{"foo"},
						fmtExpr: "$n",
					},
					{
						values: []types.AttributeValue{
							&types.AttributeValueMemberS{Value: "bar"},
						},
						fmtExpr: "$v",
					},
				},
				fmtExpr: "begins_with ($c, $c)",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := c.input.buildTree()
			if c.err != noKeyConditionError {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := string(c.err), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expectedNode, actual; !reflect.DeepEqual(a, e) {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}

func TestKeyAnd(t *testing.T) {
	cases := []struct {
		name         string
		input        KeyConditionBuilder
		expectedNode exprNode
		err          keyCondErrorMode
	}{
		{
			name:  "key and",
			input: Key("foo").Equal(Value(5)).And(Key("bar").BeginsWith("baz")),
			expectedNode: exprNode{
				children: []exprNode{
					{
						children: []exprNode{
							{
								names:   []string{"foo"},
								fmtExpr: "$n",
							},
							{
								values: []types.AttributeValue{
									&types.AttributeValueMemberN{Value: "5"},
								},
								fmtExpr: "$v",
							},
						},
						fmtExpr: "$c = $c",
					},
					{
						children: []exprNode{
							{
								names:   []string{"bar"},
								fmtExpr: "$n",
							},
							{
								values: []types.AttributeValue{
									&types.AttributeValueMemberS{Value: "baz"},
								},
								fmtExpr: "$v",
							},
						},
						fmtExpr: "begins_with ($c, $c)",
					},
				},
				fmtExpr: "($c) AND ($c)",
			},
		},
		{
			name:  "multi-attribute key and",
			input: Key("foo").Equal(Value(5)).And(Key("bar").Equal(Value("baz"))).And(Key("qux").BeginsWith("corge")),
			expectedNode: exprNode{
				children: []exprNode{
					{
						children: []exprNode{
							{
								names:   []string{"foo"},
								fmtExpr: "$n",
							},
							{
								values: []types.AttributeValue{
									&types.AttributeValueMemberN{Value: "5"},
								},
								fmtExpr: "$v",
							},
						},
						fmtExpr: "$c = $c",
					},
					{
						children: []exprNode{
							{
								names:   []string{"bar"},
								fmtExpr: "$n",
							},
							{
								values: []types.AttributeValue{
									&types.AttributeValueMemberS{Value: "baz"},
								},
								fmtExpr: "$v",
							},
						},
						fmtExpr: "$c = $c",
					},
					{
						children: []exprNode{
							{
								names:   []string{"qux"},
								fmtExpr: "$n",
							},
							{
								values: []types.AttributeValue{
									&types.AttributeValueMemberS{Value: "corge"},
								},
								fmtExpr: "$v",
							},
						},
						fmtExpr: "begins_with ($c, $c)",
					},
				},
				fmtExpr: "($c) AND ($c) AND ($c)",
			},
		},
		{
			name:  "multi-attribute non terminal condition is not equal",
			input: Key("foo").Equal(Value(5)).And(Key("bar").BeginsWith("baz")).And(Key("qux").Equal(Value("corge"))),
			err:   invalidKeyConditionFormat,
		},
		{
			name:  "first condition is not equal",
			input: Key("foo").LessThan(Value(5)).And(Key("bar").BeginsWith("baz")),
			err:   invalidKeyConditionFormat,
		},
		{
			name:  "nested key conditions",
			input: Key("foo").Equal(Value(5)).And(Key("bar").Equal(Value(1)).And(Key("baz").BeginsWith("yar"))),
			err:   invalidKeyConditionFormat,
		},
		{
			name:  "operand error",
			input: Key("").Equal(Value("yikes")),
			err:   invalidKeyConditionOperand,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := c.input.buildTree()
			if c.err != noKeyConditionError {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := string(c.err), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expectedNode, actual; !reflect.DeepEqual(a, e) {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}

func TestKeyConditionBuildChildNodes(t *testing.T) {
	cases := []struct {
		name     string
		input    KeyConditionBuilder
		expected []exprNode
		err      keyCondErrorMode
	}{
		{
			name:  "build child nodes",
			input: Key("foo").Equal(Value("bar")).And(Key("baz").LessThan(Value(10))),
			expected: []exprNode{
				{
					children: []exprNode{
						{
							names:   []string{"foo"},
							fmtExpr: "$n",
						},
						{
							values: []types.AttributeValue{
								&types.AttributeValueMemberS{Value: "bar"},
							},
							fmtExpr: "$v",
						},
					},
					fmtExpr: "$c = $c",
				},
				{
					children: []exprNode{
						{
							names:   []string{"baz"},
							fmtExpr: "$n",
						},
						{
							values: []types.AttributeValue{
								&types.AttributeValueMemberN{Value: "10"},
							},
							fmtExpr: "$v",
						},
					},
					fmtExpr: "$c < $c",
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := c.input.buildChildNodes()
			if c.err != noKeyConditionError {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := string(c.err), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expected, actual; !reflect.DeepEqual(a, e) {
					t.Errorf("expect %#v, got %#v", e, a)
				}
			}
		})
	}
}

func TestKeyAndMultiAttributeBuild(t *testing.T) {
	eq := func(name string) KeyConditionBuilder { return Key(name).Equal(Value("v")) }
	gt := Key("sk4").GreaterThan(Value(1))

	cases := map[string]struct {
		input        KeyConditionBuilder
		expectedExpr string
		expectedLen  int
		err          keyCondErrorMode
	}{
		"variadic three conditions": {
			input:        KeyAnd(eq("pk1"), eq("pk2"), gt),
			expectedExpr: "(#0 = :0) AND (#1 = :1) AND (#2 > :2)",
			expectedLen:  3,
		},
		"variadic all equality": {
			input:        KeyAnd(eq("pk1"), eq("pk2"), eq("sk1")),
			expectedExpr: "(#0 = :0) AND (#1 = :1) AND (#2 = :2)",
			expectedLen:  3,
		},
		"four partition keys and four sort keys": {
			input: KeyAnd(
				eq("pk1"), eq("pk2"), eq("pk3"), eq("pk4"),
				eq("sk1"), eq("sk2"), eq("sk3"), gt,
			),
			expectedExpr: "(#0 = :0) AND (#1 = :1) AND (#2 = :2) AND (#3 = :3) AND (#4 = :4) AND (#5 = :5) AND (#6 = :6) AND (#7 > :7)",
			expectedLen:  8,
		},
		"chained and": {
			input:        eq("pk1").And(eq("pk2")).And(eq("sk1")).And(gt),
			expectedExpr: "(#0 = :0) AND (#1 = :1) AND (#2 = :2) AND (#3 > :3)",
			expectedLen:  4,
		},
		"variadic with leading and": {
			input:        KeyAnd(KeyAnd(eq("pk1"), eq("pk2")), gt),
			expectedExpr: "(#0 = :0) AND (#1 = :1) AND (#2 > :2)",
			expectedLen:  3,
		},
		"no conditions": {
			input: KeyAnd(),
			err:   invalidKeyConditionFormat,
		},
		"single condition": {
			input: KeyAnd(eq("pk1")),
			err:   invalidKeyConditionFormat,
		},
		"single and condition": {
			input: KeyAnd(KeyAnd(eq("pk1"), eq("pk2"))),
			err:   invalidKeyConditionFormat,
		},
		"non terminal condition is not equal": {
			input: KeyAnd(eq("pk1"), gt, eq("sk1")),
			err:   invalidKeyConditionFormat,
		},
		"leading and with non equal last becomes non terminal": {
			input: KeyAnd(KeyAnd(eq("pk1"), gt), eq("sk1")),
			err:   invalidKeyConditionFormat,
		},
		"last condition is and": {
			input: KeyAnd(eq("pk1"), KeyAnd(eq("pk2"), eq("sk1"))),
			err:   invalidKeyConditionFormat,
		},
		"non leading and": {
			input: KeyAnd(eq("pk1"), KeyAnd(eq("pk2"), eq("sk1")), gt),
			err:   invalidKeyConditionFormat,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			expr, err := NewBuilder().WithKeyCondition(c.input).Build()
			if c.err != noKeyConditionError {
				if err == nil {
					t.Fatalf("expect error %q, got no error", c.err)
				}
				if e, a := string(c.err), err.Error(); !strings.Contains(a, e) {
					t.Errorf("expect %q error message to be in %q", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got unexpected Error %q", err)
			}

			if e, a := c.expectedExpr, *expr.KeyCondition(); e != a {
				t.Errorf("expect %q, got %q", e, a)
			}
			if e, a := c.expectedLen, len(expr.Names()); e != a {
				t.Errorf("expect %d names, got %d", e, a)
			}
			if e, a := c.expectedLen, len(expr.Values()); e != a {
				t.Errorf("expect %d values, got %d", e, a)
			}
		})
	}
}
