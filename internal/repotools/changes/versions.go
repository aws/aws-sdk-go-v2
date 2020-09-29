package changes

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/git"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

// VersionIncrement describes how a Change should affect a module's version.
type VersionIncrement int

const (
	// NoBump indicates the module's version should not change.
	NoBump VersionIncrement = iota
	// PatchBump indicates the module's version should be incremented by a patch version bump.
	PatchBump
	// MinorBump indicates the module's version should be incremented by a minor version bump.
	MinorBump
	// MajorBump indicates the module's major version should be incremented from v0 to v1.
	MajorBump
	// NewModule indicates that a new modules has been discovered and will be assigned a default version.
	NewModule
)

// Version is the version of a Go module.
type Version struct {
	Module     string // Module is the repo relative module path of the Go module.
	ImportPath string // ImportPath is the full module import path.
	Version    string // Version is a valid Go module semantic version, which can potentially be a pseudo-version.
	ModuleHash string // ModuleHash is the Go module checksum of the the Version's module.
}

// VersionBump describes a version increment to a module.
type VersionBump struct {
	From string // From is the old module version
	To   string // To is the new module version
}

// VersionEnclosure is a set of versions for Go modules in a given repository.
type VersionEnclosure struct {
	SchemaVersion  int                // SchemaVersion is the version of the library's types used to create this VersionEnclosure
	ModuleVersions map[string]Version // ModuleVersions is a mapping between shortened module paths and their corresponding Version.
	Packages       map[string]string  // Packages maps each package in the repo to the shortened module path that provides the package.
}

// isValid returns nil if the ModuleVersions contained in the VersionEnclosure v accurately reflect the latest tagged versions.
// Otherwise, isValid returns an error.
func (v *VersionEnclosure) isValid(git git.VcsClient) error {
	for m, encVer := range v.ModuleVersions {
		gitVer, err := taggedVersion(git, m, false)
		if err != nil {
			return err
		}

		if encVer.Version != gitVer {
			return fmt.Errorf("module %s enclosure version %s does not match git tag version %s", m, encVer, gitVer)
		}
	}

	return nil
}

func (v *VersionEnclosure) bump(module string, incr VersionIncrement) (VersionBump, error) {
	if _, ok := v.ModuleVersions[module]; !ok {
		return VersionBump{}, fmt.Errorf("the VersionEnclosure doesn't contain module %s", module)
	}

	ver := v.ModuleVersions[module]
	oldVer := ver.Version
	nextVer, err := nextVersion(oldVer, incr, "")
	if err != nil {
		return VersionBump{}, fmt.Errorf("couldn't bump module %s's version: %v", module, err)
	}

	ver.Version = nextVer

	v.ModuleVersions[module] = ver

	return VersionBump{
		From: oldVer,
		To:   nextVer,
	}, nil
}

func (v *VersionEnclosure) updateHashes(hashes map[string]string) error {
	for mod, hash := range hashes {
		if ver, ok := v.ModuleVersions[mod]; ok {
			ver.ModuleHash = hash
			v.ModuleVersions[mod] = ver
		} else {
			return fmt.Errorf("module %s is contained in hashes, but not enclosure", mod)
		}
	}

	return nil
}

// HashDiff returns all modules whose hash provided in hashes differs from the has present in VersionEnclosure v. hashes
// is a map between shortened module names and their Go checksum.
func (v *VersionEnclosure) HashDiff(hashes map[string]string) []string {
	var diff []string

	for mod, hash := range hashes {
		if v.ModuleVersions[mod].ModuleHash != hash {
			diff = append(diff, mod)
		}
	}

	return diff
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

func nextVersion(version string, bumpType VersionIncrement, minTargetVersion string) (string, error) {
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

	major, err := strconv.Atoi(strings.TrimPrefix(parts[0], "v"))
	if err != nil {
		return "", err
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
		patch++
	case MinorBump:
		patch = 0
		minor++
	case MajorBump:
		if major != 0 {
			return "", errors.New("major increment can only be applied to v0")
		}

		major = 1
		minor = 0
		patch = 0
	}

	next := fmt.Sprintf("v%d.%d.%d", major, minor, patch)

	if len(minTargetVersion) > 0 && semver.Compare(next, minTargetVersion) == -1 {
		next = minTargetVersion
	}

	return next, nil
}

func tagRepo(git git.VcsClient, releaseID, mod, version string) error {
	path, major, ok := module.SplitPathVersion(mod)
	if !ok {
		return fmt.Errorf("couldn't split module path: %s", mod)
	}

	major = strings.TrimLeft(major, "/")
	if !strings.HasPrefix(version, major) {
		return fmt.Errorf("version %s does not match module %s's major version", version, mod)
	}

	if path == rootModule || path == "" {
		path = ""
	} else {
		path += "/"
	}

	tag := path + version

	msg := fmt.Sprintf("Release %s", releaseID)

	return git.Tag(tag, msg)
}

// taggedVersion returns the latest tagged version of the given module in the specified repository.
func taggedVersion(git git.VcsClient, mod string, includePrereleases bool) (string, error) {
	path, major, ok := module.SplitPathVersion(mod)
	if !ok {
		return "", fmt.Errorf("couldn't split module path: %s", mod)
	}

	major = strings.TrimLeft(major, "/")

	var versions []string
	var err error

	if major == "" {
		// if there is no major version suffix, then the latest version could be v1 or v0.
		versions, err = versionTags(git, path, "v1", includePrereleases)
		if err != nil {
			return "", err
		}

		if len(versions) == 0 {
			versions, err = versionTags(git, path, "v0", includePrereleases)
			if err != nil {
				return "", err
			}
		}
	} else {
		versions, err = versionTags(git, path, major, includePrereleases)
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
func versionTags(git git.VcsClient, mod, major string, includePrereleases bool) ([]string, error) {
	if mod == rootModule {
		mod = ""
	} else {
		mod += "/"
	}

	versions, err := git.Tags(mod + major)
	if err != nil {
		return nil, err
	}

	versions = filterVersions(versions, mod, includePrereleases)

	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare(versions[i], versions[j]) > 0
	})

	return versions, nil
}

func filterVersions(versions []string, mod string, includePrereleases bool) []string {
	var filtered []string

	for _, v := range versions {
		v = strings.TrimPrefix(v, mod)
		prerelease := semver.Prerelease(v) != ""

		if (semver.IsValid(v) && semver.Build(v) == "") && (!prerelease || includePrereleases) {
			filtered = append(filtered, v)
		}
	}

	return filtered
}
