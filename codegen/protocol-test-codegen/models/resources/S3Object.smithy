$version: "2"

namespace com.amazonaws.sdk.benchmark

resource S3Object {
    operations: [
        PutObject, HeadObject, CopyObject, GetObject
    ]
}