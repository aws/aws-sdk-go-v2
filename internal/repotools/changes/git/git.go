package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/util"
)

// Client is a wrapper around the git CLI tool.
type Client struct {
	RepoPath string // RepoPath is the path to the git repository's root. Git commands executed by this Client will be run at this path.
}

// VcsClient is a client that interacts with a version control system.
type VcsClient interface {
	Tag(tag, message string) error
	Tags(prefix string) ([]string, error)
	Commit(unstagedPaths []string, message string) error
	Push() error
	CommitHash() (string, error)
}

// Tag creates an annotated git tag with the given message.
func (c Client) Tag(tag, message string) error {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	_, err := util.ExecAt(cmd, c.RepoPath)
	if err != nil {
		return err
	}

	return nil
}

// Tags returns all git tags with the given prefix.
func (c Client) Tags(prefix string) ([]string, error) {
	cmd := exec.Command("git", "tag", "-l", prefix+"*")
	output, err := util.ExecAt(cmd, c.RepoPath)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(output), "\n"), nil
}

// Commit stages all given paths and commits with the given message.
func (c Client) Commit(unstagedPaths []string, message string) error {
	for _, p := range unstagedPaths {
		cmd := exec.Command("git", "add", p)
		_, err := util.ExecAt(cmd, c.RepoPath)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("git", "commit", "-m", message)
	_, err := util.ExecAt(cmd, c.RepoPath)

	return err
}

// Push pushes commits and tags to the repository.
func (c Client) Push() error {
	cmd := exec.Command("git", "push", "--follow-tags")
	_, err := util.ExecAt(cmd, c.RepoPath)
	return err
}

// CommitHash returns a timestamp and commit hash for the HEAD commit of the given repository, formatted in the way
// expected for a go.mod file pseudo-version.
func (c Client) CommitHash() (string, error) {
	cmd := exec.Command("git", "show", "--quiet", "--abbrev=12", "--date=format-local:%Y%m%d%H%M%S", "--format=%cd-%h")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "TZ=UTC")

	output, err := util.ExecAt(cmd, c.RepoPath)
	if err != nil {
		return "", fmt.Errorf("couldn't make pseudo-version: %v", err)
	}

	return strings.Trim(string(output), "\n"), nil // clean up git show output and return
}
