$version: "2"

namespace com.amazonaws.sdk.benchmark

resource DynamoDBItem {
    operations: [
        PutItem, GetItem
    ]
}