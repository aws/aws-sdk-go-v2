# v1.4.1 (2026-07-01)

* **Bug Fix**: Bump smithy-go to 1.27.3, fix JSON encorder for document.Number, endpoint host label format validation and CBOR union serialization on new serde
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2026-06-29)

* **Feature**: Expand input validation to support Unicode characters and markdown table syntax.

# v1.3.0 (2026-06-10)

* **Feature**: Add support for MedicalScribeBinaryAudioEvent in the Medical Scribe streaming input. This new event type lets you send audio as a raw binary payload instead of a base64-encoded value

# v1.2.7 (2026-06-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.6 (2026-06-05.2)

* **Bug Fix**: Undo the initial wave of schema-serde releases due to several customer-reported regressions.

# v1.2.5 (2026-06-04)

* **Dependency Update**: Update to smithy-go v1.27.1 to fix several union-related deserialization bugs in schema-serde-enabled services.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2026-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2026-06-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2026-05-29)

* **Dependency Update**: Update to smithy-go v1.26.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2026-05-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2026-05-26)

* **Feature**: Adding new BDD representation of endpoint ruleset

# v1.1.1 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2026-04-24)

* **Feature**: Corrected CreateWebAppConfiguration documentation. Adding slash as an allowed character for the Ambient documentation agent to allow pronoun specifications.

# v1.0.4 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2026-03-26)

* **Bug Fix**: Fix a bug where a recorded clock skew could persist on the client even if the client and server clock ended up realigning.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2026-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2026-03-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2026-03-05)

* **Release**: New AWS service client module
* **Feature**: Connect-Health SDK is AWS's unified SDK for the Amazon Connect Health offering. It allows healthcare developers to integrate purpose-built agents - such as patient insights, ambient documentation, and medical coding - into their existing applications, including EHRs, telehealth, and revenue cycle.

