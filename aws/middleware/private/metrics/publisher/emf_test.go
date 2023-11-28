// This package is designated as private and is intended for use only by the
// AWS client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package publisher

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

type TestSerializerWithError struct{}

func (TestSerializerWithError) Serialize(obj interface{}) (string, error) {
	return "", fmt.Errorf("serialization error")
}

func TestPostRequestMetrics(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	cases := map[string]struct {
		Namespace       string
		Serializer      metrics.Serializer
		Data            metrics.MetricData
		ExpectedError   error
		ExpectedResults []string
	}{
		"emptyRequestMetricData": {
			Namespace:     "testNamespace",
			Serializer:    metrics.DefaultSerializer{},
			Data:          metrics.MetricData{},
			ExpectedError: nil,
			ExpectedResults: []string{
				emptyRequestMetricData,
			},
		},
		"serializerError": {
			Namespace:  "testNamespace",
			Serializer: TestSerializerWithError{},
			Data: metrics.MetricData{
				Attempts: []metrics.AttemptMetrics{{}},
			},
			ExpectedError: nil,
			ExpectedResults: []string{
				"error generating log entry for request metrics due to [serialization error]",
				"error generating log entry for attempt metrics due to [serialization error]",
			},
		},
		"completeRequestMetricData": {
			Namespace:  "testNamespace",
			Serializer: metrics.DefaultSerializer{},
			Data: metrics.MetricData{
				RequestStartTime:         time.Unix(1234, 0),
				RequestEndTime:           time.Unix(1434, 0),
				SerializeStartTime:       time.Unix(1234, 0),
				SerializeEndTime:         time.Unix(1434, 0),
				ResolveEndpointStartTime: time.Unix(1234, 0),
				ResolveEndpointEndTime:   time.Unix(1434, 0),
				Success:                  1,
				ClientRequestID:          "crid",
				ServiceID:                "sid",
				OperationName:            "operationname",
				PartitionID:              "partitionid",
				Region:                   "region",
				RequestContentLength:     100,
				Stream:                   metrics.StreamMetrics{},
				Attempts: []metrics.AttemptMetrics{{
					ServiceCallStart:          time.Unix(1234, 0),
					ServiceCallEnd:            time.Unix(1434, 0),
					FirstByteTime:             time.Unix(1234, 0),
					ConnRequestedTime:         time.Unix(1234, 0),
					ConnObtainedTime:          time.Unix(1434, 0),
					CredentialFetchStartTime:  time.Unix(1234, 0),
					CredentialFetchEndTime:    time.Unix(1434, 0),
					SignStartTime:             time.Unix(1234, 0),
					SignEndTime:               time.Unix(1434, 0),
					DeserializeStartTime:      time.Unix(1234, 0),
					DeserializeEndTime:        time.Unix(1434, 0),
					RetryDelay:                100,
					ResponseContentLength:     100,
					StatusCode:                200,
					RequestID:                 "reqid",
					ExtendedRequestID:         "exreqid",
					HTTPClient:                "Default",
					MaxConcurrency:            10,
					PendingConnectionAcquires: 1,
					AvailableConcurrency:      2,
					ActiveRequests:            3,
					ReusedConnection:          false,
				},
					{
						ServiceCallStart:          time.Unix(1234, 0),
						ServiceCallEnd:            time.Unix(1434, 0),
						FirstByteTime:             time.Unix(1234, 0),
						ConnRequestedTime:         time.Unix(1234, 0),
						ConnObtainedTime:          time.Unix(1434, 0),
						CredentialFetchStartTime:  time.Unix(1234, 0),
						CredentialFetchEndTime:    time.Unix(1434, 0),
						SignStartTime:             time.Unix(1234, 0),
						SignEndTime:               time.Unix(1434, 0),
						DeserializeStartTime:      time.Unix(1234, 0),
						DeserializeEndTime:        time.Unix(1434, 0),
						RetryDelay:                100,
						ResponseContentLength:     100,
						StatusCode:                200,
						RequestID:                 "reqid",
						ExtendedRequestID:         "exreqid",
						HTTPClient:                "Default",
						MaxConcurrency:            10,
						PendingConnectionAcquires: 1,
						AvailableConcurrency:      2,
						ActiveRequests:            3,
						ReusedConnection:          false,
					}},
			},
			ExpectedError: nil,
			ExpectedResults: []string{
				completeRequestMetricData,
				completeMetricDataAttempt1,
				completeMetricDataAttempt2,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			var actualResults []string

			output = func(format string, a ...interface{}) {
				actualResults = append(actualResults, fmt.Sprintf(format, a))
			}

			publisher := NewEMFPublisher(c.Namespace, c.Serializer)

			c.Data.ComputeRequestMetrics()

			err := publisher.PostRequestMetrics(&c.Data)

			if !reflect.DeepEqual(err, c.ExpectedError) {
				t.Errorf("Unexpected error, should be '%s' but was '%s'", c.ExpectedError, err)
			}

			if len(c.ExpectedResults) != len(actualResults) {
				t.Errorf("Different number of results. Expected %d but got %d", len(c.ExpectedResults), len(actualResults))
			}

			for i := range c.ExpectedResults {
				if !reflect.DeepEqual(actualResults[i], c.ExpectedResults[i]) {
					t.Errorf("Unexpected result, should be '%s' but was '%s'", c.ExpectedResults[i], actualResults[i])
				}
			}
		})
	}
}

func TestPostStreamMetrics(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	cases := map[string]struct {
		Namespace       string
		Serializer      metrics.Serializer
		Data            metrics.MetricData
		ExpectedError   error
		ExpectedResults []string
	}{
		"emptyStreamMetricData": {
			Namespace:     "testNamespace",
			Serializer:    metrics.DefaultSerializer{},
			Data:          metrics.MetricData{},
			ExpectedError: nil,
			ExpectedResults: []string{
				emptyStreamMetricData,
			},
		},
		"serializerError": {
			Namespace:     "testNamespace",
			Serializer:    TestSerializerWithError{},
			Data:          metrics.MetricData{},
			ExpectedError: nil,
			ExpectedResults: []string{
				"error generating log entry for stream metrics due to [serialization error]",
			},
		},
		"completeStreamMetricData": {
			Namespace:  "testNamespace",
			Serializer: metrics.DefaultSerializer{},
			Data: metrics.MetricData{
				RequestStartTime:         time.Unix(1234, 0),
				RequestEndTime:           time.Unix(1434, 0),
				SerializeStartTime:       time.Unix(1234, 0),
				SerializeEndTime:         time.Unix(1434, 0),
				ResolveEndpointStartTime: time.Unix(1234, 0),
				ResolveEndpointEndTime:   time.Unix(1434, 0),
				Success:                  1,
				StatusCode:               200,
				ClientRequestID:          "crid",
				ServiceID:                "sid",
				OperationName:            "operationname",
				PartitionID:              "partitionid",
				Region:                   "region",
				RequestContentLength:     100,
				Stream: metrics.StreamMetrics{
					ReadDuration: 150,
					ReadBytes:    12,
					Throughput:   80000000,
				},
			},
			ExpectedError: nil,
			ExpectedResults: []string{
				completeStreamMetricData,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			var actualResults []string

			output = func(format string, a ...interface{}) {
				actualResults = append(actualResults, fmt.Sprintf(format, a))
			}

			publisher := NewEMFPublisher(c.Namespace, c.Serializer)

			err := publisher.PostStreamMetrics(&c.Data)

			if !reflect.DeepEqual(err, c.ExpectedError) {
				t.Errorf("Unexpected error, should be '%s' but was '%s'", c.ExpectedError, err)
			}

			if len(c.ExpectedResults) != len(actualResults) {
				t.Errorf("Different number of results. Expected %d but got %d", len(c.ExpectedResults), len(actualResults))
			}

			for i := range c.ExpectedResults {
				if !reflect.DeepEqual(actualResults[i], c.ExpectedResults[i]) {
					t.Errorf("Unexpected result, should be '%s' but was '%s'", c.ExpectedResults[i], actualResults[i])
				}
			}
		})
	}
}
