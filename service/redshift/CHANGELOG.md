# v1.37.0 (2023-11-17)

* **Feature**: Updated SDK for Amazon Redshift, which you can use to configure a connection with IAM Identity Center to manage access to databases. With these, you can create a connection through a managed application. You can also change a managed application, delete it, or get information about an existing one.

# v1.36.0 (2023-11-15)

* **Feature**: The custom domain name SDK for Amazon Redshift provisioned clusters is updated with additional required parameters for modify and delete operations. Additionally, users can provide domain names with longer top-level domains.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.0 (2023-11-01)

* **Feature**: Added support for Multi-AZ deployments for Provisioned RA3 clusters that provide 99.99% SLA availability.
* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.0 (2023-10-27)

* **Feature**: added support to create a dual stack cluster

# v1.32.0 (2023-10-26)

* **Feature**: Add Redshift APIs GetResourcePolicy, DeleteResourcePolicy, PutResourcePolicy and DescribeInboundIntegrations for the new Amazon Redshift Zero-ETL integration feature, which can be used to control data ingress into Redshift namespace, and view inbound integrations.

# v1.31.0 (2023-10-24)

* **Feature**: **BREAKFIX**: Correct nullability and default value representation of various input fields across a large number of services. Calling code that references one or more of the affected fields will need to update usage accordingly. See [2162](https://github.com/aws/aws-sdk-go-v2/issues/2162).

# v1.30.0 (2023-10-16)

* **Feature**: Added support for managing credentials of provisioned cluster admin using AWS Secrets Manager.

# v1.29.7 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.6 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.1 (2023-08-01)

* No change notes available for this release.

# v1.29.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.2 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2023-06-20)

* **Feature**: Added support for custom domain names for Redshift Provisioned clusters. This feature enables customers to create a custom domain name and use ACM to generate fully secure connections to it.

# v1.27.13 (2023-06-15)

* No change notes available for this release.

# v1.27.12 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.11 (2023-05-04)

* No change notes available for this release.

# v1.27.10 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.9 (2023-04-10)

* No change notes available for this release.

# v1.27.8 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.7 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.6 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.5 (2023-02-24)

* **Documentation**: Documentation updates for Redshift API bringing it in line with IAM best practices.

# v1.27.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.27.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.2 (2023-02-03)

* **Documentation**: Corrects descriptions of the parameters for the API operations RestoreFromClusterSnapshot, RestoreTableFromClusterSnapshot, and CreateCluster.
* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade smithy to 1.27.2 and correct empty query list serialization.

# v1.27.1 (2023-01-23)

* No change notes available for this release.

# v1.27.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.26.18 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.17 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.16 (2022-11-22)

* No change notes available for this release.

# v1.26.15 (2022-11-16)

* No change notes available for this release.

# v1.26.14 (2022-11-10)

* No change notes available for this release.

# v1.26.13 (2022-10-27)

* **Documentation**: This release clarifies use for the ElasticIp parameter of the CreateCluster and RestoreFromClusterSnapshot APIs.

# v1.26.12 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.11 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.10 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.9 (2022-09-14)

* **Documentation**: This release updates documentation for AQUA features and other description updates.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.8 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.7 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.6 (2022-08-30)

* No change notes available for this release.

# v1.26.5 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2022-07-11)

* **Feature**: This release adds a new --snapshot-arn field for describe-cluster-snapshots, describe-node-configuration-options, restore-from-cluster-snapshot, authorize-snapshot-acsess, and revoke-snapshot-acsess APIs. It allows customers to give a Redshift snapshot ARN or a Redshift Serverless ARN as input.

# v1.25.2 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2022-06-08)

* **Feature**: Adds new API GetClusterCredentialsWithIAM to return temporary credentials.

# v1.24.2 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-05-06)

* **Feature**: Introduces new field 'LoadSampleData' in CreateCluster operation. Customers can now specify 'LoadSampleData' option during creation of a cluster, which results in loading of sample data in the cluster that is created.

# v1.23.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-04-19)

* **Feature**: Introduces new fields for LogDestinationType and LogExports on EnableLogging requests and Enable/Disable/DescribeLogging responses. Customers can now select CloudWatch Logs as a destination for their Audit Logs.

# v1.22.2 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-03-23)

* **Feature**: This release adds a new [--encrypted | --no-encrypted] field in restore-from-cluster-snapshot API. Customers can now restore an unencrypted snapshot to a cluster encrypted with AWS Managed Key or their own KMS key.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: API client updated
* **Feature**: Updated to latest service endpoints

# v1.16.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2021-11-30)

* **Feature**: API client updated

# v1.15.0 (2021-11-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.13.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-07-15)

* **Feature**: The ErrorCode method on generated service error types has been corrected to match the API model.
* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-07-01)

* **Feature**: API client updated

# v1.7.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2021-06-04)

* No change notes available for this release.

# v1.6.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

