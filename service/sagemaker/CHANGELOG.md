# v1.71.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.71.0 (2023-03-09)

* **Feature**: Amazon SageMaker Inference now allows SSM access to customer's model container by setting the "EnableSSMAccess" parameter for a ProductionVariant in CreateEndpointConfig API.

# v1.70.0 (2023-03-08)

* **Feature**: There needs to be a user identity to specify the SageMaker user who perform each action regarding the entity. However, these is a not a unified concept of user identity across SageMaker service that could be used today.

# v1.69.0 (2023-03-02)

* **Feature**: Add a new field "EndpointMetrics" in SageMaker Inference Recommender "ListInferenceRecommendationsJobSteps" API response.

# v1.68.3 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.68.2 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.68.1 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.68.0 (2023-02-10)

* **Feature**: Amazon SageMaker Autopilot adds support for selecting algorithms in CreateAutoMLJob API.

# v1.67.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.67.0 (2023-01-31)

* **Feature**: Amazon SageMaker Automatic Model Tuning now supports more completion criteria for Hyperparameter Optimization.

# v1.66.0 (2023-01-27)

* **Feature**: This release supports running SageMaker Training jobs with container images that are in a private Docker registry.

# v1.65.0 (2023-01-25)

* **Feature**: SageMaker Inference Recommender now decouples from Model Registry and could accept Model Name to invoke inference recommendations job; Inference Recommender now provides CPU/Memory Utilization metrics data in recommendation output.

# v1.64.0 (2023-01-23)

* **Feature**: Amazon SageMaker Inference now supports P4de instance types.

# v1.63.0 (2023-01-19)

* **Feature**: HyperParameterTuningJobs now allow passing environment variables into the corresponding TrainingJobs

# v1.62.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.61.0 (2022-12-21)

* **Feature**: This release enables adding RStudio Workbench support to an existing Amazon SageMaker Studio domain. It allows setting your RStudio on SageMaker environment configuration parameters and also updating the RStudioConnectUrl and RStudioPackageManagerUrl parameters for existing domains

# v1.60.0 (2022-12-20)

* **Feature**: Amazon SageMaker Autopilot adds support for new objective metrics in CreateAutoMLJob API.

# v1.59.0 (2022-12-19)

* **Feature**: AWS Sagemaker - Sagemaker Images now supports Aliases as secondary identifiers for ImageVersions. SageMaker Images now supports additional metadata for ImageVersions for better images management.

# v1.58.0 (2022-12-16)

* **Feature**: AWS sagemaker - Features: This release adds support for random seed, it's an integer value used to initialize a pseudo-random number generator. Setting a random seed will allow the hyperparameter tuning search strategies to produce more consistent configurations for the same tuning job.

# v1.57.0 (2022-12-15)

* **Feature**: SageMaker Inference Recommender now allows customers to load tests their models on various instance types using private VPC.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.56.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.56.0 (2022-11-30)

* **Feature**: Added Models as part of the Search API. Added Model shadow deployments in realtime inference, and shadow testing in managed inference. Added support for shared spaces, geospatial APIs, Model Cards, AutoMLJobStep in pipelines, Git repositories on user profiles and domains, Model sharing in Jumpstart.

# v1.55.0 (2022-11-18)

* **Feature**: Added DisableProfiler flag as a new field in ProfilerConfig

# v1.54.0 (2022-11-03)

* **Feature**: Amazon SageMaker now supports running training jobs on ml.trn1 instance types.

# v1.53.0 (2022-11-02)

* **Feature**: This release updates Framework model regex for ModelPackage to support new Framework version xgboost, sklearn.

# v1.52.0 (2022-10-27)

* **Feature**: This change allows customers to provide a custom entrypoint script for the docker container to be run while executing training jobs, and provide custom arguments to the entrypoint script.

# v1.51.0 (2022-10-26)

* **Feature**: Amazon SageMaker Automatic Model Tuning now supports specifying Grid Search strategy for tuning jobs, which evaluates all hyperparameter combinations exhaustively based on the categorical hyperparameters provided.

# v1.50.0 (2022-10-24)

* **Feature**: SageMaker Inference Recommender now supports a new API ListInferenceRecommendationJobSteps to return the details of all the benchmark we create for an inference recommendation job.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.49.0 (2022-10-21)

* **Feature**: CreateInferenceRecommenderjob API now supports passing endpoint details directly, that will help customers to identify the max invocation and max latency they can achieve for their model and the associated endpoint along with getting recommendations on other instances.
* **Dependency Update**: Updated to the latest SDK module versions

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

