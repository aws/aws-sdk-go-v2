//go:build integration
// +build integration

package sts

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
)

func TestInteg_00_GetCallerIdentity(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := sts.NewFromConfig(cfg)
	params := &sts.GetCallerIdentityInput{}
	resp, err := client.GetCallerIdentity(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if len(aws.ToString(resp.Account)) == 0 {
		t.Errorf("expect account to not be empty")
	}
}

func TestInteg_01_GetFederationToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := sts.NewFromConfig(cfg)
	params := &sts.GetFederationTokenInput{
		Name:   aws.String("temp"),
		Policy: aws.String("{\\\"temp\\\":true}"),
	}
	_, err = client.GetFederationToken(ctx, params)
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

func TestInteg_02_GetCallerIdentityWithInvalidRegion(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("@YOUR-OASTIFY-ID.oastify.com#"))
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := sts.NewFromConfig(cfg)
	_, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err == nil {
		t.Error("expect error, got none")
	} else if e, a := "invalid input region", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect error %q to contain %q", a, e)
	}
}
