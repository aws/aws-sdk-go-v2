# v1.42.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.42.1 (2023-08-01)

* No change notes available for this release.

# v1.42.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.41.2 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.41.1 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.41.0 (2023-06-22)

* **Feature**: Introducing Amazon Kendra Retrieve API that can be used to retrieve relevant passages or text excerpts given an input query.

# v1.40.4 (2023-06-15)

* No change notes available for this release.

# v1.40.3 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.40.2 (2023-06-06)

* No change notes available for this release.

# v1.40.1 (2023-05-04)

* No change notes available for this release.

# v1.40.0 (2023-05-02)

* **Feature**: AWS Kendra now supports configuring document fields/attributes via the GetQuerySuggestions API. You can now base query suggestions on the contents of document fields.

# v1.39.3 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.2 (2023-04-10)

* No change notes available for this release.

# v1.39.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.39.0 (2023-03-30)

* **Feature**: AWS Kendra now supports featured results for a query.

# v1.38.7 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.6 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.5 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.38.4 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.3 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.38.2 (2023-02-08)

* No change notes available for this release.

# v1.38.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.38.0 (2023-01-11)

* **Feature**: This release adds support to new document types - RTF, XML, XSLT, MS_EXCEL, CSV, JSON, MD

# v1.37.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.36.3 (2022-12-30)

* No change notes available for this release.

# v1.36.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.36.0 (2022-11-28)

* **Feature**: Amazon Kendra now supports preview of table information from HTML tables in the search results. The most relevant cells with their corresponding rows, columns are displayed as a preview in the search result. The most relevant table cell or cells are also highlighted in table preview.

# v1.35.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.35.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

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

