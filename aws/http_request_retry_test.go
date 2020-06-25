package aws_test

import (
	"context"
	"strings"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/internal/sdk"
)

func TestRequestCancelRetry(t *testing.T) {
	restoreSleep := sdk.TestingUseNoOpSleep()
	defer restoreSleep()

	var reqNum int
	cfg := unit.Config()
	s := mock.NewMockClient(cfg)

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.Clear()
	s.Handlers.UnmarshalMeta.Clear()
	s.Handlers.UnmarshalError.Clear()
	s.Handlers.Send.PushFront(func(r *aws.Request) {
		reqNum++
	})
	out := &testData{}

	ctx, cancelFn := context.WithCancel(context.Background())
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	r.SetContext(ctx)
	cancelFn() // cancelling the context associated with the request

	err := r.Send()
	if e, a := "canceled", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q to be in %q", e, a)
	}
	if e, a := 1, reqNum; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}
