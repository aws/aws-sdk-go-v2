// +build integration

package s3

import (
	"testing"
)

func TestInteg_AccessPoint_WriteToObject(t *testing.T) {
	// fmt.Printf("AccessPoints target name : %v", setupMetadata.AccessPoints.Target.Name)

	testWriteToObject(t, setupMetadata.AccessPoints.Source.ARN, nil)
}

func TestInteg_AccessPoint_CopyObject(t *testing.T) {
	// fmt.Printf("AccessPoints source name : %v", setupMetadata.AccessPoints.Source.Name)
	// fmt.Printf("AccessPoints target name : %v", setupMetadata.AccessPoints.Target.Name)

	t.Skip("API does not support access point")
	testCopyObject(t,
		setupMetadata.AccessPoints.Source.ARN,
		setupMetadata.AccessPoints.Target.ARN, nil)
}
