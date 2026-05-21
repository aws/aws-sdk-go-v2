$version: "2"

namespace com.amazonaws.sdk.benchmark

use aws.protocols#awsQuery
use smithy.test#httpRequestTests

@documentation("""
    As seen in Amazon CloudWatch.
""")
@http(method: "POST", uri: "/PutMetricData")
@httpRequestTests([
    {
        id: "awsQuery_PutMetricDataRequest_Baseline"
        protocol: awsQuery
        method: "POST"
        uri: "/"
        params: {
            Namespace: "AWS/SDK"
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsQuery_PutMetricDataRequest_S"
        protocol: awsQuery
        method: "POST"
        uri: "/"
        params: {
            Namespace: "AWS/SDK"
            MetricData: [
                { MetricName: "alpacas_found", Value: 75.0, Unit: "Percent" }
                { MetricName: "llamas_sleeping", Value: 60.0, Unit: "Percent" }
                { MetricName: "penguins_waddling", Value: 45.0, Unit: "Percent" }
            ]
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsQuery_PutMetricDataRequest_M"
        protocol: awsQuery
        method: "POST"
        uri: "/"
        params: {
            Namespace: "AWS/SDK"
            MetricData: [
                { MetricName: "alpacas_found", Value: 75.0, Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }], Timestamp: 1609459200 }
                { MetricName: "llamas_sleeping", Value: 60.0, Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }], Timestamp: 1609459200 }
                { MetricName: "penguins_waddling", Value: 45.0, Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }], Timestamp: 1609459200 }
                { MetricName: "dolphins_jumping", Value: 1024.0, Unit: "Bytes", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }], Timestamp: 1609459200 }
                { MetricName: "elephants_trumpeting", Value: 2048.0, Unit: "Bytes", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }], Timestamp: 1609459200 }
                { MetricName: "giraffes_eating", Value: 100.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }], Timestamp: 1609459200 }
                { MetricName: "zebras_running", Value: 150.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }], Timestamp: 1609459200 }
                { MetricName: "pandas_munching", Value: 50.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }], Timestamp: 1609459200 }
                { MetricName: "koalas_napping", Value: 75.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }], Timestamp: 1609459200 }
                { MetricName: "kangaroos_hopping", Value: 0.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }], Timestamp: 1609459200 }
            ]
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsQuery_PutMetricDataRequest_L"
        protocol: awsQuery
        method: "POST"
        uri: "/"
        params: {
            Namespace: "AWS/SDK"
            MetricData: [
                { MetricName: "alpacas_found", Value: 75.0, Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "llamas_sleeping", Values: [60.0, 65.0, 58.0, 62.0, 67.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "penguins_waddling", StatisticValues: { SampleCount: 10.0, Sum: 450.0, Minimum: 40.0, Maximum: 50.0 }, Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }] }
                { MetricName: "dolphins_jumping", Value: 1024.0, Unit: "Bytes", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "elephants_trumpeting", Values: [2048.0, 1800.0, 2200.0, 1950.0, 2100.0, 1750.0, 2300.0, 1900.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Bytes", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "giraffes_eating", StatisticValues: { SampleCount: 5.0, Sum: 500.0, Minimum: 80.0, Maximum: 120.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }] }
                { MetricName: "zebras_running", Value: 150.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }] }
                { MetricName: "pandas_munching", Values: [50.0, 45.0, 55.0, 48.0, 52.0, 47.0, 53.0, 49.0, 51.0, 46.0, 54.0, 50.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "koalas_napping", StatisticValues: { SampleCount: 8.0, Sum: 600.0, Minimum: 70.0, Maximum: 80.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "kangaroos_hopping", Value: 0.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "tigers_prowling", Values: [10.0, 12.0, 8.0, 11.0, 9.0, 13.0, 7.0, 14.0, 6.0, 15.0, 5.0, 16.0, 4.0, 17.0, 3.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "lions_roaring", StatisticValues: { SampleCount: 12.0, Sum: 1728.0, Minimum: 140.0, Maximum: 148.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "bears_fishing", Value: 4096.0, Unit: "Bytes", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }] }
                { MetricName: "wolves_howling", Values: [8192.0, 7500.0, 8800.0, 7200.0, 9000.0, 6800.0, 9200.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Bytes", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }] }
                { MetricName: "foxes_hunting", StatisticValues: { SampleCount: 1.0, Sum: 0.0, Minimum: 0.0, Maximum: 0.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "rabbits_hopping", Value: 25.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "VolumeId", Value: "vol-1234567890abcdef0" }] }
                { MetricName: "squirrels_gathering", Values: [35.0, 32.0, 38.0, 30.0, 40.0, 28.0, 42.0, 26.0, 44.0, 24.0, 46.0, 22.0, 48.0, 20.0, 50.0, 18.0, 52.0, 16.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "VolumeId", Value: "vol-1234567890abcdef0" }] }
                { MetricName: "owls_hooting", StatisticValues: { SampleCount: 6.0, Sum: 12288.0, Minimum: 1800.0, Maximum: 2300.0 }, Unit: "Bytes", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "VolumeId", Value: "vol-1234567890abcdef0" }] }
                { MetricName: "eagles_soaring", Value: 4096.0, Unit: "Bytes", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "VolumeId", Value: "vol-1234567890abcdef0" }] }
                { MetricName: "hawks_circling", Values: [100.0, 98.0, 102.0, 96.0, 104.0, 94.0, 106.0, 92.0, 108.0, 90.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "VolumeId", Value: "vol-1234567890abcdef0" }] }
                { MetricName: "falcons_diving", StatisticValues: { SampleCount: 4.0, Sum: 400.0, Minimum: 95.0, Maximum: 105.0 }, Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "VolumeId", Value: "vol-1234567890abcdef0" }] }
                { MetricName: "cheetahs_sprinting", Value: 5.0, Unit: "Milliseconds", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "gazelles_leaping", Values: [1000.0, 950.0, 1050.0, 900.0, 1100.0, 850.0, 1150.0, 800.0, 1200.0, 750.0, 1250.0, 700.0, 1300.0, 650.0, 1350.0, 600.0, 1400.0, 550.0, 1450.0, 500.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Bytes/Second", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "lizards_basking", StatisticValues: { SampleCount: 15.0, Sum: 975.0, Minimum: 60.0, Maximum: 70.0 }, Unit: "None", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "hummingbirds_hovering", Value: 2000.0, Unit: "Count/Second", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "bees_buzzing", Values: [150.0, 145.0, 155.0, 140.0, 160.0, 135.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "None", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "butterflies_fluttering", StatisticValues: { SampleCount: 20.0, Sum: 600.0, Minimum: 25.0, Maximum: 35.0 }, Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "ants_marching", Value: 5.0, Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "spiders_weaving", Values: [1.5, 1.3, 1.7, 1.2, 1.8, 1.1, 1.9, 1.0, 2.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "None", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "crickets_chirping", StatisticValues: { SampleCount: 10.0, Sum: 12.0, Minimum: 1.0, Maximum: 1.4 }, Unit: "None", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "fireflies_glowing", Value: 1.0, Unit: "None", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "frogs_croaking", Values: [150.0, 148.0, 152.0, 146.0, 154.0, 144.0, 156.0, 142.0, 158.0, 140.0, 160.0, 138.0, 162.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "toads_hopping", StatisticValues: { SampleCount: 5.0, Sum: 2500.0, Minimum: 480.0, Maximum: 520.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "salamanders_hiding", Value: 1000.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "turtles_swimming", Values: [50.0, 48.0, 52.0, 46.0], Counts: [1.0, 1.0, 1.0, 1.0], Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "snakes_slithering", StatisticValues: { SampleCount: 8.0, Sum: 200.0, Minimum: 20.0, Maximum: 30.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "iguanas_sunbathing", Value: 10.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "whales_singing", Values: [100.0, 95.0, 105.0, 90.0, 110.0, 85.0, 115.0, 80.0, 120.0, 75.0, 125.0, 70.0, 130.0, 65.0, 135.0, 60.0, 140.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Service", Value: "web-server" }] }
                { MetricName: "seals_barking", StatisticValues: { SampleCount: 10.0, Sum: 20.0, Minimum: 1.0, Maximum: 3.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Service", Value: "web-server" }] }
                { MetricName: "otters_playing", Value: 200.0, Unit: "Milliseconds", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Service", Value: "web-server" }] }
                { MetricName: "fish_swimming", Values: [10.0, 9.0, 11.0, 8.0, 12.0, 7.0, 13.0, 6.0, 14.0, 5.0, 15.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Database", Value: "mysql" }] }
                { MetricName: "sharks_hunting", StatisticValues: { SampleCount: 25.0, Sum: 12500.0, Minimum: 480.0, Maximum: 520.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Database", Value: "mysql" }] }
                { MetricName: "rays_gliding", Value: 50.0, Unit: "Milliseconds", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Database", Value: "mysql" }] }
                { MetricName: "octopuses_hiding", Values: [800.0, 780.0, 820.0, 760.0, 840.0, 740.0, 860.0, 720.0, 880.0, 700.0, 900.0, 680.0, 920.0, 660.0, 940.0, 640.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Cache", Value: "redis" }] }
                { MetricName: "jellyfish_floating", StatisticValues: { SampleCount: 12.0, Sum: 2400.0, Minimum: 180.0, Maximum: 220.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Cache", Value: "redis" }] }
                { MetricName: "crabs_scuttling", Value: 5.0, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Queue", Value: "sqs-queue" }] }
                { MetricName: "lobsters_crawling", Values: [100.0, 98.0, 102.0, 96.0, 104.0, 94.0, 106.0], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Queue", Value: "sqs-queue" }] }
                { MetricName: "starfish_clinging", StatisticValues: { SampleCount: 18.0, Sum: 1710.0, Minimum: 90.0, Maximum: 100.0 }, Unit: "Count", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Queue", Value: "sqs-queue" }] }
                { MetricName: "seahorses_drifting", Value: 0.5, Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "clownfish_hiding", Values: [99.5, 99.3, 99.7, 99.1, 99.9, 98.9, 99.8, 98.7, 99.6, 98.5, 99.4, 98.3, 99.2, 98.1, 99.0, 97.9, 98.8, 97.7, 98.6], Counts: [1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0], Unit: "Percent", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }
                { MetricName: "angelfish_swimming", StatisticValues: { SampleCount: 30.0, Sum: 1260.0, Minimum: 40.0, Maximum: 44.0 }, Unit: "None", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Environment", Value: "production" }] }
            ]
        }
        tags: ["serde-benchmark"]
    }
])
operation PutMetricData {
    input: PutMetricDataInput
    output: Unit
}

structure PutMetricDataInput {
    @required
    Namespace: String

    MetricData: MetricData

    EntityMetricData: EntityMetricDataList
    StrictEntityValidation: Boolean
}

list EntityMetricDataList {
    member: EntityMetricDatum
}

structure EntityMetricDatum {
    Entity: Entity
    MetricData: MetricData
}

structure Entity {
    KeyAttributes: EntityKeyAttributesMap
    Attributes: EntityAttributesMap
}

@length(min: 2, max: 4)
map EntityKeyAttributesMap {
    @length(min: 1, max: 32)
    key: String
    @length(min: 1, max: 2048)
    value: String
}

@length(min: 0, max: 10)
map EntityAttributesMap {
    @length(min: 1, max: 256)
    key: String
    @length(min: 1, max: 2048)
    value: String
}
