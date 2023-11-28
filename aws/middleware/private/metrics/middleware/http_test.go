// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"
	"net/http"
	"net/http/httptrace"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics/testutils"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
)

func TestHTTPMetrics_HandleFinalizes(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	counter := metrics.SharedConnectionCounter{}

	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
		},
	}

	ctx := metrics.InitMetricContext(context.TODO(), &counter, &testutils.NoopPublisher{})

	mw := GetHttpMetricMiddleware(&client)

	mctx := metrics.Context(ctx)
	mctx.Data().NewAttempt()

	var traceInput *httptrace.ClientTrace

	addClientTrace = func(ctx context.Context, trace *httptrace.ClientTrace) context.Context {
		traceInput = trace
		return ctx
	}

	_, _, _ = mw.HandleFinalize(ctx, middleware.FinalizeInput{}, testutils.NoopFinalizeHandler{})

	if traceInput == nil {
		t.Fatal("trace should be added to context")
	}
}

func TestHTTPMetrics_HandleFinalizes_AttemptErr(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	counter := metrics.SharedConnectionCounter{}

	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
		},
	}

	ctx := metrics.InitMetricContext(context.TODO(), &counter, &testutils.NoopPublisher{})

	mw := GetHttpMetricMiddleware(&client)

	var traceInput *httptrace.ClientTrace

	addClientTrace = func(ctx context.Context, trace *httptrace.ClientTrace) context.Context {
		traceInput = trace
		return ctx
	}

	_, _, _ = mw.HandleFinalize(ctx, middleware.FinalizeInput{}, testutils.NoopFinalizeHandler{})

	if traceInput != nil {
		t.Fatal("trace should not be added to context")
	}

}

func TestHTTPMetrics_callbacks(t *testing.T) {

	counter := metrics.SharedConnectionCounter{}

	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
		},
	}

	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &counter, &testutils.NoopPublisher{})
	mctx := metrics.Context(ctx)
	mctx.Data().NewAttempt()
	attempt, _ := mctx.Data().LatestAttempt()

	counter.AddPendingConnectionAcquire()
	counter.AddActiveRequest()

	gotFirstResponseByte(attempt, sdk.NowTime())

	getConn(attempt, &counter, sdk.NowTime(), &client, "hostPort")

	gotConn(attempt, &counter, httptrace.GotConnInfo{
		Conn:     nil,
		Reused:   false,
		WasIdle:  false,
		IdleTime: 0,
	}, sdk.NowTime())

	actualConnRequestedTime := attempt.ConnRequestedTime
	expectedConnRequestedTime := sdk.NowTime()

	if actualConnRequestedTime != expectedConnRequestedTime {
		t.Errorf("Unexpected ConnRequestedTime, should be '%s' but was '%s'", expectedConnRequestedTime, actualConnRequestedTime)
	}

	actualPendingConnectionAcquires := attempt.PendingConnectionAcquires
	expectedPendingConnectionAcquires := 1

	if actualPendingConnectionAcquires != expectedPendingConnectionAcquires {
		t.Errorf("Unexpected PendingConnectionAcquires, should be '%d' but was '%d'", expectedPendingConnectionAcquires, actualPendingConnectionAcquires)
	}

	actualActiveRequests := attempt.ActiveRequests
	expectedActiveRequests := 1

	if actualActiveRequests != expectedActiveRequests {
		t.Errorf("Unexpected ActiveRequests, should be '%d' but was '%d'", expectedActiveRequests, actualActiveRequests)
	}

	actualConnObtainedTime := attempt.ConnObtainedTime
	expectedConnObtainedTime := sdk.NowTime()

	if actualConnObtainedTime != expectedConnObtainedTime {
		t.Errorf("Unexpected ConnObtainedTime, should be '%s' but was '%s'", actualConnObtainedTime, expectedConnObtainedTime)
	}

}
