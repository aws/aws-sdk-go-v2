$version: "2.0"
namespace aws.protocoltests.json

@aws.api#service(sdkId: "Json10QueryCompatible")
@aws.protocols#awsJson1_0
@aws.auth#sigv4(name: "Json10QueryCompatible")
@aws.protocols#awsQueryCompatible
service Json10QueryCompatible {
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
