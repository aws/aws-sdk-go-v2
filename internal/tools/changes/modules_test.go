package changes

import (
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
	for i, _ := range wantMods {
		wantMods[i] = prefix + wantMods[i]
	}

	mods, _, err := discoverModules(filepath.Join("testdata", "modules"))
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
