package gomod

import (
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/semver"
	"golang.org/x/mod/modfile"
)

// UpdateRequires updates all modules discovered starting at repoRootPath using the provided tags and dependencies.
// Using force will update the module required versions regardless whether the target version less the currently
// written version.
func UpdateRequires(repoRootPath string, tags git.ModuleTags, dependencies map[string]string, force bool) error {
	discoverer := NewDiscoverer(repoRootPath)

	if err := discoverer.Discover(); err != nil {
		return fmt.Errorf("failed to discover repository modules: %v", err)
	}

	modules, err := discoverer.ModulesRel()
	if err != nil {
		return err
	}

	repoModules := make(map[string]struct {
		ModuleDir string
		File      *modfile.File
	})

	for moduleDir := range modules {
		mod, err := LoadModuleFile(filepath.Join(discoverer.Root(), moduleDir), nil, true)
		if err != nil {
			return fmt.Errorf("failed to load module file: %w", err)
		}
		rm := repoModules[mod.Module.Mod.Path]
		rm.File = mod
		rm.ModuleDir = moduleDir
		repoModules[mod.Module.Mod.Path] = rm
	}

	for _, mod := range repoModules {
		for _, require := range mod.File.Require {
			version := require.Mod.Version
			if requireMod, ok := repoModules[require.Mod.Path]; ok {
				latest, ok := tags.Latest(requireMod.ModuleDir)
				if ok {
					if force {
						version = latest
					} else if semver.Compare(latest, version) > 0 {
						version = latest
					}
				}
			} else {
				dv, ok := dependencies[require.Mod.Path]
				if ok {
					if force {
						version = dv
					} else if semver.Compare(dv, version) > 0 {
						version = dv
					}
				}
			}
			if err := mod.File.AddRequire(require.Mod.Path, version); err != nil {
				return err
			}
		}

		if err = WriteModuleFile(filepath.Join(discoverer.Root(), mod.ModuleDir), mod.File); err != nil {
			return fmt.Errorf("failed to write module file: %w", err)
		}
	}

	return nil
}
