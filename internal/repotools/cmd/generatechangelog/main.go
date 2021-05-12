package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/release"
)

const changeLogFile = "CHANGELOG.md"

var releaseManifestFile, summaryNotesFile string

func init() {
	flag.StringVar(&releaseManifestFile, "release", "", "release manifest file")
	flag.StringVar(&summaryNotesFile, "o", "", "indicates that a copy of the changelog notes should be written to the target file")
}

func main() {
	flag.Parse()

	if len(releaseManifestFile) == 0 {
		log.Fatalln("first argument should be a release manifest file")
	}

	manifest, err := loadManifest(releaseManifestFile)
	if err != nil {
		log.Fatalf("failed to load release manifest file: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}

	repoRoot, err := repotools.FindRepoRoot(cwd)
	if err != nil {
		log.Fatalf("failed to get git repository root: %v", err)
	}

	annotations, err := changelog.GetAnnotations(repoRoot)
	if err != nil {
		log.Fatalf("failed to get change annotations: %v", err)
	}

	annotations = filterUnreferencedAnnotations(manifest, annotations)

	summary, err := generateSummary(manifest, annotations)
	if err != nil {
		log.Fatalf("failed to generate summary: %v", err)
	}

	if err := writeRepoChangeLogEntry(repoRoot, summary); err != nil {
		log.Fatalf("failed to write summary CHANGELOG.md")
	}

	for moduleDir, ms := range summary.Modules {
		if moduleDir == "." {
			// The root module contains the repository changelog
			continue
		}
		if err = writeModuleChangeLog(filepath.Join(repoRoot, moduleDir), ms); err != nil {
			log.Fatalf("failed to write module CHANGELOG.md: %v", err)
		}
	}

	if len(summaryNotesFile) > 0 {
		if err := writeSummaryNotes(summaryNotesFile, summary); err != nil {
			log.Fatalf("failed to write summary notes: %v", err)
		}
	}
}

func copyToTempFile(name string) (io.ReadSeeker, func() error, error) {
	if _, err := os.Stat(name); err != nil && os.IsNotExist(err) {
		return bytes.NewReader(nil), func() error {
			return nil
		}, nil
	} else if err != nil {
		return nil, nil, err
	}

	t, err := ioutil.TempFile("", "CHANGELOG-*.md")
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() error {
		return os.Remove(t.Name())
	}

	f, err := os.Open(name)
	if err != nil {
		_ = cleanup()
		return nil, nil, err
	}
	defer f.Close()

	_, err = io.Copy(t, f)
	if err != nil {
		_ = cleanup()
		return nil, nil, err
	}

	if _, err := t.Seek(0, io.SeekStart); err != nil {
		_ = cleanup()
		return nil, nil, err
	}

	return t, cleanup, nil
}

func filterUnreferencedAnnotations(manifest release.Manifest, annotations []changelog.Annotation) (filtered []changelog.Annotation) {
	referenced := make(map[string]struct{})

	for _, module := range manifest.Modules {
		for _, id := range module.Annotations {
			referenced[id] = struct{}{}
		}
	}

	for _, annotation := range annotations {
		if _, ok := referenced[annotation.ID]; ok {
			filtered = append(filtered, annotation)
		} else {
			log.Printf("[WARN] Annotation %v will be ignored as it does not annotate any module changes", annotation.ID)
		}
	}

	return filtered
}

func loadManifest(path string) (v release.Manifest, err error) {
	fb, err := ioutil.ReadFile(path)
	if err != nil {
		return release.Manifest{}, err
	}

	if err = json.Unmarshal(fb, &v); err != nil {
		return release.Manifest{}, err
	}

	return v, nil
}

// sortAnnotations sorts from their highest numerical order to lowest
func sortAnnotations(annotations []changelog.Annotation) {
	sort.Slice(annotations, func(i, j int) bool {
		if annotations[i].Type < annotations[j].Type {
			return false
		}
		if annotations[i].Type != annotations[j].Type {
			return true
		}
		return annotations[i].Description < annotations[j].Description
	})
}

func inlineCodeBlock(v string) string {
	return "`" + v + "`"
}
