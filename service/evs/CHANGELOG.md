# v1.5.9 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.8 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.5.7 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.5.6 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.5.5 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2025-09-11)

* **Feature**: CreateEnvironment API now supports parameters (isHcxPublic & hcxNetworkAclId) for HCX migration via public internet, adding flexibility for migration scenarios. New APIs have been added for associating (AssociateEipToVlan) & disassociating (DisassociateEipFromVlan) Elastic IP (EIP) addresses.

# v1.4.4 (2025-09-10)

* No change notes available for this release.

# v1.4.3 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2025-08-26)

* **Feature**: Remove incorrect endpoint tests

# v1.3.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

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

