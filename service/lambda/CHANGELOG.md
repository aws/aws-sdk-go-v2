# v1.48.0 (2023-11-16)

* **Feature**: Adds support for logging configuration in Lambda Functions. Customers will have more control how their function logs are captured and to which cloud watch log group they are delivered also.

# v1.47.0 (2023-11-15)

* **Feature**: Add Java 21 (java21) support to AWS Lambda
* **Dependency Update**: Updated to the latest SDK module versions

# v1.46.0 (2023-11-14)

* **Feature**: Add Python 3.12 (python3.12) support to AWS Lambda

# v1.45.0 (2023-11-09.2)

* **Feature**: Add Custom runtime on Amazon Linux 2023 (provided.al2023) support to AWS Lambda.

# v1.44.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.44.0 (2023-11-08)

* **Feature**: Add Node 20 (nodejs20.x) support to AWS Lambda.

# v1.43.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.41.0 (2023-10-24)

* **Feature**: **BREAKFIX**: Correct nullability and default value representation of various input fields across a large number of services. Calling code that references one or more of the affected fields will need to update usage accordingly. See [2162](https://github.com/aws/aws-sdk-go-v2/issues/2162).

# v1.40.0 (2023-10-12)

* **Feature**: Adds support for Lambda functions to access Dual-Stack subnets over IPv6, via an opt-in flag in CreateFunction and UpdateFunctionConfiguration APIs
* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.6 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.1 (2023-08-01)

* No change notes available for this release.

# v1.39.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.0 (2023-07-25)

* **Feature**: Add Python 3.11 (python3.11) support to AWS Lambda

# v1.37.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.0 (2023-06-28)

* **Feature**: Surface ResourceConflictException in DeleteEventSourceMapping

# v1.36.0 (2023-06-20)

* **Feature**: This release adds RecursiveInvocationException to the Invoke API and InvokeWithResponseStream API.

# v1.35.2 (2023-06-15)

* No change notes available for this release.

# v1.35.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.0 (2023-06-05)

* **Feature**: Add Ruby 3.2 (ruby3.2) Runtime support to AWS Lambda.

# v1.34.1 (2023-05-04)

* No change notes available for this release.

# v1.34.0 (2023-04-27)

* **Feature**: Add Java 17 (java17) support to AWS Lambda

# v1.33.1 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.0 (2023-04-17)

* **Feature**: Add Python 3.10 (python3.10) support to AWS Lambda

# v1.32.0 (2023-04-14)

* **Feature**: This release adds SnapStart related exceptions to InvokeWithResponseStream API. IAM access related documentation is also added for this API.

# v1.31.1 (2023-04-10)

* No change notes available for this release.

# v1.31.0 (2023-04-07)

* **Feature**: This release adds a new Lambda InvokeWithResponseStream API to support streaming Lambda function responses. The release also adds a new InvokeMode parameter to Function Url APIs to control whether the response will be streamed or buffered.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-02-27)

* **Feature**: This release adds the ability to create ESMs with Document DB change streams as event source. For more information see  https://docs.aws.amazon.com/lambda/latest/dg/with-documentdb.html.

# v1.29.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.29.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.29.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-01-23)

* **Feature**: Release Lambda RuntimeManagementConfig, enabling customers to better manage runtime updates to their Lambda functions. This release adds two new APIs, GetRuntimeManagementConfig and PutRuntimeManagementConfig, as well as support on existing Create/Get/Update function APIs.

# v1.28.0 (2023-01-12)

* **Feature**: Add support for MaximumConcurrency parameter for SQS event source. Customers can now limit the maximum concurrent invocations for their SQS Event Source Mapping.

# v1.27.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.26.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2022-11-29)

* **Feature**: Adds support for Lambda SnapStart, which helps improve the startup performance of functions. Customers can now manage SnapStart based functions via CreateFunction and UpdateFunctionConfiguration APIs

# v1.25.1 (2022-11-22)

* No change notes available for this release.

# v1.25.0 (2022-11-17)

* **Feature**: Add Node 18 (nodejs18.x) support to AWS Lambda.

# v1.24.11 (2022-11-16)

* No change notes available for this release.

# v1.24.10 (2022-11-10)

* No change notes available for this release.

# v1.24.9 (2022-11-02)

* No change notes available for this release.

# v1.24.8 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.7 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.6 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.5 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.4 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.3 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.2 (2022-08-30)

* No change notes available for this release.

# v1.24.1 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-08-17)

* **Feature**: Added support for customization of Consumer Group ID for MSK and Kafka Event Source Mappings.

# v1.23.8 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.7 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.6 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.5 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.4 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.3 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-05-12)

* **Feature**: Lambda releases NodeJs 16 managed runtime to be available in all commercial regions.

# v1.22.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-04-06)

* **Feature**: This release adds new APIs for creating and managing Lambda Function URLs and adds a new FunctionUrlAuthType parameter to the AddPermission API. Customers can use Function URLs to create built-in HTTPS endpoints on their functions.

# v1.21.1 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-03-24)

* **Feature**: Adds support for increased ephemeral storage (/tmp) up to 10GB for Lambda functions. Customers can now provision up to 10 GB of ephemeral storage per function instance, a 20x increase over the previous limit of 512 MB.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-03-11)

* **Feature**: Adds PrincipalOrgID support to AddPermission API. Customers can use it to manage permissions to lambda functions at AWS Organizations level.

# v1.19.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.1 (2022-01-28)

* **Bug Fix**: Updates SDK API client deserialization to pre-allocate byte slice and string response payloads, [#1565](https://github.com/aws/aws-sdk-go-v2/pull/1565). Thanks to [Tyson Mote](https://github.com/tysonmote) for submitting this PR.

# v1.17.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-01-07)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: Updated to latest service endpoints

# v1.14.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-11-30)

* **Feature**: API client updated

# v1.13.0 (2021-11-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-11-12)

* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.11.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-09-30)

* **Feature**: API client updated

# v1.8.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-08-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

