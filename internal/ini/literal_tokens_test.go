package ini

import (
	"reflect"
	"testing"
)

// TODO: test errors
func TestNewLiteralToken(t *testing.T) {
	cases := []struct {
		name          string
		b             []rune
		expectedRead  int
		expectedToken Token
		expectedError bool
	}{
		{
			name:         "numbers",
			b:            []rune("123"),
			expectedRead: 3,
			expectedToken: newToken(TokenLit,
				[]rune("123"),
				StringType,
			),
		},
		{
			name:         "decimal",
			b:            []rune("123.456"),
			expectedRead: 7,
			expectedToken: newToken(TokenLit,
				[]rune("123.456"),
				StringType,
			),
		},
		{
			name:         "two numbers",
			b:            []rune("123 456"),
			expectedRead: 3,
			expectedToken: newToken(TokenLit,
				[]rune("123"),
				StringType,
			),
		},
		{
			name:         "number followed by alpha",
			b:            []rune("123 abc"),
			expectedRead: 3,
			expectedToken: newToken(TokenLit,
				[]rune("123"),
				StringType,
			),
		},
		{
			name:         "quoted string followed by number",
			b:            []rune(`"Hello" 123`),
			expectedRead: 7,
			expectedToken: newToken(TokenLit,
				[]rune("Hello"),
				QuotedStringType,
			),
		},
		{
			name:         "quoted string",
			b:            []rune(`"Hello World"`),
			expectedRead: 13,
			expectedToken: newToken(TokenLit,
				[]rune("Hello World"),
				QuotedStringType,
			),
		},
		{
			name:         "boolean true",
			b:            []rune("true"),
			expectedRead: 4,
			expectedToken: newToken(TokenLit,
				[]rune("true"),
				StringType,
			),
		},
		{
			name:         "boolean false",
			b:            []rune("false"),
			expectedRead: 5,
			expectedToken: newToken(TokenLit,
				[]rune("false"),
				StringType,
			),
		},
		{
			name:         "utf8 whitespace",
			b:            []rune("00"),
			expectedRead: 3,
			expectedToken: newToken(TokenLit,
				[]rune("0"),
				StringType,
			),
		},
		{
			name:         "utf8 whitespace expr",
			b:            []rune("0=00"),
			expectedRead: 1,
			expectedToken: newToken(TokenLit,
				[]rune("0"),
				StringType,
			),
		},
	}

	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tok, n, err := newLitToken(c.b)

			if e, a := c.expectedToken.ValueType, tok.ValueType; !reflect.DeepEqual(e, a) {
				t.Errorf("%d: expected %v, but received %v", i+1, e, a)
			}

			if e, a := c.expectedRead, n; e != a {
				t.Errorf("%d: expected %v, but received %v", i+1, e, a)
			}

			if e, a := c.expectedError, err != nil; e != a {
				t.Errorf("%d: expected %v, but received %v", i+1, e, a)
			}
		})
	}
}

func TestNewStringValue(t *testing.T) {
	const expect = "abc123"

	actual, err := NewStringValue(expect)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}

	if e, a := StringType, actual.Type; e != a {
		t.Errorf("expect %v type got %v", e, a)
	}
	if e, a := expect, actual.str; e != a {
		t.Errorf("expect %v string got %v", e, a)
	}
}
