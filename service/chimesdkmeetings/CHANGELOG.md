# v1.14.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.13.12 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.11 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.10 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.9 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.8 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.7 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.6 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.5 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.4 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.3 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.2 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2022-08-04)

* **Feature**: Adds support for Tags on Amazon Chime SDK WebRTC sessions

# v1.12.1 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2022-07-07)

* **Feature**: Adds support for AppKeys and TenantIds in Amazon Chime SDK WebRTC sessions

# v1.11.2 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.1 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2022-06-09)

* **Feature**: Adds support for live transcription in AWS GovCloud (US) Regions.

# v1.10.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2022-06-01)

* **Feature**: Adds support for centrally controlling each participant's ability to send and receive audio, video and screen share within a WebRTC session.  Attendee capabilities can be specified when the attendee is created and updated during the session with the new BatchUpdateAttendeeCapabilitiesExcept API.

# v1.9.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2022-04-22)

* **Feature**: Include additional exceptions types.

# v1.8.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2022-03-21)

* **Feature**: Add support for media replication to link multiple WebRTC media sessions together to reach larger and global audiences. Participants connected to a replica session can be granted access to join the primary session and can switch sessions with their existing WebRTC connection

# v1.7.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service client model to latest release.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.
* **Feature**: Updated to latest service endpoints

# v1.2.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2021-11-30)

* **Feature**: API client updated

# v1.1.0 (2021-11-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2021-11-12)

* **Release**: New AWS service client module

