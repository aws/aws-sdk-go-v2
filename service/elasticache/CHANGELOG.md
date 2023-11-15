# v1.32.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-10-24)

* **Feature**: **BREAKFIX**: Correct nullability and default value representation of various input fields across a large number of services. Calling code that references one or more of the affected fields will need to update usage accordingly. See [2162](https://github.com/aws/aws-sdk-go-v2/issues/2162).

# v1.29.5 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.4 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.3 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.1 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-08-08)

* **Feature**: Added support for cluster mode in online migration and test migration API

# v1.28.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2023-08-01)

* No change notes available for this release.

# v1.28.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.4 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.2 (2023-06-15)

* No change notes available for this release.

# v1.27.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2023-05-11)

* **Feature**: Added support to modify the cluster mode configuration for the existing ElastiCache ReplicationGroups. Customers can now modify the configuration from cluster mode disabled to cluster mode enabled.

# v1.26.10 (2023-05-04)

* No change notes available for this release.

# v1.26.9 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.8 (2023-04-10)

* No change notes available for this release.

# v1.26.7 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.6 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.26.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.2 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade smithy to 1.27.2 and correct empty query list serialization.

# v1.26.1 (2023-01-23)

* No change notes available for this release.

# v1.26.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.25.0 (2022-12-28)

* **Feature**: This release allows you to modify the encryption in transit setting, for existing Redis clusters. You can now change the TLS configuration of your Redis clusters without the need to re-build or re-provision the clusters or impact application availability.

# v1.24.3 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.2 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.1 (2022-11-22)

* No change notes available for this release.

# v1.24.0 (2022-11-16)

* **Feature**: for Redis now supports AWS Identity and Access Management authentication access to Redis clusters starting with redis-engine version 7.0

# v1.23.1 (2022-11-10)

* No change notes available for this release.

# v1.23.0 (2022-11-07)

* **Feature**: Added support for IPv6 and dual stack for Memcached and Redis clusters. Customers can now launch new Redis and Memcached clusters with IPv6 and dual stack networking support.

# v1.22.12 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.11 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.10 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.9 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.8 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.7 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.6 (2022-08-30)

* No change notes available for this release.

# v1.22.5 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-07-18)

* **Feature**: Adding AutoMinorVersionUpgrade in the DescribeReplicationGroups API

# v1.21.3 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-05-23)

* **Feature**: Added support for encryption in transit for Memcached clusters. Customers can now launch Memcached cluster with encryption in transit enabled when using Memcached version 1.6.12 or later.

# v1.20.7 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.6 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.5 (2022-04-21)

* **Documentation**: Doc only update for ElastiCache

# v1.20.4 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.3 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2022-03-23)

* **Documentation**: Doc only update for ElastiCache
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-03-14)

* **Documentation**: Doc only update for ElastiCache

# v1.20.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Documentation**: Updated service client model to latest release.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-01-14)

* **Feature**: Updated API models
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: Updated to latest service endpoints

# v1.15.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-11-30)

* **Feature**: API client updated

# v1.14.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.13.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-09-10)

* **Feature**: API client updated

# v1.10.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-08-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-07-15)

* **Feature**: The ErrorCode method on generated service error types has been corrected to match the API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-06-04)

* No change notes available for this release.

# v1.6.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

