# v1.4.7 (2024-09-17)

* **Bug Fix**: **BREAKFIX**: Only generate AccountIDEndpointMode config for services that use it. This is a compiler break, but removes no actual functionality, as no services currently use the account ID in endpoint resolution.

# v1.4.6 (2024-09-04)

* No change notes available for this release.

# v1.4.5 (2024-09-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2024-08-15)

* **Dependency Update**: Bump minimum Go version to 1.21.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2024-07-10.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2024-07-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2024-06-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2024-06-26)

* **Feature**: Support list-of-string endpoint parameter.

# v1.3.1 (2024-06-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2024-06-18)

* **Feature**: Track usage of various AWS SDK features in user-agent string.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.11 (2024-06-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.10 (2024-06-07)

* **Bug Fix**: Add clock skew correction on all service clients
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.9 (2024-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.8 (2024-05-23)

* No change notes available for this release.

# v1.2.7 (2024-05-16)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.6 (2024-05-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2024-05-08)

* **Bug Fix**: GoDoc improvement

# v1.2.4 (2024-03-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2024-03-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2024-03-07)

* **Bug Fix**: Remove dependency on go-cmp.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2024-02-23)

* **Bug Fix**: Move all common, SDK-side middleware stack ops into the service client module to prevent cross-module compatibility issues in the future.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.1.3 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.1.1 (2024-02-15)

* **Bug Fix**: Correct failure to determine the error type in awsJson services that could occur when errors were modeled with a non-string `code` field.

# v1.1.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.0.3 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.0.1 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2023-11-30)

* **Release**: New AWS service client module
* **Feature**: The AWS Marketplace Agreement Service provides an API interface that helps AWS Marketplace sellers manage their agreements, including listing, filtering, and viewing details about their agreements.
* **Dependency Update**: Updated to the latest SDK module versions

