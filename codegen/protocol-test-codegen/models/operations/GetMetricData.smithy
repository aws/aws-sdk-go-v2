$version: "2"

namespace com.amazonaws.sdk.benchmark

use aws.protocols#awsQuery
use smithy.test#httpRequestTests
use smithy.test#httpResponseTests

@documentation("""
    As seen in Amazon CloudWatch
""")
@http(method: "POST", uri: "/GetMetricData")
@httpRequestTests([
    {
        id: "awsQuery_GetMetricDataRequest_S"
        protocol: awsQuery
        method: "POST"
        uri: "/"
        params: {
            MetricDataQueries: [
                { Id: "m1", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "alpacas_found" }, Period: 300, Stat: "Average" } }
                { Id: "m2", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "llamas_sleeping" }, Period: 300, Stat: "Maximum" } }
                { Id: "m3", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "penguins_waddling" }, Period: 300, Stat: "Minimum" } }
            ]
            StartTime: 1609459200
            EndTime: 1609462800
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsQuery_GetMetricDataRequest_M"
        protocol: awsQuery
        method: "POST"
        uri: "/"
        params: {
            MetricDataQueries: [
                { Id: "m1", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "alpacas_found", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Average" } }
                { Id: "m2", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "llamas_sleeping", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Maximum" } }
                { Id: "m3", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "penguins_waddling", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m4", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "dolphins_jumping", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Average" } }
                { Id: "m5", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "elephants_trumpeting", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m6", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "giraffes_eating", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m7", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "zebras_running", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }, { Name: "Device", Value: "/dev/xvda1" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m8", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "pandas_munching", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m9", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "koalas_napping", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m10", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "kangaroos_hopping", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Maximum" } }
            ]
            StartTime: 1609459200
            EndTime: 1609462800
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsQuery_GetMetricDataRequest_L"
        protocol: awsQuery
        method: "POST"
        uri: "/"
        params: {
            MetricDataQueries: [
                { Id: "m1", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "alpacas_found", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Average" } }
                { Id: "m2", Expression: "m1 * 100", Label: "alpacas_found_percent" }
                { Id: "m3", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "llamas_sleeping" }, Period: 300, Stat: "Maximum" }, ReturnData: false }
                { Id: "m4", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "penguins_waddling", Dimensions: [{ Name: "Device", Value: "/dev/xvda1" }] }, Period: 60, Stat: "Sum" } }
                { Id: "m5", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "dolphins_jumping" }, Period: 300, Stat: "Average", Unit: "Bytes" } }
                { Id: "m6", Expression: "ANOMALY_DETECTION_FUNCTION(m5, 2)", Label: "dolphins_jumping_anomaly" }
                { Id: "m7", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "elephants_trumpeting", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m8", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "giraffes_eating" }, Period: 300, Stat: "Sum" }, AccountId: "123456789012" }
                { Id: "m9", Expression: "m7 + m8", Label: "combined_animal_activity", ReturnData: false }
                { Id: "m10", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "zebras_running", Dimensions: [{ Name: "Device", Value: "/dev/xvda1" }] }, Period: 60, Stat: "Maximum" } }
                { Id: "m11", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "pandas_munching" }, Period: 300, Stat: "Average" } }
                { Id: "m12", Expression: "IF(m11 > 50, 1, 0)", Label: "high_panda_activity" }
                { Id: "m13", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "koalas_napping", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m14", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "kangaroos_hopping" }, Period: 300, Stat: "Maximum" }, ReturnData: false }
                { Id: "m15", Expression: "RATE(m13)", Label: "koala_nap_rate" }
                { Id: "m16", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "tigers_prowling", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Average" } }
                { Id: "m17", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "lions_roaring" }, Period: 300, Stat: "Minimum" } }
                { Id: "m18", Expression: "m16 / m17", Label: "big_cat_ratio" }
                { Id: "m19", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "otters_swimming", Dimensions: [{ Name: "Device", Value: "/dev/xvda1" }] }, Period: 60, Stat: "Sum" } }
                { Id: "m20", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "seals_clapping" }, Period: 60, Stat: "Sum" }, AccountId: "123456789012" }
                { Id: "m21", Expression: "(m19 + m20) / 1024", Label: "aquatic_mammals_total" }
                { Id: "m22", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "flamingos_standing", Dimensions: [{ Name: "InstanceId", Value: "i-1234567890abcdef0" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m23", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "parrots_squawking" }, Period: 300, Stat: "Sum" } }
                { Id: "m24", Expression: "SEARCH('{AWS/SDK,InstanceId} MetricName=\"alpacas_found\"', 'Average', 300)", Label: "all_alpacas" }
                { Id: "m25", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "toucans_flying", Dimensions: [{ Name: "VolumeId", Value: "vol-1234567890abcdef0" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m26", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "owls_hooting" }, Period: 300, Stat: "Average" }, ReturnData: false }
                { Id: "m27", Expression: "m25 * 4096", Label: "estimated_toucan_bytes" }
                { Id: "m28", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "eagles_soaring", Dimensions: [{ Name: "VolumeId", Value: "vol-1234567890abcdef0" }] }, Period: 300, Stat: "Sum" } }
                { Id: "m29", MetricStat: { Metric: { Namespace: "AWS/SDK", MetricName: "hawks_circling" }, Period: 300, Stat: "Sum" } }
                { Id: "m30", Expression: "m29 / m23", Label: "avg_bird_latency" }
            ]
            StartTime: 1609459200
            EndTime: 1609462800
            MaxDatapoints: 1440
            ScanBy: "TimestampDescending"
            LabelOptions: { Timezone: "UTC" }
        }
        tags: ["serde-benchmark"]
    }
])
@httpResponseTests([
    {
        id: "awsQuery_GetMetricDataResponse_S"
        protocol: awsQuery
        code: 200
        body: """
<GetMetricDataResponse
    xmlns="https://awsquerydataplane.amazonaws.com">
    <GetMetricDataResult>
        <MetricDataResults>
            <member>
                <Id>m1</Id>
                <Label>alpacas_found</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>75</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m2</Id>
                <Label>llamas_sleeping</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>60</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m3</Id>
                <Label>penguins_waddling</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>45</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
        </MetricDataResults>
    </GetMetricDataResult>
    <ResponseMetadata>
        <RequestId>12345678-1234-1234-1234-123456789013</RequestId>
    </ResponseMetadata>
</GetMetricDataResponse>
        """
        params: {
            MetricDataResults: [
                { Id: "m1", Label: "alpacas_found", Timestamps: [1609459200], Values: [75.0], StatusCode: "Complete" }
                { Id: "m2", Label: "llamas_sleeping", Timestamps: [1609459200], Values: [60.0], StatusCode: "Complete" }
                { Id: "m3", Label: "penguins_waddling", Timestamps: [1609459200], Values: [45.0], StatusCode: "Complete" }
            ]
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsQuery_GetMetricDataResponse_M"
        protocol: awsQuery
        code: 200
        body: """
<GetMetricDataResponse
    xmlns="https://awsquerydataplane.amazonaws.com">
    <GetMetricDataResult>
        <MetricDataResults>
            <member>
                <Id>m1</Id>
                <Label>alpacas_found</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>75</member>
                    <member>72.5</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m2</Id>
                <Label>llamas_sleeping</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>60</member>
                    <member>58</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m3</Id>
                <Label>penguins_waddling</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>45</member>
                    <member>47</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m4</Id>
                <Label>dolphins_jumping</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>1024</member>
                    <member>1100</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m5</Id>
                <Label>elephants_trumpeting</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>2048</member>
                    <member>2200</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m6</Id>
                <Label>giraffes_eating</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>100</member>
                    <member>95</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m7</Id>
                <Label>zebras_running</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>150</member>
                    <member>145</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m8</Id>
                <Label>pandas_munching</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>50</member>
                    <member>48</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m9</Id>
                <Label>koalas_napping</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>75</member>
                    <member>72</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m10</Id>
                <Label>kangaroos_hopping</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>0</member>
                    <member>0</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
        </MetricDataResults>
    </GetMetricDataResult>
    <ResponseMetadata>
        <RequestId>12345678-1234-1234-1234-123456789014</RequestId>
    </ResponseMetadata>
</GetMetricDataResponse>
        """
        params: {
            MetricDataResults: [
                { Id: "m1", Label: "alpacas_found", Timestamps: [1609459200, 1609459500], Values: [75.0, 72.5], StatusCode: "Complete" }
                { Id: "m2", Label: "llamas_sleeping", Timestamps: [1609459200, 1609459500], Values: [60.0, 58.0], StatusCode: "Complete" }
                { Id: "m3", Label: "penguins_waddling", Timestamps: [1609459200, 1609459500], Values: [45.0, 47.0], StatusCode: "Complete" }
                { Id: "m4", Label: "dolphins_jumping", Timestamps: [1609459200, 1609459500], Values: [1024.0, 1100.0], StatusCode: "Complete" }
                { Id: "m5", Label: "elephants_trumpeting", Timestamps: [1609459200, 1609459500], Values: [2048.0, 2200.0], StatusCode: "Complete" }
                { Id: "m6", Label: "giraffes_eating", Timestamps: [1609459200, 1609459500], Values: [100.0, 95.0], StatusCode: "Complete" }
                { Id: "m7", Label: "zebras_running", Timestamps: [1609459200, 1609459500], Values: [150.0, 145.0], StatusCode: "Complete" }
                { Id: "m8", Label: "pandas_munching", Timestamps: [1609459200, 1609459500], Values: [50.0, 48.0], StatusCode: "Complete" }
                { Id: "m9", Label: "koalas_napping", Timestamps: [1609459200, 1609459500], Values: [75.0, 72.0], StatusCode: "Complete" }
                { Id: "m10", Label: "kangaroos_hopping", Timestamps: [1609459200, 1609459500], Values: [0.0, 0.0], StatusCode: "Complete" }
            ]
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsQuery_GetMetricDataResponse_L"
        protocol: awsQuery
        code: 200
        body: """
<GetMetricDataResponse
        xmlns="https://awsquerydataplane.amazonaws.com">
    <GetMetricDataResult>
        <MetricDataResults>
            <member>
                <Id>m1</Id>
                <Label>alpacas_found</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>75</member>
                    <member>72.5</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m2</Id>
                <Label>alpacas_found_percent</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>7500</member>
                    <member>7250</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m4</Id>
                <Label>penguins_waddling</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>450</member>
                </Values>
                <StatusCode>PartialData</StatusCode>
                <Messages>
                    <member>
                        <Code>InternalError</Code>
                        <Value>Penguin data partially unavailable due to ice storm</Value>
                    </member>
                </Messages>
            </member>
            <member>
                <Id>m5</Id>
                <Label>dolphins_jumping</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>1024</member>
                    <member>1100</member>
                    <member>980</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m6</Id>
                <Label>dolphins_jumping_anomaly</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>0</member>
                    <member>1</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m7</Id>
                <Label>elephants_trumpeting</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>2048</member>
                    <member>2200</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m8</Id>
                <Label>giraffes_eating</Label>
                <Timestamps/>
                <Values/>
                <StatusCode>InternalError</StatusCode>
                <Messages>
                    <member>
                        <Code>InternalError</Code>
                        <Value>Giraffe feeding schedule access denied</Value>
                    </member>
                </Messages>
            </member>
            <member>
                <Id>m10</Id>
                <Label>zebras_running</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>150</member>
                </Values>
                <StatusCode>Forbidden</StatusCode>
                <Messages>
                    <member>
                        <Code>AccessDenied</Code>
                        <Value>Zebra tracking permissions insufficient</Value>
                    </member>
                </Messages>
            </member>
            <member>
                <Id>m11</Id>
                <Label>pandas_munching</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>50</member>
                    <member>48</member>
                    <member>52</member>
                    <member>49</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m12</Id>
                <Label>high_panda_activity</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>1</member>
                    <member>0</member>
                    <member>1</member>
                    <member>0</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m13</Id>
                <Label>koalas_napping</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                </Timestamps>
                <Values>
                    <member>75</member>
                    <member>72</member>
                    <member>78</member>
                    <member>74</member>
                    <member>76</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m15</Id>
                <Label>koala_nap_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.25</member>
                    <member>0.24</member>
                    <member>0.26</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m16</Id>
                <Label>tigers_prowling</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                </Timestamps>
                <Values>
                    <member>10</member>
                    <member>12</member>
                    <member>8</member>
                    <member>11</member>
                    <member>9</member>
                    <member>13</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m17</Id>
                <Label>lions_roaring</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>144</member>
                    <member>142</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m18</Id>
                <Label>big_cat_ratio</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.069</member>
                    <member>0.085</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m19</Id>
                <Label>otters_swimming</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                    <member>2021-01-01T00:30:00Z</member>
                </Timestamps>
                <Values>
                    <member>4096</member>
                    <member>4200</member>
                    <member>3900</member>
                    <member>4100</member>
                    <member>4050</member>
                    <member>4150</member>
                    <member>4000</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m20</Id>
                <Label>seals_clapping</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>8192</member>
                </Values>
                <StatusCode>PartialData</StatusCode>
            </member>
            <member>
                <Id>m21</Id>
                <Label>aquatic_mammals_total</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>12</member>
                </Values>
                <StatusCode>PartialData</StatusCode>
            </member>
            <member>
                <Id>m22</Id>
                <Label>flamingos_standing</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                    <member>2021-01-01T00:30:00Z</member>
                    <member>2021-01-01T00:35:00Z</member>
                </Timestamps>
                <Values>
                    <member>0</member>
                    <member>0</member>
                    <member>0</member>
                    <member>0</member>
                    <member>0</member>
                    <member>0</member>
                    <member>0</member>
                    <member>0</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m23</Id>
                <Label>parrots_squawking</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>25</member>
                    <member>23</member>
                    <member>27</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m24</Id>
                <Label>all_alpacas</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>68.5</member>
                    <member>71.2</member>
                    <member>69.8</member>
                    <member>70.1</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m25</Id>
                <Label>toucans_flying</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>35</member>
                    <member>32</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m27</Id>
                <Label>estimated_toucan_bytes</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>143360</member>
                    <member>131072</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m28</Id>
                <Label>eagles_soaring</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                </Timestamps>
                <Values>
                    <member>4096</member>
                    <member>4200</member>
                    <member>3800</member>
                    <member>4300</member>
                    <member>4000</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m29</Id>
                <Label>hawks_circling</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                    <member>2021-01-01T00:30:00Z</member>
                    <member>2021-01-01T00:35:00Z</member>
                    <member>2021-01-01T00:40:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.025</member>
                    <member>0.023</member>
                    <member>0.027</member>
                    <member>0.024</member>
                    <member>0.026</member>
                    <member>0.025</member>
                    <member>0.028</member>
                    <member>0.022</member>
                    <member>0.024</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>m30</Id>
                <Label>avg_bird_latency</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.001</member>
                    <member>0.001</member>
                    <member>0.001</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>r1</Id>
                <Label>requests_from_bees</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                </Timestamps>
                <Values>
                    <member>1000</member>
                    <member>1050</member>
                    <member>980</member>
                    <member>1020</member>
                    <member>1100</member>
                    <member>990</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>r2</Id>
                <Label>bee_request_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>3.33</member>
                    <member>3.5</member>
                    <member>3.27</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>r3</Id>
                <Label>butterfly_response_time</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.125</member>
                    <member>0.132</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>r4</Id>
                <Label>ant_success_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>99.2</member>
                    <member>99.5</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>r5</Id>
                <Label>spider_4xx_errors</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>5</member>
                    <member>3</member>
                    <member>7</member>
                    <member>4</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>r6</Id>
                <Label>beetle_5xx_errors</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>2</member>
                    <member>1</member>
                    <member>3</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>r7</Id>
                <Label>insect_error_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.7</member>
                    <member>0.4</member>
                    <member>1</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>d1</Id>
                <Label>fish_swimming_speed</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                </Timestamps>
                <Values>
                    <member>45.2</member>
                    <member>47.8</member>
                    <member>44.1</member>
                    <member>46.5</member>
                    <member>48.2</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>d2</Id>
                <Label>shark_connections</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>15</member>
                    <member>17</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>d3</Id>
                <Label>whale_song_anomaly</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>0</member>
                    <member>1</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>d4</Id>
                <Label>octopus_response_time</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                    <member>2021-01-01T00:30:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.002</member>
                    <member>0.0025</member>
                    <member>0.0018</member>
                    <member>0.0022</member>
                    <member>0.0024</member>
                    <member>0.0019</member>
                    <member>0.0021</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>d5</Id>
                <Label>max_sea_creature_latency</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.0025</member>
                    <member>0.0028</member>
                    <member>0.0023</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>l1</Id>
                <Label>firefly_invocations</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>250</member>
                    <member>280</member>
                    <member>220</member>
                    <member>260</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>l2</Id>
                <Label>moth_duration</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                </Timestamps>
                <Values>
                    <member>125.5</member>
                    <member>132.8</member>
                    <member>118.2</member>
                    <member>128.9</member>
                    <member>135.1</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>l3</Id>
                <Label>total_bug_execution_time</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>31375</member>
                    <member>37184</member>
                    <member>26004</member>
                    <member>33514</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>l4</Id>
                <Label>cricket_errors</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>2</member>
                    <member>1</member>
                    <member>3</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>l5</Id>
                <Label>grasshopper_throttles</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>0</member>
                    <member>1</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>l6</Id>
                <Label>bug_error_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.8</member>
                    <member>0.7</member>
                    <member>1.4</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>s1</Id>
                <Label>acorn_storage_bytes</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>1073741824</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>s2</Id>
                <Label>nuts_collected</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>1024</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>s3</Id>
                <Label>avg_acorn_size</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                </Timestamps>
                <Values>
                    <member>1048576</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>dy1</Id>
                <Label>squirrel_read_capacity</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                </Timestamps>
                <Values>
                    <member>50</member>
                    <member>55</member>
                    <member>48</member>
                    <member>52</member>
                    <member>58</member>
                    <member>47</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>dy2</Id>
                <Label>total_chipmunk_capacity</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>85</member>
                    <member>92</member>
                    <member>78</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>dy3</Id>
                <Label>hamster_throttled_requests</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>0</member>
                    <member>1</member>
                    <member>0</member>
                    <member>2</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>sq1</Id>
                <Label>rabbit_messages_visible</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                </Timestamps>
                <Values>
                    <member>25</member>
                    <member>28</member>
                    <member>22</member>
                    <member>30</member>
                    <member>26</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>sq2</Id>
                <Label>bunny_message_growth_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.01</member>
                    <member>0.02</member>
                    <member>-0.01</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>sq3</Id>
                <Label>hare_messages_sent</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                    <member>2021-01-01T00:30:00Z</member>
                </Timestamps>
                <Values>
                    <member>15</member>
                    <member>18</member>
                    <member>12</member>
                    <member>20</member>
                    <member>16</member>
                    <member>14</member>
                    <member>19</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>sq4</Id>
                <Label>cottontail_messages_received</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                </Timestamps>
                <Values>
                    <member>12</member>
                    <member>15</member>
                    <member>14</member>
                    <member>18</member>
                    <member>13</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>sq5</Id>
                <Label>rabbit_message_backlog</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>3</member>
                    <member>3</member>
                    <member>-2</member>
                    <member>2</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>sn1</Id>
                <Label>pigeon_notifications_sent</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>100</member>
                    <member>105</member>
                    <member>98</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>sn2</Id>
                <Label>dove_notification_failure_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>1</member>
                    <member>0.95</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cf1</Id>
                <Label>falcon_requests</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                    <member>2021-01-01T00:30:00Z</member>
                    <member>2021-01-01T00:35:00Z</member>
                </Timestamps>
                <Values>
                    <member>5000</member>
                    <member>5200</member>
                    <member>4800</member>
                    <member>5100</member>
                    <member>5300</member>
                    <member>4900</member>
                    <member>5050</member>
                    <member>5150</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cf2</Id>
                <Label>eagle_bytes_downloaded</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>10485760</member>
                    <member>10737418</member>
                    <member>10223616</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cf3</Id>
                <Label>avg_bird_response_size</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>2097.15</member>
                    <member>2065.66</member>
                    <member>2129.92</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cf4</Id>
                <Label>crow_4xx_error_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.5</member>
                    <member>0.4</member>
                    <member>0.6</member>
                    <member>0.45</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cf5</Id>
                <Label>raven_5xx_error_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.1</member>
                    <member>0.15</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cf6</Id>
                <Label>total_bird_error_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.6</member>
                    <member>0.55</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>ag1</Id>
                <Label>monkey_api_count</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                </Timestamps>
                <Values>
                    <member>800</member>
                    <member>850</member>
                    <member>780</member>
                    <member>820</member>
                    <member>870</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>ag2</Id>
                <Label>gorilla_p95_latency</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>250</member>
                    <member>275</member>
                    <member>230</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>ag3</Id>
                <Label>chimp_4xx_errors</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                </Timestamps>
                <Values>
                    <member>8</member>
                    <member>6</member>
                    <member>10</member>
                    <member>7</member>
                    <member>9</member>
                    <member>5</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>ag4</Id>
                <Label>orangutan_5xx_errors</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>2</member>
                    <member>1</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>ag5</Id>
                <Label>primate_error_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>1.25</member>
                    <member>0.82</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>ec1</Id>
                <Label>cheetah_cpu_utilization</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                    <member>2021-01-01T00:30:00Z</member>
                </Timestamps>
                <Values>
                    <member>55.2</member>
                    <member>58.7</member>
                    <member>52.1</member>
                    <member>56.8</member>
                    <member>59.3</member>
                    <member>54.6</member>
                    <member>57.2</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>ec2</Id>
                <Label>leopard_memory_utilization</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>62.5</member>
                    <member>65.8</member>
                    <member>60.2</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>ec3</Id>
                <Label>max_wildcat_utilization</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>62.5</member>
                    <member>65.8</member>
                    <member>60.2</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>el1</Id>
                <Label>sloth_cpu_utilization</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>25.8</member>
                    <member>28.2</member>
                    <member>23.5</member>
                    <member>26.9</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>el2</Id>
                <Label>armadillo_cache_misses</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                </Timestamps>
                <Values>
                    <member>50</member>
                    <member>45</member>
                    <member>55</member>
                    <member>48</member>
                    <member>52</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>el3</Id>
                <Label>anteater_cache_hit_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>92.5</member>
                    <member>94.2</member>
                    <member>90.8</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>k1</Id>
                <Label>salmon_incoming_records</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                </Timestamps>
                <Values>
                    <member>1000</member>
                    <member>1100</member>
                    <member>950</member>
                    <member>1050</member>
                    <member>1150</member>
                    <member>980</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>k2</Id>
                <Label>trout_outgoing_records</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>980</member>
                    <member>1080</member>
                    <member>940</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>k3</Id>
                <Label>fish_record_backlog</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>20</member>
                    <member>20</member>
                    <member>10</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>rs1</Id>
                <Label>whale_cpu_utilization</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>35.8</member>
                    <member>38.2</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>rs2</Id>
                <Label>dolphin_connection_anomaly</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                </Timestamps>
                <Values>
                    <member>0</member>
                    <member>0</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cw1</Id>
                <Label>mole_disk_used_percent</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                    <member>2021-01-01T00:30:00Z</member>
                    <member>2021-01-01T00:35:00Z</member>
                    <member>2021-01-01T00:40:00Z</member>
                    <member>2021-01-01T00:45:00Z</member>
                </Timestamps>
                <Values>
                    <member>75.2</member>
                    <member>75.8</member>
                    <member>76.1</member>
                    <member>76.5</member>
                    <member>76.9</member>
                    <member>77.2</member>
                    <member>77.6</member>
                    <member>78</member>
                    <member>78.3</member>
                    <member>78.7</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cw2</Id>
                <Label>badger_mem_used_percent</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>68.5</member>
                    <member>69.2</member>
                    <member>67.8</member>
                    <member>70.1</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cw3</Id>
                <Label>groundhog_resource_alert</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                </Timestamps>
                <Values>
                    <member>0</member>
                    <member>0</member>
                    <member>0</member>
                    <member>0</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cw4</Id>
                <Label>prairie_dog_tcp_connections</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>125</member>
                    <member>132</member>
                    <member>118</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cw5</Id>
                <Label>gopher_processes_total</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                </Timestamps>
                <Values>
                    <member>245</member>
                    <member>248</member>
                    <member>242</member>
                    <member>250</member>
                    <member>247</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>cw6</Id>
                <Label>woodchuck_process_growth_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>0.01</member>
                    <member>-0.02</member>
                    <member>0.03</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>u1</Id>
                <Label>owl_api_call_count</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                    <member>2021-01-01T00:15:00Z</member>
                    <member>2021-01-01T00:20:00Z</member>
                    <member>2021-01-01T00:25:00Z</member>
                    <member>2021-01-01T00:30:00Z</member>
                </Timestamps>
                <Values>
                    <member>500</member>
                    <member>520</member>
                    <member>480</member>
                    <member>510</member>
                    <member>530</member>
                    <member>490</member>
                    <member>515</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
            <member>
                <Id>u2</Id>
                <Label>nightingale_api_call_rate</Label>
                <Timestamps>
                    <member>2021-01-01T00:00:00Z</member>
                    <member>2021-01-01T00:05:00Z</member>
                    <member>2021-01-01T00:10:00Z</member>
                </Timestamps>
                <Values>
                    <member>1.67</member>
                    <member>1.73</member>
                    <member>1.6</member>
                </Values>
                <StatusCode>Complete</StatusCode>
            </member>
        </MetricDataResults>
        <NextToken>AQICAHhQdAFQVGGp</NextToken>
        <Messages>
            <member>
                <Code>PartialData</Code>
                <Value>Some animal metrics could not be retrieved due to migration season</Value>
            </member>
        </Messages>
    </GetMetricDataResult>
    <ResponseMetadata>
        <RequestId>12345678-1234-1234-1234-123456789015</RequestId>
    </ResponseMetadata>
</GetMetricDataResponse>

"""
        params: {
            MetricDataResults: [
                { Id: "m1", Label: "alpacas_found", Timestamps: [1609459200, 1609459500], Values: [75.0, 72.5], StatusCode: "Complete" }
                { Id: "m2", Label: "alpacas_found_percent", Timestamps: [1609459200, 1609459500], Values: [7500.0, 7250.0], StatusCode: "Complete" }
                { Id: "m4", Label: "penguins_waddling", Timestamps: [1609459200], Values: [450.0], StatusCode: "PartialData", Messages: [{ Code: "InternalError", Value: "Penguin data partially unavailable due to ice storm" }] }
                { Id: "m5", Label: "dolphins_jumping", Timestamps: [1609459200, 1609459500, 1609459800], Values: [1024.0, 1100.0, 980.0], StatusCode: "Complete" }
                { Id: "m6", Label: "dolphins_jumping_anomaly", Timestamps: [1609459200, 1609459500], Values: [0.0, 1.0], StatusCode: "Complete" }
                { Id: "m7", Label: "elephants_trumpeting", Timestamps: [1609459200, 1609459500], Values: [2048.0, 2200.0], StatusCode: "Complete" }
                { Id: "m8", Label: "giraffes_eating", Timestamps: [], Values: [], StatusCode: "InternalError", Messages: [{ Code: "InternalError", Value: "Giraffe feeding schedule access denied" }] }
                { Id: "m10", Label: "zebras_running", Timestamps: [1609459200], Values: [150.0], StatusCode: "Forbidden", Messages: [{ Code: "AccessDenied", Value: "Zebra tracking permissions insufficient" }] }
                { Id: "m11", Label: "pandas_munching", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [50.0, 48.0, 52.0, 49.0], StatusCode: "Complete" }
                { Id: "m12", Label: "high_panda_activity", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [1.0, 0.0, 1.0, 0.0], StatusCode: "Complete" }
                { Id: "m13", Label: "koalas_napping", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400], Values: [75.0, 72.0, 78.0, 74.0, 76.0], StatusCode: "Complete" }
                { Id: "m15", Label: "koala_nap_rate", Timestamps: [1609459200, 1609459500, 1609459800], Values: [0.25, 0.24, 0.26], StatusCode: "Complete" }
                { Id: "m16", Label: "tigers_prowling", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700], Values: [10.0, 12.0, 8.0, 11.0, 9.0, 13.0], StatusCode: "Complete" }
                { Id: "m17", Label: "lions_roaring", Timestamps: [1609459200, 1609459500], Values: [144.0, 142.0], StatusCode: "Complete" }
                { Id: "m18", Label: "big_cat_ratio", Timestamps: [1609459200, 1609459500], Values: [0.069, 0.085], StatusCode: "Complete" }
                { Id: "m19", Label: "otters_swimming", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700, 1609461000], Values: [4096.0, 4200.0, 3900.0, 4100.0, 4050.0, 4150.0, 4000.0], StatusCode: "Complete" }
                { Id: "m20", Label: "seals_clapping", Timestamps: [1609459200], Values: [8192.0], StatusCode: "PartialData" }
                { Id: "m21", Label: "aquatic_mammals_total", Timestamps: [1609459200], Values: [12.0], StatusCode: "PartialData" }
                { Id: "m22", Label: "flamingos_standing", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700, 1609461000, 1609461300], Values: [0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0], StatusCode: "Complete" }
                { Id: "m23", Label: "parrots_squawking", Timestamps: [1609459200, 1609459500, 1609459800], Values: [25.0, 23.0, 27.0], StatusCode: "Complete" }
                { Id: "m24", Label: "all_alpacas", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [68.5, 71.2, 69.8, 70.1], StatusCode: "Complete" }
                { Id: "m25", Label: "toucans_flying", Timestamps: [1609459200, 1609459500], Values: [35.0, 32.0], StatusCode: "Complete" }
                { Id: "m27", Label: "estimated_toucan_bytes", Timestamps: [1609459200, 1609459500], Values: [143360.0, 131072.0], StatusCode: "Complete" }
                { Id: "m28", Label: "eagles_soaring", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400], Values: [4096.0, 4200.0, 3800.0, 4300.0, 4000.0], StatusCode: "Complete" }
                { Id: "m29", Label: "hawks_circling", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700, 1609461000, 1609461300, 1609461600], Values: [0.025, 0.023, 0.027, 0.024, 0.026, 0.025, 0.028, 0.022, 0.024], StatusCode: "Complete" }
                { Id: "m30", Label: "avg_bird_latency", Timestamps: [1609459200, 1609459500, 1609459800], Values: [0.001, 0.001, 0.001], StatusCode: "Complete" }
                { Id: "r1", Label: "requests_from_bees", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700], Values: [1000.0, 1050.0, 980.0, 1020.0, 1100.0, 990.0], StatusCode: "Complete" }
                { Id: "r2", Label: "bee_request_rate", Timestamps: [1609459200, 1609459500, 1609459800], Values: [3.33, 3.50, 3.27], StatusCode: "Complete" }
                { Id: "r3", Label: "butterfly_response_time", Timestamps: [1609459200, 1609459500], Values: [0.125, 0.132], StatusCode: "Complete" }
                { Id: "r4", Label: "ant_success_rate", Timestamps: [1609459200, 1609459500], Values: [99.2, 99.5], StatusCode: "Complete" }
                { Id: "r5", Label: "spider_4xx_errors", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [5.0, 3.0, 7.0, 4.0], StatusCode: "Complete" }
                { Id: "r6", Label: "beetle_5xx_errors", Timestamps: [1609459200, 1609459500, 1609459800], Values: [2.0, 1.0, 3.0], StatusCode: "Complete" }
                { Id: "r7", Label: "insect_error_rate", Timestamps: [1609459200, 1609459500, 1609459800], Values: [0.7, 0.4, 1.0], StatusCode: "Complete" }
                { Id: "d1", Label: "fish_swimming_speed", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400], Values: [45.2, 47.8, 44.1, 46.5, 48.2], StatusCode: "Complete" }
                { Id: "d2", Label: "shark_connections", Timestamps: [1609459200, 1609459500], Values: [15.0, 17.0], StatusCode: "Complete" }
                { Id: "d3", Label: "whale_song_anomaly", Timestamps: [1609459200, 1609459500], Values: [0.0, 1.0], StatusCode: "Complete" }
                { Id: "d4", Label: "octopus_response_time", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700, 1609461000], Values: [0.002, 0.0025, 0.0018, 0.0022, 0.0024, 0.0019, 0.0021], StatusCode: "Complete" }
                { Id: "d5", Label: "max_sea_creature_latency", Timestamps: [1609459200, 1609459500, 1609459800], Values: [0.0025, 0.0028, 0.0023], StatusCode: "Complete" }
                { Id: "l1", Label: "firefly_invocations", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [250.0, 280.0, 220.0, 260.0], StatusCode: "Complete" }
                { Id: "l2", Label: "moth_duration", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400], Values: [125.5, 132.8, 118.2, 128.9, 135.1], StatusCode: "Complete" }
                { Id: "l3", Label: "total_bug_execution_time", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [31375.0, 37184.0, 26004.0, 33514.0], StatusCode: "Complete" }
                { Id: "l4", Label: "cricket_errors", Timestamps: [1609459200, 1609459500, 1609459800], Values: [2.0, 1.0, 3.0], StatusCode: "Complete" }
                { Id: "l5", Label: "grasshopper_throttles", Timestamps: [1609459200, 1609459500], Values: [0.0, 1.0], StatusCode: "Complete" }
                { Id: "l6", Label: "bug_error_rate", Timestamps: [1609459200, 1609459500, 1609459800], Values: [0.8, 0.7, 1.4], StatusCode: "Complete" }
                { Id: "s1", Label: "acorn_storage_bytes", Timestamps: [1609459200], Values: [1073741824.0], StatusCode: "Complete" }
                { Id: "s2", Label: "nuts_collected", Timestamps: [1609459200], Values: [1024.0], StatusCode: "Complete" }
                { Id: "s3", Label: "avg_acorn_size", Timestamps: [1609459200], Values: [1048576.0], StatusCode: "Complete" }
                { Id: "dy1", Label: "squirrel_read_capacity", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700], Values: [50.0, 55.0, 48.0, 52.0, 58.0, 47.0], StatusCode: "Complete" }
                { Id: "dy2", Label: "total_chipmunk_capacity", Timestamps: [1609459200, 1609459500, 1609459800], Values: [85.0, 92.0, 78.0], StatusCode: "Complete" }
                { Id: "dy3", Label: "hamster_throttled_requests", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [0.0, 1.0, 0.0, 2.0], StatusCode: "Complete" }
                { Id: "sq1", Label: "rabbit_messages_visible", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400], Values: [25.0, 28.0, 22.0, 30.0, 26.0], StatusCode: "Complete" }
                { Id: "sq2", Label: "bunny_message_growth_rate", Timestamps: [1609459200, 1609459500, 1609459800], Values: [0.01, 0.02, -0.01], StatusCode: "Complete" }
                { Id: "sq3", Label: "hare_messages_sent", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700, 1609461000], Values: [15.0, 18.0, 12.0, 20.0, 16.0, 14.0, 19.0], StatusCode: "Complete" }
                { Id: "sq4", Label: "cottontail_messages_received", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400], Values: [12.0, 15.0, 14.0, 18.0, 13.0], StatusCode: "Complete" }
                { Id: "sq5", Label: "rabbit_message_backlog", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [3.0, 3.0, -2.0, 2.0], StatusCode: "Complete" }
                { Id: "sn1", Label: "pigeon_notifications_sent", Timestamps: [1609459200, 1609459500, 1609459800], Values: [100.0, 105.0, 98.0], StatusCode: "Complete" }
                { Id: "sn2", Label: "dove_notification_failure_rate", Timestamps: [1609459200, 1609459500], Values: [1.0, 0.95], StatusCode: "Complete" }
                { Id: "cf1", Label: "falcon_requests", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700, 1609461000, 1609461300], Values: [5000.0, 5200.0, 4800.0, 5100.0, 5300.0, 4900.0, 5050.0, 5150.0], StatusCode: "Complete" }
                { Id: "cf2", Label: "eagle_bytes_downloaded", Timestamps: [1609459200, 1609459500, 1609459800], Values: [10485760.0, 10737418.0, 10223616.0], StatusCode: "Complete" }
                { Id: "cf3", Label: "avg_bird_response_size", Timestamps: [1609459200, 1609459500, 1609459800], Values: [2097.15, 2065.66, 2129.92], StatusCode: "Complete" }
                { Id: "cf4", Label: "crow_4xx_error_rate", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [0.5, 0.4, 0.6, 0.45], StatusCode: "Complete" }
                { Id: "cf5", Label: "raven_5xx_error_rate", Timestamps: [1609459200, 1609459500], Values: [0.1, 0.15], StatusCode: "Complete" }
                { Id: "cf6", Label: "total_bird_error_rate", Timestamps: [1609459200, 1609459500], Values: [0.6, 0.55], StatusCode: "Complete" }
                { Id: "ag1", Label: "monkey_api_count", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400], Values: [800.0, 850.0, 780.0, 820.0, 870.0], StatusCode: "Complete" }
                { Id: "ag2", Label: "gorilla_p95_latency", Timestamps: [1609459200, 1609459500, 1609459800], Values: [250.0, 275.0, 230.0], StatusCode: "Complete" }
                { Id: "ag3", Label: "chimp_4xx_errors", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700], Values: [8.0, 6.0, 10.0, 7.0, 9.0, 5.0], StatusCode: "Complete" }
                { Id: "ag4", Label: "orangutan_5xx_errors", Timestamps: [1609459200, 1609459500], Values: [2.0, 1.0], StatusCode: "Complete" }
                { Id: "ag5", Label: "primate_error_rate", Timestamps: [1609459200, 1609459500], Values: [1.25, 0.82], StatusCode: "Complete" }
                { Id: "ec1", Label: "cheetah_cpu_utilization", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700, 1609461000], Values: [55.2, 58.7, 52.1, 56.8, 59.3, 54.6, 57.2], StatusCode: "Complete" }
                { Id: "ec2", Label: "leopard_memory_utilization", Timestamps: [1609459200, 1609459500, 1609459800], Values: [62.5, 65.8, 60.2], StatusCode: "Complete" }
                { Id: "ec3", Label: "max_wildcat_utilization", Timestamps: [1609459200, 1609459500, 1609459800], Values: [62.5, 65.8, 60.2], StatusCode: "Complete" }
                { Id: "el1", Label: "sloth_cpu_utilization", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [25.8, 28.2, 23.5, 26.9], StatusCode: "Complete" }
                { Id: "el2", Label: "armadillo_cache_misses", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400], Values: [50.0, 45.0, 55.0, 48.0, 52.0], StatusCode: "Complete" }
                { Id: "el3", Label: "anteater_cache_hit_rate", Timestamps: [1609459200, 1609459500, 1609459800], Values: [92.5, 94.2, 90.8], StatusCode: "Complete" }
                { Id: "k1", Label: "salmon_incoming_records", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700], Values: [1000.0, 1100.0, 950.0, 1050.0, 1150.0, 980.0], StatusCode: "Complete" }
                { Id: "k2", Label: "trout_outgoing_records", Timestamps: [1609459200, 1609459500, 1609459800], Values: [980.0, 1080.0, 940.0], StatusCode: "Complete" }
                { Id: "k3", Label: "fish_record_backlog", Timestamps: [1609459200, 1609459500, 1609459800], Values: [20.0, 20.0, 10.0], StatusCode: "Complete" }
                { Id: "rs1", Label: "whale_cpu_utilization", Timestamps: [1609459200, 1609459500], Values: [35.8, 38.2], StatusCode: "Complete" }
                { Id: "rs2", Label: "dolphin_connection_anomaly", Timestamps: [1609459200, 1609459500], Values: [0.0, 0.0], StatusCode: "Complete" }
                { Id: "cw1", Label: "mole_disk_used_percent", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700, 1609461000, 1609461300, 1609461600, 1609461900], Values: [75.2, 75.8, 76.1, 76.5, 76.9, 77.2, 77.6, 78.0, 78.3, 78.7], StatusCode: "Complete" }
                { Id: "cw2", Label: "badger_mem_used_percent", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [68.5, 69.2, 67.8, 70.1], StatusCode: "Complete" }
                { Id: "cw3", Label: "groundhog_resource_alert", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100], Values: [0.0, 0.0, 0.0, 0.0], StatusCode: "Complete" }
                { Id: "cw4", Label: "prairie_dog_tcp_connections", Timestamps: [1609459200, 1609459500, 1609459800], Values: [125.0, 132.0, 118.0], StatusCode: "Complete" }
                { Id: "cw5", Label: "gopher_processes_total", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400], Values: [245.0, 248.0, 242.0, 250.0, 247.0], StatusCode: "Complete" }
                { Id: "cw6", Label: "woodchuck_process_growth_rate", Timestamps: [1609459200, 1609459500, 1609459800], Values: [0.01, -0.02, 0.03], StatusCode: "Complete" }
                { Id: "u1", Label: "owl_api_call_count", Timestamps: [1609459200, 1609459500, 1609459800, 1609460100, 1609460400, 1609460700, 1609461000], Values: [500.0, 520.0, 480.0, 510.0, 530.0, 490.0, 515.0], StatusCode: "Complete" }
                { Id: "u2", Label: "nightingale_api_call_rate", Timestamps: [1609459200, 1609459500, 1609459800], Values: [1.67, 1.73, 1.60], StatusCode: "Complete" }
            ]
            NextToken: "AQICAHhQdAFQVGGp"
            Messages: [{ Code: "PartialData", Value: "Some animal metrics could not be retrieved due to migration season" }]
        }
        tags: ["serde-benchmark"]
    }
])
operation GetMetricData {
    input: GetMetricDataInput
    output: GetMetricDataOutput
}

structure GetMetricDataInput {
    @required
    MetricDataQueries: MetricDataQueries
    @required
    StartTime: Timestamp
    @required
    EndTime: Timestamp

    NextToken: String
    ScanBy: ScanBy
    MaxDatapoints: Integer
    LabelOptions: LabelOptions
}

structure GetMetricDataOutput {
    MetricDataResults: MetricDataResults
    NextToken: String
    Messages: MetricDataResultMessages
}

