package ini

import (
	"testing"
)

func TestStringValue(t *testing.T) {
	cases := []struct {
		b             []rune
		expectedRead  int
		expectedError bool
		expectedValue string
	}{
		{
			b:             []rune(`"foo"`),
			expectedRead:  5,
			expectedValue: `"foo"`,
		},
		{
			b:             []rune(`"123 !$_ 456 abc"`),
			expectedRead:  17,
			expectedValue: `"123 !$_ 456 abc"`,
		},
		{
			b:             []rune("foo"),
			expectedError: true,
		},
		{
			b:             []rune(` "foo"`),
			expectedError: true,
		},
	}

	for i, c := range cases {
		n, err := getStringValue(c.b)

		if e, a := c.expectedValue, string(c.b[:n]); e != a {
			t.Errorf("%d: expected %v, but received %v", i, e, a)
		}

		if e, a := c.expectedRead, n; e != a {
			t.Errorf("%d: expected %v, but received %v", i, e, a)
		}

		if e, a := c.expectedError, err != nil; e != a {
			t.Errorf("%d: expected %v, but received %v", i, e, a)
		}
	}
}
