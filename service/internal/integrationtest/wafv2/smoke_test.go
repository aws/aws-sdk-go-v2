// +build integration

package wafv2

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/awslabs/smithy-go"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
)

func TestInteg_00_ListWebACLs(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, _ := integrationtest.LoadConfigWithDefaultRegion("us-east-1")
	svc := wafv2.NewFromConfig(cfg)

	input := &wafv2.ListWebACLsInput{
		Limit: aws.Int32(20),
		Scope: types.ScopeRegional,
	}
	_, err := svc.ListWebACLs(ctx, input)

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestInteg_01_CreateRegexPatternSet(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, _ := integrationtest.LoadConfigWithDefaultRegion("us-east-1")
	svc := wafv2.NewFromConfig(cfg)

	input := &wafv2.CreateRegexPatternSetInput{
		Name:  aws.String("fake_name"),
		Scope: types.ScopeRegional,
	}

	_, err := svc.CreateRegexPatternSet(ctx, input)
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
