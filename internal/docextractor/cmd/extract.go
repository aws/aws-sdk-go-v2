package main

import (
	"fmt"
	"path/filepath"
	"go/token"
	"io/fs"
	"go/ast"
	"go/parser"
	"strings"
	"log"
)


// Extract will extract documentation from serviceDir and all sub-directories,
// populate items with newly-created JewelryItem(s).
// The overall strategy is to do a 
func Extract(servicePath string, serviceDir fs.DirEntry, items map[string]JewelryItem) {

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
				Types: map[string]*ast.TypeSpec{},
				Functions: map[string]*ast.FuncDecl{},
				Fields: map[string]*ast.Field{},
				Other: []*ast.GenDecl{},
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
func extractTypes(packageName string, types map[string]*ast.TypeSpec, items map[string]JewelryItem) error {
	for kt, vt := range types {
		summary := ""
		if vt.Doc != nil {
			summary = vt.Doc.Text()
		}
		typeName := vt.Name.Name

		item := JewelryItem{
			Name:    typeName,
			Summary: summary,
			Members: []JewelryItem{},
			Tags: []string{},
			OtherBlocks: map[string]string{},
			Params: []JewelryParam{},
			BreadCrumbs: []BreadCrumb{
				{
					Name: packageName,
					Kind: PACKAGE,
				},
			},
		}
		members := []JewelryItem{}

		st, ok := vt.Type.(*ast.StructType)

		if !ok {
			item.Type = INTERFACE

			bc := item.BreadCrumbs
			bc = append(bc, BreadCrumb{
				Name: typeName,
				Kind: INTERFACE,
			})
			item.BreadCrumbs = bc
			item.Signature = TypeSignature{
				Signature: fmt.Sprintf("type %v interface", typeName),
			}

		} else {
			item.Type = STRUCT
			bc := item.BreadCrumbs
			bc = append(bc, BreadCrumb{
				Name: typeName,
				Kind: STRUCT,
			})
			item.BreadCrumbs = bc
			item.Signature = TypeSignature{
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
					var fieldItem JewelryItem
					if vf.Doc == nil || vf.Doc.List == nil || vf.Doc.List[i] == nil {
						fieldItem = JewelryItem{
							Name:    fieldName,
							Tags: []string{},
							OtherBlocks: map[string]string{},
							Params: []JewelryParam{},
							Members: []JewelryItem{},
							Summary: "",
						}

					} else {
						fieldItem = JewelryItem{
							Name:    fieldName,
							Tags: []string{},
							OtherBlocks: map[string]string{},
							Params: []JewelryParam{},
							Members: []JewelryItem{},
							Summary: vf.Doc.List[i].Text,
						}
					}
					fieldItem.Type = FIELD
					fieldItem.BreadCrumbs = []BreadCrumb{
						{
							Name: packageName,
							Kind: PACKAGE,
						},
						{
							Name: typeName,
							Kind: STRUCT,
						},
						{
							Name: fieldName,
							Kind: FIELD,
						},
					}
					se, ok := vf.Type.(*ast.StarExpr)
					if ok {
						ident, ok := se.X.(*ast.Ident)
						if ok {
							fieldItem.Signature = TypeSignature{
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


func extractFunctions(packageName string, types map[string]*ast.TypeSpec, functions map[string]*ast.FuncDecl, items map[string]JewelryItem) error {
	for _, vf := range functions {
		
		// extract top-level functions
		if vf.Recv == nil {
			functionName := vf.Name.Name
			items[functionName] = JewelryItem{
				Type:    FUNCTION,
				Name:    functionName,
				Tags: []string{},
				OtherBlocks: map[string]string{},
				Params: []JewelryParam{},
				Members: []JewelryItem{},
				Summary: vf.Doc.Text(),
				BreadCrumbs: []BreadCrumb{
					{
						Name: packageName,
						Kind: PACKAGE,
					},
					{
						Name: functionName,
						Kind: FUNCTION,
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

		params := []JewelryParam{}
		returns := ""

		// extract operations
		// assumes that all receiver methods on Client are
		// service API operations except for the Options method.
		if receiverName == "Client" && methodName != "Options" {
			inputItem := items[fmt.Sprintf("%vInput", methodName)]
			input := JewelryParam{
				JewelryItem: JewelryItem{
					Name: inputItem.Name,
					Summary: inputItem.Summary,
					Type: inputItem.Type,
					Members: inputItem.Members,
					BreadCrumbs: inputItem.BreadCrumbs,
					Signature: inputItem.Signature,
				},
				IsOptional: false,
				IsReadonly: false,
			}
			params = append(params, input)
			returns = fmt.Sprintf("%vOutput", methodName)
		}

		members := i.Members
		members = append(members,
			JewelryItem{
				Type:    METHOD,
				Name:    methodName,
				Members: []JewelryItem{},
				Tags: []string{},
				OtherBlocks: map[string]string{},
				Params: params,
				Returns: returns,
				Summary: vf.Doc.Text(),
				BreadCrumbs: []BreadCrumb{
					{
						Name: packageName,
						Kind: PACKAGE,
					},
					{
						Name: receiverName,
						Kind: STRUCT,
					},
					{
						Name: methodName,
						Kind: METHOD,
					},
				},
			},
		)
		i.Members = members
		items[receiverName] = i
	}
	
	return nil
}