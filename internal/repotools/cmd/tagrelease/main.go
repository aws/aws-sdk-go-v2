package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/release"
)

var releaseFile string

func init() {
	flag.StringVar(&releaseFile, "release", "", "release manifest file path")
}

func main() {
	flag.Parse()

	if len(releaseFile) == 0 {
		log.Fatal("release manifest file path must be provided")
	}

	repoRoot, err := repotools.GetRepoRoot()
	if err != nil {
		log.Fatalf("failed to get repository root: %v", err)
	}

	manifest, err := loadManifest(releaseFile)
	if err != nil {
		log.Fatalf("failed to laod manifest: %v", err)
	}

	if len(manifest.Tags) == 0 {
		log.Println("[INFO] no modules for release")
		return
	}

	if err = git.Add(repoRoot, "-A", "."); err != nil {
		log.Fatalf("failed to add working directory changes: %v", err)
	}

	message := fmt.Sprintf("Release %s", manifest.ID)

	if err = git.Commit(repoRoot, message); err != nil {
		log.Fatalf("failed to add working directory changes: %v", err)
	}

	for _, tag := range manifest.Tags {
		if err := git.Tag(repoRoot, tag, message, "HEAD"); err != nil {
			log.Fatalf("failed to create tag %v: %v", tag, err)
		}
	}

	releaseTag := fmt.Sprintf("release-%s", manifest.ID)
	if err = git.Tag(repoRoot, releaseTag, message, "HEAD"); err != nil {
		log.Fatalf("failed to create release tag %v: %v", releaseTag, err)
	}
}

func loadManifest(path string) (manifest release.Manifest, err error) {
	fBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return release.Manifest{}, err
	}

	if err := json.Unmarshal(fBytes, &manifest); err != nil {
		return release.Manifest{}, err
	}

	return manifest, nil
}
