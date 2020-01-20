// +build codegen

package api

import (
	"bytes"
	"text/template"
)

type externalConfigBinding struct {
	ResolverMethodName string

	MustResolveEnvConfig    bool
	MustResolveSharedConfig bool
}

type serviceConfigField struct {
	Name string
	Doc  string
	Type string

	ExternalConfigMethod *externalConfigBinding
}

// HasExternalConfigBinding returns whether the service config field binds to an externally resolve
// configuration property.
func (s serviceConfigField) HasExternalConfigBinding() bool {
	return s.ExternalConfigMethod != nil
}

// ExternalConfigInterfaceName returns the name of the external service config method name
func (s serviceConfigField) ExternalConfigInterfaceName() string {
	return s.Name + "Provider"
}

// ExternalConfigResolverName returns the name of the external service config resolver
func (s serviceConfigField) ExternalConfigResolverName() string {
	return "Resolve" + s.Name
}

type serviceConfigFields []serviceConfigField

// ServiceClientGoCode returns rendered service configuration fields suitable for a structure definition
func (fs serviceConfigFields) ServiceClientGoCode() string {
	w := bytes.NewBuffer(nil)

	if err := tplServiceConfigFields.Execute(w, fs); err != nil {
		panic("failed to render serviceConfigFields " + err.Error())
	}

	return w.String()
}

var serviceSpecificConfigs = map[string]serviceConfigFields{
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
			ExternalConfigMethod: &externalConfigBinding{
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

var tplServiceConfigFields = template.Must(template.New("tplServiceConfigFields").
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

var tplExtServiceConfigResolvers = template.Must(template.New("tplExtServiceConfigResolvers").Funcs(
	map[string]interface{}{
		"ExternalConfigFields": func(a *API) serviceConfigFields {
			if !a.HasExternalServiceConfigFields() {
				return nil
			}

			var fields serviceConfigFields
			for _, field := range serviceSpecificConfigs[a.ServiceID()] {
				if field.HasExternalConfigBinding() {
					fields = append(fields, field)
				}
			}

			return fields
		},
	}).Parse(`
{{ define "resolvers" }}
	{{- $fields := ExternalConfigFields . -}}
	{{- range $_, $field := $fields }}
		// {{ $field.ExternalConfigInterfaceName }} is an interface for retrieving external configuration value for {{ $field.Name }}
		type {{ $field.ExternalConfigInterfaceName }} interface {
			{{ $field.ExternalConfigMethod.ResolverMethodName }}() (value {{ $field.Type }}, ok bool, err error)
		}

		func {{ $field.ExternalConfigResolverName }}(configs []interface{}) (value {{ $field.Type }}, ok bool, err error) {
			for _, cfg := range configs {
				if p, pOk := cfg.({{ $field.ExternalConfigInterfaceName }}); pOk {
					value, ok, err = p.{{ $field.ExternalConfigMethod.ResolverMethodName }}()
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
		// {{ $field.ExternalConfigInterfaceName }} Assertions
		var (
			{{- if $field.ExternalConfigMethod.MustResolveEnvConfig }}
				_ {{ $.ExternalConfigPackageName }}.{{ $field.ExternalConfigInterfaceName }} = &external.EnvConfig{}
			{{- end }}
			{{- if $field.ExternalConfigMethod.MustResolveSharedConfig }}
				_ {{ $.ExternalConfigPackageName }}.{{ $field.ExternalConfigInterfaceName }} = &external.SharedConfig{}
			{{- end }}
		)
	{{ end }}
{{ end }}
`))
