$version: "2"

namespace com.amazonaws.sdk.benchmark

list MetricData {
    member: MetricDatum
}

structure MetricDatum {
    @required
    MetricName: String
    Dimensions: Dimensions
    Timestamp: Timestamp
    Value: Double
    StatisticValues: StatisticSet
    Values: Values
    Counts: Counts
    Unit: StandardUnit
    StorageResolution: Integer
}

@length(min: 0, max: 30)
list Dimensions {
    member: Dimension
}

structure Dimension {
    @length(min: 1, max: 255)
    @required
    Name: String
    @length(min: 1, max: 1024)
    @required
    Value: String
}

structure StatisticSet {
    @required
    SampleCount: Double
    @required
    Sum: Double
    @required
    Minimum: Double
    @required
    Maximum: Double
}

list Values {
    member: Double
}

list Counts {
    member: Double
}

enum StandardUnit {
    SECONDS = "Seconds"
    MICROSECONDS = "Microseconds"
    MILLISECONDS = "Milliseconds"
    BYTES = "Bytes"
    KILOBYTES = "Kilobytes"
    MEGABYTES = "Megabytes"
    GIGABYTES = "Gigabytes"
    TERABYTES = "Terabytes"
    BITS = "Bits"
    KILOBITS = "Kilobits"
    MEGABITS = "Megabits"
    GIGABITS = "Gigabits"
    TERABITS = "Terabits"
    PERCENT = "Percent"
    COUNT = "Count"
    BYTES_SECOND = "Bytes/Second"
    KILOBYTES_SECOND = "Kilobytes/Second"
    MEGABYTES_SECOND = "Megabytes/Second"
    GIGABYTES_SECOND = "Gigabytes/Second"
    TERABYTES_SECOND = "Terabytes/Second"
    BITS_SECOND = "Bits/Second"
    KILOBITS_SECOND = "Kilobits/Second"
    MEGABITS_SECOND = "Megabits/Second"
    GIGABITS_SECOND = "Gigabits/Second"
    TERABITS_SECOND = "Terabits/Second"
    COUNT_SECOND = "Count/Second"
    NONE = "None"
}
