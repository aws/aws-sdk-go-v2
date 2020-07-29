package changes

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes/git"
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes/golist"
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes/util"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"golang.org/x/mod/sumdb/dirhash"
	"golang.org/x/mod/zip"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var sdkRepo = "github.com/aws/aws-sdk-go-v2"

//var sdkRepo = "github.com/aggagen/test"

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
func discoverModules(golist golist.ModuleClient, root string) ([]string, map[string]string, error) {
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

			mod := shortenModPath(modFile.Module.Mod.String())

			modPackages, err := golist.Packages(mod)
			if err != nil {
				return err
			}

			for _, p := range modPackages {
				packages[p] = mod
			}

			modules = append(modules, mod)
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
		major = "v0"
	}

	return fmt.Sprintf("%s.0.0", major), nil
}

func pseudoVersion(repoPath, mod string) (string, error) {
	commitHash, err := commitHash(repoPath)
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	tagVer, err := taggedVersion(git.Client{RepoPath: repoPath}, mod, true) // TODO: use actual repo git client
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

	output, err := util.ExecAt(cmd, repoPath)
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	return strings.Trim(string(output), "\n"), nil // clean up git show output and return
}

func goChecksum(repoPath, mod, version string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "modfile-zip")
	if err != nil {
		return "", err
	}

	defer os.Remove(tmpfile.Name())

	err = zip.CreateFromDir(tmpfile, module.Version{
		Path:    lengthenModPath(mod),
		Version: version,
	}, filepath.Join(repoPath, mod))
	if err != nil {
		return "", err
	}

	return dirhash.HashZip(tmpfile.Name(), dirhash.DefaultHash)
}
