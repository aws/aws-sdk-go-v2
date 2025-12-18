# v1.16.0 (2025-12-18)

* **Feature**: Feature to support header exchanges between Bedrock AgentCore Gateway Targets and client, along with propagating query parameter to the configured targets.

# v1.15.2 (2025-12-15)

* **Documentation**: This release updates broken links for AgentCore Policy APIs in the AWS CLI and SDK resources.

# v1.15.1 (2025-12-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2025-12-02)

* **Feature**: Supports AgentCore Evaluations, Policy, Episodic Memory Strategy, Resource Based Policy for Runtime and Gateway APIs, API Gateway Rest API Targets and enhances JWT authorizer.
* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.14.1 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.14.0 (2025-11-21)

* **Feature**: Support for agentcore gateway interceptor configurations and NONE authorizer type

# v1.13.4 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.3 (2025-11-12)

* **Bug Fix**: Further reduce allocation overhead when the metrics system isn't in-use.
* **Bug Fix**: Reduce allocation overhead when the client doesn't have any HTTP interceptors configured.
* **Bug Fix**: Remove blank trace spans towards the beginning of the request that added no additional information. This conveys a slight reduction in overall allocations.

# v1.13.2 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.13.1 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.13.0 (2025-11-03)

* **Feature**: Adds support for direct code deploy with CreateAgentRuntime and UpdateAgentRuntime

# v1.12.0 (2025-10-30)

* **Feature**: Web-Bot-Auth support for AgentCore Browser tool to help reduce captcha challenges.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2025-10-23)

* **Feature**: Fixing the service documentation name
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2025-10-13)

* **Feature**: Updated http status code in control plane apis of agentcore runtime, tools and identity. Additional included provider types for AgentCore Identity

# v1.9.0 (2025-10-10)

* **Feature**: Bedrock AgentCore release for Gateway, and Memory including Self-Managed Strategies support for Memory.

# v1.8.0 (2025-10-08)

* **Feature**: Adding support for authorizer type AWS_IAM to AgentCore Control Gateway.

# v1.7.0 (2025-10-06)

* **Feature**: Add support for VM lifecycle configuration parameters and A2A protocol

# v1.6.0 (2025-09-30)

* **Feature**: Tagging support for AgentCore Gateway

# v1.5.2 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2025-09-19)

* **Feature**: Add tagging and VPC support to AgentCore Runtime, Code Interpreter, and Browser resources. Add support for configuring request headers in Runtime. Fix AgentCore Runtime shape names.

# v1.4.4 (2025-09-10)

* No change notes available for this release.

# v1.4.3 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2025-08-26)

* **Feature**: Remove incorrect endpoint tests

# v1.3.2 (2025-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.3.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2025-08-04)

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

