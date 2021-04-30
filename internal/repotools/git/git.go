package git

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/semver"
	"golang.org/x/mod/module"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
)

// Tags returns a slice of Git tags at the repository located at path
func Tags(path string) ([]string, error) {
	output, err := Git(path, "tag", "-l")
	if err != nil {
		return nil, err
	}
	return splitOutput(string(output)), nil
}

// Fetch fetches all objects and refs for the Git repository located at path
func Fetch(path string) error {
	_, err := Git(path, "fetch", "--all")
	return err
}

// Git executes the git with the provided arguments. The command is executed in the provided
// directory path.
func Git(path string, arguments ...string) (output []byte, err error) {
	cmd := exec.Command("git", arguments...)
	if len(path) == 0 {
		path, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	cmd.Dir = path
	cmd.Env = append(os.Environ(), "PWD="+path)

	return cmd.Output()
}

// ToModuleTag converts the relative module path and semver version string to a git tag
// that can be used to identify the module version.
// For example:
//   Path: .              Version: v1.2.3 => v1.2.3
//   Path: service/s3     Version: v0.2.3 => service/s3/v0.2.3
//   Path: service/s3     Version: v1.2.3 => service/s3/v1.2.3
//   Path: service/s3/v2  Version: v2.2.3 => service/s3/v2.2.3
//   Path: service/s3/v3  Version: v2.2.3 => error
func ToModuleTag(modulePath string, version string) (string, error) {
	major := semver.Major(version)
	if len(major) == 0 {
		return "", fmt.Errorf("invalid semantic version: %v", major)
	}

	prefix, pathMajor, ok := module.SplitPathVersion(modulePath)
	if !ok {
		return "", fmt.Errorf("invalid module path version")
	}

	if err := module.CheckPathMajor(version, pathMajor); err != nil {
		return "", err
	}

	return path.Join(prefix, version), nil
}

// ModuleTags is a map of module paths to a slice of tagged Go semver versions.
// Root module tags will be placed in the map at ".". Major versions > v1 will be
// added to the map with the semver major version appended to the module path.
//
// Versions will be sorted in the slice from highest to lowest by comparing the values
// following semantic versioning rules.
//
// Example:
//   . => ["v1.2.3", "v1.0.0"]
//   v2 => ["v2.0.0"]
//   sub/module => ["v1.2.3"]
//   sub/module/v2 => ["v2.2.3"]
//
type ModuleTags map[string][]string

// Latest returns the latest tag for the given relative module path. Returns false if
// the module version is not known.
func (r ModuleTags) Latest(module string) (string, bool) {
	_, ok := r[module]
	if !ok {
		return "", false
	}
	return r[module][0], true
}

// Add adds the given tag to the ModuleTags
func (r ModuleTags) Add(tag string) bool {
	module, version, ok := parseTag(tag)
	if !ok {
		return false
	}

	index := sort.Search(len(r[module]), func(i int) bool {
		return semver.Compare(version, r[module][i]) >= 0
	})

	if index < len(r[module]) && index >= 0 {
		if semver.Compare(r[module][index], version) == 0 {
			return true
		}
	}

	r[module] = append(r[module], "")
	copy(r[module][index+1:], r[module][index:])
	r[module][index] = version

	return true
}

// ParseModuleTags parses a list of Git tags into a set of ModuleTags.
// Tags that are not semvar compliant with Go will be ignored.
func ParseModuleTags(tags []string) ModuleTags {
	modules := make(map[string][]string)

	for _, tag := range tags {
		module, version, ok := parseTag(tag)
		if !ok {
			continue
		}
		modules[module] = append(modules[module], version)
	}

	for _, versions := range modules {
		sort.Slice(versions, func(i, j int) bool {
			// We want to sort higher versions first
			return semver.Compare(versions[i], versions[j]) > 0
		})
	}

	return modules
}

func parseTag(tag string) (string, string, bool) {
	idx := strings.LastIndex(tag, "/")

	module := "."
	version := tag

	if idx != -1 {
		module = tag[:idx]
		version = tag[idx+1:]
	}

	if !semver.IsValid(version) {
		return "", "", false
	}

	major := semver.Major(version)

	majorInt, _ := strconv.Atoi(major[1:])

	if majorInt > 1 {
		module = path.Join(module, major)
	}

	return module, version, true
}
