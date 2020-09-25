package repotools

import (
	"os"
	"path/filepath"
)

// Boots was made for walking the file tree searching for modules.
type Boots struct {
	// Directories to skip when iterating.
	SkipDirs []string

	modulePaths []string
}

// Modules returns a slice of module directory absolute paths.
func (b *Boots) Modules() []string {
	return b.modulePaths
}

// Walk is the function to walk folders in the repo searching for go modules.
func (b *Boots) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	for _, skip := range b.SkipDirs {
		if path == skip {
			return filepath.SkipDir
		}
	}

	if !hasGoMod(path) {
		return nil
	}

	b.modulePaths = append(b.modulePaths, path)

	return nil
}

func hasGoMod(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "go.mod"))
	if err != nil {
		return false
	}
	return true
}
