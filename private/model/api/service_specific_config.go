// +build codegen

package api

import (
	"bytes"
	"text/template"
)

type serviceConfigField struct {
	Name string
	Doc  string
	Type string
}

type serviceConfigFields []serviceConfigField

func (fs serviceConfigFields) GoCode() string {
	w := bytes.NewBuffer(nil)

	if err := svcConfigTmpl.Execute(w, fs); err != nil {
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

var svcConfigTmpl = template.Must(template.New("svcConfigTmpl").
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
