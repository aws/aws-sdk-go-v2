# v1.2.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-03-02)

* **Feature**: This release fixes some input parameter range and patterns.

# v1.1.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.1.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.1.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.0.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2022-12-01)

* **Release**: New AWS service client module
* **Feature**: AWS introduces new Amazon EventBridge Pipes which allow you to connect sources (SQS, Kinesis, DDB, Kafka, MQ) to Targets (14+ EventBridge Targets) without any code, with filtering, batching, input transformation, and an optional Enrichment stage (Lambda, StepFunctions, ApiGateway, ApiDestinations)

