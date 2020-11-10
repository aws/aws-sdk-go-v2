// +build integration

package s3

import (
	"testing"
)

func TestInteg_AccessPoint_WriteToObject(t *testing.T) {
	t.Skip("skip till accesspoint support is merged in")
	testWriteToObject(t, setupMetadata.AccessPoints.Source.ARN, nil)
}
