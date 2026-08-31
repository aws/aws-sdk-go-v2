# v1.5.0 (2026-08-31.2)

* **Feature**: Stop registering the `SetCredentialSourceMiddleware` middleware in generated clients. Credential source user agent features are now set when the client's middleware stack is constructed.

# v1.4.1 (2026-08-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2026-08-27)

* **Feature**: Support connection read timeouts in the SDK. This is currently available on an opt-in basis by setting env `AWS_ENABLE_DEFAULT_SOCKET_TIMEOUT_2026=true`.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2026-08-26)

* **Feature**: Stop registering the `ComputeContentLength` middleware in generated clients. `Content-Length` is now set when the request body is set via `SetStream`.
* **Dependency Update**: Update to smithy-go v1.28.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.8 (2026-08-25)

* **Dependency Update**: Update to smithy-go v1.27.10.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.7 (2026-08-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.6 (2026-08-14)

* **Dependency Update**: Update to smithy-go v1.27.8.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2026-08-10)

* **Dependency Update**: Update to smithy-go v1.27.7.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2026-08-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2026-07-31.2)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.27.6 to fix various serde issues in HTTP binding services.

# v1.2.2 (2026-07-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2026-07-28)

* **Dependency Update**: Update to smithy-go v1.27.5.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2026-07-21)

* **Feature**: Add an option to clients to disable clock skew
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2026-07-13)

* No change notes available for this release.

# v1.1.0 (2026-07-06)

* **Feature**: Add request serialization snapshot tests.

# v1.0.2 (2026-07-01)

* **Bug Fix**: Bump smithy-go to 1.27.3, fix JSON encorder for document.Number, endpoint host label format validation and CBOR union serialization on new serde
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2026-06-29)

* No change notes available for this release.

# v1.0.0 (2026-06-22)

* **Release**: New AWS service client module
* **Feature**: Initial release of the AWS Lambda Core SDK with APIs to create, manage, and tag network connectors that enable Lambda compute resources to access private resources in your Amazon VPC.

