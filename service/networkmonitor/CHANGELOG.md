# v1.2.1 (2024-02-23)

* **Bug Fix**: Move all common, SDK-side middleware stack ops into the service client module to prevent cross-module compatibility issues in the future.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.1.2 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.1.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2023-12-22)

* **Release**: New AWS service client module
* **Feature**: CloudWatch Network Monitor is a new service within CloudWatch that will help network administrators and operators continuously monitor network performance metrics such as round-trip-time and packet loss between their AWS-hosted applications and their on-premises locations.

