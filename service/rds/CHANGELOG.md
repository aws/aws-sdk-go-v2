# v1.63.0 (2023-11-17)

* **Feature**: This release adds support for option groups and replica enhancements to Amazon RDS Custom.

# v1.62.4 (2023-11-15)

* **Documentation**: Updates Amazon RDS documentation for support for upgrading RDS for MySQL snapshots from version 5.7 to version 8.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.62.3 (2023-11-10)

* **Documentation**: Updates Amazon RDS documentation for zero-ETL integrations.

# v1.62.2 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.62.1 (2023-11-08)

* **Documentation**: This Amazon RDS release adds support for patching the OS of an RDS Custom for Oracle DB instance. You can now upgrade the database or operating system using the modify-db-instance command.

# v1.62.0 (2023-11-07)

* **Feature**: This Amazon RDS release adds support for the multi-tenant configuration. In this configuration, an RDS DB instance can contain multiple tenant databases. In RDS for Oracle, a tenant database is a pluggable database (PDB).

# v1.61.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Feature**: This release adds support for customized networking resources to Amazon RDS Custom.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.60.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.59.0 (2023-10-30)

* **Feature**: This release launches the CreateIntegration, DeleteIntegration, and DescribeIntegrations APIs to manage zero-ETL Integrations.

# v1.58.0 (2023-10-24)

* **Feature**: **BREAKFIX**: Correct nullability and default value representation of various input fields across a large number of services. Calling code that references one or more of the affected fields will need to update usage accordingly. See [2162](https://github.com/aws/aws-sdk-go-v2/issues/2162).

# v1.57.0 (2023-10-18)

* **Feature**: This release adds support for upgrading the storage file system configuration on the DB instance using a blue/green deployment or a read replica.

# v1.56.0 (2023-10-12)

* **Feature**: This release adds support for adding a dedicated log volume to open-source RDS instances.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.55.2 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.55.1 (2023-10-05)

* **Documentation**: Updates Amazon RDS documentation for corrections and minor improvements.

# v1.55.0 (2023-10-02)

* **Feature**: Adds DefaultCertificateForNewLaunches field in the DescribeCertificates API response.

# v1.54.0 (2023-09-05)

* **Feature**: Add support for feature integration with AWS Backup.

# v1.53.0 (2023-08-24)

* **Feature**: This release updates the supported versions for Percona XtraBackup in Aurora MySQL.

# v1.52.0 (2023-08-22)

* **Feature**: Adding parameters to CreateCustomDbEngineVersion reserved for future use.

# v1.51.0 (2023-08-21)

* **Feature**: Adding support for RDS Aurora Global Database Unplanned Failover
* **Dependency Update**: Updated to the latest SDK module versions

# v1.50.3 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.50.2 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.50.1 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.50.0 (2023-08-01)

* **Feature**: Added support for deleted clusters PiTR.

# v1.49.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Feature**: This release adds support for Aurora MySQL local write forwarding, which allows for forwarding of write operations from reader DB instances to the writer DB instance.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.48.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.48.0 (2023-07-25)

* **Feature**: This release adds support for monitoring storage optimization progress on the DescribeDBInstances API.

# v1.47.0 (2023-07-21)

* **Feature**: Adds support for the DBSystemID parameter of CreateDBInstance to RDS Custom for Oracle.

# v1.46.2 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.46.1 (2023-07-06)

* **Documentation**: Updates Amazon RDS documentation for creating DB instances and creating Aurora global clusters.

# v1.46.0 (2023-06-28)

* **Feature**: Amazon Relational Database Service (RDS) now supports joining a RDS for SQL Server instance to a self-managed Active Directory.

# v1.45.3 (2023-06-23)

* **Documentation**: Documentation improvements for create, describe, and modify DB clusters and DB instances.

# v1.45.2 (2023-06-15)

* No change notes available for this release.

# v1.45.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.45.0 (2023-05-31)

* **Feature**: This release adds support for changing the engine for Oracle using the ModifyDbInstance API

# v1.44.1 (2023-05-18)

* **Documentation**: RDS documentation update for the EngineVersion parameter of ModifyDBSnapshot

# v1.44.0 (2023-05-10)

* **Feature**: Amazon Relational Database Service (RDS) updates for the new Aurora I/O-Optimized storage type for Amazon Aurora DB clusters

# v1.43.3 (2023-05-04)

* No change notes available for this release.

# v1.43.2 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.43.1 (2023-04-19)

* **Documentation**: Adds support for the ImageId parameter of CreateCustomDBEngineVersion to RDS Custom for Oracle

# v1.43.0 (2023-04-14)

* **Feature**: This release adds support of modifying the engine mode of database clusters.

# v1.42.3 (2023-04-10)

* No change notes available for this release.

# v1.42.2 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.1 (2023-04-06)

* **Documentation**: Adds and updates the SDK examples

# v1.42.0 (2023-03-29)

* **Feature**: Add support for creating a read replica DB instance from a Multi-AZ DB cluster.

# v1.41.0 (2023-03-24)

* **Feature**: Added error code CreateCustomDBEngineVersionFault for when the create custom engine version for Custom engines fails.

# v1.40.7 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.6 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.40.4 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.3 (2023-02-15)

* **Documentation**: Database Activity Stream support for RDS for SQL Server.

# v1.40.2 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade smithy to 1.27.2 and correct empty query list serialization.

# v1.40.1 (2023-01-23)

* No change notes available for this release.

# v1.40.0 (2023-01-10)

* **Feature**: This release adds support for configuring allocated storage on the CreateDBInstanceReadReplica, RestoreDBInstanceFromDBSnapshot, and RestoreDBInstanceToPointInTime APIs.

# v1.39.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).
* **Feature**: This release adds support for specifying which certificate authority (CA) to use for a DB instance's server certificate during DB instance creation, as well as other CA enhancements.

# v1.38.0 (2022-12-28)

* **Feature**: This release adds support for Custom Engine Version (CEV) on RDS Custom SQL Server.

# v1.37.0 (2022-12-22)

* **Feature**: Add support for managing master user password in AWS Secrets Manager for the DBInstance and DBCluster.

# v1.36.0 (2022-12-19)

* **Feature**: Add support for --enable-customer-owned-ip to RDS create-db-instance-read-replica API for RDS on Outposts.

# v1.35.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.0 (2022-12-13)

* **Feature**: This deployment adds ClientPasswordAuthType field to the Auth structure of the DBProxy.

# v1.34.0 (2022-12-12)

* **Feature**: Update the RDS API model to support copying option groups during the CopyDBSnapshot operation

# v1.33.0 (2022-12-06)

* **Feature**: This release adds the BlueGreenDeploymentNotFoundFault to the AddTagsToResource, ListTagsForResource, and RemoveTagsFromResource operations.

# v1.32.0 (2022-12-05)

* **Feature**: This release adds the InvalidDBInstanceStateFault to the RestoreDBClusterFromSnapshot operation.

# v1.31.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2022-11-28)

* **Feature**: This release enables new Aurora and RDS feature called Blue/Green Deployments that makes updates to databases safer, simpler and faster.

# v1.30.1 (2022-11-22)

* No change notes available for this release.

# v1.30.0 (2022-11-16)

* **Feature**: This release adds support for container databases (CDBs) to Amazon RDS Custom for Oracle. A CDB contains one PDB at creation. You can add more PDBs using Oracle SQL. You can also customize your database installation by setting the Oracle base, Oracle home, and the OS user name and group.

# v1.29.0 (2022-11-14)

* **Feature**: This release adds support for restoring an RDS Multi-AZ DB cluster snapshot to a Single-AZ deployment or a Multi-AZ DB instance deployment.

# v1.28.1 (2022-11-10)

* No change notes available for this release.

# v1.28.0 (2022-11-01)

* **Feature**: Relational Database Service - This release adds support for configuring Storage Throughput on RDS database instances.

# v1.27.0 (2022-10-25)

* **Feature**: Relational Database Service - This release adds support for exporting DB cluster data to Amazon S3.

# v1.26.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2022-09-19)

* **Feature**: This release adds support for Amazon RDS Proxy with SQL Server compatibility.

# v1.25.6 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.5 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.4 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.3 (2022-08-30)

* No change notes available for this release.

# v1.25.2 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2022-08-26)

* **Documentation**: Removes support for RDS Custom from DBInstanceClass in ModifyDBInstance

# v1.25.0 (2022-08-23)

* **Feature**: RDS for Oracle supports Oracle Data Guard switchover and read replica backups.

# v1.24.0 (2022-08-17)

* **Feature**: Adds support for Internet Protocol Version 6 (IPv6) for RDS Aurora database clusters.

# v1.23.6 (2022-08-14)

* **Documentation**: Adds support for RDS Custom to DBInstanceClass in ModifyDBInstance

# v1.23.5 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.4 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.3 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-07-26)

* **Documentation**: Adds support for using RDS Proxies with RDS for MariaDB databases.

# v1.23.0 (2022-07-22)

* **Feature**: This release adds the "ModifyActivityStream" API with support for audit policy state locking and unlocking.

# v1.22.1 (2022-07-21)

* **Documentation**: Adds support for creating an RDS Proxy for an RDS for MariaDB database.

# v1.22.0 (2022-07-05)

* **Feature**: Adds waiters support for DBCluster.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.5 (2022-07-01)

* **Documentation**: Adds support for additional retention periods to Performance Insights.

# v1.21.4 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.3 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.1 (2022-05-06)

* **Documentation**: Various documentation improvements.

# v1.21.0 (2022-04-29)

* **Feature**: Feature - Adds support for Internet Protocol Version 6 (IPv6) on RDS database instances.

# v1.20.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-04-20)

* **Feature**: Added a new cluster-level attribute to set the capacity range for Aurora Serverless v2 instances.

# v1.19.0 (2022-04-15)

* **Feature**: Removes Amazon RDS on VMware with the deletion of APIs related to Custom Availability Zones and Media installation

# v1.18.4 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-03-15)

* **Documentation**: Various documentation improvements

# v1.18.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Documentation**: Updated service client model to latest release.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-01-14)

* **Feature**: Updated API models
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Documentation**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: Updated to latest service endpoints

# v1.13.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-11-30)

* **Feature**: API client updated

# v1.12.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.
* **Feature**: Updated service to latest API model.
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.11.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-09-17)

* **Feature**: Updated API client and endpoints to latest revision.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-07-15)

* **Feature**: The ErrorCode method on generated service error types has been corrected to match the API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2021-06-11)

* **Documentation**: Updated to latest API model.

# v1.4.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

