# v1.19.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2023-08-01)

* No change notes available for this release.

# v1.19.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2023-06-22)

* **Feature**: Adds support for Versions and Aliases. Adds 8 operations: PublishStateMachineVersion, DeleteStateMachineVersion, ListStateMachineVersions, CreateStateMachineAlias, DescribeStateMachineAlias, UpdateStateMachineAlias, DeleteStateMachineAlias, ListStateMachineAliases

# v1.17.13 (2023-06-15)

* No change notes available for this release.

# v1.17.12 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.11 (2023-05-04)

* No change notes available for this release.

# v1.17.10 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.9 (2023-04-10)

* No change notes available for this release.

# v1.17.8 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.7 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.6 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.17.4 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.3 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.17.2 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.1 (2023-01-23)

* No change notes available for this release.

# v1.17.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.16.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-12-01)

* **Feature**: This release adds support for the AWS Step Functions Map state in Distributed mode. The changes include a new MapRun resource and several new and modified APIs.

# v1.15.1 (2022-11-22)

* No change notes available for this release.

# v1.15.0 (2022-11-18)

* **Feature**: This release adds support for using Step Functions service integrations to invoke any cross-account AWS resource, even if that service doesn't support resource-based policies or cross-account calls. See https://docs.aws.amazon.com/step-functions/latest/dg/concepts-access-cross-acct-resources.html

# v1.14.5 (2022-11-16)

* No change notes available for this release.

# v1.14.4 (2022-11-10)

* No change notes available for this release.

# v1.14.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.16 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.15 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.14 (2022-08-30)

* No change notes available for this release.

# v1.13.13 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.12 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.11 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.10 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.9 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.8 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.7 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.6 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: Updated to latest service endpoints

# v1.8.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.

# v1.7.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-08-27)

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

