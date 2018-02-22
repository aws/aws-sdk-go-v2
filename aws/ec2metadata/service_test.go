package ec2metadata_test

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestClientDisableIMDS(t *testing.T) {
	env := awstesting.StashEnv()
	defer awstesting.PopEnv(env)

	os.Setenv("AWS_EC2_METADATA_DISABLED", "true")

	svc := ec2metadata.New(unit.Config())
	resp, err := svc.Region()
	if err == nil {
		t.Fatalf("expect error, got none")
	}
	if len(resp) != 0 {
		t.Errorf("expect no response, got %v", resp)
	}

	aerr := err.(awserr.Error)
	if e, a := aws.ErrCodeRequestCanceled, aerr.Code(); e != a {
		t.Errorf("expect %v error code, got %v", e, a)
	}
	if e, a := "AWS_EC2_METADATA_DISABLED", aerr.Message(); !strings.Contains(a, e) {
		t.Errorf("expect %v in error message, got %v", e, a)
	}
}
