$version: "2"

namespace com.amazonaws.sdk.benchmark

use aws.protocols#awsJson1_0
use smithy.test#httpRequestTests
use smithy.test#httpResponseTests

@documentation("""
    A response that only says "OK", if it can.
""")
@readonly
@http(method: "GET", uri: "/Healthcheck", code: 200)
@httpRequestTests([{
    id: "awsJson1_0_HealthcheckRequest_Example"
    protocol: awsJson1_0
    method: "POST"
    uri: "/"
    tags: ["serde-benchmark"]
}])
@httpResponseTests([{
    id: "awsJson1_0_HealthcheckResponse_Example"
    protocol: awsJson1_0
    code: 200
    tags: ["serde-benchmark"]
}])
operation Healthcheck {
    output := {
        ok: String
    }
}
