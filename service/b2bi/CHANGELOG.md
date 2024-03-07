# v1.0.0-preview.18 (2024-03-07)

* **Bug Fix**: Remove dependency on go-cmp.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.17 (2024-02-23)

* **Bug Fix**: Move all common, SDK-side middleware stack ops into the service client module to prevent cross-module compatibility issues in the future.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.16 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.0.0-preview.15 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.14 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.0.0-preview.13 (2024-02-15)

* **Bug Fix**: Correct failure to determine the error type in awsJson services that could occur when errors were modeled with a non-string `code` field.

# v1.0.0-preview.12 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.11 (2024-01-18)

* **Feature**: Increasing TestMapping inputFileContent file size limit to 5MB and adding file size limit 250KB for TestParsing input file. This release also includes exposing InternalServerException for Tag APIs.

# v1.0.0-preview.10 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.9 (2023-12-14)

* **Documentation**: Documentation updates for AWS B2B Data Interchange

# v1.0.0-preview.8 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.0.0-preview.7 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.6 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.0.0-preview.5 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.4 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.3 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.2 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0-preview.1 (2023-11-28)

* No change notes available for this release.

# v1.0.0-preview (2023-11-27.2)

* **Feature**: This is the initial SDK release for AWS B2B Data Interchange.

