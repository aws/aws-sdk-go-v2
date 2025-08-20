# v1.9.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.9.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2025-07-15)

* **Feature**: Adds table bucket type to ListTableBucket and GetTableBucket API operations

# v1.5.0 (2025-06-23)

* **Feature**: S3 Tables now supports sort and z-order compaction strategies for Iceberg tables in addition to binpack.

# v1.4.2 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2025-06-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2025-06-06)

* **Feature**: S3 Tables now supports getting details about a table via its table ARN.

# v1.3.0 (2025-04-16)

* **Feature**: S3 Tables now supports setting encryption configurations on table buckets and tables. Encryption configurations can use server side encryption using AES256 or KMS customer-managed keys.

# v1.2.2 (2025-04-03)

* No change notes available for this release.

# v1.2.1 (2025-03-04.2)

* **Bug Fix**: Add assurance test for operation order.

# v1.2.0 (2025-02-27)

* **Feature**: Track credential providers via User-Agent Feature ids
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2025-02-18)

* **Bug Fix**: Bump go version to 1.22
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2025-02-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2025-01-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2025-01-30)

* **Feature**: You can now use the CreateTable API operation to create tables with schemas by adding an optional metadata argument.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2025-01-24)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.22.2.

# v1.0.4 (2025-01-17)

* **Bug Fix**: Fix bug where credentials weren't refreshed during retry loop.

# v1.0.3 (2025-01-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2025-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2024-12-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2024-12-03.2)

* **Release**: New AWS service client module
* **Feature**: Amazon S3 Tables deliver the first cloud object store with built-in open table format support, and the easiest way to store tabular data at scale.

