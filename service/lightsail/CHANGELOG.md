# v1.35.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Feature**: This release adds support to upgrade the major version of a database.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.0 (2024-01-25)

* **Feature**: This release adds support for IPv6-only instance plans.

# v1.33.0 (2024-01-04)

* **Feature**: This release adds support to set up an HTTPS endpoint on an instance.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.32.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.32.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.31.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-10-24)

* **Feature**: **BREAKFIX**: Correct nullability and default value representation of various input fields across a large number of services. Calling code that references one or more of the affected fields will need to update usage accordingly. See [2162](https://github.com/aws/aws-sdk-go-v2/issues/2162).

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

# v1.27.3 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.2 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.1 (2023-06-15)

* No change notes available for this release.

# v1.27.0 (2023-06-13)

* **Feature**: This release adds pagination for the Get Certificates API operation.
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

# v1.26.0 (2023-02-28)

* **Feature**: This release adds Lightsail for Research feature support, such as GUI session access, cost estimates, stop instance on idle, and disk auto mount.

# v1.25.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.25.4 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.3 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.25.2 (2023-02-09)

* **Documentation**: Documentation updates for Lightsail

# v1.25.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).
* **Documentation**: Documentation updates for Amazon Lightsail.

# v1.24.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-11-08)

* **Feature**: This release adds support for Amazon Lightsail to automate the delegation of domains registered through Amazon Route 53 to Lightsail DNS management and to automate record creation for DNS validation of Lightsail SSL/TLS certificates.

# v1.23.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-09-23)

* **Feature**: This release adds Instance Metadata Service (IMDS) support for Lightsail instances.

# v1.22.12 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.11 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.10 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.9 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.8 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.7 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.6 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.5 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.4 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.3 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-05-26)

* **Feature**: Amazon Lightsail now supports the ability to configure a Lightsail Container Service to pull images from Amazon ECR private repositories in your account.

# v1.21.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-05-12)

* **Feature**: This release adds support to include inactive database bundles in the response of the GetRelationalDatabaseBundles request.

# v1.20.1 (2022-05-04)

* **Documentation**: Documentation updates for Lightsail

# v1.20.0 (2022-04-26)

* **Feature**: This release adds support for Lightsail load balancer HTTP to HTTPS redirect and TLS policy configuration.

# v1.19.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-04-15)

* **Feature**: This release adds support to describe the synchronization status of the account-level block public access feature for your Amazon Lightsail buckets.

# v1.18.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-02-24.2)

* **Feature**: API client updated

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
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-08-12)

* **Feature**: API client updated

# v1.8.1 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-06-04)

* **Documentation**: Updated service client to latest API model.

# v1.6.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

