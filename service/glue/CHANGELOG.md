# v1.61.3 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.61.2 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.61.1 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.61.0 (2023-08-15)

* **Feature**: AWS Glue Crawlers can now accept SerDe overrides from a custom csv classifier. The two SerDe options are LazySimpleSerDe and OpenCSVSerDe. In case, the user wants crawler to do the selection, "None" can be selected for this purpose.

# v1.60.1 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.60.0 (2023-08-02)

* **Feature**: This release includes additional Glue Streaming KAKFA SASL property types.

# v1.59.1 (2023-08-01)

* No change notes available for this release.

# v1.59.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.58.2 (2023-07-28.2)

* No change notes available for this release.

# v1.58.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.58.0 (2023-07-26)

* **Feature**: Release Glue Studio Snowflake Connector Node for SDK/CLI

# v1.57.0 (2023-07-24)

* **Feature**: Added support for Data Preparation Recipe node in Glue Studio jobs

# v1.56.0 (2023-07-21)

* **Feature**: This release adds support for AWS Glue Crawler with Apache Hudi Tables, allowing Crawlers to discover Hudi Tables in S3 and register them in Glue Data Catalog for query engines to query against.

# v1.55.0 (2023-07-17)

* **Feature**: Adding new supported permission type flags to get-unfiltered endpoints that callers may pass to indicate support for enforcing Lake Formation fine-grained access control on nested column attributes.

# v1.54.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.54.0 (2023-07-07)

* **Feature**: This release enables customers to create new Apache Iceberg tables and associated metadata in Amazon S3 by using native AWS Glue CreateTable operation.

# v1.53.0 (2023-06-29)

* **Feature**: This release adds support for AWS Glue Crawler with Iceberg Tables, allowing Crawlers to discover Iceberg Tables in S3 and register them in Glue Data Catalog for query engines to query against.

# v1.52.0 (2023-06-26)

* **Feature**: Timestamp Starting Position For Kinesis and Kafka Data Sources in a Glue Streaming Job

# v1.51.0 (2023-06-19)

* **Feature**: This release adds support for creating cross region table/database resource links

# v1.50.2 (2023-06-15)

* No change notes available for this release.

# v1.50.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.50.0 (2023-05-30)

* **Feature**: Added Runtime parameter to allow selection of Ray Runtime

# v1.49.0 (2023-05-25)

* **Feature**: Added ability to create data quality rulesets for shared, cross-account Glue Data Catalog tables. Added support for dataset comparison rules through a new parameter called AdditionalDataSources. Enhanced the data quality results with a map containing profiled metric values.

# v1.48.0 (2023-05-16)

* **Feature**: Add Support for Tags for Custom Entity Types

# v1.47.0 (2023-05-09)

* **Feature**: This release adds AmazonRedshift Source and Target nodes in addition to DynamicTransform OutputSchemas

# v1.46.0 (2023-05-08)

* **Feature**: Support large worker types G.4x and G.8x for Glue Spark

# v1.45.5 (2023-05-04)

* No change notes available for this release.

# v1.45.4 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.45.3 (2023-04-10)

* No change notes available for this release.

# v1.45.2 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.45.1 (2023-04-06)

* No change notes available for this release.

# v1.45.0 (2023-04-03)

* **Feature**: Add support for database-level federation

# v1.44.0 (2023-03-30)

* **Feature**: This release adds support for AWS Glue Data Quality, which helps you evaluate and monitor the quality of your data and includes the API for creating, deleting, or updating data quality rulesets, runs and evaluations.

# v1.43.4 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.43.3 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.43.2 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.43.1 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.43.0 (2023-02-17)

* **Feature**: Release of Delta Lake Data Lake Format for Glue Studio Service

# v1.42.0 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Feature**: Fix DirectJDBCSource not showing up in CLI code gen
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.41.0 (2023-02-08)

* **Feature**: DirectJDBCSource + Glue 4.0 streaming options

# v1.40.2 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.1 (2023-01-31)

* No change notes available for this release.

# v1.40.0 (2023-01-19)

* **Feature**: Release Glue Studio Hudi Data Lake Format for SDK/CLI

# v1.39.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.38.1 (2022-12-19)

* No change notes available for this release.

# v1.38.0 (2022-12-15)

* **Feature**: This release adds support for AWS Glue Crawler with native DeltaLake tables, allowing Crawlers to classify Delta Lake format tables and catalog them for query engines to query against.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.0 (2022-11-30)

* **Feature**: This release adds support for AWS Glue Data Quality, which helps you evaluate and monitor the quality of your data and includes the API for creating, deleting, or updating data quality rulesets, runs and evaluations.

# v1.36.0 (2022-11-29)

* **Feature**: This release allows the creation of Custom Visual Transforms (Dynamic Transforms) to be created via AWS Glue CLI/SDK.

# v1.35.0 (2022-11-18)

* **Feature**: AWSGlue Crawler - Adding support for Table and Column level Comments with database level datatypes for JDBC based crawler.

# v1.34.1 (2022-11-11)

* **Documentation**: Added links related to enabling job bookmarks.

# v1.34.0 (2022-10-27)

* **Feature**: Added support for custom datatypes when using custom csv classifier.

# v1.33.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.0 (2022-10-05)

* **Feature**: This SDK release adds support to sync glue jobs with source control provider. Additionally, a new parameter called SourceControlDetails will be added to Job model.

# v1.32.0 (2022-09-22)

* **Feature**: Added support for S3 Event Notifications for Catalog Target Crawlers.

# v1.31.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.4 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.3 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.2 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.1 (2022-08-25)

* No change notes available for this release.

# v1.30.0 (2022-08-11)

* **Feature**: Add support for Python 3.9 AWS Glue Python Shell jobs
* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.1 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2022-08-08)

* **Feature**: Add an option to run non-urgent or non-time sensitive Glue Jobs on spare capacity
* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.2 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2022-07-19)

* **Documentation**: Documentation updates for AWS Glue Job Timeout and Autoscaling

# v1.28.0 (2022-07-14)

* **Feature**: This release adds an additional worker type for Glue Streaming jobs.

# v1.27.1 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2022-06-30)

* **Feature**: This release adds tag as an input of CreateDatabase

# v1.26.1 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2022-05-17)

* **Feature**: This release adds a new optional parameter called codeGenNodeConfiguration to CRUD job APIs that allows users to manage visual jobs via APIs. The updated CreateJob and UpdateJob will create jobs that can be viewed in Glue Studio as a visual graph. GetJob can be used to get codeGenNodeConfiguration.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.2 (2022-04-26)

* **Documentation**: This release adds documentation for the APIs to create, read, delete, list, and batch read of AWS Glue custom patterns, and for Lake Formation configuration settings in the AWS Glue crawler.

# v1.24.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-04-21)

* **Feature**: This release adds APIs to create, read, delete, list, and batch read of Glue custom entity types

# v1.23.0 (2022-04-14)

* **Feature**: Auto Scaling for Glue version 3.0 and later jobs to dynamically scale compute resources. This SDK change provides customers with the auto-scaled DPU usage

# v1.22.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-03-18)

* **Feature**: Added 9 new APIs for AWS Glue Interactive Sessions: ListSessions, StopSession, CreateSession, GetSession, DeleteSession, RunStatement, GetStatement, ListStatements, CancelStatement

# v1.21.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-01-14)

* **Feature**: Updated API models
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-01-07)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.16.0 (2021-12-02)

* **Feature**: API client updated
* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.

# v1.14.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-10-11)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-07-01)

* **Feature**: API client updated

# v1.7.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-06-11)

* **Feature**: Updated to latest API model.

# v1.5.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

