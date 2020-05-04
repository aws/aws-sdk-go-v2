package rest

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

func TestURIValue(t *testing.T) {
	const uriKey = "someKey"
	const path = "/some/{someKey}/{path+}"

	type expected struct {
		path string
		raw  string
	}

	cases := map[string]struct {
		path     string
		args     []interface{}
		expected expected
	}{
		"string": {
			path: path,
			args: []interface{}{"someValue"},
			expected: expected{
				path: "/some/someValue/{path+}",
				raw:  "/some/someValue/{path+}",
			},
		},
		"float64": {
			path: path,
			args: []interface{}{3.14159},
			expected: expected{
				path: "/some/3.14159/{path+}",
				raw:  "/some/3.14159/{path+}",
			},
		},
		"bool": {
			path: path,
			args: []interface{}{true},
			expected: expected{
				path: "/some/true/{path+}",
				raw:  "/some/true/{path+}",
			},
		},
		"time": {
			path: path,
			args: []interface{}{time.Unix(0, 0), protocol.ISO8601TimeFormatName},
			expected: expected{
				path: "/some/1970-01-01T00:00:00Z/{path+}",
				raw:  "/some/1970-01-01T00%3A00%3A00Z/{path+}",
			},
		},
		"byte slice": {
			path: path,
			args: []interface{}{[]byte("baz")},
			expected: expected{
				path: "/some/baz/{path+}",
				raw:  "/some/baz/{path+}",
			},
		},
	}

	buffer := make([]byte, 1024)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			pBytes, rBytes := []byte(tt.path), []byte(tt.path)

			uv := newURIValue(&pBytes, &rBytes, &buffer, uriKey)

			if err := setURI(uv, tt.args); err != nil {
				t.Fatalf("expected no error, %v", err)
			}

			if e, a := tt.expected.path, string(pBytes); e != a {
				t.Errorf("expected %v, got %v", e, a)
			}

			if e, a := tt.expected.raw, string(rBytes); e != a {
				t.Errorf("expected %v, got %v", e, a)
			}
		})
	}
}

func setURI(uv URIValue, args []interface{}) error {
	value := args[0]

	switch value.(type) {
	case string:
		return reflectCall(reflect.ValueOf(uv.String), args)
	case float32:
		return reflectCall(reflect.ValueOf(uv.Float), args)
	case float64:
		return reflectCall(reflect.ValueOf(uv.Double), args)
	case bool:
		return reflectCall(reflect.ValueOf(uv.Boolean), args)
	case time.Time:
		return reflectCall(reflect.ValueOf(uv.Timestamp), args)
	case []byte:
		return reflectCall(reflect.ValueOf(uv.Blob), args)
	default:
		return fmt.Errorf("unhandled value type")
	}
}
