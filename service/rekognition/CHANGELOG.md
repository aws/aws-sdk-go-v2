# v1.30.6 (2023-08-22)

* No change notes available for this release.

# v1.30.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.2 (2023-08-07)

* **Documentation**: This release adds code snippets for Amazon Rekognition Custom Labels.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.1 (2023-08-01)

* No change notes available for this release.

# v1.30.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.4 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2023-06-15)

* No change notes available for this release.

# v1.29.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-06-12)

* **Feature**: This release adds support for improved accuracy with user vector in Amazon Rekognition Face Search. Adds new APIs: AssociateFaces, CreateUser, DeleteUser, DisassociateFaces, ListUsers, SearchUsers, SearchUsersByImage. Also adds new face metadata that can be stored: user vector.

# v1.28.0 (2023-05-15)

* **Feature**: This release adds a new EyeDirection attribute in Amazon Rekognition DetectFaces and IndexFaces APIs which predicts the yaw and pitch angles of a person's eye gaze direction for each face detected in the image.

# v1.27.0 (2023-05-04)

* **Feature**: This release adds a new attribute FaceOccluded. Additionally, you can now select attributes individually (e.g. ["DEFAULT", "FACE_OCCLUDED", "AGE_RANGE"] instead of ["ALL"]), which can reduce response time.

# v1.26.0 (2023-04-28)

* **Feature**: Added support for aggregating moderation labels by video segment timestamps for Stored Video Content Moderation APIs and added additional information about the job to all Stored Video Get API responses.

# v1.25.0 (2023-04-24)

* **Feature**: Added new status result to Liveness session status.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2023-04-10)

* **Feature**: This release adds support for Face Liveness APIs in Amazon Rekognition. Updates UpdateStreamProcessor to return ResourceInUseException Exception. Minor updates to API documentation.

# v1.23.7 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.6 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.23.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.23.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.22.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-12-12)

* **Feature**: Adds support for "aliases" and "categories", inclusion and exclusion filters for labels and label categories, and aggregating labels by video segment timestamps for Stored Video Label Detection APIs.

# v1.21.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-11-11)

* **Feature**: Adding support for ImageProperties feature to detect dominant colors and image brightness, sharpness, and contrast, inclusion and exclusion filters for labels and label categories, new fields to the API response, "aliases" and "categories"

# v1.20.7 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.6 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.5 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.4 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.3 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-08-16)

* **Feature**: This release adds APIs which support copying an Amazon Rekognition Custom Labels model and managing project policies across AWS account.

# v1.19.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-07-26)

* **Feature**: This release introduces support for the automatic scaling of inference units used by Amazon Rekognition Custom Labels models.

# v1.18.5 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.4 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-05-16)

* **Documentation**: Documentation updates for Amazon Rekognition.

# v1.18.0 (2022-04-27)

* **Feature**: This release adds support to configure stream-processor resources for label detections on streaming-videos. UpateStreamProcessor API is also launched with this release, which could be used to update an existing stream-processor.

# v1.17.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2022-01-07)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.12.0 (2021-12-03)

* **Feature**: API client updated

# v1.11.2 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.
* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.10.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-08-12)

* **Feature**: API client updated

# v1.6.2 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

