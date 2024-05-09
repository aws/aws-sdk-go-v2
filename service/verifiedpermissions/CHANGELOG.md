# v1.14.0 (2024-05-09)

* **Feature**: Adds policy effect and actions fields to Policy API's.

# v1.13.2 (2024-05-08)

* **Bug Fix**: GoDoc improvement

# v1.13.1 (2024-04-11)

* No change notes available for this release.

# v1.13.0 (2024-04-05)

* **Feature**: Adding BatchIsAuthorizedWithToken API which supports multiple authorization requests against a PolicyStore given a bearer token.

# v1.12.0 (2024-04-04)

* **Feature**: Adds GroupConfiguration field to Identity Source API's

# v1.11.4 (2024-04-03)

* No change notes available for this release.

# v1.11.3 (2024-03-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.2 (2024-03-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2024-03-07)

* **Bug Fix**: Remove dependency on go-cmp.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2024-03-06)

* **Feature**: Deprecating details in favor of configuration for GetIdentitySource and ListIdentitySources APIs.

# v1.10.1 (2024-02-23)

* **Bug Fix**: Move all common, SDK-side middleware stack ops into the service client module to prevent cross-module compatibility issues in the future.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.9.3 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.9.1 (2024-02-15)

* **Bug Fix**: Correct failure to determine the error type in awsJson services that could occur when errors were modeled with a non-string `code` field.

# v1.9.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.4 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.3 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.8.2 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.8.0 (2023-12-01)

* **Feature**: Adds description field to PolicyStore API's and namespaces field to GetSchema.
* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.6.1 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2023-11-17)

* **Feature**: Adding BatchIsAuthorized API which supports multiple authorization requests against a PolicyStore

# v1.5.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-10-19)

* **Feature**: Improving Amazon Verified Permissions Create experience

# v1.2.3 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-08-24)

* **Documentation**: Documentation updates for Amazon Verified Permissions.

# v1.2.0 (2023-08-22)

* **Feature**: Documentation updates for Amazon Verified Permissions. Increases max results per page for ListPolicyStores, ListPolicies, and ListPolicyTemplates APIs from 20 to 50.

# v1.1.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2023-08-01)

* No change notes available for this release.

# v1.1.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2023-06-30)

* **Documentation**: This release corrects several broken links in the documentation.

# v1.0.3 (2023-06-27)

* **Documentation**: This update fixes several broken links to the Cedar documentation.

# v1.0.2 (2023-06-23)

* **Documentation**: Added improved descriptions and new code samples to SDK documentation.

# v1.0.1 (2023-06-15)

* No change notes available for this release.

# v1.0.0 (2023-06-13)

* **Release**: New AWS service client module
* **Feature**: GA release of Amazon Verified Permissions.
* **Dependency Update**: Updated to the latest SDK module versions

