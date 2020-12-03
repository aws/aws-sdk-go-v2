package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WriteJSON marshals the given data as JSON and writes it to the specified location.
func WriteJSON(data interface{}, root, dir, name string) error {
	filePath := filepath.Join(root, dir, name+".json")
	changeBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return WriteFile(changeBytes, filePath, false)
}

// WriteFile writes the given bytes to the file at path, creating one if necessary. If appendTo is true, then the data
// will be prepended to the top of the file.
func WriteFile(data []byte, path string, appendTo bool) error {
	if appendTo {
		exists, err := FileExists(path, false)
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

// FileExists returns whether or not a file exists at the specified path. If dir is true, FileExists checks for the existence
// of the specified directory.
func FileExists(path string, dir bool) (bool, error) {
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

// FindFile recursively searches upwards from the current directory to the filesystem root for the specified file.
func FindFile(fileName string, dir bool) (string, error) {
	currPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to find file: %v", err)
	}

	for {
		if currPath == string(os.PathSeparator) || filepath.VolumeName(currPath) == currPath {
			return "", errors.New("failed to find file: reached filesystem root")
		}

		targetFilepath := filepath.Join(currPath, fileName)
		found, err := FileExists(targetFilepath, dir)
		if err != nil {
			return "", fmt.Errorf("failed to find file: %v", err)
		}

		if found {
			return targetFilepath, nil
		}

		// trimming trailing '/' causes filepath.Split to trim the last directory in currPath
		currPath = strings.TrimSuffix(currPath, string(os.PathSeparator))
		currPath, _ = filepath.Split(currPath)
	}
}

// ReplaceLine replaces any line in the file at the given filename that begins
// with linePrefix with the given replacement string.
func ReplaceLine(filename, linePrefix, replacement string) (err error) {
	var f *os.File
	f, err = os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open file %v", err)
	}
	defer func() {
		cErr := f.Close()
		if err == nil && cErr != nil {
			err = fmt.Errorf("failed to close file, %w", cErr)
		}
	}()

	var buff bytes.Buffer
	scanner := bufio.NewScanner(ioutil.NopCloser(f))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, linePrefix) {
			line = replacement
		}

		buff.WriteString(line)
		buff.WriteRune('\n')
	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file lines, %w", err)
	}

	if err = f.Truncate(0); err != nil {
		return fmt.Errorf("failed to reset file, %w", err)
	}
	if _, err = f.Seek(0, os.SEEK_SET); err != nil {
		return fmt.Errorf("failed to seek file, %w", err)
	}

	if _, err = io.Copy(f, &buff); err != nil {
		return fmt.Errorf("failed to update file, %w", err)
	}

	return nil
}

// ExecAt runs the given Cmd with is working directory set to path.
func ExecAt(cmd *exec.Cmd, path string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Dir = path

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("couldn't run cmd %s: %v: %s", cmd.String(), err, stderr.String())
	}

	return stdout.Bytes(), nil
}
