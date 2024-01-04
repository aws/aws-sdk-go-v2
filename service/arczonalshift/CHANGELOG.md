# v1.7.6 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.5 (2023-12-22)

* No change notes available for this release.

# v1.7.4 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.7.3 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.7.1 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2023-11-30.2)

* **Feature**: This release adds a new capability, zonal autoshift. You can configure zonal autoshift so that AWS shifts traffic for a resource away from an Availability Zone, on your behalf, when AWS determines that there is an issue that could potentially affect customers in the Availability Zone.

# v1.6.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.6 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.5 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.5.4 (2023-11-27)

* No change notes available for this release.

# v1.5.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-09-18)

* **Announcement**: [BREAKFIX] Change in MaxResults datatype from value to pointer type in cognito-sync service.
* **Feature**: Adds several endpoint ruleset changes across all models: smaller rulesets, removed non-unique regional endpoints, fixes FIPS and DualStack endpoints, and make region not required in SDK::Endpoint. Additional breakfix to cognito-sync field.

# v1.2.6 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2023-08-09)

* No change notes available for this release.

# v1.2.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-08-01)

* No change notes available for this release.

# v1.2.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.23 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.22 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.21 (2023-06-15)

* No change notes available for this release.

# v1.1.20 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.19 (2023-05-04)

* No change notes available for this release.

# v1.1.18 (2023-04-28)

* No change notes available for this release.

# v1.1.17 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.16 (2023-04-19)

* No change notes available for this release.

# v1.1.15 (2023-04-17)

* No change notes available for this release.

# v1.1.14 (2023-04-13)

* No change notes available for this release.

# v1.1.13 (2023-04-11)

* No change notes available for this release.

# v1.1.12 (2023-04-10)

* No change notes available for this release.

# v1.1.11 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.10 (2023-04-05)

* No change notes available for this release.

# v1.1.9 (2023-04-03)

* No change notes available for this release.

# v1.1.8 (2023-03-29)

* No change notes available for this release.

# v1.1.7 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.6 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.1.4 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.1.2 (2023-02-08)

* No change notes available for this release.

# v1.1.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.0.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2022-11-29)

* **Release**: New AWS service client module
* **Feature**: Amazon Route 53 Application Recovery Controller Zonal Shift is a new service that makes it easy to shift traffic away from an Availability Zone in a Region. See the developer guide for more information: https://docs.aws.amazon.com/r53recovery/latest/dg/what-is-route53-recovery.html

