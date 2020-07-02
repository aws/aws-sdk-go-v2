package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

func editTemplate(template []byte) ([]byte, error) {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")

		if editor == "" {
			editor = "vim"
		}
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
