# v1.1.5 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2026-03-26)

* **Bug Fix**: Fix a bug where a recorded clock skew could persist on the client even if the client and server clock ended up realigning.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2026-03-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2026-03-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2026-03-03)

* **Feature**: Adds the Resold Unified Operations support plan and removes the Resold Business support plan in the CreateRelationship and UpdateRelationship APIs
* **Dependency Update**: Bump minimum Go version to 1.24
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2026-02-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2026-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.0.2 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.0.1 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2025-11-19)

* **Release**: New AWS service client module
* **Feature**: Initial GA launch of Partner Central Channel

