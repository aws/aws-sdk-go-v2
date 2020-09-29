// +build integration

package s3

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestInteg_00_ListBuckets(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := s3.NewFromConfig(cfg)
	params := &s3.ListBucketsInput{}
	_, err = client.ListBuckets(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
