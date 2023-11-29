// This package is designated as private and is intended for use only by the
// AWS client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package publisher

import "strings"

func stripString(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return str + "\n"
}

var emptyRequestMetricData = stripString(`
[{
	"ApiCallDuration": 0,
	"ApiCallSuccessful": 0,
	"ClientRequestId": "",
	"EndpointResolutionDuration": 0,
	"HttpStatusCode":"0",
	"MarshallingDuration": 0,
	"OperationName": "",
	"RetryCount": -1,
	"ServiceId": "",
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				["ServiceId", "OperationName", "HttpStatusCode"]
			],
			"Metrics": [{
				"Name": "ApiCallDuration"
			}, {
				"Name": "ApiCallSuccessful"
			}, {
				"Name": "MarshallingDuration"
			}, {
				"Name": "EndpointResolutionDuration"
			}, {
				"Name": "RetryCount"
			}],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	}
}]
`)

var emptyStreamMetricData = stripString(`
[{
	"ClientRequestId": "",
	"HttpStatusCode": "0",
	"OperationName": "",
	"ServiceId": "",
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				["ServiceId", "OperationName", "HttpStatusCode"]
			],
			"Metrics": [],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	}
}]
`)

var completeStreamMetricData = stripString(`
[{
	"ClientRequestId": "crid",
	"HttpStatusCode": "200",
	"OperationName": "operationname",
	"ServiceId": "sid",
	"Throughput": 80000000,
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				["ServiceId", "OperationName", "HttpStatusCode"]
			],
			"Metrics": [{
				"Name": "Throughput"
			}],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	}
}]
`)

var completeRequestMetricData = stripString(`
[{
	"ApiCallDuration": 200000000000,
	"ApiCallSuccessful": 1,
	"ClientRequestId": "crid",
	"EndpointResolutionDuration": 200000000000,
	"HttpStatusCode":"200",
	"InThroughput": 0.5,
	"MarshallingDuration": 200000000000,
	"OperationName": "operationname",
	"OutThroughput": 0.5,
	"RetryCount": 1,
	"ServiceId": "sid",
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				["ServiceId", "OperationName", "HttpStatusCode"]
			],
			"Metrics": [{
				"Name": "ApiCallDuration"
			}, {
				"Name": "ApiCallSuccessful"
			}, {
				"Name": "MarshallingDuration"
			}, {
				"Name": "EndpointResolutionDuration"
			}, {
				"Name": "RetryCount"
			}, {
				"Name": "InThroughput"
			}, {
				"Name": "OutThroughput"
			}],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	}
}]
`)

var completeMetricDataAttempt1 = stripString(`
[{
	"AttemptNumber": 0,
	"AvailableConcurrency": 2,
	"AwsExtendedRequestId": "exreqid",
	"AwsRequestId": "reqid",
	"BackoffDelayDuration":100,
	"ClientRequestId": "crid",
	"ConcurrencyAcquireDuration": 200000000000,
	"HttpStatusCode": "200",
	"MaxConcurrency":10,
	"OperationName": "operationname",
	"PendingConcurrencyAcquires": 1,
	"ServiceCallDuration": 200000000000,
	"ServiceId": "sid",
	"SigningDuration": 200000000000,
	"TimeToFirstByte": 0,
	"UnmarshallingDuration": 200000000000,
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				["ServiceId", "OperationName", "HttpStatusCode"]
			],
			"Metrics": [{
				"Name": "MaxConcurrency"
			}, {
				"Name": "AvailableConcurrency"
			}, {
				"Name": "ConcurrencyAcquireDuration"
			}, {
				"Name": "PendingConcurrencyAcquires"
			}, {
				"Name": "SigningDuration"
			}, {
				"Name": "UnmarshallingDuration"
			}, {
				"Name": "TimeToFirstByte"
			}, {
				"Name": "ServiceCallDuration"
			}, {
				"Name":"BackoffDelayDuration"
			}],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	}
}]
`)

var completeMetricDataAttempt2 = stripString(`
[{
	"AttemptNumber": 1,
	"AvailableConcurrency": 2,
	"AwsExtendedRequestId": "exreqid",
    "AwsRequestId": "reqid",
	"BackoffDelayDuration":100,
	"ClientRequestId": "crid",
	"ConcurrencyAcquireDuration": 200000000000,
	"HttpStatusCode": "200",
	"MaxConcurrency":10,
	"OperationName": "operationname",
	"PendingConcurrencyAcquires": 1,
	"ServiceCallDuration": 200000000000,
	"ServiceId": "sid",
	"SigningDuration": 200000000000,
	"TimeToFirstByte": 0,
	"UnmarshallingDuration": 200000000000,
	"_aws": {
		"CloudWatchMetrics": [{
			"Dimensions": [
				["ServiceId", "OperationName", "HttpStatusCode"]
			],
			"Metrics": [{
				"Name": "MaxConcurrency"
			}, {
				"Name": "AvailableConcurrency"
			}, {
				"Name": "ConcurrencyAcquireDuration"
			}, {
				"Name": "PendingConcurrencyAcquires"
			}, {
				"Name": "SigningDuration"
			}, {
				"Name": "UnmarshallingDuration"
			}, {
				"Name": "TimeToFirstByte"
			}, {
				"Name": "ServiceCallDuration"
			}, {
				"Name": "BackoffDelayDuration"
			}],
			"Namespace": "testNamespace"
		}],
		"Timestamp": 1234000
	}
}]
`)
