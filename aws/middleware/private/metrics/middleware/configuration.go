package middleware

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/smithy-go/middleware"
)

func WithMetricMiddlewares(
	publisher metrics.MetricPublisher, client *http.Client,
) func(stack *middleware.Stack) error {
	connectionCounter := &metrics.SharedConnectionCounter{}
	return func(stack *middleware.Stack) error {
		if err := stack.Initialize.Add(GetSetupMetricCollectionMiddleware(connectionCounter, publisher), middleware.Before); err != nil {
			return err
		}
		if err := stack.Serialize.Add(GetRecordStackSerializeStartMiddleware(), middleware.Before); err != nil {
			return err
		}
		if err := stack.Serialize.Add(GetRecordStackSerializeEndMiddleware(), middleware.After); err != nil {
			return err
		}
		if err := stack.Finalize.Insert(GetRecordEndpointResolutionStartMiddleware(), "ResolveEndpointV2", middleware.Before); err != nil {
			return err
		}
		if err := stack.Finalize.Insert(GetRecordEndpointResolutionEndMiddleware(), "ResolveEndpointV2", middleware.After); err != nil {
			return err
		}
		if err := stack.Build.Add(GetWrapDataStreamMiddleware(), middleware.After); err != nil {
			return err
		}
		if err := stack.Finalize.Add(GetRegisterRequestMetricContextMiddleware(), middleware.Before); err != nil {
			return err
		}
		if err := stack.Finalize.Insert(GetRegisterAttemptMetricContextMiddleware(), "Retry", middleware.After); err != nil {
			return err
		}
		if err := stack.Finalize.Add(GetHttpMetricMiddleware(client), middleware.After); err != nil {
			return err
		}
		if err := stack.Deserialize.Add(GetRecordStackDeserializeStartMiddleware(), middleware.After); err != nil {
			return err
		}
		if err := stack.Deserialize.Add(GetRecordStackDeserializeEndMiddleware(), middleware.Before); err != nil {
			return err
		}
		if err := stack.Deserialize.Insert(GetTransportMetricsMiddleware(), "StackDeserializeStart", middleware.After); err != nil {
			return err
		}
		if err := timeGetIdentity(stack); err != nil {
			return err
		}
		if err := timeSigning(stack); err != nil {
			return err
		}
		if err := stack.Build.Add(&captureUserAgent{}, middleware.After); err != nil {
			return err
		}
		return nil
	}
}
