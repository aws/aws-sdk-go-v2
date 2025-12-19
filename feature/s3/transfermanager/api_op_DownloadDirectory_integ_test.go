//go:build integration
// +build integration

package transfermanager

import (
	"testing"
)

func TestInteg_DownloadDirectory(t *testing.T) {
	cases := map[string]downloadDirectoryTestData{
		"multi objects with prefix": {
			ObjectsSize: map[string]int64{
				"oii/bar":     2 * 1024 * 1024,
				"oiibaz/zoo":  10 * 1024 * 1024,
				"oii/baz/zoo": 10 * 1024 * 1024,
				"oi":          20 * 1024 * 1024,
			},
			KeyPrefix:               "oii",
			ExpectObjectsDownloaded: 3,
			ExpectFiles:             []string{"bar", "oiibaz/zoo", "baz/zoo"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testDownloadDirectory(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}
