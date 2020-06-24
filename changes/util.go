package changes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func writeJSON(data interface{}, root, dir, name string) error {
	filePath := filepath.Join(root, dir, name+".json")
	changeBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return writeFile(changeBytes, filePath, false)
}

func writeFile(data []byte, path string, appendTo bool) error {
	if appendTo {
		exists, err := fileExists(path, false)
		if err != nil {
			return err
		}

		if exists {
			existingData, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			data = append(data, existingData...)
		}
	}

	return ioutil.WriteFile(path, data, 0644)
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
