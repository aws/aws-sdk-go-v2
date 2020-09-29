package golist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/util"
	"golang.org/x/mod/module"
	"golang.org/x/mod/sumdb/dirhash"
	"golang.org/x/mod/zip"
)

// Client gets information about Go modules and packages.
type Client struct {
	RootPath        string              // RootPath is the path to the root of a multi-module git repository.
	ShortenModPath  func(string) string // ShortenModPath shortens a module's import path to be a relative path from the RootPath.
	LengthenModPath func(string) string // LengthenModPath lengthens a module's import path to be the full import path.
}

// ModuleClient gets dependency and package information about go modules.
type ModuleClient interface {
	Dependencies(mod string) ([]string, error)
	Packages(mod string) ([]string, error)
	Checksum(mod, version string) (string, error)
	Tidy(mod string) error
}

func (c Client) path(mod string) string {
	parts := []string{c.RootPath}
	parts = append(parts, strings.Split(mod, "/")...)

	return filepath.Join(parts...)
}

// Dependencies returns a list of all modules that the module mod depends on.
func (c Client) Dependencies(mod string) ([]string, error) {
	mod = c.ShortenModPath(mod)

	cmd := exec.Command("go", "list", "-json", "-m", "all")
	out, err := util.ExecAt(cmd, c.path(mod))
	if err != nil {
		return nil, err
	}

	return c.parseGoModuleList(out)
}

// goModule is a package as output by the `go list` command.
type goModule struct {
	Path string // Path is the module's import path.
	Main bool   // Main indicates whether the module is the main module.
}

func (c Client) parseGoModuleList(output []byte) ([]string, error) {
	var modules []string
	dec := json.NewDecoder(bytes.NewReader(output))

	for {
		var p goModule
		if err := dec.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if !p.Main && c.ShortenModPath(p.Path) != p.Path {
			modules = append(modules, c.ShortenModPath(p.Path))
		}
	}

	return modules, nil
}

// Packages returns a slice of packages that are part of the module mod.
func (c Client) Packages(mod string) ([]string, error) {
	mod = c.ShortenModPath(mod)

	packages := map[string]struct{}{}

	absRoot, err := filepath.Abs(c.RootPath)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(c.path(mod), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}

		hasGoFile := false
		for _, f := range files {
			if f.Name() == "go.mod" && c.path(mod) != path {
				return filepath.SkipDir
			} else if strings.HasSuffix(f.Name(), ".go") {
				hasGoFile = true
			}
		}

		if !hasGoFile {
			return nil
		}

		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, absRoot)
		path = strings.TrimLeft(path, "/")
		parts := strings.Split(filepath.ToSlash(path), "/")

		p := c.LengthenModPath(strings.Join(parts, "/"))

		packages[p] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, err
	}

	packageList := []string{}
	for p := range packages {
		packageList = append(packageList, p)
	}

	return packageList, nil
}

// goPackage is a package as output by the `go list` command.
type goPackage struct {
	ImportPath string // ImportPath is the package's import path.
}

func parseGoList(output []byte) ([]string, error) {
	var packages []string
	dec := json.NewDecoder(bytes.NewReader(output))

	for {
		var p goPackage
		if err := dec.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		packages = append(packages, p.ImportPath)
	}

	return packages, nil
}

// Checksum returns the given module's Go checksum at the specified version.
func (c Client) Checksum(mod, version string) (string, error) {
	mod = c.ShortenModPath(mod)

	tmpfile, err := ioutil.TempFile("", "modfile-zip")
	if err != nil {
		return "", err
	}

	defer os.Remove(tmpfile.Name())

	err = zip.CreateFromDir(tmpfile, module.Version{
		Path:    c.LengthenModPath(mod),
		Version: version,
	}, filepath.Join(c.RootPath, mod))
	if err != nil {
		return "", err
	}

	return dirhash.HashZip(tmpfile.Name(), dirhash.DefaultHash)
}

// Tidy runs go mod tidy on the specified module.
func (c Client) Tidy(mod string) error {
	cmd := exec.Command("go", "mod", "tidy")
	_, err := util.ExecAt(cmd, c.path(mod))
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %v", err)
	}

	return nil
}
