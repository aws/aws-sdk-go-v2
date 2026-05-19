# v1.0.4 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2026-03-26)

* **Bug Fix**: Fix a bug where a recorded clock skew could persist on the client even if the client and server clock ended up realigning.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2026-03-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2026-03-11)

* **Release**: New AWS service client module
* **Feature**: Introduced Amazon SimpleDB export functionality enabling domain data export to S3 in JSON format. Added three new APIs StartDomainExport, GetExport, and ListExports via SimpleDBv2 service. Supports cross-region exports and KMS encryption.

