// +build codegen

// Package api represents API abstractions for rendering service generated files.
package api

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

// SDKImportRoot is the root import path of the SDK.
const SDKImportRoot = "github.com/aws/aws-sdk-go-v2"

// An API defines a service API's definition. and logic to serialize the definition.
type API struct {
	Metadata      Metadata
	Operations    map[string]*Operation
	Shapes        map[string]*Shape
	Waiters       []Waiter
	Documentation string
	Examples      Examples
	SmokeTests    SmokeTestSuite

	// Set to true to avoid removing unused shapes
	NoRemoveUnusedShapes bool

	// Set to true to avoid renaming to 'Input/Output' postfixed shapes
	NoRenameToplevelShapes bool

	// Set to true to ignore service/request init methods (for testing)
	NoInitMethods bool

	// Set to true to not generate String() methods (e.g for generated tests)
	NoStringerMethods bool

	// Set to true to not generate API service name constants
	NoConstServiceNames bool

	// Set to true to not generate validation shapes
	NoValidataShapeMethods bool

	// Set to true to not generate struct field accessors
	NoGenStructFieldAccessors bool

	// Set to true to use the model's ServiceID instead of "Client" for client
	// struct name. For protocol tests with multiple clients in the same
	// package.
	UseServiceIDForClientStruct bool

	// Set to true to not generate (un)marshalers
	NoGenMarshalers   bool
	NoGenUnmarshalers bool

	BaseImportPath string

	initialized bool
	imports     map[string]bool
	serviceID   string
	path        string

	BaseCrosslinkURL string
}

// A Metadata is the metadata about an API's definition.
type Metadata struct {
	APIVersion          string
	EndpointPrefix      string
	SigningName         string
	ServiceAbbreviation string
	ServiceFullName     string
	SignatureVersion    string
	JSONVersion         string
	TargetPrefix        string
	Protocol            string
	UID                 string
	EndpointsID         string `json:"endpointsId"`
	ServiceID           string
}

// PackageName name of the API package
func (a *API) PackageName() string {
	return strings.ToLower(a.ServiceID())
}

// ImportPath returns the client's full import path
func (a *API) ImportPath() string {
	return path.Join(a.BaseImportPath, a.PackageName())
}

// InterfacePackageName returns the package name for the interface.
func (a *API) InterfacePackageName() string {
	return a.PackageName() + "iface"
}

var stripServiceNamePrefixes = []string{
	"Amazon",
	"AWS",
}

// StructName returns the struct name for a given API.
func (a *API) StructName() string {
	if a.UseServiceIDForClientStruct {
		// Allow many protocol test clients to be generated in the same
		// package.
		return a.ServiceID()
	}
	return "Client"
}

// UseInitMethods returns if the service's init method should be rendered.
func (a *API) UseInitMethods() bool {
	return !a.NoInitMethods
}

// NiceName returns the human friendly API name.
func (a *API) NiceName() string {
	name := a.Metadata.ServiceAbbreviation
	if len(name) == 0 {
		name = a.Metadata.ServiceFullName
	}

	return strings.TrimSpace(name)
}

// ProtocolPackage returns the package name of the protocol this API uses.
func (a *API) ProtocolPackage() string {
	switch a.Metadata.Protocol {
	case "json":
		return "jsonrpc"
	case "ec2":
		return "ec2query"
	default:
		return strings.Replace(a.Metadata.Protocol, "-", "", -1)
	}
}

// OperationNames returns a slice of API operations supported.
func (a *API) OperationNames() []string {
	i, names := 0, make([]string, len(a.Operations))
	for n := range a.Operations {
		names[i] = n
		i++
	}
	sort.Strings(names)
	return names
}

// OperationList returns a slice of API operation pointers
func (a *API) OperationList() []*Operation {
	list := make([]*Operation, len(a.Operations))
	for i, n := range a.OperationNames() {
		list[i] = a.Operations[n]
	}
	return list
}

// OperationHasOutputPlaceholder returns if any of the API operation input
// or output shapes are place holders.
func (a *API) OperationHasOutputPlaceholder() bool {
	for _, op := range a.Operations {
		if op.OutputRef.Shape.Placeholder {
			return true
		}
	}
	return false
}

// ShapeNames returns a slice of names for each shape used by the API.
func (a *API) ShapeNames() []string {
	i, names := 0, make([]string, len(a.Shapes))
	for n := range a.Shapes {
		names[i] = n
		i++
	}
	sort.Strings(names)
	return names
}

// ShapeList returns a slice of shape pointers used by the API.
//
// Will exclude error shapes from the list of shapes returned.
func (a *API) ShapeList() []*Shape {
	list := make([]*Shape, 0, len(a.Shapes))
	for _, n := range a.ShapeNames() {
		// Ignore error shapes in list
		if s := a.Shapes[n]; !s.IsError {
			list = append(list, s)
		}
	}
	return list
}

// ShapeListErrors returns a list of the errors defined by the API model
func (a *API) ShapeListErrors() []*Shape {
	list := []*Shape{}
	for _, n := range a.ShapeNames() {
		// Ignore error shapes in list
		if s := a.Shapes[n]; s.IsError {
			list = append(list, s)
		}
	}
	return list
}

func (a *API) ShapeEnumList() []*Shape {
	list := []*Shape{}
	for _, n := range a.ShapeNames() {
		// Ignore error shapes in list
		if s := a.Shapes[n]; s.IsEnum() {
			list = append(list, s)
		}
	}
	return list
}

// resetImports resets the import map to default values.
func (a *API) resetImports() {
	a.imports = map[string]bool{}
}

// importsGoCode returns the generated Go import code.
func (a *API) importsGoCode() string {
	if len(a.imports) == 0 {
		return ""
	}

	corePkgs, extPkgs := []string{}, []string{}
	for i := range a.imports {
		if strings.Contains(i, ".") {
			extPkgs = append(extPkgs, i)
		} else {
			corePkgs = append(corePkgs, i)
		}
	}
	sort.Strings(corePkgs)
	sort.Strings(extPkgs)

	code := "import (\n"
	for _, i := range corePkgs {
		code += fmt.Sprintf("\t%q\n", i)
	}
	if len(corePkgs) > 0 {
		code += "\n"
	}
	for _, i := range extPkgs {
		code += fmt.Sprintf("\t%q\n", i)
	}
	code += ")\n\n"
	return code
}

// A tplAPI is the top level template for the API
var tplAPI = template.Must(template.New("api").Parse(`
{{ range $_, $o := .OperationList }}
{{ $o.GoCode }}

{{ end }}

{{ range $_, $s := .ShapeList }}
	{{- if and (and (not $s.UsedAsOutput) (not $s.UsedAsInput)) (eq $s.Type "structure") }}
		{{ $s.GoCode }}
	{{- end }}
{{ end }}

{{ range $_, $s := .ShapeList }}
{{ if $s.IsEnum }}{{ $s.GoCode }}{{ end }}

{{ end }}
`))

// AddImport adds the import path to the generated file's import.
func (a *API) AddImport(v string) error {
	a.imports[v] = true
	return nil
}

// AddSDKImport adds a SDK package import to the generated file's import.
func (a *API) AddSDKImport(v ...string) error {
	e := make([]string, 0, 5)
	e = append(e, SDKImportRoot)
	e = append(e, v...)

	a.imports[path.Join(e...)] = true
	return nil
}

// APIGoCode renders the API in Go code. Returning it as a string
func (a *API) APIGoCode() string {
	a.resetImports()
	a.AddSDKImport("aws")
	a.AddSDKImport("internal/awsutil")
	a.AddImport("context")

	if a.OperationHasOutputPlaceholder() {
		a.AddSDKImport("private/protocol", a.ProtocolPackage())
		a.AddSDKImport("private/protocol")
	}
	if !a.NoGenMarshalers || !a.NoGenUnmarshalers {
		a.AddSDKImport("private/protocol")
	}

	var buf bytes.Buffer
	err := tplAPI.Execute(&buf, a)
	if err != nil {
		panic(err)
	}

	code := a.importsGoCode() + strings.TrimSpace(buf.String())
	return code
}

// APIOperationGoCode renders the Operation's in Go code. Returning it as a
// string.
func (a *API) APIOperationGoCode(op *Operation) string {
	a.resetImports()
	a.AddSDKImport("aws")
	a.AddSDKImport("internal/awsutil")
	a.AddImport("context")

	if op.OutputRef.Shape.Placeholder {
		a.AddSDKImport("private/protocol", a.ProtocolPackage())
		a.AddSDKImport("private/protocol")
	}
	if !a.NoGenMarshalers || !a.NoGenUnmarshalers {
		a.AddSDKImport("private/protocol")
	}

	// Need to generate code before imports are generated.
	code := op.GoCode()
	return a.importsGoCode() + code
}

// APIEnumsGoCode renders the API's enumerations in Go code. Returning them as
// a string.
func (a *API) APIEnumsGoCode() string {
	a.resetImports()

	var code strings.Builder
	code.WriteString(a.importsGoCode())

	for _, s := range a.ShapeEnumList() {
		code.WriteString(s.GoCode())
		fmt.Fprintln(&code)
	}

	return code.String()
}

// A tplAPIShapes is the top level template for the API Shapes.
var tplAPIShapes = template.Must(template.New("api").Parse(`
{{ range $_, $s := .ShapeList }}
	{{ if and (and (not $s.UsedAsInput) (not $s.UsedAsOutput)) (eq $s.Type "structure") -}}
		{{ $s.GoCode }}
	{{- end }}
{{ end }}
`))

// APIParamShapesGoCode renders the API's shape types in Go code. Returning
// them as a string.
func (a *API) APIParamShapesGoCode() string {
	a.resetImports()
	a.AddSDKImport("aws")
	a.AddSDKImport("internal/awsutil")

	if (!a.NoGenMarshalers || !a.NoGenUnmarshalers) && (a.hasNonIOShapes()) {
		a.AddSDKImport("private/protocol")
	}

	var buf bytes.Buffer
	err := tplAPIShapes.Execute(&buf, a)
	if err != nil {
		panic(err)
	}

	// TODO this is hacky, imports should only be added when needed.
	importStubs := `
	var _ aws.Config
	var _ = awsutil.Prettify
	`

	code := a.importsGoCode() + importStubs + strings.TrimSpace(buf.String())
	return code
}

var noCrossLinkServices = map[string]struct{}{
	"apigateway":        {},
	"budgets":           {},
	"cloudsearch":       {},
	"cloudsearchdomain": {},
	"elastictranscoder": {},
	"es":                {},
	"glacier":           {},
	"importexport":      {},
	"iot":               {},
	"iot-data":          {},
	"machinelearning":   {},
	"rekognition":       {},
	"sdb":               {},
	"swf":               {},
}

// GetCrosslinkURL returns the crosslinking URL for the shape based on the name and
// uid provided. Empty string is returned if no crosslink link could be determined.
func GetCrosslinkURL(baseURL, uid string, params ...string) string {
	if uid == "" || baseURL == "" {
		return ""
	}

	id := crosslinkServiceIDFromUID(uid)
	if _, ok := noCrossLinkServices[strings.ToLower(id)]; ok {
		return ""
	}

	return strings.Join(append([]string{baseURL, "goto", "WebAPI", uid}, params...), "/")
}

// crosslinkServiceIDFromUID will parse the service id from the uid and return
// the service id that was found.
func crosslinkServiceIDFromUID(uid string) string {
	found := 0
	i := len(uid) - 1
	for ; i >= 0; i-- {
		if uid[i] == '-' {
			found++
		}
		// Terminate after the date component is found, e.g. es-2017-11-11
		if found == 3 {
			break
		}
	}

	return uid[0:i]
}

var serviceIDRegex = regexp.MustCompile("[^a-zA-Z0-9 ]+")
var prefixDigitRegex = regexp.MustCompile("^[0-9]+")

// RawServiceID will return a unique identifier specific to a service.
func (a *API) RawServiceID() string {
	if len(a.Metadata.ServiceID) > 0 {
		return a.Metadata.ServiceID
	}

	name := a.NiceName()

	name = strings.Replace(name, "Amazon", "", -1)
	name = strings.Replace(name, "AWS", "", -1)
	name = serviceIDRegex.ReplaceAllString(name, "")
	name = prefixDigitRegex.ReplaceAllString(name, "")
	name = strings.TrimSpace(name)

	return name
}

// ServiceID returns the symbolized service identifier for the API model.
func (a *API) ServiceID() string {
	if len(a.serviceID) != 0 {
		return a.serviceID
	}

	name := a.RawServiceID()
	if len(name) == 0 {
		name = a.NiceName()
	}

	// Strip out prefix names not reflected in service client symbol names.
	for _, prefix := range stripServiceNamePrefixes {
		if strings.HasPrefix(name, prefix) {
			name = name[len(prefix):]
			break
		}
	}

	// Replace all Non-letter/number values with space
	runes := []rune(name)
	for i := 0; i < len(runes); i++ {
		if r := runes[i]; !(unicode.IsNumber(r) || unicode.IsLetter(r)) {
			runes[i] = ' '
		}
	}
	name = string(runes)

	// Title case name so its readable as a symbol.
	name = strings.Title(name)

	// Strip out spaces.
	name = strings.Replace(name, " ", "", -1)

	a.serviceID = name
	return a.serviceID
}

var tplServiceDoc = template.Must(template.New("service docs").Funcs(template.FuncMap{
	"GetCrosslinkURL": GetCrosslinkURL,
}).
	Parse(`
// Package {{ .PackageName }} provides the client and types for making API
// requests to {{ .NiceName }}.
{{ if .Documentation -}}
//
{{ .Documentation }}
{{ end -}}
{{ $crosslinkURL := GetCrosslinkURL $.BaseCrosslinkURL $.Metadata.UID -}}
{{ if $crosslinkURL -}}
//
// See {{ $crosslinkURL }} for more information on this service.
{{ end -}}
//
// See {{ .PackageName }} package documentation for more information.
// https://docs.aws.amazon.com/sdk-for-go/api/service/{{ .PackageName }}/
//
// Using the Client
//
// To use {{ .NiceName }} with the SDK use the New function to create
// a new service client. With that client you can make API requests to the service.
// These clients are safe to use concurrently.
//
// See the SDK's documentation for more information on how to use the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/
// 
// See aws.Config documentation for more information on configuring SDK clients.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the {{ .NiceName }} client for more information on 
// creating client for this service.
// https://docs.aws.amazon.com/sdk-for-go/api/service/{{ .PackageName }}/#New
`))

// A tplService defines the template for the service generated code.
var tplService = template.Must(template.New("service").Funcs(template.FuncMap{
	"ServiceNameValue": func(a *API) string {
		return fmt.Sprintf("%q", a.NiceName())
	},
	"ServiceNameConst": func(a *API) string {
		if a.NoConstServiceNames {
			return fmt.Sprintf("%q", a.NiceName())
		}
		return "ServiceName"
	},
	"ServiceIDValue": func(a *API) string {
		return fmt.Sprintf("%q", a.ServiceID())
	},
	"ServiceIDConst": func(a *API) string {
		if a.NoConstServiceNames {
			return fmt.Sprintf("%q", a.ServiceID())
		}
		return "ServiceID"
	},
	"EndpointsIDValue": func(a *API) string {
		return fmt.Sprintf("%q", a.Metadata.EndpointsID)
	},
	"EndpointsIDConst": func(a *API) string {
		if a.NoConstServiceNames {
			return fmt.Sprintf("%q", a.Metadata.EndpointsID)
		}
		return "EndpointsID"
	},
	"SigningName": func(a *API) string {
		if v := a.Metadata.SigningName; len(v) != 0 {
			return v
		}
		return a.Metadata.EndpointsID
	},
	"ServiceSpecificConfig": func(a *API) string {
		cfgs, ok := serviceSpecificConfigs[a.ServiceID()]
		if !ok || len(cfgs) == 0 {
			return ""
		}

		return cfgs.GoCode()
	},
}).Parse(`
// {{ .StructName }} provides the API operation methods for making requests to
// {{ .NiceName }}. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type {{ .StructName }} struct {
	*aws.Client
	{{ ServiceSpecificConfig . -}}
}

{{ if .UseInitMethods -}}
// Used for custom client initialization logic
var initClient func(*{{ .StructName }})

// Used for custom request initialization logic
var initRequest func(*{{ .StructName }}, *aws.Request)
{{ end }}

{{ if not .NoConstServiceNames -}}
const (
	ServiceName = {{ ServiceNameValue . }} // Service's name
	ServiceID   = {{ ServiceIDValue . }}   // Service's identifier
	EndpointsID = {{ EndpointsIDValue . }} // Service's Endpoint identifier
)
{{ end }}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := {{ .PackageName }}.New(myConfig)
func New(config aws.Config) *{{ .StructName }} {
    svc := &{{ .StructName }}{
    	Client: aws.NewClient(
    		config,
    		aws.Metadata{
				ServiceName:   {{ ServiceNameConst . }},
				ServiceID:     {{ ServiceIDConst . }},
				EndpointsID:   {{ EndpointsIDConst . }},
				SigningName:   "{{ SigningName . }}",
				SigningRegion: config.Region,
				APIVersion:    "{{ .Metadata.APIVersion }}",
				{{ if eq .Metadata.Protocol "json" -}}
					{{ if .Metadata.JSONVersion -}}
						JSONVersion:  "{{ .Metadata.JSONVersion }}",
					{{- end }}
					{{ if .Metadata.TargetPrefix -}}
						TargetPrefix: "{{ .Metadata.TargetPrefix }}",
					{{- end }}
				{{- end }}
    		},
    	),
    }

	// Handlers
	{{- if eq .Metadata.SignatureVersion "v2" }}
		svc.Handlers.Sign.PushBackNamed(v2.SignRequestHandler)
		svc.Handlers.Sign.PushBackNamed(defaults.BuildContentLengthHandler)
	{{- else }}
		svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	{{- end }}
	svc.Handlers.Build.PushBackNamed({{ .ProtocolPackage }}.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed({{ .ProtocolPackage }}.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed({{ .ProtocolPackage }}.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed({{ .ProtocolPackage }}.UnmarshalErrorHandler)

	{{ if .UseInitMethods -}}
		// Run custom client initialization if present
		if initClient != nil {
			initClient(svc)
		}
	{{ end }}

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *{{ .StructName }}) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	{{ if .UseInitMethods }}// Run custom request initialization if present
	if initRequest != nil {
		initRequest(c, req)
	}
	{{ end }}

	return req
}
`))

// ServicePackageDoc generates the contents of the doc file for the service.
//
// Will also read in the custom doc templates for the service if found.
func (a *API) ServicePackageDoc() string {
	a.imports = map[string]bool{}

	var buf bytes.Buffer
	if err := tplServiceDoc.Execute(&buf, a); err != nil {
		panic(err)
	}

	return buf.String()
}

// ServiceGoCode renders service go code. Returning it as a string.
func (a *API) ServiceGoCode() string {
	a.resetImports()
	a.AddSDKImport("aws")
	if a.Metadata.SignatureVersion == "v2" {
		a.AddSDKImport("aws/signer/v2")
		a.AddSDKImport("aws/defaults")
	} else {
		a.AddSDKImport("aws/signer/v4")
	}
	a.AddSDKImport("private/protocol", a.ProtocolPackage())

	var buf bytes.Buffer
	err := tplService.Execute(&buf, a)
	if err != nil {
		panic(err)
	}

	code := a.importsGoCode() + buf.String()
	return code
}

// ExampleGoCode renders service example code. Returning it as a string.
func (a *API) ExampleGoCode() string {
	exs := []string{}
	imports := map[string]bool{}
	for _, o := range a.OperationList() {
		o.imports = map[string]bool{}
		exs = append(exs, o.Example())
		for k, v := range o.imports {
			imports[k] = v
		}
	}

	code := fmt.Sprintf("import (\n%q\n%q\n%q\n\n%q\n%q\n",
		"bytes",
		"fmt",
		"time",
		SDKImportRoot+"/aws",
		a.ImportPath(),
	)
	for k := range imports {
		code += fmt.Sprintf("%q\n", k)
	}
	code += ")\n\n"
	code += "var _ time.Duration\nvar _ bytes.Buffer\n\n"
	code += strings.Join(exs, "\n\n")
	return code
}

// A tplInterface defines the template for the service interface type.
var tplInterface = template.Must(template.New("interface").Parse(`
// {{ .StructName }}API provides an interface to enable mocking the
// {{ .PackageName }}.{{ .StructName }} methods. This make unit testing your code that
// calls out to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the SDK's request pipeline.
//
//    // myFunc uses an SDK service client to make a request to
//    // {{ .NiceName }}. {{ $opts := .OperationList }}{{ $opt := index $opts 0 }}
//    func myFunc(svc {{ .InterfacePackageName }}.{{ .StructName }}API) bool {
//        // Make svc.{{ $opt.ExportedName }} request
//    }
//
//    func main() {
//        cfg, err := external.LoadDefaultAWSConfig()
//        if err != nil {
//            panic("failed to load config, " + err.Error())
//        }
//
//        svc := {{ .PackageName }}.New(cfg)
//
//        myFunc(svc)
//    }
//
// In your _test.go file:
//
//    // Define a mock struct to be used in your unit tests of myFunc.
//    type mock{{ .StructName }}Client struct {
//        {{ .InterfacePackageName }}.{{ .StructName }}PI
//    }
//    func (m *mock{{ .StructName }}Client) {{ $opt.ExportedName }}(input {{ $opt.InputRef.GoTypeWithPkgName }}) ({{ $opt.OutputRef.GoTypeWithPkgName }}, error) {
//        // mock response/functionality
//    }
//
//    func TestMyFunc(t *testing.T) {
//        // Setup Test
//        mockSvc := &mock{{ .StructName }}Client{}
//
//        myfunc(mockSvc)
//
//        // Verify myFunc's functionality
//    }
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters. Its suggested to use the pattern above for testing, or using 
// tooling to generate mocks to satisfy the interfaces.
type {{ .StructName }}API interface {
    {{ range $_, $o := .OperationList }}
        {{ $o.InterfaceSignature }}
    {{ end }}
    {{ range $_, $w := .Waiters }}
        {{ $w.InterfaceSignature }}
    {{ end }}
}

var _ {{ .StructName }}API = (*{{ .PackageName }}.{{ .StructName }})(nil)
`))

// InterfaceGoCode returns the go code for the service's API operations as an
// interface{}. Assumes that the interface is being created in a different
// package than the service API's package.
func (a *API) InterfaceGoCode() string {
	a.resetImports()
	if len(a.Waiters) != 0 {
		a.AddSDKImport("aws")
	}
	a.AddImport(a.ImportPath())

	var buf bytes.Buffer
	err := tplInterface.Execute(&buf, a)

	if err != nil {
		panic(err)
	}

	code := a.importsGoCode() + strings.TrimSpace(buf.String())
	return code
}

// NewAPIGoCodeWithPkgName returns a string of instantiating the API prefixed
// with its package name. Takes a string depicting the Config.
func (a *API) NewAPIGoCodeWithPkgName(cfg string) string {
	return fmt.Sprintf("%s.New(%s)", a.PackageName(), cfg)
}

// computes the validation chain for all input shapes
func (a *API) addShapeValidations() {
	for _, o := range a.Operations {
		resolveShapeValidations(o.InputRef.Shape)
	}
}

// Updates the source shape and all nested shapes with the validations that
// could possibly be needed.
func resolveShapeValidations(s *Shape, ancestry ...*Shape) {
	for _, a := range ancestry {
		if a == s {
			return
		}
	}

	children := []string{}
	for _, name := range s.MemberNames() {
		ref := s.MemberRefs[name]

		if s.IsRequired(name) && !s.Validations.Has(ref, ShapeValidationRequired) {
			s.Validations = append(s.Validations, ShapeValidation{
				Name: name, Ref: ref, Type: ShapeValidationRequired,
			})
		}

		if ref.Shape.Min != 0 && !s.Validations.Has(ref, ShapeValidationMinVal) {
			s.Validations = append(s.Validations, ShapeValidation{
				Name: name, Ref: ref, Type: ShapeValidationMinVal,
			})
		}

		if !ref.CanBeEmpty() && !s.Validations.Has(ref, ShapeValidationMinVal) {
			s.Validations = append(s.Validations, ShapeValidation{
				Name: name, Ref: ref, Type: ShapeValidationMinVal,
			})
		}

		switch ref.Shape.Type {
		case "map", "list", "structure":
			children = append(children, name)
		}
	}

	ancestry = append(ancestry, s)
	for _, name := range children {
		ref := s.MemberRefs[name]
		// Since this is a grab bag we will just continue since
		// we can't validate because we don't know the valued shape.
		if ref.JSONValue {
			continue
		}

		nestedShape := ref.Shape.NestedShape()

		var v *ShapeValidation
		if len(nestedShape.Validations) > 0 {
			v = &ShapeValidation{
				Name: name, Ref: ref, Type: ShapeValidationNested,
			}
		} else {
			resolveShapeValidations(nestedShape, ancestry...)
			if len(nestedShape.Validations) > 0 {
				v = &ShapeValidation{
					Name: name, Ref: ref, Type: ShapeValidationNested,
				}
			}
		}

		if v != nil && !s.Validations.Has(v.Ref, v.Type) {
			s.Validations = append(s.Validations, *v)
		}
	}
	ancestry = ancestry[:len(ancestry)-1]
}

// A tplAPIErrors is the top level template for the API
var tplAPIErrors = template.Must(template.New("api").Parse(`
const (
{{ range $_, $s := $.ShapeListErrors }}
	// {{ $s.ErrorCodeName }} for service response error code
	// {{ printf "%q" $s.ErrorName }}.
	{{ if $s.Docstring -}}
	//
	{{ $s.Docstring }}
	{{ end -}}
	{{ $s.ErrorCodeName }} = {{ printf "%q" $s.ErrorName }}
{{ end }}
)
`))

func (a *API) APIErrorsGoCode() string {
	var buf bytes.Buffer
	err := tplAPIErrors.Execute(&buf, a)

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}

// removeOperation removes an operation, its input/output shapes, as well as
// any references/shapes that are unique to this operation.
func (a *API) removeOperation(name string) {
	debugLogger.Logln("removing operation,", name)
	op := a.Operations[name]

	delete(a.Operations, name)
	delete(a.Examples, name)

	a.removeShape(op.InputRef.Shape)
	a.removeShape(op.OutputRef.Shape)
}

// removeShape removes the given shape, and all form member's reference target
// shapes. Will also remove member reference targeted shapes if those shapes do
// not have any additional references.
func (a *API) removeShape(s *Shape) {
	debugLogger.Logln("removing shape,", s.ShapeName)

	delete(a.Shapes, s.ShapeName)

	for name, ref := range s.MemberRefs {
		a.removeShapeRef(ref)
		delete(s.MemberRefs, name)
	}

	for _, ref := range []*ShapeRef{&s.MemberRef, &s.KeyRef, &s.ValueRef} {
		if ref.Shape == nil {
			continue
		}
		a.removeShapeRef(ref)
		*ref = ShapeRef{}
	}
}

// removeShapeRef removes the shape reference from its target shape. If the
// reference was the last reference to the target shape, the shape will also be
// removed.
func (a *API) removeShapeRef(ref *ShapeRef) {
	if ref.Shape == nil {
		return
	}

	ref.Shape.removeRef(ref)
	if len(ref.Shape.refs) == 0 {
		a.removeShape(ref.Shape)
	}
}

func (a *API) hasNonIOShapes() bool {
	for _, s := range a.Shapes {
		if s.IsError || s.Type != "structure" {
			continue
		}
		if !(s.UsedAsInput || s.UsedAsOutput) {
			return true
		}
	}
	return false
}
