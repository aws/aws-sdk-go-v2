// +build codegen

package api

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// MarshalShapeGoCode renders the shape's MarshalFields method with marshalers
// for each field within the shape. A string is returned of the rendered Go code.
//
// Will panic if error.
func MarshalShapeGoCode(s *Shape) string {
	w := &bytes.Buffer{}
	if err := protocolTmpl.ExecuteTemplate(w, "protocol shape", s); err != nil {
		panic(fmt.Sprintf("failed to render shape's fields marshaler, %v", err))
	}

	if err := marshalShapeTmpl.ExecuteTemplate(w, "encode shape", s); err != nil {
		panic(fmt.Sprintf("failed to render shape's fields marshaler, %v", err))
	}

	return w.String()
}

// MarshalShapeRefGoCode renders protocol encode for the shape ref's type.
//
// Will panic if error.
func MarshalShapeRefGoCode(refName string, ref *ShapeRef, context *Shape) string {
	if ref.XMLAttribute {
		return "// Skipping " + refName + " XML Attribute."
	}
	if context.IsRefPayloadReader(refName, ref) {
		if strings.HasSuffix(context.ShapeName, "Output") {
			return "// Skipping " + refName + " Output type's body not valid."
		}
	}

	mRef := marshalShapeRef{
		Name:    refName,
		Ref:     ref,
		Context: context,
	}

	switch mRef.Location() {
	case "StatusCode":
		return "// ignoring invalid encode state, StatusCode. " + refName
	}

	w := &bytes.Buffer{}
	if err := marshalShapeRefTmpl.ExecuteTemplate(w, "encode field", mRef); err != nil {
		panic(fmt.Sprintf("failed to marshal shape ref, %s, %v", ref.Shape.Type, err))
	}

	return w.String()
}

func isStreaming(s *Shape) bool {
	payloadName := s.Payload
	payload, ok := s.MemberRefs[payloadName]
	if !ok {
		return false
	}

	switch payload.Shape.Type {
	case "blob":
		return true
	}
	return false
}

func jsonVersion(s *Shape) string {
	if isStreaming(s) {
		return ""
	}

	return s.API.Metadata.JSONVersion
}

var marshalShapeTmpl = template.Must(template.New("marshalShapeTmpl").Funcs(
	map[string]interface{}{
		"MarshalShapeRefGoCode": MarshalShapeRefGoCode,
		"nestedRefsByLocation":  nestedRefsByLocation,
		"isShapeFieldsNested":   isShapeFieldsNested,
		"jsonVersion":           jsonVersion,
	},
).Parse(`
{{ define "encode shape" -}}
{{ $shapeName := $.ShapeName -}}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s {{ $shapeName }}) MarshalFields(e protocol.FieldEncoder) error {
	{{- if $.UsedAsInput -}}
	{{- $version := (jsonVersion $) -}}
	{{ if (ne  $version "") -}}
		e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/x-amz-json-{{ $version }}"), protocol.Metadata{})
	{{- end }}
	{{ end }}
	{{ $refMap := nestedRefsByLocation $ -}}
	{{ range $loc, $refs := $refMap -}}
		{{ $fieldsNested := isShapeFieldsNested $loc $ -}}
		{{ if $fieldsNested -}}
			e.SetFields(protocol.BodyTarget, "{{ $.LocationName }}", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		{{ end -}}
		{{ range $name, $ref := $refs -}}
			{{ MarshalShapeRefGoCode $name $ref $ }}
		{{ end -}}
		{{ if $fieldsNested -}}
			return nil
		}), {{ template "shape metadata" $ }})
		{{ end -}}
	{{ end -}}
	return nil
}

{{- end }}

{{ define "shape metadata" -}}
	protocol.Metadata{
		{{- if $.XMLNamespace.URI -}}
			XMLNamespaceURI: "{{ $.XMLNamespace.URI }}",
		{{- end -}}
	}
{{- end }}
`))

// marshalEnumTmpl will generate MarshalValue and MarshalValueBuf methods
// to satisfy our value interface
var marshalEnumTmpl = template.Must(template.New("marshalShapeTmpl").Funcs(
	map[string]interface{}{},
).Parse(`
{{ define "enum defs" -}}
func (enum {{ $.GoType }}) MarshalValue() (string, error) {
	return string(enum), nil
}

func (enum {{ $.GoType }}) MarshalValueBuf(b []byte) ([]byte, error) {
	b = b[0:0]
	return append(b, enum...), nil
}

{{- end }}
{{ define "list enum defs" -}}
{{ if $.UsedInMap -}}
func encode{{ $.EnumType }}Map(vs map[string]{{ $.EnumType }}) func(protocol.MapEncoder) {
	return func(me protocol.MapEncoder) {
		for k, v := range vs {
			me.MapSetFields(k, v)
		}
	}
}
{{ end -}}
func encode{{ $.EnumType }}List(vs []{{ $.EnumType }}) func(protocol.ListEncoder) {
	return func(le protocol.ListEncoder) {
		for _, v := range vs {
			le.ListAddValue(v)
		}
	}
}
{{- end }}
`))

func marshalEnumGoCode(s *Shape) string {
	w := &bytes.Buffer{}
	if err := marshalEnumTmpl.ExecuteTemplate(w, "enum defs", s); err != nil {
		panic(fmt.Sprintf("failed to render shape's fields marshaler, %v", err))
	}

	return w.String()
}

func nestedRefsByLocation(s *Shape) map[string]map[string]*ShapeRef {
	refs := map[string]map[string]*ShapeRef{}

	for refName, ref := range s.MemberRefs {
		mRef := marshalShapeRef{
			Name:    refName,
			Ref:     ref,
			Context: s,
		}

		loc := mRef.Location()

		var rs map[string]*ShapeRef
		var ok bool
		if rs, ok = refs[loc]; !ok {
			rs = map[string]*ShapeRef{}
		}
		rs[refName] = ref

		refs[loc] = rs
	}

	return refs
}
func isShapeFieldsNested(loc string, s *Shape) bool {
	return loc == "Body" && len(s.LocationName) != 0 && s.API.Metadata.Protocol == "rest-xml"
}

var marshalShapeRefTmpl = template.Must(template.New("marshalShapeRefTmpl").Funcs(template.FuncMap{
	"Collection": Collection,
}).Parse(`
{{ define "encode field" -}}
	{{ if $.IsIdempotencyToken -}}
		{{ template "idempotency token" $ }}
	{
		v := {{ $.Name }}
	{{ else -}}
	if {{ template "is ref set" $ }} {
		v := {{ template "ref value" $ }}
	{{- end }}
		{{ if $.HasAttributes -}}
			{{ template "attributes" $ }}
		{{- end }}
		metadata := {{ template "metadata" $ }}
		{{ if $.IsCollection -}}
		{{ Collection $ 0 }}
		{{ else -}}
		e.Set{{ $.MarshalerType }}(protocol.{{ $.Location }}Target, "{{ $.LocationName }}", {{ template "marshaler" $ }}, metadata)
		{{ end -}}
	}
{{- end }}

{{ define "marshaler" -}}
	{{- if $.IsShapeType "list" -}}
		{{- $helperName := $.EncodeHelperName "list" -}}
			{{ $helperName }}
	{{- else if $.IsShapeType "map" -}}
		{{- $helperName := $.EncodeHelperName "map" -}}
			{{ $helperName }}
	{{- else if $.Ref.Shape.IsEnum -}}
	{{ if $.Quoted }}protocol.QuotedValue{ValueMarshaler: {{ end }}v{{ if $.Quoted }} }  {{ end }}
	{{- else if $.IsShapeType "structure" -}}
		v
	{{- else if $.IsShapeType "timestamp" -}}
		protocol.TimeValue{V: v, Format: {{ $.TimeFormat }} }
	{{- else if $.IsShapeType "jsonvalue" -}}
		{{ if eq $.Location "Header" -}}
		protocol.JSONValue{V: v , EscapeMode: protocol.Base64Escape}
		{{- else if eq $.Location "Body" -}}
		protocol.JSONValue{V: v, EscapeMode: protocol.QuotedEscape}
		{{- else -}}
		protocol.JSONValue{V: v}
		{{- end }}
	{{- else if $.IsPayloadStream -}}
		protocol.{{ $.GoType }}{{ $.MarshalerType }}{V:v}
	{{- else -}}
	{{ if $.Quoted }}protocol.QuotedValue{ValueMarshaler: {{ end }}protocol.{{ $.GoType }}{{ $.MarshalerType }}(v){{ if $.Quoted }} } {{ end }}
	{{- end -}}
{{- end }}

{{ define "metadata" -}}
	protocol.Metadata{
		{{- if $.IsFlattened -}}
			Flatten: true,
		{{- end -}}

		{{- if $.HasAttributes -}}
			Attributes: attrs,
		{{- end -}}

		{{- if $.XMLNamespacePrefix -}}
			XMLNamespacePrefix: "{{ $.XMLNamespacePrefix }}",
		{{- end -}}

		{{- if $.XMLNamespaceURI -}}
			XMLNamespaceURI: "{{ $.XMLNamespaceURI }}",
		{{- end -}}

		{{- if $.ListLocationName -}}
			ListLocationName: "{{ $.ListLocationName }}",
		{{- end -}}

		{{- if $.MapLocationNameKey -}}
			MapLocationNameKey: "{{ $.MapLocationNameKey }}",
		{{- end -}}

		{{- if $.MapLocationNameValue -}}
			MapLocationNameValue: "{{ $.MapLocationNameValue }}",
		{{- end -}}
	}
{{- end }}

{{ define "ref value" -}}
	{{ if $.Ref.UseIndirection }}*{{ end }}s.{{ $.Name }}
{{- end}}

{{ define "is ref set" -}}
	{{ $isList := $.IsShapeType "list" -}}
	{{ $isMap := $.IsShapeType "map" -}}
	{{- if or $isList $isMap -}}
		len(s.{{ $.Name }}) > 0
	{{- else if $.Ref.Shape.IsEnum -}}
		len(s.{{ $.Name }}) > 0
	{{- else -}}
		s.{{ $.Name }} != nil
	{{- end -}}
{{- end }}

{{ define "attributes" -}}
	attrs := make([]protocol.Attribute, 0, {{ $.NumAttributes }})

	{{ range $name, $child := $.ChildrenRefs -}}
		{{ if $child.Ref.XMLAttribute -}}
			{{ if $child.Ref.Shape.IsEnum -}}
			if len(s.{{ $.Name }}.{{ $name }}) > 0 {
			{{ else -}}
			if s.{{ $.Name }}.{{ $name }} != nil {
			{{- end }}
				v := {{ if $child.Ref.UseIndirection }}*{{ end }}s.{{ $.Name }}.{{ $name }}
				attrs = append(attrs, protocol.Attribute{Name: "{{ $child.LocationName }}", Value: {{ template "marshaler" $child }}, Meta: {{ template "metadata" $child }}})
			}
		{{- end }}
	{{- end }}
{{- end }}

{{ define "idempotency token" -}}
    var {{ $.Name }} string
	if {{ template "is ref set" $ }} {
		{{ $.Name }} = {{ template "ref value" $ }}
	} else {
		{{ $.Name }} = protocol.GetIdempotencyToken()
	}
{{- end }}
`))

type marshalShapeRef struct {
	Name    string
	Ref     *ShapeRef
	Context *Shape
}

func (r marshalShapeRef) Quoted() bool {
	if r.Context.IsRefPayload(r.Name) {
		return false
	}
	if r.Context.IsRefPayloadReader(r.Name, r.Ref) {
		return false
	}

	if r.Ref.API != nil {
		switch r.Ref.API.ProtocolPackage() {
		case "jsonrpc", "restjson":
			return r.Ref.Shape.Type == "string" || r.Ref.Shape.Type == "blob"
		}
	}

	return false
}

func (r marshalShapeRef) IsCollection() bool {
	switch r.Ref.Shape.Type {
	case "list", "map":
		return true
	default:
		return false
	}
}

func Collection(ref marshalShapeRef, level int) string {
	buf := bytes.Buffer{}
	parentVar := getParentVariableName(ref.Context, level-1)

	switch ref.Ref.Shape.Type {
	case "list":
		if level > 0 {
			buf.WriteString(fmt.Sprintf("ls%d := %s.%s\n", level, parentVar, startList(ref.Context, level)))
			buf.WriteString(fmt.Sprintf("ls%d.Start()\n", level))
			buf.WriteString(fmt.Sprintf("for _, v%d := range v%d {\n", level+1, level))
		} else {
			buf.WriteString(fmt.Sprintf("ls%d := e.List(protocol.%sTarget, %q, metadata)\n", level, ref.Location(), ref.LocationName()))
			buf.WriteString(fmt.Sprintf("ls%d.Start()\n", level))
			buf.WriteString(fmt.Sprintf("for _, v%d := range v {\n", level+1))
		}
		nested := marshalShapeRef{
			Name:    "nested",
			Ref:     &ref.Ref.Shape.MemberRef,
			Context: ref.Ref.Shape,
		}
		buf.WriteString(Collection(nested, level+1))
		buf.WriteString("}\n")
		buf.WriteString(fmt.Sprintf("ls%d.End()\n", level))
	case "map":
		if level > 0 {
			buf.WriteString(fmt.Sprintf("ms%d := %s.%s\n", level, parentVar, startMap(ref.Context, level)))
			buf.WriteString(fmt.Sprintf("ms%d.Start()\n", level))
			buf.WriteString(fmt.Sprintf("for k%d, v%d := range v%d {\n", level+1, level+1, level))
		} else {
			buf.WriteString(fmt.Sprintf("ms%d := e.Map(protocol.%sTarget, %q, metadata)\n", level, ref.Location(), ref.LocationName()))
			buf.WriteString(fmt.Sprintf("ms%d.Start()\n", level))
			buf.WriteString(fmt.Sprintf("for k%d, v%d := range v {\n", level+1, level+1))
		}

		nested := marshalShapeRef{
			Name:    "nested",
			Ref:     &ref.Ref.Shape.ValueRef,
			Context: ref.Ref.Shape,
		}
		buf.WriteString(Collection(nested, level+1))
		buf.WriteString("}\n")
		buf.WriteString(fmt.Sprintf("ms%d.End()\n", level))
	default:
		switch ref.Context.Type {
		case "list":
			if ref.Ref.Shape.Type == "structure" {
				buf.WriteString(fmt.Sprintf("%s.ListAddFields(v%d)\n", parentVar, level))
			} else if ref.Ref.Shape.Type == "jsonvalue" {
				buf.WriteString(fmt.Sprintf("%s.ListAddValue(protocol.JSONValue{V: v%d, EscapeMode: protocol.QuotedEscape})\n", parentVar, level))
			} else if ref.Ref.Shape.Type == "timestamp" {
				buf.WriteString(fmt.Sprintf("%s.ListAddValue(protocol.TimeValue{V: v%d})\n", parentVar, level))
			} else {
				if ref.Quoted() {
					buf.WriteString(fmt.Sprintf("%s.ListAddValue(protocol.QuotedValue{ValueMarshaler: protocol.%s%s(v%d)})\n", parentVar, ref.GoType(), ref.MarshalerType(), level))
				} else {
					buf.WriteString(fmt.Sprintf("%s.ListAddValue(protocol.%s%s(v%d))\n", parentVar, ref.GoType(), ref.MarshalerType(), level))
				}
			}
		case "map":
			if ref.Ref.Shape.Type == "structure" {
				buf.WriteString(fmt.Sprintf("%s.MapSetFields(k%d, v%d)\n", parentVar, level, level))
			} else if ref.Ref.Shape.Type == "jsonvalue" {
				buf.WriteString(fmt.Sprintf("%s.MapSetValue(k%d, protocol.JSONValue{V: v%d, EscapeMode: protocol.QuotedEscape})\n", parentVar, level, level))
			} else if ref.Ref.Shape.Type == "timestamp" {
				buf.WriteString(fmt.Sprintf("%s.MapSetValue(k%d, protocol.TimeValue{V: v%d})\n", parentVar, level, level))
			} else {
				if ref.Quoted() {
					buf.WriteString(fmt.Sprintf("%s.MapSetValue(k%d, protocol.QuotedValue{ValueMarshaler: protocol.%s%s(v%d)})\n", parentVar, level, ref.GoType(), ref.MarshalerType(), level))
				} else {
					buf.WriteString(fmt.Sprintf("%s.MapSetValue(k%d, protocol.%s%s(v%d))\n", parentVar, level, ref.GoType(), ref.MarshalerType(), level))
				}
			}
		}
	}
	return buf.String()
}

func startList(shape *Shape, level int) string {
	if shape.Type == "map" {
		return fmt.Sprintf("List(k%d)", level)
	}
	return "List()"
}

func startMap(shape *Shape, level int) string {
	if shape.Type == "map" {
		return fmt.Sprintf("Map(k%d)", level)
	}
	return "Map()"
}

func getParentVariableName(shape *Shape, level int) string {
	if shape.Type == "list" {
		return fmt.Sprintf("ls%d", level)
	} else if shape.Type == "map" {
		return fmt.Sprintf("ms%d", level)
	}
	return ""
}

func (r marshalShapeRef) ListMemberRef() marshalShapeRef {
	return marshalShapeRef{
		Name:    r.Name + "ListMember",
		Ref:     &r.Ref.Shape.MemberRef,
		Context: r.Ref.Shape,
	}
}
func (r marshalShapeRef) MapValueRef() marshalShapeRef {
	return marshalShapeRef{
		Name:    r.Name + "MapValue",
		Ref:     &r.Ref.Shape.ValueRef,
		Context: r.Ref.Shape,
	}
}

func (r marshalShapeRef) MapKeyRef() marshalShapeRef {
	return marshalShapeRef{
		Name:    r.Name + "MapKey",
		Ref:     &r.Ref.Shape.KeyRef,
		Context: r.Ref.Shape,
	}
}

func (r marshalShapeRef) ChildrenRefs() map[string]marshalShapeRef {
	children := map[string]marshalShapeRef{}

	for name, ref := range r.Ref.Shape.MemberRefs {
		children[name] = marshalShapeRef{
			Name:    name,
			Ref:     ref,
			Context: r.Ref.Shape,
		}
	}

	return children
}
func (r marshalShapeRef) IsShapeType(typ string) bool {
	return r.Ref.Shape.Type == typ
}
func (r marshalShapeRef) IsPayloadStream() bool {
	return r.Context.IsRefPayloadReader(r.Name, r.Ref)
}
func (r marshalShapeRef) MarshalerType() string {
	switch r.Ref.Shape.Type {
	case "list":
		return "List"
	case "map":
		return "Map"
	case "structure":
		return "Fields"
	default:
		// Streams have a special case
		if r.Context.IsRefPayload(r.Name) {
			return "Stream"
		}
		return "Value"
	}
}

func getSuffixName(r marshalShapeRef) string {
	ref := r.MapKeyRef()

	switch ref.Ref.Shape.Type {
	case "blob":
		return "Blob"
	default:
		return ref.GoType()
	}
}

func (r marshalShapeRef) EncodeHelperName(typ string) string {
	if r.Ref.Shape.Type != typ {
		panic(fmt.Sprintf("ref shape %q does not match type %q", r.Ref.Shape.Type, typ))
	}

	prefixEncodeName := ""
	suffixEncodeName := ""
	var memberRef marshalShapeRef
	switch r.Ref.Shape.Type {
	case "map":
		memberRef = r.MapValueRef()
		prefixEncodeName = "EncodeMap"
		suffixEncodeName = getSuffixName(r)
	case "list":
		memberRef = r.ListMemberRef()
		prefixEncodeName = "EncodeList"
	default:
		panic("invalid")
	}

	if memberRef.Ref.Shape.IsEnum() {
		return "encode" + memberRef.Ref.Shape.EnumType() + "Enums" + "(v)"
	}

	switch memberRef.Ref.Shape.Type {
	case "list":
		return ""
	case "map":
		return ""
	case "jsonvalue":
		return ""
	case "structure":
		shapeName := memberRef.Ref.Shape.ShapeName
		return "encode" + shapeName + strings.Title(typ) + "(v)"
	case "integer", "long":
		return fmt.Sprintf("protocol.%sInt64%s(v)", prefixEncodeName, suffixEncodeName)
	case "float", "double":
		return fmt.Sprintf("protocol.%sFloat64%s(v)", prefixEncodeName, suffixEncodeName)
	case "string":
		return fmt.Sprintf("protocol.%sString%s(v)", prefixEncodeName, suffixEncodeName)
	case "blob":
		return fmt.Sprintf("protocol.%sBlob%s(v)", prefixEncodeName, suffixEncodeName)
	case "timestamp":
		return fmt.Sprintf("protocol.%Time%s(v)", prefixEncodeName, suffixEncodeName)
	default:
		return "protocol.Encode" + memberRef.GoType() + strings.Title(typ) + "(v)"
	}
}

func (r marshalShapeRef) GoType() string {
	switch r.Ref.Shape.Type {
	case "boolean":
		return "Bool"
	case "string", "character":
		return "String"
	case "integer", "long":
		return "Int64"
	case "float", "double":
		return "Float64"
	case "timestamp":
		return "Time"
	case "jsonvalue":
		return "JSONValue"
	case "blob":
		if r.Context.IsRefPayloadReader(r.Name, r.Ref) {
			if strings.HasSuffix(r.Context.ShapeName, "Output") {
				return "ReadCloser"
			}
			return "ReadSeeker"
		}
		return "Bytes"
	default:
		panic(fmt.Sprintf("unknown marshal shape ref type, %s", r.Ref.Shape.Type))
	}
}
func (r marshalShapeRef) Location() string {
	var loc string
	if l := r.Ref.Location; len(l) > 0 {
		loc = l
	} else if l := r.Ref.Shape.Location; len(l) > 0 {
		loc = l
	}

	switch loc {
	case "querystring":
		return "Query"
	case "header":
		return "Header"
	case "headers": // headers means key is header prefix
		return "Headers"
	case "uri":
		return "Path"
	case "statusCode":
		return "StatusCode"
	default:
		if len(loc) != 0 {
			panic(fmt.Sprintf("unknown marshal shape ref location, %s", loc))
		}

		if r.Context.IsRefPayload(r.Name) {
			return "Payload"
		}

		return "Body"
	}
}

func (r marshalShapeRef) LocationName() string {
	if l := r.Ref.QueryName; len(l) > 0 {
		// Special case for EC2 query
		return l
	}

	locName := r.Name
	if l := r.Ref.LocationName; len(l) > 0 {
		locName = l
	} else if l := r.Ref.Shape.LocationName; len(l) > 0 {
		locName = l
	}

	return locName
}

func (r marshalShapeRef) IsFlattened() bool {
	return r.Ref.Flattened || r.Ref.Shape.Flattened
}

func (r marshalShapeRef) XMLNamespacePrefix() string {
	if v := r.Ref.XMLNamespace.Prefix; len(v) != 0 {
		return v
	}
	return r.Ref.Shape.XMLNamespace.Prefix
}

func (r marshalShapeRef) XMLNamespaceURI() string {
	if v := r.Ref.XMLNamespace.URI; len(v) != 0 {
		return v
	}
	return r.Ref.Shape.XMLNamespace.URI
}

func (r marshalShapeRef) ListLocationName() string {
	if v := r.Ref.Shape.MemberRef.LocationName; len(v) > 0 {
		if !(r.Ref.Shape.Flattened || r.Ref.Flattened) {
			return v
		}
	}
	return ""
}

func (r marshalShapeRef) MapLocationNameKey() string {
	return r.Ref.Shape.KeyRef.LocationName
}
func (r marshalShapeRef) MapLocationNameValue() string {
	return r.Ref.Shape.ValueRef.LocationName
}
func (r marshalShapeRef) HasAttributes() bool {
	for _, ref := range r.Ref.Shape.MemberRefs {
		if ref.XMLAttribute {
			return true
		}
	}
	return false
}
func (r marshalShapeRef) NumAttributes() (n int) {
	for _, ref := range r.Ref.Shape.MemberRefs {
		if ref.XMLAttribute {
			n++
		}
	}
	return n
}
func (r marshalShapeRef) IsIdempotencyToken() bool {
	return r.Ref.IdempotencyToken || r.Ref.Shape.IdempotencyToken
}
func (r marshalShapeRef) TimeFormat() string {
	switch r.Location() {
	case "Header", "Headers":
		return "protocol.RFC822TimeFromat"
	case "Query":
		return "protocol.RFC822TimeFromat"
	default:
		switch r.Context.API.Metadata.Protocol {
		case "json", "rest-json":
			return "protocol.UnixTimeFormat"
		case "rest-xml", "ec2", "query":
			return "protocol.ISO8601TimeFormat"
		default:
			panic(fmt.Sprintf("unable to determine time format for %s ref", r.Name))
		}
	}
}

// UnmarshalShapeGoCode renders the shape's UnmarshalAWS method with unmarshalers
// for each field within the shape. A string is returned of the rendered Go code.
//
// Will panic if error.
func UnmarshalShapeGoCode(s *Shape) string {
	w := &bytes.Buffer{}
	if err := unmarshalShapeTmpl.Execute(w, s); err != nil {
		panic(fmt.Sprintf("failed to render shape's fields unmarshaler, %v", err))
	}

	return w.String()
}

var unmarshalShapeTmpl = template.Must(template.New("unmarshalShapeTmpl").Funcs(
	template.FuncMap{
		"UnmarshalShapeRefGoCode": UnmarshalShapeRefGoCode,
	},
).Parse(`
{{ $shapeName := $.ShapeName -}}

// UnmarshalAWS decodes the AWS API shape using the passed in protocol decoder.
func (s *{{ $shapeName }}) UnmarshalAWS(d protocol.FieldDecoder) {
	{{ range $name, $ref := $.MemberRefs -}}
		{{ UnmarshalShapeRefGoCode $name $ref $ }}
	{{ end }}
}

{{ if $.UsedInList -}}
func decode{{ $shapeName }}List(vsp *[]*{{ $shapeName }}) func(int, protocol.ListDecoder) {
	return func(n int, ld protocol.ListDecoder) {
		vs := make([]{{ $shapeName }}, n)
		*vsp = make([]*{{ $shapeName }}, n)
		for i := 0; i < n; i++ {
			ld.ListGetUnmarshaler(&vs[i])
			(*vsp)[i] = &vs[i]
		}
	}
}
{{- end }}

{{ if $.UsedInMap -}}
func decode{{ $shapeName }}Map(vsp *map[string]*{{ $shapeName }}) func([]string, protocol.MapDecoder) {
	return func(ks []string, md protocol.MapDecoder) {
		vs := make(map[string]*{{ $shapeName }}, n)
		for _, k range ks {
			v := &{{ $shapeName }}{}
			md.MapGetUnmarshaler(k, v)
			vs[k] = v
		}
	}
}
{{- end }}
`))

// UnmarshalShapeRefGoCode generates the Go code to unmarshal an API shape.
func UnmarshalShapeRefGoCode(refName string, ref *ShapeRef, context *Shape) string {
	if ref.XMLAttribute {
		return "// Skipping " + refName + " XML Attribute."
	}

	mRef := marshalShapeRef{
		Name:    refName,
		Ref:     ref,
		Context: context,
	}

	switch mRef.Location() {
	case "Path":
		return "// ignoring invalid decode state, Path. " + refName
	case "Query":
		return "// ignoring invalid decode state, Query. " + refName
	}

	w := &bytes.Buffer{}
	if err := unmarshalShapeRefTmpl.ExecuteTemplate(w, "decode", mRef); err != nil {
		panic(fmt.Sprintf("failed to decode shape ref, %s, %v", ref.Shape.Type, err))
	}

	return w.String()
}

var unmarshalShapeRefTmpl = template.Must(template.New("unmarshalShapeRefTmpl").Parse(`
//  Decode {{ $.Name }} {{ $.GoType }} {{ $.MarshalerType }} to {{ $.Location }} at {{ $.LocationName }}
`))

func CanGenerateProtocolShapes(s *Shape) bool {
	if s.Type == "list" {
		return true
	}
	return false
}

var protocolTmpl = template.Must(template.New("protocolTmpl").Funcs(
	map[string]interface{}{
		"CanGenerateProtocolShapes": CanGenerateProtocolShapes,
	},
).Parse(`
{{ define "protocol shape" -}}
{{ if CanGenerateProtocolShapes $ -}}
// valueMarshaler{{ .ShapeName }} ...
type valueMarshaler{{ .ShapeName }} struct {
	V {{ .ShapeName }}
}

func (p *valueMarshaler{{ .ShapeName }}) MarshalValue() (string, error) {
	return string(p.V), nil
}
{{ end -}}
{{- end }}
`))
