# v1.14.6 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.14.5 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.14.4 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.14.3 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.2 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.1 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2025-10-01)

* **Feature**: Added the UpdateCluster API action to modify cluster configurations, and Slurm custom settings for queues.

# v1.13.2 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2025-09-17)

* **Feature**: Add support for Amazon EC2 Capacity Blocks for ML

# v1.12.4 (2025-09-10)

* No change notes available for this release.

# v1.12.3 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.2 (2025-09-05)

* **Documentation**: Documentation-only update to add AccountingStorageEnforce to SlurmCustomSetting.

# v1.12.1 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2025-08-28)

* **Feature**: Remove incorrect endpoint tests

# v1.11.3 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.11.0 (2025-08-14)

* **Feature**: Updated the regex pattern and description of iamInstanceProfileArn in the CreateComputeNodeGroup and UpdateComputeNodeGroup API actions. Name and path requirements apply to the ARN of the IAM role associated with the instance profile and not the ARN of the instance profile.

# v1.10.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2025-08-01)

* **Feature**: Add support for IPv6 Networking for Clusters.

# v1.7.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2025-06-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2025-06-02)

* **Feature**: Introduces SUSPENDING and SUSPENDED states for clusters, compute node groups, and queues.

# v1.5.0 (2025-05-15)

* **Feature**: This release adds support for Slurm accounting. For more information, see the Slurm accounting topic in the AWS PCS User Guide. Slurm accounting is supported for Slurm 24.11 and later. This release also adds 24.11 as a valid value for the version parameter of the Scheduler data type.

# v1.4.2 (2025-04-24)

* **Documentation**: Documentation-only update: added valid values for the version property of the Scheduler and SchedulerRequest data types.

# v1.4.1 (2025-04-03)

* No change notes available for this release.

# v1.4.0 (2025-03-24)

* **Feature**: ClusterName/ClusterIdentifier, ComputeNodeGroupName/ComputeNodeGroupIdentifier, and QueueName/QueueIdentifier can now have 10 characters, and a minimum of 3 characters. The TagResource API action can now return ServiceQuotaExceededException.

# v1.3.1 (2025-03-04.2)

* **Bug Fix**: Add assurance test for operation order.

# v1.3.0 (2025-02-27)

* **Feature**: Track credential providers via User-Agent Feature ids
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.17 (2025-02-18)

* **Bug Fix**: Bump go version to 1.22
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.16 (2025-02-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.15 (2025-01-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.14 (2025-01-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.13 (2025-01-24)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.22.2.

# v1.2.12 (2025-01-17)

* **Bug Fix**: Fix bug where credentials weren't refreshed during retry loop.

# v1.2.11 (2025-01-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.10 (2025-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.9 (2024-12-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.8 (2024-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.7 (2024-11-18)

* **Dependency Update**: Update to smithy-go v1.22.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.6 (2024-11-07)

* **Bug Fix**: Adds case-insensitive handling of error message fields in service responses

# v1.2.5 (2024-11-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2024-10-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2024-10-24)

* **Documentation**: Documentation update: added the default value of the Slurm configuration parameter scaleDownIdleTimeInSeconds to its description.

# v1.2.2 (2024-10-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2024-10-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2024-10-04)

* **Feature**: Add support for HTTP client metrics.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.5 (2024-10-03)

* No change notes available for this release.

# v1.1.4 (2024-09-27)

* No change notes available for this release.

# v1.1.3 (2024-09-26)

* **Documentation**: AWS PCS API documentation - Edited the description of the iamInstanceProfileArn parameter of the CreateComputeNodeGroup and UpdateComputeNodeGroup actions; edited the description of the SlurmCustomSetting data type to list the supported parameters for clusters and compute node groups.

# v1.1.2 (2024-09-25)

* No change notes available for this release.

# v1.1.1 (2024-09-23)

* No change notes available for this release.

# v1.1.0 (2024-09-20)

* **Feature**: Add tracing and metrics support to service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2024-09-17)

* **Bug Fix**: **BREAKFIX**: Only generate AccountIDEndpointMode config for services that use it. This is a compiler break, but removes no actual functionality, as no services currently use the account ID in endpoint resolution.

# v1.0.2 (2024-09-04)

* No change notes available for this release.

# v1.0.1 (2024-09-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2024-08-28)

* **Release**: New AWS service client module
* **Feature**: Introducing AWS Parallel Computing Service (AWS PCS), a new service makes it easy to setup and manage high performance computing (HPC) clusters, and build scientific and engineering models at virtually any scale on AWS.

