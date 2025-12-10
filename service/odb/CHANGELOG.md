# v1.7.0 (2025-12-10)

* **Feature**: The following APIs now return CloudExadataInfrastructureArn and OdbNetworkArn fields for improved resource identification and AWS service integration - GetCloudVmCluster, ListCloudVmClusters, GetCloudAutonomousVmCluster, and ListCloudAutonomousVmClusters.

# v1.6.3 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.6.1 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.6.0 (2025-11-21)

* **Feature**: Adds AssociateIamRoleToResource and DisassociateIamRoleFromResource APIs for managing IAM roles. Enhances CreateOdbNetwork and UpdateOdbNetwork APIs with KMS, STS, and cross-region S3 parameters. Adds OCI identity domain support to InitializeService API.

# v1.5.7 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.6 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.5.5 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.5.4 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.5.3 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2025-10-23)

* **Documentation**: Doc-only update that removes duplicate values from descriptions of ODB peering APIs.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2025-10-10)

* **Feature**: This release adds APIs that allow you to specify CIDR ranges in your ODB peering connection.

# v1.4.6 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.5 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2025-09-10)

* No change notes available for this release.

# v1.4.3 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2025-08-25)

* **Feature**: Remove incorrect endpoint tests

# v1.3.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.3.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2025-07-01)

* **Release**: New AWS service client module
* **Feature**: This release adds API operations for Oracle Database@AWS. You can use the APIs to create Exadata infrastructure, ODB networks, and Exadata and Autonomous VM clusters inside AWS data centers. The infrastructure is managed by OCI. You can integrate these resources with AWS services.

