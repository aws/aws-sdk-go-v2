# v1.8.2 (2026-07-01)

* **Bug Fix**: Bump smithy-go to 1.27.3, fix JSON encorder for document.Number, endpoint host label format validation and CBOR union serialization on new serde
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2026-06-29)

* No change notes available for this release.

# v1.8.0 (2026-06-17)

* **Feature**: Adds support for Remote A2A (Agent-to-Agent) agent registration and management. Adds new Release Readiness Review and Release Testing capabilities. Adds support for Git managed skills in AWS DevOps Agent.

# v1.7.0 (2026-06-12)

* **Feature**: Adds support for Trigger CRUD APIs (CreateTrigger, GetTrigger, UpdateTrigger, DeleteTrigger, ListTriggers) for managing schedule-based automation triggers in DevOps Agent agent spaces.

# v1.6.0 (2026-06-08)

* **Feature**: Add Asset APIs for managing versioned assets and asset files in AWS DevOps Agent agent spaces.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.6 (2026-06-05.2)

* **Bug Fix**: Undo the initial wave of schema-serde releases due to several customer-reported regressions.

# v1.5.5 (2026-06-04)

* **Dependency Update**: Update to smithy-go v1.27.1 to fix several union-related deserialization bugs in schema-serde-enabled services.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2026-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2026-06-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2026-05-29)

* **Dependency Update**: Update to smithy-go v1.26.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2026-05-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2026-05-27)

* **Feature**: Adding new BDD representation of endpoint ruleset

# v1.4.0 (2026-05-19)

* **Feature**: Added a new serviceType mcpserversigv4 service and association. This provides feature to register MCP sigv4 authorization based MCPs

# v1.3.2 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2026-04-16)

* **Feature**: Deprecate the userId from the Chat operations. This update also removes  support of AllowVendedLogDeliveryForResource API from AWS SDKs.

# v1.2.0 (2026-04-10)

* **Feature**: Devops Agent now supports associate Splunk, Datadog and custom MCP server to an Agent Space.

# v1.1.0 (2026-03-31)

* **Feature**: AWS DevOps Agent service General Availability release.

# v1.0.0 (2026-03-30)

* **Release**: New AWS service client module
* **Feature**: AWS DevOps Agent General Availability.

