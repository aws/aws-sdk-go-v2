package changes

import (
	"strconv"
	"testing"
)

func TestRelease_AffectedModules(t *testing.T) {
	var testCases = []struct {
		releaseChanges  []*Change
		expectedModules []string
	}{
		{
			[]*Change{
				&Change{
					Module: "test",
				},
			},
			[]string{"test"},
		},
		{
			[]*Change{
				&Change{
					Module: "test",
				},
				&Change{
					Module: "test",
				},
			},
			[]string{"test"},
		},
		{
			[]*Change{
				&Change{
					Module: "test",
				},
				&Change{
					Module: "other",
				},
			},
			[]string{"test", "other"},
		},
	}

	for i, tt := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
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
