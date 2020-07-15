package changes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const sdkRepo = "github.com/aws/aws-sdk-go-v2"
const rootModule = "/"

// GetCurrentModule returns a shortened module path (from the root of the repository to the module, not a full import
// path) for the Go module containing the current directory.
func GetCurrentModule() (string, error) {
	path, err := findFile("go.mod", false)
	if err != nil {
		return "", err
	}

	modFile, err := getModFile(path)
	if err != nil {
		return "", err
	}

	return shortenModPath(modFile.Module.Mod.String()), nil
}

func getModFile(path string) (*modfile.File, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return modfile.ParseLax(path, data, nil)
}

func shortenModPath(modulePath string) string {
	if modulePath == sdkRepo {
		return rootModule
	}

	modulePath = strings.TrimLeft(modulePath, "/")
	return strings.TrimPrefix(modulePath, sdkRepo+"/")
}

func lengthenModPath(modulePath string) string {
	if modulePath == rootModule {
		return sdkRepo
	}

	return sdkRepo + "/" + modulePath
}

// discoverModules returns a list of all modules and a map between all packages and their providing module.
func discoverModules(root string) ([]string, map[string]string, error) {
	var modules []string
	packages := map[string]string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "go.mod" {
			modFile, err := getModFile(path)
			if err != nil {
				return err
			}

			mod := modFile.Module.Mod.String()

			modPackages, err := listPackages(filepath.Dir(path))
			if err != nil {
				return err
			}

			for _, p := range modPackages {
				packages[p] = mod
			}

			modules = append(modules, shortenModPath(mod))
		} else if info.Name() == "testdata" && path != "testdata" {
			// skip testdata directory unless it is the root directory.
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return modules, packages, nil
}

func UpdateDependencies(repoPath, mod, dependency, version string) error {
	goModPath := filepath.Join(repoPath, mod, "go.mod")

	modFile, err := getModFile(goModPath)
	if err != nil {
		return fmt.Errorf("couldn't update %s's dependency on %s: %v", mod, dependency, err)
	}

	err = modFile.AddRequire(lengthenModPath(dependency), version)
	if err != nil {
		return fmt.Errorf("couldn't update %s's dependency on %s: %v", mod, dependency, err)
	}

	out, err := modFile.Format()
	if err != nil {
		return fmt.Errorf("couldn't update %s's dependency on %s: %v", mod, dependency, err)
	}

	return writeFile(out, goModPath, false)
}

func listDependencies(path string) ([]string, error) {
	cmd := exec.Command("go", "list", "-json", "-m", "all")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOSUMDB=off")

	out, err := execAt(cmd, path)
	if err != nil {
		return nil, err
	}

	return parseGoModuleList(out)
}

// goModule is a package as output by the `go list` command.
type goModule struct {
	Path string // Path is the module's import path.
	Main bool   // Main indicates whether the module is the main module.
}

func parseGoModuleList(output []byte) ([]string, error) {
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

		if !p.Main && strings.HasPrefix(p.Path, sdkRepo) {
			modules = append(modules, shortenModPath(p.Path))
		}
	}

	return modules, nil
}

// listPackages returns a slice of packages that are part of the module whose go.mod file is in the directory specified
// by path.
func listPackages(path string) ([]string, error) {
	cmd := exec.Command("go", "list", "-json", "./...")
	out, err := execAt(cmd, path)
	if err != nil {
		return nil, err
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

// defaultVersion returns a default version for the given module based on its import path. If a version suffix /vX is
// present in the import path, the default version will be vX.0.0. Otherwise, the version will be v0.0.0.
func defaultVersion(mod string) (string, error) {
	_, major, ok := module.SplitPathVersion(mod)
	if !ok {
		return "", fmt.Errorf("couldn't split module path: %s", mod)
	}
	major = strings.TrimLeft(major, "/")

	if major == "" {
		major = "v0"
	}

	return fmt.Sprintf("%s.0.0", major), nil
}

func pseudoVersion(repoPath, mod string) (string, error) {
	commitHash, err := commitHash(repoPath)
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	tagVer, err := taggedVersion(repoPath, mod, true)
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	return formatPseudoVersion(commitHash, tagVer)
}

// formatPseudoVersion returns a Go module pseudo version as described in https://golang.org/cmd/go/#hdr-Pseudo_versions.
// taggedVersion is the latest semantic version tag, which may include a prerelease component.
func formatPseudoVersion(commitHash, taggedVersion string) (string, error) {
	if b := semver.Build(taggedVersion); b != "" {
		return "", fmt.Errorf("expected version to not have build tag, got: %s", b)
	}

	if taggedVersion == "" {
		return fmt.Sprintf("v0.0.0-%s", commitHash), nil
	}

	pre := semver.Prerelease(taggedVersion) // pre includes '-'
	if pre != "" {
		taggedVersion = strings.TrimSuffix(taggedVersion, pre)

		return fmt.Sprintf("%s%s.0.%s", taggedVersion, pre, commitHash), nil
	}

	taggedVersion, err := nextVersion(taggedVersion, PatchBump)
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	return fmt.Sprintf("%s-0.%s", taggedVersion, commitHash), nil
}

// commitHash returns a timestamp and commit hash for the HEAD commit of the given repository, formatted in the way
// expected for a go.mod file pseudo-version.
func commitHash(repoPath string) (string, error) {
	cmd := exec.Command("git", "show", "--quiet", "--abbrev=12", "--date=format-local:%Y%m%d%H%M%S", "--format=%cd-%h")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "TZ=UTC")

	output, err := execAt(cmd, repoPath)
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	return strings.Trim(string(output), "\n"), nil // clean up git show output and return
}
