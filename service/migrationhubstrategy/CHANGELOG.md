# v1.7.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.6.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2022-12-08)

* **Feature**: This release adds known application filtering, server selection for assessments, support for potential recommendations, and indications for configuration and assessment status. For more information, see the AWS Migration Hub documentation at https://docs.aws.amazon.com/migrationhub/index.html

# v1.5.20 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.19 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.18 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.17 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.16 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.15 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.14 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.13 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.12 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.11 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.10 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.9 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.8 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.7 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.6 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.0.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2021-11-19)

* **Release**: New AWS service client module
* **Dependency Update**: Updated to the latest SDK module versions

