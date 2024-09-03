# v1.17.1 (2024-09-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2024-08-23)

* **Feature**: Releasing the support for Action User Confirmation.

# v1.16.1 (2024-08-15)

* **Dependency Update**: Bump minimum Go version to 1.21.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2024-07-10.2)

* **Feature**: Introduces new data sources and chunking strategies for Knowledge bases, advanced parsing logic using FMs, session summary generation, and code interpretation (preview) for Claude V3 Sonnet and Haiku models. Also introduces Prompt Flows (preview) to link prompts, foundational models, and resources.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2024-07-10)

* **Feature**: Introduces new data sources and chunking strategies for Knowledge bases, advanced parsing logic using FMs, session summary generation, and code interpretation (preview) for Claude V3 Sonnet and Haiku models. Also introduces Prompt Flows (preview) to link prompts, foundational models, and resources.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.1 (2024-06-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.14.0 (2024-06-26)

* **Feature**: Support list-of-string endpoint parameter.

# v1.13.1 (2024-06-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2024-06-18)

* **Feature**: Track usage of various AWS SDK features in user-agent string.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.3 (2024-06-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.2 (2024-06-07)

* **Bug Fix**: Add clock skew correction on all service clients
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.1 (2024-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2024-05-30)

* **Feature**: With this release, Knowledge bases for Bedrock adds support for Titan Text Embedding v2.

# v1.11.1 (2024-05-23)

* No change notes available for this release.

# v1.11.0 (2024-05-20)

* **Feature**: This release adds support for using Guardrails with Bedrock Agents.

# v1.10.3 (2024-05-16)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.2 (2024-05-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2024-05-08)

* **Bug Fix**: GoDoc improvement

# v1.10.0 (2024-05-03)

* **Feature**: This release adds support for using Provisioned Throughput with Bedrock Agents.

# v1.9.0 (2024-05-01)

* **Feature**: This release adds support for using MongoDB Atlas as a vector store when creating a knowledge base.

# v1.8.0 (2024-04-23)

* **Feature**: Introducing the ability to create multiple data sources per knowledge base, specify S3 buckets as data sources from external accounts, and exposing levers to define the deletion behavior of the underlying vector store data.

# v1.7.0 (2024-04-22)

* **Feature**: Releasing the support for simplified configuration and return of control

# v1.6.0 (2024-04-16)

* **Feature**: For Create Agent API, the agentResourceRoleArn parameter is no longer required.

# v1.5.1 (2024-03-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2024-03-27)

* **Feature**: This changes introduces metadata documents statistics and also updates the documentation for bedrock agent.

# v1.4.3 (2024-03-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2024-03-07)

* **Bug Fix**: Remove dependency on go-cmp.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2024-02-23)

* **Bug Fix**: Move all common, SDK-side middleware stack ops into the service client module to prevent cross-module compatibility issues in the future.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2024-02-22)

* **Feature**: Add middleware stack snapshot tests.

# v1.3.2 (2024-02-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.3.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-12-21)

* **Feature**: This release introduces Amazon Aurora as a vector store on Knowledge Bases for Amazon Bedrock

# v1.1.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.1.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.1.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2023-11-28.2)

* **Release**: New AWS service client module
* **Feature**: This release introduces Agents for Amazon Bedrock
* **Dependency Update**: Updated to the latest SDK module versions

