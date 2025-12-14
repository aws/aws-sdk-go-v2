//go:build integration
// +build integration

package transfermanager

import (
	"testing"
)

func TestInteg_UploadDirectory(t *testing.T) {
	cases := map[string]uploadDirectoryTestData{
		"single file": {
			FilesSize: map[string]int64{
				"foo": 2 * 1024 * 1024,
			},
			Source:              "integ-dir",
			Recursive:           true,
			ExpectFilesUploaded: 1,
			ExpectKeys:          []string{"foo"},
		},
		"multi file non-recursive": {
			FilesSize: map[string]int64{
				"foo":        2 * 1024 * 1024,
				"bar":        10 * 1024 * 1024,
				"to/the/baz": 20 * 1024 * 1024,
			},
			Source:              "integ-dir",
			ExpectFilesUploaded: 2,
			ExpectKeys:          []string{"foo", "bar"},
		},
		"multi file recursive with prefix": {
			FilesSize: map[string]int64{
				"foo":        2 * 1024 * 1024,
				"to/bar":     10 * 1024 * 1024,
				"to/the/baz": 20 * 1024 * 1024,
			},
			Source:              "integ-dir",
			Recursive:           true,
			KeyPrefix:           "bla",
			ExpectFilesUploaded: 3,
			ExpectKeys:          []string{"bla/foo", "bla/to/bar", "bla/to/the/baz"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testUploadDirectory(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}
