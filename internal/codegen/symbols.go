package codegen

import "strings"

// SymbolizeExport takes an input value and symbolizes it for export.
func SymbolizeExport(v string) string {
	v = strings.Replace(strings.Replace(v, "-", " ", -1), "_", " ", -1)
	return strings.Replace(strings.Title(v), " ", "", -1)
}
