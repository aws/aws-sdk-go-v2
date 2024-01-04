# v1.10.7 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.6 (2023-12-20)

* No change notes available for this release.

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

# v1.9.3 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.9.1 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2023-11-17)

* **Feature**: Adds new querying capabilities for running data queries on a monitor

# v1.8.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2023-09-13)

* **Feature**: This release updates the Amazon CloudWatch Internet Monitor API domain name.

# v1.5.5 (2023-09-05)

* No change notes available for this release.

# v1.5.4 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-08-01)

* **Feature**: This release adds a new feature for Amazon CloudWatch Internet Monitor that enables customers to set custom thresholds, for performance and availability drops, for impact limited to a single city-network to trigger creation of a health event.

# v1.4.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.2 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-06-28)

* **Feature**: This release adds a new feature for Amazon CloudWatch Internet Monitor that enables customers to set custom thresholds, for performance and availability drops, for triggering when to create a health event.

# v1.2.5 (2023-06-15)

* No change notes available for this release.

# v1.2.4 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2023-05-04)

* No change notes available for this release.

# v1.2.2 (2023-05-02)

* No change notes available for this release.

# v1.2.1 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-04-17)

* **Feature**: This release includes a new configurable value, TrafficPercentageToMonitor, which allows users to adjust the amount of traffic monitored by percentage

# v1.1.2 (2023-04-10)

* No change notes available for this release.

# v1.1.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-03-31)

* **Feature**: This release adds a new feature for Amazon CloudWatch Internet Monitor that enables customers to deliver internet measurements to Amazon S3 buckets as well as CloudWatch Logs.

# v1.0.3 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2023-03-08)

* No change notes available for this release.

# v1.0.0 (2023-02-27)

* **Release**: New AWS service client module
* **Feature**: CloudWatch Internet Monitor is a a new service within CloudWatch that will help application developers and network engineers continuously monitor internet performance metrics such as availability and performance between their AWS-hosted applications and end-users of these applications

