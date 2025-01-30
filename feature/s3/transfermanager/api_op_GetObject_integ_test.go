//go:build integration
// +build integration

package transfermanager

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestInteg_GetObject(t *testing.T) {
	cases := map[string]getObjectTestData{
		"seekable body":           {Body: strings.NewReader("hello world"), ExpectBody: []byte("hello world")},
		"empty string body":       {Body: strings.NewReader(""), ExpectBody: []byte("")},
		"multipart download body": {Body: bytes.NewReader(largeObjectBuf), ExpectBody: largeObjectBuf},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testDownloadObject(t, setupMetadata.Buckets.Source.Name, c)
			c.Body.(io.Seeker).Seek(0, io.SeekStart)
			testGetObject(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}
