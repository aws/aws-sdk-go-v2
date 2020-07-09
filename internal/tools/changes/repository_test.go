package changes

import (
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo, err := NewRepository("testdata")
	if err != nil {
		t.Error(err)
	}

	if repo.Metadata.ChangePath != filepath.Join("testdata", metadataDir) {
		t.Errorf("expected Metadata.ChangePath to be %s, got %s", filepath.Join("testdata", metadataDir), repo.Metadata.ChangePath)
	}
}

func TestRepository_UpdateChangelog(t *testing.T) {
	repo := getRepository(t)

	for _, changelogType := range []string{"pending", "regular"} {
		cases := getTestChangelogCases(t, changelogType)
		for id, tt := range cases {
			pending := false
			fileName := filepath.Join("testdata", "CHANGELOG.md")
			if changelogType == "pending" {
				pending = true
				fileName = filepath.Join("testdata", "CHANGELOG_PENDING.md")
			}

			t.Run(id+changelogType, func(t *testing.T) {
				err := repo.UpdateChangelog(tt.release, pending)
				if err != nil {
					t.Fatal(err)
				}

				changelog, err := ioutil.ReadFile(fileName)
				if err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(string(changelog), tt.changelog); diff != "" {
					t.Errorf("expect changelogs to match:\n%v", diff)
				}

				err = os.Remove(fileName)
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func TestRepository_discoverVersions(t *testing.T) {
	testVersionSelector := func(r *Repository, mod string) (string, error) {
		switch mod {
		case "a":
			return "v1.0.0", nil
		case "b":
			return "v1.2.3", nil
		case "c/v2":
			return "v2.0.0", nil
		default:
			return "", errors.New("couldn't get version")
		}
	}

	var cases = map[string]struct {
		modules       []string
		selector      VersionSelector
		wantEnclosure VersionEnclosure
		wantErr       string
	}{
		"two modules": {
			modules:  []string{"a", "b"},
			selector: testVersionSelector,
			wantEnclosure: VersionEnclosure{
				SchemaVersion: SchemaVersion,
				ModuleVersions: map[string]Version{
					"a": {"a", "v1.0.0"},
					"b": {"b", "v1.2.3"},
				},
			},
		},
		"three modules": {
			modules:  []string{"a", "b", "c/v2"},
			selector: testVersionSelector,
			wantEnclosure: VersionEnclosure{
				SchemaVersion: SchemaVersion,
				ModuleVersions: map[string]Version{
					"a":    {"a", "v1.0.0"},
					"b":    {"b", "v1.2.3"},
					"c/v2": {"c/v2", "v2.0.0"},
				},
			},
		},
		"error": {
			modules:  []string{"a", "b", "error"},
			selector: testVersionSelector,
			wantErr:  "couldn't get version",
		},
	}

	for id, tt := range cases {
		repo := getRepository(t)

		t.Run(id, func(t *testing.T) {
			enc, err := repo.discoverVersions(tt.modules, tt.selector)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatal(err)
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected err to contain %s, got %v", tt.wantErr, err)
				}
			}

			if diff := cmp.Diff(enc, tt.wantEnclosure); diff != "" {
				t.Errorf("expect enclosures to match:\n%v", diff)
			}
		})
	}
}

func TestRepository_DiscoverVersions(t *testing.T) {
	t.Run("no changes", func(t *testing.T) {
		repo := getRepository(t)

		enc, err := repo.DiscoverVersions(ReleaseVersionSelector)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(enc, repo.Metadata.CurrentVersions); diff != "" {
			t.Errorf("expect enclosures to match:\n%v", diff)
		}
	})

	t.Run("module a minor bump", func(t *testing.T) {
		repo := getRepository(t)

		repo.Metadata.Changes = []Change{
			{
				ID:            "test-change",
				SchemaVersion: SchemaVersion,
				Module:        "internal/tools/changes/testdata/modules/a",
				Type:          FeatureChangeType,
				Description:   "this is a test change",
			},
		}

		enc, err := repo.DiscoverVersions(ReleaseVersionSelector)
		if err != nil {
			t.Fatal(err)
		}

		wantEnc := repo.Metadata.CurrentVersions
		wantEnc.ModuleVersions["internal/tools/changes/testdata/modules/a"] = Version{"internal/tools/changes/testdata/modules/a", "v0.1.0"}

		if diff := cmp.Diff(enc, repo.Metadata.CurrentVersions); diff != "" {
			t.Errorf("expect enclosures to match:\n%v", diff)
		}
	})

	t.Run("new module", func(t *testing.T) {
		repo := getRepository(t)
		// simulate new module by removing "a" from CurrentVersions
		delete(repo.Metadata.CurrentVersions.ModuleVersions, "internal/tools/changes/testdata/modules/a")

		enc, err := repo.DiscoverVersions(ReleaseVersionSelector)
		if err != nil {
			t.Fatal(err)
		}

		wantEnc := repo.Metadata.CurrentVersions
		wantEnc.ModuleVersions["internal/tools/changes/testdata/modules/a"] = Version{"internal/tools/changes/testdata/modules/a", "v0.0.0"}

		if diff := cmp.Diff(enc, repo.Metadata.CurrentVersions); diff != "" {
			t.Errorf("expect enclosures to match:\n%v", diff)
		}
	})
}

func getRepository(t *testing.T) *Repository {
	t.Helper()

	repo, err := NewRepository("testdata")
	if err != nil {
		panic(err)
	}

	return repo
}

type changelogCase struct {
	release   *Release
	changelog string
}

func getTestChangelogCases(t *testing.T, changelogType string) map[string]changelogCase {
	t.Helper()
	const releasesTestDir = "releases"
	const changelogsTestDir = "changelogs"

	cases := map[string]changelogCase{}

	files, err := ioutil.ReadDir(filepath.Join("testdata", releasesTestDir))
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		releaseData, err := ioutil.ReadFile(filepath.Join("testdata", releasesTestDir, f.Name()))
		if err != nil {
			t.Fatal(err)
		}

		var release Release
		err = json.Unmarshal(releaseData, &release)
		if err != nil {
			t.Fatal(err)
		}

		changelog, err := ioutil.ReadFile(filepath.Join("testdata", changelogsTestDir, changelogType, release.ID+".md"))
		if err != nil {
			t.Fatal(err)
		}

		cases[release.ID] = changelogCase{&release, string(changelog)}
	}

	return cases
}
