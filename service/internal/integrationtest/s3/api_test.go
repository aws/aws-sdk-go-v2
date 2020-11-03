// +build integration

package s3

import (
	"fmt"
	"testing"
)

func TestInteg_WriteToObject(t *testing.T) {
	fmt.Printf("bucket source name : %v", setupMetadata.Buckets.Source.Name)
	// testWriteToObject(t, setupMetadata.Buckets.Source.Name, nil)
}

func TestInteg_CopyObject(t *testing.T) {
	fmt.Printf("bucket target name : %v", setupMetadata.Buckets.Target.Name)

	// testCopyObject(t, setupMetadata.Buckets.Source.Name, setupMetadata.Buckets.Target.Name, nil)
}
