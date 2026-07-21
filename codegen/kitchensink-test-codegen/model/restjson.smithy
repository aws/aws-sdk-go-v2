$version: "2.0"
namespace aws.kitchensinktestrestjson

@aws.api#service(sdkId: "restJson1 kitchen sink")
@aws.protocols#restJson1
@aws.auth#sigv4(name: "restjson1kitchensink")
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
service RestJson1KitchenSink {
    version: "2025-03-01",
    operations: [GetResource, GetStreamingResource, SubscribeEvents],
}

// Non-streaming operation: response body should be closed after deserialization.
@http(method: "GET", uri: "/resource/{id}")
@readonly
operation GetResource {
    input: GetResourceInput,
    output: GetResourceOutput,
    errors: [ResourceNotFound],
}

structure GetResourceInput {
    @required
    @httpLabel
    id: String,
}

structure GetResourceOutput {
    name: String,
}

// Streaming (caller-owned) operation: the response body is handed back to the
// caller as Body and must NOT be closed on the success path.
@http(method: "GET", uri: "/resource/{id}/stream")
@readonly
operation GetStreamingResource {
    input: GetStreamingResourceInput,
    output: GetStreamingResourceOutput,
    errors: [ResourceNotFound],
}

structure GetStreamingResourceInput {
    @required
    @httpLabel
    id: String,
}

structure GetStreamingResourceOutput {
    @httpPayload
    body: StreamingPayload = "",
}

@streaming
blob StreamingPayload

@error("client")
@httpError(404)
structure ResourceNotFound {}

// Event stream (caller-owned) operation: the response body backs an event
// stream reader and must NOT be closed on the success path.
@http(method: "GET", uri: "/resource/{id}/events")
@readonly
operation SubscribeEvents {
    input: SubscribeEventsInput,
    output: SubscribeEventsOutput,
    errors: [ResourceNotFound],
}

structure SubscribeEventsInput {
    @required
    @httpLabel
    id: String,
}

structure SubscribeEventsOutput {
    @httpPayload
    events: Events,
}

@streaming
union Events {
    message: MessageEvent,
}

structure MessageEvent {
    body: String,
}
