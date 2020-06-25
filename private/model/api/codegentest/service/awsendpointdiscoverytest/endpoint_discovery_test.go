package awsendpointdiscoverytest

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestEndpointDiscovery(t *testing.T) {
	cfg := unit.Config()
	cfg.EnableEndpointDiscovery = true

	svc := New(cfg)
	svc.Handlers.Clear()
	svc.Handlers.Send.PushBack(mockSendDescEndpoint)

	var descCount int32
	svc.Handlers.Complete.PushBack(func(r *aws.Request) {
		if r.Operation.Name != opDescribeEndpoints {
			return
		}
		atomic.AddInt32(&descCount, 1)
	})

	for i := 0; i < 2; i++ {
		req := svc.TestDiscoveryIdentifiersRequiredRequest(
			&TestDiscoveryIdentifiersRequiredInput{
				Sdk: aws.String("sdk"),
			},
		)
		req.Handlers.Send.PushBack(func(r *aws.Request) {
			if e, a := "http://foo/", r.HTTPRequest.URL.String(); e != a {
				t.Errorf("expected %q, but received %q", e, a)
			}
		})
		if _, err := req.Send(context.Background()); err != nil {
			t.Fatal(err)
		}
	}

	if e, a := int32(1), atomic.LoadInt32(&descCount); e != a {
		t.Errorf("expect desc endpoint called %d, got %d", e, a)
	}
}

func TestAsyncEndpointDiscovery(t *testing.T) {
	t.Parallel()

	cfg := unit.Config()
	cfg.EnableEndpointDiscovery = true

	svc := New(cfg)
	svc.Handlers.Clear()

	var firstAsyncReq sync.WaitGroup
	firstAsyncReq.Add(1)
	svc.Handlers.Build.PushBack(func(r *aws.Request) {
		if r.Operation.Name == opDescribeEndpoints {
			firstAsyncReq.Wait()
		}
	})
	svc.Handlers.Send.PushBack(mockSendDescEndpoint)

	req := svc.TestDiscoveryOptionalRequest(&TestDiscoveryOptionalInput{
		Sdk: aws.String("sdk"),
	})
	req.Handlers.Send.PushBack(func(r *aws.Request) {
		if e, a := "https://endpoint/", r.HTTPRequest.URL.String(); e != a {
			t.Errorf("expected %q, but received %q", e, a)
		}
	})
	req.Handlers.Complete.PushBack(func(r *aws.Request) {
		firstAsyncReq.Done()
	})
	if _, err := req.Send(context.Background()); err != nil {
		t.Fatal(err)
	}

	var cacheUpdated bool
	for s := time.Now().Add(10 * time.Second); s.After(time.Now()); {
		// Wait for the cache to be updated before making second aws.
		if svc.endpointCache.Has(req.Operation.Name) {
			cacheUpdated = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !cacheUpdated {
		t.Fatalf("expect endpoint cache to be updated, was not")
	}

	req = svc.TestDiscoveryOptionalRequest(&TestDiscoveryOptionalInput{
		Sdk: aws.String("sdk"),
	})
	req.Handlers.Send.PushBack(func(r *aws.Request) {
		if e, a := "http://foo/", r.HTTPRequest.URL.String(); e != a {
			t.Errorf("expected %q, but received %q", e, a)
		}
	})
	if _, err := req.Send(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func mockSendDescEndpoint(r *aws.Request) {
	if r.Operation.Name != opDescribeEndpoints {
		return
	}

	out, _ := r.Data.(*DescribeEndpointsOutput)
	out.Endpoints = []Endpoint{
		{
			Address:              aws.String("http://foo"),
			CachePeriodInMinutes: aws.Int64(5),
		},
	}
	r.Data = out
}
