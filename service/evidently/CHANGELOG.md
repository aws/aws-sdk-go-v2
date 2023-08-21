# v1.12.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.1 (2023-08-01)

* No change notes available for this release.

# v1.12.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.10 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.9 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.8 (2023-06-15)

* No change notes available for this release.

# v1.11.7 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.6 (2023-05-04)

* No change notes available for this release.

# v1.11.5 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.4 (2023-04-10)

* No change notes available for this release.

# v1.11.3 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2023-03-09)

* **Feature**: Updated entity override documentation

# v1.10.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.10.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.10.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.9.5 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.4 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Feature**: This release adds support for the client-side evaluation - powered by AWS AppConfig feature.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.7 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.6 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.5 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2022-07-15)

* **Feature**: This release adds support for the new segmentation feature.

# v1.7.4 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2022-05-09)

* **Feature**: Add detail message inside GetExperimentResults API response to indicate experiment result availability

# v1.6.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.1.1 (2021-12-03)

* **Bug Fix**: Fixed a bug that prevented the resolution of the correct endpoint for some API operations.

# v1.1.0 (2021-12-02)

* **Feature**: API client updated
* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2021-11-30)

* **Release**: New AWS service client module

