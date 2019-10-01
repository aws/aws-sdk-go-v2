// +build codegen

package api

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// An Operation defines a specific API Operation.
type Operation struct {
	API                 *API `json:"-"`
	ExportedName        string
	Name                string
	Documentation       string
	HTTP                HTTPInfo
	InputRef            ShapeRef   `json:"input"`
	OutputRef           ShapeRef   `json:"output"`
	ErrorRefs           []ShapeRef `json:"errors"`
	Paginator           *Paginator
	Deprecated          bool     `json:"deprecated"`
	AuthType            AuthType `json:"authtype"`
	imports             map[string]bool
	CustomBuildHandlers []string       `json:"-"`
	Endpoint            *EndpointTrait `json:"endpoint"`
}

// EndpointTrait provides the structure of the modeled enpdoint trait, and its
// properties.
type EndpointTrait struct {
	// Specifies the hostPrefix template to prepend to the operation's request
	// endpoint host.
	HostPrefix string `json:"hostPrefix"`
}

// A HTTPInfo defines the method of HTTP request for the Operation.
type HTTPInfo struct {
	Method       string
	RequestURI   string
	ResponseCode uint
}

// HasInput returns if the Operation accepts an input paramater
func (o *Operation) HasInput() bool {
	return o.InputRef.ShapeName != ""
}

// HasOutput returns if the Operation accepts an output parameter
func (o *Operation) HasOutput() bool {
	return o.OutputRef.ShapeName != ""
}

// AuthType provides the enumeration of AuthType trait.
type AuthType string

// Enumeration values for AuthType trait
const (
	NoneAuthType           AuthType = "none"
	V4UnsignedBodyAuthType AuthType = "v4-unsigned-body"
)

// GetSigner returns the signer to use for a request.
func (o *Operation) GetSigner() string {
	buf := bytes.NewBuffer(nil)

	switch o.AuthType {
	case NoneAuthType:
		o.API.AddSDKImport("aws")

		buf.WriteString("req.Config.Credentials = aws.AnonymousCredentials")
	case V4UnsignedBodyAuthType:
		o.API.AddSDKImport("aws/signer/v4")

		buf.WriteString("req.Handlers.Sign.Remove(v4.SignRequestHandler)\n")
		buf.WriteString("handler := v4.BuildNamedHandler(\"v4.CustomSignerHandler\", v4.WithUnsignedPayload)\n")
		buf.WriteString("req.Handlers.Sign.PushFrontNamed(handler)")
	}

	buf.WriteString("\n")
	return buf.String()
}

// operationTmpl defines a template for rendering an API Operation
var operationTmpl = template.Must(template.New("operation").Funcs(template.FuncMap{
	"GetCrosslinkURL": GetCrosslinkURL,
}).Parse(`
{{ $reqType := printf "%sRequest" .ExportedName -}}
{{ $respType := printf "%sResponse" .ExportedName -}}
{{ $pagerType := printf "%sPaginator" .ExportedName -}}

{{ if .HasInput -}}
	{{ .InputRef.Shape.GoCode }}
{{- end }}

{{ if .HasOutput -}}
	{{ .OutputRef.Shape.GoCode }}
{{- end }}

const op{{ .ExportedName }} = "{{ .Name }}"

// {{ $reqType }} returns a request value for making API operation for
// {{ .API.Metadata.ServiceFullName }}.
{{ if .Documentation -}}
//
{{ .Documentation }}
{{ end -}}
//
//    // Example sending a request using {{ $reqType }}.
//    req := client.{{ $reqType }}(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
{{ $crosslinkURL := GetCrosslinkURL $.API.BaseCrosslinkURL $.API.Metadata.UID $.ExportedName -}}
{{ if ne $crosslinkURL "" -}}
//
// Please also see {{ $crosslinkURL }}
{{ end -}}
func (c *{{ .API.StructName }}) {{ $reqType }}(input {{ .InputRef.GoType }}) ({{ $reqType }}) {
	{{ if (or .Deprecated (or .InputRef.Deprecated .OutputRef.Deprecated)) -}}
		if c.Client.Config.Logger != nil {
			c.Client.Config.Logger.Log("This operation, {{ .ExportedName }}, has been deprecated")
		}
	{{ end -}}

	op := &aws.Operation{
		Name:       op{{ .ExportedName }},
		{{ if ne .HTTP.Method "" }}HTTPMethod: "{{ .HTTP.Method }}",{{ end }}
		HTTPPath: {{ if ne .HTTP.RequestURI "" }}"{{ .HTTP.RequestURI }}"{{ else }}"/"{{ end }},
		{{- if .Paginator }}
		Paginator: &aws.Paginator{
				InputTokens: {{ .Paginator.InputTokensString }},
				OutputTokens: {{ .Paginator.OutputTokensString }},
				LimitToken: "{{ .Paginator.LimitKey }}",
				TruncationToken: "{{ .Paginator.MoreResults }}",
		},
		{{ end }}
	}

	if input == nil {
		input = &{{ .InputRef.GoTypeElem }}{}
	}

	req := c.newRequest(op, input, &{{ .OutputRef.GoTypeElem }}{})
	{{ if eq .OutputRef.Shape.Placeholder true -}}
		req.Handlers.Unmarshal.Remove({{ .API.ProtocolPackage }}.UnmarshalHandler)
		req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	{{ end -}}

	{{ if ne .AuthType "" }}{{ .GetSigner }}{{ end -}}
	{{- range $_, $handler := $.CustomBuildHandlers -}}
		req.Handlers.Build.PushBackNamed({{ $handler }})
	{{ end -}}

	return {{ $reqType }}{Request: req, Input: input, Copy: c.{{ $reqType }} }
}

// {{ $reqType }} is the request type for the
// {{ .ExportedName }} API operation.
type {{ $reqType}} struct {
	*aws.Request
	Input {{ .InputRef.GoType }}
	Copy func({{ .InputRef.GoType }}) {{ $reqType }}
}

// Send marshals and sends the {{ .ExportedName }} API request.
func (r {{ $reqType }}) Send(ctx context.Context) (*{{ $respType }}, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &{{ $respType }}{
		{{ if .HasOutput -}}
			{{ .OutputRef.GoTypeElem }}: r.Request.Data.({{ .OutputRef.GoType }}),
		{{- end }}
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

{{ if .Paginator }}
	// New{{ $reqType }}Paginator returns a paginator for {{ .ExportedName }}.
	// Use Next method to get the next page, and CurrentPage to get the current
	// response page from the paginator. Next will return false, if there are
	// no more pages, or an error was encountered.
	//
	// Note: This operation can generate multiple requests to a service.
	//
	//   // Example iterating over pages.
	//   req := client.{{ $reqType }}(input)
	//   p := {{ .API.PackageName }}.New{{ $reqType }}Paginator(req)
	//
	//   for p.Next(context.TODO()) {
	//       page := p.CurrentPage()
	//   }
	//
	//   if err := p.Err(); err != nil {
	//       return err
	//   }
	//
	func New{{ .ExportedName }}Paginator(req {{ $reqType }}) {{ $pagerType }} {
		return {{ $pagerType }}{
			Pager: aws.Pager {
				NewRequest: func(ctx context.Context) (*aws.Request, error) {
					var inCpy {{ .InputRef.GoType }}
					if req.Input != nil  {
						tmp := *req.Input
						inCpy = &tmp
					}

					newReq := req.Copy(inCpy)
					newReq.SetContext(ctx)
					return newReq.Request, nil
				},
			},
		}
	}

	// {{ $pagerType }} is used to paginate the request. This can be done by
	// calling Next and CurrentPage.
	type {{ $pagerType }} struct {
		aws.Pager
	}

	func (p *{{ $pagerType}}) CurrentPage() {{ .OutputRef.GoType }} {
		return p.Pager.CurrentPage().({{ .OutputRef.GoType }})
	}
{{ end }}

// {{ $respType }} is the response type for the
// {{ .ExportedName }} API operation.
type {{ $respType }} struct {
	{{ if .HasOutput -}}
		{{ .OutputRef.GoType }}
	{{- end }}

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// {{ .ExportedName }} request.
func (r * {{ $respType }}) SDKResponseMetdata() *aws.Response {
	return r.response
}
`))

// GoCode returns a string of rendered GoCode for this Operation
func (o *Operation) GoCode() string {
	if o.Endpoint != nil && len(o.Endpoint.HostPrefix) != 0 {
		setupEndpointHostPrefix(o)
	}

	var buf bytes.Buffer
	err := operationTmpl.Execute(&buf, o)
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}

// tplInfSig defines the template for rendering an Operation's signature within an Interface definition.
var tplInfSig = template.Must(template.New("opsig").Parse(`{{ .ExportedName }}Request({{ .InputRef.GoTypeWithPkgName }}) {{ .API.PackageName }}.{{ .ExportedName }}Request
`))

// InterfaceSignature returns a string representing the Operation's interface{}
// functional signature.
func (o *Operation) InterfaceSignature() string {
	var buf bytes.Buffer
	err := tplInfSig.Execute(&buf, o)
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}

// tplExample defines the template for rendering an Operation example
var tplExample = template.Must(template.New("operationExample").Parse(`
func Example{{ .API.StructName }}_{{ .ExportedName }}() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("failed to load config, " + err.Error())
	}

	svc := {{ .API.PackageName }}.New(sess)

	{{ .ExampleInput }}
	req := svc.{{ .ExportedName }}Request(params)
	resp, err := req.Send(context.TODO())

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
`))

// Example returns a string of the rendered Go code for the Operation
func (o *Operation) Example() string {
	var buf bytes.Buffer
	err := tplExample.Execute(&buf, o)
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}

// ExampleInput return a string of the rendered Go code for an example's input parameters
func (o *Operation) ExampleInput() string {
	if len(o.InputRef.Shape.MemberRefs) == 0 {
		if strings.Contains(o.InputRef.GoTypeElem(), ".") {
			o.imports[SDKImportRoot+"/service/"+strings.Split(o.InputRef.GoTypeElem(), ".")[0]] = true
			o.imports["context"] = true
			return fmt.Sprintf("var params *%s", o.InputRef.GoTypeElem())
		}
		return fmt.Sprintf("var params *%s.%s",
			o.API.PackageName(), o.InputRef.GoTypeElem())
	}
	e := example{o, map[string]int{}}
	return "params := " + e.traverseAny(o.InputRef.Shape, false, false)
}

// A example provides
type example struct {
	*Operation
	visited map[string]int
}

// traverseAny returns rendered Go code for the shape.
func (e *example) traverseAny(s *Shape, required, payload bool) string {
	str := ""
	e.visited[s.ShapeName]++

	switch s.Type {
	case "structure":
		str = e.traverseStruct(s, required, payload)
	case "list":
		str = e.traverseList(s, required, payload)
	case "map":
		str = e.traverseMap(s, required, payload)
	case "jsonvalue":
		str = "aws.JSONValue{\"key\": \"value\"}"
		if required {
			str += " // Required"
		}
	default:
		str = e.traverseScalar(s, required, payload)
	}

	e.visited[s.ShapeName]--

	return str
}

var reType = regexp.MustCompile(`\b([A-Z])`)

// traverseStruct returns rendered Go code for a structure type shape.
func (e *example) traverseStruct(s *Shape, required, payload bool) string {
	var buf bytes.Buffer

	if s.resolvePkg != "" {
		e.imports[s.resolvePkg] = true
		buf.WriteString("&" + s.GoTypeElem() + "{")
	} else {
		buf.WriteString("&" + s.API.PackageName() + "." + s.GoTypeElem() + "{")
	}

	if required {
		buf.WriteString(" // Required")
	}
	buf.WriteString("\n")

	req := make([]string, len(s.Required))
	copy(req, s.Required)
	sort.Strings(req)

	if e.visited[s.ShapeName] < 2 {
		for _, n := range req {
			m := s.MemberRefs[n].Shape
			p := n == s.Payload && (s.MemberRefs[n].Streaming || m.Streaming)
			buf.WriteString(n + ": " + e.traverseAny(m, true, p) + ",")
			if m.Type != "list" && m.Type != "structure" && m.Type != "map" {
				buf.WriteString(" // Required")
			}
			buf.WriteString("\n")
		}

		for _, n := range s.MemberNames() {
			if s.IsRequired(n) {
				continue
			}
			m := s.MemberRefs[n].Shape
			p := n == s.Payload && (s.MemberRefs[n].Streaming || m.Streaming)
			buf.WriteString(n + ": " + e.traverseAny(m, false, p) + ",\n")
		}
	} else {
		buf.WriteString("// Recursive values...\n")
	}

	buf.WriteString("}")
	return buf.String()
}

// traverseMap returns rendered Go code for a map type shape.
func (e *example) traverseMap(s *Shape, required, payload bool) string {
	var buf bytes.Buffer

	t := ""
	if s.resolvePkg != "" {
		e.imports[s.resolvePkg] = true
		t = s.GoTypeElem()
	} else {
		t = reType.ReplaceAllString(s.GoTypeElem(), s.API.PackageName()+".$1")
	}
	buf.WriteString(t + "{")
	if required {
		buf.WriteString(" // Required")
	}
	buf.WriteString("\n")

	if e.visited[s.ShapeName] < 2 {
		m := s.ValueRef.Shape
		buf.WriteString("\"Key\": " + e.traverseAny(m, true, false) + ",")
		if m.Type != "list" && m.Type != "structure" && m.Type != "map" {
			buf.WriteString(" // Required")
		}
		buf.WriteString("\n// More values...\n")
	} else {
		buf.WriteString("// Recursive values...\n")
	}
	buf.WriteString("}")

	return buf.String()
}

// traverseList returns rendered Go code for a list type shape.
func (e *example) traverseList(s *Shape, required, payload bool) string {
	var buf bytes.Buffer
	t := ""
	if s.resolvePkg != "" {
		e.imports[s.resolvePkg] = true
		t = s.GoTypeElem()
	} else {
		t = reType.ReplaceAllString(s.GoTypeElem(), s.API.PackageName()+".$1")
	}

	buf.WriteString(t + "{")
	if required {
		buf.WriteString(" // Required")
	}
	buf.WriteString("\n")

	if e.visited[s.ShapeName] < 2 {
		m := s.MemberRef.Shape
		buf.WriteString(e.traverseAny(m, true, false) + ",")
		if m.Type != "list" && m.Type != "structure" && m.Type != "map" {
			buf.WriteString(" // Required")
		}
		buf.WriteString("\n// More values...\n")
	} else {
		buf.WriteString("// Recursive values...\n")
	}
	buf.WriteString("}")

	return buf.String()
}

// traverseScalar returns an AWS Type string representation initialized to a value.
// Will panic if s is an unsupported shape type.
func (e *example) traverseScalar(s *Shape, required, payload bool) string {
	str := ""
	switch s.Type {
	case "integer", "long":
		str = `aws.Int64(1)`
	case "float", "double":
		str = `aws.Float64(1.0)`
	case "string", "character":
		str = `aws.String("` + s.ShapeName + `")`
	case "blob":
		if payload {
			str = `bytes.NewReader([]byte("PAYLOAD"))`
		} else {
			str = `[]byte("PAYLOAD")`
		}
	case "boolean":
		str = `aws.Bool(true)`
	case "timestamp":
		str = `aws.Time(time.Now())`
	default:
		panic("unsupported shape " + s.Type)
	}

	return str
}
