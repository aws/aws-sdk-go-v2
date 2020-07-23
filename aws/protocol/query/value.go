package query

import (
	"github.com/awslabs/smithy-go/httpbinding"
	"net/url"
)

// Value represents a Query Value type.
type Value struct {
	// The query values to add the value to.
	values url.Values
	// The value's key, which will form the prefix for complex types.
	key string
	// Whether the value should be flattened or not if it's a flattenable type.
	flat bool
	httpbinding.QueryValue
}

func newValue(values url.Values, key string, flat bool) Value {
	return Value{values, key, flat, httpbinding.NewQueryValue(values, key, false)}
}

func newBaseValue(values url.Values) Value {
	return Value{values, "", false, httpbinding.NewQueryValue(nil, "", false)}
}

// Array returns a new Array encoder.
func (qv Value) Array(locationName string) *Array {
	return newArray(qv.values, qv.key, qv.flat, locationName)
}

// Object returns a new Object encoder.
func (qv Value) Object() *Object {
	return newObject(qv.values, qv.key)
}

// Map returns a new Map encoder.
func (qv Value) Map(keyLocationName string, valueLocationName string) *Map {
	return newMap(qv.values, qv.key, qv.flat, keyLocationName, valueLocationName)
}

// Base64EncodeBytes encodes v as a base64 query string value.
// This is intended to enable compatibility with the JSON encoder.
func (qv Value) Base64EncodeBytes(v []byte) {
	qv.Blob(v)
}
