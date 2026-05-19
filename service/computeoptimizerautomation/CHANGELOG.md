# v1.1.1 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2026-04-21)

* **Feature**: This release adds Smithy RPC v2 CBOR as an additional protocol alongside the existing AWS JSON 1.0. The SDK will prioritize its most performant protocol.

# v1.0.9 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.8 (2026-03-26)

* **Bug Fix**: Fix a bug where a recorded clock skew could persist on the client even if the client and server clock ended up realigning.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.7 (2026-03-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2026-03-03)

* **Dependency Update**: Bump minimum Go version to 1.24
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2026-02-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2026-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.0.1 (2025-11-25)

* No change notes available for this release.

# v1.0.0 (2025-11-21)

* **Release**: New AWS service client module
* **Feature**: Initial release of AWS Compute Optimizer Automation. Create automation rules to implement recommended actions on a recurring schedule based on your specified criteria. Supported actions include: snapshot and delete unattached EBS volumes and upgrade volume types to the latest generation.

