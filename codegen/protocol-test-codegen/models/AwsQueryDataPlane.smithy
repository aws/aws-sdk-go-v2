$version: "2"

namespace com.amazonaws.sdk.benchmark

use aws.protocols#awsQuery
use aws.api#service
use aws.auth#sigv4

@title("AWS Query Data Plane")
@sigv4(name: "awsquerydataplane")
@awsQuery
@xmlNamespace(uri: "https://awsquerydataplane.amazonaws.com")
@service(sdkId: "QueryDataPlane")
service AwsQueryDataPlane {
    version: "1999-12-31"
    operations: [Healthcheck]
    resources: [DynamoDBItem, CloudWatchMetric]
}
