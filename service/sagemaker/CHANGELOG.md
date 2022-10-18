# v1.48.0 (2022-10-18)

* **Feature**: This change allows customers to enable data capturing while running a batch transform job, and configure monitoring schedule to monitoring the captured data.

# v1.47.0 (2022-10-17)

* **Feature**: This release adds support for C7g, C6g, C6gd, C6gn, M6g, M6gd, R6g, and R6gn Graviton instance types in Amazon SageMaker Inference.

# v1.46.0 (2022-09-30)

* **Feature**: A new parameter called ExplainerConfig is added to CreateEndpointConfig API to enable SageMaker Clarify online explainability feature.

# v1.45.0 (2022-09-29)

* **Feature**: SageMaker Training Managed Warm Pools let you retain provisioned infrastructure to reduce latency for repetitive training workloads.

# v1.44.0 (2022-09-21)

* **Feature**: SageMaker now allows customization on Canvas Application settings, including enabling/disabling time-series forecasting and specifying an Amazon Forecast execution role at both the Domain and UserProfile levels.

# v1.43.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.43.0 (2022-09-15)

* **Feature**: Amazon SageMaker Automatic Model Tuning now supports specifying Hyperband strategy for tuning jobs, which uses a multi-fidelity based tuning strategy to stop underperforming hyperparameter configurations early.

# v1.42.0 (2022-09-14)

* **Feature**: Fixed a bug in the API client generation which caused some operation parameters to be incorrectly generated as value types instead of pointer types. The service API always required these affected parameters to be nilable. This fixes the SDK client to match the expectations of the the service API.
* **Feature**: SageMaker Hosting now allows customization on ML instance storage volume size, model data download timeout and inference container startup ping health check timeout for each ProductionVariant in CreateEndpointConfig API.
* **Feature**: This release adds HyperParameterTuningJob type in Search API.
* **Feature**: This release adds Mode to AutoMLJobConfig.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.41.0 (2022-09-02)

* **Feature**: This release enables administrators to attribute user activity and API calls from Studio notebooks, Data Wrangler and Canvas to specific users even when users share the same execution IAM role.  ExecutionRoleIdentityConfig at Sagemaker domain level enables this feature.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.0 (2022-08-31)

* **Feature**: SageMaker Inference Recommender now accepts Inference Recommender fields: Domain, Task, Framework, SamplePayloadUrl, SupportedContentTypes, SupportedInstanceTypes, directly in our CreateInferenceRecommendationsJob API through ContainerConfig
* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.3 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.2 (2022-08-22)

* No change notes available for this release.

# v1.39.1 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.0 (2022-08-09)

* **Feature**: Amazon SageMaker Automatic Model Tuning now supports specifying multiple alternate EC2 instance types to make tuning jobs more robust when the preferred instance type is not available due to insufficient capacity.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.0 (2022-07-19)

* **Feature**: Fixed an issue with cross account QueryLineage

# v1.37.0 (2022-07-18)

* **Feature**: Amazon SageMaker Edge Manager provides lightweight model deployment feature to deploy machine learning models on requested devices.

# v1.36.0 (2022-07-14)

* **Feature**: This release adds support for G5, P4d, and C6i instance types in Amazon SageMaker Inference and increases the number of hyperparameters that can be searched from 20 to 30 in Amazon SageMaker Automatic Model Tuning

# v1.35.0 (2022-07-07)

* **Feature**: Heterogeneous clusters: the ability to launch training jobs with multiple instance types. This enables running component of the training job on the instance type that is most suitable for it. e.g. doing data processing and augmentation on CPU instances and neural network training on GPU instances

# v1.34.1 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.0 (2022-06-29)

* **Feature**: This release adds: UpdateFeatureGroup, UpdateFeatureMetadata, DescribeFeatureMetadata APIs; FeatureMetadata type in Search API; LastModifiedTime, LastUpdateStatus, OnlineStoreTotalSizeBytes in DescribeFeatureGroup API.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.0 (2022-06-23)

* **Feature**: SageMaker Ground Truth now supports Virtual Private Cloud. Customers can launch labeling jobs and access to their private workforce in VPC mode.

# v1.32.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.0 (2022-05-27)

* **Feature**: Amazon SageMaker Notebook Instances now allows configuration of Instance Metadata Service version and Amazon SageMaker Studio now supports G5 instance types.

# v1.31.0 (2022-05-25)

* **Feature**: Amazon SageMaker Autopilot adds support for manually selecting features from the input dataset using the CreateAutoMLJob API.

# v1.30.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2022-05-03)

* **Feature**: SageMaker Autopilot adds new metrics for all candidate models generated by Autopilot experiments; RStudio on SageMaker now allows users to bring your own development environment in a custom image.

# v1.29.0 (2022-04-27)

* **Feature**: Amazon SageMaker Autopilot adds support for custom validation dataset and validation ratio through the CreateAutoMLJob and DescribeAutoMLJob APIs.

# v1.28.0 (2022-04-26)

* **Feature**: SageMaker Inference Recommender now accepts customer KMS key ID for encryption of endpoints and compilation outputs created during inference recommendation.

# v1.27.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2022-04-07)

* **Feature**: Amazon Sagemaker Notebook Instances now supports G5 instance types

# v1.26.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-01-28)

* **Feature**: Updated to latest API model.

# v1.23.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-01-07)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: API client updated

# v1.20.0 (2021-12-02)

* **Feature**: API client updated
* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.
* **Feature**: Updated service to latest API model.
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.18.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2021-10-11)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-09-17)

* **Feature**: Updated API client and endpoints to latest revision.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-09-10)

* **Feature**: API client updated

# v1.13.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-08-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-07-01)

* **Feature**: API client updated

# v1.8.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-06-11)

* **Feature**: Updated to latest API model.

# v1.6.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

