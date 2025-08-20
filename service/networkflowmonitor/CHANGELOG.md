# v1.8.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.8.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2025-08-08)

* **Feature**: Removed incorrect endpoint tests

# v1.6.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2025-07-16.2)

* **Feature**: Introducing 2 new scope status types - DEACTIVATING and DEACTIVATED.

# v1.3.0 (2025-06-30)

* **Feature**: Add ConflictExceptions to UpdateScope and DeleteScope operations for scopes being mutated.

# v1.2.3 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2025-06-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2025-04-03)

* No change notes available for this release.

# v1.2.0 (2025-03-06)

* **Feature**: This release contains 2 changes. 1: DeleteScope/GetScope/UpdateScope operations now return 404 instead of 500 when the resource does not exist. 2: Expected string format for clientToken fields of CreateMonitorInput/CreateScopeInput/UpdateMonitorInput have been updated to be an UUID based string.

# v1.1.1 (2025-03-04.2)

* **Bug Fix**: Add assurance test for operation order.

# v1.1.0 (2025-02-27)

* **Feature**: Track credential providers via User-Agent Feature ids
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.9 (2025-02-18)

* **Bug Fix**: Bump go version to 1.22
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.8 (2025-02-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.7 (2025-01-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2025-01-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2025-01-24)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.22.2.

# v1.0.4 (2025-01-17)

* **Bug Fix**: Fix bug where credentials weren't refreshed during retry loop.

# v1.0.3 (2025-01-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2025-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2024-12-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2024-12-02)

* **Release**: New AWS service client module
* **Feature**: This release adds documentation for a new feature in Amazon CloudWatch called Network Flow Monitor. You can use Network Flow Monitor to get near real-time metrics, including retransmissions and data transferred, for your actual workloads.
* **Dependency Update**: Updated to the latest SDK module versions

