# v1.21.3 (2024-09-27)

* No change notes available for this release.

# v1.21.2 (2024-09-25)

* No change notes available for this release.

# v1.21.1 (2024-09-23)

* No change notes available for this release.

# v1.21.0 (2024-09-20)

* **Feature**: Add tracing and metrics support to service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2024-09-17)

* **Bug Fix**: **BREAKFIX**: Only generate AccountIDEndpointMode config for services that use it. This is a compiler break, but removes no actual functionality, as no services currently use the account ID in endpoint resolution.

# v1.20.1 (2024-09-04)

* No change notes available for this release.

# v1.20.0 (2024-09-03)

* **Feature**: Add support to let data publisher specify a subset of the data asset that a subscriber will have access to based on the asset filters provided, when accepting a subscription request.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2024-08-30)

* **Feature**: Amazon DataZone now adds new governance capabilities of Domain Units for organization within your Data Domains, and Authorization Policies for tighter controls.

# v1.18.0 (2024-08-28)

* **Feature**: Update regex to include dot character to be consistent with IAM role creation in the authorized principal field for create and update subscription target.

# v1.17.2 (2024-08-22)

* No change notes available for this release.

# v1.17.1 (2024-08-15)

* **Dependency Update**: Bump minimum Go version to 1.21.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2024-08-05)

* **Feature**: This releases Data Product feature. Data Products allow grouping data assets into cohesive, self-contained units for ease of publishing for data producers, and ease of finding and accessing for data consumers.

# v1.16.0 (2024-07-25)

* **Feature**: Introduces GetEnvironmentCredentials operation to SDK

# v1.15.0 (2024-07-23)

* **Feature**: This release removes the deprecated dataProductItem field from Search API output.

# v1.14.0 (2024-07-22)

* **Feature**: This release adds 1/ support of register S3 locations of assets in AWS Lake Formation hybrid access mode for DefaultDataLake blueprint. 2/ support of CRUD operations for Asset Filters.

# v1.13.2 (2024-07-10.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2024-07-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2024-07-09)

* **Feature**: This release deprecates dataProductItem field from SearchInventoryResultItem, along with some unused DataProduct shapes

# v1.12.1 (2024-06-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2024-06-27)

* **Feature**: This release supports the data lineage feature of business data catalog in Amazon DataZone.

# v1.11.0 (2024-06-26)

* **Feature**: Support list-of-string endpoint parameter.

# v1.10.1 (2024-06-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2024-06-18)

* **Feature**: Track usage of various AWS SDK features in user-agent string.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2024-06-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2024-06-14)

* **Feature**: This release introduces a new default service blueprint for custom environment creation.

# v1.8.6 (2024-06-07)

* **Bug Fix**: Add clock skew correction on all service clients
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.5 (2024-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.4 (2024-05-23)

* No change notes available for this release.

# v1.8.3 (2024-05-16)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2024-05-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2024-05-08)

* **Bug Fix**: GoDoc improvement

# v1.8.0 (2024-04-03)

* **Feature**: This release supports the feature of dataQuality to enrich asset with dataQualityResult in Amazon DataZone.

# v1.7.0 (2024-04-01)

* **Feature**: This release supports the feature of AI recommendations for descriptions to enrich the business data catalog in Amazon DataZone.

# v1.6.4 (2024-03-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2024-03-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2024-03-07)

* **Bug Fix**: Remove dependency on go-cmp.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2024-02-23)

* **Bug Fix**: Move all common, SDK-side middleware stack ops into the service client module to prevent cross-module compatibility issues in the future.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.5.2 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.5.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2024-01-30)

* **Feature**: Add new skipDeletionCheck to DeleteDomain. Add new skipDeletionCheck to DeleteProject which also automatically deletes dependent objects

# v1.3.7 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.6 (2023-12-20)

* No change notes available for this release.

# v1.3.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.3.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.3.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.2.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2023-10-05)

* No change notes available for this release.

# v1.0.0 (2023-10-04)

* **Release**: New AWS service client module
* **Feature**: Initial release of Amazon DataZone

