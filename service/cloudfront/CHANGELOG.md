# v1.30.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.7 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.6 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2023-08-01)

* No change notes available for this release.

# v1.28.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2023-07-28.2)

* **Feature**: Add a new JavaScript runtime version for CloudFront Functions.

# v1.26.10 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.9 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.8 (2023-06-15)

* No change notes available for this release.

# v1.26.7 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.6 (2023-05-04)

* No change notes available for this release.

# v1.26.5 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.4 (2023-04-10)

* No change notes available for this release.

# v1.26.3 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2023-02-22)

* **Feature**: CloudFront now supports block lists in origin request policies so that you can forward all headers, cookies, or query string from viewer requests to the origin *except* for those specified in the block list.
* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.25.1 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2023-02-08)

* **Feature**: CloudFront Origin Access Control extends support to AWS Elemental MediaStore origins.

# v1.24.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.23.0 (2022-12-30)

* **Feature**: Extend response headers policy to support removing headers from viewer responses

# v1.22.2 (2022-12-16)

* **Documentation**: Updated documentation for CloudFront

# v1.22.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-12-07)

* **Feature**: Introducing UpdateDistributionWithStagingConfig that can be used to promote the staging configuration to the production.

# v1.21.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-11-18)

* **Feature**: CloudFront API support for staging distributions and associated traffic management policies.

# v1.20.7 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.6 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.5 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.4 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.3 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2022-08-31)

* **Documentation**: Update API documentation for CloudFront origin access control (OAC)
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-08-24)

* **Feature**: Adds support for CloudFront origin access control (OAC), making it possible to restrict public access to S3 bucket origins in all AWS Regions, those with SSE-KMS, and more.

# v1.19.0 (2022-08-15)

* **Feature**: Adds Http 3 support to distributions

# v1.18.8 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.7 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.6 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.5 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.4 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-05-16)

* **Feature**: Introduced a new error (TooLongCSPInResponseHeadersPolicy) that is returned when the value of the Content-Security-Policy header in a response headers policy exceeds the maximum allowed length.

# v1.17.0 (2022-04-26)

* **Feature**: CloudFront now supports the Server-Timing header in HTTP responses sent from CloudFront. You can use this header to view metrics that help you gain insights about the behavior and performance of CloudFront. To use this header, enable it in a response headers policy.

# v1.16.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

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

# v1.14.1 (2022-01-28)

* **Bug Fix**: Updates SDK API client deserialization to pre-allocate byte slice and string response payloads, [#1565](https://github.com/aws/aws-sdk-go-v2/pull/1565). Thanks to [Tyson Mote](https://github.com/tysonmote) for submitting this PR.

# v1.14.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.11.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-11-12)

* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.10.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2021-06-04)

* **Documentation**: Updated service client to latest API model.

# v1.5.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

