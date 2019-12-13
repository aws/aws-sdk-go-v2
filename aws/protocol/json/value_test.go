package json

import (
	"bytes"
	"testing"
	"time"
)

func TestValue(t *testing.T) {
	cases := map[string]struct {
		setter   func(Value)
		expected string
	}{
		"string value": {
			setter: func(value Value) {
				value.String("foo")
			},
			expected: `"foo"`,
		},
		"string escaped": {
			setter: func(value Value) {
				value.String(`{"foo":"bar"}`)
			},
			expected: `"{\"foo\":\"bar\"}"`,
		},
		"integer": {
			setter: func(value Value) {
				value.Integer(1024)
			},
			expected: `1024`,
		},
		"float": {
			setter: func(value Value) {
				value.Float(1e20)
			},
			expected: `100000000000000000000`,
		},
		"float exponent component": {
			setter: func(value Value) {
				value.Float(3e22)
			},
			expected: `3e+22`,
		},
		"boolean true": {
			setter: func(value Value) {
				value.Boolean(true)
			},
			expected: `true`,
		},
		"boolean false": {
			setter: func(value Value) {
				value.Boolean(false)
			},
			expected: `false`,
		},
		"byte slice": {
			setter: func(value Value) {
				value.ByteSlice([]byte("foo bar"))
			},
			expected: `"Zm9vIGJhcg=="`,
		},
		"byte slice nil": {
			setter: func(value Value) {
				value.ByteSlice(nil)
			},
			expected: `null`,
		},
		"time": {
			setter: func(value Value) {
				value.Time(time.Date(2019, 1, 2, 3, 4, 5, 6, time.UTC), "iso8601")
			},
			expected: `"2019-01-02T03:04:05Z"`,
		},
		"object": {
			setter: func(value Value) {
				o := value.Object()
				defer o.Close()
				o.Key("key").String("value")
			},
			expected: `{"key":"value"}`,
		},
		"array": {
			setter: func(value Value) {
				o := value.Array()
				defer o.Close()
				o.Value().String("value1")
				o.Value().String("value2")
			},
			expected: `["value1","value2"]`,
		},
		"null": {
			setter: func(value Value) {
				value.Null()
			},
			expected: `null`,
		},
	}
	scratch := make([]byte, 64)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			value := newValue(&b, &scratch)

			tt.setter(value)

			if e, a := []byte(tt.expected), b.Bytes(); bytes.Compare(e, a) != 0 {
				t.Errorf("expected %+q, but got %+q", e, a)
			}
		})
	}
}
