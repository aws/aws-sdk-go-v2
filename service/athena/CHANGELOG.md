# v1.31.6 (2023-08-24)

* No change notes available for this release.

# v1.31.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.1 (2023-08-01)

* No change notes available for this release.

# v1.31.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.5 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.4 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.3 (2023-07-03)

* No change notes available for this release.

# v1.30.2 (2023-06-15)

* No change notes available for this release.

# v1.30.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-06-08)

* **Feature**: You can now define custom spark properties at start of the session for use cases like cluster encryption, table formats, and general Spark tuning.

# v1.29.0 (2023-06-02)

* **Feature**: This release introduces the DeleteCapacityReservation API and the ability to manage capacity reservations using CloudFormation

# v1.28.0 (2023-05-18)

* **Feature**: Removing SparkProperties from EngineConfiguration object for StartSession API call

# v1.27.0 (2023-05-15)

* **Feature**: You can now define custom spark properties at start of the session for use cases like cluster encryption, table formats, and general Spark tuning.

# v1.26.1 (2023-05-04)

* No change notes available for this release.

# v1.26.0 (2023-04-28)

* **Feature**: You can now use capacity reservations on Amazon Athena to run SQL queries on fully-managed compute capacity.

# v1.25.4 (2023-04-27)

* No change notes available for this release.

# v1.25.3 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.2 (2023-04-10)

* No change notes available for this release.

# v1.25.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2023-03-30)

* **Feature**: Make DefaultExecutorDpuSize and CoordinatorDpuSize  fields optional  in StartSession

# v1.24.0 (2023-03-27)

* **Feature**: Enforces a minimal level of encryption for the workgroup for query and calculation results that are written to Amazon S3. When enabled, workgroup users can set encryption only to the minimum level set by the administrator or higher when they submit queries.

# v1.23.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2023-03-08)

* **Feature**: A new field SubstatementType is added to GetQueryExecution API, so customers have an error free way to detect the query type and interpret the result.

# v1.22.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.22.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.22.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.21.0 (2022-12-19)

* **Feature**: Add missed InvalidRequestException in GetCalculationExecutionCode,StopCalculationExecution APIs. Correct required parameters (Payload and Type) in UpdateNotebook API. Change Notebook size from 15 Mb to 10 Mb.

# v1.20.3 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2022-12-08)

* No change notes available for this release.

# v1.20.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-11-30)

* **Feature**: This release includes support for using Apache Spark in Amazon Athena.

# v1.19.1 (2022-11-08)

* No change notes available for this release.

# v1.19.0 (2022-11-07)

* **Feature**: Adds support for using Query Result Reuse

# v1.18.12 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.11 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.10 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.9 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.8 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.7 (2022-09-01)

* No change notes available for this release.

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

* **Feature**: This feature allows customers to retrieve runtime statistics for completed queries

# v1.17.0 (2022-07-14)

* **Feature**: This release updates data types that contain either QueryExecutionId, NamedQueryId or ExpectedBucketOwner. Ids must be between 1 and 128 characters and contain only non-whitespace characters. ExpectedBucketOwner must be 12-digit string.

# v1.16.1 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-06-30)

* **Feature**: This feature introduces the API support for Athena's parameterized query and BatchGetPreparedStatement API.

# v1.15.4 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.3 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-04-15)

* **Feature**: This release adds subfields, ErrorMessage, Retryable, to the AthenaError response object in the GetQueryExecution API when a query fails.

# v1.14.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service client model to latest release.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.9.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.

# v1.8.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-08-12)

* **Feature**: API client updated

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

