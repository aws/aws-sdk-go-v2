//go:build integration

package transfermanager

import (
	"bytes"
	"os"
	"path/filepath"
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
			ChecksumAlgorithm: types.ChecksumAlgorithmCrc32c,
			ChecksumType:      types.ChecksumTypeFullObject,
		},
		"multipart upload file body": {Body: largeObjectFile(t), ExpectBody: largeObjectBuf},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testPutObject(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}

// largeObjectFile returns largeObjectBuf as a temp file opened for reading.
func largeObjectFile(t *testing.T) *os.File {
	t.Helper()

	path := filepath.Join(t.TempDir(), "large-object")
	if err := os.WriteFile(path, largeObjectBuf, 0o600); err != nil {
		t.Fatalf("failed to write large object file, %v", err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open large object file, %v", err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}
