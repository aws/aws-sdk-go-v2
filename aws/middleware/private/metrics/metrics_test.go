// This package is designated as private and is intended for use only by the
// AWS client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package metrics

import (
	"context"
	"testing"
	"time"
)

type TestPublisher struct{}

func (tp *TestPublisher) PostRequestMetrics(data *MetricData) error {
	return nil
}

func (tp *TestPublisher) PostStreamMetrics(data *MetricData) error {
	return nil
}

func TestInitAndRetrieveMetricContext(t *testing.T) {
	ctx := context.Background()
	cc := SharedConnectionCounter{}
	tp := TestPublisher{}

	ctx = InitMetricContext(ctx, &cc, &tp)
	mctx := Context(ctx)

	if mctx == nil {
		t.Errorf("Metric context should not be nil")
	}
	if mctx.publisher != &tp || mctx.Publisher() != &tp {
		t.Errorf("Unexpected publisher")
	}
	if mctx.connectionCounter != &cc || mctx.ConnectionCounter() != &cc {
		t.Errorf("Unexpected connection counter")
	}
}

func TestConnectionCounter(t *testing.T) {
	cc := SharedConnectionCounter{}
	cc.AddPendingConnectionAcquire()
	cc.AddPendingConnectionAcquire()
	cc.RemovePendingConnectionAcquire()
	cc.AddActiveRequest()
	cc.AddActiveRequest()
	cc.AddActiveRequest()
	cc.RemoveActiveRequest()

	if cc.PendingConnectionAcquire() != 1 {
		t.Errorf("Unexpected count for PendingConnectionAcquire")
	}

	if cc.ActiveRequests() != 2 {
		t.Errorf("Unexpected count for ActiveRequests")
	}
}

func TestAttemptCreationAndRetrieval(t *testing.T) {
	ctx := context.TODO()
	cc := SharedConnectionCounter{}
	tp := TestPublisher{}

	ctx = InitMetricContext(ctx, &cc, &tp)

	mctx := Context(ctx)

	_, err := mctx.Data().LatestAttempt()

	if err == nil {
		t.Errorf("Expected error for uninitialized attempt")
	}

	mctx.Data().NewAttempt()

	if len(mctx.Data().Attempts) != 1 {
		t.Errorf("Unexpected number of attempts")
	}

	_, err = mctx.Data().LatestAttempt()

	if err != nil {
		t.Errorf("Unexpected error for uninitialized attempt")
	}
}

func TestMetricData_ComputeRequestMetrics(t *testing.T) {

	data := MetricData{
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
		Stream:                   StreamMetrics{},
		Attempts: []AttemptMetrics{{
			ServiceCallStart:          time.Unix(1234, 0),
			ServiceCallEnd:            time.Unix(1434, 0),
			FirstByteTime:             time.Unix(1334, 0),
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
				FirstByteTime:             time.Unix(1334, 0),
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
	}

	data.ComputeRequestMetrics()

	expectedAPICallDuration := time.Second * 200
	actualAPICallDuration := data.APICallDuration

	if expectedAPICallDuration != actualAPICallDuration {
		t.Errorf("Unexpected ApiCallDuration, should be '%s' but was '%s'", expectedAPICallDuration, actualAPICallDuration)
	}

	expectedMarshallingDuration := time.Second * 200
	actualMarshallingDuration := data.MarshallingDuration

	if expectedMarshallingDuration != actualMarshallingDuration {
		t.Errorf("Unexpected MarshallingDuration, should be '%s' but was '%s'", expectedMarshallingDuration, actualMarshallingDuration)
	}

	expectedEndpointResolutionDuration := time.Second * 200
	actualEndpointResolutionDuration := data.EndpointResolutionDuration

	if expectedEndpointResolutionDuration != actualEndpointResolutionDuration {
		t.Errorf("Unexpected EndpointResolutionDuration, should be '%s' but was '%s'", expectedEndpointResolutionDuration, actualEndpointResolutionDuration)
	}

	for idx := range data.Attempts {

		attempt := data.Attempts[idx]

		expectedServiceCallDuration := time.Second * 200
		actualServiceCallDuration := attempt.ServiceCallDuration

		if expectedServiceCallDuration != actualServiceCallDuration {
			t.Errorf("Unexpected ServiceCallDuration, should be '%s' but was '%s'", expectedServiceCallDuration, actualServiceCallDuration)
		}

		expectedTimeToFirstByte := time.Second * 100
		actualTimeToFirstByte := attempt.TimeToFirstByte

		if expectedTimeToFirstByte != actualTimeToFirstByte {
			t.Errorf("Unexpected TimeToFirstByte, should be '%s' but was '%s'", expectedTimeToFirstByte, actualTimeToFirstByte)
		}

		expectedConcurrencyAcquireDuration := time.Second * 200
		actualConcurrencyAcquireDuration := attempt.ConcurrencyAcquireDuration

		if expectedConcurrencyAcquireDuration != actualConcurrencyAcquireDuration {
			t.Errorf("Unexpected ConcurrencyAcquireDuration, should be '%s' but was '%s'", expectedConcurrencyAcquireDuration, actualConcurrencyAcquireDuration)
		}

		expectedSigningDuration := time.Second * 200
		actualSigningDuration := attempt.SigningDuration

		if expectedSigningDuration != actualSigningDuration {
			t.Errorf("Unexpected SigningDuration, should be '%s' but was '%s'", expectedSigningDuration, actualSigningDuration)
		}

		expectedUnMarshallingDuration := time.Second * 200
		actualUnMarshallingDuration := attempt.UnMarshallingDuration

		if expectedUnMarshallingDuration != actualUnMarshallingDuration {
			t.Errorf("Unexpected UnMarshallingDuration, should be '%s' but was '%s'", expectedUnMarshallingDuration, actualUnMarshallingDuration)
		}
	}
}
