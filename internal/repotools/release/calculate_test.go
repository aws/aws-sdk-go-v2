package release

import (
	"testing"
)

func Test_isModuleCarvedOut1(t *testing.T) {
	tests := map[string]struct {
		files      []string
		subModules []string
		want       bool
		wantErr    bool
	}{
		"no submodules, has go.mod, has go source": {
			files: []string{
				"a/go.mod",
				"a/foo.go",
			},
			want: false,
		},
		"no submodules, no go.mod, has go source": {
			files: []string{
				"a/foo.go",
			},
			want: true,
		},
		"no submodules, no files": {
			want: false,
		},
		"submodules, no go.mod, no go source": {
			files: []string{
				"a/b/go.mod",
				"a/b/foo.go",
				"a/c/go.mod",
				"a/c/bar.go",
			},
			subModules: []string{"a/b", "a/c"},
			want:       false,
		},
		"submodules, has go.mod, no go source": {
			files: []string{
				"a/b/go.mod",
				"a/b/foo.go",
				"a/c/go.mod",
				"a/c/bar.go",
				"a/go.mod",
			},
			subModules: []string{"a/b", "a/c"},
			want:       false,
		},
		"submodules, no go.mod, has go source": {
			files: []string{
				"a/b/go.mod",
				"a/b/foo.go",
				"a/c/go.mod",
				"a/c/bar.go",
				"a/foo.go",
			},
			subModules: []string{"a/b", "a/c"},
			want:       true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := isModuleCarvedOut(tt.files, tt.subModules)
			if (err != nil) != tt.wantErr {
				t.Errorf("isModuleCarvedOut() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isModuleCarvedOut() got = %v, want %v", got, tt.want)
			}
		})
	}
}
