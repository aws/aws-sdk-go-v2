// +build integration

package codestar

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/codestar"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
)

func TestInteg_00_ListProjects(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := codestar.NewFromConfig(cfg)
	params := &codestar.ListProjectsInput{}
	_, err = client.ListProjects(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
