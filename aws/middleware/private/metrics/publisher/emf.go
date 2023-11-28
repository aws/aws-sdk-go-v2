package publisher

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics/emf"
)

// EMFPublisher is a MetricPublisher implementation that publishes metrics to stdout using EMF format.
type EMFPublisher struct {
	namespace            string
	serializer           metrics.Serializer
	additionalDimensions map[string]string
}

var output = func(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// NewEMFPublisher creates a new EMFPublisher with the specified namespace and serializer.
func NewEMFPublisher(namespace string, serializer metrics.Serializer) *EMFPublisher {
	return &EMFPublisher{
		namespace:            namespace,
		serializer:           serializer,
		additionalDimensions: map[string]string{},
	}
}

func (p *EMFPublisher) SetAdditionalDimension(key string, value string) {
	p.additionalDimensions[key] = value
}

func (p *EMFPublisher) RemoveAdditionalDimension(key string) {
	delete(p.additionalDimensions, key)
}

func (p *EMFPublisher) populateWithAdditionalDimensions(entry *emf.Entry) {
	for k := range p.additionalDimensions {
		entry.AddDimension(k, p.additionalDimensions[k])
	}
}

// perRequestMetrics generates and returns the log entry for per-request metrics.
func (p *EMFPublisher) perRequestMetrics(data *metrics.MetricData) (string, error) {

	entry := emf.NewEntry(p.namespace, p.serializer)

	p.populateWithAdditionalDimensions(&entry)

	entry.AddDimension(metrics.ServiceIDKey, data.ServiceID)
	entry.AddDimension(metrics.OperationNameKey, data.OperationName)
	entry.AddDimension(metrics.HTTPStatusCodeKey, strconv.Itoa(data.StatusCode))

	entry.AddProperty(metrics.ClientRequestIDKey, data.ClientRequestID)

	entry.AddMetric(metrics.APICallDurationKey, float64(data.APICallDuration.Nanoseconds()))
	entry.AddMetric(metrics.APICallSuccessfulKey, float64(data.Success))
	entry.AddMetric(metrics.MarshallingDurationKey, float64(data.MarshallingDuration.Nanoseconds()))
	entry.AddMetric(metrics.EndpointResolutionDurationKey, float64(data.EndpointResolutionDuration.Nanoseconds()))

	entry.AddMetric(metrics.RetryCountKey, float64(data.RetryCount))

	// We only publish throughput if different then 0 to avoid polluting statistics
	if data.InThroughput != 0 {
		entry.AddMetric(metrics.InThroughputKey, data.InThroughput)
	}
	if data.OutThroughput != 0 {
		entry.AddMetric(metrics.OutThroughputKey, data.OutThroughput)
	}

	return entry.Build()
}

// perAttemptMetrics generates and returns the log entry for per-attempt metrics.
func (p *EMFPublisher) perAttemptMetrics(data *metrics.MetricData, attemptIndex int) (string, error) {

	attempt := data.Attempts[attemptIndex]

	entry := emf.NewEntry(p.namespace, p.serializer)

	p.populateWithAdditionalDimensions(&entry)

	entry.AddDimension(metrics.ServiceIDKey, data.ServiceID)
	entry.AddDimension(metrics.OperationNameKey, data.OperationName)
	entry.AddDimension(metrics.HTTPStatusCodeKey, strconv.Itoa(attempt.StatusCode))

	entry.AddProperty(metrics.ClientRequestIDKey, data.ClientRequestID)
	entry.AddProperty(metrics.AWSExtendedRequestIDKey, attempt.ExtendedRequestID)
	entry.AddProperty(metrics.AWSRequestIDKey, attempt.RequestID)
	entry.AddProperty(metrics.AttemptNumberKey, attemptIndex)

	entry.AddMetric(metrics.MaxConcurrencyKey, float64(attempt.MaxConcurrency))
	entry.AddMetric(metrics.AvailableConcurrencyKey, float64(attempt.AvailableConcurrency))
	entry.AddMetric(metrics.ConcurrencyAcquireDurationKey, float64(attempt.ConcurrencyAcquireDuration.Nanoseconds()))
	entry.AddMetric(metrics.PendingConcurrencyAcquiresKey, float64(attempt.PendingConnectionAcquires))
	entry.AddMetric(metrics.SigningDurationKey, float64(attempt.SigningDuration.Nanoseconds()))
	entry.AddMetric(metrics.UnmarshallingDurationKey, float64(attempt.UnMarshallingDuration.Nanoseconds()))
	entry.AddMetric(metrics.TimeToFirstByteKey, float64(attempt.TimeToFirstByte.Nanoseconds()))
	entry.AddMetric(metrics.ServiceCallDurationKey, float64(attempt.ServiceCallDuration.Nanoseconds()))
	entry.AddMetric(metrics.BackoffDelayDurationKey, float64(attempt.RetryDelay))

	return entry.Build()
}

// perStreamMetrics generates and returns the log entry for per-stream metrics.
func (p *EMFPublisher) perStreamMetrics(data *metrics.MetricData) (string, error) {

	entry := emf.NewEntry(p.namespace, p.serializer)

	p.populateWithAdditionalDimensions(&entry)

	entry.AddDimension(metrics.ServiceIDKey, data.ServiceID)
	entry.AddDimension(metrics.OperationNameKey, data.OperationName)
	entry.AddDimension(metrics.HTTPStatusCodeKey, strconv.Itoa(data.StatusCode))

	entry.AddProperty(metrics.ClientRequestIDKey, data.ClientRequestID)

	if data.Stream.Throughput > 0 {
		entry.AddMetric(metrics.StreamThroughputKey, data.Stream.Throughput)
	}

	return entry.Build()
}

// PostRequestMetrics publishes the request metrics to stdout using EMF format.
func (p *EMFPublisher) PostRequestMetrics(data *metrics.MetricData) error {
	requestMetricLogEntry, err := p.perRequestMetrics(data)
	if err != nil {
		output("error generating log entry for request metrics due to %s", err.Error())
	} else {
		output("%s\n", requestMetricLogEntry)
	}
	for idx := range data.Attempts {
		attemptMetricLogEntry, err := p.perAttemptMetrics(data, idx)
		if err != nil {
			output("error generating log entry for attempt metrics due to %s", err.Error())
		} else {
			output("%s\n", attemptMetricLogEntry)
		}
	}
	return nil
}

// PostStreamMetrics publishes the stream metrics to stdout using EMF format.
func (p *EMFPublisher) PostStreamMetrics(data *metrics.MetricData) error {
	streamMetrics, err := p.perStreamMetrics(data)
	if err != nil {
		output("error generating log entry for stream metrics due to %s", err.Error())
	} else {
		output("%s\n", streamMetrics)
	}
	return nil
}
