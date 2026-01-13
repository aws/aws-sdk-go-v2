# v1.18.4 (2026-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.18.1 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.18.0 (2025-11-20)

* **Feature**: Amazon CloudWatch Application Signals now supports un-instrumented services discovery, cross-account views, and change history, helping SRE and DevOps teams monitor and troubleshoot their large-scale distributed applications.

# v1.17.5 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.4 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.17.3 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.17.2 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.17.1 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2025-10-28)

* **Feature**: Added support for CloudWatch Synthetics Canary resources in ListAuditFindings API. This enhancement allows customers to retrieve audit findings specifically for CloudWatch Synthetics canaries and enables service-canary correlation analysis.

# v1.16.2 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.1 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2025-09-30)

* **Feature**: Amazon CloudWatch Application Signals is introducing the Application Map to give users a more comprehensive view of their service health. Users will now be able to group services, track their latest deployments, and view automated audit findings concerning service performance.

# v1.15.8 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.7 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.6 (2025-09-10)

* No change notes available for this release.

# v1.15.5 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.4 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.3 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.15.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2025-08-08)

* **Feature**: Removed incorrect endpoint tests

# v1.13.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.4 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.3 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.2 (2025-06-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2025-04-03)

* No change notes available for this release.

# v1.11.0 (2025-04-02)

* **Feature**: Application Signals now supports creating Service Level Objectives on service dependencies. Users can now create or update SLOs on discovered service dependencies to monitor their standard application metrics.

# v1.10.0 (2025-03-17)

* **Feature**: This release adds support for adding, removing, and listing SLO time exclusion windows with the BatchUpdateExclusionWindows and ListServiceLevelObjectiveExclusionWindows APIs.

# v1.9.1 (2025-03-04.2)

* **Bug Fix**: Add assurance test for operation order.

# v1.9.0 (2025-02-27)

* **Feature**: Track credential providers via User-Agent Feature ids
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2025-02-26)

* **Feature**: This release adds API support for reading Service Level Objectives and Services from monitoring accounts, from SLO and Service-scoped operations, including ListServices and ListServiceLevelObjectives.

# v1.7.11 (2025-02-18)

* **Bug Fix**: Bump go version to 1.22
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.10 (2025-02-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.9 (2025-01-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.8 (2025-01-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.7 (2025-01-24)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.22.2.

# v1.7.6 (2025-01-17)

* **Bug Fix**: Fix bug where credentials weren't refreshed during retry loop.

# v1.7.5 (2025-01-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.4 (2025-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2024-12-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2024-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2024-11-18)

* **Dependency Update**: Update to smithy-go v1.22.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2024-11-13)

* **Feature**: Amazon CloudWatch Application Signals now supports creating Service Level Objectives with burn rates. Users can now create or update SLOs with burn rate configurations to meet their specific business requirements.

# v1.6.5 (2024-11-07)

* **Bug Fix**: Adds case-insensitive handling of error message fields in service responses

# v1.6.4 (2024-11-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2024-10-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2024-10-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2024-10-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2024-10-04)

* **Feature**: Add support for HTTP client metrics.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2024-10-03)

* No change notes available for this release.

# v1.5.3 (2024-09-27)

* No change notes available for this release.

# v1.5.2 (2024-09-25)

* No change notes available for this release.

# v1.5.1 (2024-09-23)

* No change notes available for this release.

# v1.5.0 (2024-09-20)

* **Feature**: Add tracing and metrics support to service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2024-09-17)

* **Bug Fix**: **BREAKFIX**: Only generate AccountIDEndpointMode config for services that use it. This is a compiler break, but removes no actual functionality, as no services currently use the account ID in endpoint resolution.

# v1.4.0 (2024-09-05)

* **Feature**: Amazon CloudWatch Application Signals now supports creating Service Level Objectives using a new calculation type. Users can now create SLOs which are configured with request-based SLIs to help meet their specific business requirements.

# v1.3.3 (2024-09-04)

* No change notes available for this release.

# v1.3.2 (2024-09-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2024-08-15)

* **Dependency Update**: Bump minimum Go version to 1.21.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2024-07-25)

* **Feature**: CloudWatch Application Signals now supports application logs correlation with traces and operational health metrics of applications running on EC2 instances. Users can view the most relevant telemetry to troubleshoot application health anomalies such as spikes in latency, errors, and availability.

# v1.2.3 (2024-07-10.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2024-07-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2024-06-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2024-06-26)

* **Feature**: Support list-of-string endpoint parameter.

# v1.1.1 (2024-06-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2024-06-18)

* **Feature**: Track usage of various AWS SDK features in user-agent string.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2024-06-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2024-06-10)

* **Release**: New AWS service client module
* **Feature**: This is the initial SDK release for Amazon CloudWatch Application Signals. Amazon CloudWatch Application Signals provides curated application performance monitoring for developers to monitor and troubleshoot application health using pre-built dashboards and Service Level Objectives.

