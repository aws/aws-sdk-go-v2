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
const SchemaVersion = "1.0"

// Metadata is a representation of the change metadata stored in a .changes directory.
type Metadata struct {
	ChangePath string     // ChangePath is the relative path from the current directory to .changes
	Changes    []*Change  // Changes holds all pending change metadata in the .changes/next-release directory
	Releases   []*Release // Releases contains all releases in the .changes/releases directory
}

// LoadMetadata searches from the current directory upwards until it finds a .changes directory, which will be loaded
// into a new Metadata and returned.
func LoadMetadata() (*Metadata, error) {
	path, err := GetChangesPath()
	if err != nil {
		return nil, err
	}

	changes, err := loadChanges(filepath.Join(path, "next-release"))
	if err != nil {
		return nil, err
	}

	return &Metadata{
		ChangePath: path,
		Changes:    changes,
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

// GetChangeById returns the pending Change with the given id.
func (m *Metadata) GetChangeById(id string) (*Change, error) {
	_, c, err := m.getChange(id)
	return c, err
}

func (m *Metadata) getChange(id string) (int, *Change, error) {
	for i, c := range m.Changes {
		if c.Id == id {
			return i, c, nil
		}
	}

	return 0, nil, fmt.Errorf("couldn't find change with id %s", id)
}

// ListChanges returns all pending Changes.
func (m *Metadata) ListChanges() []*Change {
	return m.Changes
}

// SaveChange saves the given change to the .changes/next-release directory.
func (m *Metadata) SaveChange(c *Change) error {
	return writeFile(c, m.ChangePath, "next-release", c.Id)
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
	for _, c := range m.Changes {
		err := m.RemoveChangeById(c.Id)
		if err != nil {
			return err
		}
	}

	m.Changes = []*Change{}
	return nil
}

// CreateRelease consolidates the Metadata's pending Changes into a Release. This operation will remove all Changes from
// the Metadata and delete change files in .changes/next-release. A release file will also be created in .changes/releases.
func (m *Metadata) CreateRelease(id string, bumps []VersionBump) error {
	release := &Release{
		Id:           id,
		VersionBumps: bumps,
		Changes:      m.Changes,
	}

	if err := writeFile(release, m.ChangePath, "releases", id); err != nil {
		return err
	}

	return m.ClearChanges()
}

// deleteChangeFile deletes the file .changes/next-release/{id}.json.
func (m *Metadata) deleteChangeFile(id string) error {
	return os.Remove(filepath.Join(m.ChangePath, "next-release", id+".json"))
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

			changes = append(changes, change)
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return strings.Compare(changes[i].Id, changes[j].Id) < 0
	})

	return changes, nil
}
