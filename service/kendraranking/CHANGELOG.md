# v1.9.6 (2024-09-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.5 (2024-08-22)

* No change notes available for this release.

# v1.9.4 (2024-08-15)

* **Dependency Update**: Bump minimum Go version to 1.21.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.3 (2024-07-10.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2024-07-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2024-06-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2024-06-26)

* **Feature**: Support list-of-string endpoint parameter.

# v1.8.1 (2024-06-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2024-06-18)

* **Feature**: Track usage of various AWS SDK features in user-agent string.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.11 (2024-06-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.10 (2024-06-07)

* **Bug Fix**: Add clock skew correction on all service clients
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.9 (2024-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.8 (2024-05-23)

* No change notes available for this release.

# v1.7.7 (2024-05-16)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.6 (2024-05-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.5 (2024-05-08)

* **Bug Fix**: GoDoc improvement

# v1.7.4 (2024-03-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2024-03-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2024-03-07)

* **Bug Fix**: Remove dependency on go-cmp.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2024-02-23)

* **Bug Fix**: Move all common, SDK-side middleware stack ops into the service client module to prevent cross-module compatibility issues in the future.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.6.3 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.6.1 (2024-02-15)

* **Bug Fix**: Correct failure to determine the error type in awsJson services that could occur when errors were modeled with a non-string `code` field.

# v1.6.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.7 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.6 (2023-12-20)

* No change notes available for this release.

# v1.5.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.5.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.5.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.4.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

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

# v1.0.18 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.17 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.16 (2023-06-15)

* No change notes available for this release.

# v1.0.15 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.14 (2023-06-06)

* No change notes available for this release.

# v1.0.13 (2023-05-04)

* No change notes available for this release.

# v1.0.12 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.11 (2023-04-10)

* No change notes available for this release.

# v1.0.10 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.9 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.8 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.7 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.0.6 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2023-02-16)

* No change notes available for this release.

# v1.0.4 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.0.3 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2023-02-01)

* No change notes available for this release.

# v1.0.1 (2023-01-10)

* No change notes available for this release.

# v1.0.0 (2023-01-09)

* **Release**: New AWS service client module
* **Feature**: Introducing Amazon Kendra Intelligent Ranking, a new set of Kendra APIs that leverages Kendra semantic ranking capabilities to improve the quality of search results from other search services (i.e. OpenSearch, ElasticSearch, Solr).

