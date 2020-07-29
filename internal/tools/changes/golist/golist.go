package golist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes/util"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Client struct {
	RootPath       string
	ShortenModPath func(string) string
}

type ModuleClient interface {
	Dependencies(mod string) ([]string, error)
	Packages(mod string) ([]string, error)
}

func (c Client) path(mod string) string {
	parts := []string{c.RootPath}
	parts = append(parts, strings.Split(mod, "/")...)

	return filepath.Join(parts...)
}

func (c Client) Dependencies(mod string) ([]string, error) {
	mod = c.ShortenModPath(mod)

	cmd := exec.Command("go", "list", "-json", "-m", "all")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOSUMDB=off")

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

// ListPackages returns a slice of packages that are part of the module mod.
func (c Client) Packages(mod string) ([]string, error) {
	mod = c.ShortenModPath(mod)

	cmd := exec.Command("go", "list", "-json", "./...")
	out, err := util.ExecAt(cmd, c.path(mod))
	if err != nil {
		return nil, fmt.Errorf("failed to list packages: %v", err)
	}

	return parseGoList(out)
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
