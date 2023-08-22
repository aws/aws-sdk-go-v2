# v1.29.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.1 (2023-08-01)

* No change notes available for this release.

# v1.29.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.5 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.4 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.3 (2023-07-07)

* No change notes available for this release.

# v1.28.2 (2023-06-15)

* No change notes available for this release.

# v1.28.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2023-05-04)

* **Feature**: This release adds support for the Suricata REJECT option in midstream exception configurations.

# v1.27.0 (2023-05-03)

* **Feature**: AWS Network Firewall now supports policy level HOME_NET variable overrides.

# v1.26.3 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.2 (2023-04-10)

* No change notes available for this release.

# v1.26.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2023-04-05)

* **Feature**: AWS Network Firewall now supports IPv6-only subnets.

# v1.25.0 (2023-03-30)

* **Feature**: AWS Network Firewall added TLS inspection configurations to allow TLS traffic inspection.

# v1.24.7 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.6 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.5 (2023-03-03)

* No change notes available for this release.

# v1.24.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.24.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.24.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2023-01-17)

* **Feature**: Network Firewall now allows creation of dual stack endpoints, enabling inspection of IPv6 traffic.

# v1.23.0 (2023-01-09)

* **Feature**: Network Firewall now supports the Suricata rule action reject, in addition to the actions pass, drop, and alert.

# v1.22.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.21.0 (2022-12-28)

* **Feature**: AWS Network Firewall now provides status messages for firewalls to help you troubleshoot when your endpoint fails.

# v1.20.9 (2022-12-23)

* No change notes available for this release.

# v1.20.8 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.7 (2022-12-08)

* No change notes available for this release.

# v1.20.6 (2022-12-06)

* No change notes available for this release.

# v1.20.5 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.4 (2022-11-11)

* No change notes available for this release.

# v1.20.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-10-18)

* No change notes available for this release.

# v1.20.0 (2022-10-05)

* **Feature**: StreamExceptionPolicy configures how AWS Network Firewall processes traffic when a network connection breaks midstream

# v1.19.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.7 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.6 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.5 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-07-21)

* **Feature**: Network Firewall now supports referencing dynamic IP sets from stateful rule groups, for IP sets stored in Amazon VPC prefix lists.

# v1.17.4 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.3 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.2 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-04-28)

* **Feature**: AWS Network Firewall adds support for stateful threat signature AWS managed rule groups.

# v1.16.0 (2022-04-26)

* **Feature**: AWS Network Firewall now enables customers to use a customer managed AWS KMS key for the encryption of their firewall resources.

# v1.15.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: API client updated

# v1.10.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.

# v1.9.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-09-30)

* **Feature**: API client updated

# v1.6.0 (2021-09-17)

* **Feature**: Updated API client and endpoints to latest revision.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

