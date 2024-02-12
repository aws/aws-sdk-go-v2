# v1.10.7 (2024-01-10)

* No change notes available for this release.

# v1.10.6 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.10.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.10.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.9.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2023-10-27)

* No change notes available for this release.

# v1.7.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2023-08-28)

* **Feature**: Remove incorrect regex enforcement on pagination tokens.

# v1.6.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2023-08-01)

* No change notes available for this release.

# v1.6.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-07-20.2)

* **Feature**: Adding support for Tags on Create and Resource Tagging API.

# v1.4.4 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2023-06-15)

* No change notes available for this release.

# v1.4.2 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2023-06-01)

* No change notes available for this release.

# v1.4.0 (2023-05-30)

* **Feature**: Log sources are now versioned. AWS log sources and custom sources will now come with a version identifier that enables producers to vend multiple schema versions to subscribers. Security Lake API have been refactored to more closely align with AWS API conventions.

# v1.3.6 (2023-05-04)

* No change notes available for this release.

# v1.3.5 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.4 (2023-04-10)

* No change notes available for this release.

# v1.3.3 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.2 (2023-03-22)

* No change notes available for this release.

# v1.3.1 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-03-15)

* **Feature**: Make Create/Get/ListSubscribers APIs return resource share ARN and name so they can be used to validate the RAM resource share to accept. GetDatalake can be used to track status of UpdateDatalake and DeleteDatalake requests.

# v1.2.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.2.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.2.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.1.0 (2023-01-03)

* **Feature**: Allow CreateSubscriber API to take string input that allows setting more descriptive SubscriberDescription field. Make souceTypes field required in model level for UpdateSubscriberRequest as it is required for every API call on the backend. Allow ListSubscribers take any String as nextToken param.

# v1.0.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2022-11-29.2)

* **Release**: New AWS service client module
* **Feature**: Amazon Security Lake automatically centralizes security data from cloud, on-premises, and custom sources into a purpose-built data lake stored in your account. Security Lake makes it easier to analyze security data, so you can improve the protection of your workloads, applications, and data

