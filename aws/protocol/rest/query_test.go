package rest

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

func TestQueryValue(t *testing.T) {
	const queryKey = "someKey"

	cases := map[string]struct {
		values   url.Values
		args     []interface{}
		append   bool
		expected url.Values
	}{
		"set string": {
			values: url.Values{queryKey: []string{"foobar"}},
			args:   []interface{}{"string value"},
			expected: map[string][]string{
				queryKey: {"string value"},
			},
		},
		"set float64": {
			values: url.Values{queryKey: []string{"foobar"}},
			args:   []interface{}{3.14159},
			expected: map[string][]string{
				queryKey: {"3.14159"},
			},
		},
		"set bool": {
			values: url.Values{queryKey: []string{"foobar"}},
			args:   []interface{}{true},
			expected: map[string][]string{
				queryKey: {"true"},
			},
		},
		"set json": {
			values: url.Values{queryKey: []string{"foobar"}},
			args:   []interface{}{aws.JSONValue{"jsonKey": "jsonValue"}},
			expected: map[string][]string{
				queryKey: {`{"jsonKey":"jsonValue"}`},
			},
		},
		"set time": {
			values: url.Values{queryKey: []string{"foobar"}},
			args:   []interface{}{time.Unix(0, 0), protocol.ISO8601TimeFormatName},
			expected: map[string][]string{
				queryKey: {"1970-01-01T00:00:00Z"},
			},
		},
		"set byte slice": {
			values: url.Values{queryKey: []string{"foobar"}},
			args:   []interface{}{[]byte("baz")},
			expected: map[string][]string{
				queryKey: {"YmF6"},
			},
		},
		"add string": {
			values: url.Values{queryKey: []string{"other string"}},
			args:   []interface{}{"string value"},
			append: true,
			expected: map[string][]string{
				queryKey: {"other string", "string value"},
			},
		},
		"add float64": {
			values: url.Values{queryKey: []string{"1.61803"}},
			args:   []interface{}{3.14159},
			append: true,
			expected: map[string][]string{
				queryKey: {"1.61803", "3.14159"},
			},
		},
		"add bool": {
			values: url.Values{queryKey: []string{"false"}},
			args:   []interface{}{true},
			append: true,
			expected: map[string][]string{
				queryKey: {"false", "true"},
			},
		},
		"add json": {
			values: url.Values{queryKey: []string{`{"someKey":"someValue"}`}},
			args:   []interface{}{aws.JSONValue{"jsonKey": "jsonValue"}},
			append: true,
			expected: map[string][]string{
				queryKey: {`{"someKey":"someValue"}`, `{"jsonKey":"jsonValue"}`},
			},
		},
		"add time": {
			values: url.Values{queryKey: []string{"1991-09-17T00:00:00Z"}},
			args:   []interface{}{time.Unix(0, 0), protocol.ISO8601TimeFormatName},
			append: true,
			expected: map[string][]string{
				queryKey: {"1991-09-17T00:00:00Z", "1970-01-01T00:00:00Z"},
			},
		},
		"add byte slice": {
			values: url.Values{queryKey: []string{"YmFy"}},
			args:   []interface{}{[]byte("baz")},
			append: true,
			expected: map[string][]string{
				queryKey: {"YmFy", "YmF6"},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			if tt.values == nil {
				tt.values = url.Values{}
			}

			qv := newQueryValue(tt.values, queryKey, tt.append)

			if err := setQueryValue(qv, tt.args); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if e, a := tt.expected, qv.query; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, got %v", e, a)
			}
		})
	}
}

func setQueryValue(qv QueryValue, args []interface{}) error {
	value := args[0]

	switch value.(type) {
	case string:
		return reflectCall(reflect.ValueOf(qv.String), args)
	case float64:
		return reflectCall(reflect.ValueOf(qv.Float), args)
	case bool:
		return reflectCall(reflect.ValueOf(qv.Boolean), args)
	case aws.JSONValue:
		return reflectCall(reflect.ValueOf(qv.JSONValue), args)
	case time.Time:
		return reflectCall(reflect.ValueOf(qv.Time), args)
	case []byte:
		return reflectCall(reflect.ValueOf(qv.ByteSlice), args)
	default:
		return fmt.Errorf("unhandled query value type")
	}
}
