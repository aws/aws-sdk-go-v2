# v1.32.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.3 (2023-08-17)

* **Announcement**: BREAKFIX: corrected function spelling in environment config from GetS3DisableMultRegionAccessPoints to GetS3DisableMultiRegionAccessPoints
* **Bug Fix**: Adds DisableMRAP option to config loader, and DisableMRAP client resolver to achieve parity with other S3 options in the config loader. Additionally, added breakfix to correct spelling.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.1 (2023-08-01)

* No change notes available for this release.

# v1.32.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.10 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.9 (2023-07-18)

* No change notes available for this release.

# v1.31.8 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.7 (2023-06-15)

* No change notes available for this release.

# v1.31.6 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.5 (2023-05-04)

* No change notes available for this release.

# v1.31.4 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.3 (2023-04-10)

* No change notes available for this release.

# v1.31.2 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.1 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2023-03-15)

* **Feature**: Added support for S3 Object Lambda aliases.

# v1.30.0 (2023-03-13)

* **Feature**: Added support for cross-account Multi-Region Access Points. Added support for S3 Replication for S3 on Outposts.

# v1.29.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.29.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.1 (2023-01-23)

* No change notes available for this release.

# v1.29.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.28.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2022-11-30)

* **Feature**: Amazon S3 now supports cross-account access points. S3 bucket owners can now allow trusted AWS accounts to create access points associated with their bucket.

# v1.27.0 (2022-11-29)

* **Feature**: Added two new APIs to support Amazon S3 Multi-Region Access Point failover controls: GetMultiRegionAccessPointRoutes and SubmitMultiRegionAccessPointRoutes. The failover control APIs are supported in the following Regions: us-east-1, us-west-2, eu-west-1, ap-southeast-2, and ap-northeast-1.

# v1.26.1 (2022-11-22)

* No change notes available for this release.

# v1.26.0 (2022-11-17)

* **Feature**: Added 34 new S3 Storage Lens metrics to support additional customer use cases.

# v1.25.2 (2022-11-16)

* No change notes available for this release.

# v1.25.1 (2022-11-10)

* No change notes available for this release.

# v1.25.0 (2022-11-02)

* **Feature**: S3 on Outposts launches support for Lifecycle configuration for Outposts buckets. With S3 Lifecycle configuration, you can mange objects so they are stored cost effectively. You can manage objects using size-based rules and specify how many noncurrent versions bucket will retain.

# v1.24.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-10-04)

* **Feature**: S3 Object Lambda adds support to allow customers to intercept HeadObject and ListObjects requests and introduce their own compute. These requests were previously proxied to S3.

# v1.23.0 (2022-09-21)

* **Feature**: S3 on Outposts launches support for object versioning for Outposts buckets. With S3 Versioning, you can preserve, retrieve, and restore every version of every object stored in your buckets. You can recover from both unintended user actions and application failures.

# v1.22.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.17 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.16 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.15 (2022-08-30)

* No change notes available for this release.

# v1.21.14 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.13 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.12 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.11 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.10 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.9 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.8 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.7 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.6 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.5 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.4 (2022-04-05)

* **Documentation**: Documentation-only update for doc bug fixes for the S3 Control API docs.

# v1.21.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-02-24.2)

* **Feature**: API client updated

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
* **Feature**: Updated to latest service endpoints

# v1.15.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-11-30)

* **Feature**: API client updated

# v1.14.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-11-12)

* **Feature**: Updated service to latest API model.

# v1.13.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-09-02)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

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

# v1.8.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-06-04)

* **Feature**: The handling of AccessPoint and Outpost ARNs have been updated.
* **Feature**: Updated service client to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Feature**: Updated to latest service API model.
* **Dependency Update**: Updated to the latest SDK module versions

