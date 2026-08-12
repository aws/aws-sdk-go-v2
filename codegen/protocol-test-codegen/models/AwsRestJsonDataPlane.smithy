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
    // GetObjectStreaming is restJson1-only: it exists to contrast a @streaming
    // payload against S3Object's buffered GetObject, so it hangs off the service
    // rather than the resource shared with the restXml data plane.
    operations: [Healthcheck, GetObjectStreaming]
    resources: [S3Object, CloudWatchMetric]
}
