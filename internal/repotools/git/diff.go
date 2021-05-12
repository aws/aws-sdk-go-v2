package git

import (
	"strings"
)

// Changes outputs a list of files for the repository that have changed at the given paths between from and to commit-like references.
// If no paths are provided then all files that changed between the two tree-ish (commit or tag) references are returned.
func Changes(repository string, from, to string, paths ...string) ([]string, error) {
	arguments := []string{"diff-tree", "--name-only", "--no-commit-id", "-r", from, to}

	if len(paths) > 0 {
		arguments = append(arguments, paths...)
	}

	output, err := Git(repository, arguments...)
	if err != nil {
		return nil, err
	}
	return splitOutput(string(output)), nil
}

// Changed reports the list of files that changed for a specific commit or tag.
func Changed(repository string, commit string, paths ...string) ([]string, error) {
	arguments := []string{"diff-tree", "--name-only", "--no-commit-id", "-r", commit}

	if len(paths) > 0 {
		arguments = append(arguments, paths...)
	}

	output, err := Git(repository, arguments...)
	if err != nil {
		return nil, err
	}
	return splitOutput(string(output)), nil
}

func splitOutput(output string) []string {
	split := strings.Split(output, "\n")
	if len(split) > 0 {
		if len(split[len(split)-1]) == 0 {
			split = split[:len(split)-1]
		}
	}
	return split
}
