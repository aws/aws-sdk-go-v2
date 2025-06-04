$version: "2.0"
namespace smithy.protocoltests.rpcv2Cbor

@aws.api#service(sdkId: "RpcV2CborQueryCompatible")
@smithy.protocols#rpcv2Cbor
@aws.auth#sigv4(name: "RpcV2CborQueryCompatible")
@aws.protocols#awsQueryCompatible
service RpcV2CborQueryCompatible {
    version: "2099-01-01",
    operations: [GetItem],
}

operation GetItem {
    input: GetItemInput,
    output: GetItemOutput,
    errors: [ItemNotFound],
}

structure GetItemInput {}

structure GetItemOutput {}

@error("client")
@aws.protocols#awsQueryError(
    code: "aws.protocolstests.json#ItemNotFound"
    httpResponseCode: 404
)
structure ItemNotFound {}
