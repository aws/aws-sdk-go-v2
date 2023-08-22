# v1.10.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2023-08-01)

* No change notes available for this release.

# v1.10.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2023-07-25)

* **Feature**: This release adds support for publishing application logs to CloudWatch.

# v1.8.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2023-06-27)

* **Feature**: This release adds support to update the release label of an EMR Serverless application to upgrade it to a different version of Amazon EMR via UpdateApplication API.

# v1.7.6 (2023-06-15)

* No change notes available for this release.

# v1.7.5 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.4 (2023-05-11)

* No change notes available for this release.

# v1.7.3 (2023-05-04)

* No change notes available for this release.

# v1.7.2 (2023-05-01)

* No change notes available for this release.

# v1.7.1 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2023-04-17)

* **Feature**: The GetJobRun API has been updated to include the job's billed resource utilization. This utilization shows the aggregate vCPU, memory and storage that AWS has billed for the job run. The billed resources include a 1-minute minimum usage for workers, plus additional storage over 20 GB per worker.

# v1.6.0 (2023-04-11)

* **Feature**: This release extends GetJobRun API to return job run timeout (executionTimeoutMinutes) specified during StartJobRun call (or default timeout of 720 minutes if none was specified).

# v1.5.8 (2023-04-10)

* No change notes available for this release.

# v1.5.7 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.6 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.5.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.5.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).
* **Feature**: Adds support for customized images. You can now provide runtime images when creating or updating EMR Serverless Applications.

# v1.4.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2022-11-17)

* **Feature**: Adds support for AWS Graviton2 based applications. You can now select CPU architecture when creating new applications or updating existing ones.

# v1.3.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2022-09-29)

* **Feature**: This release adds API support to debug Amazon EMR Serverless jobs in real-time with live application UIs

# v1.2.4 (2022-09-27)

* No change notes available for this release.

# v1.2.3 (2022-09-23)

* No change notes available for this release.

# v1.2.2 (2022-09-21)

* No change notes available for this release.

# v1.2.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.10 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.9 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.8 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.7 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.6 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.5 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2022-05-27)

* **Feature**: This release adds support for Amazon EMR Serverless, a serverless runtime environment that simplifies running analytics applications using the latest open source frameworks such as Apache Spark and Apache Hive.

# v1.0.0 (2022-05-26)

* **Release**: New AWS service client module
* **Feature**: This release adds support for Amazon EMR Serverless, a serverless runtime environment that simplifies running analytics applications using the latest open source frameworks such as Apache Spark and Apache Hive.

