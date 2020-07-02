package changes

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var tmpDir string

func TestMain(m *testing.M) {
	var err error

	// testdata already has .changes, but make a temporary .changes directory so that we can easily cleanup after
	// tests have run.
	tmpDir, err = ioutil.TempDir("", "changes-test")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(tmpDir, ".changes", "next-release"), 0755)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(filepath.Join(tmpDir, ".changes", "releases"), 0755)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	err = os.RemoveAll(tmpDir)
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestLoadMetadata(t *testing.T) {
	m, err := LoadMetadata("testdata/.changes")
	if err != nil {
		t.Fatal(err)
	}

	if len(m.Changes) != 1 {
		t.Errorf("expected Metadata to have 1 change, got %d", len(m.Changes))
	}

	_, err = LoadMetadata("testdata/.changes-invalid")
	if err == nil {
		t.Fatalf("expected non-nil err, got 'nil'")
	}
}

func TestMetadata_AddChange(t *testing.T) {
	const changeID = "test-change-123456"
	m := getMetadata(t)

	newChange := &Change{
		ID:          changeID,
		Module:      "test/module",
		Type:        FeatureChangeType,
		Description: "test description",
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

	c, err := m2.GetChangeById("test-change-123456")
	if err != nil {
		t.Fatal(err)
	}

	assertChangeEqual(t, newChange, c)

	if c.SchemaVersion != SchemaVersion {
		t.Errorf("Expected SchemaVersion %d, got %d", SchemaVersion, c.SchemaVersion)
	}

	err = m.RemoveChangeById("test-change-123456")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMetadata_AddChanges(t *testing.T) {
	m := getMetadata(t)

	changes := []*Change{
		&Change{
			ID: "test-change-1",
		},
		&Change{
			ID: "test-change-2",
		},
	}

	err := m.AddChanges(changes)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range changes {
		_, err := m.GetChangeById(c.ID)
		if err != nil {
			t.Error(err)
		}

		err = m.RemoveChangeById(c.ID)
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

	_, err := m.GetChangeById("invalid-id")
	if err == nil {
		t.Errorf("Expected non-nil err, got nil")
	}

	c, err := m.GetChangeById("test-feature-1")
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

	assertChangeEqual(t, m.Changes[3], changes[0])
}

func TestMetadata_AddChangesFromTemplate(t *testing.T) {
	testCases := getTestTemplateCases(t)
	m := getMetadata(t)

	for id, tt := range testCases {
		t.Run(id, func(t *testing.T) {
			changes, err := m.AddChangesFromTemplate(tt.template)
			if len(changes) == 0 {
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
					assertChangesHas(t, changes, c)
					assertChangesHas(t, m2.Changes, c)
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
			change := &Change{
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

			newChange, err := m.UpdateChangeFromTemplate(change, tt.template)
			if len(tt.changes) != 1 {
				// this is an invalid template
				if err == nil {
					t.Errorf("expected non-nil err, got nil")
				}
			} else {
				m2 := getMetadata(t)

				// this is a valid template
				if err != nil {
					t.Errorf("expected nil err, got %v", err)
				}

				if _, err := m.GetChangeById("test-change-1"); err == nil {
					t.Errorf("old change was not removed")
				}

				if _, err := m2.GetChangeById("test-change-1"); err == nil {
					t.Errorf("old change was not removed")
				}

				assertChangeEqual(t, tt.changes[0], newChange)
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

	files, err := ioutil.ReadDir(filepath.Join(tmpDir, ".changes", "next-release"))
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

	err = os.Chdir("testdata/.changes/next-release")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalWd)

	path, err := GetChangesPath()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(path, "testdata/.changes") {
		t.Errorf("expected suffix of path to be testdata/.changes, but path is %s", path)
	}
}

func getMetadata(t *testing.T) *Metadata {
	t.Helper()

	m, _ := LoadMetadata(filepath.Join(tmpDir, ".changes"))
	return m
}

func getMockChanges(t *testing.T) []*Change {
	t.Helper()

	return []*Change{
		&Change{
			ID:          "test-feature-1",
			Module:      "test",
			Type:        FeatureChangeType,
			Description: "test description",
		},
		&Change{
			ID:          "test-bugfix-2",
			Module:      "test",
			Type:        BugFixChangeType,
			Description: "test description",
		},
		&Change{
			ID:          "test-feature-3",
			Module:      "test",
			Type:        FeatureChangeType,
			Description: "test description",
		},
		&Change{
			ID:          "other-feature-4",
			Module:      "other",
			Type:        FeatureChangeType,
			Description: "test description",
		},
	}
}

type templateCase struct {
	template []byte
	changes  []*Change
}

func getTestTemplateCases(t *testing.T) map[string]templateCase {
	t.Helper()

	templates := map[string]templateCase{}

	files, err := ioutil.ReadDir("testdata/templates")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		template, err := ioutil.ReadFile(filepath.Join("testdata", "templates", f.Name()))
		if err != nil {
			t.Fatal(err)
		}

		templates[strings.TrimSuffix(f.Name(), ".yaml")] = templateCase{
			template: template,
		}
	}

	changes, err := loadChanges("testdata/changes")
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range changes {
		if tempCase, ok := templates[c.ID]; ok {
			tempCase.changes = append(tempCase.changes, c)
			templates[c.ID] = tempCase
			c.ID = "" // we don't want to compare newly created IDs to this ID since they'll definitely differ
		}
	}

	return templates
}
