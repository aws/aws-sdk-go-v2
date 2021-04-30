package repotools

import (
	"os"
	"strings"
	"testing"
)

func TestGetEditorTool(t *testing.T) {
	cases := map[string]struct {
		SetupEnv     func()
		ExpectEditor string
		ExpectErr    string
	}{
		"default": {
			ExpectEditor: defaultEditor,
		},
		"allowed visual": {
			SetupEnv: func() {
				os.Setenv("VISUAL", `vi`)
			},
			ExpectEditor: `vi`,
		},
		"allowed editor": {
			SetupEnv: func() {
				os.Setenv("EDITOR", `emacs`)
			},
			ExpectEditor: `emacs`,
		},
		"unknown visual": {
			SetupEnv: func() {
				os.Setenv("VISUAL", `unknownCmd`)
			},
			ExpectErr: `unknown editor "unknownCmd"`,
		},
		"unknown editor": {
			SetupEnv: func() {
				os.Setenv("EDITOR", `unknownCmd`)
			},
			ExpectErr: `unknown editor "unknownCmd"`,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			origEnv := os.Environ()
			os.Clearenv()
			defer func(env []string) {
				for _, kv := range env {
					n := strings.SplitN(kv, "=", 2)
					os.Setenv(n[0], n[1])
				}
			}(origEnv)

			if c.SetupEnv != nil {
				c.SetupEnv()
			}

			editor, err := GetEditorTool()
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

			if e, a := c.ExpectEditor, editor; e != a {
				t.Errorf("expect %v editor, got %v", e, a)
			}
		})
	}
}
