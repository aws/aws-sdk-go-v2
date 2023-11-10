package ini

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidDataFiles(t *testing.T) {
	const expectedFileSuffix = "_expected"
	err := filepath.Walk(filepath.Join("testdata", "valid"),
		func(path string, info os.FileInfo, fnErr error) (err error) {
			if strings.HasSuffix(path, expectedFileSuffix) {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				t.Errorf("%s: unexpected error, %v", path, err)
			}

			defer func() {
				closeErr := f.Close()
				if err == nil {
					err = closeErr
				} else if closeErr != nil {
					err = fmt.Errorf("file close error: %v, original error: %w", closeErr, err)
				}
			}()

			v, err := Parse(f, path)
			if err != nil {
				t.Errorf("%s: unexpected parse error, %v", path, err)
			}

			expectedPath := path + "_expected"
			e := map[string]interface{}{}

			b, err := ioutil.ReadFile(expectedPath)
			if err != nil {
				// ignore files that do not have an expected file
				return nil
			}

			err = json.Unmarshal(b, &e)
			if err != nil {
				t.Errorf("unexpected error during deserialization, %v", err)
			}

			for profile, tableIface := range e {
				p, ok := v.GetSection(profile)
				if !ok {
					t.Fatal("could not find profile " + profile)
				}

				table := tableIface.(map[string]interface{})
				for k, v := range table {
					switch e := v.(type) {
					case string:
						var a string
						if p.values[k].mp != nil {
							a = fmt.Sprintf("%v", p.values[k].mp)
						} else {
							a = p.values[k].str
						}
						if e != a {
							t.Errorf("%s: expected %v, but received %v for profile %v", path, e, a, profile)
						}
					default:
						t.Errorf("unexpected type: %T", e)
					}
				}
			}

			return nil
		})
	if err != nil {
		t.Fatalf("Error while walking the file tree rooted at root, %d", err)
	}
}
