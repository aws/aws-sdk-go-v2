# v1.9.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.8.0 (2022-12-15)

* **Feature**: This release adds support for VMware vSphere tags, enabling customer to protect VMware virtual machines using tag-based policies for AWS tags mapped from vSphere tags. This release also adds support for customer-accessible gateway-hypervisor interaction log and upload bandwidth rate limit schedule.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.3 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2022-09-22)

* **Feature**: Changes include: new GetVirtualMachineApi to fetch a single user's VM, improving ListVirtualMachines to fetch filtered VMs as well as all VMs, and improving GetGatewayApi to now also return the gateway's MaintenanceStartTime.

# v1.6.12 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.11 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.10 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.9 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.8 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.7 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.6 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.5 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.4 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.3 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.2 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.1 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2022-06-01)

* **Feature**: Adds GetGateway and UpdateGatewaySoftwareNow API and adds hypervisor name to UpdateHypervisor API

# v1.5.6 (2022-05-26)

* No change notes available for this release.

# v1.5.5 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.4 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.0.0 (2021-12-02)

* **Release**: New AWS service client module
* **Dependency Update**: Updated to the latest SDK module versions

