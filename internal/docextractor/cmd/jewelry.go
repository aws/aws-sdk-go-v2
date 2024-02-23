package main


type JewelryItemKind string

const (
	PACKAGE   JewelryItemKind = "Package"
	STRUCT    JewelryItemKind = "Struct"
	INTERFACE JewelryItemKind = "Interface"
	FUNCTION  JewelryItemKind = "Function"
	METHOD    JewelryItemKind = "Method"
	FIELD     JewelryItemKind = "Field"
	OTHER     JewelryItemKind = "Other"
)

type BreadCrumb struct {
	Name string          `json:"name"`
	Kind JewelryItemKind `json:"kind"`
}

type TypeSignature struct {
	Signature string `json:"signature"`
	Location string `json:"location"`
}

type JewelryParam struct {
	JewelryItem
	IsOptional bool
	IsReadonly bool
	IsEventProperty bool
}

type JewelryItem struct {
	Name        string          `json:"name"`
	Summary     string          `json:"summary"`
	Type        JewelryItemKind `json:"type"`
	Members     []JewelryItem   `json:"members"`
	BreadCrumbs []BreadCrumb    `json:"breadcrumbs"`
	Signature TypeSignature `json:"typeSignature"`
	Tags []string `json:"tags"`
	Params []JewelryParam `json:"params"`
	Returns string `json:"returns"`
	// // optional (used only for JewelryOperations)
	// // since no out-of-box thing in Go for union types
	// Input string
	// // optional. see above.
	// Output string
	OtherBlocks map[string]string `json:"otherBlocks"`
}