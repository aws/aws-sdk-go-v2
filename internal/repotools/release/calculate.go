package release

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/gomod"
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

			hasChanges, err = gomod.IsModuleChanged(moduleDir, repositoryModules[moduleDir], changes)
			if err != nil {
				return nil, fmt.Errorf("failed to determine module changes: %w", err)
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
