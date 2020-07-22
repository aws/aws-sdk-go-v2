package query

import (
	"fmt"
	"net/url"
)

// Object represents the encoding of Query structures and unions.
type Object struct {
	// The query values to add the object to.
	values url.Values
	// The object's prefix, ending with the object's name.
	prefix string
}

func newObject(values url.Values, prefix string) *Object {
	return &Object{
		values: values,
		prefix: prefix,
	}
}

// Key adds the given named key to the Query object.
// Returns a Value encoder that should be used to encode a Query value type.
func (o *Object) Key(name string) Value {
	return o.key(name, false)
}

// FlatKey adds the given named key to the Query object.
// Returns a Value encoder that should be used to encode a Query value type. The
// value will be flattened if it is a map or array.
func (o *Object) FlatKey(name string) Value {
	return o.key(name, true)
}

func (o *Object) key(name string, flatValue bool) Value {
	var v Value
	if o.prefix != "" {
		v = newValue(o.values, fmt.Sprintf("%s.%s", o.prefix, name), flatValue)
	} else {
		v = newValue(o.values, name, flatValue)
	}
	return v
}
