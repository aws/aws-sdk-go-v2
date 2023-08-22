# v1.30.6 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.5 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.4 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.3 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.2 (2023-08-03)

* **Documentation**: Documentation changes related to Amazon EC2 Auto Scaling APIs.

# v1.30.1 (2023-08-01)

* No change notes available for this release.

# v1.30.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Feature**: You can now configure an instance refresh to set its status to 'failed' when it detects that a specified CloudWatch alarm has gone into the ALARM state. You can also choose to roll back the instance refresh automatically when the alarm threshold is met.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-07-27)

* **Feature**: This release updates validation for instance types used in the AllowedInstanceTypes and ExcludedInstanceTypes parameters of the InstanceRequirements property of a MixedInstancesPolicy.

# v1.28.10 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.9 (2023-06-15)

* No change notes available for this release.

# v1.28.8 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.7 (2023-05-04)

* No change notes available for this release.

# v1.28.6 (2023-04-28)

* No change notes available for this release.

# v1.28.5 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.4 (2023-04-21)

* No change notes available for this release.

# v1.28.3 (2023-04-10)

* No change notes available for this release.

# v1.28.2 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2023-04-04)

* **Documentation**: Documentation updates for Amazon EC2 Auto Scaling

# v1.28.0 (2023-03-30)

* **Feature**: Amazon EC2 Auto Scaling now supports Elastic Load Balancing traffic sources with the AttachTrafficSources, DetachTrafficSources, and DescribeTrafficSources APIs. This release also introduces a new activity status, "WaitingForConnectionDraining", for VPC Lattice to the DescribeScalingActivities API.

# v1.27.4 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.3 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.2 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.27.1 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2023-02-10)

* **Feature**: You can now either terminate/replace, ignore, or wait for EC2 Auto Scaling instances on standby or protected from scale in. Also, you can also roll back changes from a failed instance refresh.

# v1.26.2 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade smithy to 1.27.2 and correct empty query list serialization.

# v1.26.1 (2023-01-23)

* No change notes available for this release.

# v1.26.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.25.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2022-12-08)

* **Feature**: Adds support for metric math for target tracking scaling policies, saving you the cost and effort of publishing a custom metric to CloudWatch. Also adds support for VPC Lattice by adding the Attach/Detach/DescribeTrafficSources APIs and a new health check type to the CreateAutoScalingGroup API.

# v1.24.4 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.3 (2022-11-22)

* No change notes available for this release.

# v1.24.2 (2022-11-16)

* No change notes available for this release.

# v1.24.1 (2022-11-10)

* **Documentation**: This release adds a new price capacity optimized allocation strategy for Spot Instances to help customers optimize provisioning of Spot Instances via EC2 Auto Scaling, EC2 Fleet, and Spot Fleet. It allocates Spot Instances based on both spare capacity availability and Spot Instance price.

# v1.24.0 (2022-11-07)

* **Feature**: This release adds support for two new attributes for attribute-based instance type selection - NetworkBandwidthGbps and AllowedInstanceTypes.

# v1.23.18 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.17 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.16 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.15 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.14 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.13 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.12 (2022-08-30)

* No change notes available for this release.

# v1.23.11 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.10 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.9 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.8 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.7 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.6 (2022-07-25)

* **Documentation**: Documentation update for Amazon EC2 Auto Scaling.

# v1.23.5 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.4 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.3 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-04-19)

* **Feature**: EC2 Auto Scaling now adds default instance warm-up times for all scaling activities, health check replacements, and other replacement events in the Auto Scaling instance lifecycle.

# v1.22.4 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.3 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2022-03-08.2)

* No change notes available for this release.

# v1.22.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-02-24.2)

* **Feature**: API client updated

# v1.20.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: Updated to latest service endpoints

# v1.16.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2021-11-30)

* **Feature**: API client updated

# v1.15.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.14.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-08-12)

* **Feature**: API client updated

# v1.10.1 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-07-15)

* **Feature**: The ErrorCode method on generated service error types has been corrected to match the API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-07-01)

* **Feature**: API client updated

# v1.8.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-06-11)

* **Feature**: Updated to latest API model.

# v1.6.0 (2021-06-04)

* **Feature**: Updated service client to latest API model.

# v1.5.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

