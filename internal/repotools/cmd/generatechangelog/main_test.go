package main

import (
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/release"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_sortAnnotations(t *testing.T) {
	annotations := []changelog.Annotation{
		{Type: changelog.DependencyChangeType, Description: "b"},
		{Type: changelog.DependencyChangeType, Description: "a"},
		{Type: changelog.DocumentationChangeType, Description: "c"},
		{Type: changelog.BugFixChangeType, Description: "d"},
		{Type: changelog.FeatureChangeType, Description: "e"},
		{Type: changelog.ReleaseChangeType, Description: "f"},
		{Type: changelog.AnnouncementChangeType, Description: "g"},
	}

	want := []changelog.Annotation{
		{Type: changelog.AnnouncementChangeType, Description: "g"},
		{Type: changelog.ReleaseChangeType, Description: "f"},
		{Type: changelog.FeatureChangeType, Description: "e"},
		{Type: changelog.BugFixChangeType, Description: "d"},
		{Type: changelog.DocumentationChangeType, Description: "c"},
		{Type: changelog.DependencyChangeType, Description: "a"},
		{Type: changelog.DependencyChangeType, Description: "b"},
	}

	sortAnnotations(annotations)

	if diff := cmp.Diff(annotations, want); len(diff) > 0 {
		t.Error(diff)
	}
}

func Test_filterUnreferencedAnnotations(t *testing.T) {
	manifest := release.Manifest{
		Modules: map[string]release.ModuleManifest{
			"foo": {
				Annotations: []string{"a", "c"},
			},
			"bar": {
				Annotations: []string{"c"},
			},
		},
	}

	annotations := []changelog.Annotation{{ID: "a"}, {ID: "b"}, {ID: "c"}}
	wantFiltered := []changelog.Annotation{{ID: "a"}, {ID: "c"}}

	gotFiltered := filterUnreferencedAnnotations(manifest, annotations)

	if diff := cmp.Diff(wantFiltered, gotFiltered); len(diff) > 0 {
		t.Error(diff)
	}
}
