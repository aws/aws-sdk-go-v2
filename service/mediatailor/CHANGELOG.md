# v1.26.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2023-08-01)

* No change notes available for this release.

# v1.26.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2023-07-13)

* **Feature**: Adds categories to MediaTailor channel assembly alerts
* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2023-07-07)

* **Feature**: The AWS Elemental MediaTailor SDK for Channel Assembly has added support for EXT-X-CUE-OUT and EXT-X-CUE-IN tags to specify ad breaks in HLS outputs, including support for EXT-OATCLS, EXT-X-ASSET, and EXT-X-CUE-OUT-CONT accessory tags.

# v1.23.2 (2023-06-15)

* No change notes available for this release.

# v1.23.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2023-05-05)

* **Feature**: This release adds support for AFTER_LIVE_EDGE mode configuration for avail suppression, and adding a fill-policy setting that sets the avail suppression to PARTIAL_AVAIL or FULL_AVAIL_ONLY when AFTER_LIVE_EDGE is enabled.

# v1.22.10 (2023-05-04)

* No change notes available for this release.

# v1.22.9 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.8 (2023-04-10)

* No change notes available for this release.

# v1.22.7 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.6 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.22.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.22.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2023-02-01)

* **Feature**: The AWS Elemental MediaTailor SDK for Channel Assembly has added support for program updates, and the ability to clip the end of VOD sources in programs.

# v1.21.0 (2023-01-27)

* **Feature**: This release introduces the As Run logging type, along with API and documentation updates.

# v1.20.1 (2023-01-19)

* No change notes available for this release.

# v1.20.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.19.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-10-28)

* **Feature**: This release introduces support for SCTE-35 segmentation descriptor messages which can be sent within time signal messages.

# v1.18.0 (2022-10-25)

* **Feature**: This release is a documentation update

# v1.17.16 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.15 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.14 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.13 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.12 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.11 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.10 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.9 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.8 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.7 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.6 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.5 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.4 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.3 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-04-21)

* **Feature**: This release introduces tiered channels and adds support for live sources. Customers using a STANDARD channel can now create programs using live sources.

# v1.16.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2022-01-28)

* **Feature**: Updated to latest API model.

# v1.13.0 (2022-01-14)

* **Feature**: Updated API models
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.10.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-09-24)

* **Feature**: API client updated

# v1.7.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-09-02)

* **Feature**: API client updated

# v1.6.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

