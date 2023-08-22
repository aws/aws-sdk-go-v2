# v1.15.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.1 (2023-08-01)

* No change notes available for this release.

# v1.15.0 (2023-07-31)

* **Feature**: Add support for in-aws right sizing
* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.4 (2023-07-28.2)

* No change notes available for this release.

# v1.14.3 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.2 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.1 (2023-06-15)

* No change notes available for this release.

# v1.14.0 (2023-06-13)

* **Feature**: Added APIs to support network replication and recovery using AWS Elastic Disaster Recovery.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.2 (2023-05-04)

* No change notes available for this release.

# v1.13.1 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2023-04-17)

* **Feature**: Changed existing APIs and added new APIs to support using an account-level launch configuration template with AWS Elastic Disaster Recovery.

# v1.12.2 (2023-04-10)

* No change notes available for this release.

# v1.12.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2023-03-30)

* **Feature**: Adding a field to the replication configuration APIs to support the auto replicate new disks feature. We also deprecated RetryDataReplication.

# v1.11.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2023-02-27)

* **Feature**: New fields were added to reflect availability zone data in source server and recovery instance description commands responses, as well as source server launch status.

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

# v1.9.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2022-11-28)

* **Feature**: Non breaking changes to existing APIs, and additional APIs added to support in-AWS failing back using AWS Elastic Disaster Recovery.

# v1.8.4 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.3 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2022-09-21)

* No change notes available for this release.

# v1.8.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2022-09-14)

* **Feature**: Fixed the data type of lagDuration that is returned in Describe Source Server API
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.7 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.6 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.5 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2022-07-15)

* **Feature**: Changed existing APIs to allow choosing a dynamic volume type for replicating volumes, to reduce costs for customers.

# v1.6.4 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2022-06-16)

* No change notes available for this release.

# v1.6.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2022-05-31)

* **Feature**: Changed existing APIs and added new APIs to accommodate using multiple AWS accounts with AWS Elastic Disaster Recovery.

# v1.5.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.0.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2021-11-19)

* **Release**: New AWS service client module
* **Dependency Update**: Updated to the latest SDK module versions

