package converters

import (
	"time"
)

// DefaultRegistry is a pre-populated Registry containing converters for
// primitive numeric and boolean types, their pointer forms, time.Time and
// *time.Time, byte slices ([]byte / []uint8), strings (string / *string), and
// a generic JSON converter ("json").
//
// The keys in the registry map are the string representations returned by the
// internal getType helper (e.g. "int", "*int", "time.Time", "*time.Time").
//
// DefaultRegistry is intended as a starting point. You may call Clone() to
// obtain an isolated copy and then Add or Remove converters without affecting
// the shared defaults. Direct mutation of DefaultRegistry (calling Add or
// Remove on it) is safe only if done during program initialization; concurrent
// writes without external synchronization are not supported.
var DefaultRegistry = &Registry{
	converters: map[string]AnyAttributeConverter{
		// numbers
		"uint":    &Wrapper[uint]{Impl: &NumericConverter[uint]{}},
		"uint8":   &Wrapper[uint8]{Impl: &NumericConverter[uint8]{}},
		"uint16":  &Wrapper[uint16]{Impl: &NumericConverter[uint16]{}},
		"uint32":  &Wrapper[uint32]{Impl: &NumericConverter[uint32]{}},
		"uint64":  &Wrapper[uint64]{Impl: &NumericConverter[uint64]{}},
		"int":     &Wrapper[int]{Impl: &NumericConverter[int]{}},
		"int8":    &Wrapper[int8]{Impl: &NumericConverter[int8]{}},
		"int16":   &Wrapper[int16]{Impl: &NumericConverter[int16]{}},
		"int32":   &Wrapper[int32]{Impl: &NumericConverter[int32]{}},
		"int64":   &Wrapper[int64]{Impl: &NumericConverter[int64]{}},
		"float32": &Wrapper[float32]{Impl: &NumericConverter[float32]{}},
		"float64": &Wrapper[float64]{Impl: &NumericConverter[float64]{}},
		// numbers pointers
		"*uint":    &Wrapper[*uint]{Impl: &NumericPtrConverter[uint]{}},
		"*uint8":   &Wrapper[*uint8]{Impl: &NumericPtrConverter[uint8]{}},
		"*uint16":  &Wrapper[*uint16]{Impl: &NumericPtrConverter[uint16]{}},
		"*uint32":  &Wrapper[*uint32]{Impl: &NumericPtrConverter[uint32]{}},
		"*uint64":  &Wrapper[*uint64]{Impl: &NumericPtrConverter[uint64]{}},
		"*int":     &Wrapper[*int]{Impl: &NumericPtrConverter[int]{}},
		"*int8":    &Wrapper[*int8]{Impl: &NumericPtrConverter[int8]{}},
		"*int16":   &Wrapper[*int16]{Impl: &NumericPtrConverter[int16]{}},
		"*int32":   &Wrapper[*int32]{Impl: &NumericPtrConverter[int32]{}},
		"*int64":   &Wrapper[*int64]{Impl: &NumericPtrConverter[int64]{}},
		"*float32": &Wrapper[*float32]{Impl: &NumericPtrConverter[float32]{}},
		"*float64": &Wrapper[*float64]{Impl: &NumericPtrConverter[float64]{}},
		// other
		"bool":       &Wrapper[bool]{Impl: &BoolConverter{}},
		"*bool":      &Wrapper[*bool]{Impl: &BoolPtrConverter{}},
		"[]uint8":    &Wrapper[[]uint8]{Impl: &ByteArrayConverter{}},
		"[]byte":     &Wrapper[[]byte]{Impl: &ByteArrayConverter{}},
		"string":     &Wrapper[string]{Impl: &StringConverter{}},
		"*string":    &Wrapper[*string]{Impl: &StringPtrConverter{}},
		"time.Time":  &Wrapper[time.Time]{Impl: &TimeConverter{}},
		"*time.Time": &Wrapper[*time.Time]{Impl: &TimePtrConverter{}},
		"json":       JsonConverter{},
	},
}

// Registry maintains a mapping from a string type key to an AnyAttributeConverter.
//
// It is primarily used by the enhanced DynamoDB client to look up conversion
// strategies for Go values when serializing to / deserializing from
// DynamoDB AttributeValue types.
//
// Concurrency: A Registry is safe for concurrent read access provided no
// goroutine is mutating it. Methods that modify internal state (Add, Remove)
// are not synchronized. To customize converters at runtime without racing,
// create an independent instance using NewRegistry or Clone and mutate that
// instance before sharing it for read-only use.
//
// Keys: Converter lookup keys are the canonical type names produced by the
// internal getType helper (e.g., "int", "*int", "time.Time"). When adding
// custom converters, ensure the key matches what getType(value) would return
// for values you expect to convert.
//
// Zero value: The zero value of Registry (var r Registry) functions correctly;
// maps are lazily allocated on first Add/Converter call.
type Registry struct {
	//defaultConverter AnyAttributeConverter
	converters map[string]AnyAttributeConverter
}

// Clone creates a deep copy of the Registry's converter mapping. Converters
// themselves are not copied (the underlying converter implementations are
// assumed to be stateless or safely shareable); only the map entries are
// duplicated. The returned Registry can be mutated independently of the
// original.
func (cr *Registry) Clone() *Registry {
	r := &Registry{
		converters: make(map[string]AnyAttributeConverter),
	}

	for k, v := range cr.converters {
		r.converters[k] = v
	}

	return r
}

// Add registers (or replaces) a converter under the provided name and returns
// the Registry for fluent chaining. If the internal map has not yet been
// allocated it will be created. Replacing an existing converter is a silent
// overwrite.
func (cr *Registry) Add(name string, converter AnyAttributeConverter) *Registry {
	if cr.converters == nil {
		cr.converters = make(map[string]AnyAttributeConverter)
	}

	cr.converters[name] = converter

	return cr
}

// Remove deletes a converter by name and returns the Registry for fluent
// chaining. Removing a non-existent key is a no-op.
func (cr *Registry) Remove(name string) *Registry {
	delete(cr.converters, name)

	return cr
}

// Converter returns the converter registered for the given name. If the map
// has not yet been allocated it will be created (resulting in a nil return
// unless an entry was previously added). If no converter is found, nil is
// returned.
func (cr *Registry) Converter(name string) AnyAttributeConverter {
	if cr.converters == nil {
		cr.converters = make(map[string]AnyAttributeConverter)
	}

	return cr.converters[name]
}

// ConverterFor performs a lookup using the canonical string key derived from
// the dynamic type of x (via internal getType). If x is nil or its type has
// no registered converter, nil is returned.
func (cr *Registry) ConverterFor(x any) AnyAttributeConverter {
	return cr.Converter(getType(x))
}

// NewRegistry constructs an empty Registry with an allocated converter map.
// Use Add to populate converters or copy from DefaultRegistry via Clone.
func NewRegistry() *Registry {
	return &Registry{
		converters: make(map[string]AnyAttributeConverter),
	}
}
