package changes

import "fmt"

// VersionSelector is a function that decides what version of a Go module should be passed to code generation.
type VersionSelector func(r *Repository, module string) (string, VersionIncrement, error)

// Set parses the input s and correspondingly sets the appropriate VersionSelector.
func (v *VersionSelector) Set(s string) error {
	switch s {
	case "release":
		*v = ReleaseVersionSelector
	case "development":
		*v = DevelopmentVersionSelector
	case "tags":
		*v = TaggedVersionSelector
	default:
		return fmt.Errorf("unknown version selector: %s", s)
	}

	return nil
}

// String returns an empty string to satisfy the flag.Value interface.
func (v *VersionSelector) String() string {
	return ""
}

// ReleaseVersionSelector returns a version for the given module suitable for use during the release process.
// A version will be returned based upon what type of version bump the Change metadata for the given module requires.
// ReleaseVersionSelector will properly version modules that are not present in the versions.json file by checking git
// tags for an existing version, or by providing a default version suitable for the module's major version.
func ReleaseVersionSelector(r *Repository, module string) (string, VersionIncrement, error) {
	incr := versionIncrement(r.Metadata.GetChanges(module))

	currentVersion := r.Metadata.CurrentVersions.ModuleVersions[module].Version
	if currentVersion != "" {
		v, err := nextVersion(currentVersion, incr)
		return v, incr, err
	}

	v, err := taggedVersion(r.git, module, false)
	if err != nil {
		return "", NoBump, fmt.Errorf("couldn't find current version of %s: %v", module, err)
	}
	if v == "" {
		// there aren't version git tags for this module
		v, err = defaultVersion(module)
		return v, NewModule, err
	}

	// the module isn't in versions.json, but does have git tags
	v, err = nextVersion(v, incr)
	return v, incr, err
}

// TaggedVersionSelector returns the greatest version of module tagged in the git repository.
func TaggedVersionSelector(r *Repository, module string) (string, VersionIncrement, error) {
	v, err := taggedVersion(r.git, module, false)
	return v, NoBump, err
}

// DevelopmentVersionSelector returns a commit hash based version if the module has associated pending Changes, otherwise
// returns the latest version from the repo's metadata version enclosure.
func DevelopmentVersionSelector(r *Repository, module string) (string, VersionIncrement, error) {
	incr := versionIncrement(r.Metadata.GetChanges(module))

	if incr != NoBump {
		v, err := pseudoVersion(r.git, module)
		return v, incr, err
	}

	if v, ok := r.Metadata.CurrentVersions.ModuleVersions[module]; ok {
		return v.Version, incr, nil
	}

	return "", NoBump, fmt.Errorf("couldn't select version for module %s: module has no changes and no versions.json version", module)
}
