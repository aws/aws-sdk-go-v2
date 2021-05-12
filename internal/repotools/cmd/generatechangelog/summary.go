package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/release"
)

func executeModuleTemplate(wr io.Writer, summary moduleSummary) error {
	return moduleChangeLogTemplate.Execute(wr, summary)
}

func writeModuleChangeLog(path string, summary moduleSummary) error {
	changelogPath := filepath.Join(path, changeLogFile)

	t, cleanup, err := copyToTempFile(changelogPath)
	if err != nil {
		return err
	}
	defer func() {
		cErr := cleanup()
		if err == nil && cErr != nil {
			err = cErr
		}
	}()

	f, err := os.OpenFile(changelogPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		fErr := f.Close()
		if err == nil && fErr != nil {
			err = fErr
		}
	}()

	if err = executeModuleTemplate(f, summary); err != nil {
		return err
	}

	if _, err = io.Copy(f, t); err != nil {
		return err
	}

	return nil
}

func executeRepoChangeLogEntryTemplate(wr io.Writer, summary releaseSummary) error {
	return repoChangeLogTemplate.ExecuteTemplate(wr, "entry", summary)
}

func writeRepoChangeLogEntry(path string, summary releaseSummary) (err error) {
	changelogPath := filepath.Join(path, changeLogFile)

	t, cleanup, err := copyToTempFile(changelogPath)
	if err != nil {
		return err
	}
	defer func() {
		cErr := cleanup()
		if err == nil && cErr != nil {
			err = cErr
		}
	}()

	f, err := os.OpenFile(changelogPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		fErr := f.Close()
		if err == nil && fErr != nil {
			err = fErr
		}
	}()

	if err = executeRepoChangeLogEntryTemplate(f, summary); err != nil {
		return err
	}

	if _, err = io.Copy(f, t); err != nil {
		return err
	}

	return nil
}

func executeSummaryNotesTemplate(wr io.Writer, summary releaseSummary) error {
	return repoChangeLogTemplate.ExecuteTemplate(wr, "summary", summary)
}

func writeSummaryNotes(path string, summary releaseSummary) (err error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		fErr := f.Close()
		if err == nil && fErr != nil {
			err = fErr
		}
	}()

	if err = executeSummaryNotesTemplate(f, summary); err != nil {
		return err
	}

	return nil
}

type moduleSummary struct {
	ReleaseID   string
	ModulePath  string
	Version     string
	Annotations []changelog.Annotation
}

type releaseSummary struct {
	ReleaseID string
	General   []changelog.Annotation
	Modules   map[string]moduleSummary
}

func (r releaseSummary) IsEmptyReleaseSummary() bool {
	if len(r.General) != 0 {
		return false
	}

	for _, summary := range r.Modules {
		if len(summary.Annotations) > 0 {
			return false
		}
	}

	return true
}

var dependencyBump = changelog.Annotation{
	Type:        changelog.DependencyChangeType,
	Collapse:    true,
	Description: "Updated to the latest SDK module versions",
}

func generateSummary(manifest release.Manifest, annotations []changelog.Annotation) (releaseSummary, error) {
	summary := releaseSummary{
		ReleaseID: manifest.ID,
		Modules:   make(map[string]moduleSummary),
	}

	idToAnnotation := make(map[string]changelog.Annotation)
	for _, annotation := range annotations {
		idToAnnotation[annotation.ID] = annotation
	}

	hasDependencyBumps := false

	for modDir, mod := range manifest.Modules {
		ms := moduleSummary{
			ReleaseID:  manifest.ID,
			ModulePath: mod.ModulePath,
			Version:    mod.To,
		}
		for _, id := range mod.Annotations {
			an, ok := idToAnnotation[id]
			if !ok {
				continue
			}
			ms.Annotations = append(ms.Annotations, an)
		}
		if mod.Changes&release.DependencyUpdate != 0 {
			ms.Annotations = append(ms.Annotations, dependencyBump)
			hasDependencyBumps = true
		}
		sortAnnotations(ms.Annotations)
		summary.Modules[modDir] = ms
	}

	for _, annotation := range annotations {
		if annotation.Collapse {
			summary.General = append(summary.General, annotation)
		}
	}

	if hasDependencyBumps {
		summary.General = append(summary.General, dependencyBump)
	}

	sortAnnotations(summary.General)

	return summary, nil
}
