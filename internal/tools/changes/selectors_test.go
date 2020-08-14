package changes

import "testing"

func TestDevelopmentVersionSelector(t *testing.T) {
	repo := getRepository(t)

	_, _, err := DevelopmentVersionSelector(repo, modPrefix+"a")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTaggedVersionSelector(t *testing.T) {
	repo := getRepository(t)
	repo.git = &MockGit{
		tags: []string{modPrefix + "a/v1.2.3"},
	}

	version, _, err := TaggedVersionSelector(repo, modPrefix+"a")
	if err != nil {
		t.Fatal(err)
	}

	if version != "v1.2.3" {
		t.Errorf("expected mod a to be at version v1.2.3, got version %s", version)
	}
}
