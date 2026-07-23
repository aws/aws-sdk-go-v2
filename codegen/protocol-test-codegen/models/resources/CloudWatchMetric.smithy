$version: "2"

namespace com.amazonaws.sdk.benchmark

resource CloudWatchMetric {
    operations: [
        PutMetricData, GetMetricData
    ]
}