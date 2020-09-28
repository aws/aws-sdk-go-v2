package changes

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/golist"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/util"
)

// Repository is a representation of a git repository containing multiple Go modules.
type Repository struct {
	RootPath string    // RootPath is the path to the root of the repository.
	Metadata *Metadata // Metadata is the repository's .changes metadata.

	// Logf is a logging function. If nil, Repository will not log anything.
	Logf func(string, ...interface{})

	modules []string // modules contains the shortened module path of all modules in the repository. modules is lazily loaded by Modules().

	git    git.VcsClient
	golist golist.ModuleClient
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
		git: git.Client{
			RepoPath: path,
		},
		golist: golist.Client{
			RootPath:        path,
			ShortenModPath:  shortenModPath,
			LengthenModPath: lengthenModPath,
		},
		Logf: log.Printf,
	}

	return repo, nil
}

// Modules returns all Go modules under the given Repository. The list of module names returns are shortened to be relative
// to the root of the repository.
func (r *Repository) Modules() ([]string, error) {
	if r.modules != nil {
		return r.modules, nil
	}

	mods, err := discoverModules(r.RootPath)
	if err != nil {
		return nil, err
	}

	r.modules = mods
	return mods, err
}

// DoRelease runs the automated release process, consuming the given Repository's Metadata, updating module's go.mod files,
// creating a release JSON file, committing changes, tagging the repository, and pushing.
func (r *Repository) DoRelease(releaseID string, push, interactive bool) error {
	if err := r.Metadata.CurrentVersions.isValid(r.git); err != nil {
		return fmt.Errorf("couldn't create a release: %v", err)
	}

	enc, bumps, err := r.DiscoverVersions(ReleaseVersionSelector)
	if err != nil {
		return fmt.Errorf("couldn't discover versions for release: %v", err)
	}

	err = r.resolveDependencies(&enc, bumps) // update all dependencies on SDK modules to their latest versions.
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %v", err)
	}

	rel, err := r.Metadata.CreateRelease(releaseID, bumps, false)
	if err != nil {
		return fmt.Errorf("failed to create release metadata: %v", err)
	}

	err = r.UpdateAllChangelogs(rel, false)
	if err != nil {
		return fmt.Errorf("failed to update changelogs: %v", err)
	}

	err = r.updateVersionFiles(enc)
	if err != nil {
		return fmt.Errorf("failed to update version files: %v", err)
	}

	if err = r.commit(releaseID); err != nil {
		return fmt.Errorf("failed to commit changes to local repo: %w", err)
	}

	if interactive {
		if err = confirmationPrompt("Enter \"yes\" to proceed with tagging"); err != nil {
			return err
		}
	}

	if err = r.tag(releaseID, bumps); err != nil {
		return fmt.Errorf("failed to tag and push to repo: %v", err)
	}

	if push {
		if err = r.git.Push(); err != nil {
			return err
		}
	}

	return nil
}

func confirmationPrompt(msg string) error {
	for {
		fmt.Printf(msg + ": ")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			return err
		}
		if strings.EqualFold(response, "yes") {
			return nil
		} else if strings.EqualFold(response, "no") {
			return fmt.Errorf("stopping execution")
		}
	}
}

func (r *Repository) commit(releaseID string) error {
	return r.git.Commit([]string{"."}, fmt.Sprintf("Release %s", releaseID))
}

func (r *Repository) tag(releaseID string, bumps map[string]VersionBump) error {
	for mod, bump := range bumps {
		r.logf("tagging module %s: %s -> %s\n", mod, bump.From, bump.To)
		err := tagRepo(r.git, releaseID, mod, bump.To)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) resolveDependencies(enc *VersionEnclosure, bumps map[string]VersionBump) error {
	updates, err := r.updateEnclosure(enc, bumps)
	if err != nil {
		return err
	}

	for m, deps := range updates {
		err = updateDependencies(r.RootPath, m, deps, enc)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateEnclosure finds all necessary dependency updates, updating the given VersionEnclosure and VersionBumps. updateEnclosure
// returns a mapping between modules and a list of that module's dependencies that must be updated to the latest version.
// updateEnclosure also creates change metadata for the dependency updates.
func (r *Repository) updateEnclosure(enc *VersionEnclosure, bumps map[string]VersionBump) (map[string][]string, error) {
	mods, err := r.Modules()
	if err != nil {
		return nil, fmt.Errorf("couldn't resolve dependencies: %v", err)
	}

	depGraph, err := moduleGraph(r.golist, mods)
	if err != nil {
		return nil, err
	}

	updates := depGraph.dependencyUpdates(AffectedModules(r.Metadata.Changes))
	if len(updates) == 0 {
		return nil, nil
	}

	var depMods []string
	for m := range updates {
		if _, ok := bumps[m]; !ok {
			depMods = append(depMods, m)

			bumps[m], err = enc.bump(m, PatchBump)
			if err != nil {
				return nil, err
			}
		}
	}

	err = r.Metadata.addDependencyUpdateChange(depMods)
	if err != nil {
		return nil, err
	}

	return updates, nil
}

// UpdateAllChangelogs generates changelog entries for both the top level CHANGELOG.md and per-module CHANGELOG.md files.
func (r *Repository) UpdateAllChangelogs(release *Release, pending bool) error {
	err := r.UpdateChangelog(release, pending)
	if err != nil {
		return err
	}

	for _, m := range release.AffectedModules() {
		if m == rootModule {
			// skip updating the root module, since this would overwrite the top level consolidated CHANGELOG
			continue
		}

		err = r.UpdateModuleChangelog(release, m, pending)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateChangelog generates a new CHANGELOG entry for the given release. If pending is true, the contents of
// CHANGELOG_PENDING.md will be replaced with the new entry. Otherwise, the entry will be prepended to CHANGELOG.md.
func (r *Repository) UpdateChangelog(release *Release, pending bool) error {
	entry, err := release.RenderChangelog()
	if err != nil {
		return err
	}

	return util.WriteFile([]byte(entry), filepath.Join(r.RootPath, "CHANGELOG.md"), !pending)
}

// UpdateModuleChangelog generates a changelog entry for the specified module, updating the module's CHANGELOG.md file.
func (r *Repository) UpdateModuleChangelog(release *Release, module string, pending bool) error {
	entry, err := release.RenderChangelogForModule(module, false)
	if err != nil {
		return err
	}

	return util.WriteFile([]byte(entry), filepath.Join(r.RootPath, module, "CHANGELOG.md"), !pending)
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

// DiscoverVersions creates a VersionEnclosure containing all Go modules in the Repository. The version of each module
// is determined by the provided VersionSelector.
func (r *Repository) DiscoverVersions(selector VersionSelector) (VersionEnclosure, map[string]VersionBump, error) {
	modules, err := discoverModules(r.RootPath)
	if err != nil {
		return VersionEnclosure{}, nil, fmt.Errorf("failed to discover versions: %v", err)
	}

	packages, err := packages(r.golist, modules)
	if err != nil {
		return VersionEnclosure{}, nil, fmt.Errorf("failed to discover versions: %v", err)
	}

	enc, bumps, err := r.discoverVersions(modules, selector)
	if err != nil {
		return VersionEnclosure{}, nil, fmt.Errorf("failed to discover versions: %v", err)
	}

	hashes, err := r.ModuleHashes(enc)
	if err != nil {
		return VersionEnclosure{}, nil, fmt.Errorf("failed to discover versions: %v", err)
	}

	err = enc.updateHashes(hashes)
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
			return VersionEnclosure{}, nil, fmt.Errorf("failed to discover version of module %s: %v", m, err)
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

// ModuleHashes computes and returns a mapping between shortened module names and their corresponding Go module checksum.
// Since the version of a module is used to compute the Go checksum (i.e. two modules with the same directory contents
// have a different checksum if their versions are different), the versions of the given VersionEnclosure will be used
// for this purpose.
func (r *Repository) ModuleHashes(enc VersionEnclosure) (map[string]string, error) {
	modules, err := r.Modules()
	if err != nil {
		return nil, fmt.Errorf("couldn't compute module hashes: %v", err)
	}

	hashes := map[string]string{}

	for _, m := range modules {
		hash, err := r.golist.Checksum(m, enc.ModuleVersions[m].Version)
		if err != nil {
			return nil, fmt.Errorf("couldn't compute module hashes: %v", err)
		}

		hashes[m] = hash
	}

	return hashes, nil
}

// updateVersionFiles updates versions.json and all existing version.go files using the given VersionEnclosure.
func (r *Repository) updateVersionFiles(enc VersionEnclosure) error {
	err := r.UpdateVersionFiles(enc)
	if err != nil {
		return err
	}

	hashes, err := r.ModuleHashes(enc)
	if err != nil {
		return err
	}

	// update hashes after modifying module contents (e.g. version.go, CHANGELOG.md, and go.mod files)
	err = enc.updateHashes(hashes)
	if err != nil {
		return err
	}

	return r.Metadata.SaveEnclosure(enc)
}

// Tidy runs go mod tidy on all modules in the repository.
func (r *Repository) Tidy() error {
	modules, err := r.Modules()
	if err != nil {
		return fmt.Errorf("couldn't tidy modules: %v", err)
	}

	for _, m := range modules {
		err = r.golist.Tidy(m)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateVersionFiles updates the version.go file with the version present in the VersionEnclosure enc. If version.go
// does not exist for a module, UpdateVersionFiles does not create one.
func (r *Repository) UpdateVersionFiles(enc VersionEnclosure) error {
	const versionGoFile = "version.go"
	const prefix = "const ModuleVersion = "
	const template = `const ModuleVersion = "%s"`

	modules, err := r.Modules()
	if err != nil {
		return fmt.Errorf("couldn't update version.go files: %v", err)
	}

	for _, m := range modules {
		v := enc.ModuleVersions[m].Version
		if v == "" {
			return fmt.Errorf("couldn't update version.go file for module %s: version is empty", m)
		}

		path := filepath.Join(r.RootPath, modToPath(m), versionGoFile)

		exists, err := util.FileExists(path, false)
		if err != nil {
			return fmt.Errorf("couldn't determine if %s exists for module %s: %v", versionGoFile, m, err)
		}

		if exists {
			err = util.ReplaceLine(path, prefix, fmt.Sprintf(template, v))
			if err != nil {
				return fmt.Errorf("couldn't update %s for module %s: %v", versionGoFile, m, err)
			}
		}
	}

	return nil
}

func (r *Repository) logf(format string, a ...interface{}) {
	if r.Logf != nil {
		r.Logf(format, a...)
	}
}
