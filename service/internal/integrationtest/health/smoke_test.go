// +build integration

package health

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
)

func TestInteg_00_DescribeEntityAggregates(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-east-1")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := health.NewFromConfig(cfg)
	params := &health.DescribeEntityAggregatesInput{}
	_, err = client.DescribeEntityAggregates(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
