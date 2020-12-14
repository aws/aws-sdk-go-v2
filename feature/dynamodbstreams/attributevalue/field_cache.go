package attributevalue

import (
	"strings"
	"sync"
)

var fieldCache fieldCacher

type fieldCacher struct {
	cache sync.Map
}

func (c *fieldCacher) Load(t interface{}) (*cachedFields, bool) {
	if v, ok := c.cache.Load(t); ok {
		return v.(*cachedFields), true
	}
	return nil, false
}

func (c *fieldCacher) LoadOrStore(t interface{}, fs *cachedFields) (*cachedFields, bool) {
	v, ok := c.cache.LoadOrStore(t, fs)
	return v.(*cachedFields), ok
}

type cachedFields struct {
	fields       []field
	fieldsByName map[string]int
}

func (f *cachedFields) All() []field {
	return f.fields
}

func (f *cachedFields) FieldByName(name string) (field, bool) {
	if i, ok := f.fieldsByName[name]; ok {
		return f.fields[i], ok
	}
	for _, f := range f.fields {
		if strings.EqualFold(f.Name, name) {
			return f, true
		}
	}
	return field{}, false
}
