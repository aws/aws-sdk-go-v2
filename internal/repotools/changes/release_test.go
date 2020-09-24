package changes

import (
	"testing"
)

func TestRelease_AffectedModules(t *testing.T) {
	var testCases = map[string]struct {
		releaseChanges  []Change
		expectedModules []string
	}{
		"single module": {
			[]Change{
				{
					Module: "test",
					Type:   FeatureChangeType,
				},
			},
			[]string{"test"},
		},
		"two changes": {
			[]Change{
				{
					Module: "test",
					Type:   FeatureChangeType,
				},
				{
					Module: "test",
					Type:   FeatureChangeType,
				},
			},
			[]string{"test"},
		},
		"two modules": {
			[]Change{
				{
					Module: "test",
					Type:   FeatureChangeType,
				},
				{
					Module: "other",
					Type:   BugFixChangeType,
				},
			},
			[]string{"test", "other"},
		},
	}

	for id, tt := range testCases {
		t.Run(id, func(t *testing.T) {
			r := &Release{
				Changes: tt.releaseChanges,
			}

			modules := r.AffectedModules()
			for _, expectedMod := range tt.expectedModules {
				found := false
				for _, m := range modules {
					if m == expectedMod {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("expected modules to contain %s", expectedMod)
				}
			}
		})
	}
}
