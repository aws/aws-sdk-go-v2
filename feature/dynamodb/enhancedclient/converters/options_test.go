package converters

import (
	"strconv"
	"testing"
)

func TestGetOpt(t *testing.T) {
	cases := []struct {
		opts     []string
		name     string
		expected string
	}{
		{},
		{
			opts:     []string{"test"},
			name:     "test",
			expected: "",
		},
		{
			opts:     []string{"test="},
			name:     "test",
			expected: "",
		},
		{
			opts:     []string{"test=test"},
			name:     "test",
			expected: "test",
		},
		{
			opts:     []string{"TEST=test"},
			name:     "test",
			expected: "",
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := getOpt(c.opts, c.name)

			if actual != c.expected {
				t.Fatalf(`expected "%s", got "%s"`, c.expected, actual)
			}
		})
	}
}
