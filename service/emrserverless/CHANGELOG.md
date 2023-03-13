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

