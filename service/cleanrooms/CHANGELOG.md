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

