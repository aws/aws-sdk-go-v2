# v1.15.7 (2024-01-10)

* **Documentation**: Updates to ConfigParameter for RSS workgroup, removal of use_fips_ssl

# v1.15.6 (2024-01-05)

* **Documentation**: use_fips_ssl and require_ssl parameter support for Workgroup, UpdateWorkgroup, and CreateWorkgroup

# v1.15.5 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.4 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.15.3 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.2 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.15.1 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2023-11-30)

* **Feature**: This release adds the following support for Amazon Redshift Serverless: 1) cross-account cross-VPCs, 2) copying snapshots across Regions, 3) scheduling snapshot creation, and 4) restoring tables from a recovery point.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.3 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.2 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.13.1 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2023-11-17)

* **Feature**: Updated SDK for Amazon Redshift Serverless, which provides the ability to configure a connection with IAM Identity Center to manage user and group access to databases.

# v1.12.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2023-11-08)

* **Feature**: Added a new parameter in the workgroup that helps you control your cost for compute resources. This feature provides a ceiling for RPUs that Amazon Redshift Serverless can scale up to. When automatic compute scaling is required, having a higher value for MaxRPU can enhance query throughput.

# v1.11.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2023-10-30)

* **Feature**: Added support for custom domain names for Amazon Redshift Serverless workgroups. This feature enables customers to create a custom domain name and use ACM to generate fully secure connections to it.

# v1.8.0 (2023-10-23)

* **Feature**: This release adds support for customers to see the patch version and workgroup version in Amazon Redshift Serverless.

# v1.7.0 (2023-10-16)

* **Feature**: Added support for managing credentials of serverless namespace admin using AWS Secrets Manager.

# v1.6.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2023-09-18)

* **Announcement**: [BREAKFIX] Change in MaxResults datatype from value to pointer type in cognito-sync service.
* **Feature**: Adds several endpoint ruleset changes across all models: smaller rulesets, removed non-unique regional endpoints, fixes FIPS and DualStack endpoints, and make region not required in SDK::Endpoint. Additional breakfix to cognito-sync field.

# v1.5.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2023-08-01)

* No change notes available for this release.

# v1.5.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.16 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.15 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.14 (2023-06-15)

* No change notes available for this release.

# v1.4.13 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.12 (2023-05-04)

* No change notes available for this release.

# v1.4.11 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.10 (2023-04-10)

* No change notes available for this release.

# v1.4.9 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.8 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.7 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.6 (2023-03-09)

* No change notes available for this release.

# v1.4.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.4.4 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.4.2 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2023-01-25)

* **Documentation**: Added query monitoring rules as possible parameters for create and update workgroup operations.

# v1.4.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.3.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2022-12-02)

* **Feature**: Add Table Level Restore operations for Amazon Redshift Serverless. Add multi-port support for Amazon Redshift Serverless endpoints. Add Tagging support to Snapshots and Recovery Points in Amazon Redshift Serverless.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.13 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.12 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.11 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.10 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.9 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.8 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.7 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.6 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2022-07-11)

* **Documentation**: Removed prerelease language for GA launch.

# v1.2.1 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2022-06-29)

* **Feature**: Add new API operations for Amazon Redshift Serverless, a new way of using Amazon Redshift without needing to manually manage provisioned clusters. The new operations let you interact with Redshift Serverless resources, such as create snapshots, list VPC endpoints, delete resource policies, and more.
* **Dependency Update**: Updated to the latest SDK module versions

