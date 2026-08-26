# v1.13.0 (2026-08-26)

* **Feature**: Stop registering the `ComputeContentLength` middleware in generated clients. `Content-Length` is now set when the request body is set via `SetStream`.
* **Dependency Update**: Update to smithy-go v1.28.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.1 (2026-08-25)

* **Dependency Update**: Update to smithy-go v1.27.10.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2026-08-24)

* **Feature**: Adding private and self-signed certificate configuration support for penetration tests

# v1.11.2 (2026-08-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2026-08-14)

* **Dependency Update**: Update to smithy-go v1.27.8.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2026-08-13)

* **Feature**: Add support for setting a maximum task-hour budget cap on penetration tests and code reviews, and for revalidating previously reported findings via a new REVALIDATION job type.

# v1.10.1 (2026-08-10)

* **Dependency Update**: Update to smithy-go v1.27.7.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2026-08-07)

* **Feature**: Added enableEmailMfa input field on Actor to enable email-based MFA during penetration tests. When enabled, a server-generated mfaForwardingAddress is returned. Set up a forwarding rule in your email provider to forward MFA emails to this address so the agent can complete email-based MFA login flows

# v1.9.2 (2026-08-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2026-07-31.2)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.27.6 to fix various serde issues in HTTP binding services.

# v1.9.0 (2026-07-30)

* **Feature**: Adds support for providing a branch override when configured integrated repositories

# v1.8.2 (2026-07-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2026-07-28)

* **Dependency Update**: Update to smithy-go v1.27.5.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2026-07-27)

* **Feature**: AWS Security Agent adds a new task hours field that reflects the active work done for a task.

# v1.7.0 (2026-07-21)

* **Feature**: Add an option to clients to disable clock skew
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2026-07-13)

* No change notes available for this release.

# v1.6.0 (2026-07-06)

* **Feature**: Add request serialization snapshot tests.

# v1.5.2 (2026-07-01)

* **Bug Fix**: Bump smithy-go to 1.27.3, fix JSON encorder for document.Number, endpoint host label format validation and CBOR union serialization on new serde
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2026-06-29)

* No change notes available for this release.

# v1.5.0 (2026-06-17)

* **Feature**: Updated AWS Security Agent SDK model with new APIs for threat modeling, code review, security requirements, and additional integration providers.

# v1.4.6 (2026-06-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.5 (2026-06-05.2)

* **Bug Fix**: Undo the initial wave of schema-serde releases due to several customer-reported regressions.

# v1.4.4 (2026-06-04)

* **Dependency Update**: Update to smithy-go v1.27.1 to fix several union-related deserialization bugs in schema-serde-enabled services.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2026-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2026-06-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2026-05-29)

* **Dependency Update**: Update to smithy-go v1.26.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2026-05-28)

* **Feature**: Adding new BDD representation of endpoint ruleset
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2026-05-22)

* **Feature**: Adds support for verification scripts on penetration test findings. Customers can now download executable scripts to independently reproduce confirmed vulnerabilities, with instructions and required environment variables provided for each finding.

# v1.2.0 (2026-05-13)

* **Feature**: Add support for code reviews, a new resource type that enables automated security-focused static analysis of source code repositories.

# v1.1.0 (2026-05-04)

* **Feature**: AWS Security Agent is adding a new target domain verification method for private VPC penetration testing. Additionally, the target domain resource will now have a verification status reason field to surface additional details about domain verification

# v1.0.2 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2026-03-31)

* **Release**: New AWS service client module
* **Feature**: AWS Security Agent is a service that proactively secures applications throughout the development lifecycle with automated security reviews and on-demand penetration testing.

