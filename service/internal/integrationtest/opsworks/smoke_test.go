// +build integration

package opsworks

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/opsworks"
	"github.com/awslabs/smithy-go"
)

func TestInteg_00_DescribeStacks(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := opsworks.NewFromConfig(cfg)
	params := &opsworks.DescribeStacksInput{}
	_, err = client.DescribeStacks(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestInteg_01_DescribeLayers(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := opsworks.NewFromConfig(cfg)
	params := &opsworks.DescribeLayersInput{
		StackId: aws.String("fake_stack"),
	}
	_, err = client.DescribeLayers(ctx, params)
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
