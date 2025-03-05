$version: "2.0"
namespace aws.kitchensinktest

@aws.api#service(sdkId: "awsJson1 kitchen sink")
@aws.protocols#awsJson1_0
@aws.auth#sigv4(name: "awsjson1kitchensink")
service AwsJson1KitchenSink {
    version: "2025-03-01",
    operations: [GetItem],
}

operation GetItem {
    input: GetItemInput,
    output: GetItemOutput,
    errors: [ItemNotFound],
}

structure GetItemInput {
    item: Item,
}

structure GetItemOutput {}

structure Item {}

@error("client")
structure ItemNotFound {}
