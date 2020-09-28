package changes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/golist"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const modPrefix = "internal/tools/changes/testdata/modules/"

func TestNewRepository(t *testing.T) {
	repo, err := NewRepository("testdata")
	if err != nil {
		t.Error(err)
	}

	if repo.Metadata.ChangePath != filepath.Join("testdata", metadataDir) {
		t.Errorf("expected Metadata.ChangePath to be %s, got %s", filepath.Join("testdata", metadataDir), repo.Metadata.ChangePath)
	}
}

func TestRepository_Modules(t *testing.T) {
	repo := getRepository(t)

	mods, err := repo.Modules()
	if err != nil {
		t.Fatal(err)
	}

	wantMods := []string{
		"internal/tools/changes/testdata/modules/a",
		"internal/tools/changes/testdata/modules/b",
		"internal/tools/changes/testdata/modules/nested/c/d",
		"internal/tools/changes/testdata/modules/nested/c",
	}

	if diff := cmp.Diff(wantMods, mods); diff != "" {
		t.Errorf("expect modules to match (-want, +got):\n%v", diff)
	}
}

func TestRepository_DoRelease(t *testing.T) {
	repo, cleanup := tmpRepository(t)
	defer cleanup()

	// setup repo with changes
	err := repo.Metadata.AddChanges([]Change{
		{
			ID:            "test-change-1",
			SchemaVersion: SchemaVersion,
			Module:        "a",
			Type:          FeatureChangeType,
			Description:   "a feature change",
		},
		{
			ID:            "test-change-2",
			SchemaVersion: SchemaVersion,
			Module:        "service/...",
			Type:          BugFixChangeType,
			Description:   "all services wildcard bugfix",
			AffectedModules: []string{
				"service/c",
				"service/d",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = repo.DoRelease("test-release", true, false)
	if err != nil {
		t.Fatal(err)
	}

	assertFileContains(t, filepath.Join(repo.RootPath, "a", "go.mod"), []string{"require github.com/aws/aws-sdk-go-v2/service/c v0.1.3"})
	assertFileContains(t, filepath.Join(repo.RootPath, "b", "go.mod"), []string{"require github.com/aws/aws-sdk-go-v2/a v1.1.0"})

	assertFileContains(t, filepath.Join(repo.RootPath, "CHANGELOG.md"), []string{
		"* [a](a/CHANGELOG.md#Release-test-release) - v1.1.0\n  * Feature: a feature change",
		"Service Client Highlights\n* Bug Fix: all services wildcard bugfix",
		"Dependency Update: Updated SDK dependencies to their latest versions.",
	})

	assertFileContains(t, filepath.Join(repo.RootPath, "a", "CHANGELOG.md"), []string{
		"* Feature: a feature change",
	})

	assertFileContains(t, filepath.Join(repo.RootPath, "service", "d", "CHANGELOG.md"), []string{
		"* Bug Fix: all services wildcard bugfix",
		"v1.100.101",
	})
}

func TestRepository_TagAndPush(t *testing.T) {
	repo := getRepository(t)
	gitClient := &MockGit{}
	repo.git = gitClient

	bumps := map[string]VersionBump{
		"a": {
			From: "v1.0.0",
			To:   "v1.0.1",
		},
		"b": {
			From: "v0.0.0",
			To:   "v1.0.0",
		},
		"c/v2": {
			To: "v2.0.0",
		},
	}

	wantTags := []string{"a/v1.0.1", "b/v1.0.0", "c/v2.0.0"}

	err := repo.tag("test-release", bumps)
	if err != nil {
		t.Fatal(err)
	}

	if err = repo.git.Push(); err != nil {
		t.Fatal(err)
	}

	sort.Strings(gitClient.tags)

	if diff := cmp.Diff(wantTags, gitClient.tags); diff != "" {
		t.Errorf("expect tags to match (-want, +got):\n%v", diff)
	}
}

func TestRepository_UpdateChangelog(t *testing.T) {
	repo := getRepository(t)

	dir, err := ioutil.TempDir("", "tmp-changelog-test")
	if err != nil {
		t.Fatal(err)
	}

	repo.RootPath = dir

	cases := getTestChangelogCases(t)
	for _, tt := range cases {
		t.Run(tt.release.ID, func(t *testing.T) {
			err = repo.UpdateChangelog(tt.release, false)
			if err != nil {
				t.Fatal(err)
			}

			changelog, err := ioutil.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.changelog, string(changelog)); diff != "" {
				t.Errorf("expect changelogs to match (-want, +got):\n%v", diff)
			}

			err = os.Remove(filepath.Join(dir, "CHANGELOG.md"))
		})
	}

	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRepository_discoverVersions(t *testing.T) {
	testVersionSelector := func(r *Repository, mod string) (string, VersionIncrement, error) {
		switch mod {
		case "a":
			return "v1.0.0", NewModule, nil
		case "b":
			return "v1.2.3", PatchBump, nil
		case "c/v2":
			return "v2.0.0", NoBump, nil
		default:
			return "", NoBump, errors.New("couldn't get version")
		}
	}

	var cases = map[string]struct {
		modules       []string
		selector      VersionSelector
		wantEnclosure VersionEnclosure
		wantErr       string
		wantBumps     map[string]VersionBump
	}{
		"two modules": {
			modules:  []string{"a", "b"},
			selector: testVersionSelector,
			wantEnclosure: VersionEnclosure{
				SchemaVersion: SchemaVersion,
				ModuleVersions: map[string]Version{
					"a": {"a", sdkRepo + "/" + "a", "v1.0.0", ""},
					"b": {"b", sdkRepo + "/" + "b", "v1.2.3", ""},
				},
			},
			wantBumps: map[string]VersionBump{
				"a": {To: "v1.0.0"},
				"b": {To: "v1.2.3"},
			},
		},
		"three modules": {
			modules:  []string{"a", "b", "c/v2"},
			selector: testVersionSelector,
			wantEnclosure: VersionEnclosure{
				SchemaVersion: SchemaVersion,
				ModuleVersions: map[string]Version{
					"a":    {"a", sdkRepo + "/" + "a", "v1.0.0", ""},
					"b":    {"b", sdkRepo + "/" + "b", "v1.2.3", ""},
					"c/v2": {"c/v2", sdkRepo + "/" + "c/v2", "v2.0.0", ""},
				},
			},
			wantBumps: map[string]VersionBump{
				"a": {To: "v1.0.0"},
				"b": {To: "v1.2.3"},
				// c is NoBump
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
			enc, bumps, err := repo.discoverVersions(tt.modules, tt.selector)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatal(err)
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected err to contain %s, got %v", tt.wantErr, err)
				}
			} else {
				if diff := cmp.Diff(tt.wantEnclosure, enc); diff != "" {
					t.Errorf("expect enclosures to match (-want, +got):\n%v", diff)
				}

				if diff := cmp.Diff(tt.wantBumps, bumps); diff != "" {
					t.Errorf("expect bumps to match (-want, +got):\n%v", diff)
				}
			}
		})
	}
}

func TestRepository_DiscoverVersions(t *testing.T) {
	t.Run("no changes", func(t *testing.T) {
		repo := getRepository(t)

		enc, bumps, err := repo.DiscoverVersions(ReleaseVersionSelector)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(repo.Metadata.CurrentVersions, enc, cmpopts.IgnoreFields(Version{}, "ModuleHash")); diff != "" {
			t.Errorf("expect enclosures to match (-want +got):\n%v", diff)
		}

		if len(bumps) != 0 {
			t.Errorf("expected 0 version bumps, got %d", len(bumps))
		}
	})

	t.Run("module a minor bump", func(t *testing.T) {
		repo := getRepository(t)

		repo.Metadata.Changes = []Change{
			{
				ID:            "test-change",
				SchemaVersion: SchemaVersion,
				Module:        modPrefix + "a",
				Type:          FeatureChangeType,
				Description:   "this is a test change",
			},
		}

		enc, _, err := repo.DiscoverVersions(ReleaseVersionSelector)
		if err != nil {
			t.Fatal(err)
		}

		wantEnc := repo.Metadata.CurrentVersions
		wantEnc.ModuleVersions[modPrefix+"a"] = Version{modPrefix + "a", sdkRepo + "/" + modPrefix + "a", "v0.1.0", ""}

		if diff := cmp.Diff(repo.Metadata.CurrentVersions, enc, cmpopts.IgnoreFields(Version{}, "ModuleHash")); diff != "" {
			t.Errorf("expect enclosures to match (-want +got):\n%v", diff)
		}
	})

	t.Run("new module", func(t *testing.T) {
		repo := getRepository(t)
		// simulate new module by removing "a" from CurrentVersions
		delete(repo.Metadata.CurrentVersions.ModuleVersions, modPrefix+"a")

		enc, _, err := repo.DiscoverVersions(ReleaseVersionSelector)
		if err != nil {
			t.Fatal(err)
		}

		wantEnc := repo.Metadata.CurrentVersions
		wantEnc.ModuleVersions[modPrefix+"a"] = Version{modPrefix + "a", sdkRepo + "/" + modPrefix + "a", "v0.1.0", ""}

		if diff := cmp.Diff(repo.Metadata.CurrentVersions, enc, cmpopts.IgnoreFields(Version{}, "ModuleHash")); diff != "" {
			t.Errorf("expect enclosures to match (-want +got):\n%v", diff)
		}
	})
}

func getRepository(t *testing.T) *Repository {
	t.Helper()

	repo, err := NewRepository("testdata")
	if err != nil {
		panic(err)
	}

	repo.golist = golist.Client{
		RootPath: filepath.Join("testdata", "modules"),
		ShortenModPath: func(mod string) string {
			return strings.TrimPrefix(mod, modPrefix)
		},
		LengthenModPath: func(mod string) string {
			return sdkRepo + "/" + modPrefix + mod
		},
	}

	repo.Logf = nil // suppress logging

	return repo
}

func tmpRepository(t *testing.T) (*Repository, func() error) {
	t.Helper()

	dirName, cleanup, err := setupTmpChanges()
	if err != nil {
		panic(err)
	}

	repo, err := NewRepository(dirName)
	if err != nil {
		panic(err)
	}

	repo.Logf = nil

	repo.git = &MockGit{
		tags: []string{"a/v1.0.0", "b/v1.2.3", "service/c/v0.1.2", "service/d/v1.100.100"},
	}
	repo.golist = &mockGolist{
		dependencies: map[string][]string{
			"a":         {"service/c"},
			"b":         {"a"},
			"service/c": {"a", "b"},
		},
	}

	mods := []string{"a", "b", "newmod", "service/c", "service/d"}
	for _, m := range mods {
		err = addModuleToTmpRepo(t, dirName, m)
		if err != nil {
			panic(err)
		}
	}

	return repo, cleanup
}

func addModuleToTmpRepo(t *testing.T, dir, mod string) error {
	t.Helper()

	const goMod = `module %s

go 1.14
`

	err := os.MkdirAll(filepath.Join(dir, mod), 0755)
	if err != nil {
		return err
	}

	data := fmt.Sprintf(goMod, mod)

	return ioutil.WriteFile(filepath.Join(dir, mod, "go.mod"), []byte(data), 0755)
}

func assertFileContains(t *testing.T, path string, substrings []string) bool {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range substrings {
		if !strings.Contains(string(data), s) {
			t.Errorf("expected file %s to contain %s", path, s)
			return false
		}
	}

	return true
}

type changelogCase struct {
	release   *Release
	changelog string
}

func getTestChangelogCases(t *testing.T) []changelogCase {
	t.Helper()
	const releasesTestDir = "releases"
	const changelogsTestDir = "changelogs"

	var cases []changelogCase

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

		changelog, err := ioutil.ReadFile(filepath.Join("testdata", changelogsTestDir, release.ID+".md"))
		if err != nil {
			t.Fatal(err)
		}

		cases = append(cases, changelogCase{&release, string(changelog)})
	}

	return cases
}
