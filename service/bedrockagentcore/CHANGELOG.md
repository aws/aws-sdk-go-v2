# v1.32.1 (2026-06-29)

* No change notes available for this release.

# v1.32.0 (2026-06-22)

* **Feature**: Adds an optional extractionMode field to CreateEvent. SKIP retains the event in short-term memory but excludes it from long-term memory extraction.

# v1.31.0 (2026-06-17)

* **Feature**: AgentCore Harness service will be Generally Available at NYS 2026 with this Treb release. Harness will support invoking specific endpoints via the qualifier parameter, AWS Skills for pre-built agent capabilities, and improved validation for skill git source URLs.

# v1.30.0 (2026-06-12)

* **Feature**: Added tagging and CMK support across optimization, an explanation field in recommendation output, and an insights feature to identify failure patterns, extract user intents, and summarize execution behavior

# v1.29.0 (2026-06-09)

* **Feature**: Add RetryableConflictException (HTTP 409) to InvokeAgentRuntimeCommand and GetAgentCard to prevent orphaned VMs during concurrent session access. The SDK automatically retries this exception with backoff. Enforcement is not yet active and will be enabled in a future service update.

# v1.28.6 (2026-06-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.5 (2026-06-05.2)

* **Bug Fix**: Undo the initial wave of schema-serde releases due to several customer-reported regressions.

# v1.28.4 (2026-06-04)

* **Dependency Update**: Update to smithy-go v1.27.1 to fix several union-related deserialization bugs in schema-serde-enabled services.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.3 (2026-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.2 (2026-06-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.1 (2026-05-29)

* **Dependency Update**: Update to smithy-go v1.26.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.28.0 (2026-05-28)

* **Feature**: Added Harness support for LiteLLM model configuration for third-party model providers. Added S3 and Git skill source types. Added Responses API format for OpenAI and Bedrock models. Added runtimeUserId and runtimeClientError to InvokeHarness.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2026-05-26)

* **Feature**: Adding new BDD representation of endpoint ruleset

# v1.26.0 (2026-05-19)

* **Feature**: Add RetryableConflictException (HTTP 409) to InvokeAgentRuntime and StopRuntimeSession to prevent orphaned VMs during concurrent session access. The SDK automatically retries this exception with backoff. Enforcement is not yet active and will be enabled in a future service update.

# v1.25.0 (2026-05-07)

* **Feature**: Launching AgentCore payments - a capability that provides secure, instant microtransaction payments for AI agents to access paid APIs, MCP servers, and content. It handles payment processing for x402 protocol, payment limits, and 3P wallet integrations with Coinbase CDP and Stripe (Privy).

# v1.24.0 (2026-04-30)

* **Feature**: AgentCore Identity now supports on-behalf-of token exchange OAuth2. AgentCore Memory now supports metadata for LongTerm Memory Records.

# v1.23.0 (2026-04-29)

* **Feature**: Adds batch evaluation for running evaluators against multiple agent sessions with server-side orchestration, AI-powered recommendations for optimizing system prompts and tool descriptions, and AB testing with controlled traffic splitting and statistical significance reporting
* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2026-04-22)

* **Feature**: Adds support for Amazon Bedrock AgentCore Harness data plane APIs, enabling customers to invoke managed agent loops and execute commands on live agent sessions with streaming responses.

# v1.21.1 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2026-04-16)

* **Feature**: Introducing NamespacePath in AgentCore Memory to support hierarchical prefix based memory record retrieval.

# v1.20.0 (2026-04-09)

* **Feature**: Introducing support for SearchRegistryRecords API on AgentCoreRegistry

# v1.19.0 (2026-04-07)

* **Feature**: This release includes support for 1) InvokeBrowser API, enabling OS-level control of AgentCore Browser Tool sessions through mouse actions, keyboard input, and screenshots. 2) Added documentation noting that empty sessions are automatically deleted after one day in the ListSessions API.

# v1.18.0 (2026-04-01)

* **Feature**: Added the ability to filter out empty sessions when listing sessions. Customers can now retrieve only sessions that still contain events, eliminating the need to check each session individually. No changes required for existing integrations.

# v1.17.0 (2026-03-30)

* **Feature**: Adds Ground Truth support for AgentCore Evaluations (Evaluate)

# v1.16.0 (2026-03-27)

* **Feature**: Adding AgentCore Code Interpreter Node.js Runtime Support with an optional runtime field

# v1.15.2 (2026-03-26)

* **Bug Fix**: Fix a bug where a recorded clock skew could persist on the client even if the client and server clock ended up realigning.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.1 (2026-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2026-03-19)

* **Feature**: This release includes SDK support for the following new features on AgentCore Built In Tools.  1. Enterprise Policies for AgentCore Browser Tool. 2. Root CA Configuration Support for AgentCore Browser Tool and Code Interpreter. 3. API changes to AgentCore Browser Profile APIs

# v1.14.0 (2026-03-16)

* **Feature**: Provide support to perform deterministic operations on agent runtime through shell command executions via the new InvokeAgentRuntimeCommand API

# v1.13.3 (2026-03-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.2 (2026-03-03)

* **Dependency Update**: Bump minimum Go version to 1.24
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2026-02-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2026-02-10.2)

* **Feature**: Added AgentCore browser proxy configuration support, allowing routing of browser traffic through HTTP and HTTPS proxy servers with authentication and bypass rules.

# v1.12.0 (2026-02-05)

* **Feature**: Support Browser profile persistence (cookies and local storage) across sessions for AgentCore Browser.

# v1.11.0 (2026-01-21)

* **Feature**: Supports custom browser extensions for AgentCore Browser and increased message payloads up to 100KB per message in an Event for AgentCore Memory

# v1.10.2 (2026-01-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2025-12-02)

* **Feature**: Support for AgentCore Evaluations and Episodic memory strategy for AgentCore Memory.
* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.9.1 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.9.0 (2025-11-20)

* **Feature**: Bedrock AgentCore Memory release for redriving memory extraction jobs (StartMemoryExtractionJob and ListMemoryExtractionJob)

# v1.8.5 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.4 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.8.3 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.8.2 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.8.1 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2025-10-23)

* **Feature**: Fixing the service documentation name
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2025-10-13)

* **Feature**: Updated InvokeAgentRuntime API to accept account id optionally and added CompleteResourceTokenAuth API.

# v1.6.0 (2025-10-10)

* **Feature**: Bedrock AgentCore release for Runtime, and Memory.

# v1.5.0 (2025-10-06)

* **Feature**: Add support for batch memory management, agent card retrieval and session termination

# v1.4.5 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2025-09-10)

* No change notes available for this release.

# v1.4.2 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2025-08-27)

* **Feature**: Remove incorrect endpoint tests
* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.3.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2025-08-04)

* **Feature**: Remove superfluous field from API
* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2025-07-16.2)

* **Release**: New AWS service client module
* **Feature**: Initial release of Amazon Bedrock AgentCore SDK including Runtime, Built-In Tools, Memory, Gateway and Identity.

