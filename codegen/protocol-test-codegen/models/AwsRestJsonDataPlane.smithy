$version: "2"

namespace com.amazonaws.sdk.benchmark

use aws.protocols#restJson1
use aws.api#service
use aws.auth#sigv4

@title("AWS REST JSON Data Plane")
@sigv4(name: "awsrestjsondataplane")
@restJson1
@service(sdkId: "RestJsonDataPlane")
service AwsRestJsonDataPlane {
    version: "1999-12-31"
    operations: [Healthcheck]
    resources: [S3Object, CloudWatchMetric]
}
