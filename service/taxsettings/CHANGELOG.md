# v1.4.3 (2024-09-27)

* No change notes available for this release.

# v1.4.2 (2024-09-25)

* No change notes available for this release.

# v1.4.1 (2024-09-23)

* No change notes available for this release.

# v1.4.0 (2024-09-20)

* **Feature**: Add tracing and metrics support to service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.4 (2024-09-17)

* **Bug Fix**: **BREAKFIX**: Only generate AccountIDEndpointMode config for services that use it. This is a compiler break, but removes no actual functionality, as no services currently use the account ID in endpoint resolution.

# v1.3.3 (2024-09-04)

* No change notes available for this release.

# v1.3.2 (2024-09-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2024-08-15)

* **Dependency Update**: Bump minimum Go version to 1.21.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2024-07-18)

* **Feature**: Set default endpoint for aws partition. Requests from all regions in aws partition will be forward to us-east-1 endpoint.

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

# v1.0.2 (2024-06-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2024-06-07)

* **Bug Fix**: Add clock skew correction on all service clients
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2024-06-04)

* **Release**: New AWS service client module
* **Feature**: Initial release of AWS Tax Settings API

