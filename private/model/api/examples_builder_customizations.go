// +build codegen

package api

import (
	"bytes"
)

type wafregionalExamplesBuilder struct {
	defaultExamplesBuilder
}

func (builder wafregionalExamplesBuilder) Imports(a *API) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(`"fmt"
	"context"
	"strings"
	"time"

	"` + SDKImportRoot + `/aws"
	"` + SDKImportRoot + `/aws/awserr"
	"` + SDKImportRoot + `/aws/external"
	"` + SDKImportRoot + `/service/waf"
	"` + a.ImportPath() + `"
	`)

	return buf.String()
}
