# v1.11.0 (2026-08-31.2)

* **Feature**: Stop registering the `SetCredentialSourceMiddleware` middleware in generated clients. Credential source user agent features are now set when the client's middleware stack is constructed.

# v1.10.1 (2026-08-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2026-08-27)

* **Feature**: Support connection read timeouts in the SDK. This is currently available on an opt-in basis by setting env `AWS_ENABLE_DEFAULT_SOCKET_TIMEOUT_2026=true`.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2026-08-26)

* **Feature**: Stop registering the `ComputeContentLength` middleware in generated clients. `Content-Length` is now set when the request body is set via `SetStream`.
* **Dependency Update**: Update to smithy-go v1.28.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.8 (2026-08-25)

* **Dependency Update**: Update to smithy-go v1.27.10.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.7 (2026-08-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.6 (2026-08-14)

* **Dependency Update**: Update to smithy-go v1.27.8.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.5 (2026-08-10)

* **Dependency Update**: Update to smithy-go v1.27.7.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.4 (2026-08-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.3 (2026-07-31.2)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.27.6 to fix various serde issues in HTTP binding services.

# v1.8.2 (2026-07-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2026-07-28)

* **Dependency Update**: Update to smithy-go v1.27.5.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2026-07-27)

* **Feature**: Adds optional headquarters location to StartProfileUpdateTask, letting partners record their headquarters as an ISO 3166 country and subdivision code on their profile. When headquarters is provided, both the country and subdivision codes are required.

# v1.7.0 (2026-07-22)

* **Feature**: Adds Qualifications Association APIs that enable partners to associate a subsidiary account's qualifications with a primary account. Once associated, qualifications are shared across all connected accounts and scorecards are consolidated. Partners can start and track association and disassociation.

# v1.6.0 (2026-07-21)

* **Feature**: Add an option to clients to disable clock skew
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2026-07-13)

* No change notes available for this release.

# v1.5.0 (2026-07-06)

* **Feature**: Add request serialization snapshot tests.

# v1.4.9 (2026-07-01)

* **Bug Fix**: Bump smithy-go to 1.27.3, fix JSON encorder for document.Number, endpoint host label format validation and CBOR union serialization on new serde
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.8 (2026-06-29)

* No change notes available for this release.

# v1.4.7 (2026-06-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.6 (2026-06-05.2)

* **Bug Fix**: Undo the initial wave of schema-serde releases due to several customer-reported regressions.

# v1.4.5 (2026-06-04)

* **Dependency Update**: Update to smithy-go v1.27.1 to fix several union-related deserialization bugs in schema-serde-enabled services.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2026-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2026-06-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2026-05-29)

* **Dependency Update**: Update to smithy-go v1.26.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2026-05-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2026-05-22)

* **Feature**: Adding new BDD representation of endpoint ruleset

# v1.3.0 (2026-05-13)

* **Feature**: Added ServiceQuotaExceededExceptions for Profile operations

# v1.2.2 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2026-03-30)

* **Feature**: KYB Supplemental Form enables partners who fail business verification to submit additional details and supporting documentation through a self-service form, triggering an automated re-verification without requiring manual intervention from support teams.

# v1.1.6 (2026-03-26)

* **Bug Fix**: Fix a bug where a recorded clock skew could persist on the client even if the client and server clock ended up realigning.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.5 (2026-03-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2026-03-03)

* **Dependency Update**: Bump minimum Go version to 1.24
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2026-02-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2026-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2025-12-05)

* **Feature**: Adding Verification API's to Partner Central Account SDK.

# v1.0.1 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.0.0 (2025-12-01)

* **Release**: New AWS service client module
* **Feature**: Initial GA launch of Partner Central Account

