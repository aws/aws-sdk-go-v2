# v1.8.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2023-08-01)

* No change notes available for this release.

# v1.8.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2023-07-25)

* **Feature**: Added support for Auto-Assocate Billing Groups for CreateBillingGroup, UpdateBillingGroup, and ListBillingGroups.

# v1.6.8 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.7 (2023-06-15)

* No change notes available for this release.

# v1.6.6 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.5 (2023-05-04)

* No change notes available for this release.

# v1.6.4 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2023-04-10)

* No change notes available for this release.

# v1.6.2 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2023-03-17)

* **Feature**: This release adds a new filter to ListAccountAssociations API and a new filter to ListBillingGroups API.

# v1.5.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.5.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.5.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-01-17)

* **Feature**: This release adds support for SKU Scope for pricing plans.

# v1.4.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.3.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2022-12-06)

* **Feature**: This release adds the Tiering Pricing Rule feature.

# v1.2.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2022-11-16)

* **Feature**: This release adds a new feature BillingEntity pricing rule.

# v1.1.0 (2022-11-08)

* **Feature**: This release adds the Recurring Custom Line Item feature along with a new API ListCustomLineItemVersions.

# v1.0.19 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.18 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.17 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.16 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.15 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.14 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.13 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.12 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.11 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.10 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.9 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.8 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.7 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2022-03-16)

* **Release**: New AWS service client module
* **Feature**: This is the initial SDK release for AWS Billing Conductor. The AWS Billing Conductor is a customizable billing service, allowing you to customize your billing data to match your desired business structure.

