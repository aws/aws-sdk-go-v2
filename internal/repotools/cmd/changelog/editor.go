package main

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
)

func editTemplate(template []byte) ([]byte, error) {
	editor, err := repotools.GetEditorTool()
	if err != nil {
		return nil, err
	}

	f, err := ioutil.TempFile("", "changelog-*.toml")
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

	return filledTemplate, nil
}
