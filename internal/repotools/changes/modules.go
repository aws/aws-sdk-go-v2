package changes

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/golist"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/util"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

var sdkRepo = "github.com/aws/aws-sdk-go-v2"

const rootModule = "/"

// GetCurrentModule returns a shortened module path (from the root of the repository to the module, not a full import
// path) for the Go module containing the current directory.
func GetCurrentModule() (string, error) {
	path, err := util.FindFile("go.mod", false)
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
	if modulePath == sdkRepo || modulePath == rootModule {
		return rootModule
	}

	modulePath = strings.TrimPrefix(modulePath, "./")
	modulePath = strings.TrimLeft(modulePath, "/")
	return strings.TrimPrefix(modulePath, sdkRepo+"/")
}

func lengthenModPath(modulePath string) string {
	if modulePath == rootModule || modulePath == "" {
		return sdkRepo
	}

	modulePath = strings.TrimLeft(modulePath, "/")
	return sdkRepo + "/" + modulePath
}

func modToPath(mod string) string {
	parts := strings.Split(mod, "/")

	return filepath.Join(parts...)
}

// discoverModules returns a list of all modules within the provided root directory.
func discoverModules(root string) ([]string, error) {
	var modules []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "go.mod" {
			modFile, err := getModFile(path)
			if err != nil {
				return err
			}

			mod := shortenModPath(modFile.Module.Mod.String())
			modules = append(modules, mod)
		} else if info.Name() == "testdata" && path != "testdata" {
			// skip testdata directory unless it is the root directory.
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return modules, nil
}

// packages returns a mapping between each package in the repository and its providing module.
func packages(golist golist.ModuleClient, modules []string) (map[string]string, error) {
	packages := map[string]string{}

	for _, mod := range modules {
		modPackages, err := golist.Packages(mod)
		if err != nil {
			return nil, err
		}

		for _, p := range modPackages {
			packages[p] = mod
		}
	}

	return packages, nil
}

func updateDependencies(repoPath, mod string, dependencies []string, enc *VersionEnclosure) error {
	goModPath := filepath.Join(repoPath, mod, "go.mod")

	modFile, err := getModFile(goModPath)
	if err != nil {
		return fmt.Errorf("couldn't update %s's dependencies: %v", mod, err)
	}

	for _, dep := range dependencies {
		err = modFile.AddRequire(lengthenModPath(dep), enc.ModuleVersions[dep].Version)
		if err != nil {
			return fmt.Errorf("couldn't update %s's dependency on %s: %v", mod, dep, err)
		}
	}

	out, err := modFile.Format()
	if err != nil {
		return fmt.Errorf("couldn't update %s's dependencies: %v", mod, err)
	}

	return util.WriteFile(out, goModPath, false)
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
		return "v0.1.0", nil
	}

	return fmt.Sprintf("%s.0.0", major), nil
}

func pseudoVersion(git git.VcsClient, mod string) (string, error) {
	commitHash, err := git.CommitHash()
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	tagVer, err := taggedVersion(git, mod, true)
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

	taggedVersion, err := nextVersion(taggedVersion, PatchBump, "")
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	return fmt.Sprintf("%s-0.%s", taggedVersion, commitHash), nil
}

func pathMajorVersion(modulePath string) (major string, err error) {
	_, pathMajor, ok := module.SplitPathVersion(modulePath)
	if !ok {
		return "", fmt.Errorf("module path %s contains invalid version componenet", modulePath)
	}
	pathMajor = strings.TrimLeft(pathMajor, "/")
	return pathMajor, nil
}

func validateModulePathSemVer(modulePath, version string) error {
	pathMajor, err := pathMajorVersion(modulePath)
	if err != nil {
		return err
	}
	major := semver.Major(version)
	if (len(pathMajor) == 0 && !(major == "v0" || major == "v1")) || (len(pathMajor) > 0 && pathMajor != major) {
		return fmt.Errorf("%s module path does not match with minimum version string: %s", modulePath, version)
	}
	return nil
}
