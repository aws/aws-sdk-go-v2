package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

func editTemplate(template string) (string, error) {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = "vim"
	}

	f, err := ioutil.TempFile(".", "tmp-template-entry")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	_, err = f.Write([]byte(template))
	if err != nil {
		return "", err
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	filledTemplate, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", err
	}

	if string(filledTemplate) == template {
		return "", errors.New("template was not modified")
	}

	return string(filledTemplate), nil
}
