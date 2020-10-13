// +build integration

package s3manager_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/s3manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestGetBucketRegion(t *testing.T) {
	expectRegion := integConfig.Region

	region, err := s3manager.GetBucketRegion(context.Background(), s3.NewFromConfig(integConfig), aws.ToString(bucketName))
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := expectRegion, region; e != a {
		t.Errorf("expect %s bucket region, got %s", e, a)
	}
}
