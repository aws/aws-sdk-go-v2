// +build integration

package s3

import (
	"testing"
)

func TestInteg_AccessPoint_WriteToObject(t *testing.T) {
	testWriteToObject(t, setupMetadata.AccessPoints.Source.ARN, nil)
}

func TestInteg_AccessPoint_CopyObject(t *testing.T) {
	t.Skip("API does not support access point")
	testCopyObject(t,
		setupMetadata.AccessPoints.Source.ARN,
		setupMetadata.AccessPoints.Target.ARN, nil)
}
