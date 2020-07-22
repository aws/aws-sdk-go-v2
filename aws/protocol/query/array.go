package query

import (
	"fmt"
	"net/url"
)

// Array represents the encoding of Query lists and sets
type Array struct {
	values url.Values
	prefix string
	flat bool
	memberName string
	size int32
}

func newArray(values url.Values, prefix string, flat bool, memberName string) *Array {
	return &Array{
		values: values,
		prefix: prefix,
		flat: flat,
		memberName: memberName,
	}
}

// Value adds a new element to the Query Array. Returns a Value type used to
// encode the array element.
func (a *Array) Value() Value {
	// Query lists start a 1, so adjust the size first
	a.size++
	prefix := a.prefix
	if !a.flat {
		prefix = fmt.Sprintf("%s.%s", prefix, a.memberName)
	}
	// Lists can't have flat members
	return newValue(a.values, fmt.Sprintf("%s.%d", prefix, a.size), false)
}