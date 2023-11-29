// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
)

const (
	idleConnFieldName     = "idleConn"
	addressFieldName      = "addr"
	unkHttpClientName     = "Other"
	defaultHttpClientName = "Default"
)

type HTTPMetrics struct {
	client *http.Client
}

func GetHttpMetricMiddleware(client *http.Client) *HTTPMetrics {
	return &HTTPMetrics{
		client: client,
	}
}

func (m *HTTPMetrics) ID() string {
	return "HTTPMetrics"
}

func (m *HTTPMetrics) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, attemptError error,
) {
	ctx = m.addTraceContext(ctx)
	finalize, metadata, err := next.HandleFinalize(ctx, in)
	return finalize, metadata, err
}

var addClientTrace = func(ctx context.Context, trace *httptrace.ClientTrace) context.Context {
	return httptrace.WithClientTrace(ctx, trace)
}

func (m *HTTPMetrics) addTraceContext(ctx context.Context) context.Context {
	mctx := metrics.Context(ctx)
	counter := mctx.ConnectionCounter()

	attemptMetrics, attemptErr := mctx.Data().LatestAttempt()

	if attemptErr == nil {
		trace := &httptrace.ClientTrace{
			GotFirstResponseByte: func() {
				gotFirstResponseByte(attemptMetrics, sdk.NowTime())
			},
			GetConn: func(hostPort string) {
				getConn(attemptMetrics, counter, sdk.NowTime(), m.client, hostPort)
			},
			GotConn: func(info httptrace.GotConnInfo) {
				gotConn(attemptMetrics, counter, info, time.Now())
			},
		}

		ctx = addClientTrace(ctx, trace)
	} else {
		fmt.Println(attemptErr)
	}
	return ctx
}

func gotFirstResponseByte(attemptMetrics *metrics.AttemptMetrics, now time.Time) {
	attemptMetrics.FirstByteTime = now
}

func getConn(
	attemptMetrics *metrics.AttemptMetrics, counter *metrics.SharedConnectionCounter, now time.Time, client *http.Client, hostPort string,
) {
	attemptMetrics.ConnRequestedTime = now
	attemptMetrics.PendingConnectionAcquires = int(counter.PendingConnectionAcquire())
	attemptMetrics.ActiveRequests = int(counter.ActiveRequests())

	// Adding HTTP client metrics here since we need the hostPort to identify the correct conn queues.
	addHTTPClientMetrics(attemptMetrics, client, hostPort)
	counter.AddPendingConnectionAcquire()
}

func gotConn(
	attemptMetrics *metrics.AttemptMetrics, counter *metrics.SharedConnectionCounter, info httptrace.GotConnInfo, now time.Time,
) {
	attemptMetrics.ReusedConnection = info.Reused
	attemptMetrics.ConnObtainedTime = now
	counter.RemovePendingConnectionAcquire()
}

func addHTTPClientMetrics(attemptMetrics *metrics.AttemptMetrics, client *http.Client, hostPort string) {

	maxConnsPerHost := -1
	idleConnCountPerHost := -1
	httpClient := unkHttpClientName

	clientInterface := interface{}(client)

	switch clientInterface.(type) {
	// If not a standard HTTP client we cannot retrieve these metrics
	case *http.Client:
		transport := clientInterface.(*http.Client).Transport
		httpClient = defaultHttpClientName
		switch transport.(type) {
		case *http.Transport:

			maxConnsPerHost = transport.(*http.Transport).MaxConnsPerHost

			transportPtr := reflect.ValueOf(transport.(*http.Transport))

			if transportPtr.IsValid() && transportPtr.Kind() == reflect.Pointer {

				transportValue := transportPtr.Elem()
				idleConn := transportValue.FieldByName(idleConnFieldName)

				if idleConn.IsValid() && idleConn.Kind() == reflect.Map {

					IdleConnMap := idleConn.MapRange()
					// We iterate through all the connection queues to look for the target host
					for IdleConnMap.Next() {
						address := IdleConnMap.Key().FieldByName(addressFieldName)

						if address.IsValid() && address.Kind() == reflect.String {

							if address.String() == hostPort {
								// Number of idle connections for the requests host
								idleConnCountPerHost = IdleConnMap.Value().Len()
								break
							}
						}
					}
				}
			}
		}
	}

	attemptMetrics.HTTPClient = httpClient
	attemptMetrics.AvailableConcurrency = idleConnCountPerHost
	attemptMetrics.MaxConcurrency = maxConnsPerHost
}
