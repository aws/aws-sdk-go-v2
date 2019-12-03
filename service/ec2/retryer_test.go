package ec2

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestCustomRetryRules(t *testing.T) {

	cfg := unit.Config()
	svc := New(cfg)

	req := svc.ModifyNetworkInterfaceAttributeRequest(&types.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String("foo"),
	})

	duration := req.Request.Retryer.RetryRules(req.Request)
	if duration < time.Second*1 || duration > time.Second*2 {
		t.Errorf("expected duration to be between 1s and 2s, but received %s", duration)
	}

	req.Request.RetryCount = 15
	duration = req.Request.Retryer.RetryRules(req.Request)

	if duration < time.Second*4 || duration > time.Second*8 {
		t.Errorf("expected duration to be between 4s and 8s, but received %s", duration)
	}

}

func TestCustomRetryer_WhenRetrierSpecified(t *testing.T) {
	svc := New(aws.Config{
		Region: "us-west-2",
		Retryer: aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
			d.NumMaxRetries = 4
			d.MinThrottleDelay = 50 * time.Millisecond
			d.MinRetryDelay = 10 * time.Millisecond
			d.MaxThrottleDelay = 200 * time.Millisecond
			d.MaxRetryDelay = 300 * time.Millisecond
		}),
		EndpointResolver: unit.Config().EndpointResolver,
	})

	if _, ok := svc.Client.Retryer.(aws.DefaultRetryer); !ok {
		t.Error("expected default retryer, but received otherwise")
	}

	req := svc.AssignPrivateIpAddressesRequest(&types.AssignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String("foo"),
	})

	d := req.Request.Retryer.(aws.DefaultRetryer)

	if d.NumMaxRetries != 4 {
		t.Errorf("expected max retries to be %v, got %v", 4, d.NumMaxRetries)
	}

	if d.MinRetryDelay != 10*time.Millisecond {
		t.Errorf("expected min retry delay to be %v, got %v", "10 ms", d.MinRetryDelay)
	}

	if d.MinThrottleDelay != 50*time.Millisecond {
		t.Errorf("expected min throttle delay to be %v, got %v", "50 ms", d.MinThrottleDelay)
	}

	if d.MaxRetryDelay != 300*time.Millisecond {
		t.Errorf("expected max retry delay to be %v, got %v", "300 ms", d.MaxRetryDelay)
	}

	if d.MaxThrottleDelay != 200*time.Millisecond {
		t.Errorf("expected max throttle delay to be %v, got %v", "200 ms", d.MaxThrottleDelay)
	}
}

func TestCustomRetryer(t *testing.T) {

	cfg := unit.Config()
	svc := New(cfg)

	req := svc.AssignPrivateIpAddressesRequest(&types.AssignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String("foo"),
	})

	d := req.Request.Retryer.(aws.DefaultRetryer)

	if d.NumMaxRetries != customRetryerMaxNumRetries {
		t.Errorf("expected max retries to be %v, got %v", customRetryerMaxNumRetries, d.NumMaxRetries)
	}

	if d.MinRetryDelay != customRetryerMinRetryDelay {
		t.Errorf("expected min retry delay to be %v, got %v", customRetryerMinRetryDelay, d.MinRetryDelay)
	}

	if d.MinThrottleDelay != customRetryerMinRetryDelay {
		t.Errorf("expected min throttle delay to be %v, got %v", customRetryerMinRetryDelay, d.MinThrottleDelay)
	}

	if d.MaxRetryDelay != customRetryerMaxRetryDelay {
		t.Errorf("expected max retry delay to be %v, got %v", customRetryerMaxRetryDelay, d.MaxRetryDelay)
	}

	if d.MaxThrottleDelay != customRetryerMaxRetryDelay {
		t.Errorf("expected max throttle delay to be %v, got %v", customRetryerMaxRetryDelay, d.MaxThrottleDelay)
	}
}
