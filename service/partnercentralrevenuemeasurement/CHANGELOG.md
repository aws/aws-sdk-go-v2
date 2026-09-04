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

# v1.1.8 (2026-08-25)

* **Dependency Update**: Update to smithy-go v1.27.10.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.7 (2026-08-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.6 (2026-08-14)

* **Dependency Update**: Update to smithy-go v1.27.8.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.5 (2026-08-10)

* **Dependency Update**: Update to smithy-go v1.27.7.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2026-08-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2026-07-31.2)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.27.6 to fix various serde issues in HTTP binding services.

# v1.1.2 (2026-07-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2026-07-28)

* **Dependency Update**: Update to smithy-go v1.27.5.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2026-07-21)

* **Feature**: Add an option to clients to disable clock skew
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2026-07-13)

* No change notes available for this release.

# v1.0.0 (2026-07-07)

* **Release**: New AWS service client module
* **Feature**: Add support for AWS Partner Central Revenue Measurement API for creating, managing, and tracking revenue attributions and marketplace revenue share allocations.

