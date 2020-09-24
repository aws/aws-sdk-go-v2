package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	gitRoot, err := findGitRoot(currDir)
	if err != nil {
		log.Fatal(err)
	}

	_, ok := isGoModPresent(gitRoot)
	if !ok {
		log.Fatalf("go.mod not present at %v", gitRoot)
	}

	root := gitRoot
	toProcess := []string{root}

	registry := make(Registry)

	for len(toProcess) > 0 {
		root, toProcess = toProcess[0], toProcess[1:]

		subs, err := findSubModules(root)
		if err != nil {
			log.Fatalf("failed to find submodules: %v", err)
		}

		err = addRelativeReplaces(root, subs, registry)
		if err != nil {
			log.Fatal(err)
		}

		for i := range subs {
			toProcess = append(toProcess, subs[i])
		}
	}

	for _, module := range registry {
		if err := module.Write(); err != nil {
			log.Fatal(err)
		}
	}
}

// Registry is a map of module path to a module
type Registry map[string]*Module

// Get loads or retrieves the Module from the registry for the given path.
func (r Registry) Get(dir string) (module *Module, err error) {
	module, ok := r[dir]
	if !ok {
		m, err := loadGoMod(dir)
		if err != nil {
			return nil, err
		}
		module = &Module{File: m}
		r[dir] = module
	}

	return module, nil
}

// Module is a go.mod file that tracks whether modifications have been made.
type Module struct {
	*modfile.File
	modified bool
}

// Write writes any pending changes back to the go.mod
func (m *Module) Write() error {
	if !m.modified {
		return nil
	}

	mb, err := m.Format()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(m.Syntax.Name, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(mb))
	if err != nil {
		return err
	}

	return nil
}

// AddReplace replaces oldPath with newPath.
func (m *Module) AddReplace(oldPath, oldVers, newPath, newVers string) error {
	m.modified = true
	return m.File.AddReplace(oldPath, oldVers, newPath, newVers)
}

// addRelativeReplaces takes the given root and submodule paths and adds go.mod replace directives for any sub modules
// that refer to the given root as a dependency.
func addRelativeReplaces(root string, subs []string, registry Registry) error {
	rootModule, err := registry.Get(root)
	if err != nil {
		return err
	}

	rootPath := rootModule.Module.Mod.Path

	for _, sub := range subs {
		mod, err := registry.Get(sub)
		if err != nil {
			return err
		}

		modRelativeToRoot := convertToDotted(makeRelativeTo(mod.Module.Mod.Path, rootPath))

		for _, req := range mod.Require {
			if !strings.HasPrefix(req.Mod.Path, rootPath) {
				continue
			}

			reqRelativeToNearestParent := makeRelativeTo(req.Mod.Path, rootPath)

			var relToReq strings.Builder
			if modRelativeToRoot != "." {
				relToReq.WriteString(modRelativeToRoot)
				if reqRelativeToNearestParent != "." {
					relToReq.WriteRune('/')
					relToReq.WriteString(reqRelativeToNearestParent)
				}
				relToReq.WriteRune('/')
			} else {
				relToReq.WriteString(reqRelativeToNearestParent)
			}

			err := mod.AddReplace(req.Mod.Path, "", relToReq.String(), "")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// makeRelativeTo makes the module path relative to rootModule.
func makeRelativeTo(module string, rootModule string) string {
	relative := strings.TrimLeft(strings.TrimPrefix(module, rootModule), "/")
	if relative == "" {
		return "."
	}
	return relative
}

// convertToDotted converts a relative path to a form such as ../ or ../../ etc.
func convertToDotted(path string) string {
	if path == "." {
		return path
	}

	count := strings.Count(path, "/")

	var builder strings.Builder
	first := true
	for i := 0; i <= count; i++ {
		if !first {
			builder.WriteRune('/')
		} else {
			first = false
		}
		builder.WriteString("..")
	}

	return builder.String()
}

// loadGoMod loads the go.mod file found at the given directory path
func loadGoMod(dir string) (*modfile.File, error) {
	path := filepath.Join(dir, "go.mod")
	mb, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	parse, err := modfile.Parse(path, mb, nil)
	if err != nil {
		return nil, err
	}

	return parse, nil
}

// findGitRoot finds the .git in dir or one of dir's parent directories.
func findGitRoot(dir string) (path string, err error) {
	found := false
	for {
		err = filepath.Walk(dir, func(fPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if dir == fPath {
				return nil
			}

			if info.IsDir() && info.Name() == ".git" {
				found = true
				path = dir
				return filepath.SkipDir
			} else if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		})
		if err != nil {
			return path, err
		}
		if !found && filepath.Base(dir) != "/" {
			dir = filepath.Dir(dir)
		} else {
			break
		}
	}

	if !found {
		return "", fmt.Errorf(".git directory not found")
	}

	if path == "" {
		path = "."
	}

	return path, nil
}

func findSubModules(dir string) (modules []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if dir == path {
			return nil
		}

		if info.IsDir() {
			_, ok := isGoModPresent(path)
			if !ok {
				return nil
			}
			modules = append(modules, path)
		}

		return nil
	})
	return modules, err
}

func isGoModPresent(dir string) (string, bool) {
	path := filepath.Join(dir, "go.mod")
	_, err := os.Stat(path)
	if err != nil {
		return "", false
	}
	return path, true
}
