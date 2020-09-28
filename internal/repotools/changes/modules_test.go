package changes

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes/golist"
	"github.com/google/go-cmp/cmp"
)

func TestGetCurrentModule(t *testing.T) {
	mod, err := GetCurrentModule()
	if err != nil {
		t.Errorf("expected nil err, got %v", err)
	}

	if mod != "internal/repotools/changes" {
		t.Errorf("expected mod to be \"internal/repotools/changes\", got %s", mod)
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

	mods, err := discoverModules(goclient.RootPath)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantMods, mods); diff != "" {
		t.Errorf("expect modules to match (-want, +got):\n%v", diff)
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
			wantVersion: "v0.1.0",
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
	v, err := pseudoVersion(&MockGit{
		tags: []string{
			"internal/tools/changes/v0.1.2",
		},
	}, "internal/tools/changes")
	if err != nil {
		t.Fatal(err)
	}

	wantVer := "v0.1.3-0.1234567abcde"

	if v != wantVer {
		t.Errorf("wanted version %s, got %s", wantVer, v)
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

func (c *mockGolist) Checksum(mod, version string) (string, error) {
	return "01234567abcdef", nil
}

func (c *mockGolist) Tidy(mod string) error {
	return nil
}
