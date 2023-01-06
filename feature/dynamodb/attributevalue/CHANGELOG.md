# v1.10.9 (2023-01-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.8 (2022-12-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.7 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.6 (2022-11-22)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.5 (2022-11-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.4 (2022-11-16)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.3 (2022-11-10)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.2 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2022-09-26)

* **Feature**: Adds a String method to UnixTime, so that when structs with this field get logged it prints a human readable time.

# v1.9.19 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.18 (2022-09-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.17 (2022-09-14)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.16 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.15 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.14 (2022-08-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.13 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.12 (2022-08-18)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.11 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.10 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.9 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.8 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.7 (2022-07-22)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.6 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.5 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.4 (2022-06-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.3 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.2 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.1 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.9.0 (2022-04-15)

* **Feature**: Support has been added for specifying a custom time format when encoding and decoding DynamoDB AttributeValues. Use `EncoderOptions.EncodeTime` to specify a custom time encoding function, and use `DecoderOptions.DecodeTime` for specifying how to handle the corresponding AttributeValues using the format. Thank you [Pablo Lopez](https://github.com/plopezlpz) for this contribution.

# v1.8.4 (2022-03-31)

* **Documentation**: Fixes documentation typos in Number type's helper methods

# v1.8.3 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2022-02-24)

* **Feature**: Fixes [#645](https://github.com/aws/aws-sdk-go-v2/issues/645), [#411](https://github.com/aws/aws-sdk-go-v2/issues/411) by adding support for (un)marshaling AttributeValue maps to Go maps key types of string, number, bool, and types implementing encoding.Text(un)Marshaler interface
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Bug Fix**: Fixes [#1569](https://github.com/aws/aws-sdk-go-v2/issues/1569) inconsistent serialization of Go struct field names
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2022-01-14)

* **Feature**: Adds new MarshalWithOptions and UnmarshalWithOptions helpers allowing Encoding and Decoding options to be specified when serializing AttributeValues. Addresses issue: https://github.com/aws/aws-sdk-go-v2/issues/1494
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.5.0 (2022-01-07)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.5 (2021-12-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.4 (2021-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.3 (2021-11-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.2 (2021-11-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.1 (2021-11-12)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.2.0 (2021-08-27)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Bug Fix**: Fix unmarshaler's decoding of AttributeValueMemberN into a type that is a string alias.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.5 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.4 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.3 (2021-07-15)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.2 (2021-06-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.1 (2021-05-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Dependency Update**: Updated to the latest SDK module versions

