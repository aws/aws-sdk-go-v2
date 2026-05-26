//go:build integration

package transfermanager

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
)

func TestInteg_PutObject(t *testing.T) {
	cases := map[string]putObjectTestData{
		"seekable body":         {Body: strings.NewReader("hello world"), ExpectBody: []byte("hello world")},
		"empty string body":     {Body: strings.NewReader(""), ExpectBody: []byte("")},
		"multipart upload body": {Body: bytes.NewReader(largeObjectBuf), ExpectBody: largeObjectBuf},
		"multipart upload body with full object checksum type": {
			Body:              bytes.NewReader(largeObjectBuf),
			ExpectBody:        largeObjectBuf,
			ChecksumAlgorithm: types.ChecksumAlgorithmCrc32c, // only CRC algorithms support full object checksum
			ChecksumType:      types.ChecksumTypeFullObject,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testPutObject(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}
