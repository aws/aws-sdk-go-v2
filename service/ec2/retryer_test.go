package ec2

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestCustomRetryer(t *testing.T) {

	svc := New(unit.Config())

	if _, ok := svc.Client.Retryer.(retryer); !ok {
		t.Error("expected custom retryer, but received otherwise")
	}

	req := svc.ModifyNetworkInterfaceAttributeRequest(&ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String("foo"),
	})

	duration := svc.Client.Retryer.RetryRules(req.Request)
	if duration < time.Second*1 || duration > time.Second*2 {
		t.Errorf("expected duration to be between 1s and 2s, but received %s", duration)
	}

	req.Request.RetryCount = 15
	duration = svc.Client.Retryer.RetryRules(req.Request)

	if duration < time.Second*4 || duration > time.Second*8 {
		t.Errorf("expected duration to be between 4s and 8s, but received %s", duration)
	}

	svc = New(aws.Config{
		Region:  "us-west-2",
		Retryer: aws.DefaultRetryer{}},
	)
	if _, ok := svc.Client.Retryer.(aws.DefaultRetryer); !ok {
		t.Error("expected default retryer, but received otherwise")
	}
}
