# v1.4.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.4.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.4.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-01-25)

* **Feature**: Add returnCode, batchJobIdentifier in GetBatchJobExecution response, for user to view the batch job execution result & unique identifier from engine. Also removed unused headers from REST APIs

# v1.3.1 (2023-01-19)

* No change notes available for this release.

# v1.3.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.2.0 (2022-12-15)

* **Feature**: Adds an optional create-only `KmsKeyId` property to Environment and Application resources.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.6 (2022-12-12)

* No change notes available for this release.

# v1.1.5 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2022-11-16)

* No change notes available for this release.

# v1.1.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.9 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.8 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.7 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2022-06-08)

* **Release**: New AWS service client module
* **Feature**: AWS Mainframe Modernization service is a managed mainframe service and set of tools for planning, migrating, modernizing, and running mainframe workloads on AWS

