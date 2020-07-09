package changes

import (
	"fmt"
	"path/filepath"
)

// Repository is a representation of a git repository containing multiple Go modules.
type Repository struct {
	RootPath string    // RootPath is the path to the root of the repository.
	Metadata *Metadata // Metadata is the repository's .changes metadata.
}

// NewRepository loads the repository at the given path.
func NewRepository(path string) (*Repository, error) {
	metadata, err := LoadMetadata(filepath.Join(path, metadataDir))
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		RootPath: path,
		Metadata: metadata,
	}

	return repo, nil
}

// UpdateChangelog generates a new CHANGELOG entry for the given release. If pending is true, the contents of
// CHANGELOG_PENDING.md will be replaced with the new entry. Otherwise, the entry will be prepended to CHANGELOG.md.
func (r *Repository) UpdateChangelog(release *Release, pending bool) error {
	id := release.ID
	fileName := "CHANGELOG.md"
	if pending {
		id = "pending"
		fileName = "CHANGELOG_PENDING.md"
	}

	entry := fmt.Sprintf("# Release %s\n", id)

	for _, module := range release.AffectedModules() {
		moduleEntry, err := release.RenderChangelogForModule(module, "#")
		if err != nil {
			return fmt.Errorf("couldn't update changelog: %v", err)
		}

		entry += moduleEntry
	}

	return writeFile([]byte(entry), filepath.Join(r.RootPath, fileName), !pending)
}

// UpdatePendingChangelog updates the repository's top level CHANGELOG_PENDING.md with the Repository's Metadata's
// pending Changes.
func (r *Repository) UpdatePendingChangelog() error {
	release, err := r.Metadata.CreateRelease("pending", map[string]VersionBump{}, true)
	if err != nil {
		return fmt.Errorf("couldn't update pending changelog: %v", err)
	}

	return r.UpdateChangelog(release, true)
}

// InitializeVersions creates an initial versions.json enclosure in the Repository's .changes directory. The VersionEnclosure
// created by InitializeVersions will include all modules that have a tagged version in git and their corresponding versions.
func (r *Repository) InitializeVersions() error {
	enc, err := r.DiscoverVersions(TaggedVersionSelector)
	if err != nil {
		return fmt.Errorf("failed to initialize versions: %v", err)
	}

	return r.Metadata.SaveEnclosure(enc)
}

// VersionSelector is a function that decides what version of a Go module should be passed to code generation.
type VersionSelector func(r *Repository, module string) (string, error)

func (r *Repository) discoverVersions(modules []string, selector VersionSelector) (VersionEnclosure, error) {
	enclosure := VersionEnclosure{
		SchemaVersion:  SchemaVersion,
		ModuleVersions: map[string]Version{},
	}

	for _, m := range modules {
		v, err := selector(r, m)
		if err != nil {
			return VersionEnclosure{}, err
		}

		if v != "" {
			enclosure.ModuleVersions[m] = Version{
				Module:  m,
				Version: v,
			}
		}
	}

	return enclosure, nil
}

// DiscoverVersions creates a VersionEnclosure containing all Go modules in the Repository. The version of each module
// is determined by the provided VersionSelector.
func (r *Repository) DiscoverVersions(selector VersionSelector) (VersionEnclosure, error) {
	modules, packages, err := discoverModules(r.RootPath)
	if err != nil {
		return VersionEnclosure{}, err
	}

	enc, err := r.discoverVersions(modules, selector)
	if err != nil {
		return VersionEnclosure{}, fmt.Errorf("failed to discover versions: %v", err)
	}

	enc.Packages = packages

	return enc, nil
}

// ReleaseVersionSelector returns a version for the given module suitable for use during the release process.
// A version will be returned based upon what type of version bump the Change metadata for the given module requires.
// ReleaseVersionSelector will properly version modules that are not present in the versions.json file by checking git
// tags for an existing version, or by providing a default version suitable for the module's major version.
func ReleaseVersionSelector(r *Repository, module string) (string, error) {
	incr := versionIncrement(r.Metadata.GetChanges(module))

	currentVersion := r.Metadata.CurrentVersions.ModuleVersions[module].Version
	if currentVersion != "" {
		return nextVersion(currentVersion, incr)
	}

	v, err := taggedVersion(r.RootPath, module)
	if err != nil {
		return "", fmt.Errorf("couldn't find current version of %s: %v", module, err)
	}
	if v == "" {
		// there aren't version git tags for this module
		return defaultVersion(module)
	}

	// the module isn't in versions.json, but does have git tags
	return nextVersion(v, incr)
}

// TaggedVersionSelector returns the greatest version of module tagged in the git repository.
func TaggedVersionSelector(r *Repository, module string) (string, error) {
	return taggedVersion(r.RootPath, module)
}

// DevelopmentVersionSelector returns a commit hash based version if the module has associated pending Changes, otherwise
// returns the latest version from the repo's metadata version enclosure.
func DevelopmentVersionSelector(r *Repository, module string) (string, error) {
	incr := versionIncrement(r.Metadata.GetChanges(module))

	if incr != NoBump {
		return PseudoVersion(r.RootPath, module)
	}

	if v, ok := r.Metadata.CurrentVersions.ModuleVersions[module]; ok {
		return v.Version, nil
	}

	return "", fmt.Errorf("couldn't select version for module %s: module has no changes and versions.json version", module)
}
