package changes

import (
	"testing"
)

func TestVersionEnclosure_IsValid(t *testing.T) {
	repo := getRepository(t)

	t.Run("success - empty enclosure", func(t *testing.T) {
		enc, err := repo.DiscoverVersions(TaggedVersionSelector)
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
