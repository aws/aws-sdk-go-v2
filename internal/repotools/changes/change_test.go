package changes

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestParseChangeType(t *testing.T) {
	var testCases = map[string]struct {
		input    string
		wantType ChangeType
		wantErr  string
	}{
		"feature":       {"feature", FeatureChangeType, ""},
		"feature-case":  {"FEATURE", FeatureChangeType, ""},
		"bugfix":        {"bugfix", BugFixChangeType, ""},
		"bugfix-case":   {"BugFix", BugFixChangeType, ""},
		"major":         {"major", MajorChangeType, ""},
		"major-case":    {"mAjOr", MajorChangeType, ""},
		"invalid":       {"not-a-type", "", "unknown change type: not-a-type"},
		"invalid-empty": {"", "", "unknown change type:"},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := ParseChangeType(tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected non-nil err, got nil")
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected err to contain %s, got %s", tt.wantErr, err.Error())
				}
			} else {
				if c != tt.wantType {
					t.Errorf("expected type %s, got %s", tt.wantType, c)
				}
			}
		})
	}
}

func TestNewChanges(t *testing.T) {
	var changeTests = map[string]struct {
		modules     []string
		changeType  ChangeType
		description string
		minVersion  string
		wantErr     bool
	}{
		"valid feature 1 module": {
			modules:     []string{"a"},
			changeType:  FeatureChangeType,
			description: "this is a description",
		},
		"valid feature 2 modules": {
			modules:     []string{"a", "b"},
			changeType:  FeatureChangeType,
			description: "this is a description",
		},
		"valid bugfix 2 modules": {
			modules:     []string{"a", "b"},
			changeType:  BugFixChangeType,
			description: "this is a description",
		},
		"valid module with min version": {
			modules:     []string{"a"},
			changeType:  FeatureChangeType,
			description: "this is a description",
			minVersion:  "v0.2.0",
		},
		"invalid module with min version": {
			modules:     []string{"a/v2"},
			changeType:  FeatureChangeType,
			description: "this is a description",
			minVersion:  "v0.2.0",
			wantErr:     true,
		},
		"invalid missing description": {
			modules:    []string{"a", "b"},
			changeType: BugFixChangeType,
			wantErr:    true,
		},
		"invalid missing modules": {
			modules:     []string{},
			changeType:  FeatureChangeType,
			description: "this is a description",
			wantErr:     true,
		},
	}

	for name, tt := range changeTests {
		t.Run(name, func(t *testing.T) {
			changes, err := NewChanges(tt.modules, tt.changeType, tt.description, tt.minVersion)
			if err != nil && !tt.wantErr {
				t.Errorf("expected nil err, got %v", err)
			} else if err == nil {
				if tt.wantErr {
					t.Errorf("expected non-nil err, got nil")
				}

				if len(changes) != len(tt.modules) {
					t.Errorf("expected %d changes, got %d", len(tt.modules), len(changes))
				}

				for _, c := range changes {
					want := Change{
						ID:            c.ID,
						SchemaVersion: SchemaVersion,
						Module:        c.Module,
						Type:          tt.changeType,
						Description:   tt.description,
						MinVersion:    tt.minVersion,
					}

					if diff := cmp.Diff(want, c); diff != "" {
						t.Errorf("expect changes to match (-want, +got):\n%v", diff)
					}
				}
			}
		})
	}
}

func TestChangeToTemplate(t *testing.T) {
	const wantTemplate = `modules:
- test
type: feature
description: test description

# type may be one of "feature", "bugfix", "announcement", "dependency", or "major".
# multiple modules may be listed. A change metadata file will be created for each module.

# affected_modules should not be provided unless you are creating a wildcard change (by passing
# the wildcard and module flag to the add command).`

	template, err := ChangeToTemplate(Change{
		ID:          "test-feature-1",
		Module:      "test",
		Type:        FeatureChangeType,
		Description: "test description",
	})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	if diff := cmp.Diff([]byte(wantTemplate), template); len(diff) != 0 {
		t.Errorf("expect templates to match (-want, +got):\n%v", diff)
	}
}

func TestTemplateToChanges(t *testing.T) {
	const template = `modules:
- test
type: feature
description: test description

# type may be one of "feature" or "bugfix".
# multiple modules may be listed. A change metadata file will be created for each module.`

	changes, err := TemplateToChanges([]byte(template))
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}

	change := changes[0]

	want := Change{
		ID:            change.ID,
		SchemaVersion: SchemaVersion,
		Module:        "test",
		Type:          FeatureChangeType,
		Description:   "test description",
	}

	if diff := cmp.Diff(want, change); diff != "" {
		t.Errorf("expect changes to match (-want, +got):\n%v", diff)
	}
}

// assertHasChangeLike asserts that the given changes contains a change with the same type and description as want.
func assertHasChangeLike(t *testing.T, changes []Change, want Change) bool {
	want.SchemaVersion = SchemaVersion
	for _, c := range changes {
		want.ID = c.ID

		if diff := cmp.Diff(want, c); diff == "" {
			return true
		}
	}

	want.ID = ""
	t.Errorf("expected changes to contain %v", want)
	return true
}
