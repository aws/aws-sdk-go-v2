// +build integration

package autoscaling

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/awslabs/smithy-go"
)

func TestInteg_00_DescribeScalingProcessTypes(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := autoscaling.NewFromConfig(cfg)
	params := &autoscaling.DescribeScalingProcessTypesInput{}
	_, err = client.DescribeScalingProcessTypes(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestInteg_01_CreateLaunchConfiguration(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := autoscaling.NewFromConfig(cfg)
	params := &autoscaling.CreateLaunchConfigurationInput{
		ImageId:                 aws.String("ami-12345678"),
		InstanceType:            aws.String("m1.small"),
		LaunchConfigurationName: aws.String("hello, world"),
	}
	_, err = client.CreateLaunchConfiguration(ctx, params)
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
