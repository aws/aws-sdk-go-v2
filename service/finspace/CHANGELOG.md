# v1.12.0 (2023-08-21)

* **Feature**: Allow customers to manage outbound traffic from their Kx Environment when attaching a transit gateway by providing network acl entries. Allow the customer to choose how they want to update the databases on a cluster allowing updates to possibly be faster than usual.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2023-08-01)

* No change notes available for this release.

# v1.11.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.4 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.2 (2023-06-15)

* No change notes available for this release.

# v1.10.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2023-06-05)

* **Feature**: Releasing new Managed kdb Insights APIs

# v1.9.10 (2023-05-04)

* No change notes available for this release.

# v1.9.9 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.8 (2023-04-10)

* No change notes available for this release.

# v1.9.7 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.6 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.9.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.9.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.8.21 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.20 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.19 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.18 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.17 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.16 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.15 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.14 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.13 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.12 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.11 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.10 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.9 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.8 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.7 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.6 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2021-05-14)

* **Release**: New AWS service client module
* **Feature**: Updated to latest service API model.
* **Dependency Update**: Updated to the latest SDK module versions

