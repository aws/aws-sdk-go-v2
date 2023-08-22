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

# v1.25.4 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.2 (2023-06-15)

* No change notes available for this release.

# v1.25.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2023-06-05)

* **Feature**: Added new variable types, new DateTime data type, and new rules engine functions for interacting and working with DateTime data types.

# v1.24.0 (2023-05-31)

* **Feature**: This release enables publishing event predictions from Amazon Fraud Detector (AFD) to Amazon EventBridge. For example, after getting predictions from AFD, Amazon EventBridge rules can be configured to trigger notification through an SNS topic, send a message with SES, or trigger Lambda workflows.

# v1.23.8 (2023-05-04)

* No change notes available for this release.

# v1.23.7 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.6 (2023-04-10)

* No change notes available for this release.

# v1.23.5 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.4 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.3 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.23.1 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Feature**: This release introduces Lists feature which allows customers to reference a set of values in Fraud Detector's rules. With Lists, customers can dynamically manage these attributes in real time. Lists can be created/deleted and its contents can be modified using the Fraud Detector API.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.22.0 (2023-02-06)

* **Feature**: My AWS Service (Amazon Fraud Detector) - This release introduces Cold Start Model Training which optimizes training for small datasets and adds intelligent methods for treating unlabeled data. You can now train Online Fraud Insights or Transaction Fraud Insights models with minimal historical-data.

# v1.21.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.20.14 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.13 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.12 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.11 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.10 (2022-10-18)

* **Documentation**: Documentation Updates for Amazon Fraud Detector

# v1.20.9 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.8 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.7 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.6 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.5 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-07-21)

* **Feature**: The release introduces Account Takeover Insights (ATI) model. The ATI model detects fraud relating to account takeover. This release also adds support for new variable types: ARE_CREDENTIALS_VALID and SESSION_ID and adds new structures to Model Version APIs.

# v1.19.9 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.8 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.7 (2022-06-10)

* **Documentation**: Documentation updates for Amazon Fraud Detector (AWSHawksNest)

# v1.19.6 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-01-28)

* **Feature**: Updated to latest API model.

# v1.16.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.13.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.

# v1.12.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-10-11)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-09-10)

* **Feature**: API client updated

# v1.8.0 (2021-09-02)

* **Feature**: API client updated

# v1.7.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

