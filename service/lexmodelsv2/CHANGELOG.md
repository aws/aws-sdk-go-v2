# v1.32.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.32.1 (2023-08-01)

* No change notes available for this release.

# v1.32.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2023-07-20.2)

* **Feature**: This release updates type for Channel field in SessionSpecification and UtteranceSpecification

# v1.30.0 (2023-07-18)

* **Feature**: This release adds support for Lex Developers to view analytics for their bots.

# v1.29.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.2 (2023-06-15)

* No change notes available for this release.

# v1.29.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2023-06-06)

* **Feature**: This release adds support for Lex Developers to create test sets and to execute those test-sets against their bots.

# v1.28.9 (2023-05-04)

* No change notes available for this release.

# v1.28.8 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.7 (2023-04-10)

* No change notes available for this release.

# v1.28.6 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.5 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.4 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.3 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.28.2 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.28.0 (2023-02-09)

* **Feature**: AWS Lex now supports Network of Bots.

# v1.27.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.26.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.26.0 (2022-11-07)

* **Feature**: Amazon Lex now supports new APIs for viewing and editing Custom Vocabulary in bots.

# v1.25.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2022-09-23)

* **Feature**: This release introduces additional optional parameters promptAttemptsSpecification to PromptSpecification, which enables the users to configure interrupt setting and Audio, DTMF and Text input configuration for the initial and retry prompt played by the Bot

# v1.24.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-09-14)

* **Feature**: This release is for supporting Composite Slot Type feature in AWS Lex V2. Composite Slot Type will help developer to logically group coherent slots and maintain their inter-relationships in runtime conversation.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.3 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-08-22)

* **Feature**: This release introduces a new feature to stop a running BotRecommendation Job for Automated Chatbot Designer.

# v1.22.0 (2022-08-17)

* **Feature**: This release introduces support for enhanced conversation design with the ability to define custom conversation flows with conditional branching and new bot responses.

# v1.21.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-07-05)

* **Feature**: This release introduces additional optional parameters "messageSelectionStrategy" to PromptSpecification, which enables the users to configure the bot to play messages in orderly manner.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.7 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.6 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-03-10)

* **Feature**: This release makes slotTypeId an optional parameter in CreateSlot and UpdateSlot APIs in Amazon Lex V2 for model building. Customers can create and update slots without specifying a slot type id.

# v1.19.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-01-14)

* **Feature**: Updated API models
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: API client updated

# v1.14.0 (2021-12-02)

* **Feature**: API client updated
* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-11-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-11-12)

* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.11.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-10-11)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-09-24)

* **Feature**: API client updated

# v1.7.0 (2021-09-17)

* **Feature**: Updated API client and endpoints to latest revision.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2021-08-12)

* **Feature**: API client updated

# v1.4.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2021-07-15)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

