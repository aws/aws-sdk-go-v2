# v1.3.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.3.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Documentation**: Update for general availability of Amazon Elastic VMware Service (EVS).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Feature**: TagResource API now throws ServiceQuotaExceededException when the number of tags on the Amazon EVS resource exceeds the maximum allowed. TooManyTagsException is deprecated.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2025-06-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2025-06-04)

* **Release**: New AWS service client module
* **Feature**: Amazon Elastic VMware Service (Amazon EVS) allows you to run VMware Cloud Foundation (VCF) directly within your Amazon VPC including simplified self-managed migration experience with guided workflow in AWS console or via AWS CLI, get full access to their VCF deployment and VCF license portability.

