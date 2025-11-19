# v1.9.1 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2025-11-19)

* **Feature**: Added name filtering support to ListBillingViews API through the new names parameter to efficiently filter billing views by name.

# v1.8.6 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.8.5 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.8.4 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.8.3 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2025-09-26)

* **Feature**: Add ability to combine custom billing views to create new consolidated views.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.6 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.5 (2025-09-10)

* No change notes available for this release.

# v1.7.4 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2025-08-20)

* **Feature**: Clarify IPv4 and IPv6 endpoints
* **Bug Fix**: Remove unused deserialization code.

# v1.6.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2025-08-08)

* **Feature**: Removed incorrect endpoint tests

# v1.4.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2025-06-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2025-04-03)

* No change notes available for this release.

# v1.2.1 (2025-03-04.2)

* **Bug Fix**: Add assurance test for operation order.

# v1.2.0 (2025-02-27)

* **Feature**: Track credential providers via User-Agent Feature ids
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.8 (2025-02-18)

* **Bug Fix**: Bump go version to 1.22
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.7 (2025-02-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.6 (2025-01-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.5 (2025-01-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2025-01-24)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.22.2.

# v1.1.3 (2025-01-17)

* **Bug Fix**: Fix bug where credentials weren't refreshed during retry loop.

# v1.1.2 (2025-01-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2025-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2024-12-20)

* **Feature**: Added new API's for defining and fetching Billing Views.

# v1.0.3 (2024-12-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2024-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2024-11-18)

* **Dependency Update**: Update to smithy-go v1.22.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2024-11-13)

* **Release**: New AWS service client module
* **Feature**: Today, AWS announces the general availability of ListBillingViews API in the AWS SDKs, to enable AWS Billing Conductor (ABC) users to create proforma Cost and Usage Reports (CUR) programmatically.

