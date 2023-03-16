# v1.33.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.0 (2023-03-09)

* **Feature**: This release has two changes: add state persistence feature for embedded dashboard and console in GenerateEmbedUrlForRegisteredUser API; add properties for hidden collapsed row dimensions in PivotTableOptions.

# v1.32.2 (2023-03-03)

* No change notes available for this release.

# v1.32.1 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.32.0 (2023-02-21)

* **Feature**: S3 data sources now accept a custom IAM role.

# v1.31.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.31.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2023-02-02)

* **Feature**: QuickSight support for Radar Chart and Dashboard Publish Options

# v1.30.0 (2023-01-20)

* **Feature**: This release adds support for data bars in QuickSight table and increases pivot table field well limit.

# v1.29.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.28.3 (2022-12-16)

* No change notes available for this release.

# v1.28.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2022-11-29)

* **Feature**: This release adds new Describe APIs and updates Create and Update APIs to support the data model for Dashboards, Analyses, and Templates.

# v1.27.0 (2022-11-18)

* **Feature**: This release adds the following: 1) Asset management for centralized assets governance 2) QuickSight Q now supports public embedding 3) New Termination protection flag to mitigate accidental deletes 4) Athena data sources now accept a custom IAM role 5) QuickSight supports connectivity to Databricks

# v1.26.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2022-10-07)

* **Feature**: Amazon QuickSight now supports SecretsManager Secret ARN in place of CredentialPair for DataSource creation and update. This release also has some minor documentation updates and removes CountryCode as a required parameter in GeoSpatialColumnGroup

# v1.25.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.3 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.2 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.1 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-08-24)

* **Feature**: Added a new optional property DashboardVisual under ExperienceConfiguration parameter of GenerateEmbedUrlForAnonymousUser and GenerateEmbedUrlForRegisteredUser API operations. This supports embedding of specific visuals in QuickSight dashboards.

# v1.23.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2022-08-08)

* **Documentation**: A series of documentation updates to the QuickSight API reference.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-07-05)

* **Feature**: This release allows customers to programmatically create QuickSight accounts with Enterprise and Enterprise + Q editions. It also releases allowlisting domains for embedding QuickSight dashboards at runtime through the embedding APIs.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-05-18)

* **Feature**: API UpdatePublicSharingSettings enables IAM admins to enable/disable account level setting for public access of dashboards. When enabled, owners/co-owners for dashboards can enable public access on their dashboards. These dashboards can only be accessed through share link or embedding.

# v1.21.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-03-21)

* **Feature**: AWS QuickSight Service Features - Expand public API support for group management.

# v1.20.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Documentation**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.15.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-11-30)

* **Feature**: API client updated

# v1.14.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.1 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-09-24)

* **Feature**: API client updated

# v1.11.0 (2021-09-17)

* **Feature**: Updated API client and endpoints to latest revision.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-09-02)

* **Feature**: API client updated

# v1.9.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-05-25)

* **Feature**: API client updated

# v1.4.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

