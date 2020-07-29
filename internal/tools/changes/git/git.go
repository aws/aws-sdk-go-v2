package git

import (
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes/util"
	"os/exec"
	"strings"
)

type Client struct {
	RepoPath string
}

type VcsClient interface {
	Tag(tag, message string) error
	Tags(prefix string) ([]string, error)
	Commit(unstagedPaths []string) error
	Push() error
}

func (c Client) Tag(tag, message string) error {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	_, err := util.ExecAt(cmd, c.RepoPath)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) Tags(prefix string) ([]string, error) {
	cmd := exec.Command("git", "tag", "-l", prefix+"*")
	output, err := util.ExecAt(cmd, c.RepoPath)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(output), "\n"), nil
}

func (c Client) Commit(unstagedPaths []string) error {
	for _, p := range unstagedPaths {
		cmd := exec.Command("git", "add", p)
		_, err := util.ExecAt(cmd, c.RepoPath)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("git", "commit", "-m", `"release commit"`)
	_, err := util.ExecAt(cmd, c.RepoPath)

	return err
}

func (c Client) Push() error {
	cmd := exec.Command("git", "push", "--follow-tags")
	_, err := util.ExecAt(cmd, c.RepoPath)
	return err
}
