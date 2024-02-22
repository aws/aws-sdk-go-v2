package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"io/fs"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"path/filepath"
	"strings"
	// "slices"
	// "sort"
)

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

// for every directory that is 1 below aws-sdk-go-v2/service, call parseServiceDir
// then parseServiceDir should call parseServicePackage which takes in a packageName.
// parseServicePackage should ALWAYS use that packageName for any nested directories.
func parseServiceDir(servicePath string, serviceDir fs.DirEntry, items map[string]JewelryItem) {

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
				fmt.Printf("Skipping %v\n", path)
				return nil
			}


			fset := token.NewFileSet() // positions are relative to fset
			// why does d only contain the service package and not the types package?
			directory, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				panic(err)
			}
			// items := []interface{}{
			// 	PackageDocumentation{},
			// }
		
			// needs to be a map, so that i can check if a bad dupe item (no doc) exists
			// additionally, separation of types will make it easier to add methods to types:
			// this is because when youre iterating through the nodes, and you hit a method
			// you dont know if the type has been visited yet, so ensuring that a pass has
			// occurred through all the types is guaranteed with separating into 4 maps,
			// rather than one
			index := astIndex{
				Types: map[string]*ast.TypeSpec{},
				Functions: map[string]*ast.FuncDecl{},
				Fields: map[string]*ast.Field{},
				Other: []*ast.GenDecl{},
			}

		
			// the nodes need to be ordered alphabetically
			for _, p := range directory {
				removeTestFiles(p.Files)

		
				// add package JewelryItem
		
				packageItem, err := getPackageItem(packageName, p.Files)
				if err == nil {
					items["packageDocumentation"] = packageItem
				}
		
				indexFromAst(p, &index)

		
				// need to go through types to get the contained fields
				// in TypeSpec
				for kt, vt := range index.Types {
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
		
						// add interface type to breadcrumb
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
		
					// populate fields into type's JewelryItem
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
		
				for _, vf := range index.Functions {
		
					if vf.Recv == nil {
						// TODO: add top-level functions to items list
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
					_, ok := index.Types[receiverName]
					if !ok {
						// type doesnt exist
						continue
					}
		
					methodName := vf.Name.Name
		
					i := items[receiverName]

					params := []JewelryParam{}
					returns := ""
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
		
				// add breadcrumbs to
		
				// fmt.Println(items)
			}



			return nil
		})

		// there are single files getting through here
	clientData := map[string]JewelryItem{}
	for k, v := range items {
		if k == "packageDocumentation" {
			clientData[k] = v
			continue
		}
		clientData[v.Name] = v
		if v.Name == "Client" {
			for _, m := range v.Members {
				if m.Name == "Options" {
					continue
				}
				// m.Input = fmt.Sprintf("%vInput", m.Name)
				// m.Output = fmt.Sprintf("%vOutput", m.Name)
				clientData[m.Name] = m
			}
		}
	}
	content, err := json.Marshal(clientData)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(fmt.Sprintf("clients/%v.json", packageName), content, 0644)
	if err != nil {
		panic(err)
	}

	for _, item := range clientData {
		if item.Name == "" || item.Name == "packageDocumentation" {
			continue
		}
		content, err := json.Marshal(item)
		if err != nil {
			fmt.Println(err)
		}
		err = os.WriteFile(fmt.Sprintf("public/members/-aws-sdk-client-%v.%v.%v.json", packageName, item.Name, string(item.Type)), content, 0644)
		if err != nil {
			panic(err)
		}
	}


	typeData := map[string][]string{}
	for _, item := range clientData {
		if item.Name == "" || item.Name == "packageDocumentation" {
			continue
		}
		val, ok := typeData[string(item.Type)]
		if !ok {
			val = []string{}
		}
		val = append(val, item.Name)
		typeData[string(item.Type)] = val
	}
	content, err = json.Marshal(typeData)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(fmt.Sprintf("public/members/-aws-sdk-client-%v.json", packageName), content, 0644)
	if err != nil {
		panic(err)
	}
}


func main() {

	// initialize authoritative items
	// 
	maxDepthMacbook := 6
	// maxDepthDevDsk := 3

	serviceDirMacbook := "/Users/isvita/go-v2-isaiahvita/aws-sdk-go-v2/service"
	// serviceDirDevDsk := "../../service"

	filepath.WalkDir(serviceDirMacbook,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// handle possible path err, just in case...
				return err
			}
			if !d.IsDir() {
				return nil
			}

			currentDepth := strings.Count(path, string(os.PathSeparator))

			isInternal := strings.Count(path, "/internal") > 0
			if isInternal {
				fmt.Printf("Skipping %v\n", path)
				return nil
			}

			if currentDepth > maxDepthMacbook || isInternal {
				return fs.SkipDir
			}
			// ... process entry
			items := map[string]JewelryItem{}
			parseServiceDir(path, d, items)
			fmt.Printf("processed %v\n", path)
			return nil
    })

}
