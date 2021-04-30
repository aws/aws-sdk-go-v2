package gomod

import "testing"

func TestIsModuleChanged(t *testing.T) {
	tests := map[string]struct {
		moduleDir  string
		submodules []string
		changes    []string
		want       bool
		wantErr    bool
	}{
		"no submodules": {
			moduleDir: ".",
			changes: []string{
				"sub3/foo.go",
				"sub2/bar.go",
				"sub1/baz.go",
				"foo.go",
			},
			want: true,
		},
		"no submodules, no go changes": {
			moduleDir: ".",
			changes: []string{
				"foo.java",
			},
			want: false,
		},
		"go.mod considered": {
			moduleDir: ".",
			changes: []string{
				"go.mod",
			},
			want: true,
		},
		"repo root with submodules": {
			moduleDir:  ".",
			submodules: []string{"sub1", "sub2"},
			changes: []string{
				"sub3/foo.go",
				"sub2/bar.go",
				"sub1/baz.go",
				"foo.go",
			},
			want: true,
		},
		"submodule directory, no submodules, no changes": {
			moduleDir:  "sub1",
			submodules: nil,
			changes: []string{
				"sub3/foo.go",
				"sub2/bar.go",
				"foo.go",
			},
		},
		"submodule directory, no submodules, changes": {
			moduleDir:  "sub1",
			submodules: nil,
			changes: []string{
				"sub3/foo.go",
				"sub2/bar.go",
				"sub1/bar.go",
				"foo.go",
			},
			want: true,
		},
		"submodule directory, submodules, no changes": {
			moduleDir:  "sub1",
			submodules: []string{"sub1/subsub1", "sub1/subsub2"},
			changes: []string{
				"sub3/foo.go",
				"sub2/bar.go",
				"sub1/subsub1/foo.go",
				"sub1/subsub1/bar.go",
				"sub1/subsub2/bar.go",
				"foo.go",
			},
		},
		"submodule directory, submodules, changes": {
			moduleDir:  "sub1",
			submodules: []string{"sub1/subsub1", "sub1/subsub2"},
			changes: []string{
				"sub3/foo.go",
				"sub2/bar.go",
				"sub1/subsub1/foo.go",
				"sub1/subsub1/bar.go",
				"sub1/subsub2/bar.go",
				"sub1/notsub/foo.go",
				"foo.go",
			},
			want: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := IsModuleChanged(tt.moduleDir, tt.submodules, tt.changes)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsModuleChanged() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsModuleChanged() got = %v, want %v", got, tt.want)
			}
		})
	}
}
