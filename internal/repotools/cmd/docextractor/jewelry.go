package main

type jewelryItemKind string

const (
	jewelryItemKindPackage   jewelryItemKind = "Package"
	jewelryItemKindStruct                    = "Struct"
	jewelryItemKindInterface                 = "Interface"
	jewelryItemKindFunc                      = "Function"
	jewelryItemKindMethod                    = "Method"
	jewelryItemKindField                     = "Field"
	jewelryItemKindOther                     = "Other"
)

type breadCrumb struct {
	Name string          `json:"name"`
	Kind jewelryItemKind `json:"kind"`
}

type typeSignature struct {
	Signature string `json:"signature"`
	Location  string `json:"location"`
}

type jewelryParam struct {
	jewelryItem
	IsOptional      bool
	IsReadonly      bool
	IsEventProperty bool
}

type jewelryItem struct {
	Name        string          `json:"name"`
	Summary     string          `json:"summary"`
	Type        jewelryItemKind `json:"type"`
	Members     []jewelryItem   `json:"members"`
	BreadCrumbs []breadCrumb    `json:"breadcrumbs"`
	Signature   typeSignature   `json:"typeSignature"`
	Tags        []string        `json:"tags"`
	Params      []jewelryParam  `json:"params"`
	Returns     string          `json:"returns"`
	// // optional (used only for JewelryOperations)
	// // since no out-of-box thing in Go for union types
	// Input string
	// // optional. see above.
	// Output string
	OtherBlocks map[string]string `json:"otherBlocks"`
}
