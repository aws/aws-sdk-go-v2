// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package emf

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

func TestCreateNewEntry(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	cases := map[string]struct {
		Namespace     string
		ExpectedEntry Entry
	}{
		"success": {
			Namespace: "testNamespace",
			ExpectedEntry: Entry{
				namespace:  "testNamespace",
				serializer: metrics.DefaultSerializer{},
				metrics:    []metric{},
				dimensions: [][]string{{}},
				fields:     map[string]interface{}{},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actualEntry := NewEntry(c.Namespace, metrics.DefaultSerializer{})
			if !reflect.DeepEqual(actualEntry, c.ExpectedEntry) {
				t.Errorf("Entry contained unexpected values")
			}
		})
	}
}

func TestBuild(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	cases := map[string]struct {
		Namespace      string
		Configure      func(entry *Entry)
		Serializer     metrics.Serializer
		ExpectedError  error
		ExpectedResult string
	}{
		"completeEntry": {
			Namespace:  "testNamespace",
			Serializer: metrics.DefaultSerializer{},
			Configure: func(entry *Entry) {
				entry.AddMetric("testMetric1", 1)
				entry.AddMetric("testMetric2", 2)
				entry.AddDimension("testDimension1", "dim1")
				entry.AddDimension("testDimension2", "dim2")
				entry.AddProperty("testProperty1", "prop1")
				entry.AddProperty("testProperty2", "prop2")
			},
			ExpectedError:  nil,
			ExpectedResult: completeEntry,
		},
		"noMetrics": {
			Namespace:  "testNamespace",
			Serializer: metrics.DefaultSerializer{},
			Configure: func(entry *Entry) {
				entry.AddDimension("testDimension1", "dim1")
				entry.AddDimension("testDimension2", "dim2")
				entry.AddProperty("testProperty1", "prop1")
				entry.AddProperty("testProperty2", "prop2")
			},
			ExpectedError:  nil,
			ExpectedResult: noMetrics,
		},
		"noDimensions": {
			Namespace:  "testNamespace",
			Serializer: metrics.DefaultSerializer{},
			Configure: func(entry *Entry) {
				entry.AddMetric("testMetric1", 1)
				entry.AddMetric("testMetric2", 2)
				entry.AddProperty("testProperty1", "prop1")
				entry.AddProperty("testProperty2", "prop2")
			},
			ExpectedError:  nil,
			ExpectedResult: noDimensions,
		},
		"noProperties": {
			Namespace:  "testNamespace",
			Serializer: metrics.DefaultSerializer{},
			Configure: func(entry *Entry) {
				entry.AddMetric("testMetric1", 1)
				entry.AddMetric("testMetric2", 2)
				entry.AddDimension("testDimension1", "dim1")
				entry.AddDimension("testDimension2", "dim2")
			},
			ExpectedError:  nil,
			ExpectedResult: noProperties,
		},
		"serializationError": {
			Namespace:  "testNamespace",
			Serializer: TestSerializerWithError{},
			Configure: func(entry *Entry) {
			},
			ExpectedError:  fmt.Errorf("serialization error"),
			ExpectedResult: "",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			entry := NewEntry(c.Namespace, c.Serializer)

			c.Configure(&entry)

			result, err := entry.Build()

			if !reflect.DeepEqual(err, c.ExpectedError) {
				t.Errorf("Unexpected error, should be '%s' but was '%s'", c.ExpectedError, err)
			}

			if !reflect.DeepEqual(result, c.ExpectedResult) {
				t.Errorf("Unexpected result, should be '%s' but was '%s'", c.ExpectedResult, result)
			}
		})
	}
}
