package changes

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
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
		for _, tt := range cases {
			pending := false
			fileName := filepath.Join("testdata", "CHANGELOG.md")
			if changelogType == "pending" {
				pending = true
				fileName = filepath.Join("testdata", "CHANGELOG_PENDING.md")
			}

			t.Run(tt.release.ID+"_"+changelogType, func(t *testing.T) {
				err := repo.UpdateChangelog(tt.release, pending)
				if err != nil {
					t.Fatal(err)
				}

				changelog, err := ioutil.ReadFile(fileName)
				if err != nil {
					t.Fatal(err)
				}

				if string(changelog) != tt.changelog {
					t.Errorf("expected changelog \"%s\", got \"%s\"", tt.changelog, string(changelog))
				}

				err = os.Remove(fileName)
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func getRepository(t *testing.T) *Repository {
	return &Repository{
		"testdata",
		getMetadata(t),
	}
}

type changelogCase struct {
	release   *Release
	changelog string
}

func getTestChangelogCases(t *testing.T, changelogType string) []changelogCase {
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

		changelog, err := ioutil.ReadFile(filepath.Join("testdata", changelogsTestDir, changelogType, release.ID+".md"))
		if err != nil {
			t.Fatal(err)
		}

		cases = append(cases, changelogCase{&release, string(changelog)})
	}

	return cases
}
