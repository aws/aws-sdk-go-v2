# v1.7.7 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.6 (2023-12-28)

* No change notes available for this release.

# v1.7.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.7.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.7.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.6.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-10-30)

* **Feature**: This feature allows customers to create a combination of manifest filtering, startover and time delay configuration that applies to all egress requests by default.

# v1.3.1 (2023-10-26)

* No change notes available for this release.

# v1.3.0 (2023-10-16)

* **Feature**: This release allows customers to manage MediaPackage v2 resource using CloudFormation.

# v1.2.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-09-18)

* **Announcement**: [BREAKFIX] Change in MaxResults datatype from value to pointer type in cognito-sync service.
* **Feature**: Adds several endpoint ruleset changes across all models: smaller rulesets, removed non-unique regional endpoints, fixes FIPS and DualStack endpoints, and make region not required in SDK::Endpoint. Additional breakfix to cognito-sync field.

# v1.1.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2023-08-01)

* No change notes available for this release.

# v1.1.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2023-07-18)

* No change notes available for this release.

# v1.0.4 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2023-06-15)

* No change notes available for this release.

# v1.0.2 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2023-05-23)

* No change notes available for this release.

# v1.0.0 (2023-05-19)

* **Release**: New AWS service client module
* **Feature**: Adds support for the MediaPackage Live v2 API

