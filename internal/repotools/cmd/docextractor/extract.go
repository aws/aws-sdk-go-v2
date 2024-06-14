package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

// Extract will extract documentation from serviceDir and all sub-directories,
// populate items with newly-created JewelryItem(s).
// The overall strategy is to do a
func Extract(servicePath string, serviceDir fs.DirEntry, items map[string]jewelryItem) {

	if serviceDir.Name() == "service" {
		return
	}

	packageName := serviceDir.Name()

	filepath.WalkDir(servicePath,
		func(path string, d fs.DirEntry, e error) error {
			if !d.IsDir() {
				return nil
			}

			isInternal := strings.Count(path, "/internal") > 0
			if isInternal {
				return nil
			}

			fset := token.NewFileSet()
			directory, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				panic(err)
			}

			index := astIndex{
				Types:            map[string]*ast.TypeSpec{},
				Functions:        map[string]*ast.FuncDecl{},
				Fields:           map[string]*ast.Field{},
				StringEnumConsts: map[string]string{},
				Other:            []*ast.GenDecl{},
			}

			for _, p := range directory {
				removeTestFiles(p.Files)

				packageItem, err := getPackageItem(packageName, p.Files)
				if err == nil {
					items["packageDocumentation"] = packageItem
				}

				indexFromAst(p, &index)

				err = extractTypes(d.Name(), packageName, index, items)
				if err != nil {
					log.Fatal(err)
				}

				err = extractFunctions(packageName, index.Types, index.Functions, items)
				if err != nil {
					log.Fatal(err)
				}
			}

			return nil
		})

	serialize(packageName, items)
}

// extractType iterates through
func extractTypes(pkg, module string, index astIndex, items map[string]jewelryItem) error {
	types := index.Types
	for kt, vt := range types {
		typeName := vt.Name.Name

		item := jewelryItem{
			Package:     pkg,
			Name:        typeName,
			Summary:     formatComment(vt.Doc),
			Members:     []jewelryItem{},
			Tags:        []string{},
			OtherBlocks: map[string]string{},
			Params:      []jewelryParam{},
			BreadCrumbs: []breadCrumb{
				{
					Name: module,
					Kind: jewelryItemKindPackage,
				},
			},
		}
		members := []jewelryItem{}

		switch st := vt.Type.(type) {
		case *ast.StructType:
			item.Type = jewelryItemKindStruct
			bc := item.BreadCrumbs
			bc = append(bc, breadCrumb{
				Name: typeName,
				Kind: jewelryItemKindStruct,
			})
			item.BreadCrumbs = bc
			item.Signature = typeSignature{
				Signature: fmt.Sprintf("type %v struct", typeName),
			}
			for _, vf := range st.Fields.List {
				namesNum := len(vf.Names)
				for i := 0; i < namesNum; i++ {
					if !isExported(vf.Names[i].Name) {
						break
					}
					fieldName := vf.Names[i].Name
					fieldItem := jewelryItem{
						Name:        fieldName,
						Tags:        []string{},
						OtherBlocks: map[string]string{},
						Params:      []jewelryParam{},
						Members:     []jewelryItem{},
						Summary:     formatComment(vf.Doc),
					}
					fieldItem.Type = jewelryItemKindField
					fieldItem.BreadCrumbs = []breadCrumb{
						{
							Name: module,
							Kind: jewelryItemKindPackage,
						},
						{
							Name: typeName,
							Kind: jewelryItemKindStruct,
						},
						{
							Name: fieldName,
							Kind: jewelryItemKindField,
						},
					}
					fieldItem.Signature = typeSignature{
						Signature: toSignature(vf.Type, module),
						// Location is unused - links have to be embedded in signature
					}

					members = append(members, fieldItem)
				}
			}
		case *ast.Ident:
			if st.Name != "string" {
				continue
			}

			// probably an enum, map its variants to members
			item.Type = jewelryItemKindEnum
			item.BreadCrumbs = append(item.BreadCrumbs, breadCrumb{
				Name: typeName,
				Kind: jewelryItemKindEnum,
			})
			for name, value := range index.StringEnumConsts {
				if strings.HasPrefix(name, typeName) { // good enough
					members = append(members, jewelryItem{
						Name:      name,
						Signature: typeSignature{Signature: typeName},
						Summary:   value,
					})
				}
			}
		case *ast.InterfaceType:
			if !isProbablyUnion(typeName, st) {
				continue
			}

			item.Type = jewelryItemKindUnion
			item.BreadCrumbs = append(item.BreadCrumbs, breadCrumb{
				Name: typeName,
				Kind: jewelryItemKindUnion,
			})
			for name, typ := range index.Types {
				if strings.HasPrefix(name, typeName+"Member") {
					members = append(members, jewelryItem{
						Name: name,
						Signature: typeSignature{
							Signature: fmt.Sprintf("[%s](-aws-sdk-client-%s!%s:Struct)", name, module, name),
						},
						Summary: formatComment(typ.Doc),
					})
				}
			}
		default:
			continue
		}

		item.Members = members
		items[kt] = item
	}
	return nil
}

func isProbablyUnion(name string, i *ast.InterfaceType) bool {
	for _, field := range i.Methods.List {
		if len(field.Names) == 0 {
			continue
		}
		if field.Names[0].Name == "is"+name {
			return true
		}
	}
	return false
}

// We've already converted the model's HTML to Go docs, now for ref docs we
// must convert back. We can't use the model's original docs directly because
// that doesn't include extra content we may inject at codegen.
//
// Practically this is just a matter of converting lists and paragraphs, since
// that's really all Go docs do in terms of formatting. Links are left as-is.
// Note we don't bother with <ul> on lists since our API ref docs generator
// doesn't require them to render.
func formatComment(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}

	var inlist bool
	var html, currp, currli string
	for _, c := range cg.List {
		line := c.Text
		if line == "//" {
			flushp(&html, &currp)
			continue
		}

		line = strings.TrimPrefix(line, "// ")
		if strings.HasPrefix(line, "  ") && !inlist {
			inlist = true
			flushp(&html, &currp)
		}
		if !strings.HasPrefix(line, "  ") && inlist {
			inlist = false
			flushli(&html, &currli)
		}

		if inlist {
			if strings.HasPrefix(line, "  -") {
				if len(currli) > 0 {
					flushli(&html, &currli)
				}
				currli = strings.TrimPrefix(line, "  -") + " "
			} else if strings.HasPrefix(line, "  ") {
				currli += strings.TrimPrefix(line, "  ") + " "
			}
		} else {
			currp += line + " "
		}
	}

	if len(currp) > 0 {
		flushp(&html, &currp)
	}
	if len(currli) > 0 {
		flushli(&html, &currli)
	}

	return html
}

func flushp(dst, src *string) {
	if len(*src) == 0 {
		return
	}
	*dst = *dst + "<p>" + *src + "</p>"
	*src = ""
}

func flushli(dst, src *string) {
	if len(*src) == 0 {
		return
	}
	*dst = *dst + "<li>" + *src + "</li>"
	*src = ""
}

func toSignature(v ast.Expr, pkg string) string {
	switch vv := v.(type) {
	case *ast.Ident:
		return fmt.Sprintf("%s", vv.Name)
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", toSignature(vv.X, pkg))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", toSignature(vv.Elt, pkg))
	case *ast.MapType:
		return fmt.Sprintf("map[string]%s", toSignature(vv.Value, pkg))
	case *ast.SelectorExpr:
		spkg := vv.X.(*ast.Ident).Name
		if spkg == "types" {
			return fmt.Sprintf("[%s.%s](-aws-sdk-client-%s!%s:Struct)", spkg, vv.Sel.Name, pkg, vv.Sel.Name)
		}
		// FUTURE: handle links to runtime
		return fmt.Sprintf("%s.%s", spkg, vv.Sel.Name)
	case *ast.FuncType:
		return toFuncSignature(vv, pkg)
	default:
		return fmt.Sprintf("[unhandled %T]", v)
	}
}

func toFuncSignature(v *ast.FuncType, pkg string) string {
	xpr := "func ("
	if v.Params != nil {
		for i, param := range v.Params.List {
			xpr += toSignature(param.Type, pkg)
			if i < len(v.Params.List)-1 {
				xpr += ", "
			}
		}
	}
	xpr += ")"
	if v.Results != nil {
		xpr += " "
		for i, param := range v.Results.List {
			xpr += toSignature(param.Type, pkg)
			if i < len(v.Params.List)-1 {
				xpr += ", "
			}
		}
	}
	return xpr
}

func extractFunctions(packageName string, types map[string]*ast.TypeSpec, functions map[string]*ast.FuncDecl, items map[string]jewelryItem) error {
	for _, vf := range functions {

		// extract top-level functions
		if vf.Recv == nil {
			functionName := vf.Name.Name
			items[functionName] = jewelryItem{
				Type:        jewelryItemKindFunc,
				Name:        functionName,
				Tags:        []string{},
				OtherBlocks: map[string]string{},
				Params:      []jewelryParam{},
				Members:     []jewelryItem{},
				Summary:     formatComment(vf.Doc),
				BreadCrumbs: []breadCrumb{
					{
						Name: packageName,
						Kind: jewelryItemKindPackage,
					},
					{
						Name: functionName,
						Kind: jewelryItemKindFunc,
					},
				},
			}
			continue
		}
		var receiverName string
		switch r := vf.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			rName, _ := r.X.(*ast.Ident)
			receiverName = rName.Name
		case *ast.Ident:
			receiverName = r.Name
		}

		// grab existing type
		_, ok := types[receiverName]
		if !ok {
			// type doesnt exist
			continue
		}

		methodName := vf.Name.Name

		i := items[receiverName]

		params := []jewelryParam{}
		returns := ""

		// extract operations
		// assumes that all receiver methods on Client are
		// service API operations except for the Options method.
		if receiverName == "Client" && methodName != "Options" {
			inputItem := items[fmt.Sprintf("%vInput", methodName)]
			input := jewelryParam{
				jewelryItem: jewelryItem{
					Name:        inputItem.Name,
					Summary:     inputItem.Summary,
					Type:        inputItem.Type,
					Members:     inputItem.Members,
					BreadCrumbs: inputItem.BreadCrumbs,
					Signature:   inputItem.Signature,
				},
				IsOptional: false,
				IsReadonly: false,
			}
			params = append(params, input)
			returns = fmt.Sprintf("%vOutput", methodName)
		}

		members := i.Members

		// without proper runtime documentation, we have to bridge the gap to
		// event payloads for now
		if vf.Name.Name == "GetStream" && isInputOutput(receiverName) {
			stream := strings.TrimSuffix(receiverName, whichSuffix(receiverName)) + "EventStream"
			members = append(members, jewelryItem{
				Name:    "(event stream payload)",
				Summary: "The event streaming payload union for this structure.",
				Signature: typeSignature{
					Signature: fmt.Sprintf("[%s](-aws-sdk-client-%s!%s:Union)", stream, packageName, stream),
				},
			})
		} else {
			members = append(members,
				jewelryItem{
					Type:        jewelryItemKindMethod,
					Name:        methodName,
					Members:     []jewelryItem{},
					Tags:        []string{},
					OtherBlocks: map[string]string{},
					Params:      params,
					Returns:     returns,
					Summary:     formatComment(vf.Doc),
					BreadCrumbs: []breadCrumb{
						{
							Name: packageName,
							Kind: jewelryItemKindPackage,
						},
						{
							Name: receiverName,
							Kind: jewelryItemKindStruct,
						},
						{
							Name: methodName,
							Kind: jewelryItemKindMethod,
						},
					},
				},
			)
		}

		i.Members = members
		items[receiverName] = i
	}

	return nil
}

// whether there's "Input" or "Output" on a structure
func isInputOutput(name string) bool {
	return strings.HasSuffix(name, "Input") || strings.HasSuffix(name, "Output")
}

// "Input" or "Output" on a structure
func whichSuffix(name string) string {
	if strings.HasSuffix(name, "Input") {
		return "Input"
	} else if strings.HasSuffix(name, "Output") {
		return "Output"
	}
	panic("expected -Input or -Output suffix")
}
