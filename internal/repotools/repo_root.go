package repotools

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindRepoRoot returns the absolute path to the root directory of the
// repository, or error. If the dir passed in is a relative path it will be
// used relative to the current working directory of the executable.
func FindRepoRoot(dir string) (string, error) {
	if len(dir) == 0 {
		dir = "."
	}

	if !filepath.IsAbs(dir) {
		var err error
		dir, err = JoinWorkingDirectory(dir)
		if err != nil {
			return "", err
		}
	}

	var found bool
	for {
		if dir == string(filepath.Separator) {
			break
		}

		_, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil {
			found = true
			break
		}

		dir = filepath.Dir(dir)
	}

	if !found {
		return "", fmt.Errorf(".git directory not found")
	}

	return dir, nil
}

// JoinWorkingDirectory will return an absolute file system path of the passed
// in dir path with the current working directory.
func JoinWorkingDirectory(dir string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory, %w", err)
	}

	return filepath.Join(wd, dir), nil
}
