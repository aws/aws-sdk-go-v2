package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/gomod"
	"golang.org/x/mod/modfile"
)

func main() {
	repoRoot, err := repotools.GetRepoRoot()
	if err != nil {
		log.Fatalf("failed to get repository root: %v", err)
	}

	registry := NewRegistry()

	discoverer := gomod.NewDiscoverer(repoRoot)
	if err = discoverer.Discover(); err != nil {
		log.Fatalf("failed to discover modules: %v", err)
	}

	// Load Discovered Modules into Registry
	var modules []string

	for moduleDir := range discoverer.Modules() {
		m := registry.MustLoad(moduleDir)
		modules = append(modules, m.Module.Mod.Path)
	}

	var modulePath string
	for len(modules) > 0 {
		modulePath, modules = modules[0], modules[1:]

		err = addRelativeReplaces(repoRoot, modulePath, registry)
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
	dirToModule     map[string]*Module
	modulePathToDir map[string]string
}

// NewRegistry returns a new module registry.
func NewRegistry() *Registry {
	return &Registry{
		dirToModule:     map[string]*Module{},
		modulePathToDir: map[string]string{},
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
func (r *Registry) MustGet(path string) (string, *Module) {
	modulePath, module, err := r.Get(path)
	if err != nil {
		panic(err)
	}
	return modulePath, module
}

// Get retrieves the module identified by the give module path.
func (r *Registry) Get(path string) (string, *Module, error) {
	dir, ok := r.modulePathToDir[path]
	if !ok {
		return "", nil, fmt.Errorf("module not found")
	}

	module, ok := r.dirToModule[dir]
	if !ok {
		return "", nil, fmt.Errorf("module missing or not loaded")
	}

	return dir, module, nil
}

// Has returns whether the given module path is in the registry.
func (r *Registry) Has(path string) bool {
	_, ok := r.modulePathToDir[path]
	return ok
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
		m, err := gomod.LoadModuleFile(dir, nil, false)
		if err != nil {
			return nil, err
		}
		module = &Module{File: m}
		r.dirToModule[dir] = module
		r.modulePathToDir[module.Module.Mod.Path] = dir
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

	m.Cleanup()

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

type toReplace struct {
	ModulePath   string
	RelativePath string
}

// addRelativeReplaces takes the given root and submodule paths and adds go.mod replace directives for any sub modules
// that refer to the given root as a dependency.
func addRelativeReplaces(repoRoot, modulePath string, registry *Registry) error {
	modDir, mod := registry.MustGet(modulePath)

	modDirToRoot, err := filepath.Rel(modDir, repoRoot)
	if err != nil {
		return err
	}

	var toDrop []string
	for _, replace := range mod.Replace {
		if !registry.Has(replace.Old.Path) {
			continue
		}
		toDrop = append(toDrop, replace.Old.Path)
	}

	for _, drop := range toDrop {
		if err := mod.DropReplace(drop, ""); err != nil {
			return err
		}
	}

	seen := make(map[string]struct{})
	var replaces []toReplace
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

		if !registry.Has(req.Mod.Path) {
			continue
		}

		reqDir, reqMod := registry.MustGet(req.Mod.Path)

		reqFromRoot, err := filepath.Rel(repoRoot, reqDir)
		if err != nil {
			return err
		}

		relPathToReq := filepath.Join(modDirToRoot, reqFromRoot)
		if !strings.HasSuffix(relPathToReq, string(filepath.Separator)) {
			relPathToReq += string(filepath.Separator)
		}

		replaces = append(replaces, toReplace{
			ModulePath:   req.Mod.Path,
			RelativePath: relPathToReq,
		})

		toProcess = append(toProcess, reqMod.Require...)
	}

	sort.Slice(replaces, func(i, j int) bool {
		return replaces[i].ModulePath < replaces[j].ModulePath
	})

	for _, replace := range replaces {
		err = mod.AddReplace(replace.ModulePath, "", replace.RelativePath, "")
		if err != nil {
			return err
		}
	}

	return nil
}
