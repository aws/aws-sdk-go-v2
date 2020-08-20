// +build integration,disabled

package s3manager_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestInteg_GetBucketRegion(t *testing.T) {
	expectRegion := integCfg.Region

	ctx := context.Background()
	region, err := GetBucketRegion(ctx, integCfg,
		aws.StringValue(bucketName), integCfg.Region)

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := expectRegion, region; e != a {
		t.Errorf("expect %s bucket region, got %s", e, a)
	}
}
