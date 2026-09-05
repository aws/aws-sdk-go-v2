# v1.7.0 (2026-09-04)

* **Feature**: Stop registering the `spanRetryLoop` middleware in generated clients. The retry loop's tracing span is now opened by the retry middleware itself.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2026-08-31.2)

* **Feature**: Stop registering the `SetCredentialSourceMiddleware` middleware in generated clients. Credential source user agent features are now set when the client's middleware stack is constructed.

# v1.5.1 (2026-08-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2026-08-27)

* **Feature**: Support connection read timeouts in the SDK. This is currently available on an opt-in basis by setting env `AWS_ENABLE_DEFAULT_SOCKET_TIMEOUT_2026=true`.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2026-08-26)

* **Feature**: Stop registering the `ComputeContentLength` middleware in generated clients. `Content-Length` is now set when the request body is set via `SetStream`.
* **Dependency Update**: Update to smithy-go v1.28.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.8 (2026-08-25)

* **Dependency Update**: Update to smithy-go v1.27.10.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.7 (2026-08-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.6 (2026-08-14)

* **Dependency Update**: Update to smithy-go v1.27.8.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.5 (2026-08-10)

* **Dependency Update**: Update to smithy-go v1.27.7.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.4 (2026-08-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.3 (2026-07-31.2)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.27.6 to fix various serde issues in HTTP binding services.

# v1.3.2 (2026-07-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2026-07-28)

* **Dependency Update**: Update to smithy-go v1.27.5.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2026-07-21)

* **Feature**: Add an option to clients to disable clock skew
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2026-07-13)

* No change notes available for this release.

# v1.2.0 (2026-07-06)

* **Feature**: Add request serialization snapshot tests.

# v1.1.3 (2026-07-01)

* **Bug Fix**: Bump smithy-go to 1.27.3, fix JSON encorder for document.Number, endpoint host label format validation and CBOR union serialization on new serde
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2026-06-29)

* No change notes available for this release.

# v1.1.1 (2026-06-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2026-06-04)

* **Feature**: Adding new BDD representation of endpoint ruleset
* **Dependency Update**: Update to smithy-go v1.27.1 to fix several union-related deserialization bugs in schema-serde-enabled services.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2026-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2026-06-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2026-05-29)

* **Dependency Update**: Update to smithy-go v1.26.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2026-05-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2026-04-07)

* **Release**: New AWS service client module
* **Feature**: Support for S3 Files, a new shared file system that connects any AWS compute directly with your data in Amazon S3. It provides fast, direct access to all of your S3 data as files with full file system semantics and low-latency performance, without your data ever leaving S3.

