package changes

import (
	"fmt"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
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

	return strings.TrimPrefix(modulePath, sdkRepo+"/")
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
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return modules, packages, nil
}

func listPackages(path string) ([]string, error) {
	cmd := exec.Command("go", "list", "./...")
	out, err := execAt(cmd, path)
	if err != nil {
		return nil, err
	}

	packages := strings.Split(string(out), "\n")

	if len(packages) > 0 && packages[len(packages)-1] == "" {
		packages = packages[:len(packages)-1] // remove the empty package caused by the last newline from go list
	}

	return packages, nil
}

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

func PseudoVersion() (string, error) {
	cmd := exec.Command("git", "show", "--quiet", "--abbrev=12", "--date='format-local:%Y%m%d%H%M%S'", "--format='%cd-%h'")
	//cmd := exec.Command("git", "--no-pager", "show", "--quiet", "--abbrev=12",
	//	"--date='format-local:%Y%m%d%H%M%S'", `--format="%cd-%h"`)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "TZ=UTC")

	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	fmt.Println(string(output))
	return "", nil
}
