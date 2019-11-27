// +build codegen

package api

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
	return ""
}
