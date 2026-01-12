# v1.15.17 (2026-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.16 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.15 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.15.14 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.15.13 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.12 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.15.11 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.15.10 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.15.9 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.8 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.7 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.6 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.5 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.4 (2025-09-10)

* No change notes available for this release.

# v1.15.3 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.2 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.1 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2025-08-25)

* **Feature**: Remove incorrect endpoint tests

# v1.14.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.14.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.4 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.3 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.2 (2025-06-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2025-04-03)

* No change notes available for this release.

# v1.11.0 (2025-03-10)

* **Feature**: PrivateCA Connector for Active Directory now supports dual stack endpoints. This release adds the IpAddressType option to the VpcInformation on a Connector which determines whether the endpoint supports IPv4 only or IPv4 and IPv6 traffic.

# v1.10.1 (2025-03-04.2)

* **Bug Fix**: Add assurance test for operation order.

# v1.10.0 (2025-02-27)

* **Feature**: Track credential providers via User-Agent Feature ids
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.16 (2025-02-18)

* **Bug Fix**: Bump go version to 1.22
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.15 (2025-02-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.14 (2025-01-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.13 (2025-01-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.12 (2025-01-24)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.22.2.

# v1.9.11 (2025-01-17)

* **Bug Fix**: Fix bug where credentials weren't refreshed during retry loop.

# v1.9.10 (2025-01-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.9 (2025-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.8 (2024-12-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.7 (2024-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.6 (2024-11-18)

* **Dependency Update**: Update to smithy-go v1.22.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.5 (2024-11-07)

* **Bug Fix**: Adds case-insensitive handling of error message fields in service responses

# v1.9.4 (2024-11-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.3 (2024-10-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2024-10-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2024-10-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2024-10-04)

* **Feature**: Add support for HTTP client metrics.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.4 (2024-10-03)

* No change notes available for this release.

# v1.8.3 (2024-09-27)

* No change notes available for this release.

# v1.8.2 (2024-09-25)

* No change notes available for this release.

# v1.8.1 (2024-09-23)

* No change notes available for this release.

# v1.8.0 (2024-09-20)

* **Feature**: Add tracing and metrics support to service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.7 (2024-09-17)

* **Bug Fix**: **BREAKFIX**: Only generate AccountIDEndpointMode config for services that use it. This is a compiler break, but removes no actual functionality, as no services currently use the account ID in endpoint resolution.

# v1.7.6 (2024-09-04)

* No change notes available for this release.

# v1.7.5 (2024-09-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.4 (2024-08-15)

* **Dependency Update**: Bump minimum Go version to 1.21.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2024-07-10.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2024-07-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2024-06-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2024-06-26)

* **Feature**: Support list-of-string endpoint parameter.

# v1.6.1 (2024-06-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2024-06-18)

* **Feature**: Track usage of various AWS SDK features in user-agent string.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.11 (2024-06-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.10 (2024-06-07)

* **Bug Fix**: Add clock skew correction on all service clients
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.9 (2024-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.8 (2024-05-23)

* No change notes available for this release.

# v1.5.7 (2024-05-16)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.6 (2024-05-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.5 (2024-05-08)

* **Bug Fix**: GoDoc improvement

# v1.5.4 (2024-03-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2024-03-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2024-03-07)

* **Bug Fix**: Remove dependency on go-cmp.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2024-02-23)

* **Bug Fix**: Move all common, SDK-side middleware stack ops into the service client module to prevent cross-module compatibility issues in the future.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.4.2 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.4.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.6 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.3.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.3.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.2.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2023-08-30)

* **Release**: New AWS service client module
* **Feature**: The Connector for AD allows you to use a fully-managed AWS Private CA as a drop-in replacement for your self-managed enterprise CAs without local agents or proxy servers. Enterprises that use AD to manage Windows environments can reduce their private certificate authority (CA) costs and complexity.

