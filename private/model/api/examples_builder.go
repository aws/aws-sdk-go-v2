// +build codegen

package api

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type examplesBuilder interface {
	BuildShape(*ShapeRef, map[string]interface{}, bool, bool) string
	BuildList(string, string, *ShapeRef, []interface{}) string
	BuildComplex(string, string, *ShapeRef, map[string]interface{}, bool) string
	Imports(*API) string
}

type defaultExamplesBuilder struct{}

// BuildShape will recursively build the referenced shape based on the json object
// provided.
// isMap will dictate how the field name is specified. If isMap is true, we will expect
// the member name to be quotes like "Foo".
func (builder defaultExamplesBuilder) BuildShape(ref *ShapeRef, shapes map[string]interface{}, isMap, parentCollection bool) string {
	order := make([]string, len(shapes))
	for k := range shapes {
		order = append(order, k)
	}
	sort.Strings(order)

	ret := ""
	for _, name := range order {
		if name == "" {
			continue
		}
		shape := shapes[name]

		// If the shape isn't a map, we want to export the value, since every field
		// defined in our shapes are exported.
		if len(name) > 0 && !isMap && strings.ToLower(name[0:1]) == name[0:1] {
			name = strings.Title(name)
		}

		memName := name
		if isMap {
			memName = fmt.Sprintf("%q", memName)
		}

		switch v := shape.(type) {
		case map[string]interface{}:
			ret += builder.BuildComplex(name, memName, ref, v, parentCollection)
		case []interface{}:
			ret += builder.BuildList(name, memName, ref, v)
		default:
			ret += builder.BuildScalar(name, memName, ref, v, parentCollection)
		}
	}
	return ret
}

// BuildList will construct a list shape based off the service's definition
// of that list.
func (builder defaultExamplesBuilder) BuildList(name, memName string, ref *ShapeRef, v []interface{}) string {
	ret := ""

	if len(v) == 0 || ref == nil {
		return ""
	}

	t := ""
	dataType := ""
	format := ""
	isComplex := false
	isEnum := false
	passRef := ref
	isMap := false

	if ref.Shape.MemberRefs[name] != nil {
		isEnum = ref.Shape.MemberRefs[name].Shape.MemberRef.Shape.IsEnum()
		t = builder.GoType(&ref.Shape.MemberRefs[name].Shape.MemberRef, true)
		dataType = ref.Shape.MemberRefs[name].Shape.MemberRef.Shape.Type
		passRef = ref.Shape.MemberRefs[name]
		if dataType == "map" {
			t = fmt.Sprintf("map[string]%s", builder.GoType(&ref.Shape.MemberRefs[name].Shape.MemberRef.Shape.ValueRef, true))
			passRef = &ref.Shape.MemberRefs[name].Shape.MemberRef.Shape.ValueRef
			isMap = true
		}
	} else if ref.Shape.MemberRef.Shape != nil && ref.Shape.MemberRef.Shape.MemberRefs[name] != nil {
		isEnum = ref.Shape.MemberRef.Shape.MemberRefs[name].Shape.MemberRef.Shape.IsEnum()
		t = builder.GoType(&ref.Shape.MemberRef.Shape.MemberRefs[name].Shape.MemberRef, true)
		dataType = ref.Shape.MemberRef.Shape.MemberRefs[name].Shape.MemberRef.Shape.Type
		passRef = &ref.Shape.MemberRef.Shape.MemberRefs[name].Shape.MemberRef
	} else {
		isEnum = ref.Shape.MemberRef.Shape.IsEnum()
		t = builder.GoType(&ref.Shape.MemberRef, true)
		dataType = ref.Shape.MemberRef.Shape.Type
		passRef = &ref.Shape.MemberRef
	}

	switch v[0].(type) {
	case string:
		format = "%s"
	case bool:
		format = "%t"
	case float64:
		if dataType == "integer" || dataType == "int64" {
			format = "%d"
		} else {
			format = "%f"
		}
	default:
		if ref.Shape.MemberRefs[name] != nil {
		} else {
			passRef = ref.Shape.MemberRef.Shape.MemberRefs[name]

			// if passRef is nil that means we are either in a map or within a nested array
			if passRef == nil {
				passRef = &ref.Shape.MemberRef
			}
		}
		isComplex = true
	}

	ret += fmt.Sprintf("%s: []%s {\n", memName, t)
	for _, elem := range v {
		if isComplex {
			ret += fmt.Sprintf("{\n%s\n},\n", builder.BuildShape(passRef, elem.(map[string]interface{}), isMap, isMap))
		} else if isEnum {
			ret += fmt.Sprintf("%s,\n", getEnumName(passRef, t, elem.(string)))
		} else {
			if dataType == "integer" || dataType == "int64" || dataType == "long" {
				elem = int(elem.(float64))
			}
			ret += fmt.Sprintf("%s,\n", getValue(t, fmt.Sprintf(format, elem), true))
		}
	}

	ret += "},\n"
	return ret
}

// BuildScalar will build atomic Go types.
func (builder defaultExamplesBuilder) BuildScalar(name, memName string, ref *ShapeRef, shape interface{}, parentCollection bool) string {
	if ref == nil || ref.Shape == nil {
		return ""
	} else if ref.Shape.MemberRefs[name] == nil {
		if ref.Shape.MemberRef.Shape != nil && ref.Shape.MemberRef.Shape.MemberRefs[name] != nil {
			if ref.Shape.MemberRef.Shape.MemberRefs[name].Shape.IsEnum() {
				refTemp := ref.Shape.MemberRef.Shape.MemberRefs[name]
				return fmt.Sprintf("%s: %s, \n", memName, getEnumName(refTemp, builder.GoType(refTemp, true), shape.(string)))
			}
			return correctType(memName, ref.Shape.MemberRef.Shape.MemberRefs[name].Shape.Type, shape, parentCollection)
		}

		if ref.Shape.Type != "structure" && ref.Shape.Type != "map" {
			if ref.Shape.IsEnum() {
				return fmt.Sprintf("%s: %s, \n", memName, getEnumName(ref, builder.GoType(ref, true), shape.(string)))
			}
			return correctType(memName, ref.Shape.Type, shape, parentCollection)
		}
		return ""
	}

	switch v := shape.(type) {
	case bool:
		return convertToCorrectType(memName, ref.Shape.MemberRefs[name].Shape.Type, fmt.Sprintf("%t", v), parentCollection)
	case int:
		if ref.Shape.MemberRefs[name].Shape.Type == "timestamp" {
			return parseTimeString(ref, memName, fmt.Sprintf("%d", v))
		}
		return convertToCorrectType(memName, ref.Shape.MemberRefs[name].Shape.Type, fmt.Sprintf("%d", v), parentCollection)
	case float64:
		dataType := ref.Shape.MemberRefs[name].Shape.Type
		if dataType == "integer" || dataType == "int64" || dataType == "long" {
			return convertToCorrectType(memName, ref.Shape.MemberRefs[name].Shape.Type, fmt.Sprintf("%d", int(shape.(float64))), parentCollection)
		}
		return convertToCorrectType(memName, ref.Shape.MemberRefs[name].Shape.Type, fmt.Sprintf("%f", v), parentCollection)
	case string:
		t := ref.Shape.MemberRefs[name].Shape.Type

		if ref.Shape.MemberRefs[name].Shape.IsEnum() {
			return fmt.Sprintf("%s: %s,\n", memName, getEnumName(ref.Shape.MemberRefs[name], builder.GoType(ref.Shape.MemberRefs[name], false), shape.(string)))
		}

		switch t {
		case "timestamp":
			return parseTimeString(ref, memName, fmt.Sprintf("%s", v))
		case "blob":
			if (ref.Shape.MemberRefs[name].Streaming || ref.Shape.MemberRefs[name].Shape.Streaming) && ref.Shape.Payload == name {
				return fmt.Sprintf("%s: aws.ReadSeekCloser(strings.NewReader(%q)),\n", memName, v)
			}

			return fmt.Sprintf("%s: []byte(%q),\n", memName, v)
		default:
			return convertToCorrectType(memName, t, v, parentCollection)
		}
	default:
		panic(fmt.Errorf("Unsupported scalar type: %v", reflect.TypeOf(v)))
	}
	return ""
}

func (builder defaultExamplesBuilder) BuildComplex(name, memName string, ref *ShapeRef, v map[string]interface{}, asValue bool) string {
	t := ""
	if ref == nil {
		return builder.BuildShape(nil, v, true, false)
	}

	member := ref.Shape.MemberRefs[name]

	if member != nil && member.Shape != nil {
		t = ref.Shape.MemberRefs[name].Shape.Type
	} else {
		t = ref.Shape.Type
	}

	switch t {
	case "structure":
		passRef := ref.Shape.MemberRefs[name]
		// passRef will be nil if the entry is a map. In that case
		// we want to pass the reference, because the previous call
		// passed the value reference.
		if passRef == nil {
			passRef = ref
		}

		mem := "&"
		if asValue {
			mem = ""
		}

		return fmt.Sprintf(`%s: %s%s{
				%s
			},
			`, memName, mem, builder.GoType(passRef, true), builder.BuildShape(passRef, v, false, false))
	case "map":
		return fmt.Sprintf(`%s: %s{
				%s
			},
			`, name, builder.GoType(ref.Shape.MemberRefs[name], true), builder.BuildShape(&ref.Shape.MemberRefs[name].Shape.ValueRef, v, true, true))
	}

	return ""
}

func (builder defaultExamplesBuilder) GoType(ref *ShapeRef, elem bool) string {
	if ref.Shape.Type == "list" {
		ref = &ref.Shape.MemberRef
	}

	if elem {
		return ref.Shape.GoTypeWithPkgNameElem()
	}

	return ref.GoTypeWithPkgName()
}

func (builder defaultExamplesBuilder) Imports(a *API) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(`"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	`)

	buf.WriteString(fmt.Sprintf("\"%s/%s\"", "github.com/aws/aws-sdk-go-v2/service", a.PackageName()))
	return buf.String()
}

func getEnumName(ref *ShapeRef, t, name string) string {
	pkg := getPkgName(ref.Shape)
	for i, enum := range ref.Shape.Enum {
		if name == enum {
			return fmt.Sprintf("%s.%s", pkg, ref.Shape.EnumConsts[i])
		}
	}

	return fmt.Sprintf("%s(%q)", t, name)
}
