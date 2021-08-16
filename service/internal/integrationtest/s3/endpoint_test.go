//go:build integration
// +build integration

package s3

import (
	"strings"
	"testing"
)

func TestInteg_AccessPoint_WriteToObject(t *testing.T) {
	cases := map[string]writeToObjectTestData{
		"seekable body": {Body: strings.NewReader("hello world"), ExpectBody: []byte("hello world")},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testWriteToObject(t, setupMetadata.AccessPoints.Source.ARN, c)
		})
	}
}
