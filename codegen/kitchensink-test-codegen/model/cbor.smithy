$version: "2.0"
namespace aws.kitchensinktestcbor

@aws.api#service(sdkId: "rpcv2Cbor kitchen sink")
@smithy.protocols#rpcv2Cbor
@aws.auth#sigv4(name: "rpcv2cborkitchensink")
@smithy.rules#endpointRuleSet({
    version: "1.0",
    parameters: {
        Region: {
            type: "string",
            builtIn: "AWS::Region",
            required: true,
            documentation: "The AWS region"
        }
    },
    rules: [
        {
            type: "endpoint",
            documentation: "Default regional endpoint",
            conditions: [],
            endpoint: {
                url: "https://example.{Region}.amazonaws.com"
            }
        }
    ]
})
service RpcV2CborKitchenSink {
    version: "2025-03-01",
    operations: [CborGetItem, CborPutCompressedData],
}

operation CborGetItem {
    input: CborGetItemInput,
    output: CborGetItemOutput,
    errors: [CborItemNotFound],
}

structure CborGetItemInput {
    id: String,
}

structure CborGetItemOutput {}

@error("client")
structure CborItemNotFound {}

// Request-compression operation: the request body is gzip-compressed by a
// Serialize-step middleware after serialization. Content length must be
// computed from the compressed body, so this guards against it being computed
// (inline in the serializer) from the uncompressed body.
@requestCompression(encodings: ["gzip"])
operation CborPutCompressedData {
    input: CborPutCompressedDataInput,
    output: CborPutCompressedDataOutput,
    errors: [CborItemNotFound],
}

structure CborPutCompressedDataInput {
    data: String,
}

structure CborPutCompressedDataOutput {}
