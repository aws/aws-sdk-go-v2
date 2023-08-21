# v1.5.5 (2023-08-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2023-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2023-08-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2023-08-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2023-08-01)

* No change notes available for this release.

# v1.5.0 (2023-07-31)

* **Feature**: Adds support for smithy-modeled endpoint resolution. A new rules-based endpoint resolution will be added to the SDK which will supercede and deprecate existing endpoint resolution. Specifically, EndpointResolver will be deprecated while BaseEndpoint and EndpointResolverV2 will take its place. For more information, please see the Endpoints section in our Developer Guide.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2023-07-28)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2023-07-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2023-06-15)

* No change notes available for this release.

# v1.4.1 (2023-06-13)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2023-05-30)

* **Feature**: Campaigns now support selecting Timestream or S3 as the data destination, Signal catalogs now support "Deprecation" keyword released in VSS v2.1 and "Comment" keyword released in VSS v3.0

# v1.3.10 (2023-05-04)

* No change notes available for this release.

# v1.3.9 (2023-04-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.8 (2023-04-10)

* No change notes available for this release.

# v1.3.7 (2023-04-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.6 (2023-03-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.5 (2023-03-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.4 (2023-02-22)

* **Bug Fix**: Prevent nil pointer dereference when retrieving error codes.

# v1.3.3 (2023-02-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.2 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.3.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.2.1 (2022-12-30)

* **Documentation**: Update documentation - correct the epoch constant value of default value for expiryTime field in CreateCampaign request.

# v1.2.0 (2022-12-16)

* **Feature**: Updated error handling for empty resource names in "UpdateSignalCatalog" and "GetModelManifest" operations.

# v1.1.1 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2022-12-09)

* **Feature**: Deprecated assignedValue property for actuators and attributes.  Added a message to invalid nodes and invalid decoder manifest exceptions.

# v1.0.4 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.1 (2022-10-13)

* **Documentation**: Documentation update for AWS IoT FleetWise

# v1.0.0 (2022-09-26)

* **Release**: New AWS service client module
* **Feature**: General availability (GA) for AWS IoT Fleetwise. It adds AWS IoT Fleetwise to AWS SDK. For more information, see https://docs.aws.amazon.com/iot-fleetwise/latest/APIReference/Welcome.html.

