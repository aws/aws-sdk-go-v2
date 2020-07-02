package changes

import (
	"bytes"
	"strconv"
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
	var changeTests = []struct {
		modules     []string
		changeType  ChangeType
		description string
		wantErr     bool
	}{
		{[]string{"a"}, FeatureChangeType, "this is a description", false},
		{[]string{"a", "b"}, FeatureChangeType, "this is a description", false},
		{[]string{"a", "b"}, BugFixChangeType, "this is a description", false},
		{[]string{"a", "b"}, BugFixChangeType, "", true},
		{[]string{}, FeatureChangeType, "this is a description", true},
	}

	for i, tt := range changeTests {
		t.Run("NewChangesCase_"+strconv.Itoa(i), func(t *testing.T) {
			changes, err := NewChanges(tt.modules, tt.changeType, tt.description)
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
					assertChangeEqual(t, &Change{
						ID:          c.ID,
						Module:      c.Module,
						Type:        tt.changeType,
						Description: tt.description,
					}, c)
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

# type may be one of "feature" or "bugfix".
# multiple modules may be listed. A change metadata file will be created for each module.`

	template, err := ChangeToTemplate(&Change{
		ID:          "test-feature-1",
		Module:      "test",
		Type:        FeatureChangeType,
		Description: "test description",
	})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	if bytes.Compare(template, []byte(wantTemplate)) != 0 {
		t.Errorf("expected template \"%s\", got \"%s\"", string(wantTemplate), string(template))
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

	assertChangeEqual(t, &Change{
		Module:      "test",
		Type:        FeatureChangeType,
		Description: "test description",
	}, change)
}

func assertChangeEqual(t *testing.T, want, got *Change) bool {
	if !changeEquals(t, want, got) {
		t.Errorf("expected change like %v, got %v", want, got)
		return false
	}

	return true
}

func changeEquals(t *testing.T, want, got *Change) bool {
	t.Helper()

	if want.Module != "" && want.Module != got.Module {
		return false
	}
	if want.Type != "" && want.Type != got.Type {
		return false
	}
	if want.Description != "" && want.Description != got.Description {
		return false
	}
	if want.ID != "" && want.ID != got.ID {
		return false
	}

	return true
}

func assertChangesHas(t *testing.T, changes []*Change, want *Change) bool {
	for _, c := range changes {
		if changeEquals(t, want, c) {
			return true
		}
	}

	t.Errorf("expected changes to contain %v", want)
	return true
}
