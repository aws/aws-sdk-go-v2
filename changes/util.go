package changes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type VersionedSchema interface {
	SetSchemaVersion(string)
}

func writeFile(data VersionedSchema, root, dir, name string) error {
	data.SetSchemaVersion(SchemaVersion)

	filePath := filepath.Join(root, dir, name+".json")
	changeBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, changeBytes, 0644)
}

func fileExists(path string, dir bool) (bool, error) {
	if f, err := os.Stat(path); err == nil {
		if f.IsDir() != dir {
			return false, nil
		}

		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func findFile(fileName string, dir bool) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pathParts := strings.Split(path, string(os.PathSeparator))
	prefix := ""
	if strings.HasPrefix(path, string(os.PathSeparator)) {
		prefix = string(os.PathSeparator)
	}

	for i := len(pathParts); i > 0; i-- {
		path = prefix + filepath.Join(append(pathParts[:i], fileName)...)

		found, err := fileExists(path, dir)
		if err != nil {
			return "", err
		}

		if found {
			return path, nil
		}
	}

	return "", fmt.Errorf("couldn't find %s in current or parent directories", fileName)
}
