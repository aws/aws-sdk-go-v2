$version: "2"

namespace com.amazonaws.sdk.benchmark

use aws.protocols#restXml
use aws.api#service
use aws.auth#sigv4

@title("AWS REST XML Data Plane")
@sigv4(name: "awsrestxmldataplane")
@restXml
@xmlNamespace(uri: "https://awsrestxmldataplane.amazonaws.com")
@service(sdkId: "RestXmlDataPlane")
service AwsRestXmlDataPlane {
    version: "1999-12-31"
    operations: [Healthcheck]
    resources: [S3Object, CloudWatchMetric]
}
