//go:build gofuzzbeta
// +build gofuzzbeta

package ini

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func FuzzParse(f *testing.F) {
	corpus, err := loadCorpus(
		// Also loads "testdata/corpus/FuzzParse" automatically,
		// for previous known panics if present and contains any files
		filepath.Join("testdata", "valid", "*"),
		filepath.Join("testdata", "invalid", "*"),
	)
	if err != nil {
		f.Fatalf("failed to load corpus, %v", err)
	}

	for name, c := range corpus {
		f.Add(name, c)
	}

	f.Fuzz(func(t *testing.T, name string, c []byte) {
		_, err := Parse(bytes.NewReader(c), name)
		if err != nil {
			t.Logf("parse failed for %v, %v", name, err.Error())
		}
	})
}

func loadCorpus(globs ...string) (map[string][]byte, error) {
	corpus := map[string][]byte{}
	for _, g := range globs {
		paths, err := filepath.Glob(g)
		if err != nil {
			return nil, fmt.Errorf("unable to glob corpus, %w", err)
		}
		if len(paths) == 0 {
			return nil, fmt.Errorf("no test corpus files found for glob %v", g)
		}

		for _, p := range paths {
			c, err := ioutil.ReadFile(p)
			if err != nil {
				return nil, fmt.Errorf("failed to read corpus file, %w", err)
			}
			corpus[p] = c
		}
	}

	return corpus, nil
}
