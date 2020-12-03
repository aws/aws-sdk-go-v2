package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func newTestFile(t *testing.T, content []byte) (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("failed to open test file, %w", err)
	}
	filename := f.Name()

	t.Cleanup(func() {
		if err := os.Remove(filename); err != nil {
			t.Errorf("failed to cleanup test file, %v", err)
		}
	})

	if _, err = f.Write(content); err != nil {
		return "", fmt.Errorf("failed to write test file, %w", err)
	}

	if err = f.Close(); err != nil {
		return "", fmt.Errorf("failed to close test file, %w", err)
	}

	return filename, nil
}

func TestReplaceLines(t *testing.T) {
	cases := map[string]struct {
		Setup   func(*testing.T) (string, error)
		Prefix  string
		Replace string

		ExpectErr     string
		ExpectContent string
	}{
		"file not exists": {
			ExpectErr: "failed to open",
			Setup: func(t *testing.T) (string, error) {
				return filepath.Join("testdata", "not_exists"), nil
			},
		},

		"replace lines": {
			Prefix:  "foo - ",
			Replace: "bar - content",
			Setup: func(t *testing.T) (string, error) {
				content := `
content
foo - something
else
`
				return newTestFile(t, []byte(content))
			},
			ExpectContent: `
content
bar - content
else
`,
		},
		"no change": {
			Prefix:  "not found - ",
			Replace: "bar - content",
			Setup: func(t *testing.T) (string, error) {
				content := `
content
foo - something
else
`
				return newTestFile(t, []byte(content))
			},
			ExpectContent: `
content
foo - something
else
`,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			filename, err := c.Setup(t)
			if err != nil {
				t.Fatalf("failed to setup test case, %v", err)
			}

			err = ReplaceLine(filename, c.Prefix, c.Replace)
			if len(c.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect error to contain %v, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			actual, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatalf("failed to reopen test file, %v", err)
			}

			if diff := cmp.Diff(c.ExpectContent, string(actual)); len(diff) != 0 {
				t.Errorf("expect match\n%s", diff)
			}
		})
	}
}
