package gomod

import (
	"path"
	"strings"
)

// IsModuleChanged returns whether the given set of changes applies to the module.
// The submodules argument is must be lexicographically sorted list of submodule locations that are located
// under moduleDir.
func IsModuleChanged(moduleDir string, submodules []string, changes []string) (bool, error) {
	if moduleDir == "." {
		moduleDir = ""
	}

	isChildPathCache := make(map[string]bool)

	hasChanges := false

	for i := 0; i < len(changes) && !hasChanges; i++ {
		dir, fileName := path.Split(changes[i])

		if len(dir) == 0 && moduleDir != "" {
			continue
		}

		if len(moduleDir) > 0 && !strings.HasPrefix(dir, moduleDir) {
			continue
		}

		if len(dir) == 0 && (IsGoSource(fileName) || IsGoMod(fileName)) {
			hasChanges = true
			continue
		} else if !(IsGoSource(fileName) || IsGoMod(fileName)) {
			continue
		}
		dir = path.Clean(dir)

		if len(submodules) == 0 {
			hasChanges = true
			continue
		}

		if isChild, ok := isChildPathCache[dir]; !ok {
			if IsSubmodulePath(dir, submodules) {
				isChildPathCache[dir] = true
			} else {
				isChildPathCache[dir] = false
				hasChanges = true
			}
		} else if !isChild {
			hasChanges = true
		}
	}

	return hasChanges, nil
}

// IsGoSource returns whether a given file name is a Go source code file ending in `.go`
func IsGoSource(name string) bool {
	return !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

// IsGoMod returns whether a given file name is `go.mod`.
func IsGoMod(name string) bool {
	return name == "go.mod"
}
