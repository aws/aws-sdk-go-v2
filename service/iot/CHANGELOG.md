# v1.39.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.1 (2023-08-01)

* No change notes available for this release.

# v1.39.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.4 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.2 (2023-06-15)

* No change notes available for this release.

# v1.38.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.0 (2023-06-06)

* **Feature**: Adding IoT Device Management Software Package Catalog APIs to register, store, and report system software packages, along with their versions and metadata in a centralized location.

# v1.37.2 (2023-05-23)

* No change notes available for this release.

# v1.37.1 (2023-05-04)

* No change notes available for this release.

# v1.37.0 (2023-04-28)

* **Feature**: This release allows AWS IoT Core users to specify a TLS security policy when creating and updating AWS IoT Domain Configurations.

# v1.36.1 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.0 (2023-04-20)

* **Feature**: Support additional OTA states in GetOTAUpdate API

# v1.35.4 (2023-04-10)

* No change notes available for this release.

# v1.35.3 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.0 (2023-03-02)

* **Feature**: A recurring maintenance window is an optional configuration used for rolling out the job document to all devices in the target group observing a predetermined start time, duration, and frequency that the maintenance window occurs.

# v1.34.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.34.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.34.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.0 (2023-01-31)

* **Feature**: Added support for IoT Rules Engine Cloudwatch Logs action batch mode.

# v1.33.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.32.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.0 (2022-11-28)

* **Feature**: Job scheduling enables the scheduled rollout of a Job with start and end times and a customizable end behavior when end time is reached. This is available for continuous and snapshot jobs. Added support for MQTT5 properties to AWS IoT TopicRule Republish Action.

# v1.31.0 (2022-11-11)

* **Feature**: This release add new api listRelatedResourcesForAuditFinding and new member type IssuerCertificates for Iot device device defender Audit.

# v1.30.0 (2022-10-31)

* **Feature**: This release adds the Amazon Location action to IoT Rules Engine.

# v1.29.4 (2022-10-25)

* No change notes available for this release.

# v1.29.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.6 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.5 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.4 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.3 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.2 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2022-08-04)

* **Feature**: The release is to support attach a provisioning template to CACert for JITP function,  Customer now doesn't have to hardcode a roleArn and templateBody during register a CACert to enable JITP.

# v1.27.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2022-07-20)

* **Feature**: GA release the ability to enable/disable IoT Fleet Indexing for Device Defender and Named Shadow information, and search them through IoT Fleet Indexing APIs. This includes Named Shadow Selection as a part of the UpdateIndexingConfiguration API.

# v1.26.0 (2022-07-07)

* **Feature**: This release adds support to register a CA certificate without having to provide a verification certificate. This also allows multiple AWS accounts to register the same CA in the same region.

# v1.25.5 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.4 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.3 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2022-05-12)

* **Documentation**: Documentation update for China region ListMetricValues for IoT

# v1.25.0 (2022-05-05)

* **Feature**: AWS IoT Jobs now allows you to create up to 100,000 active continuous and snapshot jobs by using concurrency control.

# v1.24.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-04-04)

* **Feature**: AWS IoT - AWS IoT Device Defender adds support to list metric datapoints collected for IoT devices through the ListMetricValues API

# v1.23.3 (2022-03-30)

* **Documentation**: Doc only update for IoT that fixes customer-reported issues.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-01-07)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: API client updated

# v1.18.0 (2021-12-02)

* **Feature**: API client updated
* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2021-11-30)

* **Feature**: API client updated

# v1.16.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.

# v1.15.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-09-24)

* **Feature**: API client updated

# v1.12.0 (2021-09-17)

* **Feature**: Updated API client and endpoints to latest revision.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-09-02)

* **Feature**: API client updated

# v1.10.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-05-25)

* **Feature**: API client updated

# v1.5.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Feature**: Updated to latest service API model.
* **Dependency Update**: Updated to the latest SDK module versions

