# v1.5.0 (2026-09-04)

* **Feature**: Stop registering the `spanRetryLoop` middleware in generated clients. The retry loop's tracing span is now opened by the retry middleware itself.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2026-08-31.2)

* **Feature**: Stop registering the `SetCredentialSourceMiddleware` middleware in generated clients. Credential source user agent features are now set when the client's middleware stack is constructed.

# v1.3.1 (2026-08-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2026-08-27)

* **Feature**: Support connection read timeouts in the SDK. This is currently available on an opt-in basis by setting env `AWS_ENABLE_DEFAULT_SOCKET_TIMEOUT_2026=true`.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2026-08-26)

* **Feature**: Stop registering the `ComputeContentLength` middleware in generated clients. `Content-Length` is now set when the request body is set via `SetStream`.
* **Dependency Update**: Update to smithy-go v1.28.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2026-08-25)

* **Dependency Update**: Update to smithy-go v1.27.10.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2026-08-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2026-08-19)

* **Feature**: Adds throttling exceptions to operation outputs that were previously inconsistent with other operations.

# v1.0.1 (2026-08-14)

* **Dependency Update**: Update to smithy-go v1.27.8.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2026-08-11)

* **Release**: New AWS service client module
* **Feature**: Adds SDK support for AWS IAM account access manager, a feature that enables mapping of IAM roles to the users and groups in AWS IAM Identity Center.

