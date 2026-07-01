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

