# v1.47.0 (2024-02-28)

* **Feature**: AWS WAF now supports configurable time windows for request aggregation with rate-based rules. Customers can now select time windows of 1 minute, 2 minutes or 10 minutes, in addition to the previously supported 5 minutes.

# v1.46.1 (2024-02-23)

* **Bug Fix**: Move all common, SDK-side middleware stack ops into the service client module to prevent cross-module compatibility issues in the future.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.46.0 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.45.3 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.45.2 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.45.1 (2024-02-15)

* **Bug Fix**: Correct failure to determine the error type in awsJson services that could occur when errors were modeled with a non-string `code` field.

# v1.45.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.44.0 (2024-02-06)

* **Feature**: You can now delete an API key that you've created for use with your CAPTCHA JavaScript integration API.

# v1.43.6 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.43.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.43.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.43.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.43.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.43.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.43.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.42.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.41.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.1 (2023-10-27)

* **Documentation**: Updates the descriptions for the calls that manage web ACL associations, to provide information for customer-managed IAM policies.

# v1.40.0 (2023-10-24)

* **Feature**: **BREAKFIX**: Correct nullability and default value representation of various input fields across a large number of services. Calling code that references one or more of the affected fields will need to update usage accordingly. See [2162](https://github.com/aws/aws-sdk-go-v2/issues/2162).

# v1.39.3 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.2 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.1 (2023-09-28)

* **Documentation**: Correct and improve the documentation for the FieldToMatch option JA3 fingerprint.

# v1.39.0 (2023-09-25)

* **Feature**: You can now perform an exact match against the web request's JA3 fingerprint.

# v1.38.0 (2023-09-06)

* **Feature**: The targeted protection level of the Bot Control managed rule group now provides optional, machine-learning analysis of traffic statistics to detect some bot-related activity. You can enable or disable the machine learning functionality through the API.

# v1.37.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.37.1 (2023-08-01)

* No change notes available for this release.

# v1.37.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.0 (2023-07-19)

* **Feature**: Added the URI path to the custom aggregation keys that you can specify for a rate-based rule.

# v1.35.2 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.1 (2023-06-15)

* No change notes available for this release.

# v1.35.0 (2023-06-13)

* **Feature**: You can now detect and block fraudulent account creation attempts with the new AWS WAF Fraud Control account creation fraud prevention (ACFP) managed rule group AWSManagedRulesACFPRuleSet.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.0 (2023-06-02)

* **Feature**: Added APIs to describe managed products. The APIs retrieve information about rule groups that are managed by AWS and by AWS Marketplace sellers.

# v1.33.1 (2023-06-01)

* **Documentation**: Corrected the information for the header order FieldToMatch setting

# v1.33.0 (2023-05-30)

* **Feature**: This SDK release provides customers the ability to use Header Order as a field to match.

# v1.32.0 (2023-05-16)

* **Feature**: My AWS Service (placeholder) - You can now rate limit web requests based on aggregation keys other than IP addresses, and you can aggregate using combinations of keys. You can also rate limit all requests that match a scope-down statement, without further aggregation.

# v1.31.1 (2023-05-04)

* No change notes available for this release.

# v1.31.0 (2023-04-28)

* **Feature**: You can now associate a web ACL with a Verified Access instance.

# v1.30.2 (2023-04-25)

* No change notes available for this release.

# v1.30.1 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.30.0 (2023-04-20)

* **Feature**: You can now create encrypted API keys to use in a client application integration of the JavaScript CAPTCHA API . You can also retrieve a list of your API keys and the JavaScript application integration URL.

# v1.29.0 (2023-04-11)

* **Feature**: For web ACLs that protect CloudFront protections, the default request body inspection size is now 16 KB, and you can use the new association configuration to increase the inspection size further, up to 64 KB. Sizes over 16 KB can incur additional costs.

# v1.28.2 (2023-04-10)

* No change notes available for this release.

# v1.28.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2023-04-04)

* **Feature**: This release rolls back association config feature for webACLs that protect CloudFront protections.

# v1.27.0 (2023-04-03)

* **Feature**: For web ACLs that protect CloudFront protections, the default request body inspection size is now 16 KB, and you can use the new association configuration to increase the inspection size further, up to 64 KB. Sizes over 16 KB can incur additional costs.

# v1.26.3 (2023-03-30)

* No change notes available for this release.

# v1.26.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2023-02-23)

* **Feature**: You can now associate an AWS WAF v2 web ACL with an AWS App Runner service.

# v1.25.3 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.25.2 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2023-02-16)

* **Documentation**: Added a notice for account takeover prevention (ATP). The interface incorrectly lets you to configure ATP response inspection in regional web ACLs in Region US East (N. Virginia), without returning an error. ATP response inspection is only available in web ACLs that protect CloudFront distributions.

# v1.25.0 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Feature**: For protected CloudFront distributions, you can now use the AWS WAF Fraud Control account takeover prevention (ATP) managed rule group to block new login attempts from clients that have recently submitted too many failed login attempts.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.24.3 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.2 (2023-01-18)

* **Documentation**: Improved the visibility of the guidance for updating AWS WAF resources, such as web ACLs and rule groups.

# v1.24.1 (2023-01-12)

* No change notes available for this release.

# v1.24.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.23.4 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.3 (2022-12-12)

* **Documentation**: Documents the naming requirement for logging destinations that you use with web ACLs.

# v1.23.2 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-11-07)

* **Documentation**: The geo match statement now adds labels for country and region. You can match requests at the region level by combining a geo match statement with label match statements.

# v1.23.0 (2022-10-27)

* **Feature**: This release adds the following: Challenge rule action, to silently verify client browsers; rule group rule action override to any valid rule action, not just Count; token sharing between protected applications for challenge/CAPTCHA token; targeted rules option for Bot Control managed rule group.

# v1.22.11 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.10 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.9 (2022-09-23)

* **Documentation**: Add the default specification for ResourceType in ListResourcesForWebACL.

# v1.22.8 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.7 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.6 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.5 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.4 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.3 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-08-03)

* **Feature**: You can now associate an AWS WAF web ACL with an Amazon Cognito user pool.

# v1.21.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-07-15)

* **Feature**: This SDK release provide customers ability to add sensitivity level for WAF SQLI Match Statements.

# v1.20.5 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.4 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.3 (2022-06-16)

* No change notes available for this release.

# v1.20.2 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-04-29)

* **Feature**: You can now inspect all request headers and all cookies. You can now specify how to handle oversize body contents in your rules that inspect the body.

# v1.19.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-04-08)

* **Feature**: Add a new CurrentDefaultVersion field to ListAvailableManagedRuleGroupVersions API response; add a new VersioningSupported boolean to each ManagedRuleGroup returned from ListAvailableManagedRuleGroups API response.

# v1.18.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2021-11-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.
* **Feature**: Updated service to latest API model.

# v1.12.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-09-24)

* **Feature**: API client updated

# v1.9.0 (2021-09-17)

* **Feature**: Updated API client and endpoints to latest revision.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-08-27)

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

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

