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
		"multi file with prefix and custom delimiter": {
			ObjectsSize: map[string]int64{
				"yee#bar":        2 * 1024 * 1024,
				"yee#baz#":       0,
				"yee#baz#zoo":    10 * 1024 * 1024,
				"yee#oii@zoo":    10 * 1024 * 1024,
				"yee#yee#..#bla": 2 * 1024 * 1024,
				"ye":             20 * 1024 * 1024,
			},
			KeyPrefix:               "yee#",
			Delimiter:               "#",
			ExpectObjectsDownloaded: 4,
			ExpectFiles:             []string{"bar", "baz/zoo", "oii@zoo", "bla"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testDownloadDirectory(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}
