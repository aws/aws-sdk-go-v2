$version: "2"

namespace com.amazonaws.sdk.benchmark

use smithy.protocols#rpcv2Cbor
use aws.api#service
use aws.auth#sigv4

@title("Smithy RPC v2 CBOR Data Plane")
@sigv4(name: "smithyrpcv2cbordataplane")
@rpcv2Cbor
@service(sdkId: "RpcCborDataPlane")
service SmithyRpcV2CborDataPlane {
    version: "1999-12-31"
    operations: [Healthcheck]
    resources: [DynamoDBItem, CloudWatchMetric]
}
