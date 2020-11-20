// +build integration

package s3

import (
	"bytes"
	"strings"
	"testing"
)

func TestInteg_WriteToObject(t *testing.T) {
	cases := map[string]writeToObjectTestData{
		"seekable body":     {Body: strings.NewReader("hello world"), ExpectBody: []byte("hello world")},
		"empty string body": {Body: strings.NewReader(""), ExpectBody: []byte("")},
		"nil body":          {Body: nil, ExpectBody: []byte("")},
		"unseekable body":   {Body: bytes.NewBuffer([]byte("hello world")), ExpectError: "failed to compute payload hash"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testWriteToObject(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}

func TestInteg_CopyObject(t *testing.T) {
	testCopyObject(t, setupMetadata.Buckets.Source.Name, setupMetadata.Buckets.Target.Name, nil)
}
