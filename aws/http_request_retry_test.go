package aws_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

func TestRequestCancelRetry(t *testing.T) {
	restoreSleep := mockSleep()
	defer restoreSleep()

	c := make(chan struct{})

	reqNum := 0
	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 1
	})

	s := mock.NewMockClient(cfg)

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.Clear()
	s.Handlers.UnmarshalMeta.Clear()
	s.Handlers.UnmarshalError.Clear()
	s.Handlers.Send.PushFront(func(r *aws.Request) {
		reqNum++
	})
	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	r.HTTPRequest.Cancel = c
	close(c)

	err := r.Send()
	if e, a := "canceled", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q to be in %q", e, a)
	}
	if e, a := 1, reqNum; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func mockSleep() func() {
	origSleep := sdk.Sleep
	sdk.Sleep = func(time.Duration) {}

	origCtxSleep := sdk.SleepWithContext
	sdk.SleepWithContext = func(context.Context, time.Duration) error { return nil }

	return func() {
		sdk.Sleep = origSleep
		sdk.SleepWithContext = origCtxSleep
	}
}
