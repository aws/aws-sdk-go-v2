// +build codegen

package api

import (
	"bytes"
	"fmt"
	"text/template"
)

type externalConfigBinding struct {
	ResolverMethodName string

	MustResolveEnvConfig    bool
	MustResolveSharedConfig bool
}

type clientConfigField struct {
	Name string
	Doc  string
	Type string

	ExternalConfigBinding *externalConfigBinding
}

// HasExternalConfigBinding returns whether the service config field binds to an externally resolve
// configuration property.
func (s clientConfigField) HasExternalConfigBinding() bool {
	return s.ExternalConfigBinding != nil
}

// ExternalConfigProviderName returns the name of the external service config method name
func (s clientConfigField) ExternalConfigProviderName() string {
	return s.Name + "Provider"
}

// ExternalConfigResolverName returns the name of the external service config resolver
func (s clientConfigField) ExternalConfigResolverName() string {
	return "Resolve" + s.Name
}

// ExternalConfigProviderSignature returns the signature of the external config resolver method
func (s clientConfigField) ExternalConfigProviderSignature() string {
	return fmt.Sprintf("%s() (value %s, ok bool, err error)", s.ExternalConfigBinding.ResolverMethodName, s.Type)
}

type clientConfigFields []clientConfigField

// ClientGoCode returns rendered service configuration fields suitable for a structure definition
func (fs clientConfigFields) ClientGoCode() string {
	w := bytes.NewBuffer(nil)

	if err := tplClientConfigFields.Execute(w, fs); err != nil {
		panic("failed to render clientConfigFields " + err.Error())
	}

	return w.String()
}

var clientSpecificConfigs = map[string]clientConfigFields{
	"S3": {
		{
			Name: "Disable100Continue", Type: "bool",
			Doc: `Disables the S3 client from using the Expect: 100-Continue header to wait for
the service to respond with a 100 status code before sending the HTTP request
body.

You should disable 100-Continue if you experience issues with proxies or third
party S3 compatible services.

See http.Transport's ExpectContinueTimeout for information on adjusting the
continue wait timeout. https://golang.org/pkg/net/http/#Transport`,
		},
		{
			Name: "ForcePathStyle", Type: "bool",
			Doc: `Forces the client to use path-style addressing for S3 API operations. By
default the S3 client will use virtual hosted bucket addressing when possible.
The S3 client will automatically fall back to path-style when the bucket name
is not DNS compatible.

With ForcePathStyle 

	https://s3.us-west-2.amazonaws.com/BUCKET/KEY

Without ForcePathStyle

	https://BUCKET.s3.us-west-2.amazonaws.com/KEY

See http://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html`,
		},
		{
			Name: "UseAccelerate", Type: "bool",
			Doc: `Enables S3 Accelerate feature for API operation that support S3 Accelerate.
For all operations compatible with S3 Accelerate will use the accelerate
endpoint for requests. Requests not compatible will fall back to normal S3
requests.

The bucket must be enable for accelerate to be used with S3 client with
accelerate enabled. If the bucket is not enabled for accelerate an error will
be returned. The bucket name must be DNS compatible to also work with
accelerate.
			`,
		},
		{
			Name: "UseARNRegion", Type: "bool",
			Doc: `Set this to ` + "`true`" + ` to use the region specified
in the ARN, when an ARN is provided as an argument to a bucket parameter.`,
			ExternalConfigBinding: &externalConfigBinding{
				ResolverMethodName:      "GetS3UseARNRegion",
				MustResolveEnvConfig:    true,
				MustResolveSharedConfig: true,
			},
		},
	},

	"DynamoDB": {
		{
			Name: "DisableComputeChecksums", Type: "bool",
			Doc: `Disables the computation and validation of request and response checksums.`,
		},
	},

	"SQS": {
		{
			Name: "DisableComputeChecksums", Type: "bool",
			Doc: `Disables the computation and validation of request and response checksums.`,
		},
	},
}

var tplClientConfigFields = template.Must(template.New("tplClientConfigFields").
	Funcs(template.FuncMap{
		"Commentify": commentify,
	}).
	Parse(`
// Service specific configurations. (codegen: service_specific_config.go)
{{ range $i, $field := $ }}
	{{ Commentify $field.Doc }}
	{{ $field.Name }} {{ $field.Type }}
{{ end }}
`))

var tplExternalConfigResolvers = template.Must(template.New("tplExternalConfigResolvers").Funcs(
	map[string]interface{}{
		"ExternalConfigFields": externalConfigFields,
	}).Parse(`
{{ define "resolvers" }}
	{{- $fields := ExternalConfigFields . -}}
	{{- range $_, $field := $fields }}
		// {{ $field.ExternalConfigProviderName }} is an interface for retrieving external configuration value for {{ $field.Name }}
		type {{ $field.ExternalConfigProviderName }} interface {
			{{ $field.ExternalConfigProviderSignature }}
		}

		// {{ $field.ExternalConfigResolverName }} extracts the first instance of a {{ $field.Name }} from the config slice.
		// Additionally returns a boolean to indicate if the value was found in provided configs, and error if one is encountered.
		func {{ $field.ExternalConfigResolverName }}(configs []interface{}) (value {{ $field.Type }}, ok bool, err error) {
			for _, cfg := range configs {
				if p, pOk := cfg.({{ $field.ExternalConfigProviderName }}); pOk {
					value, ok, err = p.{{ $field.ExternalConfigBinding.ResolverMethodName }}()
					if err != nil {
						return value, false, err
					}
					if ok {
						break
					}
				}
			}

			return value, ok, err
		}
	{{ end }}
{{ end }}
{{ define "tests" }}
	{{- $fields := ExternalConfigFields . -}}
	{{- range $i, $field := $fields }}
		// {{ $field.ExternalConfigProviderName }} Assertions
		var (
			{{- if $field.ExternalConfigBinding.MustResolveEnvConfig }}
				_ svcExternal.{{ $field.ExternalConfigProviderName }} = &external.EnvConfig{}
			{{- end }}
			{{- if $field.ExternalConfigBinding.MustResolveSharedConfig }}
				_ svcExternal.{{ $field.ExternalConfigProviderName }} = &external.SharedConfig{}
			{{- end }}
		)
	{{ end }}
{{ end }}
`))

func externalConfigFields(a *API) clientConfigFields {
	if !a.HasExternalClientConfigFields() {
		return nil
	}

	var fields clientConfigFields
	for _, field := range clientSpecificConfigs[a.ServiceID()] {
		if field.HasExternalConfigBinding() {
			fields = append(fields, field)
		}
	}

	return fields
}
