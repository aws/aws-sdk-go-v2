# v1.7.4 (2026-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2026-06-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2026-05-29)

* **Dependency Update**: Update to smithy-go v1.26.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2026-05-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2026-05-27)

* **Feature**: Adding new BDD representation of endpoint ruleset

# v1.6.0 (2026-05-19)

* **Feature**: This release is to deprecate 'inboundLinksCount' field in GetResponderGateway response and introduce the new field 'linksRequestedCount' to replace it.

# v1.5.0 (2026-05-13)

* **Feature**: Customers can now configure custom domain names for their RTB Fabric gateways. This enables partners to use their own branded domain for RTB traffic instead of the default rtbfabric endpoint

# v1.4.2 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2026-04-10)

* **Feature**: Adds optional health check configuration for Responder Gateways with ASG Managed Endpoints. When provided, RTB Fabric continuously probes customers' instance IPs and routes traffic only to healthy ones, reducing errors during deployments, scaling events, and instance failures.

# v1.3.0 (2026-04-07)

* **Feature**: AWS RTB Fabric External Responder gateways now support HTTP in addition to HTTPS for inbound external links. Gateways can accept bid requests on port 80 or serve both protocols simultaneously via listener configuration, giving customers flexible transport options for their bidding infrastructure

# v1.2.10 (2026-03-26)

* **Bug Fix**: Fix a bug where a recorded clock skew could persist on the client even if the client and server clock ended up realigning.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.9 (2026-03-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.8 (2026-03-03)

* **Dependency Update**: Bump minimum Go version to 1.24
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.7 (2026-02-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.6 (2026-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.2.3 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.2.2 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

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

