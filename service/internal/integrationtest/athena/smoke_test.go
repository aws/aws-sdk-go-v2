// +build integration

package athena

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/athena"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
)

func TestInteg_00_ListNamedQueries(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := athena.NewFromConfig(cfg)
	params := &athena.ListNamedQueriesInput{}
	_, err = client.ListNamedQueries(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
