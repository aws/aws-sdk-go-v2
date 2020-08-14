package git

import (
	"strings"
	"testing"
)

func TestCommitHash(t *testing.T) {
	client := Client{
		RepoPath: ".",
	}

	commitHash, err := client.CommitHash()
	if err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(commitHash, "-")
	if len(parts) != 2 {
		t.Errorf("expected commit hash to have 2 parts separated by '-', got %s", commitHash)
	}

	if len(parts[1]) != 12 {
		t.Errorf("expected commit hash length to be 12, got %d", len(parts[1]))
	}
}
