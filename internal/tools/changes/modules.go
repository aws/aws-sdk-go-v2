package changes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const sdkRepo = "github.com/aws/aws-sdk-go-v2"
const RootModule = "/"

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
		return RootModule
	}

	modulePath = strings.TrimLeft(modulePath, "/")
	return strings.TrimPrefix(modulePath, sdkRepo+"/")
}

func lengthenModPath(modulePath string) string {
	if modulePath == RootModule {
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

	tagVer, err := taggedVersion(repoPath, mod)
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	return formatPseudoVersion(commitHash, tagVer)
}

func formatPseudoVersion(commitHash, taggedVersion string) (string, error) {
	// https://golang.org/cmd/go/#hdr-Pseudo_versions
	if taggedVersion == "" {
		return fmt.Sprintf("v0.0.0-%s", commitHash), nil
	}

	// TODO: Handle prereleases psuedo-version

	taggedVersion, err := nextVersion(taggedVersion, PatchBump)
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	return fmt.Sprintf("%s-0.%s", taggedVersion, commitHash), nil
}

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
