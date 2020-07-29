package changes

import (
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes/golist"
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetCurrentModule(t *testing.T) {
	mod, err := GetCurrentModule()
	if err != nil {
		t.Errorf("expected nil err, got %v", err)
	}

	if mod != "internal/tools/changes" {
		t.Errorf("expected mod to be \"internal/tools/changes\", got %s", mod)
	}
}

func TestDiscoverModules(t *testing.T) {
	const prefix = "internal/tools/changes/testdata/modules/"
	wantMods := []string{
		"a",
		"b",
		"nested/c/d",
		"nested/c",
	}
	for i := range wantMods {
		wantMods[i] = prefix + wantMods[i]
	}

	goclient := golist.Client{
		RootPath: filepath.Join("testdata", "modules"),
		ShortenModPath: func(mod string) string {
			return strings.TrimPrefix(mod, "internal/tools/changes/testdata/modules/")
		},
	}

	mods, _, err := discoverModules(goclient, filepath.Join("testdata", "modules"))
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantMods, mods); diff != "" {
		t.Errorf("expect modules to match:\n%v", diff)
	}
}

func TestDefaultVersion(t *testing.T) {
	var cases = map[string]struct {
		mod         string
		wantVersion string
		wantErr     string
	}{
		"v0": {
			mod:         "example.com/module",
			wantVersion: "v0.0.0",
		},
		"v2": {
			mod:         "example.com/module/v2",
			wantVersion: "v2.0.0",
		},
		"v100": {
			mod:         "example.com/module/v100",
			wantVersion: "v100.0.0",
		},
		"invalid": {
			mod:     "example.com/module/v2.0",
			wantErr: "couldn't split module path",
		},
	}

	for id, tt := range cases {
		t.Run(id, func(t *testing.T) {
			v, err := defaultVersion(tt.mod)
			if tt.wantErr != "" {
				if err == nil {
					t.Error("expected non-nil err, got nil")
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected err to contain %s, got %v", tt.wantErr, err)
				}
			}

			if tt.wantVersion != v {
				t.Errorf("expected version to be %s, got %s", tt.wantVersion, v)
			}
		})
	}
}

func TestPseudoVersion(t *testing.T) {
	_, err := pseudoVersion(".", "internal/tools/changes")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFormatPseudoVersion(t *testing.T) {
	var cases = map[string]struct {
		hash        string
		tagVersion  string
		wantVersion string
	}{
		"no tag": {
			hash:        "20200709182313-123456789012",
			tagVersion:  "",
			wantVersion: "v0.0.0-20200709182313-123456789012",
		},
		"v0.0.0 tag": {
			hash:        "20200709182313-123456789012",
			tagVersion:  "v0.0.0",
			wantVersion: "v0.0.1-0.20200709182313-123456789012",
		},
		"v1.2.3 tag": {
			hash:        "20200709182313-123456789012",
			tagVersion:  "v1.2.3",
			wantVersion: "v1.2.4-0.20200709182313-123456789012",
		},
		"v1.2.3 pre tag": {
			hash:        "20200709182313-123456789012",
			tagVersion:  "v1.2.3-pre",
			wantVersion: "v1.2.3-pre.0.20200709182313-123456789012",
		},
		"v2.0.0 beta tag": {
			hash:        "20200709182313-123456789012",
			tagVersion:  "v2.0.0-beta",
			wantVersion: "v2.0.0-beta.0.20200709182313-123456789012",
		},
		"invalid": {
			hash:        "20200709182313-123456789012",
			tagVersion:  "v2.0.0-beta+build.tag",
			wantVersion: "",
		},
	}

	for id, tt := range cases {
		t.Run(id, func(t *testing.T) {
			pseudoV, err := formatPseudoVersion(tt.hash, tt.tagVersion)
			if tt.wantVersion == "" {
				if err == nil {
					t.Error("expected non-nil err, got nil")
				}
			} else if pseudoV != tt.wantVersion {
				t.Errorf("expected pseudo-version to be %s, got %s", tt.wantVersion, pseudoV)
			}
		})
	}
}

func TestCommitHash(t *testing.T) {
	commitHash, err := commitHash(".")
	if err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(commitHash, "-")
	if len(parts) != 2 {
		t.Errorf("expected commit hash to have 2 parts separated by '-', got %s", commitHash)
	}

	if len(parts[1]) != 12 {
		t.Errorf("expected commit hash length to be 12, got %d", len(parts[1]))
	}
}

type mockGolist struct {
	dependencies map[string][]string
	packages     map[string][]string
}

func (c *mockGolist) Dependencies(mod string) ([]string, error) {
	return c.dependencies[mod], nil
}

func (c *mockGolist) Packages(mod string) ([]string, error) {
	return c.packages[mod], nil
}
