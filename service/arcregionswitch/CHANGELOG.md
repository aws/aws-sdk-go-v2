# v1.2.18 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.17 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.2.16 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.2.15 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.14 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.2.13 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.2.12 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.2.11 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.10 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.9 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.8 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.7 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.6 (2025-09-10)

* No change notes available for this release.

# v1.2.5 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.2.0 (2025-08-14)

* **Feature**: Endpoint rule test and documentation update.

# v1.1.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2025-08-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2025-08-01)

* **Release**: New AWS service client module
* **Feature**: This is the initial SDK release for Region switch

