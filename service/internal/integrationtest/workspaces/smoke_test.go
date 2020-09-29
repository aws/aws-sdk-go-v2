// +build integration

package workspaces

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
	"github.com/awslabs/smithy-go"
)

func TestInteg_00_DescribeWorkspaces(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, _ := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	svc := workspaces.NewFromConfig(cfg)

	input := &workspaces.DescribeWorkspacesInput{}
	_, err := svc.DescribeWorkspaces(ctx, input)

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestInteg_01_DescribeWorkspaces(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, _ := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	svc := workspaces.NewFromConfig(cfg)

	input := &workspaces.DescribeWorkspacesInput{
		DirectoryId: aws.String("fake-id"),
	}
	_, err := svc.DescribeWorkspaces(ctx, input)

	if err == nil {
		t.Errorf("expect request to fail")
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
