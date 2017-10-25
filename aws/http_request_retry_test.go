// +build go1.5

package aws_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"

	"github.com/stretchr/testify/assert"
)

func TestRequestCancelRetry(t *testing.T) {
	c := make(chan struct{})

	reqNum := 0
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("http://endpoint")
	cfg.Retryer = aws.DefaultRetryer{NumMaxRetries: 10}
	s := mock.NewMockClient(cfg)

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.Clear()
	s.Handlers.UnmarshalMeta.Clear()
	s.Handlers.UnmarshalError.Clear()
	s.Handlers.Send.PushFront(func(r *aws.Request) {
		reqNum++
		r.Error = errors.New("net/http: request canceled")
	})
	out := &testData{}
	r := s.NewRequest(&aws.Operation{Name: "Operation"}, nil, out)
	r.HTTPRequest.Cancel = c
	close(c)

	err := r.Send()
	fmt.Println("request error", err)
	assert.True(t, strings.Contains(err.Error(), "canceled"))
	assert.Equal(t, 1, reqNum)
}
