// +build integration

package s3

import (
	"testing"
)

func TestInteg_AccessPoint_WriteToObject(t *testing.T) {
	testWriteToObject(t, setupMetadata.AccessPoints.Source.ARN, nil)
}
