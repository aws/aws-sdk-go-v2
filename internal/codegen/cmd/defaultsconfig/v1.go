package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"text/template"

	codegen "github.com/aws/aws-sdk-go-v2/internal/codegen/defaults"
)

func renderRetryMode(context *generationContext, base interface{}, modifier Modifier) (string, error) {
	value, err := applyBaseModifier(base, modifier)
	if err != nil {
		return "", err
	}

	v, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("expect retryMode as string, got %T", base)
	}

	return fmt.Sprintf("aws.RetryMode(%q)", v), nil
}

func renderPtrTimeMillisecond(context *generationContext, base interface{}, modifier Modifier) (string, error) {
	context.AddImport("time", "")

	value, err := applyBaseModifier(base, modifier)
	if err != nil {
		return "", err
	}

	i, ok := value.(int64)
	if !ok {
		return "", fmt.Errorf("expect int64, got %T", base)
	}

	return fmt.Sprintf("aws.Duration(%d*time.Millisecond)", i), nil
}

type ConfigurationValue struct {
	FieldName string
	Setter    func(context *generationContext, base interface{}, modifier Modifier) (string, error)
}

var supportedConfigKeys = map[string]ConfigurationValue{
	"retryMode": {
		FieldName: "RetryMode",
		Setter:    renderRetryMode,
	},
	"connectTimeoutInMillis": {
		FieldName: "ConnectTimeout",
		Setter:    renderPtrTimeMillisecond,
	},
	"tlsNegotiationTimeoutInMillis": {
		FieldName: "TLSNegotiationTimeout",
		Setter:    renderPtrTimeMillisecond,
	},
}

// ModifierMap is a map of configuration keys to their respective modifier descriptor.
type ModifierMap map[string]Modifier

// Modifier is a union type of either a multiply, add, or override modification descriptor.
type Modifier struct {
	Multiply *json.Number
	Add      *int64
	Override interface{}
}

// IsZero returns whether the type is the zero-value.
func (m *Modifier) IsZero() bool {
	return m.Multiply == nil &&
		m.Add == nil &&
		m.Override == nil
}

// UnmarshalJSON is a custom json.Unmarshaler that enforces a modifier unmarshas to a union and satisfies
// the appropriate type constraints.
func (m *Modifier) UnmarshalJSON(data []byte) error {
	modifiers := make(map[string]interface{})

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(&modifiers); err != nil {
		return err
	}

	if len(modifiers) != 1 {
		return fmt.Errorf("expect one modifier key, got %d", len(modifiers))
	}

	var key string
	for key = range modifiers {
	}

	switch {
	case strings.EqualFold(key, "multiply"):
		var v json.Number
		v, ok := modifiers[key].(json.Number)
		if !ok {
			return fmt.Errorf("expect %T, got %T", json.Number(""), v)
		}
		m.Multiply = &v
	case strings.EqualFold(key, "add"):
		v, ok := modifiers[key].(json.Number)
		if !ok {
			return fmt.Errorf("expect %T, got %T", json.Number(""), v)
		}
		i, err := numberToInt64(v)
		if err != nil {
			return fmt.Errorf("%s is not an int64", v.String())
		}
		m.Add = &i
	case strings.EqualFold(key, "override"):
		switch v := modifiers[key].(type) {
		case json.Number:
			m.Override = v
		case string:
			m.Override = v
		default:
			return fmt.Errorf("expect JSON Number or String, got %T", v)
		}
	default:
		return fmt.Errorf("unsupported modifier key: %s", key)
	}

	return nil
}

// BaseConfiguration is the set of configuration keys and their deserialized json value.
type BaseConfiguration map[string]interface{}

// SchemaVersion is a versioned JSON document.
type SchemaVersion struct {
	Version int `json:"version"`
}

// SDKDefaultConfig is a default configuration descriptor.
type SDKDefaultConfig struct {
	Base          BaseConfiguration      `json:"base"`
	Modes         map[string]ModifierMap `json:"modes"`
	Documentation struct {
		Modes         map[string]string `json:"modes"`
		Configuration map[string]string `json:"configuration"`
	} `json:"documentation"`

	SchemaVersion
}

func numberToBigFloat(value json.Number) (*big.Float, error) {
	float, ok := (&big.Float{}).SetString(value.String())
	if !ok {
		return nil, fmt.Errorf("failed to parse number")
	}
	return float, nil
}

func numberToInt64(value json.Number) (int64, error) {
	float, err := numberToBigFloat(value)
	if err != nil {
		return 0, err
	}
	if i, accuracy := float.Int64(); accuracy == big.Exact {
		return i, nil
	}
	return 0, fmt.Errorf("failed to represent number as int64")
}

func getModeConstant(v string) string {
	return fmt.Sprintf("aws.%s", "DefaultsMode"+codegen.SymbolizeExport(v))
}

func isSupportedConfigKey(v string) (ok bool) {
	_, ok = supportedConfigKeys[v]
	return ok
}

func getConfigKeyValue(v string) ConfigurationValue {
	return supportedConfigKeys[v]
}

func applyBaseModifier(base interface{}, modifier Modifier) (v interface{}, err error) {
	switch bv := base.(type) {
	case string:
		bv, err = applyStringModifier(bv, modifier)
		if err != nil {
			return "", err
		}
		return bv, nil
	case json.Number:
		i, err := numberToInt64(bv)
		if err != nil {
			return "", err
		}
		i, err = applyInt64Modifier(i, modifier)
		if err != nil {
			return "", err
		}
		return i, nil
	default:
		return nil, fmt.Errorf("unexpected type %T", base)
	}
}

func isKeyInModifierMap(m ModifierMap, key string) bool {
	_, ok := m[key]
	return ok
}

func renderImport(i Import) (string, error) {
	if len(i.Package) == 0 {
		return "", fmt.Errorf("invalid empty package import")
	}

	if len(i.Alias) == 0 {
		return fmt.Sprintf("%q", i.Package), nil
	}

	return fmt.Sprintf("%s %q", i.Package, i.Alias), nil
}

var tmpl = template.Must(template.
	New("generate").
	Funcs(map[string]interface{}{
		"symbolizeExport":      codegen.SymbolizeExport,
		"isSupportedConfigKey": isSupportedConfigKey,
		"getConfigKeyValue":    getConfigKeyValue,
		"isKeyInModifierMap":   isKeyInModifierMap,
		"getModeConstant":      getModeConstant,
		"renderImport":         renderImport,
		"listModes": func(v map[string]ModifierMap) string {
			var modes []string
			for mode := range v {
				modes = append(modes, mode)
			}
			sort.Strings(modes)
			return strings.Join(modes, ", ")
		},
	}).
	Parse(`
{{- define "header" -}}
// Code generated by github.com/aws/aws-sdk-go-v2/internal/codegen/cmd/defaultsconfig. DO NOT EDIT.

package {{ $.PackageName }}

import (
  {{- range $key, $_ := $.Imports }}
  {{ renderImport $key }}
  {{- end }}
)
{{ end -}}
{{- define "config" }}
{{- $.AddImport "fmt" "" -}}
{{- $.AddSDKImport "aws" "" -}}

// {{ $.ResolverName }} returns the default Configuration descriptor for the given mode.
//
// Supports the following modes: {{ listModes $.Config.Modes }}
func {{ $.ResolverName }}(mode aws.DefaultsMode) (Configuration, error) {
	var mv aws.DefaultsMode
	mv.SetFromString(string(mode))

	switch mv {
	{{- range $mode, $modeModifiers := $.Config.Modes }}
	case {{ getModeConstant $mode }}:
		settings := Configuration{
			{{ range $key, $value := $.Config.Base }}
				{{- if isSupportedConfigKey $key }}
					{{- $configValue := (getConfigKeyValue $key) -}}
					{{ $configValue.FieldName }}: {{ call $configValue.Setter $ $value (index $modeModifiers $key) }},
				{{ end -}}
			{{- end }}
		}
		return settings, nil
	{{- end }}
	default:
		return Configuration{}, fmt.Errorf("unsupported defaults mode: %v", mode)
	}
}
{{ end }}
`))

func applyInt64Modifier(v int64, modifier Modifier) (int64, error) {
	if modifier.IsZero() {
		return v, nil
	}

	if modifier.Override != nil {
		jn, ok := modifier.Override.(json.Number)
		if !ok {
			return 0, fmt.Errorf("expect override to be JSON Number, got %T", modifier.Override)
		}
		i, err := numberToInt64(jn)
		if err != nil {
			return 0, err
		}
		return i, nil
	} else if modifier.Add != nil {
		return v + *modifier.Add, nil
	} else if modifier.Multiply != nil {
		fv, err := numberToBigFloat(*modifier.Multiply)
		if err != nil {
			return 0, err
		}
		fv = fv.Mul(fv, (&big.Float{}).SetInt64(v))
		if i, acc := fv.Int64(); acc == big.Exact {
			return i, nil
		}
		return 0, fmt.Errorf("failed to compute an int64 using multiply modifier")
	}
	return 0, fmt.Errorf("unexpected modifier option")
}

func applyStringModifier(s string, modifier Modifier) (string, error) {
	if modifier.IsZero() {
		return s, nil
	}

	if modifier.Override == nil {
		return "", fmt.Errorf("string only supports override modifier")
	}
	ov, ok := (modifier.Override).(string)
	if !ok {
		return "", fmt.Errorf("expect override value to be string")
	}
	return ov, nil
}

func validateStringArgument(v *string, message string) {
	if v == nil || len(*v) == 0 {
		log.Fatal(message)
	}
}
