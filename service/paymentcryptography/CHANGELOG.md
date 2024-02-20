# v1.8.2 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.8.1 (2024-02-15)

* **Bug Fix**: Correct failure to determine the error type in awsJson services that could occur when errors were modeled with a non-string `code` field.

# v1.8.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2024-01-16)

* **Feature**: Provide an additional option for key exchange using RSA wrap/unwrap in addition to tr-34/tr-31 in ImportKey and ExportKey operations. Added new key usage (type) TR31_M1_ISO_9797_1_MAC_KEY, for use with Generate/VerifyMac dataplane operations  with ISO9797 Algorithm 1 MAC calculations.

# v1.6.3 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.6.1 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2023-12-06)

* **Feature**: AWS Payment Cryptography IPEK feature release
* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.5.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.4.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-09-18)

* **Announcement**: [BREAKFIX] Change in MaxResults datatype from value to pointer type in cognito-sync service.
* **Feature**: Adds several endpoint ruleset changes across all models: smaller rulesets, removed non-unique regional endpoints, fixes FIPS and DualStack endpoints, and make region not required in SDK::Endpoint. Additional breakfix to cognito-sync field.

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

# v1.0.4 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2023-06-15)

* No change notes available for this release.

# v1.0.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2023-06-08)

* **Release**: New AWS service client module
* **Feature**: Initial release of AWS Payment Cryptography Control Plane service for creating and managing cryptographic keys used during card payment processing.

