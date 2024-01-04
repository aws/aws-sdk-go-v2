# v1.8.6 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.8.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.8.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2023-11-29)

* **Feature**: AWS Clean Rooms now provides differential privacy to protect against user-identification attempts and machine learning modeling to allow two parties to identify similar users in their data.
* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.4 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.7.2 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2023-11-14)

* **Feature**: This feature provides the ability for the collaboration creator to configure either the member who can run queries or a different member in the collaboration to be billed for query compute costs.

# v1.6.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-08-30)

* **Feature**: This release decouples member abilities in a collaboration. With this change, the member who can run queries no longer needs to be the same as the member who can receive results.

# v1.3.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2023-08-01)

* No change notes available for this release.

# v1.3.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Feature**: This release introduces custom SQL queries - an expanded set of SQL you can run. This release adds analysis templates, a new resource for storing pre-defined custom SQL queries ahead of time. This release also adds the Custom analysis rule, which lets you approve analysis templates for querying.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-06-29)

* **Feature**: This release adds support for the OR operator in RSQL join match conditions and the ability to control which operators (AND, OR) are allowed in a join match condition.

# v1.1.6 (2023-06-15)

* No change notes available for this release.

# v1.1.5 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2023-05-04)

* No change notes available for this release.

# v1.1.3 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2023-04-10)

* No change notes available for this release.

# v1.1.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-03-21)

* **Feature**: GA Release of AWS Clean Rooms, Added Tagging Functionality
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.0.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.0.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2023-01-12)

* **Release**: New AWS service client module
* **Feature**: Initial release of AWS Clean Rooms

