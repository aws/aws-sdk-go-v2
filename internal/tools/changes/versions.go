package changes

import (
	"fmt"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"os/exec"
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
	Module  string // Module is the full module path of the Go module.
	Version string // Version is a valid Go module semantic version, which can potentially be a pseudo-version.
}

// VersionEnclosure is a set of versions for Go modules in a given repository.
type VersionEnclosure struct {
	SchemaVersion  string             // SchemaVersion is the version of the library's types used to create this VersionEnclosure
	ModuleVersions map[string]Version // ModuleVersions is a mapping between full module paths and their corresponding Version.
	Packages       map[string]string  // Packages maps each package in the repo to the module that provides the package.
}

// isValid checks whether the ModuleVersions contained in the VersionEnclosure v accurately reflect the latest tagged versions.
func (v VersionEnclosure) isValid() error {
	for m, encVer := range v.ModuleVersions {
		gitVer, err := taggedVersion("", m)
		if err != nil {
			return err
		}

		if encVer.Version != gitVer {
			return fmt.Errorf("module %s enclosure version %s does not match git tag version %s", m, encVer, gitVer)
		}
	}

	return nil
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

func taggedVersion(root, mod string) (string, error) {
	path, major, ok := module.SplitPathVersion(mod)
	if !ok {
		return "", fmt.Errorf("couldn't split module path: %s", mod)
	}

	major = strings.TrimLeft(major, "/")

	var versions []string
	var err error

	if major == "" {
		// if there is no major version suffix, then the latest version could be v1 or v0.
		versions, err = versionTags(root, path, "v1")
		if err != nil {
			return "", err
		}

		if len(versions) == 0 {
			versions, err = versionTags(root, path, "v0")
			if err != nil {
				return "", err
			}
		}
	} else {
		versions, err = versionTags(root, path, major)
		if err != nil {
			return "", err
		}
	}

	if len(versions) == 0 {
		return "", nil
	}

	return versions[0], nil
}

func versionTags(root, mod, major string) ([]string, error) {
	if mod == RootModule {
		mod = ""
	} else {
		mod += "/"
	}

	// --sort=-v:refnam flag sorts the tags by descending version
	cmd := exec.Command("git", "tag", "--sort=-v:refname", "-l", mod+major+"*")
	output, err := execAt(cmd, root)
	if err != nil {
		return nil, err
	}

	var versions []string

	for _, v := range strings.Split(string(output), "\n") {
		if semver.IsValid(v) && semver.Prerelease(v) == "" && semver.Build(v) == "" {
			versions = append(versions, strings.TrimPrefix(v, mod))
		}
	}

	return versions, nil
}
