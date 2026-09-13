package converters

import (
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestGetType(t *testing.T) {
	cases := []struct {
		input    any
		expected string
	}{
		{
			input:    uint(0),
			expected: "uint",
		},
		{
			input:    uint8(0),
			expected: "uint8",
		},
		{
			input:    uint16(0),
			expected: "uint16",
		},
		{
			input:    uint32(0),
			expected: "uint32",
		},
		{
			input:    uint64(0),
			expected: "uint64",
		},
		{
			input:    int(0),
			expected: "int",
		},
		{
			input:    int8(0),
			expected: "int8",
		},
		{
			input:    int16(0),
			expected: "int16",
		},
		{
			input:    int32(0),
			expected: "int32",
		},
		{
			input:    int64(0),
			expected: "int64",
		},
		{
			input:    float32(0),
			expected: "float32",
		},
		{
			input:    float64(0),
			expected: "float64",
		},
		{
			input:    aws.Uint(uint(0)),
			expected: "*uint",
		},
		{
			input:    aws.Uint8(uint8(0)),
			expected: "*uint8",
		},
		{
			input:    aws.Uint16(uint16(0)),
			expected: "*uint16",
		},
		{
			input:    aws.Uint32(uint32(0)),
			expected: "*uint32",
		},
		{
			input:    aws.Uint64(uint64(0)),
			expected: "*uint64",
		},
		{
			input:    aws.Int(int(0)),
			expected: "*int",
		},
		{
			input:    aws.Int8(int8(0)),
			expected: "*int8",
		},
		{
			input:    aws.Int16(int16(0)),
			expected: "*int16",
		},
		{
			input:    aws.Int32(int32(0)),
			expected: "*int32",
		},
		{
			input:    aws.Int64(int64(0)),
			expected: "*int64",
		},
		{
			input:    aws.Float32(float32(0)),
			expected: "*float32",
		},
		{
			input:    aws.Float64(float64(0)),
			expected: "*float64",
		},
		{
			input:    time.Time{},
			expected: "time.Time",
		},
		{
			input:    &time.Time{},
			expected: "*time.Time",
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := getType(c.input)

			if actual != c.expected {
				t.Fatalf(`expected "%s", got "%s" for %T`, c.expected, actual, c.input)
			}
		})
	}
}

func BenchmarkGetType(b *testing.B) {
	x := int8(8)
	for c := 0; c < b.N; c++ {
		_ = getType(x)
	}
}
