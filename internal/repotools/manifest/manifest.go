package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// NoManifestFound is an erro returned when a build manifest file was not found in the designated directory path.
type NoManifestFound struct {
	Path string
}

// Error returns the error description
func (n *NoManifestFound) Error() string {
	return "build artifact manifest not found"
}

// Manifest is a description of a Smithy build artifact produced
// by the github.com/aws/smithy-go plugin.
type Manifest struct {
	Module       string            `json:"module"`
	Go           string            `json:"go"`
	Dependencies map[string]string `json:"dependencies"`
	Files        []string          `json:"files"`
	Unstable     bool              `json:"unstable"`
}

// ValidateManifest validates that the build artifact description
// has the minimum required information to produce a valid Go module description.
func ValidateManifest(manifest Manifest) error {
	if len(manifest.Go) == 0 {
		return fmt.Errorf("missing Go minimum version")
	}
	if len(manifest.Module) == 0 {
		return fmt.Errorf("missing module path")
	}
	return nil
}

// LoadManifest loads the manifest description from the file located at the given path.
func LoadManifest(path string) (Manifest, error) {
	mf, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		return Manifest{}, &NoManifestFound{Path: path}
	} else if err != nil {
		return Manifest{}, fmt.Errorf("failed to open manifest: %w", err)
	}
	defer mf.Close()
	return ReadManifest(mf)
}

// ReadManifest parses the manifest bytes from the provided reader and returns the manifest description.
func ReadManifest(reader io.Reader) (m Manifest, err error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return Manifest{}, err
	}
	if err = json.Unmarshal(data, &m); err != nil {
		return Manifest{}, err
	}
	if err = ValidateManifest(m); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// SmithyArtifactPaths is a slice of smithy-go build artifacts.
// See the Walk method which can be used for finding the generated Go
// source code from the Smithy build plugins projection.
type SmithyArtifactPaths []string

// Walk is a filepath.WalkFunc compatible method that can be used for finding
// smithy-go plugin build artifacts.
func (a *SmithyArtifactPaths) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	pluginOutput := filepath.Join(path, "go-codegen")
	stat, err := os.Stat(pluginOutput)
	if err != nil {
		return nil
	}

	if !stat.IsDir() {
		return nil
	}

	*a = append(*a, pluginOutput)

	return filepath.SkipDir
}
