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

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"golang.org/x/mod/modfile"
)

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	gitRoot, err := repotools.FindRepoRoot(currDir)
	if err != nil {
		log.Fatal(err)
	}

	_, ok := isGoModPresent(gitRoot)
	if !ok {
		log.Fatalf("go.mod not present at %v", gitRoot)
	}

	registry := NewRegistry()

	rootModulePath := registry.MustLoad(gitRoot).Module.Mod.Path

	subPaths, err := findSubModules(gitRoot)
	if err != nil {
		log.Fatalf("failed to find submodules: %v", err)
	}

	// Load Discovered Modules into Registry
	var modules []string
	for _, sub := range subPaths {
		m := registry.MustLoad(sub)
		modules = append(modules, m.Module.Mod.Path)
	}

	var module string
	for len(modules) > 0 {
		module, modules = modules[0], modules[1:]

		err = addRelativeReplaces(rootModulePath, module, registry)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, module := range registry.Modules() {
		if err := module.Write(); err != nil {
			log.Fatal(err)
		}
	}
}

// Registry is a map of module path to a module
type Registry struct {
	dirToModule map[string]*Module
	pathToDir   map[string]string
}

// NewRegistry returns a new module registry.
func NewRegistry() *Registry {
	return &Registry{
		dirToModule: map[string]*Module{},
		pathToDir:   map[string]string{},
	}
}

// Modules returns the modules that were registered.
func (r *Registry) Modules() (m []*Module) {
	for _, module := range r.dirToModule {
		m = append(m, module)
	}
	return m
}

// MustGet retrieves the module identified by the given module path. Panics on failure.
func (r *Registry) MustGet(path string) *Module {
	module, err := r.Get(path)
	if err != nil {
		panic(err)
	}
	return module
}

// Get retrieves the module identified by the give module path.
func (r *Registry) Get(path string) (*Module, error) {
	dir, ok := r.pathToDir[path]
	if !ok {
		return nil, fmt.Errorf("module not found")
	}

	module, ok := r.dirToModule[dir]
	if !ok {
		return nil, fmt.Errorf("module missing or not loaded")
	}

	return module, nil
}

// MustLoad loads or retrieves the Module from the registry for the given path. Panics on failure.
func (r *Registry) MustLoad(dir string) *Module {
	module, err := r.Load(dir)
	if err != nil {
		panic(err)
	}
	return module
}

// Load loads or retrieves the Module from the registry for the given directory path.
func (r *Registry) Load(dir string) (module *Module, err error) {
	module, ok := r.dirToModule[dir]
	if !ok {
		m, err := loadGoMod(dir)
		if err != nil {
			return nil, err
		}
		module = &Module{File: m}
		r.dirToModule[dir] = module
		r.pathToDir[module.Module.Mod.Path] = dir
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
func addRelativeReplaces(repoModule, module string, registry *Registry) error {
	mod := registry.MustGet(module)

	modRelativeToRoot := convertToDotted(makeRelativeTo(mod.Module.Mod.Path, repoModule))

	seen := make(map[string]struct{})
	var toProcess []*modfile.Require
	var req *modfile.Require
	toProcess = append(toProcess, mod.Require...)

	for len(toProcess) > 0 {
		req, toProcess = toProcess[0], toProcess[1:]

		if _, ok := seen[req.Mod.Path]; ok {
			continue
		} else {
			seen[req.Mod.Path] = struct{}{}
		}

		if !strings.HasPrefix(req.Mod.Path, repoModule) {
			continue
		}

		reqMod := registry.MustGet(req.Mod.Path)

		reqFromRoot := makeRelativeTo(req.Mod.Path, repoModule)
		if reqFromRoot == "." {
			reqFromRoot = ""
		} else {
			reqFromRoot += "/"
		}

		err := mod.AddReplace(req.Mod.Path, "", fmt.Sprintf("%s/%s", modRelativeToRoot, reqFromRoot), "")
		if err != nil {
			return err
		}

		toProcess = append(toProcess, reqMod.Require...)
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
