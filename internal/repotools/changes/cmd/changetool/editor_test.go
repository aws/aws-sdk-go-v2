package main

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
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
			origEnv := awstesting.StashEnv()
			defer awstesting.PopEnv(origEnv)

			if c.SetupEnv != nil {
				c.SetupEnv()
			}

			editor, err := getEditorTool()
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
