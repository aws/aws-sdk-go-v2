package ec2metadata_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestClientDisableIMDS(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_EC2_METADATA_DISABLED", "true")

	cfg := unit.Config()
	cfg.LogLevel = aws.LogDebugWithHTTPBody
	cfg.Logger = t

	svc := ec2metadata.New(cfg)
	resp, err := svc.GetUserData(context.Background())
	if err == nil {
		t.Fatalf("expect error, got none")
	}
	if len(resp) != 0 {
		t.Errorf("expect no response, got %v", resp)
	}

	var ce *aws.RequestCanceledError
	if !errors.As(err, &ce) {
		t.Fatalf("expect %T error, got %v", ce, err)
	}

	if e, a := "AWS_EC2_METADATA_DISABLED", ce.Err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v in error message, got %v", e, a)
	}
}
