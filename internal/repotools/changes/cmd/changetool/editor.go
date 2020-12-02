package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const defaultEditor = "vim"

var allowedEditors = map[string]struct{}{
	"vi":    {},
	"vim":   {},
	"gvim":  {},
	"nano":  {},
	"edit":  {},
	"gedit": {},
	"emacs": {},
}

func getEditorTool() (string, error) {
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

func editTemplate(template []byte) ([]byte, error) {
	editor, err := getEditorTool()
	if err != nil {
		return nil, err
	}

	f, err := ioutil.TempFile("", "change-*.yml")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	_, err = f.Write(template)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	filledTemplate, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}

	if bytes.Compare(filledTemplate, template) == 0 {
		return nil, errors.New("template was not modified")
	}

	return filledTemplate, nil
}
