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

