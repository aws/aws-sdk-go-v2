package middleware

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics/testutils"
	"github.com/aws/smithy-go/middleware"
)

type TestResult struct {
	Body io.ReadCloser
}

func TestWrapDataStream_HandleBuild(t *testing.T) {

	cases := map[string]struct {
		ConnectionCounter  *metrics.SharedConnectionCounter
		Publisher          testutils.MetricDataRecorderPublisher
		Result             *TestResult
		Input              middleware.BuildInput
		ExpectedStreamData string
		ExpectedMetricData metrics.StreamMetrics
	}{
		"success": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         testutils.MetricDataRecorderPublisher{},
			Input:             middleware.BuildInput{},
			Result: &TestResult{
				Body: &testutils.TestReadCloser{Data: []byte("testString")},
			},
			ExpectedStreamData: "testString",
			ExpectedMetricData: metrics.StreamMetrics{
				ReadDuration: 0,
				ReadBytes:    10,
			},
		},
		"emptyData": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         testutils.MetricDataRecorderPublisher{},
			Input:             middleware.BuildInput{},
			Result: &TestResult{
				Body: &testutils.TestReadCloser{Data: []byte("")},
			},
			ExpectedStreamData: "",
			ExpectedMetricData: metrics.StreamMetrics{
				ReadDuration: 0,
				ReadBytes:    0,
			},
		},
		"nilBody": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         testutils.MetricDataRecorderPublisher{},
			Input:             middleware.BuildInput{},
			Result: &TestResult{
				Body: nil,
			},
		},
		"nilResult": {
			ConnectionCounter: &metrics.SharedConnectionCounter{},
			Publisher:         testutils.MetricDataRecorderPublisher{},
			Input:             middleware.BuildInput{},
			Result:            nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			ctx := context.TODO()
			ctx = metrics.InitMetricContext(ctx, c.ConnectionCounter, &c.Publisher)
			mw := GetWrapDataStreamMiddleware()

			out, _, _ := mw.HandleBuild(ctx, c.Input, &testutils.StreamingBodyBuildHandler{Result: c.Result})

			result := out.Result.(*TestResult)

			if result == nil || result.Body == nil {
				return
			}

			readData, _ := io.ReadAll(result.Body)
			actualStreamData := string(readData)
			actualMetricData := c.Publisher.Data.Stream

			if actualStreamData != c.ExpectedStreamData {
				t.Errorf("Unexpected Data, should be '%s' but was '%s'", c.ExpectedStreamData, actualStreamData)
			}
			if !reflect.DeepEqual(actualMetricData, c.ExpectedMetricData) {
				t.Errorf("Unexpected Metric Data, should be '%+v' but was '%+v'", c.ExpectedMetricData, actualMetricData)
			}

		})
	}

}

func TestWrapDataStream_WrongResultType(t *testing.T) {
	pub := &testutils.NoopPublisher{}
	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, pub)
	mw := GetWrapDataStreamMiddleware()

	r1 := TestResult{
		Body: &testutils.TestReadCloser{Data: []byte("testString")},
	}

	_, _, _ = mw.HandleBuild(ctx, middleware.BuildInput{}, &testutils.StreamingBodyBuildHandler{Result: r1})
	_, _, _ = mw.HandleBuild(ctx, middleware.BuildInput{}, &testutils.StreamingBodyBuildHandler{Result: pub})

}
