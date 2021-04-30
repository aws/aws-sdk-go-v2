package changelog

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
)

const changelogDir = ".changelog"

// Annotation represents a change to one or more Go modules.
type Annotation struct {
	// The unique identifier for this Annotation.
	ID string `json:"id" toml:"id" comment:"Annotation Identifier (DO NOT CHANGE)"`

	// Indicates what category of Annotation was made. For example "feature" or "bugfix".
	Type ChangeType `json:"type" toml:"type" comment:"Valid Types: announcement, release, feature, bugfix, documentation, or dependency"`

	// Collapse indicates that the change description should collapsed into a single item when summarizing changes across modules
	Collapse bool `json:"collapse,omitempty" toml:"collapse" comment:"annotation should collapse as a summary in the CHANGELOG"`

	// A human readable description of this Annotation meant to be included in a CHANGELOG.
	Description string `json:"description" toml:"description" comment:"single-line string or markdown list"`

	// The modules this change applies to
	Modules []string `json:"modules"  toml:"modules" comment:"one or more relative module paths"`
}

// ValidationError is an error that indicates that one ore more issues are present for an annotation.
type ValidationError struct {
	Issues []string
}

// Error returns the error string
func (v *ValidationError) Error() string {
	var sb strings.Builder
	sb.WriteString("invalid change annotation:\n")
	for _, issue := range v.Issues {
		sb.WriteRune('\t')
		sb.WriteString(issue)
		sb.WriteRune('\n')
	}
	sb.WriteRune('\n')
	return sb.String()
}

// Validate returns an error if annotation does not mean the minimum requirements.
func Validate(annotation Annotation) error {
	var invalid []string

	if len(annotation.ID) == 0 {
		invalid = append(invalid, "annotation id is required")
	}

	if annotation.Type == UnknownChangeType {
		invalid = append(invalid, fmt.Sprintf("invalid change type"))
	}

	if len(annotation.Description) == 0 {
		invalid = append(invalid, "description is required")
	}

	if len(annotation.Modules) < 1 {
		invalid = append(invalid, "at least one module is required")
	}

	if len(invalid) > 0 {
		return &ValidationError{Issues: invalid}
	}

	return nil
}

// WriteAnnotation writes the annotation to changelog metadata directory.
// Path should be the location of the repository root.
func WriteAnnotation(path string, annotation Annotation) (err error) {
	dir := filepath.Join(path, changelogDir)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	marshal, err := json.MarshalIndent(annotation, "", "    ")
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s.json", strings.ReplaceAll(annotation.ID, "-", ""))

	f, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		fErr := f.Close()
		if err == nil && fErr != nil {
			err = fErr
		}
	}()

	_, err = io.Copy(f, bytes.NewReader(marshal))

	return err
}

// RemoveAnnotation removes the annotation from the changelog metadata directory.
func RemoveAnnotation(path string, annotation Annotation) (err error) {
	dir := filepath.Join(path, changelogDir)

	name := fmt.Sprintf("%s.json", strings.ReplaceAll(annotation.ID, "-", ""))

	if err := os.Remove(filepath.Join(dir, name)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// GetAnnotations returns the list of annotations that are currently present.
// Path should be the location of the repository root.
func GetAnnotations(path string) (annotations []Annotation, err error) {
	dir := filepath.Join(path, changelogDir)

	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == dir {
			return err
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		if !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}

		annotation, err := LoadAnnotationFile(path)
		if err != nil {
			return err
		}

		annotations = append(annotations, annotation)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return annotations, nil
}

// LoadAnnotationFile loads the annotation file at the given path.
func LoadAnnotationFile(path string) (a Annotation, err error) {
	fBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Annotation{}, err
	}

	if err = json.Unmarshal(fBytes, &a); err != nil {
		return Annotation{}, err
	}

	return a, nil
}

// LoadAnnotation loads the annotation id for the given repository path.
func LoadAnnotation(path string, id string) (a Annotation, err error) {
	path = filepath.Join(path, changelogDir, strings.ReplaceAll(id, "-", "")+".json")
	return LoadAnnotationFile(path)
}

// GetVersionIncrement returns the highest version increment from a set of annotations.
func GetVersionIncrement(annotations []Annotation) (v SemVerIncrement) {
	for _, annotation := range annotations {
		vi := annotation.Type.VersionIncrement()
		if vi > v {
			v = vi
		}
	}
	return v
}

// NewAnnotation creates a new annotation with a populated identifier.
func NewAnnotation() (Annotation, error) {
	var b [16]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		return Annotation{}, err
	}
	return Annotation{
		ID: repotools.UUIDVersion4(b),
	}, nil
}
