# v1.33.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.1 (2023-08-01)

* No change notes available for this release.

# v1.33.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.2 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.0 (2023-06-27)

* **Feature**: This release adds support to bypass SSO with the SAPOData connector when connecting to an SAP instance.

# v1.31.0 (2023-06-20)

* **Feature**: This release adds new API to reset connector metadata cache

# v1.30.2 (2023-06-15)

* No change notes available for this release.

# v1.30.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-06-01)

* **Feature**: Added ability to select DataTransferApiType for DescribeConnector and CreateFlow requests when using Async supported connectors. Added supportedDataTransferType to DescribeConnector/DescribeConnectors/ListConnector response.

# v1.29.1 (2023-05-04)

* No change notes available for this release.

# v1.29.0 (2023-05-02)

* **Feature**: This release adds new API to cancel flow executions.

# v1.28.0 (2023-04-28)

* **Feature**: Adds Jwt Support for Salesforce Credentials.

# v1.27.0 (2023-04-24)

* **Feature**: Increased the max length for RefreshToken and AuthCode from 2048 to 4096.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2023-04-17)

* **Feature**: This release adds a Client Token parameter to the following AppFlow APIs: Create/Update Connector Profile, Create/Update Flow, Start Flow, Register Connector, Update Connector Registration. The Client Token parameter allows idempotent operations for these APIs.

# v1.25.4 (2023-04-10)

* No change notes available for this release.

# v1.25.3 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2023-02-23)

* **Feature**: This release enables the customers to choose whether to use Private Link for Metadata and Authorization call when using a private Salesforce connections

# v1.24.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.24.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.24.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2023-01-19)

* **Feature**: Adding support for Salesforce Pardot connector in Amazon AppFlow.

# v1.23.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.22.0 (2022-12-16)

* **Feature**: This release updates the ListConnectorEntities API action so that it returns paginated responses that customers can retrieve with next tokens.

# v1.21.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-11-22)

* **Feature**: Adding support for Amazon AppFlow to transfer the data to Amazon Redshift databases through Amazon Redshift Data API service. This feature will support the Redshift destination connector on both public and private accessible Amazon Redshift Clusters and Amazon Redshift Serverless.

# v1.20.0 (2022-11-18)

* **Feature**: AppFlow provides a new API called UpdateConnectorRegistration to update a custom connector that customers have previously registered. With this API, customers no longer need to unregister and then register a connector to make an update.

# v1.19.0 (2022-11-17)

* **Feature**: AppFlow simplifies the preparation and cataloging of SaaS data into the AWS Glue Data Catalog where your data can be discovered and accessed by AWS analytics and ML services. AppFlow now also supports data field partitioning and file size optimization to improve query performance and reduce cost.

# v1.18.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-10-13)

* **Feature**: With this update, you can choose which Salesforce API is used by Amazon AppFlow to transfer data to or from your Salesforce account. You can choose the Salesforce REST API or Bulk API 2.0. You can also choose for Amazon AppFlow to pick the API automatically.

# v1.17.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.10 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.9 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.8 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.7 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.6 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.5 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.4 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.3 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.2 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-05-27)

* **Feature**: Adding the following features/changes: Parquet output that preserves typing from the source connector, Failed executions threshold before deactivation for scheduled flows, increasing max size of access and refresh token from 2048 to 4096

# v1.15.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-04-14)

* **Feature**: Enables users to pass custom token URL parameters for Oauth2 authentication during create connector profile

# v1.14.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service client model to latest release.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.9.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-08-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

