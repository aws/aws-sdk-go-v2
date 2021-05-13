package gomod

import (
	"strconv"
	"testing"
)

func TestIsSubmodulePath(t *testing.T) {
	tests := []struct {
		path       string
		submodules []string
		want       bool
	}{
		{
			path: "foo",
			want: false,
		},
		{
			path: "a/b",
			submodules: []string{
				"a",
				"b",
			},
			want: true,
		},
		{
			path: "b",
			submodules: []string{
				"a",
				"c",
			},
			want: false,
		},
		{
			path: "a",
			submodules: []string{
				"b",
				"c",
			},
			want: false,
		},
		{
			path: "a/b",
			submodules: []string{
				"a/b/c",
				"c",
			},
			want: false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := IsSubmodulePath(tt.path, tt.submodules); got != tt.want {
				t.Errorf("IsSubmodulePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
