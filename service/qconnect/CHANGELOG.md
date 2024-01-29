# v1.3.1 (2024-01-29)

* No change notes available for this release.

# v1.3.0 (2024-01-10)

* **Feature**: QueryAssistant and GetRecommendations will be discontinued starting June 1, 2024. To receive generative responses after March 1, 2024 you will need to create a new Assistant in the Connect console and integrate the Amazon Q in Connect JavaScript library (amazon-q-connectjs) into your applications.

# v1.2.4 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.3 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.2.2 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.2.0 (2023-12-01)

* **Feature**: This release adds the PutFeedback API and allows providing feedback against the specified assistant for the specified target.
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
* **Feature**: Amazon Q in Connect, an LLM-enhanced evolution of Amazon Connect Wisdom. This release adds generative AI support to Amazon Q Connect QueryAssistant and GetRecommendations APIs.
* **Dependency Update**: Updated to the latest SDK module versions

