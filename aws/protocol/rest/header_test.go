package rest

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

func TestHeaderValue(t *testing.T) {
	const keyName = "test-key"
	const expectedKeyName = "Test-Key"

	cases := map[string]struct {
		header   http.Header
		args     []interface{}
		append   bool
		expected http.Header
	}{
		"set string": {
			header: http.Header{expectedKeyName: []string{"foobar"}},
			args:   []interface{}{"string value"},
			expected: map[string][]string{
				expectedKeyName: {"string value"},
			},
		},
		"set float64": {
			header: http.Header{expectedKeyName: []string{"foobar"}},
			args:   []interface{}{3.14159},
			expected: map[string][]string{
				expectedKeyName: {"3.14159"},
			},
		},
		"set bool": {
			header: http.Header{expectedKeyName: []string{"foobar"}},
			args:   []interface{}{true},
			expected: map[string][]string{
				expectedKeyName: {"true"},
			},
		},
		"set json": {
			header: http.Header{expectedKeyName: []string{"foobar"}},
			args:   []interface{}{aws.JSONValue{"jsonKey": "jsonValue"}},
			expected: map[string][]string{
				expectedKeyName: {"eyJqc29uS2V5IjoianNvblZhbHVlIn0="},
			},
		},
		"set time": {
			header: http.Header{expectedKeyName: []string{"foobar"}},
			args:   []interface{}{time.Unix(0, 0), protocol.ISO8601TimeFormatName},
			expected: map[string][]string{
				expectedKeyName: {"1970-01-01T00:00:00Z"},
			},
		},
		"set byte slice": {
			header: http.Header{expectedKeyName: []string{"foobar"}},
			args:   []interface{}{[]byte("baz")},
			expected: map[string][]string{
				expectedKeyName: {"YmF6"},
			},
		},
		"add string": {
			header: http.Header{expectedKeyName: []string{"other string"}},
			args:   []interface{}{"string value"},
			append: true,
			expected: map[string][]string{
				expectedKeyName: {"other string", "string value"},
			},
		},
		"add float64": {
			header: http.Header{expectedKeyName: []string{"1.61803"}},
			args:   []interface{}{3.14159},
			append: true,
			expected: map[string][]string{
				expectedKeyName: {"1.61803", "3.14159"},
			},
		},
		"add bool": {
			header: http.Header{expectedKeyName: []string{"false"}},
			args:   []interface{}{true},
			append: true,
			expected: map[string][]string{
				expectedKeyName: {"false", "true"},
			},
		},
		"add json": {
			header: http.Header{expectedKeyName: []string{`eyJzb21lS2V5Ijoic29tZVZhbHVlIn0=`}},
			args:   []interface{}{aws.JSONValue{"jsonKey": "jsonValue"}},
			append: true,
			expected: map[string][]string{
				expectedKeyName: {"eyJzb21lS2V5Ijoic29tZVZhbHVlIn0=", "eyJqc29uS2V5IjoianNvblZhbHVlIn0="},
			},
		},
		"add time": {
			header: http.Header{expectedKeyName: []string{"1991-09-17T00:00:00Z"}},
			args:   []interface{}{time.Unix(0, 0), protocol.ISO8601TimeFormatName},
			append: true,
			expected: map[string][]string{
				expectedKeyName: {"1991-09-17T00:00:00Z", "1970-01-01T00:00:00Z"},
			},
		},
		"add byte slice": {
			header: http.Header{expectedKeyName: []string{"YmFy"}},
			args:   []interface{}{[]byte("baz")},
			append: true,
			expected: map[string][]string{
				expectedKeyName: {"YmFy", "YmF6"},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			if tt.header == nil {
				tt.header = http.Header{}
			}

			hv := newHeaderValue(tt.header, keyName, tt.append)

			if err := setHeader(hv, tt.args); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if e, a := tt.expected, hv.header; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, got %v", e, a)
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	const prefix = "x-amzn-meta-"
	cases := map[string]struct {
		headers  http.Header
		values   map[string]string
		append   bool
		expected http.Header
	}{
		"set": {
			headers: http.Header{
				"X-Amzn-Meta-Foo": {"bazValue"},
			},
			values: map[string]string{
				"foo":   "fooValue",
				" bar ": "barValue",
			},
			expected: http.Header{
				"X-Amzn-Meta-Foo": {"fooValue"},
				"X-Amzn-Meta-Bar": {"barValue"},
			},
		},
		"add": {
			headers: http.Header{
				"X-Amzn-Meta-Foo": {"bazValue"},
			},
			values: map[string]string{
				"foo":   "fooValue",
				" bar ": "barValue",
			},
			append: true,
			expected: http.Header{
				"X-Amzn-Meta-Foo": {"bazValue", "fooValue"},
				"X-Amzn-Meta-Bar": {"barValue"},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			headers := Headers{header: tt.headers, prefix: prefix}

			var f func(key string) HeaderValue
			if tt.append {
				f = headers.AddHeader
			} else {
				f = headers.SetHeader
			}

			for key, value := range tt.values {
				f(key).String(value)
			}

			if e, a := tt.expected, tt.headers; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but got %v", e, a)
			}
		})
	}
}

func setHeader(hv HeaderValue, args []interface{}) error {
	value := args[0]

	switch value.(type) {
	case string:
		return reflectCall(reflect.ValueOf(hv.String), args)
	case float64:
		return reflectCall(reflect.ValueOf(hv.Float), args)
	case bool:
		return reflectCall(reflect.ValueOf(hv.Boolean), args)
	case aws.JSONValue:
		return reflectCall(reflect.ValueOf(hv.JSONValue), args)
	case time.Time:
		return reflectCall(reflect.ValueOf(hv.Time), args)
	case []byte:
		return reflectCall(reflect.ValueOf(hv.ByteSlice), args)
	default:
		return fmt.Errorf("unhandled header value type")
	}
}
