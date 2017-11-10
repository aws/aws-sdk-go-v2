package protocol

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetValue(t *testing.T) {
	cases := []struct {
		val           interface{}
		expectedErr   error
		expectedValue string
		expectedEqual bool
	}{
		{
			"",
			&ErrValueNotSet{},
			"",
			true,
		},
		{
			"1",
			nil,
			"1",
			true,
		},
		{
			"Foo",
			nil,
			"Foo",
			true,
		},
		{
			"Foo",
			nil,
			"Bar",
			false,
		},
	}

	for _, c := range cases {
		v, err := GetValue(reflect.ValueOf(c.val))
		if err != c.expectedErr {
			t.Errorf("expected %v, but received %v", err, c.expectedErr)
		}

		if c.expectedEqual == (v != c.expectedValue) {
			t.Errorf("expected %v, but received %v, when they should be %v", v, c.expectedValue, c.expectedEqual)
		}
	}
}

func TestNotSetError(t *testing.T) {
	cases := []struct {
		err      error
		expected bool
	}{
		{
			nil,
			false,
		},
		{
			&ErrValueNotSet{},
			true,
		},
		{
			fmt.Errorf(""),
			false,
		},
	}

	for _, c := range cases {
		if actual := IsNotSetError(c.err); actual != c.expected {
			t.Errorf("expected %v, but received %v", actual, c.expected)
		}
	}
}
