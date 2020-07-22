package query

import (
	"fmt"
	"net/url"
)

// Map represents the encoding of Query maps.
type Map struct {
	// The query values to add the map to.
	values url.Values
	// The map's prefix, ending with the map's name.
	prefix string
	// Whether the map is flat or not.
	flat bool
	// The location name of the key. In most cases this should be "key".
	keyLocationName string
	// The location name of the value. In most cases this should be "value".
	valueLocationName string
	// Elements are stored in values, so we keep track of the list size here.
	size int32
}

func newMap(values url.Values, prefix string, flat bool, keyLocationName string, valueLocationName string) *Map {
	return &Map{
		values:            values,
		prefix:            prefix,
		flat:              flat,
		keyLocationName:   keyLocationName,
		valueLocationName: valueLocationName,
	}
}

// Key adds the given named key to the Query map.
// Returns a Value encoder that should be used to encode a Query value type.
func (m *Map) Key(name string) Value {
	// Query lists start a 1, so adjust the size first
	m.size++
	var key string
	var value string
	if m.flat {
		key = fmt.Sprintf("%s.%d.%s", m.prefix, m.size, m.keyLocationName)
		value = fmt.Sprintf("%s.%d.%s", m.prefix, m.size, m.valueLocationName)
	} else {
		key = fmt.Sprintf("%s.entry.%d.%s", m.prefix, m.size, m.keyLocationName)
		value = fmt.Sprintf("%s.entry.%d.%s", m.prefix, m.size, m.valueLocationName)
	}

	// The key can only be a string, so we just go ahead and set it here
	newValue(m.values, key, false).String(name)

	// Maps can't have flat members
	return newValue(m.values, value, false)
}
