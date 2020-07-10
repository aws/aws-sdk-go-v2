package changes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SchemaVersion defines the current JSON schema version for persistent data types (Change, Release, ...)
const SchemaVersion = 1

const metadataDir = ".changes"
const pendingDir = "next-release"
const releaseDir = "releases"
const versionsFile = "versions.json"

// Metadata is a representation of the change metadata stored in a .changes directory.
type Metadata struct {
	ChangePath      string           // ChangePath is the relative path from the current directory to .changes
	Changes         []*Change        // Changes holds all pending change metadata in the .changes/next-release directory
	Releases        []*Release       // Releases contains all releases in the .changes/releases directory
	CurrentVersions VersionEnclosure // CurrentVersions is the .changes/versions.json enclosure of current module versions
}

// LoadMetadata loads the .changes directory at the given path.
func LoadMetadata(path string) (*Metadata, error) {
	changes, err := loadChanges(filepath.Join(path, "next-release"))
	if err != nil {
		return nil, err
	}

	v, err := loadVersions(filepath.Join(path, versionsFile))
	if err != nil {
		return nil, err
	}

	return &Metadata{
		ChangePath:      path,
		Changes:         changes,
		CurrentVersions: v,
	}, nil
}

// AddChange adds the given Change to the Metadata's Changes and saves the Change to the next-release directory.
func (m *Metadata) AddChange(c *Change) error {
	err := m.SaveChange(c)
	if err != nil {
		return err
	}

	m.Changes = append(m.Changes, c)
	return nil
}

// AddChangesFromTemplate parses the given YAML template, adding the resulting Changes to Metadata's Changes and saving
// the Changes to the next-release directory. AddChangesFromTemplate returns the created Changes.
func (m *Metadata) AddChangesFromTemplate(template []byte) ([]*Change, error) {
	changes, err := TemplateToChanges(template)
	if err != nil {
		return nil, err
	}

	return changes, m.AddChanges(changes)
}

// AddChanges adds the given Changes to Metadata's Changes and saves the Changes to the next-release directory.
func (m *Metadata) AddChanges(changes []*Change) error {
	for _, c := range changes {
		err := m.AddChange(c)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetChangeById returns the pending Change with the given id.
func (m *Metadata) GetChangeById(id string) (*Change, error) {
	_, c, err := m.getChange(id)
	return c, err
}

func (m *Metadata) getChange(id string) (int, *Change, error) {
	for i, c := range m.Changes {
		if c.ID == id {
			return i, c, nil
		}
	}

	return 0, nil, fmt.Errorf("couldn't find change with id %s", id)
}

// ListChanges returns all pending Changes with a module matching the given module. If module is empty, returns all Changes.
func (m *Metadata) GetChanges(module string) []*Change {
	if module == "" {
		return m.Changes
	}

	var changes []*Change

	for _, c := range m.Changes {
		if module == c.Module {
			changes = append(changes, c)
		}
	}
	return changes
}

// SaveChange saves the given change to the .changes/next-release directory.
func (m *Metadata) SaveChange(c *Change) error {
	c.SchemaVersion = SchemaVersion
	return writeJSON(c, m.ChangePath, "next-release", c.ID)
}

// UpdateChangeFromTemplate removes oldChange and creates a new Change from the given template.
func (m *Metadata) UpdateChangeFromTemplate(oldChange *Change, template []byte) (*Change, error) {
	newChanges, err := TemplateToChanges(template)
	if err != nil {
		return nil, fmt.Errorf("failed to modify change: %v\n", err)
	} else if len(newChanges) != 1 {
		return nil, fmt.Errorf("failed to modify change: modules cannot be added to a change during modification")
	}

	err = m.SaveChange(newChanges[0])
	if err != nil {
		return nil, fmt.Errorf("failed to modify change: %v\n", err)
	}

	err = m.RemoveChangeById(oldChange.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove old change with id %s: %v\n", oldChange.ID, err)
	}

	return newChanges[0], nil
}

// RemoveChangeById removes the Change with the specified id from the Metadata's Changes and also removes the Change
// from the .changes/next-release directory.
func (m *Metadata) RemoveChangeById(id string) error {
	i, _, err := m.getChange(id)
	if err != nil {
		return fmt.Errorf("failed to remove change: %v", err)
	}

	err = m.deleteChangeFile(id)
	if err != nil {
		return fmt.Errorf("failed to remove changes: %v", err)
	}

	m.Changes = append(m.Changes[:i], m.Changes[i+1:]...)
	return nil
}

// ClearChanges removes all Changes from the Metadata's Changes and deletes the Change files in the
// .changes/next-release directory.
func (m *Metadata) ClearChanges() error {
	var ids []string
	for _, c := range m.Changes {
		ids = append(ids, c.ID)
	}

	for _, id := range ids {
		err := m.RemoveChangeById(id)
		if err != nil {
			return err
		}
	}

	m.Changes = []*Change{}
	return nil
}

// CreateRelease consolidates the Metadata's pending Changes into a Release. This operation will remove all Changes from
// the Metadata and delete change files in .changes/next-release. A release file will also be created in
// .changes/releases. If dryRun is true, CreateRelease will return a Release, but not modify change or release files.
func (m *Metadata) CreateRelease(id string, bumps map[string]VersionBump, dryRun bool) (*Release, error) {
	release := &Release{
		ID:            id,
		SchemaVersion: SchemaVersion,
		VersionBumps:  bumps,
		Changes:       m.Changes,
	}

	if !dryRun {
		if err := writeJSON(release, m.ChangePath, "releases", id); err != nil {
			return nil, err
		}

		return release, m.ClearChanges()
	}

	return release, nil
}

// deleteChangeFile deletes the file .changes/next-release/{id}.json.
func (m *Metadata) deleteChangeFile(id string) error {
	return os.Remove(filepath.Join(m.ChangePath, "next-release", id+".json"))
}

// SaveEnclosure updates the Metadata's enclosure and updates the versions.json file.
func (m *Metadata) SaveEnclosure(enc VersionEnclosure) error {
	err := writeJSON(enc, m.ChangePath, "", "versions")
	if err != nil {
		return err
	}

	m.CurrentVersions = enc
	return nil
}

// GetChangesPath searches upward from the current directory for a .changes directory, returning a relative path from
// the current directory to the .changes directory.
func GetChangesPath() (string, error) {
	return findFile(".changes", true)
}

// loadChanges unmarshals and returns all Changes in the given directory.
func loadChanges(changesDir string) ([]*Change, error) {
	files, err := ioutil.ReadDir(changesDir)
	if err != nil {
		return nil, err
	}

	var changes []*Change

	for _, f := range files {
		if !f.IsDir() {
			changeData, err := ioutil.ReadFile(filepath.Join(changesDir, f.Name()))
			if err != nil {
				return nil, err
			}

			change := &Change{}
			err = json.Unmarshal(changeData, change)
			if err != nil {
				return nil, err
			}

			if change.SchemaVersion != SchemaVersion {
				return nil, fmt.Errorf("change with id %s has Schema Version %s, but verison %s was expeced",
					change.ID, change.SchemaVersion, SchemaVersion)
			}

			changes = append(changes, change)
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return strings.Compare(changes[i].ID, changes[j].ID) < 0
	})

	return changes, nil
}

func loadVersions(path string) (VersionEnclosure, error) {
	versionsData, err := ioutil.ReadFile(path)
	if err != nil {
		return VersionEnclosure{}, fmt.Errorf("couldn't load version enclosure at %s: %v", path, err)
	}

	var enclosure VersionEnclosure
	err = json.Unmarshal(versionsData, &enclosure)
	if err != nil {
		return VersionEnclosure{}, fmt.Errorf("couldn't load version enclosure at %s: %v", path, err)
	}

	return enclosure, nil
}
