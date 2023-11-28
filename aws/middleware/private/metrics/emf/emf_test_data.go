// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package emf

import "strings"

func stripString(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return str
}

var completeEntry = stripString(`
{
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				["testDimension1", "testDimension2"]
			],
			"Metrics": [{
				"Name": "testMetric1"
			}, {
				"Name": "testMetric2"
			}],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	},
	"testDimension1": "dim1",
	"testDimension2": "dim2",
	"testMetric1": 1,
	"testMetric2": 2,
	"testProperty1": "prop1",
	"testProperty2": "prop2"
}
`)

var noMetrics = stripString(`
{
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				["testDimension1", "testDimension2"]
			],
			"Metrics": [],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	},
	"testDimension1": "dim1",
	"testDimension2": "dim2",
	"testProperty1": "prop1",
	"testProperty2": "prop2"
}
`)

var noProperties = stripString(`
{
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				["testDimension1", "testDimension2"]
			],
			"Metrics": [{
				"Name": "testMetric1"
			}, {
				"Name": "testMetric2"
			}],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	},
	"testDimension1": "dim1",
	"testDimension2": "dim2",
	"testMetric1": 1,
	"testMetric2": 2
}
`)

var noDimensions = stripString(`
{
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				[]
			],
			"Metrics": [{
				"Name": "testMetric1"
			}, {
				"Name": "testMetric2"
			}],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	},
	"testMetric1": 1,
	"testMetric2": 2,
	"testProperty1": "prop1",
	"testProperty2": "prop2"
}
`)
