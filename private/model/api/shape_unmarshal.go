// +build codegen

package api

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func UnmarshalShapeGoCode(s *Shape) string {
	w := &bytes.Buffer{}
	if s.API.Metadata.Protocol == "rest-xml" {
		if err := unmarshalShapeXMLTmpl.ExecuteTemplate(w, "main body", s); err != nil {
			panic(fmt.Sprintf("failed to render shape's fields XML unmarshaler, %v", err))
		}
		if err := unmarshalShapeRESTTmpl.Execute(w, s); err != nil {
			panic(fmt.Sprintf("failed to render shape's fields REST unmarshaler, %v", err))
		}
		if err := unmarshalShapePayloadTmpl.Execute(w, s); err != nil {
			panic(fmt.Sprintf("failed to render shape's fields Payload unmarshaler, %v", err))
		}
	}
	if s.API.Metadata.Protocol == "rest-json" {
		// TO DO
	}
	return w.String()
}

var unmarshalShapeXMLTmpl = template.Must(template.New("unmarshalShapeXMLTmpl").Funcs(
	template.FuncMap{
		"setImports": setImports,
		"setSDKImports": setSDKImports,
		"getStructFieldOuterTagName": getStructFieldOuterTagName,
		"findLocationNameLocal": findLocationNameLocal,
		"templateMap": templateMap,
	},
).Parse(`
{{ define "main body" -}}
{{ $shapeName := $.ShapeName -}}
{{ if $.NeedUnmarshalXMLEntry -}}
{{ setImports $ "encoding/xml" -}}

// UnmarshalAWSXML decodes the AWS API shape using the passed in *xml.Decoder.
func (s *{{ $.GoTypeElem }}) UnmarshalAWSXML(d *xml.Decoder) (err error) {
	defer func() {
		if err != nil {
			*s = {{ $.GoTypeElem }}{}
		}
	}()
	{{ setImports $ "fmt" -}}
	{{ setImports $ "io" -}}
	{{ $payloadRef := $.GetPayloadShapeRef -}}
	for {
		tok, err := d.Token()
		if tok == nil || err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}, %s", err)
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		{{ if $payloadRef -}}
			if s.{{ $.Payload }} == nil {
				s.{{ $.Payload }} = &{{ $payloadRef.GoTypeElem }}{}
			}
		{{ end -}}
		err = s.{{ if $payloadRef }}{{ $.Payload }}.{{ end }}unmarshalAWSXML(d, start)
		if err != nil {
			return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}, %s", err)
		}
		return nil
	}	
}

{{ end -}}

{{ if $.NeedUnmarshalXMLStructHelper -}}
{{ setImports $ "encoding/xml" -}}

func (s *{{ $.GoTypeElem }}) unmarshalAWSXML(d *xml.Decoder, head xml.StartElement) (err error) {
	defer func() {
		if err != nil {
			*s = {{ $.GoTypeElem }}{}
		}
	}()
	{{ $xmlAttribute := $.GetXMLAttributeField -}}
	{{ if $xmlAttribute -}}
		{{ $xmlRef := index $.MemberRefs $xmlAttribute -}}
		for index, attr := range head.Attr {
			if attr.Name.Local == "{{ findLocationNameLocal $xmlRef.GetLocationName }}" {
				v := head.Attr[index].Value
				{{ template "parse scalar" (templateMap "shapeRef" $xmlRef) -}}
				s.{{ $xmlAttribute }} = {{ if and (not (eq $xmlRef.Shape.Type "blob")) (not $xmlRef.Shape.IsEnum) -}}&{{ end }}value
			}
		}
	{{ end -}}
	name := ""
	for {
		tok, err := d.Token()
		if tok == nil || err != nil {
			{{ if $.IsInputOrOutput -}}
				{{ setImports $ "io" -}}
				{{ setImports $ "fmt" -}}
				if err == io.EOF {
					return nil
				}
			{{ else -}}
				return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
			{{ end -}}
		}
		if end, ok := tok.(xml.EndElement); ok {
			name = end.Name.Local
			if name == head.Name.Local {
				return nil
			}
		}
		if start, ok := tok.(xml.StartElement); ok {
			switch name = start.Name.Local; name {
			{{ range $name, $ref := $.MemberRefs -}}
				{{ if and (not (eq $name $xmlAttribute )) (getStructFieldOuterTagName $ref $name) -}}
					case "{{ getStructFieldOuterTagName $ref $name }}":
					{{ if eq $ref.Shape.Type "list" -}}
						if s.{{ $name }} == nil {
							s.{{ $name }} = make({{ $ref.GoTypeElem }}, 0)
						}
						{{ if or $ref.Flattened $ref.Shape.Flattened -}}
							{{ if not $ref.Shape.MemberRef.Shape.NestedShape -}}
								tok, err = d.Token()
								if tok == nil || err != nil {
									{{ setImports $ "fmt" -}}
									return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
								}
								v, _ := tok.(xml.CharData)
								{{ template "parse scalar" (templateMap "shapeRef" $ref.Shape.MemberRef) -}}
								s.{{ $name }} = append(s.{{ $name }}, value)
							{{ else if $ref.Shape.MemberRef.Shape.NestedShape -}}
								err := unmarshalAWSXMLFlattenedList{{ $ref.Shape.ShapeName }}(&s.{{ $name }}, d, start)
								if err != nil {
									{{ setImports $ "fmt" -}}
									return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
								}
							{{ end -}}
						{{ else -}}
							{{ if not $ref.Shape.MemberRef.Shape.NestedShape -}}
								for {
									tok, err = d.Token()
									if tok == nil || err != nil {
										{{ setImports $ "fmt" -}}
										return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
									}
									if end, ok := tok.(xml.EndElement); ok {
										name := end.Name.Local
										if name == "{{ getStructFieldOuterTagName $ref $name }}" {
											break
										}
									}
									if start, ok := tok.(xml.StartElement); ok {
										switch name = start.Name.Local; name {
										case "{{ $ref.GetUnflattenListInnerTagName }}":
											tok, err = d.Token()
											if tok == nil || err != nil {
												{{ setImports $ "fmt" -}}
												return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
											}
											v, _ := tok.(xml.CharData)
											{{ template "parse scalar" (templateMap "shapeRef" $ref.Shape.MemberRef) -}}
											s.{{ $name }} = append(s.{{ $name }}, value)
										default:
											err := d.Skip()
											if err != nil {
												return err
											}
										}
									}
								}
							{{ else if $ref.Shape.MemberRef.Shape.NestedShape -}}
								err := unmarshalAWSXMLList{{ $ref.Shape.ShapeName }}(&s.{{ $name }}, d, start)
								if err != nil {
									{{ setImports $ "fmt" -}}
									return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
								}
							{{ end -}}
						{{ end -}}
					{{ else if eq $ref.Shape.Type "map" -}}
						if s.{{ $name }} == nil {
							s.{{ $name }} = make({{ $ref.GoTypeElem }})
						}
						{{ if or $ref.Flattened $ref.Shape.Flattened -}}
							{{ if not $ref.Shape.ValueRef.Shape.NestedShape -}}
								curKey := ""
								for {
									tok, err = d.Token()
									if tok == nil || err != nil {
										{{ setImports $ "fmt" -}}
										return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
									}
									if end, ok := tok.(xml.EndElement); ok {
										name := end.Name.Local
										if name == "{{ getStructFieldOuterTagName $ref $name }}" {
											break
										}
									}
									if start, ok := tok.(xml.StartElement); ok {
										switch name = start.Name.Local; name {
										case "{{ $ref.GetMapInnerKeyTagName }}":
											tok, err = d.Token()
											if tok == nil || err != nil {
												{{ setImports $ "fmt" -}}
												return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
											}
											v, _ := tok.(xml.CharData)
											curKey = string(v)
										case "{{ $ref.GetMapInnerValueTagName }}":
											tok, err = d.Token()
											if tok == nil || err != nil {
												{{ setImports $ "fmt" -}}
												return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
											}
											v, _ := tok.(xml.CharData)
											{{ template "parse scalar" (templateMap "shapeRef" $ref.Shape.ValueRef) -}}
											s.{{ $name }}[curKey] = value
										default:
											err := d.Skip()
											if err != nil {
												return err
											}
										}
									}
								}
							{{ else if $ref.Shape.ValueRef.Shape.NestedShape -}}
								err := unmarshalAWSXMLFlattenedMap{{ $ref.Shape.ShapeName }}(&s.{{ $name }}, d, start)
								if err != nil {
									{{ setImports $ "fmt" -}}
									return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
								}
							{{ end -}}
						{{ else -}}
							{{ if not $ref.Shape.ValueRef.Shape.NestedShape -}}
								curKey := ""
								for {
									tok, err := d.Token()
									if tok == nil || err != nil {
										return err
									}
									if end, ok := tok.(xml.EndElement); ok {
										name := end.Name.Local
										if name == "{{ getStructFieldOuterTagName $ref $name }}" {
											break
										}
									}
									if start, ok := tok.(xml.StartElement); ok {
										switch name = start.Name.Local; name {
										case "{{ $ref.GetMapInnerKeyTagName }}":
											tok, err := d.Token()
											if tok == nil || err != nil {
												return err
											}
											v, _ := tok.(xml.CharData)
											curKey = string(v)
										case "{{ $ref.GetMapInnerValueTagName }}":
											tok, err = d.Token()
											if tok == nil || err != nil {
												return err
											}
											v, _ := tok.(xml.CharData)
											{{ template "parse scalar" (templateMap "shapeRef" $ref.Shape.ValueRef) -}}
											s.{{ $name }}[curKey] = value
										case "entry":
											continue
										default:
											err := d.Skip()
											if err != nil {
												return err
											}
										}
									}
								}	
							{{ else if $ref.Shape.ValueRef.Shape.NestedShape -}}
								if s.{{ $name }} == nil {
									s.{{ $name }} = make({{ $ref.GoTypeElem }})
								}
								err := unmarshalAWSXMLMap{{ $ref.Shape.ShapeName }}(&s.{{ $name }}, d, start)
								if err != nil {
									{{ setImports $ "fmt" -}}
									return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)

								}
							{{ end -}}
						{{ end -}}
					{{ else if eq $ref.Shape.Type "structure" -}}
						value := {{ $ref.Shape.GoTypeElem }}{}
						err := value.unmarshalAWSXML(d, start)
						if err != nil {
								{{ setImports $ "fmt" -}}
								return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
						}
						s.{{ $name }} = &value
					{{ else if not $ref.Shape.NestedShape -}}
						tok, err = d.Token()
						if tok == nil || err != nil {
							{{ setImports $ "fmt" -}}
							return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
						}
						v, _ := tok.(xml.CharData)
						{{ template "parse scalar" (templateMap "shapeRef" $ref) -}}
						s.{{ $name }} = {{ if and (not (eq $ref.Shape.Type "blob")) (not $ref.Shape.IsEnum) -}}&{{ end }}value
					{{ end -}}
				{{ end -}}
			{{ end -}}
			default:
				err := d.Skip()
				if err != nil {
					{{ setImports $ "fmt" -}}
					return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
				}
			}
		}
	}
}

{{ end -}}

{{ if $.NeedUnmarshalXMLListHelper -}}
{{ setImports $ "encoding/xml" -}}

func unmarshalAWSXMLList{{ $shapeName }}(s *{{ $.GoTypeElem }}, d *xml.Decoder, head xml.StartElement) (err error) {
	defer func() {
		if err != nil {
			*s = make({{ $.GoTypeElem }}, 0)
		}
	}()
	name := ""
	for {
		tok, err := d.Token()
		if tok == nil || err != nil {
			return err
		}
		if end, ok := tok.(xml.EndElement); ok {
			name = end.Name.Local
			if name == head.Name.Local {
				break
			}
		}
		if start, ok := tok.(xml.StartElement); ok {
			switch name = start.Name.Local; name {
			case "{{ $.GetUnflattenListInnerTagName }}":
				{{ if eq $.MemberRef.Shape.Type "structure" -}}
					value := {{ $.MemberRef.GoTypeElem }}{}
					err := value.unmarshalAWSXML(d, start)
					if err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
				{{ else if eq $.MemberRef.Shape.Type "list" -}}
					value := make({{ $.MemberRef.GoTypeElem }}, 0)
					{{ if or $.MemberRef.Flattened $.MemberRef.Shape.Flattened -}}
						err := unmarshalAWSXMLFlattenedList{{ $.MemberRef.Shape.ShapeName }}(&value, d, start)
					{{ else -}}
						err := unmarshalAWSXMLList{{ $.MemberRef.Shape.ShapeName }}(&value, d, start)
					{{ end -}}
					if err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
				{{ else if eq $.MemberRef.Shape.Type "map" -}}
					value := make({{ $.MemberRef.GoTypeElem }})
					{{ if or $.MemberRef.Flattened $.MemberRef.Shape.Flattened -}}
						err := unmarshalAWSXMLFlattenedMap{{ $.MemberRef.Shape.ShapeName }}(&value, d, start)
					{{ else -}}
						err := unmarshalAWSXMLMap{{ $.MemberRef.Shape.ShapeName }}(&value, d, start)
					{{ end -}}
					if err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
				{{ else if not $.MemberRef.Shape.NestedShape -}}
					tok, err = d.Token()
					if tok == nil || err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
					v, _ := tok.(xml.CharData)
					{{ template "parse scalar" (templateMap "shapeRef" $.MemberRef) -}}
				{{ end -}}
				*s = append(*s, value)
			default:
				err := d.Skip()
				if err != nil {
					{{ setImports $ "fmt" -}}
					return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
				}
			}
		}
	}
	return nil
}

func unmarshalAWSXMLFlattenedList{{ $shapeName }}(s *{{ $.GoTypeElem }}, d *xml.Decoder, head xml.StartElement) (err error) {
	defer func() {
		if err != nil {
			*s = make({{ $.GoTypeElem }}, 0)
		}
	}()
	{{ if eq $.MemberRef.Shape.Type "structure" -}}
		value := {{ $.MemberRef.GoTypeElem }}{}
		err = value.unmarshalAWSXML(d, head)
		if err != nil {
			return err
		}
		*s = append(*s, value)
	{{ else if eq $.MemberRef.Shape.Type "list" -}}
		value := make({{ $.MemberRef.GoTypeElem }}, 0)
		{{ if or $.MemberRef.Flattened $.MemberRef.Shape.Flattened -}}
			err = unmarshalAWSXMLFlattenedList{{ $.MemberRef.Shape.ShapeName }}(&value, d, head)
		{{ else -}}
			err = unmarshalAWSXMLList{{ $.MemberRef.Shape.ShapeName }}(&value, d, head)
		{{ end -}}
		if err != nil {
			return err
		}
		*s = append(*s, value)
	{{ else if eq $.MemberRef.Shape.Type "map" -}}
		value := make({{ $.MemberRef.GoTypeElem }})
		{{ if or $.MemberRef.Flattened $.MemberRef.Shape.Flattened -}}
			err = unmarshalAWSXMLFlattenedMap{{ $.MemberRef.Shape.ShapeName }}(&value, d, head)
		{{ else -}}
			err = unmarshalAWSXMLMap{{ $.MemberRef.Shape.ShapeName }}(&value, d, head)
		{{ end -}}
		if err != nil {
			return err
		}
		*s = append(*s, value)
	{{ else if $.MemberRef.Shape.NestedShape -}}
		for {
			tok, err = d.Token()
			if tok == nil || err != nil {
				return err
			}
			if end, ok := tok.(xml.EndElement); ok {
				name := end.Name.Local
				if name == head.Name.Local {
					break
				}
			}
			v, _ := tok.(xml.CharData)
			{{ template "parse scalar" (templateMap "shapeRef" $.MemberRef) -}}
			*s = append(*s, value)
		}
	{{ end -}}
	return nil
}

{{ end -}}

{{ if $.NeedUnmarshalXMLMapHelper -}}
{{ setImports $ "encoding/xml" -}}

func unmarshalAWSXMLMap{{ $shapeName }}(s *{{ $.GoTypeElem }}, d *xml.Decoder, head xml.StartElement) (err error) {
	defer func() {
		if err != nil {
			*s = make({{ $.GoTypeElem }})
		}
	}()
	name := ""
	curKey := ""
	for {
		tok, err := d.Token()
		if tok == nil || err != nil {
			{{ setImports $ "fmt" -}}
			return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
		}
		if end, ok := tok.(xml.EndElement); ok {
			name = end.Name.Local
			if name == head.Name.Local {
				break
			}
		}
		if start, ok := tok.(xml.StartElement); ok {
			switch name = start.Name.Local; name {
			case "{{ $.GetMapInnerKeyTagName }}":
				tok, err := d.Token()
				if tok == nil || err != nil {
					{{ setImports $ "fmt" -}}
					return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
				}
				v, _ := tok.(xml.CharData)
				curKey = string(v)
			case "{{ $.GetMapInnerValueTagName }}":
				{{ if eq $.ValueRef.Shape.Type "structure" -}}
					value := {{ $.ValueRef.GoTypeElem }}{}
					err := value.unmarshalAWSXML(d, start)
					if err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
				{{ else if eq $.ValueRef.Shape.Type "list" -}}
					value := make({{ $.ValueRef.GoTypeElem }}, 0)
					{{ if or $.ValueRef.Flattened $.ValueRef.Shape.Flattened -}}
						err := unmarshalAWSXMLFlattenedList{{ $.ValueRef.Shape.ShapeName }}(&value, d, start)
					{{ else -}}
						err := unmarshalAWSXMLList{{ $.ValueRef.Shape.ShapeName }}(&value, d, start)
					{{ end -}}
					if err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
				{{ else if eq $.ValueRef.Shape.Type "map" -}}
					value := make({{ $.ValueRef.GoTypeElem }})
					{{ if or $.ValueRef.Flattened $.ValueRef.Shape.Flattened -}}
						err := unmarshalAWSXMLFlattenedMap{{ $.ValueRef.Shape.ShapeName }}(&value, d, start)
					{{ else -}}
						err := unmarshalAWSXMLMap{{ $.ValueRef.Shape.ShapeName }}(&value, d, start)
					{{ end -}}
					if err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
				{{ else if not $.ValueRef.Shape.NestedShape -}}
					tok, err = d.Token()
					if tok == nil || err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
					v, _ := tok.(xml.CharData)
					{{ template "parse scalar" (templateMap "shapeRef" $.ValueRef) -}}
				{{ end -}}
				(*s)[curKey] = value
			case "entry":
				continue
			default:
				err := d.Skip()
				if err != nil {
					{{ setImports $ "fmt" -}}
					return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
				}
			}
		}
	}
	return nil
}

func unmarshalAWSXMLFlattenedMap{{ $shapeName }}(s *{{ $.GoTypeElem }}, d *xml.Decoder, head xml.StartElement) (err error) {
	defer func() {
		if err != nil {
			*s = make({{ $.GoTypeElem }})
		}
	}()
	name := ""
	curKey := ""
	for {
		tok, err := d.Token()
		if tok == nil || err != nil {
			{{ setImports $ "fmt" -}}
			return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
		}
		if end, ok := tok.(xml.EndElement); ok {
			name = end.Name.Local
			if name == head.Name.Local {
				break
			}
		}
		if start, ok := tok.(xml.StartElement); ok {
			switch name = start.Name.Local; name {
			case "{{ $.GetMapInnerKeyTagName }}":
				tok, err := d.Token()
				if tok == nil || err != nil {
					{{ setImports $ "fmt" -}}
					return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
				}
				v, _ := tok.(xml.CharData)
				curKey = string(v)
			case "{{ $.GetMapInnerValueTagName }}":
				{{ if eq $.ValueRef.Shape.Type "structure" -}}
					value := {{ $.ValueRef.GoTypeElem }}{}
					err := value.unmarshalAWSXML(d, start)
					if err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
				{{ else if eq $.ValueRef.Shape.Type "list" -}}
					value := make({{ $.ValueRef.GoTypeElem }}, 0)
					{{ if or $.ValueRef.Flattened $.ValueRef.Shape.Flattened -}}
						err := unmarshalAWSXMLFlattenedList{{ $.ValueRef.Shape.ShapeName }}(&value, d, start)
					{{ else -}}
						err := unmarshalAWSXMLList{{ $.ValueRef.Shape.ShapeName }}(&value, d, start)
					{{ end -}}
					if err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
				{{ else if eq $.ValueRef.Shape.Type "map" -}}
					value := make({{ $.ValueRef.GoTypeElem }})
					{{ if or $.ValueRef.Flattened $.ValueRef.Shape.Flattened -}}
						err := unmarshalAWSXMLFlattenedMap{{ $.ValueRef.Shape.ShapeName }}(&value, d, start)
					{{ else -}}
						err := unmarshalAWSXMLMap{{ $.ValueRef.Shape.ShapeName }}(&value, d, start)
					{{ end -}}
					if err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
				{{ else if not $.ValueRef.Shape.NestedShape -}}
					tok, err = d.Token()
					if tok == nil || err != nil {
						{{ setImports $ "fmt" -}}
						return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
					}
					v, _ := tok.(xml.CharData)
					{{ template "parse scalar" (templateMap "shapeRef" $.ValueRef) -}}
				{{ end -}}
				(*s)[curKey] = value
			default:
				err := d.Skip()
				if err != nil {
					{{ setImports $ "fmt" -}}
					return fmt.Errorf("fail to UnmarshalAWSXML {{ $.GoTypeElem }}.%s, %s", name, err)
				}
			}
		}
	}
	return nil
}
{{ end -}}
{{ end -}}

{{ define "parse scalar" -}}
	{{ $ref := index $ "shapeRef" -}}
	{{ if or (eq $ref.Shape.Type "integer") (eq $ref.Shape.Type "long") -}}
		{{ setImports $ref.Shape "strconv" -}}
		value, _ := strconv.ParseInt(string(v), 10, 64)
	{{ else if or (eq $ref.Shape.Type "float") (eq $ref.Shape.Type "double") -}}
		{{ setImports $ref.Shape "strconv" -}}
		value, _ := strconv.ParseFloat(string(v), 64)
	{{ else if eq $ref.Shape.Type "boolean" -}}
		{{ setImports $ref.Shape "strconv" -}}
		value, _ := strconv.ParseBool(string(v))
	{{ else if eq $ref.Shape.Type "blob" -}}
		{{ setImports $ref.Shape "encoding/base64" -}}
		value, _ := base64.StdEncoding.DecodeString(string(v))
	{{ else if eq $ref.Shape.Type "timestamp" -}}
		{{ setImports $ref.Shape "time" -}}
		value, _ := time.Parse({{ $ref.CheckTimeFormatShapeRef }}, string(v))
	{{ else if $ref.Shape.IsEnum -}}
		value := {{ $ref.Shape.GoTypeElem }}(v)
	{{ else if or (eq $ref.Shape.Type "string") (eq $ref.Shape.Type "character") -}}
		value := string(v)
	{{ end -}}
{{ end -}}
`))

var unmarshalShapeRESTTmpl = template.Must(template.New("unmarshalShapeRESTTmpl").Funcs(
	template.FuncMap{
		"setImports": setImports,
		"setSDKImports": setSDKImports,
		"getHeaderKeyName": getHeaderKeyName,
	},
).Parse(`
{{ $shapeName := $.ShapeName -}}
{{ if $.NeedUnmarshalREST -}}
{{ setImports $ "net/http" -}}

// UnmarshalAWSREST decodes the AWS API shape using the passed in *http.Response.
func (s *{{ $.GoTypeElem }}) UnmarshalAWSREST(r *http.Response) (err error) {
	defer func() {
		if err != nil {
			*s = {{ $.GoTypeElem }}{}
		}
	}()
	{{ $statusCode := ($.GetStatusCodeField) -}}
	{{ if $statusCode -}}
		{{ $shapeRef := index $.MemberRefs $statusCode -}}
		{{ if eq $shapeRef.Shape.Type "integer" -}}
			s.{{ $statusCode }} = &(r.StatusCode)
		{{ else -}}
			{{ setImports $ "errors" -}}
			return errors.New("rest protocol doesn't support unmarshaling this type of field.")
		{{ end -}}
	{{ end -}}
	for k, v := range r.Header {
		switch {
		{{ range $name, $ref := $.MemberRefs -}}
			{{ $location := $ref.GetLocation -}}
			{{ if or (eq $location "header") (eq $location "headers") -}}
				{{ if not $ref.Shape.NestedShape -}}
					{{ setImports $ "strings" -}}
					case strings.EqualFold(k, "{{ getHeaderKeyName $ref $name }}"):
					{{ if or (eq $ref.Shape.Type "integer") (eq $ref.Shape.Type "long") -}}
						{{ setImports $ "strconv" -}}
						{{ setImports $ "fmt" -}}
						value, err := strconv.ParseInt(v[0], 10, 64)
						if err != nil {
							return fmt.Errorf("fail to UnmarshalAWSREST {{ $.GoTypeElem }}.{{ $name }}, %s", err)
						}
						s.{{ $name }} = &value
					{{ else if or (eq $ref.Shape.Type "float") (eq $ref.Shape.Type "double") -}}
						{{ setImports $ "strconv" -}}
						{{ setImports $ "fmt" -}}
						value, err := strconv.ParseFloat(v[0], 64)
						if err != nil {
							return fmt.Errorf("fail to UnmarshalAWSREST {{ $.GoTypeElem }}.{{ $name }}, %s", err)
						}
						s.{{ $name }} = &value
					{{ else if eq $ref.Shape.Type "boolean" -}}
						{{ setImports $ "strconv" -}}
						{{ setImports $ "fmt" -}}
						value, err := strconv.ParseBool(v[0])
						if err != nil {
							return fmt.Errorf("fail to UnmarshalAWSREST {{ $.GoTypeElem }}.{{ $name }}, %s", err)
						}
						s.{{ $name }} = &value
					{{ else if eq $ref.Shape.Type "blob" -}}
						{{ setImports $ "encoding/base64" -}}
						value, _ := base64.StdEncoding.DecodeString(v[0])
						s.{{ $name }} = value
					{{ else if eq $ref.Shape.Type "timestamp" -}}
						{{ setImports $ "time" -}}
						{{ setImports $ "fmt" -}}
						value, err := time.Parse({{ $ref.CheckTimeFormatShapeRef }}, v[0])
						if err != nil {
							return fmt.Errorf("fail to UnmarshalAWSREST {{ $.GoTypeElem }}.{{ $name }}, %s", err)
						}
						s.{{ $name }} = &value
					{{ else if $ref.Shape.IsEnum -}}
						value := {{ $ref.Shape.GoTypeElem }}(v[0])
						s.{{ $name }} = value
					{{ else if or (eq $ref.Shape.Type "string") (eq $ref.Shape.Type "character") -}}
						value := v[0]
						s.{{ $name }} = &value
					{{ else if eq $ref.Shape.Type "jsonvalue" -}}
						{{ setSDKImports $ "private/protocol" -}}
						escaping := protocol.Base64Escape
						value, err := protocol.DecodeJSONValue(v[0], escaping)
						if err != nil {
							return fmt.Errorf("fail to UnmarshalAWSREST {{ $.GoTypeElem }}.{{ $name }}, %s", err)
						}
						s.{{ $name }} = &value
					{{ else -}}
						{{ setImports $ "errors" -}}
						return errors.New("rest protocol doesn't support unmarshaling this type of field.")
					{{ end -}}
				{{ else if $ref.Shape.NestedShape -}}
					{{ if eq $ref.Shape.Type "map" -}}
						{{ setImports $ "strings" -}}
						case strings.HasPrefix(strings.ToLower(k), "{{ $ref.GetLocationName }}"):
						{{ if eq $ref.Shape.ValueRef.Shape.Type "string" -}}
							if s.{{ $name }} == nil {
								s.{{ $name }} = map[{{ $ref.Shape.KeyRef.GoTypeElem }}]{{ $ref.Shape.ValueRef.GoTypeElem }}{}
							}
							s.{{ $name }}[k[len("{{ $ref.GetLocationName }}"):]] = v[0]
						{{ else if not eq $ref.Shape.ValueRef.Shape.Type "string" -}}
							{{ setImports $ "errors" -}}
							return errors.New("rest protocol doesn't support unmarshaling this type of field.")
						{{ end -}}
					{{ else -}}
						{{ setImports $ "errors" -}}
						return errors.New("rest protocol doesn't support unmarshaling this type of field.")
					{{ end -}}
				{{ end -}}
			{{ end -}}
		{{ end -}}
		}
	}
	return nil
}

{{ end -}}
`))

var unmarshalShapePayloadTmpl = template.Must(template.New("unmarshalShapePayloadTmpl").Funcs(
	template.FuncMap{
		"setImports": setImports,
		"setSDKImports": setSDKImports,
	},
).Parse(`
{{ $shapeName := $.ShapeName -}}
{{ if $.NeedUnmarshalPayload -}}
{{ setImports $ "io" -}}

// UnmarshalAWSPayload decodes the AWS API shape using the passed in io.ReadCloser.
func (s *{{ $.GoTypeElem }}) UnmarshalAWSPayload(r io.ReadCloser) (err error) {
	defer func() {
		if err != nil {
			*s = {{ $.GoTypeElem }}{}
		}
	}()
	{{ $ref := index $.MemberRefs $.Payload -}}
	{{ if eq $ref.Shape.Type "blob" -}}
		{{ $goType := $.GoStructType $ref.ShapeName $ref -}}
		{{ if eq $goType "[]byte" -}}
			{{ setImports $ "bytes" -}}
			buf := new(bytes.Buffer)
			buf.ReadFrom(r)
			s.{{ $.Payload }} = []byte(buf.String())
		{{ else if eq $goType "io.ReadCloser" -}}
			s.{{ $.Payload }} = r
		{{ else -}}
			{{ setImports $ "errors" -}}
			return errors.New("payload protocol doesn't support unmarshaling this type of field.")
		{{ end -}}
	{{ else if eq $ref.Shape.Type "string" -}}
		{{ setImports $ "bytes" -}}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		str := buf.String()
		s.{{ $.Payload }} = &str
	{{ else -}}
		{{ setImports $ "errors" -}}
		return errors.New("payload protocol doesn't support unmarshaling this type of field.")
	{{ end -}}
	return nil
}

{{ end -}}
`))

func (s *Shape) NeedUnmarshalXMLEntry() bool {
	if !s.UsedAsOutput {
		return false
	}
	if s.NeedUnmarshalPayload() {
		return false
	}
	array := s.MemberNames()
	for _, name := range array {
		if s.MemberRefs[name].GetLocation() != "header" && s.MemberRefs[name].GetLocation() != "headers" && s.MemberRefs[name].GetLocation() != "statusCode" {
			return true
		}
	}
	return false
}

func (s *Shape) NeedUnmarshalXMLStructHelper() bool {
	if s.Type != "structure" {
		return false
	}
	if !s.IsInputOrOutput() {
		return true
	}
	if s.UsedAsOutput && s.Payload == "" {
		return true
	}
	return false
}

func (s *Shape) NeedUnmarshalXMLListHelper() bool {
	return s.Type == "list"
}

func (s *Shape) NeedUnmarshalXMLMapHelper() bool {
	return s.Type == "map"
}

func (s *Shape) NeedUnmarshalREST() bool {
	if !s.UsedAsOutput {
		return false
	}
	array := s.MemberNames()
	for _, name := range array {
		if s.MemberRefs[name].GetLocation() == "header" || s.MemberRefs[name].GetLocation() == "headers" || s.MemberRefs[name].GetLocation() == "statusCode" {
			return true
		}
	}
	return false
}

func (s *Shape) NeedUnmarshalPayload() bool {
	if !s.UsedAsOutput {
		return false
	}
	ref, ok := s.MemberRefs[s.Payload]
	if !ok {
		return false
	}
	if ref.Shape.Type == "structure" {
		return false
	}
	return true
}

func (s *Shape) GetPayloadShapeName() string {
	ref, ok := s.MemberRefs[s.Payload]
	if !ok {
		return ""
	}
	return ref.ShapeName
}

func (s *Shape) GetPayloadShapeRef() *ShapeRef {
	ref, ok := s.MemberRefs[s.Payload]
	if !ok {
		return nil
	}
	return ref
}

func (s *Shape) GetStatusCodeField() string {
	for key, value := range s.MemberRefs {
		if value.Location == "statusCode" || value.Shape.Location == "statusCode" {
			return key
		}
	}
	return ""
}

func (s *Shape) GetMapInnerKeyTagName() string {
	if s.Type != "map" {
		return ""
	}
	if s.KeyRef.LocationName == "" && s.KeyRef.Shape.LocationName == "" {
		return "key"
	}
	locationName := s.KeyRef.LocationName
	if len(locationName) == 0 {
		locationName = s.KeyRef.Shape.LocationName
	}
	return locationName
}

func (s *Shape) GetMapInnerValueTagName() string {
	if s.Type != "map" {
		return ""
	}
	if s.ValueRef.LocationName == "" && s.ValueRef.Shape.LocationName == "" {
		return "value"
	}
	locationName := s.ValueRef.LocationName
	if len(locationName) == 0 {
		locationName = s.ValueRef.Shape.LocationName
	}
	return locationName
}

func (s *Shape) GetXMLAttributeField() string {
	array := s.MemberNames()
	for _, name := range array {
		if s.MemberRefs[name].XMLAttribute {
			return name
		}
	}
	return ""
}

func (s *Shape) IsInputOrOutput() bool {
	return s.UsedAsInput || s.UsedAsOutput
}

func (s *Shape) GetUnflattenListInnerTagName() string {
	if s.Type != "list" {
		return ""
	}
	if s.MemberRef.LocationName == "" && s.MemberRef.Shape.LocationName == "" {
		return "member"
	}
	location := s.MemberRef.LocationName
	if len(location) == 0 {
		location = s.MemberRef.Shape.LocationName
	}
	return location
}

func (ref *ShapeRef) GetLocation() string {
	location := ref.Location
	if len(location) == 0 {
		location = ref.Shape.Location
	}
	return location
}

func (ref *ShapeRef) GetLocationName() string {
	locationName := ref.LocationName
	if len(locationName) == 0 {
		locationName = ref.Shape.LocationName
	}
	return locationName
}

func (ref ShapeRef) CheckTimeFormatShapeRef() string {
	switch ref.GetLocation() {
	case "header", "headers":
		return "protocol.RFC822TimeFromat"
	case "query":
		return "protocol.RFC822TimeFromat"
	default:
		switch ref.API.Metadata.Protocol {
		case "json", "rest-json":
			return "protocol.UnixTimeFormat"
		case "rest-xml", "ec2", "query":
			return "protocol.ISO8601TimeFormat"
		default:
			panic(fmt.Sprintf("unable to determine time format for %s ref", ref.ShapeName))
		}
	}
}

func (ref *ShapeRef) GetUnflattenListInnerTagName() string {
	if ref.Shape.Type != "list" {
		return ""
	}
	if ref.Shape.MemberRef.LocationName == "" && ref.Shape.MemberRef.Shape.LocationName == "" {
		return "member"
	}
	locationName := ref.Shape.MemberRef.LocationName
	if len(locationName) == 0 {
		locationName = ref.Shape.MemberRef.Shape.LocationName
	}
	return locationName
}

func (ref *ShapeRef) GetMapInnerKeyTagName() string {
	if ref.Shape.Type != "map" {
		return ""
	}
	if ref.Shape.KeyRef.LocationName == "" && ref.Shape.KeyRef.Shape.LocationName == "" {
		return "key"
	}
	locationName := ref.Shape.KeyRef.LocationName
	if len(locationName) == 0 {
		locationName = ref.Shape.KeyRef.Shape.LocationName
	}
	return locationName
}

func (ref *ShapeRef) GetMapInnerValueTagName() string {
	if ref.Shape.Type != "map" {
		return ""
	}
	if ref.Shape.ValueRef.LocationName == "" && ref.Shape.ValueRef.Shape.LocationName == "" {
		return "value"
	}
	locationName := ref.Shape.ValueRef.LocationName
	if len(locationName) == 0 {
		locationName = ref.Shape.ValueRef.Shape.LocationName
	}
	return locationName
}

func setImports(s *Shape, importPath string) string {
	s.API.AddImport(importPath)
	return ""
}

func setSDKImports(s *Shape, importPath string) string {
	s.API.AddSDKImport(importPath)
	return ""
}

func getHeaderKeyName(ref *ShapeRef, fieldName string) string {
	if ref.LocationName == "" {
		return fieldName
	}
	return ref.LocationName
}

func getStructFieldOuterTagName(ref *ShapeRef, fieldName string) string {
	if ref.Deprecated || ref.Shape.Deprecated {
		return ""
	}
	if ref.LocationName == "" && ref.Shape.LocationName == "" {
		return fieldName
	}
	locationName := ref.LocationName
	if len(locationName) == 0 {
		locationName = ref.Shape.LocationName
	}
	return locationName
}

func findLocationNameLocal(locationName string) string {
	i := strings.Index(locationName, ":")
	return locationName[i + 1:]
}

func templateMap(args ...interface{}) map[string]interface{} {
	if len(args) % 2 != 0 {
		panic(fmt.Sprintf("invalid map call, non-even args %v", args))
	}
	m := map[string]interface{}{}
	for i := 0; i < len(args); i += 2 {
		k, ok := args[i].(string)
		if !ok {
			panic(fmt.Sprintf("invalid map call, arg is not string, %T, %v", args[i], args[i]))
		}
		m[k] = args[i + 1]
	}
	return m
}
