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
				Types:     map[string]*ast.TypeSpec{},
				Functions: map[string]*ast.FuncDecl{},
				Fields:    map[string]*ast.Field{},
				Other:     []*ast.GenDecl{},
			}

			for _, p := range directory {
				removeTestFiles(p.Files)

				packageItem, err := getPackageItem(packageName, p.Files)
				if err == nil {
					items["packageDocumentation"] = packageItem
				}

				indexFromAst(p, &index)

				err = extractTypes(packageName, index.Types, items)
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
func extractTypes(packageName string, types map[string]*ast.TypeSpec, items map[string]jewelryItem) error {
	for kt, vt := range types {
		summary := ""
		if vt.Doc != nil {
			summary = vt.Doc.Text()
		}
		typeName := vt.Name.Name

		item := jewelryItem{
			Name:        typeName,
			Summary:     summary,
			Members:     []jewelryItem{},
			Tags:        []string{},
			OtherBlocks: map[string]string{},
			Params:      []jewelryParam{},
			BreadCrumbs: []breadCrumb{
				{
					Name: packageName,
					Kind: jewelryItemKindPackage,
				},
			},
		}
		members := []jewelryItem{}

		st, ok := vt.Type.(*ast.StructType)

		if !ok {
			item.Type = jewelryItemKindInterface

			bc := item.BreadCrumbs
			bc = append(bc, breadCrumb{
				Name: typeName,
				Kind: jewelryItemKindInterface,
			})
			item.BreadCrumbs = bc
			item.Signature = typeSignature{
				Signature: fmt.Sprintf("type %v interface", typeName),
			}

		} else {
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
		}

		if ok && st.Fields != nil && st.Fields.List != nil {
			for _, vf := range st.Fields.List {
				namesNum := len(vf.Names)
				for i := 0; i < namesNum; i++ {
					if !isExported(vf.Names[i].Name) {
						break
					}
					fieldName := vf.Names[i].Name
					var fieldItem jewelryItem
					if vf.Doc == nil || vf.Doc.List == nil || vf.Doc.List[i] == nil {
						fieldItem = jewelryItem{
							Name:        fieldName,
							Tags:        []string{},
							OtherBlocks: map[string]string{},
							Params:      []jewelryParam{},
							Members:     []jewelryItem{},
							Summary:     "",
						}

					} else {
						fieldItem = jewelryItem{
							Name:        fieldName,
							Tags:        []string{},
							OtherBlocks: map[string]string{},
							Params:      []jewelryParam{},
							Members:     []jewelryItem{},
							Summary:     vf.Doc.List[i].Text,
						}
					}
					fieldItem.Type = jewelryItemKindField
					fieldItem.BreadCrumbs = []breadCrumb{
						{
							Name: packageName,
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
					se, ok := vf.Type.(*ast.StarExpr)
					if ok {
						ident, ok := se.X.(*ast.Ident)
						if ok {
							fieldItem.Signature = typeSignature{
								Signature: ident.Name,
							}
						}
					}
					members = append(members, fieldItem)
				}
			}
		}
		item.Members = members
		items[kt] = item
	}
	return nil
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
				Summary:     vf.Doc.Text(),
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
		members = append(members,
			jewelryItem{
				Type:        jewelryItemKindMethod,
				Name:        methodName,
				Members:     []jewelryItem{},
				Tags:        []string{},
				OtherBlocks: map[string]string{},
				Params:      params,
				Returns:     returns,
				Summary:     vf.Doc.Text(),
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
		i.Members = members
		items[receiverName] = i
	}

	return nil
}
