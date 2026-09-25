//go:build integration

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
		// The directory path reaches the code that must not assume all parts
		// share the first part's size (#3526) now that it calls DownloadObject.
		// All parts but the last are >= 5MB to satisfy the S3 minimum part size.
		"unequal part sizes": {
			ObjectsSize: map[string]int64{
				"unequal/multipart":  12 * 1024 * 1024,
				"unequal/singlepart": 1 * 1024 * 1024,
			},
			PartSizes: map[string][]int64{
				"unequal/multipart": {6 * 1024 * 1024, 5 * 1024 * 1024, 1 * 1024 * 1024},
			},
			KeyPrefix:               "unequal",
			ExpectObjectsDownloaded: 2,
			ExpectFiles:             []string{"multipart", "singlepart"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testDownloadDirectory(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}
