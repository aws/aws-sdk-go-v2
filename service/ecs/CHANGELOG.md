# v1.33.1 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.0 (2023-11-13)

* **Feature**: Adds a Client Token parameter to the ECS RunTask API. The Client Token parameter allows for idempotent RunTask requests.

# v1.32.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.4 (2023-10-17)

* **Documentation**: Documentation only updates to address Amazon ECS tickets.

# v1.30.3 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.2 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.1 (2023-09-05)

* **Documentation**: Documentation only update for Amazon ECS.

# v1.30.0 (2023-08-31)

* **Feature**: This release adds support for an account-level setting that you can use to configure the number of days for AWS Fargate task retirement.

# v1.29.6 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.5 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.4 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.3 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2023-08-04)

* **Documentation**: This is a documentation update to address various tickets.

# v1.29.1 (2023-08-01)

* No change notes available for this release.

# v1.29.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.2 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2023-06-30)

* **Feature**: Added new field  "credentialspecs" to the ecs task definition to support gMSA of windows/linux in both domainless and domain-joined mode

# v1.27.4 (2023-06-19)

* **Documentation**: Documentation only update to address various tickets.

# v1.27.3 (2023-06-15)

* No change notes available for this release.

# v1.27.2 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.1 (2023-05-18)

* **Documentation**: Documentation only release to address various tickets.

# v1.27.0 (2023-05-04)

* **Feature**: Documentation update for new error type NamespaceNotFoundException for CreateCluster and UpdateCluster

# v1.26.3 (2023-05-02)

* **Documentation**: Documentation only update to address Amazon ECS tickets.

# v1.26.2 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2023-04-21)

* **Documentation**: Documentation update to address various Amazon ECS tickets.

# v1.26.0 (2023-04-19)

* **Feature**: This release supports the Account Setting "TagResourceAuthorization" that allows for enhanced Tagging security controls.

# v1.25.1 (2023-04-14)

* **Documentation**: This release supports  ephemeral storage for AWS Fargate Windows containers.

# v1.25.0 (2023-04-10)

* **Feature**: This release adds support for enabling FIPS compliance on Amazon ECS Fargate tasks

# v1.24.4 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.3 (2023-04-05)

* **Documentation**: This is a document only updated to add information about Amazon Elastic Inference (EI).

# v1.24.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2023-02-23)

* **Feature**: This release supports deleting Amazon ECS task definitions that are in the INACTIVE state.

# v1.23.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.23.4 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.3 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.23.2 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2023-01-23)

* No change notes available for this release.

# v1.23.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.22.0 (2022-12-19)

* **Feature**: This release adds support for alarm-based rollbacks in ECS, a new feature that allows customers to add automated safeguards for Amazon ECS service rolling updates.

# v1.21.0 (2022-12-15)

* **Feature**: This release adds support for container port ranges in ECS, a new capability that allows customers to provide container port ranges to simplify use cases where multiple ports are in use in a container. This release updates TaskDefinition mutation APIs and the Task description APIs.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-12-02)

* **Documentation**: Documentation updates for Amazon ECS
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-11-28)

* **Feature**: This release adds support for ECS Service Connect, a new capability that simplifies writing and operating resilient distributed applications. This release updates the TaskDefinition, Cluster, Service mutation APIs with Service connect constructs and also adds a new ListServicesByNamespace API.

# v1.19.2 (2022-11-22)

* No change notes available for this release.

# v1.19.1 (2022-11-16)

* No change notes available for this release.

# v1.19.0 (2022-11-10)

* **Feature**: This release adds support for task scale-in protection with updateTaskProtection and getTaskProtection APIs. UpdateTaskProtection API can be used to protect a service managed task from being terminated by scale-in events and getTaskProtection API to get the scale-in protection status of a task.

# v1.18.26 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.25 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.24 (2022-10-13)

* **Documentation**: Documentation update to address tickets.

# v1.18.23 (2022-10-04)

* **Documentation**: Documentation updates to address various Amazon ECS tickets.

# v1.18.22 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.21 (2022-09-16)

* **Documentation**: This release supports new task definition sizes.

# v1.18.20 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.19 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.18 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.17 (2022-08-30)

* No change notes available for this release.

# v1.18.16 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.15 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.14 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.13 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.12 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.11 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.10 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.9 (2022-06-21)

* **Documentation**: Amazon ECS UpdateService now supports the following parameters: PlacementStrategies, PlacementConstraints and CapacityProviderStrategy.

# v1.18.8 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.7 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.6 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.5 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.4 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-03-22)

* **Documentation**: Documentation only update to address tickets

# v1.18.1 (2022-03-15)

* **Documentation**: Documentation only update to address tickets

# v1.18.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service client model to latest release.
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
* **Documentation**: API client updated
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
* **Feature**: Updated service to latest API model.
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.11.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-12)

* **Feature**: API client updated

# v1.7.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-06-04)

* **Feature**: Updated service client to latest API model.

# v1.4.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Feature**: Updated to latest service API model.
* **Dependency Update**: Updated to the latest SDK module versions

