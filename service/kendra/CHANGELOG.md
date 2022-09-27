# v1.35.0 (2022-09-27)

* **Feature**: My AWS Service (placeholder) - Amazon Kendra now provides a data source connector for DropBox. For more information, see https://docs.aws.amazon.com/kendra/latest/dg/data-source-dropbox.html

# v1.34.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.34.0 (2022-09-14)

* **Feature**: This release enables our customer to choose the option of Sharepoint 2019 for the on-premise Sharepoint connector.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.3 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.2 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.1 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.33.0 (2022-08-19)

* **Feature**: This release adds support for a new authentication type - Personal Access Token (PAT) for confluence server.

# v1.32.0 (2022-08-17)

* **Feature**: This release adds Zendesk connector (which allows you to specify Zendesk SAAS platform as data source), Proxy Support for Sharepoint and Confluence Server (which allows you to specify the proxy configuration if proxy is required to connect to your Sharepoint/Confluence Server as data source).

# v1.31.4 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.3 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.2 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.31.0 (2022-07-21)

* **Feature**: Amazon Kendra now provides Oauth2 support for SharePoint Online. For more information, see https://docs.aws.amazon.com/kendra/latest/dg/data-source-sharepoint.html

# v1.30.0 (2022-07-14)

* **Feature**: This release adds AccessControlConfigurations which allow you to redefine your document level access control without the need for content re-indexing.

# v1.29.1 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.29.0 (2022-06-30)

* **Feature**: Amazon Kendra now provides a data source connector for alfresco

# v1.28.2 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2022-06-02)

* **Feature**: Amazon Kendra now provides a data source connector for GitHub. For more information, see https://docs.aws.amazon.com/kendra/latest/dg/data-source-github.html

# v1.27.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2022-05-12)

* **Feature**: Amazon Kendra now provides a data source connector for Jira. For more information, see https://docs.aws.amazon.com/kendra/latest/dg/data-source-jira.html

# v1.26.0 (2022-05-05)

* **Feature**: AWS Kendra now supports hierarchical facets for a query. For more information, see https://docs.aws.amazon.com/kendra/latest/dg/filtering.html

# v1.25.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2022-04-19)

* **Feature**: Amazon Kendra now provides a data source connector for Quip. For more information, see https://docs.aws.amazon.com/kendra/latest/dg/data-source-quip.html

# v1.24.0 (2022-04-06)

* **Feature**: Amazon Kendra now provides a data source connector for Box. For more information, see https://docs.aws.amazon.com/kendra/latest/dg/data-source-box.html

# v1.23.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.23.0 (2022-03-14)

* **Feature**: Amazon Kendra now provides a data source connector for Slack. For more information, see https://docs.aws.amazon.com/kendra/latest/dg/data-source-slack.html

# v1.22.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service client model to latest release.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.0 (2022-01-14)

* **Feature**: Updated API models
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.17.0 (2021-12-02)

* **Feature**: API client updated
* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2021-11-19)

* **Announcement**: Fix API modeling bug incorrectly generating `DocumentAttributeValue` type as a union instead of a structure. This update corrects this bug by correcting the `DocumentAttributeValue` type to be a `struct` instead of an `interface`. This change also removes the `DocumentAttributeValueMember` types. To migrate to this change your application using service/kendra will need to be updated to use struct members in `DocumentAttributeValue` instead of `DocumentAttributeValueMember` types.
* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-11-12)

* **Feature**: Service clients now support custom endpoints that have an initial URI path defined.

# v1.14.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-10-21)

* **Feature**: API client updated
* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2021-10-11)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-08-04)

* **Feature**: Updated to latest API model.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-07-01)

* **Feature**: API client updated

# v1.7.0 (2021-06-25)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-06-11)

* **Feature**: Updated to latest API model.

# v1.5.0 (2021-06-04)

* **Feature**: Updated service client to latest API model.

# v1.4.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

