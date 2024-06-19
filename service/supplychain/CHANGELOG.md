# v1.4.1 (2024-06-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2024-06-18)

* **Feature**: Track usage of various AWS SDK features in user-agent string.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.7 (2024-06-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.6 (2024-06-07)

* **Bug Fix**: Add clock skew correction on all service clients
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.5 (2024-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.4 (2024-05-23)

* No change notes available for this release.

# v1.3.3 (2024-05-16)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.2 (2024-05-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2024-05-08)

* **Bug Fix**: GoDoc improvement

# v1.3.0 (2024-04-10)

* **Feature**: This release includes API SendDataIntegrationEvent for AWS Supply Chain

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

# v1.1.2 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.1.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2024-01-12)

* **Release**: New AWS service client module
* **Feature**: This release includes APIs CreateBillOfMaterialsImportJob and GetBillOfMaterialsImportJob.

