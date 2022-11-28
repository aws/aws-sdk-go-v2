# v1.22.0 (2022-11-28)

* **Feature**: Adds cross-account support to the GetMetricData API. Adds cross-account support to the ListMetrics API through the usage of the IncludeLinkedAccounts flag and the new OwningAccounts field.

# v1.21.11 (2022-11-22)

* No change notes available for this release.

# v1.21.10 (2022-11-16)

* No change notes available for this release.

# v1.21.9 (2022-11-10)

* No change notes available for this release.

# v1.21.8 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.7 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.6 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.5 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.4 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.3 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2022-08-30)

* No change notes available for this release.

# v1.21.1 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-08-18)

* **Feature**: Add support for managed Contributor Insights Rules

# v1.20.1 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-08-09)

* **Feature**: Various quota increases related to dimensions and custom metrics
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-07-21)

* **Feature**: Adding support for the suppression of Composite Alarm actions

# v1.18.6 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.5 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.4 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-04-14)

* **Documentation**: Updates documentation for additional statistics in CloudWatch Metric Streams.

# v1.18.0 (2022-04-13)

* **Feature**: Adds support for additional statistics in CloudWatch Metric Streams.

# v1.17.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

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

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: Updated to latest service endpoints

# v1.12.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-11-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.10.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-07-15)

* **Feature**: The ErrorCode method on generated service error types has been corrected to match the API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

