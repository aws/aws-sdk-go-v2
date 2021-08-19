//go:build integration
// +build integration

package timestreamwrite

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	tw "github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"testing"
	"time"
)

func TestInteg_00_CreateDatabase(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-east-1")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	// Create an Amazon timestreamwrite service client
	client := tw.NewFromConfig(cfg)

	// CreateDatabase
	output, err := client.CreateDatabase(ctx, &tw.CreateDatabaseInput{
		DatabaseName: aws.String("testDB"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// DeleteDatabase
	_, err = client.DeleteDatabase(ctx, &tw.DeleteDatabaseInput{
		DatabaseName: output.Database.DatabaseName,
	})
	if err != nil {
		t.Fatal(err)
	}
}
