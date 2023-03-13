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

