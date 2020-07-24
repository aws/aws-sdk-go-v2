package changes

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestModuleGraph(t *testing.T) {
	repo := getRepository(t)
	oldSdk := sdkRepo
	defer func() { sdkRepo = oldSdk }()
	sdkRepo = strings.TrimSuffix(sdkRepo+"/"+modPrefix, "/modules/")

	repoMods, err := repo.Modules()
	if err != nil {
		t.Fatal(err)
	}

	graph, err := moduleGraph(repo.RootPath, repoMods)
	if err != nil {
		t.Fatal(err)
	}

	wantGraph := ModuleGraph{
		"a": []string{"b"},
	}
	if diff := cmp.Diff(graph, wantGraph); diff != "" {
		t.Errorf("expect dependencies to match:\n%v", diff)
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

			if diff := cmp.Diff(updates, tt.wantUpdates); diff != "" {
				t.Errorf("expect dependencies to match:\n%v", diff)
			}
		})
	}
}
