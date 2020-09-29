package changes

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestVersionEnclosure_IsValid(t *testing.T) {
	const invalidErr = "does not match git tag version"

	var testCases = map[string]struct {
		tags    []string
		enc     VersionEnclosure
		wantErr string
	}{
		"valid 1 module": {
			[]string{
				"a/v1.0.0",
			},
			VersionEnclosure{
				ModuleVersions: map[string]Version{
					"a": {
						Module:  "a",
						Version: "v1.0.0",
					},
				},
			},
			"",
		},
		"valid multi module": {
			[]string{
				"a/v1.0.0",
				"b/v2.1.3",
			},
			VersionEnclosure{
				ModuleVersions: map[string]Version{
					"a": {
						Module:  "a",
						Version: "v1.0.0",
					},
					"b/v2": {
						Module:  "b/v2",
						Version: "v2.1.3",
					},
				},
			},
			"",
		},
		"invalid missing": {
			[]string{
				"b/v2.1.3",
			},
			VersionEnclosure{
				ModuleVersions: map[string]Version{
					"a": {
						Module:  "a",
						Version: "v1.0.0",
					},
					"b/v2": {
						Module:  "b/v2",
						Version: "v2.1.3",
					},
				},
			},
			invalidErr,
		},
		"invalid wrong tag": {
			[]string{
				"a/v1.0.0",
				"a/v1.0.1",
				"b/v2.1.3",
			},
			VersionEnclosure{
				ModuleVersions: map[string]Version{
					"a": {
						Module:  "a",
						Version: "v1.0.0",
					},
					"b/v2": {
						Module:  "b/v2",
						Version: "v2.1.3",
					},
				},
			},
			invalidErr,
		},
	}

	for id, tt := range testCases {
		t.Run(id, func(t *testing.T) {
			gitClient := MockGit{
				tags: tt.tags,
			}

			err := tt.enc.isValid(&gitClient)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error to contain %s, got nil err", tt.wantErr)
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected err to contain %s, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestVersionTags(t *testing.T) {
	var testCases = map[string]struct {
		tags       []string
		mod        string
		major      string
		includePre bool
		wantTags   []string
	}{
		"root module": {
			[]string{"v1.0.0", "a/v1.0.0", "a/v1.0.1", "a/v1.1.0", "a/v0.0.0"},
			rootModule,
			"v1",
			false,
			[]string{"v1.0.0"},
		},
		"one module": {
			[]string{"a/v1.0.0", "a/v1.0.1", "a/v1.1.0", "a/v0.0.0"},
			"a",
			"v1",
			false,
			[]string{"v1.1.0", "v1.0.1", "v1.0.0"},
		},
		"one module - v0": {
			[]string{"a/v1.0.0", "a/v1.0.1", "a/v1.1.0", "a/v0.0.0", "a/v0.0.1"},
			"a",
			"v0",
			false,
			[]string{"v0.0.1", "v0.0.0"},
		},
		"two modules": {
			[]string{"b/v0.0.0", "a/v1.0.0", "a/v1.0.1", "a/v1.1.0", "a/v0.0.0"},
			"a",
			"v1",
			false,
			[]string{"v1.1.0", "v1.0.1", "v1.0.0"},
		},
		"v2": {
			[]string{"a/v1.0.0", "b/v1.0.0", "b/v2.0.0", "b/v2.1.3"},
			"b",
			"v2",
			false,
			[]string{"v2.1.3", "v2.0.0"},
		},
		"ignore prereleases": {
			[]string{"a/v1.0.0", "b/v1.0.0", "b/v2.0.0", "b/v2.1.3", "b/v2.1.4-pre"},
			"b",
			"v2",
			false,
			[]string{"v2.1.3", "v2.0.0"},
		},
		"include prereleases": {
			[]string{"a/v1.0.0", "b/v1.0.0", "b/v2.0.0", "b/v2.1.3", "b/v2.1.4-pre"},
			"b",
			"v2",
			true,
			[]string{"v2.1.4-pre", "v2.1.3", "v2.0.0"},
		},
	}

	for id, tt := range testCases {
		t.Run(id, func(t *testing.T) {
			gitClient := MockGit{
				tags: tt.tags,
			}

			tags, err := versionTags(&gitClient, tt.mod, tt.major, tt.includePre)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.wantTags, tags); diff != "" {
				t.Errorf("expect tags to match (-want, +got):\n%v", diff)
			}
		})
	}
}

func TestTagRepo(t *testing.T) {
	var testCases = map[string]struct {
		mod     string
		version string
		wantTag string
		wantErr string
	}{
		"root v0":             {"/", "v0.0.0", "v0.0.0", ""},
		"root v1":             {"/", "v1.0.0", "v1.0.0", ""},
		"root v1 patch":       {"/", "v1.0.4", "v1.0.4", ""},
		"root v2":             {"/v2", "v2.0.0", "v2.0.0", ""},
		"a v0":                {"a", "v0.0.0", "a/v0.0.0", ""},
		"a v1":                {"a", "v1.2.3", "a/v1.2.3", ""},
		"a v2":                {"a/v2", "v2.2.3", "a/v2.2.3", ""},
		"a major mismatch":    {"a/v2", "v3.2.3", "", "version v3.2.3 does not match module a/v2's major version"},
		"root major mismatch": {"/v2", "v3.2.3", "", "version v3.2.3 does not match module /v2's major version"},
		"major mismatch":      {"a/v2", "v3.2.3", "", "version v3.2.3 does not match module a/v2's major version"},
	}

	for id, tt := range testCases {
		t.Run(id, func(t *testing.T) {
			gitClient := MockGit{
				tags: []string{},
			}

			err := tagRepo(&gitClient, "test-release", tt.mod, tt.version)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected err to contain %s, got nil err", tt.wantErr)
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected err to contain %s, got %v", tt.wantErr, err)
				}
			} else {
				if diff := cmp.Diff([]string{tt.wantTag}, gitClient.tags); diff != "" {
					t.Errorf("expect tags to match (-want, +got):\n%v", diff)
				}
			}
		})
	}
}

func TestNextVersion(t *testing.T) {
	var testCases = map[string]struct {
		version     string
		incr        VersionIncrement
		minVersion  string
		wantVersion string
		wantErr     string
	}{
		"none": {
			version:     "v1.0.0",
			incr:        NoBump,
			wantVersion: "v1.0.0",
		},
		"patch": {
			version:     "v1.0.0",
			incr:        PatchBump,
			wantVersion: "v1.0.1",
		},
		"minor": {
			version:     "v1.0.0",
			incr:        MinorBump,
			wantVersion: "v1.1.0",
		},
		"major": {
			version:     "v0.0.4",
			incr:        MajorBump,
			wantVersion: "v1.0.0",
		},
		"v2 none": {
			version:     "v2.2.2",
			incr:        NoBump,
			wantVersion: "v2.2.2",
		},
		"v2 patch": {
			version:     "v2.2.2",
			incr:        PatchBump,
			wantVersion: "v2.2.3",
		},
		"v2 minor": {
			version:     "v2.2.2",
			incr:        MinorBump,
			wantVersion: "v2.3.0",
		},
		"prerelease": {
			version: "v1.0.0-pre",
			incr:    NoBump,
			wantErr: "has a prerelease or build component",
		},
		"bad major": {
			version: "v4.0.0",
			incr:    MajorBump,
			wantErr: "major increment can only be applied to v0",
		},
		"bad version": {
			version: "vx.y.z",
			incr:    MajorBump,
			wantErr: "version vx.y.z is not valid",
		},
		"min version < calculated patch": {
			version:     "v1.2.3",
			incr:        PatchBump,
			minVersion:  "v1.2.2",
			wantVersion: "v1.2.4",
		},
		"min version > calculated patch": {
			version:     "v1.2.3",
			incr:        PatchBump,
			minVersion:  "v1.3.0",
			wantVersion: "v1.3.0",
		},
		"min version < calculated minor": {
			version:     "v1.2.3",
			incr:        MinorBump,
			minVersion:  "v1.2.5",
			wantVersion: "v1.3.0",
		},
		"min version > calculated minor": {
			version:     "v1.2.3",
			incr:        MinorBump,
			minVersion:  "v1.5.0",
			wantVersion: "v1.5.0",
		},
		"min version < calculated major": {
			version:     "v0.2.3",
			incr:        MajorBump,
			minVersion:  "v0.5.0",
			wantVersion: "v1.0.0",
		},
		"min version > calculated major": {
			version:     "v0.2.3",
			incr:        MajorBump,
			minVersion:  "v1.5.0",
			wantVersion: "v1.5.0",
		},
	}

	for id, tt := range testCases {
		t.Run(id, func(t *testing.T) {
			gotVersion, err := nextVersion(tt.version, tt.incr, tt.minVersion)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected err to contain %s, got nil err", tt.wantErr)
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected err to contain %s, got %v", tt.wantErr, err)
				}
			} else {
				if tt.wantVersion != gotVersion {
					t.Errorf("expected version %s, got version %s", tt.wantVersion, gotVersion)
				}
			}
		})
	}
}

func TestVersionEnclosure_Bump(t *testing.T) {
	enc := VersionEnclosure{
		SchemaVersion: SchemaVersion,
		ModuleVersions: map[string]Version{
			"a": {
				Module:     "a",
				ImportPath: "a",
				Version:    "v1.0.0",
			},
			"b": {
				Module:     "b",
				ImportPath: "b",
				Version:    "v2.1.3",
			},
		},
		Packages: map[string]string{},
	}

	wantBump := VersionBump{
		From: "v1.0.0",
		To:   "v1.0.1",
	}

	bump, err := enc.bump("a", PatchBump)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantBump, bump); diff != "" {
		t.Errorf("expect bumps to match (-want, +got):\n%v", diff)
	}
}

type MockGit struct {
	tags []string
}

func (c *MockGit) Tag(tag, message string) error {
	c.tags = append(c.tags, tag)
	return nil
}

func (c *MockGit) Tags(prefix string) ([]string, error) {
	var ret []string

	for _, t := range c.tags {
		if strings.HasPrefix(t, prefix) {
			ret = append(ret, t)
		}
	}

	return ret, nil
}

func (c *MockGit) Commit(unstagedPaths []string, message string) error {
	return nil
}

func (c *MockGit) Push() error {
	return nil
}

func (c *MockGit) CommitHash() (string, error) {
	return "1234567abcde", nil
}
