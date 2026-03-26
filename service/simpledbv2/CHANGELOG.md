# v1.0.2 (2026-03-26)

* **Bug Fix**: Fix a bug where a recorded clock skew could persist on the client even if the client and server clock ended up realigning.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2026-03-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2026-03-11)

* **Release**: New AWS service client module
* **Feature**: Introduced Amazon SimpleDB export functionality enabling domain data export to S3 in JSON format. Added three new APIs StartDomainExport, GetExport, and ListExports via SimpleDBv2 service. Supports cross-region exports and KMS encryption.

