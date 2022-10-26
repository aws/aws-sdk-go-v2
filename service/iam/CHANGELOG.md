# v1.18.23 (2022-10-26)

* **Documentation**: Doc only update that corrects instances of CLI not using an entity.

# v1.18.22 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.21 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.20 (2022-10-13)

* **Documentation**: Documentation updates for the AWS Identity and Access Management API Reference.

# v1.18.19 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.18 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.17 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.16 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.15 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.14 (2022-08-24)

* **Documentation**: Documentation updates for AWS Identity and Access Management (IAM).

# v1.18.13 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.12 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.11 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.10 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.9 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.8 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.7 (2022-06-08)

* **Documentation**: Documentation updates for AWS Identity and Access Management (IAM).

# v1.18.6 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

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
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.12.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-09-24)

* **Feature**: API client updated

# v1.9.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-07-15)

* **Feature**: The ErrorCode method on generated service error types has been corrected to match the API model.
* **Documentation**: Updated service model to latest revision.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2021-06-04)

* **Documentation**: Updated service client to latest API model.

# v1.5.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

