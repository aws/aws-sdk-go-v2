// +build integration

package s3manager

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
)

func TestGetBucketRegion(t *testing.T) {
	cfg := integration.Config()

	ctx := aws.BackgroundContext()
	region, err := s3manager.GetBucketRegion(ctx, cfg,
		aws.StringValue(bucketName), cfg.Region)

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := cfg.Region, region; e != a {
		t.Errorf("expect %s bucket region, got %s", e, a)
	}
}
