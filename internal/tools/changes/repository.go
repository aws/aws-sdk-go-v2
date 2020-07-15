package changes

import (
	"fmt"
	"path/filepath"
)

// Repository is a representation of a git repository containing multiple Go modules.
type Repository struct {
	RootPath string    // RootPath is the path to the root of the repository.
	Metadata *Metadata // Metadata is the repository's .changes metadata.
	modules  []string  // modules contains the shortened module path of all modules in the repository. modules is lazily loaded by Modules()
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

func (r *Repository) Modules() ([]string, error) {
	if r.modules != nil {
		return r.modules, nil
	}

	mods, _, err := discoverModules(r.RootPath)
	if err != nil {
		return nil, err
	}

	r.modules = mods
	return mods, err
}

func (r *Repository) DoRelease() error {
	if err := r.Metadata.CurrentVersions.isValid(r.RootPath); err != nil {
		return fmt.Errorf("couldn't create a release: %v", err)
	}

	enc, bumps, err := r.DiscoverVersions(ReleaseVersionSelector)
	if err != nil {
		return err
	}

	err = r.resolveDependencies(&enc, bumps)
	if err != nil {
		return err
	}

	rel, err := r.Metadata.CreateRelease("2020-07-13", bumps, false)
	if err != nil {
		return err
	}

	err = r.UpdateChangelog(rel, false)
	if err != nil {
		return err
	}

	err = r.Metadata.SaveEnclosure(enc)
	if err != nil {
		return err
	}

	for _, m := range rel.AffectedModules() {
		if m == rootModule {
			continue
		}

		err = r.UpdateModuleChangelog(rel, m, false)
		if err != nil {
			return err
		}
	}

	// TODO commit changelogs and go.mod files, maybe keep track in the Repo struct

	err = commit(r.RootPath, []string{"CHANGELOG.md", ".changes", "services"})
	if err != nil {
		return err
	}

	for mod, bump := range bumps {
		err = tagRepo(r.RootPath, mod, bump.To)
		if err != nil {
			return err
		}
	}

	return push(r.RootPath)
}

// dependencies returns a map between a module and all modules that depend on that module.
func (r *Repository) dependencies() (map[string][]string, error) {
	modules, err := r.Modules()
	if err != nil {
		return nil, err
	}

	deps := map[string][]string{}

	for _, m := range modules {
		mDeps, err := listDependencies(filepath.Join(r.RootPath, m))
		if err != nil {
			return nil, err
		}

		fmt.Println(m, mDeps)
		for _, d := range mDeps {
			if depList, ok := deps[d]; ok {
				deps[d] = append(depList, m)
			} else {
				deps[d] = []string{m}
			}
		}
	}

	return deps, nil
}

func (r *Repository) resolveDependencies(enc *VersionEnclosure, bumps map[string]VersionBump) error {
	seen := make(map[string]struct{})
	deps, err := r.dependencies()
	if err != nil {
		return err
	}

	fmt.Println(deps)

	updatedModules := AffectedModules(r.Metadata.Changes)

	for len(updatedModules) > 0 {
		m := updatedModules[0]
		seen[m] = struct{}{}

		if _, ok := bumps[m]; !ok {
			changes, err := NewChanges([]string{m}, DependencyChangeType, "Updated SDK dependencies to the latest versions.")
			if err != nil {
				return err
			}
			if len(changes) != 1 {
				return fmt.Errorf("expected len(changes) to be 1, got %d", len(changes))
			}

			err = r.Metadata.AddChange(changes[0])
			if err != nil {
				return err // TODO Wrap errs
			}

			// TODO: consolidate this logic?
			oldVer := r.Metadata.CurrentVersions.ModuleVersions[m].Version
			nextVer, err := nextVersion(oldVer, PatchBump)
			if err != nil {
				return err
			}

			bumps[m] = VersionBump{
				From: oldVer,
				To:   nextVer,
			}

			encVer := enc.ModuleVersions[m]
			encVer.Version = nextVer

			enc.ModuleVersions[m] = encVer
		}

		for _, d := range deps[m] {
			if _, ok := seen[d]; !ok {
				err = UpdateDependencies(r.RootPath, d, m, enc.ModuleVersions[m].Version)
				if err != nil {
					return err
				}

				fmt.Printf("updating %s modfile to depend on %s\n", d, m)
				updatedModules = append(updatedModules, d)
			}
		}

		updatedModules = updatedModules[1:]
	}

	return nil
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

func (r *Repository) UpdateModuleChangelog(release *Release, module string, pending bool) error {
	entry, err := release.RenderChangelogForModule(module, "")
	if err != nil {
		return err
	}

	return writeFile([]byte(entry), filepath.Join(r.RootPath, module, "CHANGELOG.md"), !pending)
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
	enc, _, err := r.DiscoverVersions(TaggedVersionSelector)
	if err != nil {
		return fmt.Errorf("failed to initialize versions: %v", err)
	}

	return r.Metadata.SaveEnclosure(enc)
}

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

// DiscoverVersions creates a VersionEnclosure containing all Go modules in the Repository. The version of each module
// is determined by the provided VersionSelector.
func (r *Repository) DiscoverVersions(selector VersionSelector) (VersionEnclosure, map[string]VersionBump, error) {
	modules, packages, err := discoverModules(r.RootPath)
	if err != nil {
		return VersionEnclosure{}, nil, err
	}

	enc, bumps, err := r.discoverVersions(modules, selector)
	if err != nil {
		return VersionEnclosure{}, nil, fmt.Errorf("failed to discover versions: %v", err)
	}

	enc.Packages = packages

	return enc, bumps, nil
}

func (r *Repository) discoverVersions(modules []string, selector VersionSelector) (VersionEnclosure, map[string]VersionBump, error) {
	enclosure := VersionEnclosure{
		SchemaVersion:  SchemaVersion,
		ModuleVersions: map[string]Version{},
	}

	bumps := map[string]VersionBump{}

	for _, m := range modules {
		v, incr, err := selector(r, m)
		if err != nil {
			return VersionEnclosure{}, nil, err
		}

		if v != "" {
			enclosure.ModuleVersions[m] = Version{
				Module:     m,
				ImportPath: lengthenModPath(m),
				Version:    v,
			}
		}

		if incr != NoBump {
			bumps[m] = VersionBump{
				From: r.Metadata.CurrentVersions.ModuleVersions[m].Version,
				To:   v,
			}
		}
	}

	return enclosure, bumps, nil
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

	v, err := taggedVersion(r.RootPath, module, false)
	if err != nil {
		return "", NoBump, fmt.Errorf("couldn't find current version of %s: %v", module, err)
	}
	if v == "" {
		// there aren't version git tags for this module
		v, err = defaultVersion(module)
		return v, PatchBump, err // TODO: PatchBump? or something like ModuleDiscovered
	}

	// the module isn't in versions.json, but does have git tags
	v, err = nextVersion(v, incr)
	return v, incr, err
}

// TaggedVersionSelector returns the greatest version of module tagged in the git repository.
func TaggedVersionSelector(r *Repository, module string) (string, VersionIncrement, error) {
	v, err := taggedVersion(r.RootPath, module, false)
	return v, NoBump, err
}

// DevelopmentVersionSelector returns a commit hash based version if the module has associated pending Changes, otherwise
// returns the latest version from the repo's metadata version enclosure.
func DevelopmentVersionSelector(r *Repository, module string) (string, VersionIncrement, error) {
	incr := versionIncrement(r.Metadata.GetChanges(module))

	if incr != NoBump {
		v, err := pseudoVersion(r.RootPath, module)
		return v, incr, err
	}

	if v, ok := r.Metadata.CurrentVersions.ModuleVersions[module]; ok {
		return v.Version, incr, nil
	}

	return "", NoBump, fmt.Errorf("couldn't select version for module %s: module has no changes and no versions.json version", module)
}
