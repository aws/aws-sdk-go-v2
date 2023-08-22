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

