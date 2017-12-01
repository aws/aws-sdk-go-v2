// +build codegen

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

var tplCreateServiceCase = template.Must(template.New("create-serviice").Funcs(template.FuncMap{
	"ServiceImport": serviceImport,
}).Parse(`
package main

import (
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	{{ range $_, $service := $.Packages -}}
	{{ ServiceImport $service }}
	{{- end }}
)

type service struct {
	name  string
	value reflect.Value
}

func createServices(cfg aws.Config) []service {
	{{ range $_, $service := $.Packages -}}
	{{ $.CustomConfiguration $service }}
	{{- end }}

	return []service{
		{{ range $_, $service := $.Packages -}}
		{{ $.ServiceClient $service }}
		{{- end }}
	}
}

`))

const (
	sdkPath = "github.com/aws/aws-sdk-go-v2/service"
)

type serviceInfo struct {
	clientNames map[string]string
	basePath    string
	Packages    []string
}

// getServices will crawl the basepath and any folder found will be assumed
// to be a service and stored in the serviceInfo
func (si *serviceInfo) getServices() {
	filepath.Walk(si.basePath, func(path string, info os.FileInfo, err error) error {
		if filepath.Dir(path) != si.basePath {
			return nil
		}
		if !info.IsDir() {
			return nil
		}

		serviceName := info.Name()
		si.Packages = append(si.Packages, serviceName)
		return nil
	})
}

// serviceImport will create the proper import path for the given service package
func serviceImport(p string) string {
	return fmt.Sprintf("%q\n", filepath.Join(sdkPath, p))
}

type customConfig interface {
	GoCode(*serviceInfo) string
}

var customConfigs = map[string]customConfig{
	"s3":  s3CustomConfig{},
	"sqs": sqsCustomConfig{},
}

type s3CustomConfig struct{}

func (config s3CustomConfig) GoCode(si *serviceInfo) string {
	si.clientNames["s3"] = "s3Client"
	buf := bytes.NewBuffer(nil)
	buf.WriteString("s3Client := s3.New(cfg)\n")
	buf.WriteString("s3Client.ForcePathStyle = true\n")
	return buf.String()
}

type sqsCustomConfig struct{}

func (config sqsCustomConfig) GoCode(si *serviceInfo) string {
	si.clientNames["sqs"] = "sqsClient"
	buf := bytes.NewBuffer(nil)
	buf.WriteString("sqsClient := sqs.New(cfg)\n")
	buf.WriteString("sqsClient.DisableComputeChecksums = true\n")
	return buf.String()
}

// CustomConfigurations is used to setup any custom configuration on a given
// service client
func (si *serviceInfo) CustomConfiguration(p string) string {
	if c, ok := customConfigs[p]; ok {
		return c.GoCode(si)
	}

	return ""
}

// Services will construct each service as a reflect.Value and its name.
func (si *serviceInfo) ServiceClient(p string) string {
	if name, ok := si.clientNames[p]; ok {
		return fmt.Sprintf("{name: %q, value: reflect.ValueOf(%s)},\n", p, name)
	}

	return fmt.Sprintf("{name: %q, value: reflect.ValueOf(%s.New(cfg))},\n", p, p)
}

func main() {
	si := &serviceInfo{
		clientNames: map[string]string{},
		basePath:    os.Args[1],
		Packages:    []string{},
	}

	si.getServices()

	var buf bytes.Buffer
	if err := tplCreateServiceCase.Execute(&buf, si); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("./create_service.go", buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}
