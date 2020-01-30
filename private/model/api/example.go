// +build codegen

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/private/util"
)

// ExamplesGoCode will return a code representation of the entry within the
// examples.json file.
func (a *API) ExamplesGoCode() string {
	a.resetImports()
	a.AddImport("context")
	a.AddImport("fmt")

	a.AddSDKImport("aws")
	a.AddSDKImport("aws/awserr")
	a.AddSDKImport("aws/external")
	a.AddImport(a.ImportPath())

	var extra = `
var _ aws.Config
`

	code := a.Examples.GoCode()
	if len(code) == 0 {
		return ""
	}
	if a.Examples.UsedTimestamp() {
		a.AddImport("time")
		extra += `
func parseTime(layout, value string) time.Time {
   t, err := time.Parse(layout, value)
   if err != nil {
	   panic(err)
   }
   return t
}
`
	}

	return a.importsGoCode() + extra + code
}

// ExamplesDefinition is the structural representation of the examples-1.json file
type ExamplesDefinition struct {
	*API     `json:"-"`
	Examples Examples `json:"examples"`
}

// Example is a single entry within the examples-1.json file.
type Example struct {
	API           *API                   `json:"-"`
	Operation     *Operation             `json:"-"`
	OperationName string                 `json:"-"`
	Index         string                 `json:"-"`
	Builder       *ExamplesBuilder       `json:"-"`
	VisitedErrors map[string]struct{}    `json:"-"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	ID            string                 `json:"id"`
	Comments      Comments               `json:"comments"`
	Input         map[string]interface{} `json:"input"`
	Output        map[string]interface{} `json:"output"`
}

type Comments struct {
	Input  map[string]interface{} `json:"input"`
	Output map[string]interface{} `json:"output"`
}

var exampleFuncMap = template.FuncMap{
	"commentify":           commentify,
	"wrap":                 wrap,
	"generateExampleInput": generateExampleInput,
	"generateTypes":        generateTypes,
}

var exampleCustomizations = map[string]template.FuncMap{}

var exampleTmpls = template.Must(template.New("example").Funcs(exampleFuncMap).Parse(`
{{ generateTypes . }}
{{ commentify (wrap .Title 80) }}
//
{{ commentify (wrap .Description 80) }}
func Example{{ .API.StructName }}_{{ .MethodName }}() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("failed to load config, " + err.Error())
	}

	svc := {{ .API.PackageName }}.New(cfg)
	input := {{ generateExampleInput . }}
	req := svc.{{ .OperationName }}Request(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
				{{ range $_, $ref := .Operation.ErrorRefs -}}
					{{ if not ($.HasVisitedError $ref) -}}
			case {{ .API.PackageName }}.{{ $ref.Shape.ErrorCodeName }}:
				fmt.Println({{ .API.PackageName }}.{{ $ref.Shape.ErrorCodeName }}, aerr.Error())
					{{ end -}}
				{{ end -}}
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
`))

type Examples map[string][]Example

func (exs Examples) UsedTimestamp() bool {
	for _, es := range exs {
		for _, e := range es {
			if e.Builder.HasTimestamp {
				return true
			}
		}
	}

	return false
}

// Names will return the name of the example. This will also be the name of the operation
// that is to be tested.
func (exs Examples) Names() []string {
	names := make([]string, 0, len(exs))
	for k := range exs {
		names = append(names, k)
	}

	sort.Strings(names)
	return names
}

func (exs Examples) GoCode() string {
	buf := bytes.NewBuffer(nil)
	for _, opName := range exs.Names() {
		examples := exs[opName]
		for _, ex := range examples {
			buf.WriteString(util.GoFmt(ex.GoCode()))
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

// ExampleCode will generate the example code for the given Example shape.
// TODO: Can delete
func (ex Example) GoCode() string {
	var buf bytes.Buffer
	m := exampleFuncMap
	if fMap, ok := exampleCustomizations[ex.API.PackageName()]; ok {
		m = fMap
	}
	tmpl := exampleTmpls.Funcs(m)
	if err := tmpl.ExecuteTemplate(&buf, "example", &ex); err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}

func generateExampleInput(ex Example) string {
	if ex.Operation.HasInput() {
		return fmt.Sprintf("&%s{\n%s\n}\n",
			ex.Builder.GoType(&ex.Operation.InputRef, true),
			ex.Builder.BuildShape(&ex.Operation.InputRef, ex.Input, false, false),
		)
	}
	return ""
}

// generateTypes will generate no types for default examples, but customizations may
// require their own defined types.
func generateTypes(ex Example) string {
	return ""
}

// correctType will cast the value to the correct type when printing the string.
// This is due to the json decoder choosing numbers to be floats, but the shape may
// actually be an int. To counter this, we pass the shape's type and properly do the
// casting here.
func correctType(memName string, t string, value interface{}, asValue bool) string {
	if value == nil {
		return ""
	}

	v := ""
	switch value.(type) {
	case string:
		v = value.(string)
	case int:
		v = fmt.Sprintf("%d", value.(int))
	case float64:
		if t == "integer" || t == "long" || t == "int64" {
			v = fmt.Sprintf("%d", int(value.(float64)))
		} else {
			v = fmt.Sprintf("%f", value.(float64))
		}
	case bool:
		v = fmt.Sprintf("%t", value.(bool))
	}

	return convertToCorrectType(memName, t, v, asValue)
}

func getValue(t, v string, asValue bool) string {
	if t[0] == '*' {
		t = t[1:]
	}

	if asValue {
		switch t {
		case "string":
			return fmt.Sprintf("%q", v)
		case "integer", "long", "int64", "float", "float64", "double", "boolean":
			return fmt.Sprintf("%s", v)
		default:
			panic("Unsupported type: " + t)
		}
	}

	switch t {
	case "string":
		return fmt.Sprintf("aws.String(%q)", v)
	case "integer", "long", "int64":
		return fmt.Sprintf("aws.Int64(%s)", v)
	case "float", "float64", "double":
		return fmt.Sprintf("aws.Float64(%s)", v)
	case "boolean":
		return fmt.Sprintf("aws.Bool(%s)", v)
	default:
		panic("Unsupported type: " + t)
	}
}

// AttachExamples will create a new ExamplesDefinition from the examples file
// and reference the API object.
func (a *API) AttachExamples(filename string) error {
	p := ExamplesDefinition{API: a}

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	err = json.NewDecoder(f).Decode(&p)
	if err != nil {
		return fmt.Errorf("failed to decode %s, err: %v", filename, err)
	}

	return p.setup()
}

func (p *ExamplesDefinition) setup() error {
	builder := NewExamplesBuilder()
	keys := p.Examples.Names()
	for _, n := range keys {
		examples := p.Examples[n]
		for i, e := range examples {
			n = p.ExportableName(n)
			e.OperationName = n
			e.API = p.API
			e.Index = fmt.Sprintf("shared%02d", i)

			e.Builder = builder

			e.VisitedErrors = map[string]struct{}{}
			op := p.API.Operations[e.OperationName]
			e.OperationName = p.ExportableName(e.OperationName)
			e.Operation = op
			p.Examples[n][i] = e
		}
	}

	p.API.Examples = p.Examples

	return nil
}

func (ex *Example) HasVisitedError(errRef *ShapeRef) bool {
	errName := errRef.Shape.ErrorCodeName()
	_, ok := ex.VisitedErrors[errName]
	ex.VisitedErrors[errName] = struct{}{}
	return ok
}

func (ex *Example) MethodName() string {
	return fmt.Sprintf("%sRequest_%s", ex.OperationName, ex.Index)
}
