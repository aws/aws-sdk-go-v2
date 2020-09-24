package main

import (
	"os"
	"path/filepath"
)

// Boots was made for walking the file tree searching for modules.
type Boots struct {
	modulePaths []string
	skipDirs    []string
}

func (b *Boots) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	for _, skip := range b.skipDirs {
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
