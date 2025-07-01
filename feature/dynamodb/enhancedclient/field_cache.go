package enhancedclient

import (
	"reflect"
	"strings"
	"sync"
)

var fieldCache = &fieldCacher{}

type fieldCacheKey struct {
	typ  reflect.Type
	opts structFieldOptions
}

type fieldCacher struct {
	cache sync.Map
}

func (c *fieldCacher) Load(key fieldCacheKey) (*CachedFields, bool) {
	if v, ok := c.cache.Load(key); ok {
		return v.(*CachedFields), true
	}
	return nil, false
}

func (c *fieldCacher) LoadOrStore(key fieldCacheKey, fs *CachedFields) (*CachedFields, bool) {
	v, ok := c.cache.LoadOrStore(key, fs)
	return v.(*CachedFields), ok
}

type CachedFields struct {
	fields       []Field
	fieldsByName map[string]int
}

func (f *CachedFields) All() []Field {
	return f.fields
}

func (f *CachedFields) FieldByName(name string) (Field, bool) {
	if i, ok := f.fieldsByName[name]; ok {
		return f.fields[i], ok
	}
	for _, f := range f.fields {
		if strings.EqualFold(f.Name, name) {
			return f, true
		}
	}
	return Field{}, false
}
