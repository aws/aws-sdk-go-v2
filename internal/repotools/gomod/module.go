package gomod

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
)

const (
	goModuleFile   = "go.mod"
	testDataFolder = "testdata"
)

// GetModulePath retrieves the module path from the provide file description.
func GetModulePath(file *modfile.File) (string, error) {
	if file.Module == nil {
		return "", fmt.Errorf("module directive not present")
	}
	return file.Module.Mod.Path, nil
}

// LoadModuleFile loads the Go module file located at the provided directory path.
func LoadModuleFile(path string, fix modfile.VersionFixer, lax bool) (*modfile.File, error) {
	path = filepath.Join(path, goModuleFile)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ReadModule(path, f, fix, lax)
}

// ReadModule parses the module file bytes from the provided reader.
func ReadModule(path string, f io.Reader, fix modfile.VersionFixer, lax bool) (parse *modfile.File, err error) {
	fBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if lax {
		parse, err = modfile.ParseLax(path, fBytes, fix)
	} else {
		parse, err = modfile.Parse(path, fBytes, fix)
	}
	if err != nil {
		return nil, err
	}

	return parse, nil
}

// WriteModuleFile writes the Go module description to the provided directory path.
func WriteModuleFile(path string, file *modfile.File) (err error) {
	modPath := filepath.Join(path, goModuleFile)

	var mf *os.File
	mf, err = os.OpenFile(modPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		fErr := mf.Close()
		if fErr != nil && err == nil {
			err = fErr
		}
	}()

	var fb []byte
	fb, err = file.Format()
	if err != nil {
		return err
	}

	_, err = io.Copy(mf, bytes.NewReader(fb))

	return err
}

// Discoverer is used for discovering all modules and submodules at the provided path.
type Discoverer struct {
	path    string
	modules map[string][]string
}

// NewDiscoverer constructs a new Discover for the given path.
func NewDiscoverer(path string) *Discoverer {
	return &Discoverer{
		path: path,
	}
}

// Root returns the root path of the module discovery.
func (d *Discoverer) Root() string {
	return d.path
}

// Modules returns the modules discovered after executing Discover.
func (d *Discoverer) Modules() (v map[string][]string) {
	v = make(map[string][]string)
	for modulePath, children := range d.modules {
		var c []string
		if children != nil {
			c := make([]string, 0, len(children))
			copy(c, children)
		}
		v[modulePath] = c
	}
	return v
}

// ModulesRel returns the modules discovered after executing Discover. The returned module directory paths
// will be made relative to the provided base path.
func (d *Discoverer) ModulesRel() (v map[string][]string, err error) {
	v = make(map[string][]string)
	for modulePath, children := range d.modules {
		rel, err := filepath.Rel(d.path, modulePath)
		if err != nil {
			return nil, err
		}
		var c []string
		if len(children) > 0 {
			c = make([]string, 0, len(children))
			for i := range children {
				rel, err := filepath.Rel(d.path, children[i])
				if err != nil {
					return nil, err
				}
				c = append(c, rel)
			}
		}
		v[rel] = c
	}
	return v, nil
}

// Discover will find all modules starting from the path provided when constructing the Discoverer.
// Does not iterate into testdata folders.
func (d *Discoverer) Discover() error {
	d.modules = make(map[string][]string)

	present, err := IsGoModPresent(d.path)
	if err != nil {
		return err
	}

	err = filepath.Walk(d.path, d.walkChildModules(d.path, present))
	if err != nil {
		return err
	}

	for modulePath := range d.modules {
		if len(d.modules) > 0 {
			sort.Strings(d.modules[modulePath])
		}
	}

	return nil
}

func (d *Discoverer) walkChildModules(parentPath string, isParentModule bool) func(path string, fs os.FileInfo, err error) error {
	if isParentModule {
		d.modules[parentPath] = nil
	}

	return func(path string, fs os.FileInfo, err error) error {
		if err != nil || path == parentPath {
			return err
		}

		if !fs.IsDir() {
			return nil
		}

		if fs.Name() == testDataFolder {
			return filepath.SkipDir
		}

		present, err := IsGoModPresent(path)
		if err != nil {
			return err
		}

		if !present {
			return nil
		}

		if isParentModule {
			d.modules[parentPath] = append(d.modules[parentPath], path)
		}

		err = filepath.Walk(path, d.walkChildModules(path, true))
		if err != nil {
			return err
		}

		return filepath.SkipDir
	}
}

// IsGoModPresent returns whether there is a go.mod file located in the provided directory path
func IsGoModPresent(path string) (bool, error) {
	_, err := os.Stat(filepath.Join(path, goModuleFile))
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// IsSubmodulePath determines if the given path falls within any of the submodules. Submodules MUST be a
// sorted ascending list of paths.
func IsSubmodulePath(path string, submodules []string) bool {
	i := sort.Search(len(submodules), func(i int) bool {
		return path <= submodules[i]
	})

	// Search returns where we would insert the given path, so we need to check if the returned index
	// module matches our path, or if the previous index entry is a prefix to our current directory since
	// nested directory paths would be sorted higher lexicographically
	if (i < len(submodules) && path == submodules[i]) || i > 0 && strings.HasPrefix(path, submodules[i-1]) {
		return true
	}

	return false
}
