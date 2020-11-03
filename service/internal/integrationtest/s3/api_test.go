// +build integration

package s3

import (
	"testing"
)

func TestInteg_WriteToObject(t *testing.T) {
	testWriteToObject(t, setupMetadata.Buckets.Source.Name, nil)
}

func TestInteg_CopyObject(t *testing.T) {
	testCopyObject(t, setupMetadata.Buckets.Source.Name, setupMetadata.Buckets.Target.Name, nil)
}
