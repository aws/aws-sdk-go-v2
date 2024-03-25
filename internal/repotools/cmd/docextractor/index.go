package main

import (
	"go/ast"
	"go/token"
)

type astIndex struct {
	Types map[string]*ast.TypeSpec
	Functions map[string]*ast.FuncDecl
	Fields map[string]*ast.Field
	Other []*ast.GenDecl
}


func indexFromAst(p *ast.Package, index *astIndex) {
	ast.Inspect(p, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:

			// remove unexported items
			if !isExported(x.Name.Name) {
				break
			}
			name := x.Name.Name
			index.Functions[name] = x

		// use TypeSpec (over like StructType) because
		// StructType doesnt have the name of the thing for some reason
		// and TypeSpec contains the StructType obj as a field.
		case *ast.TypeSpec:

			if !isExported(x.Name.Name) {
				break
			}

			// if a type exists AND it has a doc comment, then
			// dont add anything -- were good.
			// if not, then just add whatever.
			name := x.Name.Name
			if _, ok := index.Types[name]; ok && index.Types[name].Doc.Text() != "" {
				break
			}

			index.Types[name] = x
		case *ast.Field:
			namesNum := len(x.Names)
			for i := 0; i < namesNum; i++ {
				if !isExported(x.Names[i].Name) {
					break
				}
				name := x.Names[i].Name
				index.Fields[name] = x
			}
		case *ast.GenDecl:

			// for some reason, the same type will show up in the AST node list
			// one with documentation and one without documentation
			if x.Tok == token.TYPE {
				xt, _ := x.Specs[0].(*ast.TypeSpec)

				name := xt.Name.Name
				if !isExported(name) {
					break
				}

				// if a type exists AND it has a doc comment, then
				// dont add anything -- were good.
				// if not, then just add whatever.
				if _, ok := index.Types[name]; ok && index.Types[name].Doc.Text() != "" {
					break
				}

				// its a comment group, and each item in the list
				// is a line
				// summary := ""
				if x.Doc != nil && x.Doc.List != nil {
					xt.Doc = x.Doc
					// for _, line := range x.Doc.List {
					// 	summary += line.Text
					// }
				}

				index.Types[name] = xt
			} else {
				index.Other = append(index.Other, x)
			}
		}
		return true
	})
}