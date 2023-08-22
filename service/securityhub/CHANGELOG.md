# v1.36.1 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.0 (2023-08-18)

* **Feature**: Added Inspector Lambda code Vulnerability section to ASFF, including GeneratorDetails, EpssScore, ExploitAvailable, and CodeVulnerabilities.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.4 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.3 (2023-08-10)

* No change notes available for this release.

# v1.35.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.1 (2023-08-01)

* No change notes available for this release.

# v1.35.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.0 (2023-07-25)

* **Feature**: Add support for CONTAINS and NOT_CONTAINS comparison operators for Automation Rules string filters and map filters

# v1.33.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.2 (2023-07-05)

* **Documentation**: Documentation updates for AWS Security Hub

# v1.33.1 (2023-06-15)

* No change notes available for this release.

# v1.33.0 (2023-06-13)

* **Feature**: Add support for Security Hub Automation Rules
* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.0 (2023-05-30)

* **Feature**: Added new resource detail objects to ASFF, including resources for AwsGuardDutyDetector, AwsAmazonMqBroker, AwsEventSchemasRegistry, AwsAppSyncGraphQlApi and AwsStepFunctionStateMachine.

# v1.31.1 (2023-05-26)

* No change notes available for this release.

# v1.31.0 (2023-05-04)

* **Feature**: Add support for Finding History.

# v1.30.3 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.2 (2023-04-10)

* No change notes available for this release.

# v1.30.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-03-27)

* **Feature**: Added new resource detail objects to ASFF, including resources for AwsEksCluster, AWSS3Bucket, AwsEc2RouteTable and AwsEC2Instance.

# v1.29.3 (2023-03-22)

* No change notes available for this release.

# v1.29.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-02-24)

* **Feature**: New Security Hub APIs and updates to existing APIs that help you consolidate control findings and enable and disable controls across all supported standards

# v1.28.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.28.4 (2023-02-21)

* **Documentation**: Documentation updates for AWS Security Hub

# v1.28.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.28.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2023-01-31)

* **Feature**: New fields have been added to the AWS Security Finding Format. Compliance.SecurityControlId is a unique identifier for a security control across standards. Compliance.AssociatedStandards contains all enabled standards in which a security control is enabled.

# v1.27.1 (2023-01-13)

* No change notes available for this release.

# v1.27.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.26.0 (2022-12-15)

* **Feature**: Added new resource details objects to ASFF, including resources for AwsEc2LaunchTemplate, AwsSageMakerNotebookInstance, AwsWafv2WebAcl and AwsWafv2RuleGroup.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2022-11-29)

* **Feature**: Adding StandardsManagedBy field to DescribeStandards API response

# v1.24.0 (2022-11-17)

* **Feature**: Added SourceLayerArn and SourceLayerHash field for security findings.  Updated AwsLambdaFunction Resource detail

# v1.23.8 (2022-11-11)

* **Documentation**: Documentation updates for Security Hub

# v1.23.7 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.6 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.5 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.4 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.3 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-08-22)

* **Feature**: Added new resource details objects to ASFF, including resources for AwsBackupBackupVault, AwsBackupBackupPlan and AwsBackupRecoveryPoint. Added FixAvailable, FixedInVersion and Remediation  to Vulnerability.

# v1.22.7 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.6 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.5 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.4 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.3 (2022-07-26)

* **Documentation**: Documentation updates for AWS Security Hub

# v1.22.2 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-06-16)

* **Feature**: Added Threats field for security findings. Added new resource details for ECS Container, ECS Task, RDS SecurityGroup, Kinesis Stream, EC2 TransitGateway, EFS AccessPoint, CloudFormation Stack, CloudWatch Alarm, VPC Peering Connection and WAF Rules

# v1.21.4 (2022-06-10)

* No change notes available for this release.

# v1.21.3 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.1 (2022-05-06)

* **Documentation**: Documentation updates for Security Hub API reference

# v1.21.0 (2022-04-25)

* **Feature**: Security Hub now lets you opt-out of auto-enabling the defaults standards (CIS and FSBP) in accounts that are auto-enabled with Security Hub via Security Hub's integration with AWS Organizations.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-04-05)

* **Feature**: Added additional ASFF details for RdsSecurityGroup AutoScalingGroup, ElbLoadBalancer, CodeBuildProject and RedshiftCluster.

# v1.19.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-01-28)

* **Feature**: Updated to latest API model.

# v1.16.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: API client updated

# v1.13.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-10-11)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-09-02)

* **Feature**: API client updated

# v1.9.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

