package repotools

import (
	"fmt"
	"os"
)

const defaultEditor = "vim"

var allowedEditors = map[string]struct{}{
	"vi":    {},
	"vim":   {},
	"gvim":  {},
	"nvim":  {},
	"nano":  {},
	"edit":  {},
	"gedit": {},
	"emacs": {},
}

// GetEditorTool returns the editor tool to use for interactive file edits.
func GetEditorTool() (string, error) {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")

		if editor == "" {
			editor = defaultEditor
		}
	}

	if _, ok := allowedEditors[editor]; !ok {
		return "", fmt.Errorf("unknown editor %q not allowed, %v", editor, allowedEditors)
	}

	return editor, nil
}
