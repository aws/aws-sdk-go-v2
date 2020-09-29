// +build integration

package redshift

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/awslabs/smithy-go"
)

func TestInteg_00_DescribeClusterVersions(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := redshift.NewFromConfig(cfg)
	params := &redshift.DescribeClusterVersionsInput{}
	_, err = client.DescribeClusterVersions(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestInteg_01_DescribeClusters(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := redshift.NewFromConfig(cfg)
	params := &redshift.DescribeClustersInput{
		ClusterIdentifier: aws.String("fake-cluster"),
	}
	_, err = client.DescribeClusters(ctx, params)
	if err == nil {
		t.Fatalf("expect request to fail")
	}

	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expect error to be API error, was not, %v", err)
	}
	if len(apiErr.ErrorCode()) == 0 {
		t.Errorf("expect non-empty error code")
	}
	if len(apiErr.ErrorMessage()) == 0 {
		t.Errorf("expect non-empty error message")
	}
}
