// +build integration

package costandusagereportservice

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costandusagereportservice"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
)

func TestInteg_00_DescribeReportDefinitions(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-east-1")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := costandusagereportservice.NewFromConfig(cfg)
	params := &costandusagereportservice.DescribeReportDefinitionsInput{}
	_, err = client.DescribeReportDefinitions(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
