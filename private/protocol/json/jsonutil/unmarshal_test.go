package jsonutil

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestUnmarshal(t *testing.T) {
	type scalarMaps struct {
		Bools    map[string]bool
		Strings  map[string]string
		Int64s   map[string]int64
		Float64s map[string]float64
		Times    map[string]time.Time
	}

	cases := map[string]struct {
		Input  json.RawMessage
		Output interface{}
		Expect interface{}
	}{
		"scalar maps": {
			Input: json.RawMessage(`
{
"Bools": {"a": true},
"Strings": {"b": "123"},
"Int64s": {"c": 456},
"Float64s": {"d": 789},
"Times": {"e": 1257894000}
}`),
			Output: &scalarMaps{},
			Expect: &scalarMaps{
				Bools:    map[string]bool{"a": true},
				Strings:  map[string]string{"b": "123"},
				Int64s:   map[string]int64{"c": 456},
				Float64s: map[string]float64{"d": 789},
				Times:    map[string]time.Time{"e": time.Unix(1257894000, 0).UTC()},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := UnmarshalJSON(c.Output, bytes.NewReader(c.Input))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := c.Expect, c.Output; !reflect.DeepEqual(e, a) {
				t.Errorf("expect:\n%#v\nactual:\n%#v\n", e, a)
			}
		})
	}
}
