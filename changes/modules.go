package changes

import (
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"strings"
)

func GetCurrentModule() (string, error) {
	path, err := findFile("go.mod", false)
	if err != nil {
		return "", err
	}

	modFile, err := getModFile(path)
	if err != nil {
		return "", err
	}

	return shortenModPath(modFile.Module.Mod.String()), nil
}

func getModFile(path string) (*modfile.File, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return modfile.ParseLax(path, data, nil)
}

func shortenModPath(modulePath string) string {
	return strings.TrimPrefix(modulePath, "github.com/aggagen/")
}
