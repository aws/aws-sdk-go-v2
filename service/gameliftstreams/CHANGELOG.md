# v1.8.3 (2025-11-11)

* **Bug Fix**: Return validation error if input region is not a valid host label.

# v1.8.2 (2025-11-04)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.23.2 which should convey some passive reduction of overall allocations, especially when not using the metrics system.

# v1.8.1 (2025-10-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2025-10-28)

* **Feature**: Add stream group expiration date and expired status

# v1.7.0 (2025-10-24)

* **Feature**: Add status reasons for TERMINATED stream sessions

# v1.6.9 (2025-10-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.8 (2025-10-17)

* **Documentation**: Updates documentation to clarify valid application binaries for an Amazon GameLift Streams application and provide descriptions of stream session error status reasons

# v1.6.7 (2025-10-16)

* **Dependency Update**: Bump minimum Go version to 1.23.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.6 (2025-09-26)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.5 (2025-09-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.4 (2025-09-10)

* No change notes available for this release.

# v1.6.3 (2025-09-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2025-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2025-08-27)

* **Dependency Update**: Update to smithy-go v1.23.0.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2025-08-21)

* **Feature**: The default application in a stream group can now be changed at any time using UpdateStreamGroup to update the DefaultApplicationIdentifier.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2025-08-20)

* **Bug Fix**: Remove unused deserialization code.

# v1.5.0 (2025-08-11)

* **Feature**: Add support for configuring per-service Options via callback on global config.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2025-08-08)

* **Feature**: Removed incorrect endpoint tests

# v1.3.1 (2025-08-07)

* **Documentation**: Adds Proton 9.0-2 to the list of runtime environment options available when creating an Amazon GameLift Streams application

# v1.3.0 (2025-08-04)

* **Feature**: Support configurable auth scheme preferences in service clients via AWS_AUTH_SCHEME_PREFERENCE in the environment, auth_scheme_preference in the config file, and through in-code settings on LoadDefaultConfig and client constructor methods.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2025-07-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2025-07-28)

* **Feature**: Add support for HTTP interceptors.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.5 (2025-07-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2025-06-17)

* **Dependency Update**: Update to smithy-go v1.22.4.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2025-06-10)

* **Documentation**: Documentation updates for Amazon GameLift Streams to address formatting errors, correct resource ID examples, and update links to other guides
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2025-06-06)

* No change notes available for this release.

# v1.1.1 (2025-04-03)

* No change notes available for this release.

# v1.1.0 (2025-03-25)

* **Feature**: Minor updates to improve developer experience.

# v1.0.1 (2025-03-06)

* No change notes available for this release.

# v1.0.0 (2025-03-05)

* **Release**: New AWS service client module
* **Feature**: New Service: Amazon GameLift Streams delivers low-latency game streaming from AWS global infrastructure to virtually any device with a browser at up to 1080p resolution and 60 fps.

