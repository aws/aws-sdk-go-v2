# v1.2.0 (2026-08-27)

* **Feature**: Support connection read timeouts in the SDK. This is currently available on an opt-in basis by setting env `AWS_ENABLE_DEFAULT_SOCKET_TIMEOUT_2026=true`.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2026-08-26)

* **Feature**: Stop registering the `ComputeContentLength` middleware in generated clients. `Content-Length` is now set when the request body is set via `SetStream`.
* **Dependency Update**: Update to smithy-go v1.28.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2026-08-25)

* **Dependency Update**: Update to smithy-go v1.27.10.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2026-08-20)

* **Documentation**: Documentation update for the CreateSubscription API to correct the default value of the approval mode parameter. The default value for paid subscriptions is MANUAL, not IMMEDIATE as previously documented. The default value remains IMMEDIATE for FREE tier subscriptions.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2026-08-14)

* **Dependency Update**: Update to smithy-go v1.27.8.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2026-08-10)

* **Dependency Update**: Update to smithy-go v1.27.7.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2026-08-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2026-07-31.2)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.27.6 to fix various serde issues in HTTP binding services.

# v1.0.0 (2026-07-30)

* **Release**: New AWS service client module
* **Feature**: Adds support for Public PricingPlanManager SDK

