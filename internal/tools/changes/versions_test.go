package changes

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestVersionEnclosure_IsValid(t *testing.T) {
	repo := getRepository(t)

	t.Run("success - empty enclosure", func(t *testing.T) {
		enc, _, err := repo.DiscoverVersions(TaggedVersionSelector)
		if err != nil {
			t.Fatal(err)
		}

		err = enc.isValid(repo.RootPath)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("failure - non-existent tag", func(t *testing.T) {
		enc := VersionEnclosure{
			SchemaVersion: SchemaVersion,
			ModuleVersions: map[string]Version{
				"module": {"module", "module", "v123.123.123"},
			},
		}

		err := enc.isValid(repo.RootPath)
		if err == nil {
			t.Error("expected non-nil err, got nil")
		}
	})
}

func TestVersionEnclosure_Bump(t *testing.T) {
	enc := VersionEnclosure{
		SchemaVersion: SchemaVersion,
		ModuleVersions: map[string]Version{
			"a": {
				Module:     "a",
				ImportPath: "a",
				Version:    "v1.0.0",
			},
			"b": {
				Module:     "b",
				ImportPath: "b",
				Version:    "v2.1.3",
			},
		},
		Packages: map[string]string{},
	}

	wantBump := VersionBump{
		From: "v1.0.0",
		To:   "v1.0.1",
	}

	bump, err := enc.bump("a", PatchBump)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(bump, wantBump); diff != "" {
		t.Errorf("expect bumps to match:\n%v", diff)
	}
}
