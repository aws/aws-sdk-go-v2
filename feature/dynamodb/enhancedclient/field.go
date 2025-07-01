package enhancedclient

import (
	"reflect"
	"sort"
)

type Field struct {
	Tag

	Name        string
	NameFromTag bool

	Index []int
	Type  reflect.Type
}

func buildField(pIdx []int, i int, sf reflect.StructField, fieldTag Tag) Field {
	f := Field{
		Name: sf.Name,
		Type: sf.Type,
		Tag:  fieldTag,
	}
	if len(fieldTag.Name) != 0 {
		f.NameFromTag = true
		f.Name = fieldTag.Name
	}

	f.Index = make([]int, len(pIdx)+1)
	copy(f.Index, pIdx)
	f.Index[len(pIdx)] = i

	return f
}

type structFieldOptions struct {
	// Support other custom struct Tag keys, such as `yaml`, `json`, or `toml`.
	// Note that values provided with a custom TagKey must also be supported
	// by the (un)marshalers in this package.
	//
	// Tag key `dynamodbav` will always be read, but if custom Tag key
	// conflicts with `dynamodbav` the custom Tag key value will be used.
	TagKey string
}

// unionStructFields returns a list of CachedFields for the given type. Type info is cached
// to avoid repeated calls into the reflect package
func unionStructFields(t reflect.Type, opts structFieldOptions) *CachedFields {
	key := fieldCacheKey{
		typ:  t,
		opts: opts,
	}

	if cached, ok := fieldCache.Load(key); ok {
		return cached
	}

	f := enumFields(t, opts)
	sort.Sort(fieldsByName(f))
	f = visibleFields(f)

	fs := &CachedFields{
		fields:       f,
		fieldsByName: make(map[string]int, len(f)),
	}
	for i, f := range fs.fields {
		fs.fieldsByName[f.Name] = i
	}

	cached, _ := fieldCache.LoadOrStore(key, fs)
	return cached
}

// enumFields will recursively iterate through a structure and its nested
// anonymous CachedFields.
//
// Based on the enoding/json struct Field enumeration of the Go Stdlib
// https://golang.org/src/encoding/json/encode.go typeField func.
func enumFields(t reflect.Type, opts structFieldOptions) []Field {
	// Fields to explore
	current := []Field{}
	next := []Field{{Type: t}}

	// count of queued names
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	visited := map[reflect.Type]struct{}{}
	fields := []Field{}

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if _, ok := visited[f.Type]; ok {
				continue
			}
			visited[f.Type] = struct{}{}

			for i := 0; i < f.Type.NumField(); i++ {
				sf := f.Type.Field(i)

				fieldTag := Tag{}
				fieldTag.parseAVTag(sf.Tag)
				// Because MarshalOptions.TagKey must be explicitly set.
				if opts.TagKey != "" && opts.TagKey != defaultTagKey {
					fieldTag.parseStructTag(opts.TagKey, sf.Tag)
				}

				if sf.PkgPath != "" && !sf.Anonymous && fieldTag.Getter == "" && fieldTag.Setter == "" {
					// Ignore unexported and non-anonymous CachedFields
					// unexported but anonymous Field may still be used if
					// the type has exported nested CachedFields
					// or if they have a getter and setter
					continue
				}

				if fieldTag.Ignore {
					continue
				}

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}

				structField := buildField(f.Index, i, sf, fieldTag)
				structField.Type = ft

				if !sf.Anonymous || fieldTag.Name != "" || ft.Kind() != reflect.Struct {
					fields = append(fields, structField)
					if count[f.Type] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, structField)
					}
					continue
				}

				// Record new anon struct to explore next round
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, structField)
				}
			}
		}
	}

	return fields
}

// visibleFields will return a slice of CachedFields which are visible based on
// Go's standard visiblity rules with the exception of ties being broken
// by depth and struct Tag naming.
//
// Based on the enoding/json Field filtering of the Go Stdlib
// https://golang.org/src/encoding/json/encode.go typeField func.
func visibleFields(fields []Field) []Field {
	// Delete all CachedFields that are hidden by the Go rules for embedded CachedFields,
	// except that CachedFields with JSON tags are promoted.

	// The CachedFields are sorted in primary order of name, secondary order
	// of Field index length. Loop over names; for each name, delete
	// hidden CachedFields by choosing the one dominant Field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per name.
		// Find the sequence of CachedFields with the name of this first Field.
		fi := fields[i]
		name := fi.Name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.Name != name {
				break
			}
		}
		if advance == 1 { // Only one Field with this name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(fieldsByIndex(fields))

	return fields
}

// dominantField looks through the CachedFields, all of which are known to
// have the same name, to find the single Field that dominates the
// others using Go's embedding rules, modified by the presence of
// JSON tags. If there are multiple top-level CachedFields, the boolean
// will be false: This condition is an error in Go and we skip all
// the CachedFields.
//
// Based on the enoding/json Field filtering of the Go Stdlib
// https://golang.org/src/encoding/json/encode.go dominantField func.
func dominantField(fields []Field) (Field, bool) {
	// The CachedFields are sorted in increasing index-length order. The winner
	// must therefore be one with the shortest index length. Drop all
	// longer entries, which is easy: just truncate the slice.
	length := len(fields[0].Index)
	tagged := -1 // Index of first tagged Field.
	for i, f := range fields {
		if len(f.Index) > length {
			fields = fields[:i]
			break
		}
		if f.NameFromTag {
			if tagged >= 0 {
				// Multiple tagged CachedFields at the same level: conflict.
				// Return no Field.
				return Field{}, false
			}
			tagged = i
		}
	}
	if tagged >= 0 {
		return fields[tagged], true
	}
	// All remaining CachedFields have the same length. If there's more than one,
	// we have a conflict (two CachedFields named "X" at the same level) and we
	// return no Field.
	if len(fields) > 1 {
		return Field{}, false
	}
	return fields[0], true
}

// fieldsByName sorts Field by name, breaking ties with depth,
// then breaking ties with "name came from json Tag", then
// breaking ties with index sequence.
//
// Based on the enoding/json Field filtering of the Go Stdlib
// https://golang.org/src/encoding/json/encode.go fieldsByName type.
type fieldsByName []Field

func (x fieldsByName) Len() int { return len(x) }

func (x fieldsByName) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x fieldsByName) Less(i, j int) bool {
	if x[i].Name != x[j].Name {
		return x[i].Name < x[j].Name
	}
	if len(x[i].Index) != len(x[j].Index) {
		return len(x[i].Index) < len(x[j].Index)
	}
	if x[i].NameFromTag != x[j].NameFromTag {
		return x[i].NameFromTag
	}
	return fieldsByIndex(x).Less(i, j)
}

// fieldsByIndex sorts Field by index sequence.
//
// Based on the enoding/json Field filtering of the Go Stdlib
// https://golang.org/src/encoding/json/encode.go fieldsByIndex type.
type fieldsByIndex []Field

func (x fieldsByIndex) Len() int { return len(x) }

func (x fieldsByIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x fieldsByIndex) Less(i, j int) bool {
	for k, xik := range x[i].Index {
		if k >= len(x[j].Index) {
			return false
		}
		if xik != x[j].Index[k] {
			return xik < x[j].Index[k]
		}
	}
	return len(x[i].Index) < len(x[j].Index)
}
