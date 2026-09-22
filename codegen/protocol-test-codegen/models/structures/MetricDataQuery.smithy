$version: "2"

namespace com.amazonaws.sdk.benchmark

list MetricDataQueries {
    member: MetricDataQuery
}

structure MetricDataQuery {
    @length(min: 1, max: 255)
    @required
    @clientOptional
    Id: String

    MetricStat: MetricStat
    @length(min: 1, max: 2048)
    Expression: String
    Label: String
    ReturnData: Boolean
    @range(min: 1)
    Period: Integer
    AccountId: String
}

structure MetricStat {
    @required
    @clientOptional
    Metric: Metric
    @required
    @clientOptional
    Period: Integer
    @required
    @clientOptional
    Stat: String
    Unit: StandardUnit
}

structure Metric {
    Namespace: String
    MetricName: String
    Dimensions: Dimensions
}

list MetricDataResults {
    member: MetricDataResult
}

structure MetricDataResult {
    Id: String
    Label: String
    Timestamps: Timestamps
    Values: Values
    StatusCode: StatusCode
    Messages: MetricDataResultMessages
}

list Timestamps {
    member: Timestamp
}

enum ScanBy {
    TIMESTAMP_DESCENDING = "TimestampDescending"
    TIMESTAMP_ASCENDING = "TimestampAscending"
}

structure LabelOptions {
    Timezone: String
}

enum StatusCode {
    COMPLETE = "Complete"
    INTERNAL_ERROR = "InternalError"
    PARTIAL_DATA = "PartialData"
    FORBIDDEN = "Forbidden"
}

list MetricDataResultMessages {
    member: MessageData
}

structure MessageData {
    Code: String
    Value: String
}