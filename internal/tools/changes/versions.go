package changes

import (
	"fmt"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// VersionIncrement describes how a Change should affect a module's version.
type VersionIncrement int

const (
	NoBump    VersionIncrement = iota // NoBump indicates the module's version should not change.
	PatchBump                         // PatchBump indicates the module's version should be incremented by a patch version bump.
	MinorBump                         // MinorBump indicates the module's version should be incremented by a minor version bump.
)

// Version is the version of a Go module.
type Version struct {
	Module     string // Module is the repo relative module path of the Go module.
	ImportPath string // ImportPath is the full module import path.
	Version    string // Version is a valid Go module semantic version, which can potentially be a pseudo-version.
}

// VersionEnclosure is a set of versions for Go modules in a given repository.
type VersionEnclosure struct {
	SchemaVersion  int                // SchemaVersion is the version of the library's types used to create this VersionEnclosure
	ModuleVersions map[string]Version // ModuleVersions is a mapping between shortened module paths and their corresponding Version.
	Packages       map[string]string  // Packages maps each package in the repo to the shortened module path that provides the package.
}

// isValid returns nil if the ModuleVersions contained in the VersionEnclosure v accurately reflect the latest tagged versions.
// Otherwise, isValid returns an error.
func (v VersionEnclosure) isValid(repoPath string) error {
	for m, encVer := range v.ModuleVersions {
		gitVer, err := taggedVersion(repoPath, m, false)
		if err != nil {
			return err
		}

		if encVer.Version != gitVer {
			return fmt.Errorf("module %s enclosure version %s does not match git tag version %s", m, encVer, gitVer)
		}
	}

	return nil
}

func versionIncrement(changes []Change) VersionIncrement {
	maxBump := NoBump
	for _, c := range changes {
		bump := c.Type.VersionIncrement()
		if bump > maxBump {
			maxBump = bump
		}
	}

	return maxBump
}

func nextVersion(version string, bumpType VersionIncrement) (string, error) {
	if !semver.IsValid(version) {
		return "", fmt.Errorf("version %s is not valid", version)
	}
	if semver.Prerelease(version) != "" || semver.Build(version) != "" {
		return "", fmt.Errorf("version %s has a prerelease or build component", version)
	}

	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("expected 3 semver parts, got %d", len(parts))
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", err
	}

	switch bumpType {
	case NoBump:
		return version, nil
	case PatchBump:
		patch += 1
	case MinorBump:
		patch = 0
		minor += 1
	}

	return fmt.Sprintf("%s.%d.%d", parts[0], minor, patch), nil
}

// taggedVersion returns the latest tagged version of the given module in the specified repository.
func taggedVersion(repoPath, mod string, includePrereleases bool) (string, error) {
	path, major, ok := module.SplitPathVersion(mod)
	if !ok {
		return "", fmt.Errorf("couldn't split module path: %s", mod)
	}

	major = strings.TrimLeft(major, "/")

	var versions []string
	var err error

	if major == "" {
		// if there is no major version suffix, then the latest version could be v1 or v0.
		versions, err = versionTags(repoPath, path, "v1", includePrereleases)
		if err != nil {
			return "", err
		}

		if len(versions) == 0 {
			versions, err = versionTags(repoPath, path, "v0", includePrereleases)
			if err != nil {
				return "", err
			}
		}
	} else {
		versions, err = versionTags(repoPath, path, major, includePrereleases)
		if err != nil {
			return "", err
		}
	}

	if len(versions) == 0 {
		return "", nil
	}

	return versions[0], nil
}

// versionTags gets all semantic version git tags for the given module major version, ignoring prerelease versions.
func versionTags(repoPath, mod, major string, includePrereleases bool) ([]string, error) {
	if mod == rootModule {
		mod = ""
	} else {
		mod += "/"
	}

	cmd := exec.Command("git", "tag", "--sort=-v:refname", "-l", mod+major+"*")
	output, err := execAt(cmd, repoPath)
	if err != nil {
		return nil, err
	}

	var versions []string

	for _, v := range strings.Split(string(output), "\n") {
		v = strings.TrimPrefix(v, mod)
		prerelease := semver.Prerelease(v) != ""

		if (semver.IsValid(v) && semver.Build(v) == "") && (!prerelease || includePrereleases) {
			versions = append(versions, strings.TrimPrefix(v, mod))
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare(versions[i], versions[j]) > 0
	})

	return versions, nil
}
