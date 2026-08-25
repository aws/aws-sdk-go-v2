# v1.7.1 (2026-08-25)

* **Dependency Update**: Update to smithy-go v1.27.10.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2026-08-24)

* **Feature**: Added support for the GetFixture API, enabling customers to retrieve the details of a fixture from its fixture ID, and added the access role ARN to the CreateFeed, GetFeed, and UpdateFeed responses.

# v1.6.2 (2026-08-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2026-08-14)

* **Dependency Update**: Update to smithy-go v1.27.8.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2026-08-10)

* **Feature**: Added support for the SearchFixtures API and DataSourceConfiguration, enabling customers to map fixture event data onto clipping outputs for improved feature accuracy.
* **Dependency Update**: Update to smithy-go v1.27.7.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2026-08-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2026-07-31.2)

* **Feature**: AWS Elemental Inference now supports graphic composition on cropped video outputs, enabling branded graphics and other visual elements to be overlaid as part of the inference workflow.
* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.27.6 to fix various serde issues in HTTP binding services.

# v1.4.2 (2026-07-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2026-07-28)

* **Dependency Update**: Update to smithy-go v1.27.5.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2026-07-21)

* **Feature**: Add an option to clients to disable clock skew
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2026-07-13)

* No change notes available for this release.

# v1.3.0 (2026-07-06)

* **Feature**: Add request serialization snapshot tests.

# v1.2.0 (2026-07-02)

* **Feature**: Adding new BDD representation of endpoint ruleset

# v1.1.9 (2026-07-01)

* **Bug Fix**: Bump smithy-go to 1.27.3, fix JSON encorder for document.Number, endpoint host label format validation and CBOR union serialization on new serde
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.8 (2026-06-29)

* No change notes available for this release.

# v1.1.7 (2026-06-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.6 (2026-06-05.2)

* **Bug Fix**: Undo the initial wave of schema-serde releases due to several customer-reported regressions.

# v1.1.5 (2026-06-04)

* **Dependency Update**: Update to smithy-go v1.27.1 to fix several union-related deserialization bugs in schema-serde-enabled services.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2026-06-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2026-06-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2026-05-29)

* **Dependency Update**: Update to smithy-go v1.26.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2026-05-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2026-05-27)

* **Feature**: Added support for smart subtitles in Elemental Inference, enabling automatic generation of subtitles for media content. Available in English, Spanish, French, German, Italian, and Portuguese.

# v1.0.5 (2026-04-29)

* **Dependency Update**: Update to smithy-go v1.25.1.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2026-04-17)

* **Dependency Update**: Bump smithy-go to 1.25.0 to support endpointBdd trait
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2026-03-26)

* **Bug Fix**: Fix a bug where a recorded clock skew could persist on the client even if the client and server clock ended up realigning.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2026-03-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2026-03-03)

* **Dependency Update**: Bump minimum Go version to 1.24
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2026-02-24)

* **Release**: New AWS service client module
* **Feature**: Initial GA launch for AWS Elemental Inference including capabilities of Smart Crop and Live Event Clipping

