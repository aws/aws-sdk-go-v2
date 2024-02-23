package main

import (
	"fmt"
	"go/ast"
	"regexp"
)

func getPackageItem(packageName string, files map[string]*ast.File) (JewelryItem, error) {
	packageDoc := ""
	docRe := regexp.MustCompile(`.*/doc.go`)

	for k, f := range files {
		matched := docRe.Match([]byte(k))
		if !matched {
			continue
		}
		if f.Doc != nil && f.Doc.List != nil {
			for _, line := range f.Doc.List {
				packageDoc += line.Text
			}
		}
		return JewelryItem{
			Tags: []string{},
			OtherBlocks: map[string]string{},
			Members: []JewelryItem{},
			Params: []JewelryParam{},
			BreadCrumbs: []BreadCrumb{},
			Summary: packageDoc,
		}, nil
	}
	return JewelryItem{}, fmt.Errorf("no doc.go")
}

func hasDocFile(p *ast.Package) bool {
	docRe := regexp.MustCompile(`.*/doc.go`)

	for k, _ := range p.Files {
		matched := docRe.Match([]byte(k))
		if !matched {
			continue
		}
		return true
	}
	return false	
}

func isExported(name string) bool {
	if name == "" {
		return false
	}
	firstChar := name[0]
	if firstChar >= 'a' && firstChar <= 'z' {
		return false
	}
	return true
}

func removeTestFiles(files map[string]*ast.File) error {
	testRe := regexp.MustCompile(`.*_test.go`)
	for key := range files {
		matched := testRe.Match([]byte(key))
		if !matched {
			continue
		}
		delete(files, key)
	}
	return nil
}