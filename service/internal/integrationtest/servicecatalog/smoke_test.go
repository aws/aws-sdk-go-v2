// +build integration

package servicecatalog

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
)

func TestInteg_00_ListAcceptedPortfolioShares(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := servicecatalog.NewFromConfig(cfg)
	params := &servicecatalog.ListAcceptedPortfolioSharesInput{}
	_, err = client.ListAcceptedPortfolioShares(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
