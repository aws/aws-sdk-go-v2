# v1.22.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.21.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-12-12)

* **Feature**: This release allows custom strings in PartyType and Gender through 2 new attributes in the CreateProfile and UpdateProfile APIs: PartyTypeString and GenderString.

# v1.20.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-11-14)

* **Feature**: This release enhances the SearchProfiles API by providing functionality to search for profiles using multiple keys and logical operators.

# v1.19.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-09-14)

* **Feature**: Added isUnstructured in response for Customer Profiles Integration APIs
* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.8 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.7 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.6 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.5 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.4 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-06-30)

* **Feature**: This release adds the optional MinAllowedConfidenceScoreForMerging parameter to the CreateDomain, UpdateDomain, and GetAutoMergingPreview APIs in Customer Profiles. This parameter is used as a threshold to influence the profile auto-merging step of the Identity Resolution process.

# v1.17.8 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.7 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.6 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.5 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.4 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.3 (2022-03-28)

* No change notes available for this release.

# v1.17.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-12-21)

* **Feature**: API client updated

# v1.12.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-11-30)

* **Feature**: API client updated

# v1.11.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-19)

* **Feature**: API client updated
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

# v1.6.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

