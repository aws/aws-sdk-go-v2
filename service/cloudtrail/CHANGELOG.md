# v1.38.0 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.37.3 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.2 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.37.1 (2024-02-15)

* **Bug Fix**: Correct failure to determine the error type in awsJson services that could occur when errors were modeled with a non-string `code` field.

# v1.37.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.0 (2024-01-18)

* **Feature**: This release adds a new API ListInsightsMetricData to retrieve metric data from CloudTrail Insights.

# v1.35.7 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.6 (2023-12-20)

* No change notes available for this release.

# v1.35.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.35.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.35.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.2 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.1 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.34.0 (2023-11-27)

* **Feature**: CloudTrail Lake now supports federating event data stores. giving users the ability to run queries against their event data using Amazon Athena.

# v1.33.1 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.0 (2023-11-15)

* **Feature**: The Lake Repricing feature lets customers configure a BillingMode for an event data store. The BillingMode determines the cost for ingesting and storing events and the default and maximum retention period for the event data store.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.0 (2023-11-09.2)

* **Feature**: The Insights in Lake feature lets customers enable CloudTrail Insights on a source CloudTrail Lake event data store and create a destination event data store to collect Insights events based on unusual management event activity in the source event data store.

# v1.31.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-08-25)

* **Feature**: Add ThrottlingException with error code 429 to handle CloudTrail Delegated Admin request rate exceeded on organization resources.

# v1.28.6 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.5 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.4 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.3 (2023-08-10)

* **Documentation**: Documentation updates for CloudTrail.

# v1.28.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2023-08-01)

* No change notes available for this release.

# v1.28.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.3 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.2 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.1 (2023-06-15)

* No change notes available for this release.

# v1.27.0 (2023-06-13)

* **Feature**: This feature allows users to view dashboards for CloudTrail Lake event data stores.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2023-06-02)

* **Feature**: This feature allows users to start and stop event ingestion on a CloudTrail Lake event data store.

# v1.25.0 (2023-05-18)

* **Feature**: Add ConflictException to PutEventSelectors, add (Channel/EDS)ARNInvalidException to Tag APIs. These exceptions provide customers with more specific error messages instead of internal errors.

# v1.24.8 (2023-05-04)

* No change notes available for this release.

# v1.24.7 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.6 (2023-04-10)

* No change notes available for this release.

# v1.24.5 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.4 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.3 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.2 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.24.1 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Feature**: This release adds an InsufficientEncryptionPolicyException type to the StartImport endpoint
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.23.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2023-01-31)

* **Feature**: Add new "Channel" APIs to enable users to manage channels used for CloudTrail Lake integrations, and "Resource Policy" APIs to enable users to manage the resource-based permissions policy attached to a channel.

# v1.22.1 (2023-01-23)

* No change notes available for this release.

# v1.22.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.21.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-12-13)

* **Feature**: Merging mainline branch for service model into mainline release branch. There are no new APIs.

# v1.20.4 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.3 (2022-11-22)

* No change notes available for this release.

# v1.20.2 (2022-11-16)

* No change notes available for this release.

# v1.20.1 (2022-11-10)

* No change notes available for this release.

# v1.20.0 (2022-11-07)

* **Feature**: This release includes support for configuring a delegated administrator to manage an AWS Organizations organization CloudTrail trails and event data stores, and AWS Key Management Service encryption of CloudTrail Lake event data stores.

# v1.19.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-10-19)

* **Feature**: This release includes support for exporting CloudTrail Lake query results to an Amazon S3 bucket.

# v1.18.2 (2022-10-07)

* No change notes available for this release.

# v1.18.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-09-19)

* **Feature**: This release includes support for importing existing trails into CloudTrail Lake.

# v1.17.0 (2022-09-14)

* **Feature**: This release adds CloudTrail getChannel and listChannels APIs to allow customer to view the ServiceLinkedChannel configurations.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.12 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.11 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.10 (2022-08-30)

* No change notes available for this release.

# v1.16.9 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.8 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.7 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.6 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.5 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.4 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.3 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.2 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-04-27)

* **Feature**: Increases the retention period maximum to 2557 days. Deprecates unused fields of the ListEventDataStores API response. Updates documentation.

# v1.15.6 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.5 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.4 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.3 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.2 (2022-03-14)

* No change notes available for this release.

# v1.15.1 (2022-03-10)

* No change notes available for this release.

# v1.15.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service client model to latest release.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2022-01-28)

* **Documentation**: Updated to latest API model.

# v1.13.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2022-01-07)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.10.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-11-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.

# v1.8.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-09-02)

* **Feature**: API client updated

# v1.5.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.2 (2021-06-11)

* **Documentation**: Updated to latest API model.

# v1.3.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

