package changes

import "testing"

func TestDevelopmentVersionSelector(t *testing.T) {
	repo := getRepository(t)

	_, _, err := DevelopmentVersionSelector(repo, modPrefix+"a")
	if err != nil {
		t.Fatal(err)
	}
}
