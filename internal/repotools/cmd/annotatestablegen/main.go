package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/gomod"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/semver"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const (
	// reservedChangeID is the string "AWSSDK@\xff\xbfGOAUTO\x00"
	// Note: future needs can increment the last \x00 value up to \xFF for 256 pseudo-random values for automation
	// release notes.
	// TODO: centralize these somewhere if we need to extend this past this one UUID.
	reservedChangeID = "41575353-444b-40ff-bf47-4f4155544f00"

	description = "New AWS service client module"

	generatedFile = "generated.json"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}

	repoRoot, err := repotools.FindRepoRoot(cwd)
	if err != nil {
		log.Fatalf("failed to get git repo root: %v", err)
	}

	tags, err := git.Tags(repoRoot)
	if err != nil {
		log.Fatalf("failed to get git repository tags: %v", err)
	}

	moduleTags := git.ParseModuleTags(tags)

	discoverer := gomod.NewDiscoverer(repoRoot)

	if err := discoverer.Discover(); err != nil {
		log.Fatalf("failed to discover go modules: %v", err)
	}

	modules, err := discoverer.ModulesRel()
	if err != nil {
		log.Fatalf("failed to get modules relative to repo: %v", err)
	}

	var toRelease []string

	for modDir := range modules {
		fullPath := filepath.Join(repoRoot, modDir)

		isGenerated, err := isGeneratedModule(fullPath)
		if err != nil {
			log.Fatalf("failed to determine if module is generated: %v", err)
		}
		if !isGenerated {
			continue
		}

		stable, err := isGeneratedModuleStable(fullPath)
		if err != nil {
			log.Fatalf("failed to determine if generated module is stable: %v", err)
		}

		latest, ok := moduleTags.Latest(modDir)
		if !ok && stable {
			toRelease = append(toRelease, modDir)
		} else if ok && stable && len(semver.Prerelease(latest)) > 0 {
			toRelease = append(toRelease, modDir)
		}
	}

	if len(toRelease) == 0 {
		log.Printf("[INFO] no generated modules require release annotation")
		return
	}

	var annotation changelog.Annotation

	annotation, err = changelog.LoadAnnotation(repoRoot, reservedChangeID)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to load module annotation: %v", err)
	}
	if err != nil {
		annotation, err = changelog.NewAnnotation()
		if err != nil {
			log.Fatalf("failed to generated annotation: %v", err)
		}
		annotation.ID = reservedChangeID
		annotation.Type = changelog.ReleaseChangeType
		annotation.Description = description
	} else if annotation.Type != changelog.ReleaseChangeType {
		log.Fatalf("annotation type does not match the expected type")
	}

	sort.Strings(annotation.Modules)

	for _, modDir := range toRelease {
		annotation.Modules = repotools.AppendIfNotPresent(annotation.Modules, modDir)
	}

	if err := changelog.WriteAnnotation(repoRoot, annotation); err != nil {
		log.Fatalf("failed to write annotation: %v", err)
	}
}

type generated struct {
	Unstable bool `json:"unstable"`
}

func isGeneratedModuleStable(dir string) (bool, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(dir, generatedFile))
	if err != nil {
		return false, err
	}

	var generated generated
	if err := json.Unmarshal(bytes, &generated); err != nil {
		return false, err
	}

	return !generated.Unstable, nil
}

func isGeneratedModule(dir string) (bool, error) {
	_, err := os.Stat(filepath.Join(dir, generatedFile))
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
