# v1.21.3 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.1 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2023-08-10)

* **Feature**: This release enables configuring security groups for Network Load Balancers

# v1.20.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2023-08-01)

* No change notes available for this release.

# v1.20.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.15 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.14 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.13 (2023-06-15)

* No change notes available for this release.

# v1.19.12 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.11 (2023-05-04)

* No change notes available for this release.

# v1.19.10 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.9 (2023-04-10)

* No change notes available for this release.

# v1.19.8 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.7 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.6 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.19.4 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.3 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade smithy to 1.27.2 and correct empty query list serialization.

# v1.19.2 (2023-02-02)

* **Documentation**: The GWLB Flex Health Check project updates the default values of healthy-threshold-count from 3 to 5 and unhealthy-threshold-count from 3 to 2

# v1.19.1 (2023-01-23)

* No change notes available for this release.

# v1.19.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.18.28 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.27 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.26 (2022-11-22)

* No change notes available for this release.

# v1.18.25 (2022-11-17)

* **Documentation**: Provides new target group attributes to turn on/off cross zone load balancing and configure target group health for Network Load Balancers and Application Load Balancers. Provides improvements to health check configuration for Network Load Balancers.

# v1.18.24 (2022-11-16)

* No change notes available for this release.

# v1.18.23 (2022-11-10)

* No change notes available for this release.

# v1.18.22 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.21 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.20 (2022-10-07)

* **Documentation**: Gateway Load Balancer adds a new feature (target_failover) for customers to rebalance existing flows to a healthy target after marked unhealthy or deregistered. This allows graceful patching/upgrades of target appliances during maintenance windows, and helps reduce unhealthy target failover time.

# v1.18.19 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.18 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.17 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.16 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.15 (2022-08-30)

* No change notes available for this release.

# v1.18.14 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.13 (2022-08-25)

* **Documentation**: Documentation updates for ELBv2.  Gateway Load Balancer now supports Configurable Flow Stickiness, enabling you to configure the hashing used to maintain stickiness of flows to a specific target appliance.

# v1.18.12 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.11 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.10 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.9 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.8 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.7 (2022-06-29)

* **Documentation**: This release adds two attributes for ALB. One, helps to preserve the host header and the other helps to modify, preserve, or remove the X-Forwarded-For header in the HTTP request.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.6 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: Updated to latest service endpoints

# v1.13.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-11-30)

* **Feature**: API client updated

# v1.12.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.11.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-10-11)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-09-30)

* **Feature**: API client updated

# v1.7.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-07-15)

* **Feature**: The ErrorCode method on generated service error types has been corrected to match the API model.
* **Documentation**: Updated service model to latest revision.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

