# v1.1.2 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

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
* **Feature**: AWS Marketplace Deployment is a new service that provides essential features that facilitate the deployment of software, data, and services procured through AWS Marketplace.
* **Dependency Update**: Updated to the latest SDK module versions

