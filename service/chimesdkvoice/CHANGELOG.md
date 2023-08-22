# v1.8.3 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2023-08-09)

* **Feature**: Updating CreatePhoneNumberOrder, UpdatePhoneNumber and BatchUpdatePhoneNumbers APIs, adding phone number name

# v1.7.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2023-08-01)

* No change notes available for this release.

# v1.7.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.4 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2023-06-15)

* No change notes available for this release.

# v1.6.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2023-05-30)

* **Feature**: Added optional CallLeg field to StartSpeakerSearchTask API request

# v1.5.2 (2023-05-04)

* No change notes available for this release.

# v1.5.1 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2023-04-13)

* **Feature**: This release adds tagging support for Voice Connectors and SIP Media Applications

# v1.4.2 (2023-04-10)

* No change notes available for this release.

# v1.4.1 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-03-27)

* **Feature**: Documentation updates for Amazon Chime SDK Voice.

# v1.3.0 (2023-03-23)

* **Feature**: This release adds Amazon Chime SDK call analytics. Call analytics include voice analytics, which provides speaker search and voice tone analysis. These capabilities can be used with Amazon Transcribe and Transcribe Call Analytics to generate machine-learning-powered insights from real-time audio.

# v1.2.2 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2023-02-22)

* **Feature**: This release introduces support for Voice Connector media metrics in the Amazon Chime SDK Voice namespace
* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.1.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.1.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.0.3 (2022-12-23)

* No change notes available for this release.

# v1.0.2 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2022-11-18)

* **Release**: New AWS service client module
* **Feature**: Amazon Chime Voice Connector, Voice Connector Group and PSTN Audio Service APIs are now available in the Amazon Chime SDK Voice namespace. See https://docs.aws.amazon.com/chime-sdk/latest/dg/sdk-available-regions.html for more details.

