$version: "2.0"
namespace aws.kitchensinktest

@aws.api#service(sdkId: "awsJson1 kitchen sink")
@aws.protocols#awsJson1_0
@aws.auth#sigv4(name: "awsjson1kitchensink")
@smithy.rules#endpointRuleSet({
    version: "1.0",
    parameters: {
        Region: {
            type: "string",
            builtIn: "AWS::Region",
            required: true,
            documentation: "The AWS region"
        },
        Id: {
            type: "string",
            documentation: "The item id"
        }
    },
    rules: [
        {
            type: "endpoint",
            documentation: "id-specific endpoint",
            conditions: [
                { fn: "isSet", argv: [{ ref: "Id" }] }
            ],
            endpoint: {
                url: "https://{Id}.example.{Region}.amazonaws.com"
            }
        },
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
@smithy.rules#endpointTests({
    version: "1.0",
    testCases: [
        {
            documentation: "id-specific endpoint",
            params: {
                Region: "us-east-1",
                Id: "some-id"
            },
            operationInputs: [{
                operationName: "GetItem",
                "builtInParams": {
                    "AWS::Region": "us-east-1"
                }
                operationParams: {
                    id: "some-id",
                }
            }],
            expect: {
                endpoint: {
                    url: "https://some-id.example.us-east-1.amazonaws.com"
                }
            }
        },
        {
            documentation: "Default endpoint when id is not set",
            params: {
                Region: "us-west-2"
            },
            expect: {
                endpoint: {
                    url: "https://example.us-west-2.amazonaws.com"
                }
            }
        }
    ]
})
@smithy.rules#endpointBdd(
    version: "1.1"
    parameters: {
        Region: {
            builtIn: "AWS::Region"
            required: true
            documentation: "The AWS region"
            type: "string"
        }
        Id: {
            required: false
            documentation: "The item id"
            type: "string"
        }
    }
    conditions: [
        {
            fn: "isSet"
            argv: [
                {
                    ref: "Id"
                }
            ]
        }
    ]
    results: [
        {
            conditions: []
            endpoint: {
                url: "https://{Id}.example.{Region}.amazonaws.com"
                properties: {}
                headers: {}
            }
            type: "endpoint"
        }
        {
            documentation: "Default regional endpoint"
            conditions: []
            endpoint: {
                url: "https://example.{Region}.amazonaws.com"
                properties: {}
                headers: {}
            }
            type: "endpoint"
        }
    ]
    root: 2
    nodeCount: 2
    nodes: "/////wAAAAH/////AAAAAAX14QEF9eEC"
)
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
    @smithy.rules#contextParam(name: "Id")
    id: String,
}

structure GetItemOutput {}

structure Item {}

@error("client")
structure ItemNotFound {}
