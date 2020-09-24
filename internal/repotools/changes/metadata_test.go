package changes

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/util"
	"github.com/google/go-cmp/cmp"
)

var tmpDir string // tmpDir is a temporary directory metadata tests use.

func TestMain(m *testing.M) {
	dirName, cleanup, err := setupTmpChanges()
	if err != nil {
		panic(err)
	}

	tmpDir = dirName

	code := m.Run()

	err = cleanup()
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func setupTmpChanges() (name string, cleanup func() error, err error) {
	// testdata already has .changes, but make a temporary .changes directory so that we can easily cleanup after
	// tests have run.
	dirName, err := ioutil.TempDir("", "changes-test")
	if err != nil {
		return "", nil, err
	}

	err = os.MkdirAll(filepath.Join(dirName, metadataDir, pendingDir), 0755)
	if err != nil {
		return "", nil, err
	}

	err = os.Mkdir(filepath.Join(dirName, metadataDir, releaseDir), 0755)
	if err != nil {
		return "", nil, err
	}

	// create empty versions.json
	err = util.WriteJSON(VersionEnclosure{
		SchemaVersion:  SchemaVersion,
		ModuleVersions: map[string]Version{},
		Packages:       map[string]string{},
	}, dirName, metadataDir, "versions")
	if err != nil {
		return "", nil, err
	}

	return dirName, func() error {
		return os.RemoveAll(dirName)
	}, nil
}

func TestLoadMetadata(t *testing.T) {
	m, err := LoadMetadata(filepath.Join("testdata", metadataDir))
	if err != nil {
		t.Fatal(err)
	}

	if len(m.Changes) != 1 {
		t.Errorf("expected Metadata to have 1 change, got %d", len(m.Changes))
	}

	_, err = LoadMetadata(filepath.Join("testdata", ".changes-invalid"))
	if err == nil {
		t.Fatalf("expected non-nil err, got 'nil'")
	}
}

func TestMetadata_AddChange(t *testing.T) {
	const changeID = "test-change-123456"
	m := getMetadata(t)

	newChange := Change{
		ID:            changeID,
		SchemaVersion: SchemaVersion,
		Module:        "test/module",
		Type:          FeatureChangeType,
		Description:   "test description",
	}

	err := m.AddChange(newChange)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, c := range m.Changes {
		if c.ID == changeID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find a change with id %s, but didn't", changeID)
	}

	m2 := getMetadata(t)

	c, err := m2.GetChangeByID("test-change-123456")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(c, newChange); diff != "" {
		t.Errorf("expect changes to match:\n%v", diff)
	}

	if c.SchemaVersion != SchemaVersion {
		t.Errorf("Expected SchemaVersion %d, got %d", SchemaVersion, c.SchemaVersion)
	}

	err = m.RemoveChangeByID("test-change-123456")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMetadata_AddChanges(t *testing.T) {
	m := getMetadata(t)

	changes := []Change{
		{
			ID: "test-change-1",
		},
		{
			ID: "test-change-2",
		},
	}

	err := m.AddChanges(changes)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range changes {
		_, err := m.GetChangeByID(c.ID)
		if err != nil {
			t.Error(err)
		}

		err = m.RemoveChangeByID(c.ID)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestMetadata_CreateRelease(t *testing.T) {
	m := getMetadata(t)
	m.Changes = getMockChanges(t)

	release, err := m.CreateRelease("test-release", map[string]VersionBump{
		"test": {
			From: "v1.0.0",
			To:   "v1.0.1",
		},
		"other": {
			From: "v1.3.2",
			To:   "v1.4.0",
		},
	}, true)
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	if release.SchemaVersion != SchemaVersion {
		t.Errorf("Expected SchmeVersion to be %d, got %d", SchemaVersion, release.SchemaVersion)
	}

	if len(release.Changes) != len(m.Changes) {
		t.Errorf("Expected release to have %d changes, got %d", len(m.Changes), len(release.Changes))
	}
}

func TestMetadata_GetChangeById(t *testing.T) {
	m := getMetadata(t)
	m.Changes = getMockChanges(t)

	_, err := m.GetChangeByID("invalid-id")
	if err == nil {
		t.Errorf("Expected non-nil err, got nil")
	}

	c, err := m.GetChangeByID("test-feature-1")
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	if c.ID != "test-feature-1" {
		t.Errorf("expected ID \"test-feature-1\", got %s", c.ID)
	}
}

func TestMetadata_GetChanges(t *testing.T) {
	m := getMetadata(t)
	m.Changes = getMockChanges(t)

	changes := m.GetChanges("other")
	if len(changes) != 1 {
		t.Fatalf("expected 1 Change, got %d", len(changes))
	}

	if diff := cmp.Diff(m.Changes[3], changes[0]); diff != "" {
		t.Errorf("expect changes to match:\n%v", diff)
	}
}

func TestMetadata_AddChangesFromTemplate(t *testing.T) {
	testCases := getTestTemplateCases(t)
	m := getMetadata(t)

	for id, tt := range testCases {
		t.Run(id, func(t *testing.T) {
			changes, err := m.AddChangesFromTemplate(tt.template)
			if !tt.expectChanges {
				if err == nil {
					t.Errorf("expected non-nil err, got nil")
				}
			} else {
				// test that we can load the new changes
				m2 := getMetadata(t)
				if len(m2.Changes) != len(tt.changes) {
					t.Errorf("expected %d changes, got %d", len(tt.changes), len(m2.Changes))
				}

				for _, c := range tt.changes {
					assertHasChangeLike(t, changes, c)
					assertHasChangeLike(t, m2.Changes, c)
				}

				err = m.ClearChanges()
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestMetadata_UpdateChangeFromTemplate(t *testing.T) {
	testCases := getTestTemplateCases(t)

	for id, tt := range testCases {
		t.Run(id, func(t *testing.T) {
			m := getMetadata(t)
			// add a test change to update
			change := Change{
				ID:            "test-change-1",
				SchemaVersion: 1,
				Module:        "test/module",
				Type:          FeatureChangeType,
				Description:   "test change",
			}

			err := m.AddChange(change)
			if err != nil {
				t.Fatal(err)
			}

			newChanges, err := m.UpdateChangeFromTemplate(change, tt.template)
			if !tt.expectChanges {
				if err == nil {
					t.Errorf("expected non-nil err, got nil")
				}
			} else {
				m2 := getMetadata(t)

				if err != nil {
					t.Errorf("expected nil err, got %v", err)
				}

				if len(m2.Changes) != len(tt.changes) {
					t.Errorf("expected %d changes, got %d", len(tt.changes), len(m2.Changes))
				}

				for _, c := range tt.changes {
					assertHasChangeLike(t, newChanges, c)
					assertHasChangeLike(t, m2.Changes, c)
				}
			}

			err = m.ClearChanges()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMetadata_ClearChanges(t *testing.T) {
	m := getMetadata(t)

	err := m.AddChanges(getMockChanges(t))
	if err != nil {
		t.Fatal(err)
	}

	err = m.ClearChanges()
	if err != nil {
		t.Error(err)
	}

	files, err := ioutil.ReadDir(filepath.Join(tmpDir, metadataDir, pendingDir))
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 0 {
		t.Errorf("expected next-release to contain 0 files, but found %d files", len(files))
	}
}

func TestGetChangesPath(t *testing.T) {
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(filepath.Join("testdata", metadataDir, pendingDir))
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalWd)

	path, err := GetChangesPath()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(path, filepath.Join("testdata", metadataDir)) {
		t.Errorf("expected suffix of path to be %s, but path is %s", filepath.Join("testdata", metadataDir), path)
	}
}

func getMetadata(t *testing.T) *Metadata {
	t.Helper()

	m, err := LoadMetadata(filepath.Join(tmpDir, metadataDir))
	if err != nil {
		panic(err)
	}

	return m
}

func getMockChanges(t *testing.T) []Change {
	t.Helper()

	return []Change{
		{
			ID:          "test-feature-1",
			Module:      "test",
			Type:        FeatureChangeType,
			Description: "test description",
		},
		{
			ID:          "test-bugfix-2",
			Module:      "test",
			Type:        BugFixChangeType,
			Description: "test description",
		},
		{
			ID:          "test-feature-3",
			Module:      "test",
			Type:        FeatureChangeType,
			Description: "test description",
		},
		{
			ID:          "other-feature-4",
			Module:      "other",
			Type:        FeatureChangeType,
			Description: "test description",
		},
	}
}

type templateCase struct {
	template      []byte
	changes       []Change
	expectChanges bool
}

func getTestTemplateCases(t *testing.T) map[string]templateCase {
	t.Helper()
	const templateDir = "templates"
	const changesDir = "changes"

	templates := map[string]templateCase{}

	files, err := ioutil.ReadDir(filepath.Join("testdata", templateDir))
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		template, err := ioutil.ReadFile(filepath.Join("testdata", templateDir, f.Name()))
		if err != nil {
			t.Fatal(err)
		}

		templates[strings.TrimSuffix(f.Name(), ".yaml")] = templateCase{
			template: template,
		}
	}

	changes, err := loadChanges(filepath.Join("testdata", changesDir))
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range changes {
		if tempCase, ok := templates[c.ID]; ok {
			id := c.ID
			tempCase.changes = append(tempCase.changes, c)
			tempCase.expectChanges = true
			templates[id] = tempCase
		}
	}

	return templates
}
