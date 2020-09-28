package changes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestModuleGraph(t *testing.T) {
	var cases = map[string]struct {
		dependencies map[string][]string
		modules      []string
		wantGraph    ModuleGraph
	}{
		"simple": {
			map[string][]string{
				"a": {"b"},
			},
			[]string{"a", "b"},
			ModuleGraph{
				"b": []string{"a"},
			},
		},
		"cyclical": {
			map[string][]string{
				"a": {"b"},
				"b": {"a"},
			},
			[]string{"a", "b"},
			ModuleGraph{
				"a": []string{"b"},
				"b": []string{"a"},
			},
		},
		"chain": {
			map[string][]string{
				"a": {"b"},
				"b": {"c"},
				"c": {"d"},
			},
			[]string{"a", "b", "c", "d"},
			ModuleGraph{
				"d": []string{"c"},
				"c": []string{"b"},
				"b": []string{"a"},
			},
		},
		"two groups": {
			map[string][]string{
				"a": {"b"},
				"c": {"d"},
			},
			[]string{"a", "b", "c", "d"},
			ModuleGraph{
				"b": []string{"a"},
				"d": []string{"c"},
			},
		},
	}

	for id, tt := range cases {
		t.Run(id, func(t *testing.T) {
			goClient := mockGolist{dependencies: tt.dependencies}

			graph, err := moduleGraph(&goClient, tt.modules)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.wantGraph, graph); diff != "" {
				t.Errorf("expect dependencies to match (-want, +got):\n%v", diff)
			}
		})
	}
}

func TestModuleGraph_DependencyUpdates(t *testing.T) {
	var cases = map[string]struct {
		modGraph       ModuleGraph
		updatedModules []string
		wantUpdates    map[string][]string
	}{
		"no dependencies": {
			ModuleGraph{},
			[]string{"a", "b"},
			map[string][]string{},
		},
		"two modules": {
			ModuleGraph{
				"a": []string{"b"},
			},
			[]string{"a"},
			map[string][]string{
				"b": {"a"},
			},
		},
		"four module chain": {
			ModuleGraph{
				"a": []string{"b"},
				"b": []string{"c"},
				"c": []string{"d"},
			},
			[]string{"a"},
			map[string][]string{
				"b": {"a"},
				"c": {"b"},
				"d": {"c"},
			},
		},
		"cyclic dependency": {
			ModuleGraph{
				"a": []string{"b"},
				"b": []string{"a"},
			},
			[]string{"a"},
			map[string][]string{
				"a": {"b"},
				"b": {"a"},
			},
		},
		"multiple dependency updates": {
			ModuleGraph{
				"a": []string{"c"}, // c depends on both a and b
				"b": []string{"c"},
			},
			[]string{"a", "b"},
			map[string][]string{
				"c": {"a", "b"},
			},
		},
	}

	for id, tt := range cases {
		t.Run(id, func(t *testing.T) {
			updates := tt.modGraph.dependencyUpdates(tt.updatedModules)

			if diff := cmp.Diff(tt.wantUpdates, updates); diff != "" {
				t.Errorf("expect dependencies to match (-want, +got):\n%v", diff)
			}
		})
	}
}
