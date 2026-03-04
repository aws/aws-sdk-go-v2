//go:build !windows
// +build !windows

package shareddefaults_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/shareddefaults"
)

func TestExpandHomePath(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("HOME", "/home/user")

	cases := map[string]struct {
		Input  string
		Expect string
	}{
		"empty": {
			Input:  "",
			Expect: "",
		},
		"absolute path unchanged": {
			Input:  "/absolute/path/to/config",
			Expect: "/absolute/path/to/config",
		},
		"relative path unchanged": {
			Input:  "relative/path/to/config",
			Expect: "relative/path/to/config",
		},
		"tilde only": {
			Input:  "~",
			Expect: "/home/user",
		},
		"tilde with slash": {
			Input:  "~/.aws/config",
			Expect: filepath.Join("/home/user", ".aws", "config"),
		},
		"tilde with nested path": {
			Input:  "~/projects/my app/.aws/config",
			Expect: filepath.Join("/home/user", "projects", "my app", ".aws", "config"),
		},
		"tilde not at start unchanged": {
			Input:  "/some/~/path",
			Expect: "/some/~/path",
		},
		"tilde without slash unchanged": {
			Input:  "~username/.aws/config",
			Expect: "~username/.aws/config",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := shareddefaults.ExpandHomePath(c.Input)
			if c.Expect != actual {
				t.Errorf("expect %q, got %q", c.Expect, actual)
			}
		})
	}
}
