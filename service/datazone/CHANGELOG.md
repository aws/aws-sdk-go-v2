# v1.4.0 (2024-01-30)

* **Feature**: Add new skipDeletionCheck to DeleteDomain. Add new skipDeletionCheck to DeleteProject which also automatically deletes dependent objects

# v1.3.7 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.6 (2023-12-20)

* No change notes available for this release.

# v1.3.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.3.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.3.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.2.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2023-10-05)

* No change notes available for this release.

# v1.0.0 (2023-10-04)

* **Release**: New AWS service client module
* **Feature**: Initial release of Amazon DataZone

