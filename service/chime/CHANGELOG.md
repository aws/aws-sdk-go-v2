# v1.41.5 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.41.4 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.41.3 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.41.2 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.41.1 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.41.0 (2025-10-30)

* **Feature**: Update endpoint ruleset parameters casing
* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.7 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.6 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.5 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.4 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.3 (2025-09-10)

* No change notes available for this release.

# v1.40.2 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.1 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.0 (2025-08-27)

* **Feature**: Remove incorrect endpoint tests
* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.39.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.5 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.4 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.3 (2025-06-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.2 (2025-04-03)

* No change notes available for this release.

# v1.36.1 (2025-03-04.2)

* **Bug Fix**: Add assurance test for operation order.

# v1.36.0 (2025-02-27)

* **Feature**: Track credential providers via User-Agent Feature ids
* **Dependency Update**: Updated to the latest SDK module versions

