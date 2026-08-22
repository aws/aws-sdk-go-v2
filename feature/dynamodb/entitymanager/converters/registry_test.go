package converters

import (
	"strconv"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if r.converters == nil {
		t.Fatal("NewRegistry did not initialize converters map")
	}
}

func TestRegistry_Clone(t *testing.T) {
	r := DefaultRegistry.Clone()
	if r == nil {
		t.Fatal("Clone returned nil")
	}
	if r.converters == nil {
		t.Fatal("Clone did not initialize converters map")
	}
	if len(r.converters) != len(DefaultRegistry.converters) {
		t.Errorf("Clone did not copy all converters: got %d, want %d", len(r.converters), len(DefaultRegistry.converters))
	}
}

func TestRegistry_Add(t *testing.T) {
	r := &Registry{}
	initial := len(r.converters)
	r.Add("mock", &mockConverter{})
	if len(r.converters) != initial+1 {
		t.Errorf("Add did not increase converter count: got %d, want %d", len(r.converters), initial+1)
	}
	if r.Converter("mock") == nil {
		t.Error("Add did not register converter under 'mock'")
	}
}

func TestRegistry_Remove(t *testing.T) {
	r := DefaultRegistry.Clone()
	initial := len(r.converters)
	r.Remove("json")
	if len(r.converters) != initial-1 {
		t.Errorf("Remove did not decrease converter count: got %d, want %d", len(r.converters), initial-1)
	}
	if r.Converter("json") != nil {
		t.Error("Remove did not remove 'json' converter")
	}
}

func TestRegistry_Converter(t *testing.T) {
	r := DefaultRegistry.Clone()
	cases := []struct {
		name string
		ok   bool
	}{
		{"404", false},
		{"json", true},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			conv := r.Converter(c.name)
			if (conv != nil) != c.ok {
				t.Errorf("Converter(%q) presence = %v, want %v", c.name, conv != nil, c.ok)
			}
		})
	}
}

func TestRegistry_ConverterFor(t *testing.T) {
	r := DefaultRegistry.Clone()

	t.Run("known type", func(t *testing.T) {
		var s string
		conv := r.ConverterFor(s)
		if conv == nil {
			t.Errorf("expected converter for string, got nil")
		}
	})

	t.Run("unknown type", func(t *testing.T) {
		type custom struct{}
		var c custom
		conv := r.ConverterFor(c)
		if conv != nil {
			t.Errorf("expected no converter for custom type, got one")
		}
	})
}
