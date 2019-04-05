// +build codegen

package api

import (
	"bytes"
	"fmt"
)

type examplesBuilder interface {
	BuildShape(*ShapeRef, map[string]interface{}, bool, bool) string
	BuildList(string, string, *ShapeRef, []interface{}) string
	BuildComplex(string, string, *ShapeRef, map[string]interface{}, bool) string
	GoType(*ShapeRef, bool) string
	Imports(*API) string
}

type defaultExamplesBuilder struct {
	ShapeValueBuilder
}

func (builder defaultExamplesBuilder) Imports(a *API) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(`"fmt"
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	`)

	buf.WriteString(fmt.Sprintf("\"%s/%s\"", "github.com/aws/aws-sdk-go-v2/service", a.PackageName()))
	return buf.String()
}
