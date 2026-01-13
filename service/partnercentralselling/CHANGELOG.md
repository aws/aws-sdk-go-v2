# v1.13.1 (2026-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2025-12-10)

* **Feature**: Adds support for the new Project.AwsPartition field on Opportunity and AWS Opportunity Summary. Use this field to specify the AWS partition where the opportunity will be deployed.

# v1.12.0 (2025-12-08)

* **Feature**: Deal Sizing Service for AI-based deal size estimation with AWS service-level breakdown, supporting Expansion and Migration deals across Technology, and Reseller partner cohorts, including Pricing Calculator AddOn for MAP deals and funding incentives.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.11.0 (2025-12-01)

* **Feature**: New Features:
Lead Management APIs for capturing and nurturing leads
Lead invitation support for partner collaboration
Lead-to-opportunity conversion operations
AWS Marketplace OfferSets support for opportunities

# v1.10.16 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.10.15 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.14 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.10.13 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.10.12 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.10.11 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.10 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.9 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.8 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.7 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.6 (2025-09-10)

* No change notes available for this release.

# v1.10.5 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.4 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.3 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.10.0 (2025-08-13)

* **Feature**: Add Tagging Support for Opportunity resources

# v1.9.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2025-08-08)

* **Feature**: Removed incorrect endpoint tests

# v1.7.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2025-06-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2025-05-21)

* **Feature**: Modified validation to allow expectedCustomerSpend array with zero elements in Partner Opportunity operations.

# v1.4.2 (2025-04-03)

* No change notes available for this release.

# v1.4.1 (2025-03-04.2)

* **Bug Fix**: Add assurance test for operation order.

# v1.4.0 (2025-02-27)

* **Feature**: Track credential providers via User-Agent Feature ids
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.6 (2025-02-18)

* **Bug Fix**: Bump go version to 1.22
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.5 (2025-02-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.4 (2025-01-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.3 (2025-01-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.2 (2025-01-24)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.22.2.

# v1.3.1 (2025-01-17)

* **Bug Fix**: Fix bug where credentials weren't refreshed during retry loop.

# v1.3.0 (2025-01-15)

* **Feature**: Add Tagging support for ResourceSnapshotJob resources
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2025-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2024-12-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2024-12-05)

* **Feature**: Introducing the preview of new partner central selling APIs designed to transform how AWS partners collaborate and co-sell with multiple partners. This enables multiple partners to seamlessly engage and jointly pursue customer opportunities, fostering a new era of collaborative selling.

# v1.1.2 (2024-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2024-11-18)

* **Dependency Update**: Update to smithy-go v1.22.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2024-11-15)

* **Feature**: Announcing AWS Partner Central API for Selling: This service launch Introduces new APIs for co-selling opportunity management and related functions. Key features include notifications, a dynamic sandbox for testing, and streamlined validations.

# v1.0.0 (2024-11-14)

* **Release**: New AWS service client module
* **Feature**: Announcing AWS Partner Central API for Selling: This service launch Introduces new APIs for co-selling opportunity management and related functions. Key features include notifications, a dynamic sandbox for testing, and streamlined validations.

