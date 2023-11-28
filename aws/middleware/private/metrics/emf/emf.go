// Package emf implements an EMF metrics publisher.
//
// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.
package emf

import (
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

const (
	emfIdentifier        = "_aws"
	timestampKey         = "Timestamp"
	cloudWatchMetricsKey = "CloudWatchMetrics"
	namespaceKey         = "Namespace"
	dimensionsKey        = "Dimensions"
	metricsKey           = "Metrics"
)

// Entry represents a log entry in the EMF format.
type Entry struct {
	namespace  string
	serializer metrics.Serializer
	metrics    []metric
	dimensions [][]string
	fields     map[string]interface{}
}

type metric struct {
	Name string
}

// NewEntry creates a new Entry with the specified namespace and serializer.
func NewEntry(namespace string, serializer metrics.Serializer) Entry {
	return Entry{
		namespace:  namespace,
		serializer: serializer,
		metrics:    []metric{},
		dimensions: [][]string{{}},
		fields:     map[string]interface{}{},
	}
}

// Build constructs the EMF log entry as a JSON string.
func (e *Entry) Build() (string, error) {

	entry := map[string]interface{}{}

	entry[emfIdentifier] = map[string]interface{}{
		timestampKey: sdk.NowTime().UnixNano() / 1e6,
		cloudWatchMetricsKey: []map[string]interface{}{
			{
				namespaceKey:  e.namespace,
				dimensionsKey: e.dimensions,
				metricsKey:    e.metrics,
			},
		},
	}

	for k, v := range e.fields {
		entry[k] = v
	}

	jsonEntry, err := e.serializer.Serialize(entry)
	if err != nil {
		return "", err
	}
	return jsonEntry, nil
}

// AddDimension adds a CW Dimension to the EMF entry.
func (e *Entry) AddDimension(key string, value string) {
	// Dimensions are a list of lists. We only support a single list.
	e.dimensions[0] = append(e.dimensions[0], key)
	e.fields[key] = value
}

// AddMetric adds a CW Metric to the EMF entry.
func (e *Entry) AddMetric(key string, value float64) {
	e.metrics = append(e.metrics, metric{key})
	e.fields[key] = value
}

// AddProperty adds a CW Property to the EMF entry.
// Properties are not published as metrics, but they are available in logs and in CW insights.
func (e *Entry) AddProperty(key string, value interface{}) {
	e.fields[key] = value
}
