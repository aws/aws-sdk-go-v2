package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/release"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_executeModuleTemplate(t *testing.T) {
	tests := map[string]struct {
		summary moduleSummary
		wantWr  string
		wantErr bool
	}{
		"no annotations": {
			summary: moduleSummary{
				ReleaseID: "2021-05-05",
				Version:   "v1.0.0",
			},
			wantWr: `# v1.0.0 (2021-05-05)

* No change notes available for this release.

`,
		},

		"annotations": {
			summary: moduleSummary{
				ReleaseID: "2021-05-05",
				Version:   "v1.0.0",
				Annotations: func() (v []changelog.Annotation) {
					v = []changelog.Annotation{
						{
							Type:        changelog.DependencyChangeType,
							Description: "Updated foo to a new version.",
						},
						{
							Type:        changelog.DocumentationChangeType,
							Description: "Fixed a documentation bug.",
						},
						{
							Type:        changelog.BugFixChangeType,
							Description: "Fixed a broken thing.",
						},
						{
							Type:        changelog.FeatureChangeType,
							Description: "A new fancy feature.",
						},
						{
							Type:        changelog.ReleaseChangeType,
							Description: "New stable version.",
						},
						{
							Type:        changelog.AnnouncementChangeType,
							Description: "Something you should know.",
						},
					}
					sortAnnotations(v)
					return v
				}(),
			},
			wantWr: `# v1.0.0 (2021-05-05)

* **Announcement**: Something you should know.
* **Release**: New stable version.
* **Feature**: A new fancy feature.
* **Bug Fix**: Fixed a broken thing.
* **Documentation**: Fixed a documentation bug.
* **Dependency Update**: Updated foo to a new version.

`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			wr := bytes.NewBuffer(nil)
			err := executeModuleTemplate(wr, tt.summary)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeModuleTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotWr := wr.String()
			if diff := cmp.Diff(tt.wantWr, gotWr); len(diff) > 0 {
				t.Errorf(diff)
			}
		})
	}
}

func Test_executeRepoChangeLogEntryTemplate(t *testing.T) {
	tests := map[string]struct {
		summary releaseSummary
		wantWr  string
		wantErr bool
	}{
		"general highlights and module highlights": {
			summary: releaseSummary{
				ReleaseID: "2021-05-05",
				General: func() (v []changelog.Annotation) {
					v = []changelog.Annotation{
						{
							Type:        changelog.FeatureChangeType,
							Collapse:    true,
							Description: "a",
						},
					}
					sortAnnotations(v)
					return v
				}(),
				Modules: map[string]moduleSummary{
					".": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b",
						Version:    "v1.1.0",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.FeatureChangeType,
								Collapse:    true,
								Description: "a",
							},
							{
								Type:        changelog.FeatureChangeType,
								Description: "b",
							},
						},
					},
					"c/d": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b/c/d",
						Version:    "v1.1.0",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.FeatureChangeType,
								Collapse:    true,
								Description: "a",
							},
						},
					},
					"e/f": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b/e/f",
						Version:    "v1.0.1",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.BugFixChangeType,
								Description: "c",
							},
						},
					},
				},
			},
			wantWr: `# Release (2021-05-05)

## General Highlights
* **Feature**: a

## Module Highlights
* ` + "`a/b`" + `: v1.1.0
  * **Feature**: b
* ` + "`a/b/e/f`" + `: [v1.0.1](e/f/CHANGELOG.md#v101-2021-05-05)
  * **Bug Fix**: c

`,
		},
		"general highlights only": {
			summary: releaseSummary{
				ReleaseID: "2021-05-05",
				General: func() (v []changelog.Annotation) {
					v = []changelog.Annotation{
						{
							Type:        changelog.FeatureChangeType,
							Collapse:    true,
							Description: "a",
						},
					}
					sortAnnotations(v)
					return v
				}(),
				Modules: map[string]moduleSummary{
					".": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b",
						Version:    "v1.1.0",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.FeatureChangeType,
								Collapse:    true,
								Description: "a",
							},
						},
					},
				},
			},
			wantWr: `# Release (2021-05-05)

## General Highlights
* **Feature**: a

`,
		},
		"module highlights only": {
			summary: releaseSummary{
				ReleaseID: "2021-05-05",
				Modules: map[string]moduleSummary{
					"c/d": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b/c/d",
						Version:    "v1.1.0",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.FeatureChangeType,
								Description: "a",
							},
						},
					},
				},
			},
			wantWr: `# Release (2021-05-05)

## Module Highlights
* ` + "`a/b/c/d`" + `: [v1.1.0](c/d/CHANGELOG.md#v110-2021-05-05)
  * **Feature**: a

`,
		},
		"no release notes": {
			summary: releaseSummary{
				ReleaseID: "2021-05-05",
				Modules: map[string]moduleSummary{
					"c/d": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b/c/d",
						Version:    "v1.1.0",
					},
				},
			},
			wantWr: `# Release (2021-05-05)

* No change notes available for this release.

`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			wr := &bytes.Buffer{}
			err := executeRepoChangeLogEntryTemplate(wr, tt.summary)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeRepoChangeLogEntryTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.wantWr, wr.String()); len(diff) > 0 {
				t.Error(diff)
			}
		})
	}
}

func Test_executeSummaryNotesTemplate(t *testing.T) {
	tests := map[string]struct {
		summary releaseSummary
		wantWr  string
		wantErr bool
	}{
		"general highlights and module highlights": {
			summary: releaseSummary{
				ReleaseID: "2021-05-05",
				General: func() (v []changelog.Annotation) {
					v = []changelog.Annotation{
						{
							Type:        changelog.FeatureChangeType,
							Collapse:    true,
							Description: "a",
						},
					}
					sortAnnotations(v)
					return v
				}(),
				Modules: map[string]moduleSummary{
					".": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b",
						Version:    "v1.1.0",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.FeatureChangeType,
								Collapse:    true,
								Description: "a",
							},
							{
								Type:        changelog.FeatureChangeType,
								Description: "b",
							},
						},
					},
					"c/d": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b/c/d",
						Version:    "v1.1.0",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.FeatureChangeType,
								Collapse:    true,
								Description: "a",
							},
						},
					},
					"e/f": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b/e/f",
						Version:    "v1.0.1",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.BugFixChangeType,
								Description: "c",
							},
						},
					},
				},
			},
			wantWr: `## General Highlights
* **Feature**: a

## Module Highlights
* ` + "`a/b`" + `: v1.1.0
  * **Feature**: b
* ` + "`a/b/e/f`" + `: [v1.0.1](e/f/CHANGELOG.md#v101-2021-05-05)
  * **Bug Fix**: c

`,
		},
		"general highlights only": {
			summary: releaseSummary{
				ReleaseID: "2021-05-05",
				General: func() (v []changelog.Annotation) {
					v = []changelog.Annotation{
						{
							Type:        changelog.FeatureChangeType,
							Collapse:    true,
							Description: "a",
						},
					}
					sortAnnotations(v)
					return v
				}(),
				Modules: map[string]moduleSummary{
					".": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b",
						Version:    "v1.1.0",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.FeatureChangeType,
								Collapse:    true,
								Description: "a",
							},
						},
					},
				},
			},
			wantWr: `## General Highlights
* **Feature**: a

`,
		},
		"module highlights only": {
			summary: releaseSummary{
				ReleaseID: "2021-05-05",
				Modules: map[string]moduleSummary{
					"c/d": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b/c/d",
						Version:    "v1.1.0",
						Annotations: []changelog.Annotation{
							{
								Type:        changelog.FeatureChangeType,
								Description: "a",
							},
						},
					},
				},
			},
			wantWr: `## Module Highlights
* ` + "`a/b/c/d`" + `: [v1.1.0](c/d/CHANGELOG.md#v110-2021-05-05)
  * **Feature**: a

`,
		},
		"no release notes": {
			summary: releaseSummary{
				ReleaseID: "2021-05-05",
				Modules: map[string]moduleSummary{
					"c/d": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b/c/d",
						Version:    "v1.1.0",
					},
				},
			},
			wantWr: `* No change notes available for this release.

`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			wr := &bytes.Buffer{}
			err := executeSummaryNotesTemplate(wr, tt.summary)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeSummaryNotesTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.wantWr, wr.String()); len(diff) > 0 {
				t.Error(diff)
			}
		})
	}
}

func Test_generateSummary(t *testing.T) {
	tests := map[string]struct {
		manifest    release.Manifest
		annotations []changelog.Annotation
		want        releaseSummary
		wantErr     bool
	}{
		"summary": {
			manifest: release.Manifest{
				ID: "2021-05-05",
				Modules: map[string]release.ModuleManifest{
					".": {
						ModulePath:  "a/b",
						From:        "v1.0.0",
						To:          "v1.0.1",
						Changes:     release.SourceChange,
						Annotations: []string{"A"},
					},
					"c/d": {
						ModulePath:  "a/b/c/d",
						From:        "v1.0.0",
						To:          "v1.1.0",
						Changes:     release.SourceChange | release.DependencyUpdate,
						Annotations: []string{"A", "B"},
					},
				},
			},
			annotations: []changelog.Annotation{
				{
					ID:          "B",
					Collapse:    true,
					Type:        changelog.FeatureChangeType,
					Description: "B feature",
				},
				{
					ID:          "A",
					Type:        changelog.AnnouncementChangeType,
					Description: "An Announcement",
				},
			},
			want: releaseSummary{
				ReleaseID: "2021-05-05",
				General: []changelog.Annotation{
					{
						ID:          "B",
						Collapse:    true,
						Type:        changelog.FeatureChangeType,
						Description: "B feature",
					},
					dependencyBump,
				},
				Modules: map[string]moduleSummary{
					".": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b",
						Version:    "v1.0.1",
						Annotations: []changelog.Annotation{
							{
								ID:          "A",
								Type:        changelog.AnnouncementChangeType,
								Description: "An Announcement",
							},
						},
					},
					"c/d": {
						ReleaseID:  "2021-05-05",
						ModulePath: "a/b/c/d",
						Version:    "v1.1.0",
						Annotations: []changelog.Annotation{
							{
								ID:          "A",
								Type:        changelog.AnnouncementChangeType,
								Description: "An Announcement",
							},
							{
								ID:          "B",
								Collapse:    true,
								Type:        changelog.FeatureChangeType,
								Description: "B feature",
							},
							dependencyBump,
						},
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := generateSummary(tt.manifest, tt.annotations)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); len(diff) > 0 {
				t.Error(diff)
			}
		})
	}
}
