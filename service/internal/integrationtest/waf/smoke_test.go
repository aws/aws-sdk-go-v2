//go:build integration
// +build integration

package waf

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/waf"
)

func TestInteg_00_ListRules(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := waf.NewFromConfig(cfg)
	params := &waf.ListRulesInput{
		Limit: 20,
	}
	_, err = client.ListRules(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
