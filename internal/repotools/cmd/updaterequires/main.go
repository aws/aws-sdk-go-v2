package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/gomod"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/release"
	"io/ioutil"
	"log"
)

var config = struct {
	ReleaseManifestPath string

	Force bool
}{}

func init() {
	flag.StringVar(&config.ReleaseManifestPath, "release", "", "file path to a release manifest containing module tags to be released overlayed")
	flag.BoolVar(&config.Force, "force", false, "force module versions regardless of the current recorded version")
}

func main() {
	flag.Parse()

	repoRootPath, err := repotools.GetRepoRoot()
	if err != nil {
		log.Fatalf("failed to get repository details: %v", err)
	}

	tags, err := getRepoTags(repoRootPath)
	if err != nil {
		log.Fatalf("failed to retrieve git tags: %v", err)
	}

	dependencies, err := getDependencies(repoRootPath)
	if err != nil {
		log.Fatalf("failed to load dependency versions: %v", err)
	}

	if len(config.ReleaseManifestPath) > 0 {
		if err := applyOverlayTags(config.ReleaseManifestPath, tags); err != nil {
			log.Fatalf("failed to apply tag overlay: %v", err)
		}
	}

	if err := gomod.UpdateRequires(repoRootPath, tags, dependencies, config.Force); err != nil {
		log.Fatalf("failed to update module dependencies: %v", err)
	}
}

func applyOverlayTags(path string, tags git.ModuleTags) error {
	fBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var manifest release.Manifest

	if err := json.Unmarshal(fBytes, &manifest); err != nil {
		return err
	}

	for _, tag := range manifest.Tags {
		if len(tag) == 0 {
			continue
		}
		tags.Add(tag)
	}

	return nil
}

func getDependencies(path string) (map[string]string, error) {
	loadConfig, err := repotools.LoadConfig(path)
	if err != nil {
		return nil, err
	}
	return loadConfig.Dependencies, nil
}

func getRepoTags(path string) (git.ModuleTags, error) {
	if err := git.Fetch(path); err != nil {
		return nil, err
	}

	tags, err := git.Tags(path)
	if err != nil {
		return nil, err
	}

	return git.ParseModuleTags(tags), nil
}
