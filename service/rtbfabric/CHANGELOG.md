# v1.2.1 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.2.0 (2025-11-11)

* **Feature**: Added LogSettings and LinkAttribute fields to external links
* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.1.2 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.1.1 (2025-10-30)

* **Documentation**: RTB Fabric documentation update.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2025-10-24)

* **Feature**: Add support for custom rate limits.

# v1.0.1 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2025-10-22)

* **Release**: New AWS service client module
* **Feature**: Update for general availability of AWS RTB Fabric service.

