package attributevalue

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected types.AttributeValue
		err      bool
	}{
		{
			name:     "types.AttributeValueMemberS",
			input:    `{"S":"test"}`,
			expected: &types.AttributeValueMemberS{Value: "test"},
		},
		{
			name:     "types.AttributeValueMemberN",
			input:    `{"N":"1.5"}`,
			expected: &types.AttributeValueMemberN{Value: "1.5"},
		},
		{
			name:     "types.AttributeValueMemberB",
			input:    `{"B":"dGVzdA=="}`,
			expected: &types.AttributeValueMemberB{Value: []byte("test")},
		},
		{
			name:     "types.AttributeValueMemberBOOL true",
			input:    `{"BOOL":true}`,
			expected: &types.AttributeValueMemberBOOL{Value: true},
		},
		{
			name:     "types.AttributeValueMemberBOOL false",
			input:    `{"BOOL":false}`,
			expected: &types.AttributeValueMemberBOOL{Value: false},
		},
		{
			name:     "types.AttributeValueMemberNULL true",
			input:    `{"NULL":true}`,
			expected: &types.AttributeValueMemberNULL{Value: true},
		},
		{
			name:     "types.AttributeValueMemberNULL false",
			input:    `{"NULL":false}`,
			expected: &types.AttributeValueMemberNULL{Value: false},
		},
		{
			name:     "types.AttributeValueMemberSS",
			input:    `{"SS":["test"]}`,
			expected: &types.AttributeValueMemberSS{Value: []string{"test"}},
		},
		{
			name:     "types.AttributeValueMemberNS",
			input:    `{"NS":["4.2"]}`,
			expected: &types.AttributeValueMemberNS{Value: []string{"4.2"}},
		},
		{
			name:     "types.AttributeValueMemberBS",
			input:    `{"BS":["dGVzdA=="]}`,
			expected: &types.AttributeValueMemberBS{Value: [][]byte{[]byte("test")}},
		},
		{
			name: "types.AttributeValueMemberM",
			input: `{
  "M": {
    "types.AttributeValueMemberB": {"B": "dGVzdA=="},
    "types.AttributeValueMemberBOOL:false": {"BOOL": false},
    "types.AttributeValueMemberBOOL:true": {"BOOL": true},
    "types.AttributeValueMemberBS": {"BS": ["dGVzdA=="]},
    "types.AttributeValueMemberM": {
      "M": {
        "types.AttributeValueMemberM": {
          "M": {
            "types.AttributeValueMemberM": {
              "M": {
                "types.AttributeValueMemberS": {
                  "S": "test"
                }
              }
            }
          }
        }
      }
    },
    "types.AttributeValueMemberN": {"N": "1.5"},
    "types.AttributeValueMemberNS": {"NS": ["4.2"]},
    "types.AttributeValueMemberNULL:false": {"NULL": false},
    "types.AttributeValueMemberNULL:true": {"NULL": true},
    "types.AttributeValueMemberS": {"S": "test"},
    "types.AttributeValueMemberSS": {"SS": ["test"]}
  }
}`,
			expected: &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"types.AttributeValueMemberS": &types.AttributeValueMemberS{
						Value: "test",
					},
					"types.AttributeValueMemberN": &types.AttributeValueMemberN{
						Value: "1.5",
					},
					"types.AttributeValueMemberB": &types.AttributeValueMemberB{
						Value: []byte("test"),
					},
					"types.AttributeValueMemberBOOL:true": &types.AttributeValueMemberBOOL{
						Value: true,
					},
					"types.AttributeValueMemberBOOL:false": &types.AttributeValueMemberBOOL{
						Value: false,
					},
					"types.AttributeValueMemberNULL:true": &types.AttributeValueMemberNULL{
						Value: true,
					},
					"types.AttributeValueMemberNULL:false": &types.AttributeValueMemberNULL{
						Value: false,
					},
					"types.AttributeValueMemberSS": &types.AttributeValueMemberSS{
						Value: []string{"test"},
					},
					"types.AttributeValueMemberNS": &types.AttributeValueMemberNS{
						Value: []string{"4.2"},
					},
					"types.AttributeValueMemberBS": &types.AttributeValueMemberBS{
						Value: [][]byte{[]byte("test")},
					},
					"types.AttributeValueMemberM": &types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"types.AttributeValueMemberM": &types.AttributeValueMemberM{
								Value: map[string]types.AttributeValue{
									"types.AttributeValueMemberM": &types.AttributeValueMemberM{
										Value: map[string]types.AttributeValue{
											"types.AttributeValueMemberS": &types.AttributeValueMemberS{
												Value: "test",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "types.AttributeValueMemberL",
			input: `{"L":[{"B": "dGVzdA=="},{"M":{"S":{"S":"test"}}}]}`,
			expected: &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberB{
						Value: []byte("test"),
					},
					&types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"S": &types.AttributeValueMemberS{
								Value: "test",
							},
						},
					},
				},
			},
		},
		{
			name:     "broken input",
			input:    `{"L":`,
			expected: nil,
			err:      true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v types.AttributeValue
			var err error
			v, err = UnmarshalJSON([]byte(c.input))
			if c.err {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if !reflect.DeepEqual(v, c.expected) {
				t.Fatalf("expected %v, got %v", c.expected, v)
			}
		})
	}
}

func TestUnmarshalListJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []types.AttributeValue
		err      bool
	}{
		{
			name:  "1",
			input: `[{"S":"test"}]`,
			expected: []types.AttributeValue{
				&types.AttributeValueMemberS{
					Value: "test",
				},
			},
		},
		{
			name:     "broken input",
			input:    `[`,
			expected: nil,
			err:      true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v []types.AttributeValue
			var err error
			v, err = UnmarshalListJSON([]byte(c.input))
			if c.err {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if !reflect.DeepEqual(v, c.expected) {
				t.Fatalf("expected %v, got %v", c.expected, v)
			}
		})
	}
}

func TestUnmarshalMapJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected map[string]types.AttributeValue
		err      bool
	}{
		{
			name:  "1",
			input: `{"test":{"S":"test"}}`,
			expected: map[string]types.AttributeValue{
				"test": &types.AttributeValueMemberS{
					Value: "test",
				},
			},
		},
		{
			name:     "bad input",
			input:    `["asd"]`,
			expected: nil,
			err:      true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v map[string]types.AttributeValue
			var err error
			v, err = UnmarshalMapJSON([]byte(c.input))
			if c.err {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if !reflect.DeepEqual(v, c.expected) {
				t.Fatalf("expected %v, got %v", c.expected, v)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    types.AttributeValue
		expected string
	}{
		{
			name: "&types.AttributeValueMemberB",
			input: &types.AttributeValueMemberB{
				Value: []byte("test"),
			},
			expected: `{"B":"dGVzdA=="}`,
		},
		{
			name:     "&types.AttributeValueMemberB empty",
			input:    &types.AttributeValueMemberB{},
			expected: `{"B":null}`,
		},
		{
			name:     "&types.AttributeValueMemberBOOL empty",
			input:    &types.AttributeValueMemberBOOL{},
			expected: `{"BOOL":false}`,
		},
		{
			name: "&types.AttributeValueMemberBOOL true",
			input: &types.AttributeValueMemberBOOL{
				Value: true,
			},
			expected: `{"BOOL":true}`,
		},
		{
			name:     "&types.AttributeValueMemberBS empty",
			input:    &types.AttributeValueMemberBS{},
			expected: `{"BS":[]}`,
		},
		{
			name: "&types.AttributeValueMemberBS",
			input: &types.AttributeValueMemberBS{
				Value: [][]byte{
					[]byte("test"),
				},
			},
			expected: `{"BS":["dGVzdA=="]}`,
		},
		{
			name:     "&types.AttributeValueMemberL empty",
			input:    &types.AttributeValueMemberL{},
			expected: `{"L":[]}`,
		},
		{
			name: "&types.AttributeValueMemberL",
			input: &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberB{
						Value: []byte("test"),
					},
					&types.AttributeValueMemberBOOL{
						Value: true,
					},
					&types.AttributeValueMemberBS{
						Value: [][]byte{
							[]byte("test"),
						},
					},
				},
			},
			expected: `{"L":[{"B":"dGVzdA=="},{"BOOL":true},{"BS":["dGVzdA=="]}]}`,
		},
		{
			name:     "&types.AttributeValueMemberM empty",
			input:    &types.AttributeValueMemberM{},
			expected: `{"M":{}}`,
		},
		{
			name: "&types.AttributeValueMemberM",
			input: &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					// we only use 1 key here because go does not guarantee key order and test will fail randomly
					"testList": &types.AttributeValueMemberL{
						Value: []types.AttributeValue{
							&types.AttributeValueMemberB{
								Value: []byte("test"),
							},
							&types.AttributeValueMemberBOOL{
								Value: true,
							},
							&types.AttributeValueMemberBS{
								Value: [][]byte{
									[]byte("test"),
								},
							},
						},
					},
				},
			},
			expected: `{"M":{"testList":{"L":[{"B":"dGVzdA=="},{"BOOL":true},{"BS":["dGVzdA=="]}]}}}`,
		},
		{
			name:     "&types.AttributeValueMemberN empty",
			input:    &types.AttributeValueMemberN{},
			expected: `{"N":""}`,
		},
		{
			name: "&types.AttributeValueMemberN 1",
			input: &types.AttributeValueMemberN{
				Value: "123",
			},
			expected: `{"N":"123"}`,
		},
		{
			name: "&types.AttributeValueMemberN 2",
			input: &types.AttributeValueMemberN{
				Value: "123.456",
			},
			expected: `{"N":"123.456"}`,
		},
		{
			name: "&types.AttributeValueMemberN max max int",
			input: &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d.%d", math.MaxInt, math.MaxInt),
			},
			expected: fmt.Sprintf(`{"N":"%d.%d"}`, math.MaxInt, math.MaxInt),
		},
		{
			name: "&types.AttributeValueMemberN min max int",
			input: &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d.%d", math.MinInt, math.MaxInt),
			},
			expected: fmt.Sprintf(`{"N":"%d.%d"}`, math.MinInt, math.MaxInt),
		},
		{
			name:     "&types.AttributeValueMemberNS empty",
			input:    &types.AttributeValueMemberNS{},
			expected: `{"NS":[]}`,
		},
		{
			name: "&types.AttributeValueMemberNS",
			input: &types.AttributeValueMemberNS{
				Value: []string{"123", "456"},
			},
			expected: `{"NS":["123","456"]}`,
		},
		{
			name:     "&types.AttributeValueMemberNULL empty",
			input:    &types.AttributeValueMemberNULL{},
			expected: `{"NULL":false}`,
		},
		{
			name: "&types.AttributeValueMemberNULL true",
			input: &types.AttributeValueMemberNULL{
				Value: true,
			},
			expected: `{"NULL":true}`,
		},
		{
			name:     "&types.AttributeValueMemberS empty",
			input:    &types.AttributeValueMemberS{},
			expected: `{"S":""}`,
		},
		{
			name: "&types.AttributeValueMemberS",
			input: &types.AttributeValueMemberS{
				Value: "test",
			},
			expected: `{"S":"test"}`,
		},
		{
			name:     "&types.AttributeValueMemberSS empty",
			input:    &types.AttributeValueMemberSS{},
			expected: `{"SS":[]}`,
		},
		{
			name: "&types.AttributeValueMemberSS",
			input: &types.AttributeValueMemberSS{
				Value: []string{"test", "foo"},
			},
			expected: `{"SS":["test","foo"]}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := MarshalJSON(c.input)
			if err != nil {
				t.Fatalf("unexpected error, got: %v", err)
			}

			if string(actual) != c.expected {
				t.Errorf("expected %s, got %s", c.expected, string(actual))
			}
		})
	}
}

func TestMarshalListJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    []types.AttributeValue
		expected string
	}{
		{
			name:     "nil",
			input:    nil,
			expected: `[]`,
		},
		{
			name:     "empty list",
			input:    []types.AttributeValue{},
			expected: `[]`,
		},
		{
			name: "list of complex items",
			input: []types.AttributeValue{
				&types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberB{
							Value: []byte("test"),
						},
						&types.AttributeValueMemberBOOL{
							Value: true,
						},
						&types.AttributeValueMemberBS{
							Value: [][]byte{
								[]byte("test"),
							},
						},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						// we only use 1 key here because go does not guarantee key order and test will fail randomly
						"testList": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{
								&types.AttributeValueMemberB{
									Value: []byte("test"),
								},
								&types.AttributeValueMemberBOOL{
									Value: true,
								},
								&types.AttributeValueMemberBS{
									Value: [][]byte{
										[]byte("test"),
									},
								},
							},
						},
					},
				},
			},
			expected: `[{"L":[{"B":"dGVzdA=="},{"BOOL":true},{"BS":["dGVzdA=="]}]},{"M":{"testList":{"L":[{"B":"dGVzdA=="},{"BOOL":true},{"BS":["dGVzdA=="]}]}}}]`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := MarshalListJSON(c.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(result) != c.expected {
				t.Errorf("expected %v, got %v", c.expected, string(result))
			}
		})
	}
}

func TestMarshalMapJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    map[string]types.AttributeValue
		expected string
	}{
		{
			name:     "nil",
			input:    nil,
			expected: `{}`,
		},
		{
			name:     "empty list",
			input:    map[string]types.AttributeValue{},
			expected: `{}`,
		},
		{
			name: "list of complex items",
			input: map[string]types.AttributeValue{
				// we only use 1 key here because go does not guarantee key order and test will fail randomly
				"testList": &types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberB{
							Value: []byte("test"),
						},
						&types.AttributeValueMemberBOOL{
							Value: true,
						},
						&types.AttributeValueMemberBS{
							Value: [][]byte{
								[]byte("test"),
							},
						},
					},
				},
			},
			expected: `{"testList":{"L":[{"B":"dGVzdA=="},{"BOOL":true},{"BS":["dGVzdA=="]}]}}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := MarshalMapJSON(c.input)
			if err != nil {
				t.Fatalf("unexpected error, got: %v", err)
			}

			if string(result) != c.expected {
				t.Errorf("expected %v, got %v", c.expected, string(result))
			}
		})
	}
}
