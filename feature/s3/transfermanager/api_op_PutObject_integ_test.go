//go:build integration
// +build integration

package transfermanager

import (
	"bytes"
	"strings"
	"testing"
)

func TestInteg_PutObject(t *testing.T) {
	cases := map[string]putObjectTestData{
		"seekable body":         {Body: strings.NewReader("hello world"), ExpectBody: []byte("hello world")},
		"empty string body":     {Body: strings.NewReader(""), ExpectBody: []byte("")},
		"multipart upload body": {Body: bytes.NewReader(largeObjectBuf), ExpectBody: largeObjectBuf},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testPutObject(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}
