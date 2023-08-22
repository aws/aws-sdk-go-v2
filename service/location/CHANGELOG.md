# v1.26.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2023-08-01)

* No change notes available for this release.

# v1.26.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.2 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2023-07-06)

* **Feature**: This release adds support for authenticating with Amazon Location Service's Places & Routes APIs with an API Key. Also, with this release developers can publish tracked device position updates to Amazon EventBridge.

# v1.24.0 (2023-06-15)

* **Feature**: Amazon Location Service adds categories to places, including filtering on those categories in searches. Also, you can now add metadata properties to your geofences.

# v1.23.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2023-05-30)

* **Feature**: This release adds API support for political views for the maps service APIs: CreateMap, UpdateMap, DescribeMap.

# v1.22.7 (2023-05-04)

* No change notes available for this release.

# v1.22.6 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.5 (2023-04-10)

* No change notes available for this release.

# v1.22.4 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.3 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2023-03-07)

* **Documentation**: Documentation update for the release of 3 additional map styles for use with Open Data Maps: Open Data Standard Dark, Open Data Visualization Light & Open Data Visualization Dark.

# v1.22.0 (2023-02-23)

* **Feature**: This release adds support for using Maps APIs with an API Key in addition to AWS Cognito. This includes support for adding, listing, updating and deleting API Keys.

# v1.21.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.21.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.21.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2023-01-10)

* **Feature**: This release adds support for two new route travel models, Bicycle and Motorcycle which can be used with Grab data source.

# v1.20.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.19.5 (2022-12-15)

* **Documentation**: This release adds support for a new style, "VectorOpenDataStandardLight" which can be used with the new data source, "Open Data Maps (Preview)".
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.4 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.3 (2022-10-25)

* **Documentation**: Added new map styles with satellite imagery for map resources using HERE as a data provider.

# v1.19.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-09-27)

* **Feature**: This release adds place IDs, which are unique identifiers of places, along with a new GetPlace operation, which can be used with place IDs to find a place again later. UnitNumber and UnitType are also added as new properties of places.

# v1.18.6 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.5 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.4 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.3 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-08-09)

* **Feature**: Amazon Location Service now allows circular geofences in BatchPutGeofence, PutGeofence, and GetGeofence  APIs.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.6 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.5 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.4 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.3 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.2 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-05-06)

* **Feature**: Amazon Location Service now includes a MaxResults parameter for ListGeofences requests.

# v1.16.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-03-22)

* **Feature**: Amazon Location Service now includes a MaxResults parameter for GetDevicePositionHistory requests.

# v1.15.1 (2022-03-15)

* **Documentation**: New HERE style "VectorHereExplore" and "VectorHereExploreTruck".

# v1.15.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2022-01-28)

* **Feature**: Updated to latest API model.
* **Bug Fix**: Updates SDK API client deserialization to pre-allocate byte slice and string response payloads, [#1565](https://github.com/aws/aws-sdk-go-v2/pull/1565). Thanks to [Tyson Mote](https://github.com/tysonmote) for submitting this PR.

# v1.12.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: API client updated

# v1.9.2 (2021-12-03)

* **Bug Fix**: Fixed a bug that prevented the resolution of the correct endpoint for some API operations.
* **Bug Fix**: Fixed an issue that caused some operations to not be signed using sigv4, resulting in authentication failures.

# v1.9.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-11-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-10-11)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Bug Fix**: Fixed an issue that caused one or more API operations to fail when attempting to resolve the service endpoint. ([#1349](https://github.com/aws/aws-sdk-go-v2/pull/1349))
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2021-06-04)

* **Feature**: Updated service client to latest API model.

# v1.1.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

