// +build integration

package s3manager_test

import (
	"context"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/service/s3/s3manager"
)

func TestInteg_GetBucketRegion(t *testing.T) {
	expectRegion := integCfg.Region

	ctx := context.Background()
	region, err := s3manager.GetBucketRegion(ctx, integCfg,
		aws.StringValue(bucketName), integCfg.Region)

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := expectRegion, region; e != a {
		t.Errorf("expect %s bucket region, got %s", e, a)
	}
}
