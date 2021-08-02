package release

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/gomod"
	"log"
	"path"
	"path/filepath"
	"sort"
)

// ModuleFinder is a type that
type ModuleFinder interface {
	Root() string

	ModulesRel() (map[string][]string, error)
}

// Calculate calculates the modules to be released and their next versions based on the Git history, previous tags,
// module configuration, and associated changelog annotaitons.
func Calculate(finder ModuleFinder, tags git.ModuleTags, config repotools.Config, annotations []changelog.Annotation) (map[string]*Module, error) {
	rootDir := finder.Root()

	repositoryModules, err := finder.ModulesRel()
	if err != nil {
		log.Fatalf("failed to modules: %v", err)
	}

	moduleAnnotations := make(map[string][]changelog.Annotation)
	for _, annotation := range annotations {
		for _, am := range annotation.Modules {
			moduleAnnotations[am] = append(moduleAnnotations[am], annotation)
		}
	}

	modules := make(map[string]*Module)
	var repositoryModuleTombstonePaths []string

	for moduleDir := range tags {
		if _, ok := repositoryModules[moduleDir]; !ok {
			repositoryModuleTombstonePaths = append(repositoryModuleTombstonePaths, moduleDir)
		}
	}

	for moduleDir := range repositoryModules {
		moduleFile, err := gomod.LoadModuleFile(filepath.Join(rootDir, moduleDir), nil, true)
		if err != nil {
			return nil, fmt.Errorf("failed to load module file: %w", err)
		}

		modulePath, err := gomod.GetModulePath(moduleFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read module path: %w", err)
		}

		var latestVersion string
		var hasChanges bool

		latestVersion, ok := tags.Latest(moduleDir)
		if ok {
			startTag, err := git.ToModuleTag(moduleDir, latestVersion)
			if err != nil {
				log.Fatalf("failed to convert module path and version to tag: %v", err)
			}

			changes, err := git.Changes(finder.Root(), startTag, "HEAD", moduleDir)
			if err != nil {
				log.Fatalf("failed to get git changes: %v", err)
			}

			subModulePaths := repositoryModules[moduleDir]

			ignoredModulePaths := make([]string, 0, len(subModulePaths)+len(repositoryModuleTombstonePaths))
			ignoredModulePaths = append(ignoredModulePaths, subModulePaths...)

			if len(repositoryModuleTombstonePaths) > 0 {
				ignoredModulePaths = append(ignoredModulePaths, repositoryModuleTombstonePaths...)
				// IsModuleChanged expects the provided list of ignored modules paths to be sorted
				sort.Strings(ignoredModulePaths)
			}

			hasChanges, err = gomod.IsModuleChanged(moduleDir, ignoredModulePaths, changes)
			if err != nil {
				return nil, fmt.Errorf("failed to determine module changes: %w", err)
			}

			if !hasChanges {
				// Check if any of the submodules have been "carved out" of this module since the last tagged release
				for _, subModuleDir := range subModulePaths {
					if _, ok := tags.Latest(subModuleDir); ok {
						continue
					}

					treeFiles, err := git.LsTree(rootDir, startTag, subModuleDir)
					if err != nil {
						return nil, fmt.Errorf("failed to list git tree: %v", err)
					}

					carvedOut, err := isModuleCarvedOut(treeFiles, repositoryModules[subModuleDir])
					if err != nil {
						return nil, err
					}
					if carvedOut {
						hasChanges = true
						break
					}
				}
			}
		}

		var changeReason ModuleChange
		if hasChanges && len(latestVersion) > 0 {
			changeReason |= SourceChange
		} else if len(latestVersion) == 0 {
			changeReason |= NewModule
		}

		modules[modulePath] = &Module{
			File:              moduleFile,
			RelativeRepoPath:  moduleDir,
			Latest:            latestVersion,
			Changes:           changeReason,
			ChangeAnnotations: moduleAnnotations[moduleDir],
			ModuleConfig:      config.Modules[moduleDir],
		}
	}

	if err := CalculateDependencyUpdates(modules); err != nil {
		return nil, err
	}

	for moduleDir := range modules {
		if modules[moduleDir].Changes == 0 || config.Modules[moduleDir].NoTag {
			delete(modules, moduleDir)
		}
	}

	return modules, nil
}

// isModuleCarvedOut takes a list of files for a (new) submodule directory. The list of files are the files that are located
// in the submodule directory path from the parent's previous tagged release. Returns true the new submodule has been
// carved out of the parent module directory it is located under. This is determined by looking through the file list
// and determining if Go source is present but no `go.mod` file existed.
func isModuleCarvedOut(files []string, subModules []string) (bool, error) {
	hasGoSource := false
	hasGoMod := false

	isChildPathCache := make(map[string]bool)

	for _, file := range files {
		dir, fileName := path.Split(file)
		dir = path.Clean(dir)

		isGoMod := gomod.IsGoMod(fileName)
		isGoSource := gomod.IsGoSource(fileName)

		if !(isGoMod || isGoSource) {
			continue
		}

		if isChild, ok := isChildPathCache[dir]; (isChild && ok) || (!ok && gomod.IsSubmodulePath(dir, subModules)) {
			isChildPathCache[dir] = true
			continue
		} else {
			isChildPathCache[dir] = false
		}

		if isGoSource {
			hasGoSource = true
		} else if isGoMod {
			hasGoMod = true
		}

		if hasGoMod && hasGoSource {
			break
		}
	}

	return !hasGoMod && hasGoSource, nil
}
