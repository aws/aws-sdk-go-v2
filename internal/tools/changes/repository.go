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
		entry += release.RenderChangelogForModule(module, "#")
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
