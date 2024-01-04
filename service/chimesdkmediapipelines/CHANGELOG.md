# v1.13.6 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.5 (2023-12-08)

* **Bug Fix**: Reinstate presence of default Retryer in functional options, but still respect max attempts set therein.

# v1.13.4 (2023-12-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.3 (2023-12-06)

* **Bug Fix**: Restore pre-refactor auth behavior where all operations could technically be performed anonymously.

# v1.13.2 (2023-12-01)

* **Bug Fix**: Correct wrapping of errors in authentication workflow.
* **Bug Fix**: Correctly recognize cache-wrapped instances of AnonymousCredentials at client construction.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.1 (2023-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2023-11-29)

* **Feature**: Expose Options() accessor on service clients.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.5 (2023-11-28.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.4 (2023-11-28)

* **Bug Fix**: Respect setting RetryMaxAttempts in functional options at client construction.

# v1.12.3 (2023-11-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.2 (2023-11-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.1 (2023-11-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.12.0 (2023-11-01)

* **Feature**: Adds support for configured endpoints via environment variables and the AWS shared configuration file.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2023-10-31)

* **Feature**: **BREAKING CHANGE**: Bump minimum go version to 1.19 per the revised [go version support policy](https://aws.amazon.com/blogs/developer/aws-sdk-for-go-aligns-with-go-release-policy-on-supported-runtimes/).
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2023-10-24)

* **Feature**: **BREAKFIX**: Correct nullability and default value representation of various input fields across a large number of services. Calling code that references one or more of the affected fields will need to update usage accordingly. See [2162](https://github.com/aws/aws-sdk-go-v2/issues/2162).

# v1.9.2 (2023-10-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2023-10-06)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2023-09-25)

* **Feature**: Adds support for sending WebRTC audio to Amazon Kineses Video Streams.

# v1.8.0 (2023-09-01)

* **Feature**: This release adds support for the Voice Analytics feature for customer-owned KVS streams as part of the Amazon Chime SDK call analytics.

# v1.7.0 (2023-08-31)

* **Feature**: This release adds support for feature Voice Enhancement for Call Recording as part of Amazon Chime SDK call analytics.

# v1.6.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2023-08-01)

* No change notes available for this release.

# v1.6.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-07-24)

* **Feature**: AWS Media Pipeline compositing enhancement and Media Insights Pipeline auto language identification.

# v1.4.5 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2023-06-15)

* No change notes available for this release.

# v1.4.3 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2023-05-04)

* No change notes available for this release.

# v1.4.1 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-04-20)

* **Feature**: This release adds support for specifying the recording file format in an S3 recording sink configuration.

# v1.3.2 (2023-04-10)

* No change notes available for this release.

# v1.3.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-03-23)

* **Feature**: This release adds Amazon Chime SDK call analytics. Call analytics include voice analytics, which provides speaker search and voice tone analysis. These capabilities can be used with Amazon Transcribe and Transcribe Call Analytics to generate machine-learning-powered insights from real-time audio.

# v1.2.6 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.2.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.2.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.1.9 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.8 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.7 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.6 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.5 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2022-08-18)

* **Feature**: The Amazon Chime SDK now supports live streaming of real-time video from the Amazon Chime SDK sessions to streaming platforms such as Amazon IVS and Amazon Elemental MediaLive. We have also added support for concatenation to create a single media capture file.

# v1.0.9 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.8 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.7 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.6 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.5 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.4 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2022-05-02)

* No change notes available for this release.

# v1.0.0 (2022-04-27)

* **Release**: New AWS service client module
* **Feature**: For Amazon Chime SDK meetings, the Amazon Chime Media Pipelines SDK allows builders to capture audio, video, and content share streams. You can also capture meeting events, live transcripts, and data messages. The pipelines save the artifacts to an Amazon S3 bucket that you designate.

