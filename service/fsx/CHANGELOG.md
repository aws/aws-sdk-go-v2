# v1.32.3 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.2 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.1 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.0 (2023-08-09)

* **Feature**: For FSx for Lustre, add new data repository task type, RELEASE_DATA_FROM_FILESYSTEM, to release files that have been archived to S3. For FSx for Windows, enable support for configuring and updating SSD IOPS, and for updating storage type. For FSx for OpenZFS, add new deployment type, MULTI_AZ_1.

# v1.31.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.1 (2023-08-01)

* No change notes available for this release.

# v1.31.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-07-13)

* **Feature**: Amazon FSx for NetApp ONTAP now supports SnapLock, an ONTAP feature that enables you to protect your files in a volume by transitioning them to a write once, read many (WORM) state.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.3 (2023-06-23)

* **Documentation**: Update to Amazon FSx documentation.

# v1.29.2 (2023-06-15)

* No change notes available for this release.

# v1.29.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-06-12)

* **Feature**: Amazon FSx for NetApp ONTAP now supports joining a storage virtual machine (SVM) to Active Directory after the SVM has been created.

# v1.28.13 (2023-05-25)

* No change notes available for this release.

# v1.28.12 (2023-05-04)

* No change notes available for this release.

# v1.28.11 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.10 (2023-04-14)

* No change notes available for this release.

# v1.28.9 (2023-04-10)

* No change notes available for this release.

# v1.28.8 (2023-04-07)

* **Documentation**: Amazon FSx for Lustre now supports creating data repository associations on Persistent_1 and Scratch_2 file systems.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.7 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.6 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.28.4 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.3 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.28.2 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2023-01-23)

* No change notes available for this release.

# v1.28.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.27.0 (2022-12-23)

* **Feature**: Fix a bug where a recent release might break certain existing SDKs.

# v1.26.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2022-11-29)

* **Feature**: This release adds support for 4GB/s / 160K PIOPS FSx for ONTAP file systems and 10GB/s / 350K PIOPS FSx for OpenZFS file systems (Single_AZ_2). For FSx for ONTAP, this also adds support for DP volumes, snapshot policy, copy tags to backups, and Multi-AZ route table updates.

# v1.25.4 (2022-10-31)

* No change notes available for this release.

# v1.25.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2022-10-19)

* No change notes available for this release.

# v1.25.0 (2022-09-29)

* **Feature**: This release adds support for Amazon File Cache.

# v1.24.14 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.13 (2022-09-19)

* No change notes available for this release.

# v1.24.12 (2022-09-14)

* **Documentation**: Documentation update for Amazon FSx.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.11 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.10 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.9 (2022-08-29)

* **Documentation**: Documentation updates for Amazon FSx for NetApp ONTAP.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.8 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.7 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.6 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.5 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.4 (2022-07-29)

* **Documentation**: Documentation updates for Amazon FSx

# v1.24.3 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.2 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-05-25)

* **Feature**: This release adds root squash support to FSx for Lustre to restrict root level access from clients by mapping root users to a less-privileged user/group with limited permissions.

# v1.23.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-04-13)

* **Feature**: This release adds support for deploying FSx for ONTAP file systems in a single Availability Zone.

# v1.22.0 (2022-04-05)

* **Feature**: Provide customers more visibility into file system status by adding new "Misconfigured Unavailable" status for Amazon FSx for Windows File Server.

# v1.21.0 (2022-03-30)

* **Feature**: This release adds support for modifying throughput capacity for FSx for ONTAP file systems.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service client model to latest release.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-01-28)

* **Feature**: Updated to latest API model.

# v1.17.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.14.0 (2021-12-02)

* **Feature**: API client updated
* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.

# v1.12.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-10-11)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-09-02)

* **Feature**: API client updated

# v1.8.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-06-11)

* **Feature**: Updated to latest API model.

# v1.5.0 (2021-06-04)

* **Feature**: Updated service client to latest API model.

# v1.4.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

